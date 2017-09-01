## Шаг 12. Делаем хранилище консистентным. Внедряем LRU

Для того, чтобы хранилище было консистентным нам хватит примитива `sync/Mutex`. Это обычный лок. С его помощью мы будем блокировать наше хранилище когда мы добаляем или удаляем элементы из него.

Плюс нам нужно еще сделать Expire механизм. Для этого мы модифицируем структуру `Driver` и добавим туда `Expiration` Ну и для хранения последних точек мы добавим к водителю LRU.
```Go
type Driver struct {
		ID           int
		LastLocation Location
		Expiration   int64
		Locations    *lru.LRU
}
```
Также сделаем метод `Expired()` для водителя, чтобы знать, нужно ли водителя удалять или нет
```Go
// Expired returns true if the item has expired.
func (d *Driver) Expired() bool {
	if d.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > d.Expiration
}

```
Расширяем хранилище
```Go
	DriverStorage struct {
		mu        *sync.RWMutex # для синхронизации
		drivers   map[int]*Driver
		locations *rtreego.Rtree
		lruSize   int # для того, чтобы инициализировать хранилище по каждому водителю
	}
```
## Новый New 
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
Поэтому перед тем как добавить элемент в БД, мы попробуем узнать есть ли он или нет. Если его нет, то мы проинициализируем LRU кеш ну и обновим все данные в итоге. Также откажемся от ключей. Они нам не нужны и мы будем только добавлять водителей. У нас же есть его ID для того, чтобы избегать дублирования
```Go
// Set an Driver to the storage, replacing any existing item.
func (s *DriverStorage) Set(driver *Driver)  {
	s.mu.Lock()
	defer s.mu.Unlock()

	d, ok := s.drivers[driver.ID]
	if !ok {
		d = driver
		cache, err := lru.New(s.lruSize)
		if err != nil {
			return errors.Wrap(err, "could not create LRU")
		}
		d.Locations = cache
		s.locations.Insert(d)
	}
	d.LastLocation = driver.LastLocation
	d.Locations.Add(time.Now().UnixNano(), d.LastLocation)
	d.Expiration = driver.Expiration

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
Мы сделали консистентное хранилище данных. В [следующей](../step13/README.md) части мы его внедрим к нам в API
