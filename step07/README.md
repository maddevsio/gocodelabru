# Шаг 7. Проектируем HTTP API
Для API нужен будет отдельный пакет. В него будем 
`api/api.go`
```Go
type DBAPI struct {
	database  *storage.DriverStorage
	waitGroup sync.WaitGroup
	echo      *echo.Echo
}
```

## API Методы

1. POST /driver/ - добавить водителя
2. GET /driver/:id - получить информацию о водителе
3. DELETE /driver/:id - удалить водителя
4. GET /driver/:lat/:lon/nearest - получить ближайших водителей.

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
```Go
func (a *DBAPI) WaitStop() {
	a.waitGroup.WaitStop()
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
Для запуска

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

```Go
type (
	DefaultResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	DriverResponse struct {
		Success bool            `json:"success"`
		Message string          `json:"message"`
		Driver  *storage.Driver `json:"driver"`
	}
	NearestDriverResponse struct {
		Success bool              `json:"success"`
		Message string            `json:"message"`
		Drivers []*storage.Driver `json:"drivers"`
	}
)
```

## Поздравления
[следующая](../step08/README.md)
