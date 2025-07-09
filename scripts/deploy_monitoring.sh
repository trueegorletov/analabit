#!/bin/bash
# Production deployment script for Analabit monitoring setup
# To be executed on the production server

set -e

echo "===== Analabit Monitoring Setup Deployment ====="
echo "Starting deployment at $(date)"

# 1. Stop current production services
echo "Stopping current production services..."
cd /opt/analabit
docker-compose -f docker-compose.prod.yml down

# 2. Fetch the new monitoring-enabled branch
echo "Fetching the new monitoring-enabled branch..."
git fetch origin monitoring-setup
git checkout monitoring-setup

# 3. Update .env file with secure passwords for Prometheus and Grafana
echo "Updating .env file with secure Prometheus and Grafana credentials..."

# Generate secure passwords if they don't exist
if ! grep -q "GRAFANA_USER" .env; then
    echo "Adding Grafana admin credentials to .env..."
    GRAFANA_PASS=$(openssl rand -base64 32 | tr -dc 'a-zA-Z0-9!@#$%^&*()_+?><:{}[]' | head -c 24)
    echo "" >> .env
    echo "# Monitoring credentials" >> .env
    echo "GRAFANA_USER=admin" >> .env
    echo "GRAFANA_PASSWORD=${GRAFANA_PASS}" >> .env
    echo "Grafana admin password generated and added to .env"
    echo "Password: ${GRAFANA_PASS} (save this securely!)"
else
    echo "Grafana credentials already exist in .env, updating with strong password..."
    GRAFANA_PASS=$(openssl rand -base64 32 | tr -dc 'a-zA-Z0-9!@#$%^&*()_+?><:{}[]' | head -c 24)
    sed -i "s/^GRAFANA_PASSWORD=.*$/GRAFANA_PASSWORD=${GRAFANA_PASS}/" .env
    echo "Grafana admin password updated"
    echo "New password: ${GRAFANA_PASS} (save this securely!)"
fi

# 4. Start the monitoring setup
echo "Starting services with monitoring enabled..."
docker-compose -f docker-compose.prod.yml up -d

# 5. Test the monitoring setup
echo "Testing the monitoring setup..."
sleep 15  # Wait for services to start

# Check if API is accessible
if curl -sf http://localhost:8080/health > /dev/null; then
    echo "✅ API service is running"
else
    echo "❌ API service is not accessible"
    exit 1
fi

# Check if Prometheus is accessible
if curl -sf http://localhost:9090/-/healthy > /dev/null; then
    echo "✅ Prometheus is running"
else
    echo "❌ Prometheus is not accessible"
    exit 1
fi

# Check if Grafana is accessible
if curl -sf http://localhost:3500/api/health > /dev/null; then
    echo "✅ Grafana is running"
else
    echo "❌ Grafana is not accessible"
    exit 1
fi

echo "All services are running. Now configuring nginx..."

# 6. Update nginx configuration for Prometheus and Grafana
NGINX_CONF="/etc/nginx/sites-available/analabit"

# Check if we need to update nginx config
if ! grep -q "location /prometheus" "$NGINX_CONF"; then
    echo "Updating nginx configuration for Prometheus and Grafana..."
    
    # Create a backup of the current config
    sudo cp "$NGINX_CONF" "${NGINX_CONF}.bak.$(date +%Y%m%d%H%M%S)"
    
    # Add the monitoring endpoints to nginx config before the closing '}'
    sudo tee -a "$NGINX_CONF" > /dev/null << 'EOF'

    # Prometheus endpoint (restricted to authenticated users)
    location /prometheus/ {
        auth_basic "Prometheus";
        auth_basic_user_file /etc/nginx/.prometheus_htpasswd;
        proxy_pass http://localhost:9090/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Grafana endpoint
    location /grafana/ {
        proxy_pass http://localhost:3500/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
EOF

    # Create a secure htpasswd file for Prometheus
    PROM_USER="prometheus"
    PROM_PASS=$(openssl rand -base64 32 | tr -dc 'a-zA-Z0-9!@#$%^&*()_+?><:{}[]' | head -c 24)
    
    echo "Creating htpasswd file for Prometheus with username: $PROM_USER"
    sudo htpasswd -bc /etc/nginx/.prometheus_htpasswd "$PROM_USER" "$PROM_PASS"
    
    echo "Prometheus credentials:"
    echo "Username: $PROM_USER"
    echo "Password: $PROM_PASS"
    echo "IMPORTANT: Save these credentials securely in a password manager!"
    
    # Test nginx config
    sudo nginx -t
    
    # Reload nginx
    sudo systemctl reload nginx
    
    echo "Nginx configuration updated and reloaded"
else
    echo "Nginx already configured for monitoring endpoints"
fi

echo "===== Deployment completed successfully at $(date) ====="
echo "Monitoring URLs:"
echo "- Grafana: https://analabit.ru/grafana/ (login with admin and password from .env)"
echo "- Prometheus: https://analabit.ru/prometheus/ (login with prometheus user and generated password)"
