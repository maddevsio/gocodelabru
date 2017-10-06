# Шаг 3. Проектируем HTTP API
Нам нужны следующие HTTP API методы, исходя из описания задачи в задачи [Шаге 0](../step00/README.md)
## API Методы

1. POST /driver/ - добавить водителя
2. GET /driver/:id - получить информацию о водителе
3. DELETE /driver/:id - удалить водителя
4. GET /driver/:lat/:lon/nearest - получить ближайших водителей.

Для построения API мы будем использовать фреймворк [echo](http://echo.labstack.com)

Нам нужны будут пустые методы для этого, которые мы имплементируем позже.

```Go
func addDriver(c echo.Context) error {
	return nil
}
func getDriver(c echo.Context) error {
	return nil
}
func deleteDriver(c echo.Context) error {
	return nil
}
func nearestDrivers(c echo.Context) error {
	return nil
}
```

Получается примерно такой код.
```Go
package main

import (
	"log"

	"github.com/labstack/echo"
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
	return nil
}
func getDriver(c echo.Context) error {
	return nil
}
func deleteDriver(c echo.Context) error {
	return nil
}
func nearestDrivers(c echo.Context) error {
	return nil
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
		Driver  int `json:"driver"`
	}
	// Для возврата ближайших водителей
	NearestDriverResponse struct {
		Success bool              `json:"success"`
		Message string            `json:"message"`
		Drivers []int `json:"drivers"`
	}
)
```

## Поздравляю!
У нас есть основные структуры для получения/отправления данных и методы "заглушки". В [следующей](../step04/README.md) части мы реализуем их.
