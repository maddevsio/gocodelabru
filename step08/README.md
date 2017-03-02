# Шаг 8. Делаем HTTP API

##AddDreiver
```Go
func (a *DBAPI) addDriver(c echo.Context) error {
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

## GE
```Go
func (a *DBAPI) getDriver(c echo.Context) error {
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

## DELETE
```func (a *DBAPI) deleteDriver(c echo.Context) error {
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

## Nearest
```Go
func (a *DBAPI) nearestDrivers(c echo.Context) error {
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

## Поздравления
[следующий](../step09/README.md)
