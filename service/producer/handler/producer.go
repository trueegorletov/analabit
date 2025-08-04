package handler

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/streadway/amqp"
	"go-micro.dev/v5/errors"

	"github.com/trueegorletov/analabit/core"
	"github.com/trueegorletov/analabit/core/drainer"
	"github.com/trueegorletov/analabit/core/registry"
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/source/flaresolverr"
	"github.com/trueegorletov/analabit/service/producer/proto"

	"go.uber.org/multierr"
)

type Producer struct{}

// loadSpbstuFromPayload loads SPbSTU data from a pre-serialized UploadPayload gob file
// and converts it back to a VarsityCalculator for use in calculations
func loadSpbstuFromPayload(minioClient *minio.Client, bucketName, objectName string, ctx context.Context) (*source.Varsity, error) {
	slog.Info("Loading SPbSTU fallback from payload", "objectName", objectName)

	// Download the gob file from MinIO
	obj, err := minioClient.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get SPbSTU fallback object %s: %w", objectName, err)
	}
	defer obj.Close()

	// Deserialize the UploadPayload
	var payload core.UploadPayload
	if err := gob.NewDecoder(obj).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode SPbSTU fallback payload from object %s: %w", objectName, err)
	}

	// Create a new VarsityCalculator
	vc := core.NewVarsityCalculator(payload.VarsityCode, payload.VarsityName)

	// Add headings
	for _, headingDTO := range payload.Headings {
		// Parse the FullCode to extract just the heading code (remove "spbstu:" prefix)
		headingCode := headingDTO.Code
		if len(headingCode) > 7 && headingCode[:7] == "spbstu:" {
			headingCode = headingCode[7:]
		}

		capacities := core.Capacities{
			Regular:        headingDTO.RegularCapacity,
			TargetQuota:    headingDTO.TargetQuotaCapacity,
			DedicatedQuota: headingDTO.DedicatedQuotaCapacity,
			SpecialQuota:   headingDTO.SpecialQuotaCapacity,
		}
		vc.AddHeading(headingCode, capacities, headingDTO.Name)
	}

	// Add applications
	for _, appDTO := range payload.Applications {
		// Parse the heading code from FullCode
		headingCode := appDTO.HeadingCode
		if len(headingCode) > 7 && headingCode[:7] == "spbstu:" {
			headingCode = headingCode[7:]
		}

		vc.AddApplication(
			headingCode,
			appDTO.StudentID,
			appDTO.RatingPlace,
			appDTO.Priority,
			appDTO.CompetitionType,
			appDTO.Score,
		)
	}

	// Set original submitted status for students
	for _, studentDTO := range payload.Students {
		if studentDTO.OriginalSubmitted {
			vc.SetOriginalSubmitted(studentDTO.ID)
		}
	}

	// Normalize applications after loading all data
	vc.NormalizeApplications()

	// Create a dummy VarsityDefinition for SPbSTU
	def := source.VarsityDefinition{
		Code:           "spbstu",
		Name:           payload.VarsityName,
		HeadingSources: []source.HeadingSource{}, // Empty since we're loading from fallback
	}

	// Create Varsity with the loaded calculator
	varsity := &source.Varsity{
		VarsityDefinition: &def,
		VarsityCalculator: vc,
		MSUInternalIDs:    make(map[string]string), // Empty for SPbSTU
	}

	slog.Info("Successfully loaded SPbSTU fallback", "headings", len(payload.Headings), "students", len(payload.Students), "applications", len(payload.Applications))
	return varsity, nil
}

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

	// Initialize FlareSolverr session management for this iteration
	if err := flaresolverr.StartForIteration(); err != nil {
		slog.Warn("Failed to initialize FlareSolverr sessions", "error", err)
		// Continue without FlareSolverr - some sources might still work
	}

	// Ensure cleanup happens regardless of success or failure
	defer func() {
		if err := flaresolverr.StopForIteration(); err != nil {
			slog.Error("Failed to cleanup FlareSolverr sessions", "error", err)
		}
	}()

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

