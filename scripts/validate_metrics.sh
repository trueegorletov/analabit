#!/bin/bash
set -euo pipefail

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[!]${NC} $1"
}

# Configuration
API_URL="${API_URL:-http://localhost:8080}"
PROMETHEUS_URL="${PROMETHEUS_URL:-http://localhost:9090}"
GRAFANA_URL="${GRAFANA_URL:-http://localhost:3000}"

echo "=== Analabit Metrics Validation ==="
echo "Testing monitoring stack..."

# Test 1: Check API service metrics endpoint
echo -e "\n1. Testing API service metrics endpoint..."
if curl -s -f "${API_URL}/metrics" > /dev/null; then
    print_status "API metrics endpoint is accessible"
    
    # Check for specific metrics
    if curl -s "${API_URL}/metrics" | grep -q "http_requests_total"; then
        print_status "Found http_requests_total metric"
    else
        print_error "http_requests_total metric not found"
        exit 1
    fi
    
    if curl -s "${API_URL}/metrics" | grep -q "http_request_duration_seconds"; then
        print_status "Found http_request_duration_seconds metric"
    else
        print_error "http_request_duration_seconds metric not found"
        exit 1
    fi
else
    print_error "API metrics endpoint is not accessible at ${API_URL}/metrics"
    exit 1
fi

# Test 2: Check API health endpoint
echo -e "\n2. Testing API health endpoint..."
if curl -s -f "${API_URL}/health" > /dev/null; then
    print_status "API health endpoint is accessible"
else
    print_error "API health endpoint is not accessible at ${API_URL}/health"
    exit 1
fi

# Test 3: Check Prometheus endpoint
echo -e "\n3. Testing Prometheus..."
if curl -s -f "${PROMETHEUS_URL}/-/healthy" > /dev/null; then
    print_status "Prometheus is healthy"
else
    print_error "Prometheus is not accessible at ${PROMETHEUS_URL}"
    exit 1
fi

# Test 4: Query Prometheus for API metrics
echo -e "\n4. Testing Prometheus API queries..."
QUERY_URL="${PROMETHEUS_URL}/api/v1/query"

# Test http_requests_total metric
if curl -s -f "${QUERY_URL}?query=http_requests_total" | grep -q '"status":"success"'; then
    print_status "Successfully queried http_requests_total from Prometheus"
else
    print_warning "Could not query http_requests_total from Prometheus (may need time to collect data)"
fi

# Test request rate query
if curl -s -f "${QUERY_URL}?query=rate(http_requests_total[1m])" | grep -q '"status":"success"'; then
    print_status "Successfully queried request rate from Prometheus"
else
    print_warning "Could not query request rate from Prometheus (may need time to collect data)"
fi

# Test 5: Check Grafana
echo -e "\n5. Testing Grafana..."
if curl -s -f "${GRAFANA_URL}/api/health" > /dev/null; then
    print_status "Grafana is accessible"
    
    # Test Grafana login (basic check)
    if curl -s "${GRAFANA_URL}/login" | grep -q "Grafana"; then
        print_status "Grafana login page is accessible"
    else
        print_warning "Grafana login page may not be properly configured"
    fi
else
    print_error "Grafana is not accessible at ${GRAFANA_URL}"
    exit 1
fi

# Test 6: Generate some test traffic to API
echo -e "\n6. Generating test traffic for metrics..."
for i in {1..10}; do
    curl -s "${API_URL}/health" > /dev/null || true
    curl -s "${API_URL}/v1/varsities" > /dev/null || true
    curl -s "${API_URL}/v1/headings" > /dev/null || true
done
print_status "Generated test traffic to API endpoints"

# Wait a bit for metrics to be scraped
sleep 5

# Test 7: Verify metrics have data after traffic generation
echo -e "\n7. Verifying metrics collection after traffic..."
METRICS_WITH_DATA=$(curl -s "${QUERY_URL}?query=rate(http_requests_total[1m])" | grep -o '"value":\[.*\]' | grep -v ',"0"' | wc -l)
if [ "$METRICS_WITH_DATA" -gt 0 ]; then
    print_status "Metrics are collecting data from API traffic"
else
    print_warning "Metrics may not be collecting data yet (check Prometheus targets)"
fi

echo -e "\n=== Validation Complete ==="
print_status "All basic checks passed!"
echo ""
echo "Access URLs:"
echo "  API: ${API_URL}"
echo "  API Metrics: ${API_URL}/metrics"
echo "  Prometheus: ${PROMETHEUS_URL}"
echo "  Grafana: ${GRAFANA_URL}"
echo ""
echo "Next steps:"
echo "  1. Log into Grafana with configured credentials"
echo "  2. Check the 'Analabit API Metrics' dashboard"
echo "  3. Verify Prometheus targets at ${PROMETHEUS_URL}/targets"
echo "  4. Generate more API traffic to see metrics in action"
