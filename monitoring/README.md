# Monitoring Configuration Documentation

This directory contains all monitoring-related configuration for the Analabit platform.

## Directory Structure

```
monitoring/
├── prometheus.yml              # Prometheus main configuration
├── grafana/
│   ├── provisioning/
│   │   ├── datasources/        # Auto-configured data sources
│   │   │   └── prometheus.yml  # Prometheus data source config
│   │   └── dashboards/         # Dashboard provisioning config
│   │       └── dashboards.yml  # Dashboard provider configuration
│   └── dashboards/             # Pre-built dashboards
│       └── analabit-api-metrics.json  # Main API metrics dashboard
├── README.md                   # This documentation file
└── DEPLOYMENT_CHECKLIST.md     # Production deployment checklist
```

## Architecture Overview

The monitoring stack is designed to be:

- **Secure** - Access controlled via authentication with strong passwords
- **Lightweight** - Only the API service is instrumented to minimize overhead
- **Informative** - Key metrics are visualized in Grafana dashboards
- **Production-ready** - Configured for both development and production environments
- **Non-intrusive** - Avoids port conflicts with existing services (NextJS on port 3000)

```
                                   ┌──────────────┐
                                   │              │
                                   │   Grafana    │───┐
                                   │  (port 3500) │   │
                                   └──────▲───────┘   │
                                          │           │
                                          │           │
                                   ┌──────┴───────┐   │
                                   │              │   │
                       ┌──────────►  Prometheus  │   │
                       │           │              │   │
                       │           └──────────────┘   │
                       │                              │
┌──────────────┐       │                         ┌────▼─────────┐
│              │       │                         │              │
│  API Service ├───────┘                         │    Nginx     │
│              │                                 │  (Basic Auth)│
└──────────────┘                                 └──────────────┘
                                                       ▲
                                                       │
                                                       │
                                                       ▼
                                                    Users
```

## Prometheus Configuration

### Scrape Targets

- **analabit-api**: Scrapes the Fiber API service metrics at `/metrics`
  - Target: API service (configured to connect via container networking or direct host access)
  - Scrape interval: 5 seconds
  - Metrics path: `/metrics`
  - Important: Only the public-facing API service is monitored, internal services (producer, aggregator) are not exposed

### Key Settings

- Global scrape interval: 15 seconds
- Evaluation interval: 15 seconds
- Storage: Local TSDB in `/prometheus` volume

## API Metrics Implementation

The API service (`service/api/main.go`) implements a custom Prometheus metrics middleware with:

1. **Request Counter**: `http_requests_total`
   - Labels: method, endpoint, status_code
   - Tracks total requests by endpoint, method, and response code

2. **Request Duration**: `http_request_duration_seconds`
   - Labels: method, endpoint
   - Histogram of request duration in seconds
   - Used for latency percentiles (p50, p95, p99)

3. **Custom Metrics Endpoint**: `/metrics`
   - Exposes all metrics in Prometheus format
   - Secured in production via nginx with basic auth

## Grafana Configuration

### Data Sources

- **Prometheus**: Auto-provisioned to connect to the Prometheus service
  - URL: `http://prometheus:9090`
  - Default data source: Yes

### Dashboards

#### Analabit API Metrics Dashboard

Located at: `grafana/dashboards/analabit-api-metrics.json`

**Panels**:

1. **Request Rate (QPS)**: Shows requests per second by method, path, and status
   - Query: `rate(http_requests_total{job="analabit-api"}[1m])`

2. **Total Requests per Second**: Single stat showing overall QPS
   - Query: `sum(rate(http_requests_total{job="analabit-api"}[1m]))`

3. **Request Duration Percentiles**: 50th, 95th, and 99th percentile latencies
   - Queries: `histogram_quantile(0.50|0.95|0.99, ...)`

4. **Response Status Codes**: Stacked view of response codes over time
   - Query: `rate(http_requests_total{job="analabit-api"}[1m])`

5. **Top API Endpoints**: Table showing busiest endpoints
   - Query: `topk(10, sum by (path) (rate(http_requests_total{job="analabit-api"}[5m])))`

### Authentication

- **Development**: `admin` / `admin`
- **Production**: 
  - Uses `GRAFANA_USER` / `GRAFANA_PASSWORD` environment variables
  - Passwords are cryptographically secure, randomly generated
  - Additional security with nginx basic auth
  - All credentials are generated using the `setup_monitoring_auth.sh` script

### Port Configuration

- **Grafana**: Runs on port 3500 instead of default 3000 to avoid conflict with NextJS frontend
- **Prometheus**: Runs on port 9090
- **External Access**: All monitoring services are accessed through nginx on HTTPS
  - Grafana: https://analabit.ru/grafana/
  - Prometheus: https://analabit.ru/prometheus/

## Available Metrics

