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
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/streadway/amqp"
	"go-micro.dev/v5/errors"

	"analabit/core"
	"analabit/core/drainer"
	"analabit/core/registry"
	"analabit/service/producer/proto"
)

type Producer struct{}

// concurrency guard to ensure only one Produce execution at a time per service instance
var (
	produceRunning int32 // 0 = not running, 1 = running (atomic flag)
)

func (p *Producer) Produce(ctx context.Context, req *proto.ProduceRequest, rsp *proto.ProduceResponse) error {
	log.Println("Received Produce request")

	// Ensure single-flight execution
	if !atomic.CompareAndSwapInt32(&produceRunning, 0, 1) {
		slog.Warn("Produce request ignored – another Produce is already running")
		return errors.InternalServerError("producer.produce.busy", "another Produce request is already being processed")
	}
	// Reset flag when function returns
	defer atomic.StoreInt32(&produceRunning, 0)

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

	slog.Info("Starting crawl and cache phase")
	result, err := registry.CrawlWithOptions(registry.AllDefinitions, params)
	if err != nil {
		log.Printf("failed to crawl or cache: %v", err)
		return err
	}
	varsities := result.LoadedVarsities
	slog.Info("Crawl completed", "varsitiesLoaded", len(varsities))

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

	// Ensure bucket exists (create if not exists)
	exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
	if errBucketExists != nil {
		log.Printf("failed to check bucket existence: %v", errBucketExists)
		return errors.InternalServerError("producer.produce.minio", "failed to check bucket existence: %v", errBucketExists)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("failed to create bucket: %v", err)
			return errors.InternalServerError("producer.produce.minio", "failed to create bucket: %v", err)
		}
		log.Printf("Successfully created bucket %s", bucketName)
	}

	// Upload the GOB files to the bucket with unique object names
	_, err = minioClient.PutObject(ctx, bucketName, calcResultsObjectName, &calcResultsBuf, int64(calcResultsBuf.Len()), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		log.Printf("failed to upload calculation results: %v", err)
		return errors.InternalServerError("producer.produce.minio", "failed to upload calculation results: %v", err)
	}

	_, err = minioClient.PutObject(ctx, bucketName, drainedResultsObjectName, &drainedResultsBuf, int64(drainedResultsBuf.Len()), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		log.Printf("failed to upload drained results: %v", err)
		return errors.InternalServerError("producer.produce.minio", "failed to upload drained results: %v", err)
	}
	slog.Info("Uploads finished – sending notification via RabbitMQ")

	// Send a notification to RabbitMQ with the bucket and object names
	conn, err := amqp.Dial(Cfg.RabbitURL)
	if err != nil {
		log.Printf("failed to connect to rabbitmq: %v", err)
		return errors.InternalServerError("producer.produce.rabbitmq", "failed to connect to rabbitmq: %v", err)
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
