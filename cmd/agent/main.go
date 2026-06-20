package main

import (
	"github.com/daryakovzhun/collect-metrics/internal/app"
	"log"
)

func main() {
	err := app.RunAgent()
	if err != nil {
		log.Fatal(err)
	}
}
