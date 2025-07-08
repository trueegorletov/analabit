#!/bin/bash
set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root
if [[ $EUID -eq 0 ]]; then
   print_error "This script should not be run as root"
   exit 1
fi

# Set up directories
APP_DIR="/opt/analabit"
DATA_DIR="/var/lib/analabit"
LOG_DIR="/var/log/analabit"

print_status "Setting up application directories..."
sudo mkdir -p $APP_DIR $DATA_DIR $LOG_DIR
sudo chown -R analabit:analabit $APP_DIR $DATA_DIR $LOG_DIR

# Generate secure passwords
print_status "Generating secure passwords..."
POSTGRES_PASSWORD=$(openssl rand -base64 32)
RABBITMQ_PASSWORD=$(openssl rand -base64 32)
MINIO_PASSWORD=$(openssl rand -base64 32)

# Create production environment file
print_status "Creating production environment file..."
cat > $APP_DIR/.env.prod << EOF
POSTGRES_USER=analabit_user
POSTGRES_PASSWORD=$POSTGRES_PASSWORD
POSTGRES_DB=analabit_db

RABBITMQ_USER=analabit_rabbit
RABBITMQ_PASSWORD=$RABBITMQ_PASSWORD

MINIO_ROOT_USER=analabit_minio
MINIO_ROOT_PASSWORD=$MINIO_PASSWORD

APP_ENV=production
APP_PORT=8080
LOG_LEVEL=info
EOF

# Secure the environment file
chmod 600 $APP_DIR/.env.prod

print_status "Created secure environment file at $APP_DIR/.env.prod"
print_warning "Passwords have been generated and stored securely."
print_warning "Make sure to backup this file and store credentials securely!"

# Create systemd service
print_status "Creating systemd service..."
sudo tee /etc/systemd/system/analabit.service > /dev/null << EOF
[Unit]
Description=Analabit Application
After=docker.service
Requires=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
User=analabit
Group=analabit
WorkingDirectory=$APP_DIR
ExecStart=/usr/local/bin/docker-compose -f docker-compose.prod.yml --env-file .env.prod up -d
ExecStop=/usr/local/bin/docker-compose -f docker-compose.prod.yml --env-file .env.prod down
ExecReload=/usr/local/bin/docker-compose -f docker-compose.prod.yml --env-file .env.prod restart
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

# Enable and start the service
sudo systemctl daemon-reload
sudo systemctl enable analabit.service

print_status "Systemd service created and enabled"
print_status "Setup complete! Next steps:"
print_status "1. Copy your application files to $APP_DIR"
print_status "2. Configure nginx"
print_status "3. Set up SSL certificates"
print_status "4. Start the service with: sudo systemctl start analabit"
