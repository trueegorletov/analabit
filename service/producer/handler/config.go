package handler

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	CacheDir               string `env:"ANALABIT_CACHE_DIR" envDefault:"./cache"`
	CacheTTLMinutes        int    `env:"ANALABIT_CACHE_TTL_MINUTES" envDefault:"100"`
	DrainStages            []int  `env:"ANALABIT_DRAIN_SIM_STAGES" envSeparator:"," envDefault:"16,33,50,66"`
	DrainIterations        int    `env:"ANALABIT_DRAIN_SIM_ITERATIONS" envDefault:"100"`
	MinioEndpoint          string `env:"ANALABIT_MINIO_ENDPOINT" envDefault:"minio:9000"`
	MinioAccessKey         string `env:"ANALABIT_MINIO_ACCESS_KEY" envDefault:"minioadmin"`
	MinioSecretKey         string `env:"ANALABIT_MINIO_SECRET_KEY" envDefault:"minioadmin"`
	MinioUseSSL            bool   `env:"ANALABIT_MINIO_USE_SSL" envDefault:"false"`
	RabbitURL              string `env:"ANALABIT_RABBIT_URL" envDefault:"amqp://guest:guest@rabbitmq:5672/"`
	MinioBucketName        string `env:"ANALABIT_MINIO_BUCKET_NAME" envDefault:"analabit-results"`
	SelfQueryPeriodMinutes int    `env:"ANALABIT_SELF_QUERY_PERIOD_MINUTES" envDefault:"100"`
}

var Cfg Config

func init() {
	if err := env.Parse(&Cfg); err != nil {
		log.Fatalf("failed to parse env Cfg: %v", err)
	}
}
