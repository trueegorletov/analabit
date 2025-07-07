package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	// Server configuration
	ServerPort string `env:"APP_PORT" envDefault:"8080"`

	// Database configuration
	DatabaseHost     string `env:"DATABASE_HOST" envDefault:"localhost"`
	DatabasePort     int    `env:"DATABASE_PORT" envDefault:"5432"`
	DatabaseUser     string `env:"DATABASE_USER" envDefault:"postgres"`
	DatabaseDBName   string `env:"DATABASE_DBNAME" envDefault:"postgres"`
	DatabasePassword string `env:"DATABASE_PASSWORD" envDefault:"postgres"`
	DatabaseSSLMode  string `env:"DATABASE_SSLMODE" envDefault:"disable"`

	// Legacy postgres connection string support
	PostgresConnStrings string `env:"POSTGRES_CONN_STRINGS"`

	// Message queue configuration
	RabbitURL string `env:"RABBITMQ_URL" envDefault:"amqp://guest:guest@rabbitmq:5672/"`

	// MinIO configuration
	MinioEndpoint   string `env:"MINIO_ENDPOINT" envDefault:"minio:9000"`
	MinioAccessKey  string `env:"MINIO_ACCESS_KEY_ID" envDefault:"minioadmin"`
	MinioSecretKey  string `env:"MINIO_SECRET_ACCESS_KEY" envDefault:"minioadmin"`
	MinioUseSSL     bool   `env:"MINIO_USE_SSL" envDefault:"false"`
	MinioBucketName string `env:"MINIO_BUCKET_NAME" envDefault:"analabit-results"`

	// Logging configuration
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}

var AppConfig Config

func LoadConfig() error {
	if err := env.Parse(&AppConfig); err != nil {
		log.Fatalf("failed to parse environment config: %v", err)
		return err
	}
	return nil
}
