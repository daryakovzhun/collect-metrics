package main

import (
	"context"
	"github.com/daryakovzhun/collect-metrics/internal/app"
	"log"
)

func main() {
	err := app.RunAgent(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
