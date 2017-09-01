# Шаг 8. Пишем первый бенчмарк и зачем он

Как вы можете догадаться, метод, который вы реализовали не совсем эффективный.
А для того, чтобы убедиться в этом, напишем бенчмарк в `storage/storage_test.go`
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
		s.Nearest(1000, 123, 123)
	}
}
```
И проверим работу на 100, 1000, 10000 элементах в хранилище.
```
cd storage
go test -bench=.
```
Для 100 элементов
```
BenchmarkNearest-4         50000             24002 ns/op
PASS
ok      github.com/maddevsio/gocodelabru/step09/storage 1.460s
```
Для 1000 элементов
```
BenchmarkNearest-4          5000            272552 ns/op
PASS
ok      github.com/maddevsio/gocodelabru/step09/storage 1.402s
```
Для 10000 элементов
```
BenchmarkNearest-4           500           2799431 ns/op
PASS
ok      github.com/maddevsio/gocodelabru/step09/storage 1.714s
```

Можно сделать вывод, что чем больше элеменов в хранилище, тем дольше мы ищем ближайших водителей.

## Поздравляю!
Вы теперь знаете как писать бенчмарки В [следующем](../step09/README.md) шаге мы будем оптимизировать работу метода выдачи ближайших водителей
