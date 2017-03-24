package storage

import (
	"errors"

	"github.com/dhconnelly/rtreego"
)

type (
	// Location used for storing driver's location
	Location struct {
		Lat float64
		Lon float64
	}
	// Driver model to store driver data
	Driver struct {
		ID           int
		LastLocation Location
	}
)

func (d *Driver) Bounds() *rtreego.Rect {
	return rtreego.Point{d.LastLocation.Lat, d.LastLocation.Lon}.ToRect(0.01)
}

// DriverStorage is main storage for our project
type DriverStorage struct {
	drivers   map[int]*Driver
	locations *rtreego.Rtree
}

// New creates new instance of DriverStorage
func New() *DriverStorage {
	d := &DriverStorage{}
	d.drivers = make(map[int]*Driver)
	d.locations = rtreego.NewTree(2, 25, 50)
	return d
}

// Set sets driver to storage by key
func (d *DriverStorage) Set(key int, driver *Driver) {
	_, ok := d.drivers[key]
	if !ok {
		d.locations.Insert(driver)
	}
	d.drivers[key] = driver
}

// Delete removes driver from storage by key
func (d *DriverStorage) Delete(key int) error {

	driver, ok := d.drivers[key]
	if !ok {
		return errors.New("driver does not exist")
	}
	if d.locations.Delete(driver) {
		delete(d.drivers, key)
		return nil
	}
	return errors.New("could not remove driver")
}

// Get gets driver from storage and an error if nothing found
func (d *DriverStorage) Get(key int) (*Driver, error) {
	driver, ok := d.drivers[key]
	if !ok {
		return nil, errors.New("Driver does not exist")
	}
	return driver, nil
}

// Nearest returns nearest drivers by locaion
func (d *DriverStorage) Nearest(count int, lat, lon float64) []*Driver {
	point := rtreego.Point{lat, lon}
	results := d.locations.NearestNeighbors(count, point)
	var drivers []*Driver
	for _, item := range results {
		if item == nil {
			continue
		}
		drivers = append(drivers, item.(*Driver))
	}
	return drivers
}
