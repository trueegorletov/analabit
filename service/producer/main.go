package main

import (
	"context"
	"log"
	"time"

	"analabit/service/producer/handler"
	"analabit/service/producer/proto"

	"github.com/caarlos0/env/v11"
	micro "go-micro.dev/v5"
)

// startSelfQuery contains the logic for the self-triggering mechanism.
// It will be executed as a go-micro AfterStart hook.
func startSelfQuery(p *handler.Producer) func() error {
	return func() error {
		// Only start the goroutine if a period is configured.
		if handler.Cfg.SelfQueryPeriodMinutes == -1 {
			log.Println("Self-query disabled: period is -1.")
			return nil
		}

		go func() {
			period := time.Duration(handler.Cfg.SelfQueryPeriodMinutes) * time.Minute

			// The initial call no longer uses an RPC client. It calls the workflow directly.
			// This runs after the service is already registered, avoiding the old race condition.
			log.Println("Self-query hook started. Making initial Produce call directly.")
			err := p.Produce(context.Background(), &proto.ProduceRequest{}, &proto.ProduceResponse{})

			if err != nil {
				// If the first call fails, it's a significant issue.
				log.Printf("FATAL: Self-query initial Produce call failed, automatic runs will not start: %v", err)
				return // Do not start the ticker if the first run fails.
			}
			log.Println("Self-query initial Produce call succeeded. Entering periodic schedule.")

			// Start the ticker for subsequent periodic runs.
			ticker := time.NewTicker(period)
			defer ticker.Stop()

			for range ticker.C {
				log.Printf("Self-query ticker triggered. Making periodic Produce call directly.")
				err := p.Produce(context.Background(), &proto.ProduceRequest{}, &proto.ProduceResponse{})
				if err != nil {
					log.Printf("ERROR: Self-query periodic Produce call failed: %v", err)
				}
			}
		}()

		return nil
	}
}

func main() {
	// Parse config once
	if err := env.Parse(&handler.Cfg); err != nil {
		log.Fatalf("failed to parse env config: %v", err)
	}

	// Create handler instance to be shared
	producerHandler := new(handler.Producer)

	// Create a new service
	service := micro.NewService(
		micro.Name("go.micro.service.producer"),
	)

	// Initialise the service, which includes AfterStart hooks
	service.Init(
		// The self-query mechanism is now a lifecycle hook that calls the handler's method directly.
		micro.AfterStart(startSelfQuery(producerHandler)),
	)

	// Register handler
	if err := proto.RegisterProducerHandler(service.Server(), producerHandler); err != nil {
		log.Fatal(err)
	}

	// Run the service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
