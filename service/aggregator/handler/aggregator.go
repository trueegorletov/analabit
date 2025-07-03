package handler

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"analabit/core"
	"analabit/core/drainer"
	"analabit/core/ent"
	"analabit/core/upload"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/streadway/amqp"
	"go.uber.org/multierr"

	"github.com/caarlos0/env/v11"
	_ "github.com/lib/pq" // for postgres driver
)

type Aggregator struct{}

// StartSubscriber starts a long-running goroutine to listen for RabbitMQ messages.
func (a *Aggregator) StartSubscriber() {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse env config: %v", err)
	}

	rabbitURL := cfg.RabbitURL
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	q, err := ch.QueueDeclare(
		"buckets", // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	go func() {
		log.Println("RabbitMQ consumer started. Waiting for messages.")
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var notification map[string]string
			if err := json.Unmarshal(d.Body, &notification); err != nil {
				log.Printf("Error unmarshalling notification: %v", err)
				continue
			}
			bucketName, ok := notification["bucket_name"]
			if !ok {
				log.Println("Invalid notification: missing bucket_name")
				continue
			}
			log.Printf("Processing bucket: %s", bucketName)
			if err := a.processBucket(context.Background(), bucketName); err != nil {
				log.Printf("Failed to process bucket %s: %v", bucketName, err)
			}
		}
		log.Println("RabbitMQ consumer stopped.")
		conn.Close()
		ch.Close()
	}()
}

func (a *Aggregator) processBucket(ctx context.Context, bucketName string) error {
	var cfg config

	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("failed to parse env config: %w", err)
	}

	endpoint := cfg.MinioEndpoint
	accessKeyID := cfg.MinioAccessKey
	secretAccessKey := cfg.MinioSecretKey
	useSSL := cfg.MinioUseSSL

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize minio client: %w", err)
	}

	// 2. Download GOB files
	calcResultsObj, err := minioClient.GetObject(ctx, bucketName, "calculation_results.gob", minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to get calculation_results.gob: %w", err)
	}
	defer calcResultsObj.Close()

	drainedResultsObj, err := minioClient.GetObject(ctx, bucketName, "drained_results.gob", minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to get drained_results.gob: %w", err)
	}
	defer drainedResultsObj.Close()

	// 3. Deserialize data
	var calcResults []core.CalculationResult
	if err := gob.NewDecoder(calcResultsObj).Decode(&calcResults); err != nil {
		return fmt.Errorf("failed to decode calculation results: %w", err)
	}

	var drainedResults [][]drainer.DrainedResult
	if err := gob.NewDecoder(drainedResultsObj).Decode(&drainedResults); err != nil {
		return fmt.Errorf("failed to decode drained results: %w", err)
	}

	// 4. Connect to PostgreSQL and upload
	pgConnStrings := cfg.PostgresConnStrings
	if pgConnStrings == "" {
		log.Println("POSTGRES_CONN_STRINGS is not set. Skipping database upload.")
		return nil
	}

	var allErrors error
	for _, connStr := range strings.Split(pgConnStrings, ";") {
		if connStr == "" {
			continue
		}
		log.Printf("Connecting to database...")
		client, err := ent.Open("postgres", connStr)
		if err != nil {
			err = fmt.Errorf("failed to connect to postgres with conn string %q: %w", connStr, err)
			multierr.AppendInto(&allErrors, err)
			continue
		}

		log.Printf("Uploading data to database...")

		if err := upload.Primary(ctx, client, nil, calcResults); err != nil {
			err = fmt.Errorf("failed to upload primary results with conn string %q: %w", connStr, err)
			multierr.AppendInto(&allErrors, err)
		}

		for _, resultsPart := range drainedResults {
			if err := upload.DrainedResults(ctx, client, nil, resultsPart); err != nil {
				err = fmt.Errorf("failed to upload drained results with conn string %q: %w", connStr, err)
				multierr.AppendInto(&allErrors, err)
			} else {
				log.Printf("Successfully uploaded data for bucket %s to database", bucketName)
			}
		}
		client.Close()
	}

	return allErrors
}
