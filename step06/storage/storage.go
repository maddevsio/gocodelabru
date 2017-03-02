package storage

import (
	"sync"

	"github.com/dhconnelly/rtreego"
	"github.com/maddevsio/gocodelabru/step06/storage/lru"
)

type (
	Location struct {
		Lat float64
		Lon float64
	}
	Driver struct {
		ID           int
		LastLocation Location
		Locations    *lru.LRU
	}
	DriverStorage struct {
		mu        *sync.RWMutex
		drivers   map[int]*Driver
		locations *rtreego.Rtree
	}
)
