#!/bin/bash
set -euo pipefail

# Upload script for Analabit project
# This script copies files from local to production server

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Configuration
PROD_HOST="prod"
PROD_USER="analabit"
LOCAL_BACKEND_DIR="/home/yegor/Prestart/analabit"
LOCAL_FRONTEND_DIR="/home/yegor/Prestart/webui"
REMOTE_BACKEND_DIR="/opt/analabit"
REMOTE_FRONTEND_DIR="/opt/analabit-webui"

# Check if directories exist
if [[ ! -d "$LOCAL_BACKEND_DIR" ]]; then
    print_warning "Backend directory not found: $LOCAL_BACKEND_DIR"
    exit 1
fi

if [[ ! -d "$LOCAL_FRONTEND_DIR" ]]; then
    print_warning "Frontend directory not found: $LOCAL_FRONTEND_DIR"
    exit 1
fi

print_status "Uploading backend files..."
rsync -av --exclude='.git' --exclude='node_modules' --exclude='__pycache__' \
  "$LOCAL_BACKEND_DIR/" "$PROD_USER@$PROD_HOST:$REMOTE_BACKEND_DIR/"

print_status "Uploading frontend files..."
rsync -av --exclude='.git' --exclude='node_modules' --exclude='.next' --exclude='__pycache__' \
  "$LOCAL_FRONTEND_DIR/" "$PROD_USER@$PROD_HOST:$REMOTE_FRONTEND_DIR/"

print_status "Making scripts executable..."
ssh "$PROD_USER@$PROD_HOST" "chmod +x $REMOTE_BACKEND_DIR/scripts/*.sh"

print_status "Files uploaded successfully!"
print_status "Next steps:"
print_status "1. SSH to the server: ssh $PROD_USER@$PROD_HOST"
print_status "2. Run setup script: cd $REMOTE_BACKEND_DIR && ./scripts/setup-production.sh"
