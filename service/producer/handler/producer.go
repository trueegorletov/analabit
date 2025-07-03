package handler

import (
	"analabit/core/source"
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/streadway/amqp"
	"go-micro.dev/v5/errors"

	"analabit/core"
	"analabit/core/drainer"
	"analabit/core/registry"
	"analabit/service/producer/proto"

	"go.uber.org/multierr"
)

type Producer struct{}

// concurrency guard to ensure only one Produce execution at a time per service instance
var (
	produceRunning int32 // 0 = not running, 1 = running (atomic flag)
)

// produceLock is a channel-based semaphore to ensure only one produce workflow runs at a time.
var produceLock = make(chan struct{}, 1)

// Produce is the public RPC endpoint. It's a thin wrapper that acquires a lock
// and calls the main production workflow.
func (p *Producer) Produce(ctx context.Context, req *proto.ProduceRequest, rsp *proto.ProduceResponse) error {
	log.Println("Received Produce request")

	select {
	case produceLock <- struct{}{}:
		log.Println("Acquired produce lock, starting workflow.")
		defer func() {
			<-produceLock
			log.Println("Released produce lock.")
		}()
	default:
		slog.Warn("Produce request ignored – another workflow is already in progress.")
		return errors.InternalServerError("producer.produce.busy", "another Produce request is already being processed")
	}

	// Run the actual workflow
	return p.runProduceWorkflow(ctx, req, rsp)
}

