# Шаг 10. Собираем все вместе и немного о Makefile
```Go
func main() {
	bindAddr := flag.String("bind_addr", ":8080", "Set bind address")
	size := flag.Int("lru_size", 20, "Set lru size per driver")
	flag.Parse()
	a := api.New(*bindAddr, *size)
	a.Start()
	a.WaitStop()
}

```

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
```
## Поздравления
[следующий](../step10/README.md)
