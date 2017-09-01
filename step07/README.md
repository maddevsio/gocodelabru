## Шаг 7. Строим сторадж

В этом шаге мы будем делать минимальное хранилище данных и решать проблему выдачи ближайшего водителя. При этом нам нужно

1. Придумать начальную архитектуру
2. Сделать его консистентным


Напомню, нам нужно хранить следующие данные
```Go
type (
  Location struct {
    Lat float64
    Lon float64
  }
  Driver struct {
    ID int 
    LastLocation Location
  }
)
```
Так и запишем их в `storage/storage.go`

Напомню, что нам нужно реализовать следующие фичи:

1. New() - для инциализации стораджа
2. Set(key, value) - для добавления или обновления элемента
3. Delete(key) - для удаления
4. Nearest(lat, lon) - для получения блищайших элементов
5. Get(key) - для получения водителя

Сделаем структуру для хранения водителей
```Go
type DriverStorage struct {
  drivers map[int]*Driver
}
```

Напишем методы для нашего хранилища. 
```Go
package storage

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
	return
}

// Delete removes driver from storage by key
func (d *DriverStorage) Delete(key int) error {
	return nil
}

// Get gets driver from storage and an error if nothing found
func (d *DriverStorage) Get(key int) (*Driver, error) {
	return nil, nil
}

// Nearest returns nearest drivers by locaion
func (d *DriverStorage) Nearest(lat, lon float64) ([]*Driver, error) {
	return nil, nil
}
```

Реализуем каждый из методов

## Set

```Go
// Set sets driver to storage by key
func (d *DriverStorage) Set(key int, driver *Driver) {
	d.drivers[key] = driver
}
```

## Delete

```Go
// Delete removes driver from storage by key
func (d *DriverStorage) Delete(key int) error {
	driver, ok := d.drivers[key]
	if !ok {
		return errors.New("Driver does not exist")
	}
	delete(d.drivers, key)
	return nil
}
```

## Get

```Go
// Get gets driver from storage and an error if nothing found
func (d *DriverStorage) Get(key int) (*Driver, error) {
	driver, ok := d.drivers[key]
	if !ok {
		return nil, errors.New("Driver does not exist")
	}
	return driver, nil
}
```
Но нужны еще также и тесты чтобы убедиться, что наш код работает.
Для того, чтобы писать меньше кода, мы поставим пакет assert
```
go get github.com/stretchr/testify/assert
```
После этого напишем тест
```Go
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
```

### Реализуем Nearest метод
В целом простая логика работы. Но нужно реализовать еще работу Nearest метода.
Как вариант можно сделать следующий алгоритм работы.

1. Мы проходим по всем водителям
2. Вычисляем расстояние до водителя
3. Если расстояние меньше заданого радиуса, то добавляем в массив результатов.


Предлагаю этот метод реализовать вам самим.
Скелет метода

```Go
// Nearest returns nearest drivers by locaion
func (d *DriverStorage) Nearest(radius, lat, lon float64) ([]*Driver, error) {
	return nil, nil
}
```
Как вычислить расстояние между двумя точками. Дистанция возвращается в метрах.

```Go
// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}
// Distance function returns the distance (in meters) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// point coordinates are supplied in degrees and converted into rad. in the func
//
// distance returned is METERS!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180
	r = 6378100 // Earth radius in METERS
	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)
	return 2 * r * math.Asin(math.Sqrt(h))
}
```

Тест для метода

```Go
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
```

```
go get github.com/stretchr/testify/assert
```

## Поздравляю! 
Вы сами реализовали метод, который ищет ближайших водителей. В [следующем](../step08/README.md) уроке мы будем разбираться, насколько наш метод эффективный и что можно сделать с этим.
