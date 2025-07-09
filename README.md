# Analabit - University Admission Analytics Platform

Analabit is a comprehensive platform for analyzing university admission data and providing insights into admission probabilities.

## Architecture

The platform consists of several microservices:

- **API Service** (`service/api`): Public-facing REST API built with Go Fiber
- **Producer Service** (`service/producer`): Internal service for data crawling and processing
- **Aggregator Service** (`service/aggregator`): Internal service for data aggregation
- **Core Libraries** (`core/`): Shared business logic and data models

## Monitoring & Metrics

The platform includes comprehensive monitoring with Prometheus and Grafana:

- **Metrics Collection**: The API service is instrumented with Prometheus metrics to track request counts, latency, and status codes
- **Secure Dashboard**: Grafana dashboard for visualizing API performance
- **Prometheus**: Time-series database for storing metrics
- **No Internal Exposure**: Only the public API service exposes metrics; internal services (producer, aggregator) do not expose metrics to the internet

### Quick Start with Monitoring

1. **Set Environment Variables** (for production):
   ```bash
   export GRAFANA_USER="your-admin-username"
   export GRAFANA_PASSWORD="your-secure-password"
   ```

2. **Start the Development Environment**:
   ```bash
   # Start dependencies and monitoring stack
   docker compose -f docker-compose.dev.yml up -d
   
   # Start the application services
   ./scripts/dev.sh
   ```

3. **Start the Production Environment**:
   ```bash
   # Ensure environment variables are set in .env.prod
   docker compose -f docker-compose.prod.yml --env-file .env.prod up -d
   ```

### Access Monitoring

- **API**: http://localhost:8080
- **API Metrics**: http://localhost:8080/metrics
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000

### Default Credentials

**Development**:
- Grafana: `admin` / `admin`

**Production**:
- Grafana: Use `GRAFANA_USER` / `GRAFANA_PASSWORD` environment variables

### Validate Monitoring Setup

Run the validation script to ensure everything is working:

```bash
./scripts/validate_metrics.sh
```

### Available Metrics

The API service exposes the following Prometheus metrics:

- `http_requests_total`: Total number of HTTP requests
- `http_request_duration_seconds`: Request duration histogram
- `http_request_size_bytes`: Request size histogram
- `http_response_size_bytes`: Response size histogram

### Dashboards

Pre-configured Grafana dashboards include:

- **Analabit API Metrics**: Comprehensive API performance dashboard
  - Request rate (QPS)
  - Response latencies (50th, 95th, 99th percentiles)
  - Status code distribution
  - Top endpoints by traffic

### Monitoring Architecture

```
Client → API Service → Prometheus ← Grafana
            ↓
        Metrics (/metrics)
```

**Note**: Only the API service exposes metrics as it's the only public-facing service. Internal microservices (producer, aggregator) are not instrumented to maintain security and architectural simplicity.

## Development

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- Node.js (for frontend components)

### Local Development

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd analabit
   ```

2. **Start dependencies**:
   ```bash
   docker compose -f docker-compose.dev.yml up -d
   ```

3. **Run the application**:
   ```bash
   ./scripts/dev.sh
   ```

### API Endpoints

The API provides the following endpoints:

- `GET /health` - Health check
- `GET /metrics` - Prometheus metrics
- `GET /v1/varsities` - List universities
- `GET /v1/headings` - List study programs
- `GET /v1/headings/:id` - Get specific program
- `GET /v1/applications` - List applications
- `GET /v1/students/:id` - Get student data
- `GET /v1/results` - Get calculation results

## Deployment

### Production Deployment

1. **Set up environment variables**:
   ```bash
   cp .env.example .env.prod
   # Edit .env.prod with your configuration
   ```

2. **Deploy**:
   ```bash
   docker compose -f docker-compose.prod.yml --env-file .env.prod up -d
   ```

3. **Validate deployment**:
   ```bash
   ./scripts/validate_metrics.sh
   ```

### Environment Variables

**Required for Production**:
- `GRAFANA_USER`: Grafana admin username
- `GRAFANA_PASSWORD`: Grafana admin password
- `POSTGRES_PASSWORD`: PostgreSQL password
- `POSTGRES_USER`: PostgreSQL username
- `POSTGRES_DB`: PostgreSQL database name
- `ANALABIT_DB_USER`: Application database user
- `ANALABIT_DB_PASSWORD`: Application database password
- `ANALABIT_DB_NAME`: Application database name
- `RABBITMQ_USER`: RabbitMQ username
- `RABBITMQ_PASSWORD`: RabbitMQ password
- `MINIO_ROOT_USER`: MinIO admin username
- `MINIO_ROOT_PASSWORD`: MinIO admin password

## Monitoring Best Practices

1. **Retention**: Prometheus data is retained according to default settings. Adjust `--storage.tsdb.retention.time` if needed.

2. **Security**: In production, Prometheus and Grafana are bound to localhost. Use a reverse proxy (nginx) for external access with proper authentication.

3. **Alerting**: Consider adding Prometheus AlertManager for critical alerts.

4. **Backup**: Regularly backup Grafana dashboards and Prometheus data volumes.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.