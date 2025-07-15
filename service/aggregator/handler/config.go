package handler

import (
	"fmt"
	"strings"
)

// Config holds environment configuration for the aggregator service
// Uses caarlos0/env struct tags for binding

type config struct {
	RabbitURL       string `env:"RABBITMQ_URL" envDefault:"amqp://guest:guest@rabbitmq:5672/"`
	MinioEndpoint   string `env:"MINIO_ENDPOINT" envDefault:"minio:9000"`
	MinioAccessKey  string `env:"MINIO_ACCESS_KEY_ID" envDefault:"minioadmin"`
	MinioSecretKey  string `env:"MINIO_SECRET_ACCESS_KEY" envDefault:"minioadmin"`
	MinioUseSSL     bool   `env:"MINIO_USE_SSL" envDefault:"false"`
	MinioBucketName string `env:"MINIO_BUCKET_NAME" envDefault:"analabit-results"`

	// Primary database connection
	DatabaseHost     string `env:"DATABASE_HOST" envDefault:"localhost"`
	DatabasePort     string `env:"DATABASE_PORT" envDefault:"5432"`
	DatabaseUser     string `env:"DATABASE_USER" envDefault:"postgres"`
	DatabasePassword string `env:"DATABASE_PASSWORD" envDefault:""`
	DatabaseDBName   string `env:"DATABASE_DBNAME" envDefault:"postgres"`
	DatabaseSSLMode  string `env:"DATABASE_SSLMODE" envDefault:"disable"`

	// Optional replica database connection strings (comma-separated)
	PostgresReplicaConnStrings string `env:"POSTGRES_REPLICA_CONN_STRINGS"`

	// Cleanup configuration
	CleanupRetentionRuns int    `env:"CLEANUP_RETENTION_RUNS" envDefault:"5"`
	CleanupBackupDir     string `env:"CLEANUP_BACKUP_DIR" envDefault:"./backups"`
}

// GetPrimaryConnString builds the primary database connection string from individual components
func (c *config) GetPrimaryConnString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DatabaseHost, c.DatabasePort, c.DatabaseUser, c.DatabasePassword, c.DatabaseDBName, c.DatabaseSSLMode)
}

// GetReplicaConnStrings returns a slice of replica connection strings, excluding empty strings
func (c *config) GetReplicaConnStrings() []string {
	if c.PostgresReplicaConnStrings == "" {
		return nil
	}

	parts := strings.Split(c.PostgresReplicaConnStrings, ",")
	var replicas []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			replicas = append(replicas, trimmed)
		}
	}
	return replicas
}

// GetAllConnStrings returns primary connection string followed by replica connection strings
func (c *config) GetAllConnStrings() []string {
	connStrings := []string{c.GetPrimaryConnString()}
	connStrings = append(connStrings, c.GetReplicaConnStrings()...)
	return connStrings
}

// ParseConnStr parses a Postgres connection string into individual components
func ParseConnStr(connStr string) (user, password, host, port, dbname string, err error) {
	parts := strings.Split(connStr, " ")
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := kv[0]
		value := kv[1]
		switch key {
		case "user":
			user = value
		case "password":
			password = value
		case "host":
			host = value
		case "port":
			port = value
		case "dbname":
			dbname = value
		}
	}
	if user == "" || host == "" || dbname == "" || port == "" {
		return "", "", "", "", "", fmt.Errorf("missing required connection parameters")
	}
	return user, password, host, port, dbname, nil
}
