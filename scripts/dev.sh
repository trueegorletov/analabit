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
fi

# Define path to air, ensure it's executable, and add to PATH
AIR_BIN_PATH="$(go env GOPATH)/bin"
chmod +x "$AIR_BIN_PATH/air"
export PATH="$AIR_BIN_PATH:$PATH"

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

export POSTGRES_CONN_STRINGS="host=192.168.0.103 port=5433 user=postgres password=postgres dbname=postgres sslmode=disable"
export DATABASE_HOST="192.168.0.103"
export DATABASE_PORT="5433"
export DATABASE_USER="postgres"
export DATABASE_PASSWORD="postgres"
export DATABASE_DBNAME="postgres"
export DATABASE_SSLMODE="disable"
export CACHE_TTL_MINUTES="6000"
export SELF_QUERY_PERIOD_MINUTES="1500"
export DRAIN_SIM_ITERATIONS="10"
export VARSITIES_LIST="itmo,hse_spb,rzgmu"

# FlareSolverr configuration for bypassing DDoS-Guard protection
export FLARESOLVERR_URL="http://localhost:8191"

# Launch all Go micro-services with live reload
goreman -f Procfile.dev start
