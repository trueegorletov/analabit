package handler

import (
	"context"
	"log"

	"analabit/service/producer/proto"
)

type Producer struct{}

func (p *Producer) Produce(ctx context.Context, req *proto.ProduceRequest, rsp *proto.ProduceResponse) error {
	log.Println("Received Produce request")

	// 1. Data Collection & Processing
	/** //TODO: implement it properly, like in root.go of cli/cmd
	sources := registry.AllDefinitions
	varsities := source.LoadFromDefinitions(sources)

	// 2. Calculates results.
	// We pass nil for ent.Client and cache, as the producer service does not interact with a database
	// and does not require caching for this operation.
	calcResults, err := core.Calculate(context.Background(), nil, nil, allSources, nil)
	if err != nil {
		log.Printf("failed to calculate results: %v", err)
		return errors.InternalServerError("producer.produce.calculate", "failed to calculate results: %v", err)
	}

	// 3. Drains the results.
	drainedResults, err := drainer.Drain(calcResults)
	if err != nil {
		log.Printf("failed to drain results: %v", err)
		return errors.InternalServerError("producer.produce.drain", "failed to drain results: %v", err)
	}

	// 4. Serializes both calculation and drained results into GOBs.
	var calcResultsBuf bytes.Buffer
	if err := gob.NewEncoder(&calcResultsBuf).Encode(calcResults); err != nil {
		log.Printf("failed to encode calculation results: %v", err)
		return errors.InternalServerError("producer.produce.gob", "failed to encode calculation results: %v", err)
	}

	var drainedResultsBuf bytes.Buffer
	if err := gob.NewEncoder(&drainedResultsBuf).Encode(drainedResults); err != nil {
		log.Printf("failed to encode drained results: %v", err)
		return errors.InternalServerError("producer.produce.gob", "failed to encode drained results: %v", err)
	}

	// 5. Creates a unique bucket in MinIO.
	// TODO: Get from config/env
	endpoint := "minio:9000"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Printf("failed to initialize minio client: %v", err)
		return errors.InternalServerError("producer.produce.minio", "failed to initialize minio client: %v", err)
	}

	bucketName := uuid.New().String()
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Check if bucket already exists
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("bucket %s already exists", bucketName)
		} else {
			log.Printf("failed to create bucket: %v", err)
			return errors.InternalServerError("producer.produce.minio", "failed to create bucket: %v", err)
		}
	} else {
		log.Printf("Successfully created bucket %s", bucketName)
	}

	// 6. Uploads the GOB files to the bucket.
	calcResultsObjectName := "calculation_results.gob"
	_, err = minioClient.PutObject(ctx, bucketName, calcResultsObjectName, &calcResultsBuf, int64(calcResultsBuf.Len()), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		log.Printf("failed to upload calculation results: %v", err)
		return errors.InternalServerError("producer.produce.minio", "failed to upload calculation results: %v", err)
	}

	drainedResultsObjectName := "drained_results.gob"
	_, err = minioClient.PutObject(ctx, bucketName, drainedResultsObjectName, &drainedResultsBuf, int64(drainedResultsBuf.Len()), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		log.Printf("failed to upload drained results: %v", err)
		return errors.InternalServerError("producer.produce.minio", "failed to upload drained results: %v", err)
	}

	// 7. Sends a notification to RabbitMQ with the bucket name.
	// TODO: Get from config/env
	rabbitURL := "amqp://guest:guest@rabbitmq:5672/"
	conn, err := amqp.Dial(rabbitURL)
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

	notification := map[string]string{"bucket_name": bucketName}
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
	log.Printf("Successfully produced data and stored in bucket %s", bucketName)
	*/

	return nil
}
