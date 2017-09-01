# Шаг 9. Оптимизируем Nearest. Используем R-tree
В прошлом шаге мы написали свой метод, который возвращает ближайщих водителей и он работает медлено. В этой части мы будем его оптимизировать. Чтобы оптимизировать его - используем R-tree

## R-tree
![](./400px-R-tree.svg.png)

R-tree выглядит как показано на картинке. Это древовидная структура данных. Она хороша для понимания, если вы знакомы с B-деревом. R-tree нужно для индексации пространственных данных(координаты, города на карте). Также она решает нашу проблему. У нее можно спросить "Дай мне 10 ближайших водителей рядом со мной". Она идеально подходит для нас.
[Подробнее](https://ru.wikipedia.org/wiki/R-%D0%B4%D0%B5%D1%80%D0%B5%D0%B2%D0%BE_(%D1%81%D1%82%D1%80%D1%83%D0%BA%D1%82%D1%83%D1%80%D0%B0_%D0%B4%D0%B0%D0%BD%D0%BD%D1%8B%D1%85))

Мы не будем ее делать сами, потому что уже есть готовая реализация и мы возьмем ее [отсюда](https://github.com/dhconnelly/rtreego).


Установим ее
```
go get github.com/dhconnelly/rtreego
```

Внедрим в наше хранилище
```Go
// DriverStorage is main storage for our project
type DriverStorage struct {
	drivers   map[int]*Driver
	locations *rtreego.Rtree
}

```
Ну и теперь нам нужно адаптировать все наши методы, чтобы они работали с Rtree

Начнем с того, чтобы построить пространственный индекс, нам нужно знать границы точки. Это можно сделать, если мы сможем построить минимальный ограничивающий прямоугольник. Для чего он и что это такое, можно прочитать [тут](https://en.wikipedia.org/wiki/Minimum_bounding_rectangle)
R-tree принимает в нашем случае Spatial интерфейс, который должен имплементировать метод `Bounds()` который как раз таки и должен возвращать прямоугольник.

Мы в наше хранилище будем класть экземпляры `Driver` Поэтому имплементируем ему метод `Bounds()`

## Bounds()
```Go
// Bounds method needs for correct working of rtree
// Lat - Y, Lon - X on coordinate system
func (d *Driver) Bounds() *rtreego.Rect {
	return rtreego.Point{d.LastLocation.Lat, d.LastLocation.Lon}.ToRect(0.01)
}
```

## New
```Go
// New creates new instance of DriverStorage
func New() *DriverStorage {
	d := &DriverStorage{}
	d.drivers = make(map[int]*Driver)
	d.locations = rtreego.NewTree(2, 25, 50)

	return d
}
```

## Set
```Go
// Set sets driver to storage by key
func (d *DriverStorage) Set(key int, driver *Driver) {
	_, ok := d.drivers[key]
	if !ok {
		d.locations.Insert(driver)
	}
	d.drivers[key] = driver
}

```

## Delete
```Go
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


```

Теперь нам нужно адаптировать `Nearest()` метод. Судя по [документации](https://godoc.org/github.com/dhconnelly/rtreego) есть метод `NearestNeighbors()` которому нужно передать количество элементов, которые нужно вернуть ближайшими. У этого метода также нет радиуса.
Поэтому метод `Nearest()` будет выглядеть следующим образом
```Go
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
```
И тест на него исправится
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
	drivers := s.Nearest(1, 42.876420, 74.588332)
	assert.Equal(t, len(drivers), 1)
}
```

Адаптируем и наш бенчмарк
```Go
func BenchmarkNearest(b *testing.B) {
	s := New()
	for i := 0; i < 100; i++ {
		s.Set(i, &Driver{
			ID: i,
			LastLocation: Location{
				Lat: float64(i),
				Lon: float64(i),
			},
		})
	}
	for i := 0; i < b.N; i++ {
		s.Nearest(10, 123, 123)
	}
}
```
И проверим его для 100, 1000 и 10000 элементов
100
```
BenchmarkNearest-4        200000              6649 ns/op
PASS
ok      github.com/maddevsio/gocodelabru/step10/storage 1.418s
```
1000
```
 go test -bench=.
BenchmarkNearest-4         20000             76832 ns/op
PASS
ok      github.com/maddevsio/gocodelabru/step10/storage 1.745s
```
10000
```
BenchmarkNearest-4          5000            210245 ns/op
PASS
ok      github.com/maddevsio/gocodelabru/step10/storage 9.951s
```

Ну и как видим, работать стало все быстрее.

## Поздравляю!
Вы узнали что такое R-tree и внедрили в проект. В [следующей](../step10/README.md) части мы начнем решать задачу с хранением нескольких последних координат.
