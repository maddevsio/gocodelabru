# Шаг 1. Что нужно знать о тестировании и написании тестов в Go.

Тестирование нужно для слабаков, которые не могут написать с первой попытки работающий код. :trollface:
Поэтому в мире существует много инстументов для того, чтобы писать тесты. Go не исключение. В Go есть пакет `testing` функционала которого хватит всем.

В Go тестовые файлы обычно лежат в той же самой папке, что и обычные файлы. Их можно узнать по присутсвию `_test.go`  в имени файла. При этом компилятор поймет, что файлы `_test.go` не надо включать в билд, о при запуске инструмета `go test` эти файлы как раз таки и будут использоваться

Например посмотрим, как тестируют пакет `math` Go.

```Go
package math

import "testing"

func TestAverage(t *testing.T) {
  var v float64
  v = Average([]float64{1,2})
  if v != 1.5 {
    t.Error("Expected 1.5, got ", v)
  }
}
```

Запуск тестов происходит командой `go test`

```
$ cd /usr/local/go/src/math
$ go test
PASS
ok  	math	0.010s
```
Если хотите
А еще Go сообщество пропогандирует вместо копипаста в тестах использовать так называемые table tests. В этом случае у нас есть пары исходного значения и результата. А тесты прогоняем в цикле

```Go
package math

import "testing"

type testpair struct {
  values []float64
  average float64
}

var tests = []testpair{
  { []float64{1,2}, 1.5 },
  { []float64{1,1,1,1,1,1}, 1 },
  { []float64{-1,1}, 0 },
}

func TestAverage(t *testing.T) {
  for _, pair := range tests {
    v := Average(pair.values)
    if v != pair.average {
      t.Error(
        "For", pair.values,
        "expected", pair.average,
        "got", v,
      )
    }
  }
}
```
[Документация](http://godoc.org/testing) к пакету `testing`

А если вы запаритесь писать постоянно 
```Go
if smth != anoher {
   t.Error("Error")
}
```
То есть пакет [testify/assert](https://godoc.org/github.com/stretchr/testify/assert)


## Поздравляю!

Вы теперь знаете как тестировать и чем тестировать в Go. В нашем проекте мы будем писать тесты. Без них никуда. Продолжение в [следующей](../step02/README.md) части
