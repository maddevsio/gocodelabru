## Шаг 6. Имплементируем сторадж

Начнем с того, чтобы построить пространственный индекс, нам нужно знать границы точки. Это можно сделать, если мы сможем построить minimum bounding rectangle.
R-tree принимает в нашем случае Spatial объект, который должен имплементировать метод `Bounds()` который как раз таки и должен возвращать прямоугольник.

Мы в наше хранилище будем класть `Driver` Поэтому имплементируем ему метод `Bounds()`

## Bounds()
```Go
// Bounds method needs for correct working of rtree
// Lat - Y, Lon - X on coordinate system
func (d *Driver) Bounds() *rtreego.Rect {
	return rtreego.Point{d.LastLocation.Lat, d.LastLocation.Lon}.ToRect(0.01)
}
```
Еще нам нужно сделать Expire механизм. Модифицируем структуру `Driver` и добавим туда `Expiration`

```Go
type	Driver struct {
		ID           int
		LastLocation Location
		Expiration   int64
		Locations    *lru.LRU
}
```
Также сделаем метод `Expired()` для водителя, чтобы знать, нужно ли водителя удалять или нет
```
// Expired returns true if the item has expired.
func (d *Driver) Expired() bool {
	if d.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > d.Expiration
}

```

## Как создать сторадж или реализуем New
```Go
// New initializes now storage
func New(lruSize int) *DriverStorage {
	s := new(DriverStorage)
	s.drivers = make(map[int]*Driver)
	s.locations = rtreego.NewTree(2, 25, 50)
	s.mu = new(sync.RWMutex)
	s.lruSize = lruSize
	return s
}
```

## Set
```Go
// Set an Driver to the storage, replacing any existing item.
func (s *DriverStorage) Set(driver *Driver) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	d, ok := s.drivers[driver.ID]
	if ok {
		deleted := s.locations.Delete(d)
		if !deleted {
			return fmt.Errorf("failed to remove driver %d from r-tree", d.ID)
		}

	}
	d = d
	cache, err := lru.New(s.lruSize)
	if err != nil {
		return errors.Wrap(err, "could not create LRU")
	}
	d.Locations = cache
	d.LastLocation = driver.LastLocation
	d.Locations.Add(time.Now().UnixNano(), d.LastLocation)
	d.Expiration = driver.Expiration
	s.locations.Insert(d)
	s.drivers[driver.ID] = driver
	return nil
}
```
