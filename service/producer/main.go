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

func main() {
	// Parse config once (the handler package's init already did this, but we parse again to allow overrides for tests)
	if err := env.Parse(&handler.Cfg); err != nil {
		log.Fatalf("failed to parse env config: %v", err)
	}

	cfg := handler.Cfg

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

	// Start self-query goroutine if enabled
	if cfg.SelfQueryPeriodMinutes != -1 {
		go func() {
			client := proto.NewProducerService("go.micro.service.producer", service.Client())

			period := time.Duration(cfg.SelfQueryPeriodMinutes) * time.Minute

			// Adaptive backoff parameters for the bootstrap phase.
			backoff := 10 * time.Second        // start with 10 s
			const maxBackoff = 1 * time.Minute // cap at 1 min while bootstrapping

			for {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
				_, err := client.Produce(ctx, &proto.ProduceRequest{})
				cancel()

				if err == nil {
					log.Println("Self-query initial Produce succeeded – entering periodic schedule")
					break
				}

				log.Printf("Self-query waiting for dependencies (next attempt in %s): %v", backoff, err)
				time.Sleep(backoff)

				// Exponentially increase backoff until maxBackoff or desired period, whichever is smaller.
				backoff *= 2
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
				if backoff > period {
					backoff = period
				}
			}

			ticker := time.NewTicker(period)
			defer ticker.Stop()
			for range ticker.C {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
				_, err := client.Produce(ctx, &proto.ProduceRequest{})
				if err != nil {
					log.Printf("Self-query Produce failed: %v", err)
				}
				cancel()
			}
		}()
	}

	// Run the service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
