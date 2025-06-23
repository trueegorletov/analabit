package main

import (
	"log"

	"analabit/service/aggregator/handler"
	"analabit/service/aggregator/proto"
	"go-micro.dev/v5"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.service.aggregator"),
	)

	service.Init()

	aggregatorHandler := new(handler.Aggregator)

	aggregatorHandler.StartSubscriber()

	if err := proto.RegisterAggregatorHandler(service.Server(), aggregatorHandler); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
