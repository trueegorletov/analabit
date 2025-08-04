package handler

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	CacheDir               string   `env:"CACHE_DIR" envDefault:"./cache"`
	CacheTTLMinutes        int      `env:"CACHE_TTL_MINUTES" envDefault:"30"`
	DrainStages            []int    `env:"DRAIN_SIM_STAGES" envSeparator:"," envDefault:"25,50,75,90"`
	DrainIterations        int      `env:"DRAIN_SIM_ITERATIONS" envDefault:"100"`
	MinioEndpoint          string   `env:"MINIO_ENDPOINT" envDefault:"minio:9000"`
	MinioAccessKey         string   `env:"MINIO_ACCESS_KEY_ID" envDefault:"minioadmin"`
	MinioSecretKey         string   `env:"MINIO_SECRET_ACCESS_KEY" envDefault:"minioadmin"`
	MinioUseSSL            bool     `env:"MINIO_USE_SSL" envDefault:"false"`
	RabbitURL              string   `env:"RABBITMQ_URL" envDefault:"amqp://guest:guest@rabbitmq:5672/"`
	MinioBucketName        string   `env:"MINIO_BUCKET_NAME" envDefault:"analabit-results"`
	VarsitiesList          []string `env:"VARSITIES_LIST" envSeparator:"," envDefault:"all"`
	VarsitiesExcluded      []string `env:"VARSITIES_EXCLUDED" envSeparator:"," envDefault:"spbstu"`
	SelfQueryPeriodMinutes int      `env:"SELF_QUERY_PERIOD_MINUTES" envDefault:"45"`
	// SPbSTU fallback configuration
	SpbstuFallbackEnabled bool   `env:"SPBSTU_FALLBACK_ENABLED" envDefault:"false"`
	SpbstuFallbackGobName string `env:"SPBSTU_FALLBACK_GOB_NAME" envDefault:"payload_spbstu_a9dc55c5-addd-4269-a3b9-b40b175dfa52.gob"`
}

var Cfg Config

func init() {
	if err := env.Parse(&Cfg); err != nil {
		log.Fatalf("failed to parse env Cfg: %v", err)
	}
}
