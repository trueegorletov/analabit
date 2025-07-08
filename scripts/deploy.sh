#!/bin/bash
set -euo pipefail

# Deployment script for Analabit
# This script should be run on the production server

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as analabit user
if [[ $USER != "analabit" ]]; then
   print_error "This script should be run as the analabit user"
   exit 1
fi

# Configuration
BACKEND_DIR="/opt/analabit"
FRONTEND_DIR="/opt/analabit-webui"
NGINX_CONFIG="/etc/nginx/sites-available/analabit"
NGINX_ENABLED="/etc/nginx/sites-enabled/analabit"

# Update and deploy backend
deploy_backend() {
    print_status "Deploying backend..."
    
    if [[ ! -d "$BACKEND_DIR" ]]; then
        print_error "Backend directory not found: $BACKEND_DIR"
        exit 1
    fi
    
    cd "$BACKEND_DIR"
    
    # Pull latest images
    docker-compose -f docker-compose.prod.yml --env-file .env.prod pull
    
    # Restart services
    docker-compose -f docker-compose.prod.yml --env-file .env.prod up -d
    
    # Wait for services to be ready
    print_status "Waiting for backend services to be ready..."
    sleep 30
    
    # Health check
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        print_status "Backend deployment successful!"
    else
        print_error "Backend health check failed!"
        exit 1
    fi
}

# Update and deploy frontend
deploy_frontend() {
    print_status "Deploying frontend..."
    
    if [[ ! -d "$FRONTEND_DIR" ]]; then
        print_error "Frontend directory not found: $FRONTEND_DIR"
        exit 1
    fi
    
    cd "$FRONTEND_DIR"
    
    # Pull latest changes
    git pull origin main
    
    # Install dependencies
    npm ci --production
    
    # Build application
    npm run build
    
    # Restart service
    sudo systemctl restart analabit-webui
    
    # Wait for service to be ready
    print_status "Waiting for frontend service to be ready..."
    sleep 10
    
    # Health check
    if curl -f http://localhost:3000 > /dev/null 2>&1; then
        print_status "Frontend deployment successful!"
    else
        print_error "Frontend health check failed!"
        exit 1
    fi
}

# Configure nginx
configure_nginx() {
    print_status "Configuring nginx..."
    
    # Copy nginx configuration
    sudo cp "$BACKEND_DIR/nginx.conf" "$NGINX_CONFIG"
    
    # Enable site
    sudo ln -sf "$NGINX_CONFIG" "$NGINX_ENABLED"
    
    # Test nginx configuration
    sudo nginx -t
    
    # Restart nginx
    sudo systemctl restart nginx
    
    print_status "Nginx configuration updated!"
}

# Main deployment function
main() {
    print_status "Starting Analabit deployment..."
    
    # Deploy components
    deploy_backend
    deploy_frontend
    configure_nginx
    
    # Clean up old Docker images
    print_status "Cleaning up old Docker images..."
    docker image prune -f
    
    print_status "Deployment completed successfully!"
    print_status "Services status:"
    systemctl status analabit --no-pager || true
    systemctl status analabit-webui --no-pager || true
    systemctl status nginx --no-pager || true
}

# Run main function
main "$@"
