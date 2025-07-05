package main

import (
	"log"

	"analabit/service/aggregator/handler"

	micro "go-micro.dev/v5"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.service.aggregator"),
	)

	service.Init()

	aggregatorHandler := new(handler.Aggregator)

	aggregatorHandler.StartSubscriber()

	// No RPC endpoints defined in proto; skip registering handler to avoid reflection error.

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
