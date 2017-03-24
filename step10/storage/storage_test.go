package storage

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	s := New()
	driver := &Driver{
		ID: 1,
		LastLocation: Location{
			Lat: 1,
			Lon: 1,
		},
	}
	s.Set(driver.ID, driver)
	d, err := s.Get(driver.ID)
	assert.NoError(t, err)
	assert.Equal(t, d, driver)
	err = s.Delete(driver.ID)
	assert.NoError(t, err)
	d, err = s.Get(driver.ID)
	assert.Equal(t, err, errors.New("Driver does not exist"))
}

func TestNearest(t *testing.T) {
	s := New()
	s.Set(123, &Driver{
		ID: 123,
		LastLocation: Location{
			Lat: 1,
			Lon: 1,
		},
	})
	s.Set(666, &Driver{
		ID: 666,
		LastLocation: Location{
			Lat: 42.875799,
			Lon: 74.588279,
		},
	})
	drivers := s.Nearest(1000, 42.876420, 74.588332)
	assert.Equal(t, len(drivers), 1)
}
