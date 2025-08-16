# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Quick Start
```bash
make dev              # Start development environment with hot reload
make test             # Run all tests
make proto            # Generate protobuf files
```

### Development Environment
- Uses `goreman` (process supervisor) with `Procfile.dev` to run multiple Go services
- Services auto-reload using `air` (hot reload for Go)
- Dependencies run in Docker: PostgreSQL (port 5433), RabbitMQ (5672), MinIO (9000), Prometheus (9090), Grafana (3500), FlareSolverr (8191)

### Service Ports in Development
- API: 8080
- Producer: 8082  
- Aggregator: 8083
- IDMSU: 8081

### Testing
```bash
go test ./...                    # Run all tests
go test ./core/calculate_test.go # Run specific test file
```

### Dependencies Check
```bash
make check-deps                  # Verify required tools (pdftotext)
```

## Architecture Overview

### Core Components

**Monorepo Structure:**
- `core/` - Shared business logic and data models
- `service/` - Microservices (api, producer, aggregator, idmsu)
- `cli/` - Command-line interface for data management
- `codegen/` - University-specific data generators
- `sample_data/` - Test data and examples

**Data Flow:**
1. **Data Sources** (`core/source/`) - University-specific crawlers and parsers for admission data
2. **Registry** (`core/registry/`) - Definitions for 13+ supported universities (HSE, MSU, MIPT, ITMO, etc.)
3. **Processing** (`core/calculate.go`) - Student admission calculations and ranking
4. **Storage** (`core/ent/`) - Ent ORM with PostgreSQL database
5. **API** (`service/api/`) - REST API for frontend consumption

### Key Data Models (Ent Schema)
- `Application` - Student applications with scores, rankings, competition types
- `Heading` - Study programs/majors with capacities and metadata  
- `Varsity` - Universities and their configurations
- `Run` - Data collection/calculation cycles with timestamps
- `DrainedResult` - Final admission simulation results

### University Data Sources
Each university in `core/source/` has custom HTTP crawlers and parsers:
- **File-based**: HSE (Excel), MIPT (HTML tables)
- **HTTP-based**: MSU, ITMO, SPbSU (with rate limiting and FlareSolverr for DDoS protection)
- **Competition Types**: Regular, BVI, TargetQuota, DedicatedQuota, SpecialQuota

### MSU ID Resolution Service (`service/idmsu/`)
Specialized microservice for mapping MSU internal IDs to canonical Gosuslugi IDs:
- Sophisticated matching algorithm with confidence scoring (see `docs/idmsu_matching_algorithm.md`)
- Layered caching strategy with hourly background refresh
- Handles BVI (strict positional), Regular (two-pass scoring), DedicatedQuota (many-to-one) competition types

### Microservices Communication
- **Producer** - Data collection orchestration with RabbitMQ
- **Aggregator** - Data processing and storage coordination  
- **API** - Public REST endpoints with Prometheus metrics
- **IDMSU** - MSU-specific ID resolution with caching

## Development Environment Setup

### Prerequisites
- Go 1.24+
- Docker and Docker Compose
- `pdftotext` (poppler-utils package)

### Environment Variables (set by scripts/dev.sh)
- Database: PostgreSQL on localhost:5433
- RabbitMQ: localhost:5672  
- MinIO: localhost:9000
- FlareSolverr: localhost:8191 (for DDoS protection)

### Hot Reload Configuration
- `.air.toml` configuration for Go services
- Services restart automatically on code changes

## Testing Strategy

### Test Files Location
- Unit tests alongside source files: `*_test.go`
- Integration tests in service directories
- Sample data in `sample_data/` for testing parsers

### Key Test Areas
- Parser validation for each university format
- Calculation algorithm correctness
- ID resolution confidence scoring
- API endpoint functionality

## Monitoring and Observability

### Metrics (Prometheus + Grafana)
- HTTP request metrics (rate, duration, status codes) 
- Custom API metrics at `/metrics` endpoint
- Grafana dashboard at port 3500
- Production: HTTPS with nginx reverse proxy and basic auth

### Health Checks
- API service: `/health` endpoint
- Individual service health via Procfile.dev

## Common Patterns

### Adding New University Support
1. Create parser in `core/source/<university>/`
2. Add registry definition in `core/registry/<university>/`
3. Update `core/registry/defs.go` AllDefinitions
4. Add sample data in `sample_data/<university>/`
5. Create tests for parser functionality

### Database Migrations
- Ent schema changes automatically generate migrations
- Custom migrations in `core/migrations/`
- Run via `client.Schema.Create()` in service startup

### Error Handling
- Universities have different rate limits and anti-bot measures
- Use FlareSolverr for sites with DDoS protection
- Implement retries with exponential backoff
- Cache parsed data to reduce API calls

## Production Deployment

### Docker Compose
- `docker-compose.prod.yml` for production services
- `docker-compose.dev.yml` for development dependencies
- Separate Dockerfiles per service

### Scripts
- `scripts/deploy.sh` - Production deployment
- `scripts/setup-production.sh` - Initial server setup
- `scripts/deploy_monitoring.sh` - Monitoring stack deployment

### Security
- Nginx reverse proxy with SSL
- Basic authentication for monitoring endpoints
- Environment-specific configurations
- No hardcoded credentials or API keys