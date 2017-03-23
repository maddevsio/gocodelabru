package storage

import "errors"

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

// DriverStorage is main storage for our project
type DriverStorage struct {
	drivers map[int]*Driver
}

// New creates new instance of DriverStorage
func New() *DriverStorage {
	d := &DriverStorage{}
	d.drivers = make(map[int]*Driver)
	return d
}

// Set sets driver to storage by key
func (d *DriverStorage) Set(key int, driver *Driver) {
	d.drivers[key] = driver
}

// Delete removes driver from storage by key
func (d *DriverStorage) Delete(key int) error {
	driver, ok := d.drivers[key]
	if !ok {
		return errors.New("Driver does not exist")
	}
	delete(d.drivers, key)
	return nil
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
func (d *DriverStorage) Nearest(lat, lon float64) ([]*Driver, error) {
	return nil, nil
}
