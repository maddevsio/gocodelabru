## Шаг 12. Делаем хранилище консистентным. Внедряем LRU

Начнем с того, чтобы построить пространственный индекс, нам нужно знать границы точки. Это можно сделать, если мы сможем построить minimum bounding rectangle.
R-tree принимает в нашем случае Spatial объект, который должен имплементировать метод `Bounds()` который как раз таки и должен возвращать прямоугольник.

Мы в наше хранилище будем класть экземпляры `Driver` Поэтому имплементируем ему метод `Bounds()`

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
type Driver struct {
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
Этот метод у нас работает таким же образом. Мы и добавляем и обновляем данные.
Метод этот возвращает ошибку, потому что в R-tree нет метода Update. Зато есть `Delete()` и `Insert()`.
Поэтому перед тем как добавить элемент в БД, мы попробуем узнать есть ли он или нет. Если его нет, то мы проинициализируем LRU кеш ну и обновим все данные в итоге. Ошибку вернем только в случае, если мы не смогли удалить данные из нашего индекса.
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
	if !ok {
		d = driver
		cache, err := lru.New(s.lruSize)
		if err != nil {
			return errors.Wrap(err, "could not create LRU")
		}
		d.Locations = cache
	}
	d.LastLocation = driver.LastLocation
	d.Locations.Add(time.Now().UnixNano(), d.LastLocation)
	d.Expiration = driver.Expiration
	s.locations.Insert(d)
	s.drivers[driver.ID] = driver
	return nil
}
```
## Delete
Метод нужен для удаления данных. Метод вернет ошибку, если мы пытаемся удалить данные, которых нет в БД
```Go
// Delete deletes a driver from storage. Does nothing if the driver is not in the storage
func (s *DriverStorage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	d, ok := s.drivers[id]
	if !ok {
		return errors.New("does not exist")
	}
	deleted := s.locations.Delete(d)
	if deleted {
		delete(s.drivers, d.ID)
		return nil
	}
	return errors.New("could not remove item")
}
```

## Get
Для получения водителя по ключу. Вернет ошибку, если данных по ключу не существует.
```Go
// Get returns driver by key
func (s *DriverStorage) Get(id int) (*Driver, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.drivers[id]
	if !ok {
		return nil, errors.New("does not exist")
	}
	return d, nil
}
```

Протестируем все методы выше

```Go
func TestDriverStorage(t *testing.T) {
	s := New(10)
	s.Set(&Driver{
		ID: 123,
		LastLocation: Location{
			Lat: 1,
			Lon: 1,
		},
		Expiration: time.Now().Add(15).UnixNano(),
	})
	d, err := s.Get(123)
	assert.NoError(t, err)
	assert.Equal(t, d.ID, 123)
	err = s.Delete(123)
	assert.NoError(t, err)
	d, err = s.Get(123)
	assert.Equal(t, err.Error(), "does not exist")

}
```
## Nearest
В основу реализации ближайших водителей легли следующие мысли. 
Мы хотим возвращать ближайших водителей и привязываться к радиусу в этом моменте было бы бесполезно по следующим соображениям. Ближайший водитель может быть как за 100 метров, так и за пять километров. Поэтому для того, чтобы решить эту задачу более эффективно, нам всегда нужно получать N ближайших водителей.
```Go
// Nearest returns nearest drivers
func (s *DriverStorage) Nearest(point rtreego.Point, count int) []*Driver {
	s.mu.Lock()
	defer s.mu.Unlock()

	results := s.locations.NearestNeighbors(count, point)
	var drivers []*Driver
	for _, item := range results {
		if item == nil {
			continue
		}
		drivers = append(drivers, item.(*Driver))
	}
	return drivers

}
```
И тест на него, который покажет, что метод работает полностью. Потому что он вернет действительно ближайших водителей и равно столько, сколько нужно. Точки взяты где-то в центре, рядом с Бишкекпарком.
```Go
func TestNearest(t *testing.T) {
	s := New(10)
	s.Set(&Driver{
		ID: 123,
		LastLocation: Location{
			Lat: 42.875799,
			Lon: 74.588279,
		},
		Expiration: time.Now().Add(15).UnixNano(),
	})
	s.Set(&Driver{
		ID: 321,
		LastLocation: Location{
			Lat: 42.875508,
			Lon: 74.588107,
		},
		Expiration: time.Now().Add(15).UnixNano(),
	})
	s.Set(&Driver{
		ID: 666,
		LastLocation: Location{
			Lat: 42.876106,
			Lon: 74.588204,
		},
		Expiration: time.Now().Add(15).UnixNano(),
	})
	s.Set(&Driver{
		ID: 2319,
		LastLocation: Location{
			Lat: 42.874942,
			Lon: 74.585908,
		},
		Expiration: time.Now().Add(15).UnixNano(),
	})
	s.Set(&Driver{
		ID: 991,
		LastLocation: Location{
			Lat: 42.875744,
			Lon: 74.584503,
		},
		Expiration: time.Now().Add(15).UnixNano(),
	})

	drivers := s.Nearest(rtreego.Point{42.876420, 74.588332}, 3)
	assert.Equal(t, len(drivers), 3)
	assert.Equal(t, drivers[0].ID, 123)
	assert.Equal(t, drivers[1].ID, 321)
	assert.Equal(t, drivers[2].ID, 666)
}
```

## DeleteExpired
Так как мы решили сделать еще Expire механизм, то нам нужно удалять водителей, которые протухли. Реализация простая, мы просто проходим по всем элементам.
```Go
// DeleteExpired removes all expired items from storage
func (s *DriverStorage) DeleteExpired() {
	now := time.Now().UnixNano()
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range s.drivers {
		if v.Expiration > 0 && now > v.Expiration {
			deleted := s.locations.Delete(v)
			if deleted {
				delete(s.drivers, v.ID)
			}
		}
	}
}
```
Протестируем его.
```Go
func TestExpire(t *testing.T) {
	s := New(10)
	driver := &Driver{
		ID: 123,
		LastLocation: Location{
			Lat: 42.876420,
			Lon: 74.588332,
		},
		Expiration: time.Now().Add(-15).UnixNano(),
	}
	s.Set(driver)
	s.DeleteExpired()
	d, err := s.Get(123)
	assert.Error(t, err)
	assert.NotEqual(t, d, driver)

}
```

## Поздравляю!
На текущий момент вы сделали структуру данных, которая решает задачу хранения. В [следующей](../step07/README.md) части мы будем делать HTTP API к нему.
