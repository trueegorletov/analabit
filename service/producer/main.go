package main

import (
	"log"

	"analabit/service/producer/handler"
	"analabit/service/producer/proto"
	"go-micro.dev/v5"
)

func main() {
	// Create a new service
	service := micro.NewService(
		micro.Name("go.micro.service.producer"),
	)

	// Initialise the service
	service.Init()

	// Register handler
	if err := proto.RegisterProducerHandler(service.Server(), new(handler.Producer)); err != nil {
		log.Fatal(err)
	}

	// Run the service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
