package handler

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/database"
	"github.com/trueegorletov/analabit/core/ent"
	"github.com/trueegorletov/analabit/core/upload"

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

	connStrings := cfg.GetAllConnStrings()
	if len(connStrings) == 0 {
		log.Println("No database connection strings configured. Skipping database upload.")
		return nil
	}

	var allErrors error
	// Process each database connection string (primary + replicas)
	for i, connStr := range connStrings {
		if connStr == "" {
			continue
		}

		dbType := "primary"
		if i > 0 {
			dbType = fmt.Sprintf("replica-%d", i)
		}

		log.Printf("Connecting to %s database...", dbType)
		client, err := ent.Open("postgres", connStr)
		if err != nil {
			err = fmt.Errorf("failed to connect to %s postgres with conn string %q: %w", dbType, connStr, err)
			multierr.AppendInto(&allErrors, err)
			continue
		}

		// Ensure schema is up-to-date once per database connection.
		if err := client.Schema.Create(ctx); err != nil {
			err = fmt.Errorf("failed to run schema migrations for %s database %q: %w", dbType, connStr, err)
			multierr.AppendInto(&allErrors, err)
			client.Close()
			continue
		}
		log.Printf("Schema check complete for %s database.", dbType)

		// Create one Run per RabbitMQ notification before processing objects
		run, err := client.Run.Create().
			SetPayloadMeta(map[string]any{
				"bucket_name":    bucketName,
				"object_count":   len(objectNames),
				"object_names":   objectNames,
				"conn_string_id": fmt.Sprintf("%s-%s", dbType, connStr), // More descriptive ID
			}).
			Save(ctx)
		if err != nil {
			err = fmt.Errorf("failed to create run for bucket %s with %s database %q: %w", bucketName, dbType, connStr, err)
			multierr.AppendInto(&allErrors, err)
			client.Close()
			continue
		}

		log.Printf("Created run %d for bucket %s with %d objects (%s database)", run.ID, bucketName, len(objectNames), dbType)

		// 1. Download and deserialize all payloads from the bucket
		payloads := make([]*core.UploadPayload, 0, len(objectNames))
		for _, objectName := range objectNames {
			log.Printf("Downloading object %s for run %d...", objectName, run.ID)
			obj, err := minioClient.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
			if err != nil {
				err = fmt.Errorf("failed to get object %s: %w", objectName, err)
				multierr.AppendInto(&allErrors, err)
				continue // Skip to the next object
			}

			var p core.UploadPayload
			if err := gob.NewDecoder(obj).Decode(&p); err != nil {
				obj.Close()
				err = fmt.Errorf("failed to decode payload from object %s: %w", objectName, err)
				multierr.AppendInto(&allErrors, err)
				continue // Skip to the next object
			}
			obj.Close()
			payloads = append(payloads, &p)
		}

		// If there were errors downloading/decoding, we might not want to proceed.
		// For now, we'll continue and let the logic handle potentially empty slices.

		// 2. Consolidate payloads into a single payload to resolve inconsistencies
		log.Printf("Consolidating %d payloads for run %d...", len(payloads), run.ID)
		consolidatedPayload := consolidatePayloads(payloads)

		// 3. Perform the unified upload with the consolidated payload
		if err := upload.Primary(ctx, client, run.ID, consolidatedPayload); err != nil {
			err = fmt.Errorf("failed to upload consolidated payload with %s database %q: %w", dbType, connStr, err)
			multierr.AppendInto(&allErrors, err)
		} else {
			log.Printf("Successfully uploaded consolidated payload for run %d (%s database)", run.ID, dbType)

			// Build synthetic drained results for stage=0 TODO: move to producer service
			synthetic := make([]core.DrainedResultDTO, 0, len(consolidatedPayload.Calculations))
			for _, calc := range consolidatedPayload.Calculations {
				if len(calc.Admitted) == 0 {
					continue
				}
				last := calc.Admitted[len(calc.Admitted)-1]
				var passingScore, larp int
				for _, app := range consolidatedPayload.Applications {
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
			for _, dtos := range consolidatedPayload.Drained {
				drainedDTOs = append(drainedDTOs, dtos...)
			}

			// Upload drained results with runID
			if err := upload.DrainedResults(ctx, client, run.ID, drainedDTOs); err != nil {
				err = fmt.Errorf("failed to upload drained results with %s database %q: %w", dbType, connStr, err)
				multierr.AppendInto(&allErrors, err)
			} else {
				log.Printf("Successfully uploaded drained results for run %d (%s database)", run.ID, dbType)
			}
		}

		// If no errors occurred for this run, refresh materialized views and mark it as finished
		if allErrors == nil {
			// Create database client wrapper for materialized view operations
			dbClient, err := database.NewClient(client)
			if err != nil {
				err = fmt.Errorf("failed to create database client for run %d: %w", run.ID, err)
				multierr.AppendInto(&allErrors, err)
			} else {
				// Run cleanup job first
				if cleanupErr := dbClient.PerformBackupAndCleanup(ctx, cfg.CleanupRetentionRuns, cfg.CleanupBackupDir); cleanupErr != nil {
					if database.IsBackupError(cleanupErr) {
						// Backup failed but cleanup succeeded - log warning and continue
						log.Printf("Warning: Backup failed for run %d but cleanup succeeded: %v", run.ID, cleanupErr)
					} else {
						// Cleanup itself failed - this is a serious error
						log.Printf("Error: Cleanup job failed for run %d: %v", run.ID, cleanupErr)
						multierr.AppendInto(&allErrors, cleanupErr)
					}
				}

				// Refresh materialized views after cleanup
				if err := upload.RefreshMaterializedViews(ctx, dbClient); err != nil {
					err = fmt.Errorf("failed to refresh materialized views for run %d: %w", run.ID, err)
					multierr.AppendInto(&allErrors, err)
				} else {
					// Mark run as finished only after successful view refresh
					_, updateErr := client.Run.UpdateOneID(run.ID).SetFinished(true).Save(ctx)
					if updateErr != nil {
						err = fmt.Errorf("failed to mark run %d as finished: %w", run.ID, updateErr)
						multierr.AppendInto(&allErrors, err)
					} else {
						log.Printf("Run %d completed successfully, cleanup performed, materialized views refreshed, and marked as finished (%s database)", run.ID, dbType)
					}
				}
			}
		}

		client.Close()
	}

	return allErrors
}

// consolidatePayloads merges multiple UploadPayloads into a single one,
// resolving inconsistencies like multiple original submissions.
func consolidatePayloads(payloads []*core.UploadPayload) *core.UploadPayload {
	if len(payloads) == 0 {
		return &core.UploadPayload{}
	}

	// 1. Combine all data and build lookup maps.
	merged := &core.UploadPayload{
		VarsityCode:  "consolidated",
		VarsityName:  "Consolidated",
		Headings:     make([]core.HeadingDTO, 0),
		Students:     make([]core.StudentDTO, 0),
		Applications: make([]core.ApplicationDTO, 0),
		Calculations: make([]core.CalculationResultDTO, 0),
		Drained:      make(map[int][]core.DrainedResultDTO),
	}
	studentOriginals := make(map[string]string) // studentID -> varsityCode where original is submitted
	headingToVarsity := make(map[string]string) // headingCode -> varsityCode
	seenHeadings := make(map[string]struct{})
	seenStudents := make(map[string]struct{})

	for _, p := range payloads {
		for _, h := range p.Headings {
			if _, seen := seenHeadings[h.Code]; !seen {
				merged.Headings = append(merged.Headings, h)
				seenHeadings[h.Code] = struct{}{}
			}
			headingToVarsity[h.Code] = p.VarsityCode
		}
		// The OriginalSubmitted flag on StudentDTO is now the source of truth.
		for _, s := range p.Students {
			if _, seen := seenStudents[s.ID]; !seen {
				merged.Students = append(merged.Students, s)
				seenStudents[s.ID] = struct{}{}
			}
			// Find the single true original submission for each student.
			if s.OriginalSubmitted {
				if existingVarsity, ok := studentOriginals[s.ID]; ok {
					log.Printf("Conflict: Student %s has original submission in both %s and %s. Using first one.", s.ID, existingVarsity, p.VarsityCode)
				} else {
					studentOriginals[s.ID] = p.VarsityCode
				}
			}
		}

		merged.Applications = append(merged.Applications, p.Applications...)
		merged.Calculations = append(merged.Calculations, p.Calculations...)
		for percent, drainedResults := range p.Drained {
			merged.Drained[percent] = append(merged.Drained[percent], drainedResults...)
		}
	}

	// 2. Correct the OriginalSubmitted flag on all applications for a student.
	// This ensures that if a student submitted an original to HSE, all their applications
	// to HSE have OriginalSubmitted = true, and all others have it as false.
	for i := range merged.Applications {
		app := &merged.Applications[i]
		originalVarsity, studentHasOriginal := studentOriginals[app.StudentID]
		appVarsity := headingToVarsity[app.HeadingCode]

		if studentHasOriginal {
			app.OriginalSubmitted = (originalVarsity == appVarsity)
		} else {
			app.OriginalSubmitted = false
		}
	}

	// 3. Filter calculations. A student can only be admitted to a university
	// where they submitted their original documents (if they submitted one at all).
	filteredCalculations := make([]core.CalculationResultDTO, 0, len(merged.Calculations))
	for _, calc := range merged.Calculations {
		calcVarsity := headingToVarsity[calc.HeadingCode]
		filteredAdmitted := make([]core.StudentDTO, 0, len(calc.Admitted))
		for _, student := range calc.Admitted {
			originalVarsity, hasOriginal := studentOriginals[student.ID]
			if !hasOriginal || (hasOriginal && originalVarsity == calcVarsity) {
				// Admit if:
				// 1. The student has not submitted an original anywhere.
				// 2. The student's original is at the same university as the calculation.
				filteredAdmitted = append(filteredAdmitted, student)
			}
		}
		calc.Admitted = filteredAdmitted
		filteredCalculations = append(filteredCalculations, calc)
	}
	merged.Calculations = filteredCalculations

	return merged
}
