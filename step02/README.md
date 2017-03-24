# Шаг 2. Hello world

Давайте напишем что-нибудь работающее. Например простое Hello world приложение. Откроем `main.go` в редакторе и напишем код.

```Go
package main

import "fmt"

func main() {
  fmt.Println("Hello world")
}
```

Теперь проект нужно собрать и запустить.
Есть несколько вариантов запустить проект.
``` 
$ go run main.go
Hello world
```

```
$ go build -o helloworld main.go
$ ./helloworld
Hello world
```

Учитывая то, что мы пишем веб приложение, давайте сделаем hello world в вебе.
```Go
package main

import (
        "fmt"
        "log"
        "net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "<h1>Hello world</h1>")
}
func main() {
        http.HandleFunc("/", hello)
        log.Fatal(http.ListenAndServe(":9911", nil))
}
```
В этом приложении мы сделали простой вебсервер, который при запуске будет слушать 9911 порт и на любой урл будет возвращать нам `Hello world`

Запустите и проверьте его работу сами

## Поздравляю!

У вас получилось что-то работающее. Продолжение в [следующей](../step03/README.md) части
