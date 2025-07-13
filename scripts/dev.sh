#!/usr/bin/env bash
set -euo pipefail

# Start dependency containers (RabbitMQ, MinIO, Postgres)
docker compose -f docker-compose.dev.yml up -d

echo "Dependencies are up."

# Ensure goreman is installed
if ! command -v goreman &> /dev/null; then
  echo "Installing goreman (process supervisor for Procfile)..."
  go install github.com/mattn/goreman@latest
  export PATH="$(go env GOPATH)/bin:$PATH"
fi

# Ensure air is installed
if ! command -v air &> /dev/null; then
  echo "Installing air (hot-reload for Go)..."
  curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
  export PATH="$(go env GOPATH)/bin:$PATH"
fi

# Export environment so that services connect to local container ports
export RABBITMQ_URL="amqp://guest:guest@localhost:5672/"
export ANALABIT_RABBIT_URL="$RABBITMQ_URL"
export MINIO_ENDPOINT="localhost:9000"
export ANALABIT_MINIO_ENDPOINT="$MINIO_ENDPOINT"
export MINIO_ACCESS_KEY_ID="minioadmin"
export MINIO_SECRET_ACCESS_KEY="minioadmin"
export ANALABIT_MINIO_ACCESS_KEY="minioadmin"
export ANALABIT_MINIO_SECRET_KEY="minioadmin"
export MINIO_USE_SSL="false"
export ANALABIT_MINIO_USE_SSL="false"
export MINIO_BUCKET_NAME="analabit-results"
export ANALABIT_MINIO_BUCKET_NAME="$MINIO_BUCKET_NAME"

export POSTGRES_CONN_STRINGS="host=localhost port=5433 user=postgres password=postgres dbname=postgres sslmode=disable"
export DATABASE_HOST="localhost"
export DATABASE_PORT="5433"
export DATABASE_USER="postgres"
export DATABASE_PASSWORD="postgres"
export DATABASE_DBNAME="postgres"
export DATABASE_SSLMODE="disable"
export CACHE_TTL_MINUTES="300"
export SELF_QUERY_PERIOD_MINUTES="150"
export DRAIN_SIM_ITERATIONS="10"
export VARSITIES_EXCLUDED="mirea,rzgmu,spbstu"

# FlareSolverr configuration for bypassing DDoS-Guard protection
export FLARESOLVERR_URL="http://localhost:8191"
export FLARESOLVERR_TIMEOUT_MS="60000"

# Launch all Go micro-services with live reload
goreman -f Procfile.dev start 
