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

## Поздравляю! 
Мы сделали архитектуру для нашего хранилища. Разобрались как сделать консистентность данных. Реализовывать его будем в [следующем](../step06/README.md) уроке
