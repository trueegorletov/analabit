package handler

// Config holds environment configuration for the aggregator service
// Uses caarlos0/env struct tags for binding

type config struct {
	RabbitURL        string `env:"RABBITMQ_URL" envDefault:"amqp://guest:guest@rabbitmq:5672/"`
	MinioEndpoint    string `env:"MINIO_ENDPOINT" envDefault:"minio:9000"`
	MinioAccessKey   string `env:"MINIO_ACCESS_KEY_ID" envDefault:"minioadmin"`
	MinioSecretKey   string `env:"MINIO_SECRET_ACCESS_KEY" envDefault:"minioadmin"`
	MinioUseSSL      bool   `env:"MINIO_USE_SSL" envDefault:"false"`
	MinioBucketName  string `env:"MINIO_BUCKET_NAME" envDefault:"analabit-results"`
	DatabaseHost     string `env:"DATABASE_HOST" envDefault:"localhost"`
	DatabasePort     int    `env:"DATABASE_PORT" envDefault:"5432"`
	DatabaseUser     string `env:"DATABASE_USER" envDefault:"postgres"`
	DatabaseDBName   string `env:"DATABASE_DBNAME" envDefault:"postgres"`
	DatabasePassword string `env:"DATABASE_PASSWORD" envDefault:"postgres"`
	DatabaseSSLMode  string `env:"DATABASE_SSLMODE" envDefault:"disable"`
}