The fiberprometheus middleware automatically provides:

### HTTP Metrics

- `http_requests_total{method, path, status_code}`: Total number of HTTP requests
- `http_request_duration_seconds{method, path}`: Request duration histogram
- `http_request_size_bytes{method, path}`: Request body size histogram  
- `http_response_size_bytes{method, path}`: Response body size histogram

### Labels

- `method`: HTTP method (GET, POST, etc.)
- `path`: URL path (e.g., `/v1/varsities`)
- `status_code`: HTTP status code (200, 404, 500, etc.)
- `job`: Prometheus job name (`analabit-api`)

## Useful PromQL Queries

### Request Rate
```promql
# Total requests per second
sum(rate(http_requests_total[1m]))

# Requests per second by endpoint
sum by (path) (rate(http_requests_total[1m]))

# Error rate (4xx/5xx responses)
sum(rate(http_requests_total{status_code=~"[45].."}[1m])) / sum(rate(http_requests_total[1m]))
```

### Latency
```promql
# 95th percentile latency
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))

# Average latency by endpoint
sum by (path) (rate(http_request_duration_seconds_sum[1m])) / sum by (path) (rate(http_request_duration_seconds_count[1m]))
```

### Traffic Analysis
```promql
# Top 10 endpoints by request count
topk(10, sum by (path) (rate(http_requests_total[5m])))

# Requests by HTTP method
sum by (method) (rate(http_requests_total[1m]))
```

## Customization

### Adding New Dashboards

1. Create JSON dashboard file in `grafana/dashboards/`
2. Restart Grafana container or wait for auto-reload (10 seconds)

### Modifying Prometheus Config

1. Edit `prometheus.yml`
2. Restart Prometheus container or reload config:
   ```bash
   curl -X POST http://localhost:9090/-/reload
   ```

### Custom Metrics

To add custom metrics to the API service:

```go
import "github.com/prometheus/client_golang/prometheus"

var customCounter = prometheus.NewCounter(prometheus.CounterOpts{
    Name: "custom_operations_total",
    Help: "Total number of custom operations",
})

func init() {
    prometheus.MustRegister(customCounter)
}
```

## Security Considerations

1. **Production Binding**: Services bind to `127.0.0.1` in production
2. **Authentication**: Grafana requires login with strong passwords
3. **Network**: Use reverse proxy for external access
4. **Data**: No sensitive data should be included in metric labels
5. **Port Configuration**: Grafana runs on port 3500 to avoid conflict with NextJS frontend on port 3000
6. **Strong Passwords**: All production passwords are cryptographically secure and randomly generated
7. **Dual Authentication**: Both Grafana auth and nginx basic auth provide layered security
8. **Credential Management**: Secure scripts for generating and storing credentials
9. **Defense Against Attacks**: Protection against previous security threats like cryptocurrency miners
10. **Secure Communications**: All monitoring traffic passes through HTTPS

## Troubleshooting

### Prometheus Not Scraping

1. Check targets at http://localhost:9090/targets
2. Verify API service is running and `/metrics` is accessible
3. Check Docker network connectivity

### Grafana Not Loading Dashboards

1. Check dashboard provisioning logs in Grafana container
2. Verify JSON dashboard syntax
3. Ensure proper volume mounts in Docker Compose

### No Metrics Data

1. Generate API traffic: `curl http://localhost:8080/health`
2. Wait for scrape interval (5-15 seconds)
3. Check Prometheus query browser

### Port Conflicts

1. Check for port conflicts: `sudo netstat -tuln | grep 3500`
2. If Grafana fails to start, check logs: `docker logs analabit_grafana_1`
3. Modify port in docker-compose.yml and nginx configuration if needed
4. Restart services and nginx: `docker-compose restart grafana && sudo systemctl restart nginx`

### Authentication Issues

1. Regenerate Grafana password: `sed -i "s/^GRAFANA_PASSWORD=.*$/GRAFANA_PASSWORD=$(openssl rand -base64 32 | tr -dc 'a-zA-Z0-9!@#$%^&*()_+?><:{}[]' | head -c 24)/" .env`
2. Regenerate nginx basic auth: `sudo htpasswd -bc /etc/nginx/.htpasswd prometheus $(openssl rand -base64 32 | tr -dc 'a-zA-Z0-9!@#$%^&*()_+?><:{}[]' | head -c 24)`
3. Restart services: `docker-compose -f docker-compose.prod.yml restart && sudo systemctl restart nginx`

## Performance

### Resource Usage

- **Prometheus**: ~100MB RAM for basic setup
- **Grafana**: ~50MB RAM for basic setup
- **Disk**: ~1GB/day for default retention

### Optimization

1. Increase scrape intervals for high-traffic services
2. Use recording rules for complex queries
3. Configure appropriate retention policies
4. Consider federation for multiple Prometheus instances
