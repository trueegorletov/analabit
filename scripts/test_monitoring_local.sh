#!/bin/bash
# Local testing script for monitoring setup

set -e

echo "===== Analabit Monitoring Local Testing ====="
echo "Starting local testing at $(date)"

# Ensure we have all the necessary directories
mkdir -p monitoring/grafana/dashboards
mkdir -p monitoring/grafana/provisioning/datasources
mkdir -p monitoring/grafana/provisioning/dashboards

# 1. Create a local .env file if it doesn't exist
if [ ! -f ".env" ]; then
    echo "Creating local .env file for testing..."
    cat << EOT > .env
# Local test environment variables
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=postgres

RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=minioadmin

# Analabit specific database
ANALABIT_DB_NAME=analabit_test
ANALABIT_DB_USER=analabit_test
ANALABIT_DB_PASSWORD=analabit_test

# Monitoring credentials
GRAFANA_USER=admin
GRAFANA_PASSWORD=admin

APP_ENV=development
APP_PORT=8080
LOG_LEVEL=debug
EOT
    echo ".env file created for local testing."
else
    echo ".env file already exists."
    # Check if Grafana credentials exist in .env, add if missing
    if ! grep -q "GRAFANA_USER" .env; then
        echo "Adding Grafana credentials to .env..."
        echo "" >> .env
        echo "# Monitoring credentials" >> .env
        echo "GRAFANA_USER=admin" >> .env
        echo "GRAFANA_PASSWORD=admin" >> .env
    fi
fi

# 2. Start the services with docker-compose
echo "Starting services with monitoring enabled..."
docker-compose -f docker-compose.dev.yml down
docker-compose -f docker-compose.dev.yml up -d

# 3. Wait for services to start up
echo "Waiting for services to start (20s)..."
sleep 20

# 4. Test the API service metrics endpoint
echo "Testing API metrics endpoint..."
curl -s http://localhost:8080/metrics | grep -q "go_" && echo "✅ API metrics endpoint is working" || echo "❌ API metrics endpoint is not accessible"

# 5. Test Prometheus
echo "Testing Prometheus..."
curl -s http://localhost:9090/-/healthy | grep -q "Prometheus" && echo "✅ Prometheus is working" || echo "❌ Prometheus is not accessible"

# 6. Test Prometheus target scraping
echo "Checking if Prometheus is scraping API metrics..."
TARGETS=$(curl -s http://localhost:9090/api/v1/targets | grep -o '"analabit-api"')
if [ ! -z "$TARGETS" ]; then
    echo "✅ Prometheus is configured to scrape API metrics"
else
    echo "❌ Prometheus is not scraping API metrics - check prometheus.yml configuration"
fi

# 7. Test Grafana
echo "Testing Grafana..."
curl -s http://localhost:3000/api/health | grep -q "ok" && echo "✅ Grafana is working" || echo "❌ Grafana is not accessible"

# 8. Test Grafana datasource
echo "Testing Grafana Prometheus datasource..."
curl -s -u admin:admin http://localhost:3000/api/datasources | grep -q "Prometheus" && echo "✅ Prometheus datasource is configured in Grafana" || echo "❌ Prometheus datasource is not configured in Grafana"

# 9. Test Grafana dashboard
echo "Testing Grafana dashboard provisioning..."
curl -s -u admin:admin http://localhost:3000/api/search?query=analabit | grep -q "dashboard" && echo "✅ Analabit dashboard is provisioned" || echo "❌ Analabit dashboard is not provisioned"

echo "===== Local testing completed at $(date) ====="
echo "Local monitoring URLs:"
echo "- API Metrics: http://localhost:8080/metrics"
echo "- Prometheus: http://localhost:9090/"
echo "- Grafana: http://localhost:3000/ (login with admin/admin)"
