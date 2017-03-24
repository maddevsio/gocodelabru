# Шаг 13. Внедряем хранилище в API
Получается у нас такая структура. `sync.WaitGroup` нам нужна для того, чтобы синхронизировать несколько горутин, которые мы запустим позднее.

```Go
type API struct {
	database  *storage.DriverStorage
	waitGroup sync.WaitGroup
	echo      *echo.Echo
	bindAddr  string
}
```

## Новый New
```Go
func New(bindAddr string, lruSize int) *API {
	a := &API{}
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
func (a *API) WaitStop() {
	a.waitGroup.Wait()
}
```

### Remove expired

```Go
func (a *API) removeExpired() {
	for range time.Tick(1) {
		a.database.DeleteExpired()
	}
}
```
Этот метод - заблокирует наш основной поток. Вы можете попробовать его вызвать и после этого запустить например веб-сервер. Веб сервер в этом случае у нас не запусится. Вот тут к нам на выручку и приходит `sync.WaitGroup`

## Start
В этом методе мы просто запустим веб-сервер и удаление протухших водителей в двух горутинах. Заблокируем основной поток с помощью метода `WaitStop()`

```Go
func (a *API) Start() {
	a.waitGroup.Add(1)
	go func() {
		a.echo.Start(a.bindAddr)
		a.waitGroup.Done()
	}()
	a.waitGroup.Add(1)
	go a.removeExpired()
}

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
Ну и в `main.go` нужно добавить следующее
```Go
func main() {
	bindAddr := flag.String("bind_addr", ":8080", "Set bind address")
	size := flag.Int("lru_size", 20, "Set lru size per driver")
	flag.Parse()
	a := api.New(*bindAddr, *size)
	a.Start()
	a.WaitStop()
}
```
## Поздравляю!
Вы все сделали Переходите к  [следующему](../step14/README.md) шагу
