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
	// Parse config
	cfg := handler.Cfg

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to parse env config: %v", err)
	}

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

			// --- Readiness phase ---
			readyCheckTicker := time.NewTicker(1 * time.Second)
			logProgressTicker := time.NewTicker(10 * time.Second)
			defer readyCheckTicker.Stop()
			defer logProgressTicker.Stop()

			waitingStart := time.Now()
			for {
				select {
				case <-readyCheckTicker.C:
					ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
					_, err := client.Produce(ctx, &proto.ProduceRequest{})
					cancel()
					if err == nil {
						log.Printf("Self-query readiness check succeeded after %s", time.Since(waitingStart).Round(time.Second))
						goto READY
					}
				case <-logProgressTicker.C:
					log.Printf("Waiting for dependencies to be ready before starting self-query (%s elapsed)…", time.Since(waitingStart).Round(time.Second))
				}
			}

		READY:
			// --- Periodic execution phase ---
			period := time.Duration(cfg.SelfQueryPeriodMinutes) * time.Minute
			for {
				time.Sleep(period)
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
