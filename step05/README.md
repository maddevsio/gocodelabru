## Шаг 5. Дробим на мелкие части
Вот когда у нас уже есть какой-то код, то можно вынести все из main.go. Сейчас первый кандидат - весь код, который относится к API.
Давайте его и вынесем.

У нас есть следующий код
```Go
package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type (
	Location struct {
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lon"`
	}
	Payload struct {
		Timestamp int64    `json:"timestamp"`
		DriverID  int      `json:"driver_id"`
		Location  Location `json:"location"`
	}
	// Структура для возврата ответа по умолчанию
	DefaultResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	// Для возврата ответа, когда мы запрашиваем водителя
	DriverResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Driver  int    `json:"driver"`
	}
	// Для возврата ближайших водителей
	NearestDriverResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Drivers []int  `json:"drivers"`
	}
)

func main() {
	e := echo.New()
	g := e.Group("/api")
	g.POST("/driver/", addDriver)
	g.GET("/driver/:id", getDriver)
	g.DELETE("/driver/:id", deleteDriver)
	g.GET("/driver/:lat/:lon/nearest", nearestDrivers)
	log.Fatal(e.Start(":9111"))
}
func addDriver(c echo.Context) error {
	p := &Payload{}
	if err := c.Bind(p); err != nil {
		return c.JSON(http.StatusUnsupportedMediaType, &DefaultResponse{
			Success: false,
			Message: "Set content-type application/json or check your payload data",
		})
	}
	return c.JSON(http.StatusOK, &DefaultResponse{
		Success: false,
		Message: "Added",
	})
}
func getDriver(c echo.Context) error {
	driverID := c.Param("id")
	id, err := strconv.Atoi(driverID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: "could not convert string to integer",
		})
	}
	return c.JSON(http.StatusOK, &DriverResponse{
		Success: true,
		Message: "found",
		Driver:  id,
	})
}

func deleteDriver(c echo.Context) error {
	driverID := c.Param("id")
	_, err := strconv.Atoi(driverID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: "could not convert string to integer",
		})
	}
	return c.JSON(http.StatusOK, &DefaultResponse{
		Success: true,
		Message: "removed",
	})
}
func nearestDrivers(c echo.Context) error {
	lat := c.Param("lat")
	lon := c.Param("lon")
	if lat == "" || lon == "" {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: "empty coordinates",
		})
	} 
	_, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: "failed convert float",
		})
	}
	_, err = strconv.ParseFloat(lon, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: "failed convert float",
		})
	}
	// TODO: Add nearest
	return c.JSON(http.StatusOK, &NearestDriverResponse{
		Success: false,
		Message: "found",
	})
}
```
К нему не хватает также и тестов. Плюс к этому, желательно как-то конфигурировать порт, на котором будет стартовать наш веб-сервер.

Создадим папку `api` и переместим весь код, который относится к API в `api/api.go`. Но для начала, нам нужна будет структура, которая понадобиться нам в будущем. Например для того, чтобы подключить в API нашу БД, которую мы сделаем позже. А сейчас там будем хранить копию echo, и адрес на котором будет слушать приложение
```Go
type API struct {
  echo *echo.Echo
  bindAddr string
}
```
На первых парах нам хватит такой простой структуры. При переносе нужно будет делать это методами класса. Ну и сделать метод для создания новой копии API.

## New или создаем API
В этом методе, мы инициализируем все наши зависимости и настраиваем роуты.
```Go
func New(bindAddr string) *API {
	a := &API{}
	a.echo = echo.New()
  	a.bindAddr = bindAddr
	g := a.echo.Group("/api")
	g.POST("/driver/", a.addDriver)
	g.GET("/driver/:id", a.getDriver)
	g.DELETE("/driver/:id", a.deleteDriver)
	g.GET("/driver/:lat/:lon/nearest", a.nearestDrivers)
	return a
}
```

## Start
Этот метод нужен будет для того, чтобы стартовать наше приложение
```Go
func (a *API) Start() error {
	return a.echo.Start(a.bindAddr)
}
```

## Остальные методы
Все остальные методы нужно привести к виду `func (a *API) addDriver(c echo.Context) error`

```Go
func (a *API) addDriver(c echo.Context) error {
	p := &Payload{}
	if err := c.Bind(p); err != nil {
		return c.JSON(http.StatusUnsupportedMediaType, &DefaultResponse{
			Success: false,
			Message: "Set content-type application/json or check your payload data",
		})
	}
	return c.JSON(http.StatusOK, &DefaultResponse{
		Success: false,
		Message: "Added",
	})
}
func (a *API) getDriver(c echo.Context) error {
	driverID := c.Param("id")
	id, err := strconv.Atoi(driverID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: "could not convert string to integer",
		})
	}
	return c.JSON(http.StatusOK, &DriverResponse{
		Success: true,
		Message: "found",
		Driver:  id,
	})
}

func (a *API) deleteDriver(c echo.Context) error {
	driverID := c.Param("id")
	_, err := strconv.Atoi(driverID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: "could not convert string to integer",
		})
	}
	return c.JSON(http.StatusOK, &DefaultResponse{
		Success: true,
		Message: "removed",
	})
}
func (a *API) nearestDrivers(c echo.Context) error {
	lat := c.Param("lat")
	lon := c.Param("lon")
	if lat == "" || lon == "" {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: "empty coordinates",
		})
	}
	_, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: "failed convert float",
		})
	}
	_, err = strconv.ParseFloat(lon, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: "failed convert float",
		})
	}
	// TODO: Add nearest
	return c.JSON(http.StatusOK, &NearestDriverResponse{
		Success: false,
		Message: "found",
	})
}
```

А все наши модели перенесем в `api/models.go`
```Go
package api

type (
	Location struct {
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lon"`
	}
	Payload struct {
		Timestamp int64    `json:"timestamp"`
		DriverID  int      `json:"driver_id"`
		Location  Location `json:"location"`
	}
	// Структура для возврата ответа по умолчанию
	DefaultResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	// Для возврата ответа, когда мы запрашиваем водителя
	DriverResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Driver  int    `json:"driver"`
	}
	// Для возврата ближайших водителей
	NearestDriverResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Drivers []int  `json:"drivers"`
	}
)

```

в main.go у нас остался только этот код
```Go
func main() {
	e := echo.New()
	g := e.Group("/api")
	g.POST("/driver/", addDriver)
	g.GET("/driver/:id", getDriver)
	g.DELETE("/driver/:id", deleteDriver)
	g.GET("/driver/:lat/:lon/nearest", nearestDrivers)
	log.Fatal(e.Start(":9111"))
}
```
Который нужно модифицировать примерно в это
```Go 
func main() {
	a := api.New(":9111")
	log.Fatal(a.Start())
}
```

## Поздравляю! 
Мы раздробили нашу програму на несколько частей. В [следующей](../step06/README.md) части мы поработаем с флагами, конфигурацией приложения и Makefile