const maxWorkerCount = 8

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

	// SPbSTU fallback logic: if enabled and SPbSTU is not in loaded varsities, load it from fallback
	if Cfg.SpbstuFallbackEnabled {
		spbstuFound := false
		for _, v := range varsities {
			if v.Code == "spbstu" {
				spbstuFound = true
				break
			}
		}

		if !spbstuFound {
			slog.Info("SPbSTU fallback enabled and SPbSTU not found in loaded varsities, loading from fallback")

			// Initialize MinIO client for fallback loading
			minioClient, err := minio.New(Cfg.MinioEndpoint, &minio.Options{
				Creds:  credentials.NewStaticV4(Cfg.MinioAccessKey, Cfg.MinioSecretKey, ""),
				Secure: Cfg.MinioUseSSL,
			})
			if err != nil {
				slog.Warn("Failed to initialize MinIO client for SPbSTU fallback", "error", err)
			} else {
				// Load SPbSTU from fallback gob file
				spbstuVarsity, err := loadSpbstuFromPayload(minioClient, Cfg.MinioBucketName, Cfg.SpbstuFallbackGobName, ctx)
				if err != nil {
					slog.Warn("Failed to load SPbSTU from fallback, continuing without it", "error", err)
				} else {
					// Add SPbSTU to the varsities list
					varsities = append(varsities, spbstuVarsity)
					slog.Info("Successfully added SPbSTU from fallback to varsities list")
				}
			}
		} else {
			slog.Info("SPbSTU fallback enabled but SPbSTU was found in loaded varsities, using normally loaded version")
		}
	}

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

	// Initialize drainedResults map
	for _, v := range varsities {
		drainedResults[v.Code] = make(map[int][]drainer.DrainedResult)
	}

	// Define job structure for parallel drainer execution
	type drainerJob struct {
		varsity *source.Varsity
		stage   int
		code    string
	}

	// Calculate total jobs and determine worker count
	totalJobs := len(varsities) * len(params.DrainStages)
	workerCount := totalJobs
	if workerCount > maxWorkerCount {
		workerCount = maxWorkerCount
	}

	slog.Info("Starting parallel drain simulations", "totalJobs", totalJobs, "workers", workerCount)

	// Create job channel and populate it
	jobs := make(chan drainerJob, totalJobs)
	for _, v := range varsities {
		for _, stage := range params.DrainStages {
			jobs <- drainerJob{
				varsity: v,
				stage:   stage,
				code:    v.Code,
			}
		}
	}
	close(jobs)

	// Synchronization primitives
	var wgDrainer sync.WaitGroup
	var muDrainer sync.Mutex
	drainerErrors := make(chan error, totalJobs)

	// Start worker goroutines
	for i := 0; i < workerCount; i++ {
		wgDrainer.Add(1)
		go func(workerID int) {
			defer wgDrainer.Done()
			for job := range jobs {
				// Create drainer instance and run simulation
				drainerInstance := drainer.New(job.varsity, job.stage)
				drainedResultSlice := drainerInstance.Run(params.DrainIterations)

				// Safely write results to shared map
				muDrainer.Lock()
				drainedResults[job.code][job.stage] = drainedResultSlice
				muDrainer.Unlock()

				slog.Info("Completed drainer job", "worker", workerID, "varsity", job.code, "stage", job.stage)
			}
		}(i)
	}

	// Wait for all jobs to complete
	wgDrainer.Wait()
	close(drainerErrors)

	// Check for any errors during drainer execution
	var allDrainerErrs error
	for err := range drainerErrors {
		allDrainerErrs = multierr.Append(allDrainerErrs, err)
	}
	if allDrainerErrs != nil {
		slog.Error("One or more drainer simulations failed", "errors", allDrainerErrs)
	}

	slog.Info("Drain simulations completed")

	// --- Per-Varsity Payload Creation and Upload ---
	var allObjectNames []string
	var wgUpload sync.WaitGroup
	var muUpload sync.Mutex
	uploadErrors := make(chan error, len(varsities))

	slog.Info("Encoding results and uploading to object storage for each varsity")
	minioClient, err := minio.New(Cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(Cfg.MinioAccessKey, Cfg.MinioSecretKey, ""),
		Secure: Cfg.MinioUseSSL,
	})
	if err != nil {
		log.Printf("failed to initialize minio client: %v", err)
		return errors.InternalServerError("producer.produce.minio", "failed to initialize minio client: %v", err)
	}

	bucketName := Cfg.MinioBucketName
	if bucketName == "" {
		return errors.InternalServerError("producer.produce.minio", "Minio bucket name is not set in Cfg")
	}

	// Ensure bucket exists once before starting parallel uploads
	err = ensureBucketExists(minioClient, bucketName)
	if err != nil {
		return err // The error is already logged and formatted
	}

	for _, v := range varsities {
		wgUpload.Add(1)
		go func(v *source.Varsity) {
			defer wgUpload.Done()

			// 1. Prepare DTOs
			drainedDTOs := make(map[int][]core.DrainedResultDTO)
			if drainedStages, ok := drainedResults[v.Code]; ok {
				for stage, results := range drainedStages {
					drainedDTOs[stage] = drainer.NewDrainedResultDTOs(results)
				}
			}

			// 2. Create Payload
			payload := core.NewUploadPayloadFromCalculator(v.VarsityCalculator, primaryResults[v.Code], drainedDTOs, v.MSUInternalIDs)

			// 3. Encode Payload
			var payloadBuf bytes.Buffer
			if err := gob.NewEncoder(&payloadBuf).Encode(payload); err != nil {
				log.Printf("failed to encode payload for varsity %s: %v", v.Code, err)
				uploadErrors <- fmt.Errorf("failed to encode payload for varsity %s: %w", v.Code, err)
				return
			}

			// 4. Upload to MinIO
			objectUUID := uuid.New().String()
			objectName := fmt.Sprintf("payload_%s_%s.gob", v.Code, objectUUID)

			err := uploadObjectWithRetry(minioClient, bucketName, objectName, &payloadBuf)
			if err != nil {
				log.Printf("failed to upload payload for varsity %s: %v", v.Code, err)
				uploadErrors <- fmt.Errorf("failed to upload payload for varsity %s: %w", v.Code, err)
				return
			}

			// 5. Collect object name
			muUpload.Lock()
			allObjectNames = append(allObjectNames, objectName)
			muUpload.Unlock()
			slog.Info("Successfully uploaded payload", "varsity", v.Code, "object", objectName)
		}(v)
	}

	wgUpload.Wait()
	close(uploadErrors)

	// Check for any errors during upload
	var allUploadErrs error
	for err := range uploadErrors {
		allUploadErrs = multierr.Append(allUploadErrs, err)
	}
	if allUploadErrs != nil {
		return errors.InternalServerError("producer.produce.upload", "one or more uploads failed: %v", allUploadErrs)
	}

	slog.Info("All uploads finished – sending notification via RabbitMQ")

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

	notification := map[string]interface{}{
		"bucket_name":     bucketName,
		"payload_objects": allObjectNames,
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
	rsp.PayloadObjects = allObjectNames
	log.Printf("Successfully produced data and stored in bucket %s", bucketName)
	slog.Info("Produce request processing completed successfully")

	return nil
}