// runProduceWorkflow contains the core logic for crawling, calculating, and uploading results.
// This can be called either by the public RPC endpoint or an internal ticker.
func (p *Producer) runProduceWorkflow(ctx context.Context, req *proto.ProduceRequest, rsp *proto.ProduceResponse) error {
	varsitiesList := req.GetVarsitiesList()
	if len(varsitiesList) == 0 {
		varsitiesList = []string{"all"}
	}
	varsitiesExcluded := req.GetVarsitiesExcluded()

	cacheTTL := req.GetCacheTtlMinutes()
	if cacheTTL == 0 {
		cacheTTL = int32(Cfg.CacheTTLMinutes)
	}
	drainStages := req.GetDrainStages()
	if len(drainStages) == 0 {
		drainStages = make([]int32, len(Cfg.DrainStages))
		for i, v := range Cfg.DrainStages {
			drainStages[i] = int32(v)
		}
	}
	drainIterations := req.GetDrainIterations()
	if drainIterations == 0 {
		drainIterations = int32(Cfg.DrainIterations)
	}

	params := registry.CrawlOptions{
		VarsitiesList:    varsitiesList,
		VarsitiesExclude: varsitiesExcluded,
		CacheDir:         Cfg.CacheDir,
		CacheTTLMinutes:  int(cacheTTL),
		DrainStages:      toIntSlice(drainStages),
		DrainIterations:  int(drainIterations),
	}

	slog.Info("Producer configured", "varsitiesList", params.VarsitiesList, "varsitiesExclude", params.VarsitiesExclude, "cacheTTL", params.CacheTTLMinutes, "drainStages", params.DrainStages, "drainIterations", params.DrainIterations)

	slog.Info("Starting crawl and cache phase")
	result, err := registry.CrawlWithOptions(registry.AllDefinitions, params)
	if err != nil {
		log.Printf("failed to crawl or cache: %v", err)
		return err
	}
	varsities := result.LoadedVarsities
	slog.Info("Crawl completed", "varsitiesLoaded", len(varsities))

	loadedVarsityCodes := make([]string, len(varsities))
	for i, v := range varsities {
		loadedVarsityCodes[i] = v.Code
	}
	slog.Info("Loaded varsity codes", "codes", loadedVarsityCodes)

	slog.Info("Starting primary calculations")
	primaryResults := make(map[string][]core.CalculationResult)
	for _, v := range varsities {
		clonedVarsity := v.Clone()
		results := clonedVarsity.VarsityCalculator.CalculateAdmissions()
		primaryResults[v.Code] = results
	}
	slog.Info("Calculations completed – starting drain simulations")
	drainedResults := make(map[string]map[int][]drainer.DrainedResult)
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, v := range varsities {
		drainedResults[v.Code] = make(map[int][]drainer.DrainedResult)
	}
	for _, v := range varsities {
		for _, stage := range params.DrainStages {
			wg.Add(1)
			go func(v *source.Varsity, stage int) {
				defer wg.Done()
				drainerInstance := drainer.New(v, stage)
				drainedResultSlice := drainerInstance.Run(params.DrainIterations)
				mu.Lock()
				drainedResults[v.Code][stage] = drainedResultSlice
				mu.Unlock()
			}(v, stage)
		}
	}
	wg.Wait()
	slog.Info("Drain simulations completed")

	for varsityCode, stages := range drainedResults {
		stageSummary := make(map[int]int)
		for stage, results := range stages {
			stageSummary[stage] = len(results)
		}
		slog.Info("Drained results summary for varsity", "varsityCode", varsityCode, "stages", stageSummary)
	}

	var calcResultsBuf bytes.Buffer
	if err := gob.NewEncoder(&calcResultsBuf).Encode(primaryResults); err != nil {
		log.Printf("failed to encode calculation results: %v", err)
		return errors.InternalServerError("producer.produce.gob", "failed to encode calculation results: %v", err)
	}

	var drainedResultsBuf bytes.Buffer
	if err := gob.NewEncoder(&drainedResultsBuf).Encode(drainedResults); err != nil {
		log.Printf("failed to encode drained results: %v", err)
		return errors.InternalServerError("producer.produce.gob", "failed to encode drained results: %v", err)
	}

	slog.Info("Encoding results and uploading to object storage")
	minioClient, err := minio.New(Cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(Cfg.MinioAccessKey, Cfg.MinioSecretKey, ""),
		Secure: Cfg.MinioUseSSL,
	})
	if err != nil {
		log.Printf("failed to initialize minio client: %v", err)
		return errors.InternalServerError("producer.produce.minio", "failed to initialize minio client: %v", err)
	}

	// Use bucket name from Cfg
	bucketName := Cfg.MinioBucketName
	if bucketName == "" {
		return errors.InternalServerError("producer.produce.minio", "Minio bucket name is not set in Cfg")
	}

	// Generate unique object names for this calculation
	objectUUID := uuid.New().String()
	calcResultsObjectName := "calculation_results_" + objectUUID + ".gob"
	drainedResultsObjectName := "drained_results_" + objectUUID + ".gob"

	// --- Retry logic for MinIO operations ---
	timeouts := []time.Duration{3 * time.Minute, 5 * time.Minute, 8 * time.Minute}
	var lastErr error

	// 1. Ensure bucket exists
	var bucketOpSuccess bool
	for i, timeout := range timeouts {
		log.Printf("Attempt %d/%d: Checking/creating MinIO bucket '%s' (timeout: %v)", i+1, len(timeouts), bucketName, timeout)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		exists, bucketErr := minioClient.BucketExists(ctx, bucketName)
		if bucketErr != nil {
			lastErr = bucketErr
			cancel()
			continue
		}
		if !exists {
			bucketErr = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
			if bucketErr != nil {
				lastErr = bucketErr
				cancel()
				continue
			}
			log.Printf("Successfully created bucket %s", bucketName)
		}

		bucketOpSuccess = true
		cancel()
		break
	}

	if !bucketOpSuccess {
		log.Printf("failed to ensure bucket exists after %d attempts: %v", len(timeouts), lastErr)
		return errors.InternalServerError("producer.produce.minio", "failed to ensure bucket exists after multiple retries: %v", lastErr)
	}

	// 2. Upload objects
	var uploadSuccess bool
	for i, timeout := range timeouts {
		log.Printf("Attempt %d/%d: Uploading results to MinIO (timeout: %v)", i+1, len(timeouts), timeout)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		var currentAttemptErr error
		var putErr error
		_, putErr = minioClient.PutObject(ctx, bucketName, calcResultsObjectName, &calcResultsBuf, int64(calcResultsBuf.Len()), minio.PutObjectOptions{ContentType: "application/octet-stream"})
		if putErr != nil {
			currentAttemptErr = multierr.Append(currentAttemptErr, putErr)
		}

		_, putErr = minioClient.PutObject(ctx, bucketName, drainedResultsObjectName, &drainedResultsBuf, int64(drainedResultsBuf.Len()), minio.PutObjectOptions{ContentType: "application/octet-stream"})
		if putErr != nil {
			currentAttemptErr = multierr.Append(currentAttemptErr, putErr)
		}

		if currentAttemptErr == nil {
			uploadSuccess = true
			cancel()
			break
		}
		lastErr = currentAttemptErr
		cancel()
	}

	if !uploadSuccess {
		log.Printf("failed to upload objects after %d attempts: %v", len(timeouts), lastErr)
		return errors.InternalServerError("producer.produce.minio", "failed to upload objects after multiple retries: %v", lastErr)
	}

	slog.Info("Uploads finished – sending notification via RabbitMQ")

	// --- Resilient RabbitMQ Connection with Retries ---
	var conn *amqp.Connection
	var dialErr error

	// Retry strategy: 1 initial + 3 short (5s) + 7 long (15s) = 11 attempts total
	// Total wait time before failure: (3 * 5s) + (7 * 15s) = 15s + 105s = 2 minutes.
	for i := 0; i < 11; i++ {
		var attemptTimeout time.Duration
		if i < 4 { // First 4 attempts (1 initial + 3 retries)
			attemptTimeout = 10 * time.Second
		} else { // Next 7 attempts
			attemptTimeout = 20 * time.Second
		}

		log.Printf("Attempt %d/11: Connecting to RabbitMQ...", i+1)
		conn, dialErr = amqp.DialConfig(Cfg.RabbitURL, amqp.Config{
			Dial: amqp.DefaultDial(attemptTimeout),
		})
		if dialErr == nil {
			log.Println("Successfully connected to RabbitMQ.")
			break // Success
		}
		log.Printf("RabbitMQ connection attempt %d failed: %v", i+1, dialErr)

		if i < 10 { // Don't sleep after the last attempt
			var sleepDuration time.Duration
			if i < 3 { // For the first 3 retries
				sleepDuration = 5 * time.Second
			} else { // For the next 7 retries
				sleepDuration = 15 * time.Second
			}
			log.Printf("Retrying in %v...", sleepDuration)
			time.Sleep(sleepDuration)
		}
	}

	if dialErr != nil { // If connection is still nil after all retries
		log.Printf("failed to connect to rabbitmq after multiple retries: %v", dialErr)
		return errors.InternalServerError("producer.produce.rabbitmq", "failed to connect to rabbitmq: %v", dialErr)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("failed to open a channel: %v", err)
		return errors.InternalServerError("producer.produce.rabbitmq", "failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"buckets", // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Printf("failed to declare a queue: %v", err)
		return errors.InternalServerError("producer.produce.rabbitmq", "failed to declare a queue: %v", err)
	}

	notification := map[string]string{
		"bucket_name":            bucketName,
		"calc_results_object":    calcResultsObjectName,
		"drained_results_object": drainedResultsObjectName,
	}
	body, err := json.Marshal(notification)
	if err != nil {
		log.Printf("failed to marshal notification: %v", err)
		return errors.InternalServerError("producer.produce.rabbitmq", "failed to marshal notification: %v", err)
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Printf("failed to publish a message: %v", err)
		return errors.InternalServerError("producer.produce.rabbitmq", "failed to publish a message: %v", err)
	}

	rsp.BucketName = bucketName
	rsp.CalcResultsObject = calcResultsObjectName
	rsp.DrainedResultsObject = drainedResultsObjectName
	log.Printf("Successfully produced data and stored in bucket %s as %s and %s", bucketName, calcResultsObjectName, drainedResultsObjectName)
	slog.Info("Produce request processing completed successfully")

	return nil
}

func toIntSlice(in []int32) []int {
	out := make([]int, len(in))
	for i, v := range in {
		out[i] = int(v)
	}
	return out
}
