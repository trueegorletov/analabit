package handler

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"analabit/core"
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
	var conn *amqp.Connection
	var err error
	const maxRetries = 10
	for i := 1; i <= maxRetries; i++ {
		conn, err = amqp.Dial(rabbitURL)
		if err == nil {
			break
		}
		log.Printf("RabbitMQ connection attempt %d/%d failed: %v", i, maxRetries, err)
		time.Sleep(time.Duration(i) * time.Second) // simple linear backoff
	}
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ after %d attempts: %v", maxRetries, err)
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
			var notification map[string]interface{}
			if err := json.Unmarshal(d.Body, &notification); err != nil {
				log.Printf("Error unmarshalling notification: %v", err)
				continue
			}
			bucketName, ok := notification["bucket_name"].(string)
			if !ok {
				log.Println("Invalid notification: missing or invalid bucket_name")
				continue
			}

			payloadObjects, ok := notification["payload_objects"].([]interface{})
			if !ok {
				log.Println("Invalid notification: missing or invalid payload_objects")
				continue
			}

			objectNames := make([]string, len(payloadObjects))
			for i, v := range payloadObjects {
				objectNames[i] = v.(string)
			}

			log.Printf("Processing bucket: %s with %d objects", bucketName, len(objectNames))
			if err := a.processBucket(context.Background(), bucketName, objectNames); err != nil {
				log.Printf("Failed to process bucket %s: %v", bucketName, err)
			}
		}
		log.Println("RabbitMQ consumer stopped.")
		conn.Close()
		ch.Close()
	}()
}

func (a *Aggregator) processBucket(ctx context.Context, bucketName string, objectNames []string) error {
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

	pgConnStrings := cfg.PostgresConnStrings
	if pgConnStrings == "" {
		log.Println("POSTGRES_CONN_STRINGS is not set. Skipping database upload.")
		return nil
	}

	var allErrors error
	// Process each database connection string
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

		// Ensure schema is up-to-date once per database connection.
		if err := client.Schema.Create(ctx); err != nil {
			err = fmt.Errorf("failed to run schema migrations for conn string %q: %w", connStr, err)
			multierr.AppendInto(&allErrors, err)
			client.Close()
			continue
		}
		log.Printf("Schema check complete for %q.", connStr)

		// Create one Run per RabbitMQ notification before processing objects
		run, err := client.Run.Create().
			SetPayloadMeta(map[string]any{
				"bucket_name":    bucketName,
				"object_count":   len(objectNames),
				"object_names":   objectNames,
				"conn_string_id": connStr, // Consider hashing this for security
			}).
			Save(ctx)
		if err != nil {
			err = fmt.Errorf("failed to create run for bucket %s with conn string %q: %w", bucketName, connStr, err)
			multierr.AppendInto(&allErrors, err)
			client.Close()
			continue
		}

		log.Printf("Created run %d for bucket %s with %d objects", run.ID, bucketName, len(objectNames))

		// Process each object for the current database connection
		var payload core.UploadPayload // Reuse this variable
		for _, objectName := range objectNames {
			log.Printf("Processing object %s for run %d...", objectName, run.ID)

			// 1. Download GOB file
			obj, err := minioClient.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
			if err != nil {
				err = fmt.Errorf("failed to get object %s: %w", objectName, err)
				multierr.AppendInto(&allErrors, err)
				continue // Skip to the next object
			}

			// 2. Deserialize data
			if err := gob.NewDecoder(obj).Decode(&payload); err != nil {
				obj.Close()
				err = fmt.Errorf("failed to decode payload from object %s: %w", objectName, err)
				multierr.AppendInto(&allErrors, err)
				continue // Skip to the next object
			}
			obj.Close()

			// 3. Perform the unified upload with runID
			if err := upload.Primary(ctx, client, run.ID, &payload); err != nil {
				err = fmt.Errorf("failed to upload payload from object %s with conn string %q: %w", objectName, connStr, err)
				multierr.AppendInto(&allErrors, err)
			} else {
				log.Printf("Successfully uploaded payload for object %s in run %d", objectName, run.ID)
				// Build synthetic drained results for stage=0 TODO: move to producer service
				synthetic := make([]core.DrainedResultDTO, 0, len(payload.Calculations))
				for _, calc := range payload.Calculations {
					if len(calc.Admitted) == 0 {
						continue
					}
					last := calc.Admitted[len(calc.Admitted)-1]
					var passingScore, larp int
					for _, app := range payload.Applications {
						if app.HeadingCode == calc.HeadingCode && app.StudentID == last.ID {
							passingScore = app.Score
							larp = app.RatingPlace
							break
						}
					}
					synthetic = append(synthetic, core.DrainedResultDTO{
						HeadingCode:                calc.HeadingCode,
						DrainedPercent:             0,
						AvgPassingScore:            passingScore,
						AvgLastAdmittedRatingPlace: larp,
					})
				}
				// Combine synthetic and simulated drained results
				var drainedDTOs []core.DrainedResultDTO
				drainedDTOs = append(drainedDTOs, synthetic...)
				for _, dtos := range payload.Drained {
					drainedDTOs = append(drainedDTOs, dtos...)
				}
				// Upload drained results with runID
				if err := upload.DrainedResults(ctx, client, run.ID, drainedDTOs); err != nil {
					err = fmt.Errorf("failed to upload drained results from object %s with conn string %q: %w", objectName, connStr, err)
					multierr.AppendInto(&allErrors, err)
				} else {
					log.Printf("Successfully uploaded drained results for object %s in run %d", objectName, run.ID)
				}
			}
		}

		client.Close()
	}

	return allErrors
}
