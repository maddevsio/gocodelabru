# Шаг 13. Внедряем хранилище в API
# Шаг 7. Проектируем HTTP API
Нам нужны следующие HTTP API методы, исходя из описания задачи в задачи [Шаге 0](../step00/README.md)
## API Методы

1. POST /driver/ - добавить водителя
2. GET /driver/:id - получить информацию о водителе
3. DELETE /driver/:id - удалить водителя
4. GET /driver/:lat/:lon/nearest - получить ближайших водителей.

Для построения API мы будем использовать фреймворк [echo](http://echo.labstack.com)

### Идеи, которые реализуем в api пакете.
Нам нужно будет запустить две блокирующие операции паралельно. Для этого нам нужны горутины. А для того, чтобы наш основной поток не закончился раньше времени нам поможет `sync.WaitGroup`. А еще где-то нужно хранить копию нашей БД.
Вырисовывается такая структурка для `api/api.go`
```Go
type DBAPI struct {
	database  *storage.DriverStorage
	waitGroup sync.WaitGroup
	echo      *echo.Echo
	bindAddr  string
}
```

Нам нужны будут пустые методы для этого, которые мы имплементируем позже.

```
func (a *DBAPI) addDriver(c echo.Context) error {
	return nil
}
func (a *DBAPI) getDriver(c echo.Context) error {
	return nil
}
func (a *DBAPI) deleteDriver(c echo.Context) error {
	return nil
}
func (a *DBAPI) nearestDrivers(c echo.Context) error {
	return nil
}
```

## New или создаем API
В этом методе, мы инициализируем все наши зависимости и настраиваем роуты.
```Go
func New(bindAddr string, lruSize int) *DBAPI {
	a := &DBAPI{}
	a.database = storage.New(lruSize)
	a.echo = echo.New()
	g := a.echo.Group("/api")
	g.POST("/driver/", a.addDriver)
	g.GET("/driver/:id", a.getDriver)
	g.DELETE("/driver/:id", a.deleteDriver)
	g.GET("/driver/:lat/:lon/nearest", a.nearestDrivers)
	return a
}
```

## WaitStop
В качестве обертки над приватной вейтгруппой
```Go
func (a *DBAPI) WaitStop() {
	a.waitGroup.Wait()
}
```

### Remove expired

```Go
func (a *DBAPI) removeExpired() {
	for range time.Tick(1) {
		a.database.DeleteExpired()
	}
}
```

## Start 
В этом методе мы просто запустим веб-сервер и удаление протухших водителей в двух горутинах. Заблокируем основной поток с помощью метода `WaitStop()`

```Go
func (a *DBAPI) Start() {
	a.waitGroup.Add(1)
	go func() {
		a.echo.Start(a.bindAddr)
		a.waitGroup.Done()
	}()
	a.waitGroup.Add(1)
	go a.removeExpired()
}

```

### Запросы и ответы
Нам понадобятся следующие структуры для получения запросов
```Go
type (
    Location struct {
        Latitude float64 `json:"lat"`
        Longitude float64 `json:"lon"`
    }

    Payload struct {
      Timestamp int64 `json:"timestamp"`
      DriverID int `json:"driver_id"`
      Location Location `json:"location"`
    }
)
```
Для возврата ответов используем следующее
```Go
type (
	// Структура для возврата ответа по умолчанию
	DefaultResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	// Для возврата ответа, когда мы запрашиваем водителя
	DriverResponse struct {
		Success bool            `json:"success"`
		Message string          `json:"message"`
		Driver  *storage.Driver `json:"driver"`
	}
	// Для возврата ближайших водителей
	NearestDriverResponse struct {
		Success bool              `json:"success"`
		Message string            `json:"message"`
		Drivers []*storage.Driver `json:"drivers"`
	}
)
```


## addDriver
```Go
func (a *API) addDriver(c echo.Context) error {
	p := &Payload{}
	if err := c.Bind(p); err != nil {
		return c.JSON(http.StatusUnsupportedMediaType, &DefaultResponse{
			Success: false,
			Message: "Set content-type application/json or check your payload data",
		})
	}
	driver := &storage.Driver{}
	driver.ID = p.DriverID
	driver.LastLocation = storage.Location{
		Lat: p.Location.Latitude,
		Lon: p.Location.Longitude,
	}
	if err := a.database.Set(driver); err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, &DefaultResponse{
		Success: false,
		Message: "Added",
	})
}

```

## getDriver
```Go
func (a *API) getDriver(c echo.Context) error {
	driverID := c.Param("id")
	id, err := strconv.Atoi(driverID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: "could not convert string to integer",
		})
	}
	d, err := a.database.Get(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, &DriverResponse{
		Success: true,
		Message: "found",
		Driver:  d,
	})
}
```

## deleteDriver
```Go
func (a *API) deleteDriver(c echo.Context) error {
	driverID := c.Param("id")
	id, err := strconv.Atoi(driverID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: "could not convert string to integer",
		})
	}
	if err := a.database.Delete(id); err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, &DefaultResponse{
		Success: true,
		Message: "removed",
	})
}
```

## NearestDrivers
```Go
func (a *API) nearestDrivers(c echo.Context) error {
	lat := c.Param("lat")
	lon := c.Param("lon")
	if lat == "" || lon == "" {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: "empty coordinates",
		})
	}
	lt, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: "failed convert float",
		})
	}
	ln, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: "failed convert float",
		})
	}
	drivers := a.database.Nearest(rtreego.Point{lt, ln}, 10)
	return c.JSON(http.StatusOK, &NearestDriverResponse{
		Success: false,
		Message: "found",
		Drivers: drivers,
	})
}

```
Ну, правда на этом этапе. чтобы не плодить не нужные структуры, поменяем в `storage/storage.go` файлы, добавив `json` теги
```Go
	Location struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}
	Driver struct {
		ID           int      `json:"id"`
		LastLocation Location `json:"location"`
		Expiration   int64    `json:"-"`
		Locations    *lru.LRU `json:"-"`
	}

```
## Поздравляю!
Вы сделали апи. В [следующем](../step14/README.md) шаге мы напишем тесты на нее
