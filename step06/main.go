package main

import (
	"log"

	"github.com/maddevsio/gocodelabru/step06/api"
)

func main() {
	a := api.New(":9111")
	log.Fatal(a.Start())
}
