package api

import (
	"sync"
	"time"

	"github.com/labstack/echo"
	"github.com/maddevsio/gocodelabru/step07/storage"
)

type DBAPI struct {
	database  *storage.DriverStorage
	waitGroup sync.WaitGroup
	echo      *echo.Echo
	bindAddr  string
}

func New(bindAddr string, lruSize int) *DBAPI {
	a := &DBAPI{}
	a.database = storage.New(lruSize)
	a.echo = echo.New()
	g := a.echo.Group("/api")
	g.POST("/driver/", a.addDriver)
	g.GET("/driver/:id", a.getDriver)
	g.DELETE("/driver/:id", a.deleteDriver)
	g.GET("/driver/:lat/:lon/nearest", a.nearestDrivers)
	a.bindAddr = bindAddr
	return a
}

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

func (a *DBAPI) Start() {
	a.waitGroup.Add(1)
	go func() {
		a.echo.Start(a.bindAddr)
		a.waitGroup.Done()
	}()
	a.waitGroup.Add(1)
	go a.removeExpired()
}

func (a *DBAPI) removeExpired() {
	for range time.Tick(1) {
		a.database.DeleteExpired()
	}
}

func (a *DBAPI) WaitStop() {
	a.waitGroup.WaitStop()
}