// ensureBucketExists checks if a bucket exists and creates it if not, with retry logic.
func ensureBucketExists(minioClient *minio.Client, bucketName string) error {
	timeouts := []time.Duration{1 * time.Minute, 2 * time.Minute, 3 * time.Minute}
	var lastErr error

	for i, timeout := range timeouts {
		log.Printf("Attempt %d/%d: Checking/creating MinIO bucket '%s' (timeout: %v)", i+1, len(timeouts), bucketName, timeout)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		exists, err := minioClient.BucketExists(ctx, bucketName)
		if err != nil {
			lastErr = err
			continue
		}
		if !exists {
			err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
			if err != nil {
				lastErr = err
				continue
			}
			log.Printf("Successfully created bucket %s", bucketName)
		}
		return nil // Success
	}

	log.Printf("failed to ensure bucket exists after %d attempts: %v", len(timeouts), lastErr)
	return errors.InternalServerError("producer.produce.minio", "failed to ensure bucket exists after multiple retries: %v", lastErr)
}

// uploadObjectWithRetry uploads an object to MinIO with retry logic.
func uploadObjectWithRetry(minioClient *minio.Client, bucket, objectName string, data *bytes.Buffer) error {
	timeouts := []time.Duration{3 * time.Minute, 5 * time.Minute, 8 * time.Minute}
	var lastErr error

	for i, timeout := range timeouts {
		log.Printf("Attempt %d/%d: Uploading %s to MinIO (timeout: %v)", i+1, len(timeouts), objectName, timeout)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// We need to create a new reader for each attempt as the reader is consumed.
		reader := bytes.NewReader(data.Bytes())

		_, err := minioClient.PutObject(ctx, bucket, objectName, reader, int64(reader.Len()), minio.PutObjectOptions{ContentType: "application/octet-stream"})
		if err == nil {
			return nil // Success
		}
		lastErr = err
	}

	log.Printf("failed to upload object %s after %d attempts: %v", objectName, len(timeouts), lastErr)
	return fmt.Errorf("failed to upload object after multiple retries: %w", lastErr)
}

func toIntSlice(in []int32) []int {
	out := make([]int, len(in))
	for i, v := range in {
		out[i] = int(v)
	}
	return out
}
