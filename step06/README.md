# Шаг 6. Собираем все вместе и немного о Makefile
Пришла пора для того, чтобы конфигурировать наше приложение
Для конфигурации мы воспользуемся пакетом [flag](https://godoc.org/flag)
```Go
import "flag"

func main() {
	bindAddr := flag.String("bind_addr", ":8080", "Set bind address")
	flag.Parse()
	a := api.New(*bindAddr)
	log.Fatal(a.Start())
}
```
В GO есть несколько правил хорошего тона:

1. В твоем пакете должны быть тесты
2. Твои библиотеки не должны писать логи
3. Сообщения об ошибках написаны в lowercase
4. Твой код должен быть отформатирован

Для того, чтобы код всегда был отформатирован, помимо триггеров на сохранение, можно запустить команду `go fmt ./...` и он отформатирует все файлы.
Для примера поделюсь с вами простым `Makefile` 
```make
TARGET=codelab

all: fmt clean build

clean:
	rm -rf $(TARGET)

depends:
	go get -u -v

build:
	go build -v -o $(TARGET) main.go

fmt:
	go fmt ./...

test:
	go test -v ./...
run:
	go run main.go
```


## Поздравляю!
Выполнив make вы отформатируете весь код и пересоберете проект, а `make run` запустит вам проект. В [следующем](../step07/README.md) шаге мы будем писать тесты на наше API и узнаем о coverage
