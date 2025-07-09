#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Setting up secure authentication for monitoring services...${NC}"

# Generate strong random password for Grafana
GRAFANA_USER="admin"
GRAFANA_PASSWORD=$(openssl rand -base64 20 | tr -dc 'a-zA-Z0-9!@#$%^&*()_+?><:{}[]' | head -c 20)

# Generate htpasswd credentials for nginx basic auth
PROMETHEUS_USER="prometheus_admin"
PROMETHEUS_PASSWORD=$(openssl rand -base64 20 | tr -dc 'a-zA-Z0-9!@#$%^&*()_+?><:{}[]' | head -c 20)

# Check if .env file exists
if [ ! -f /opt/analabit/.env ]; then
    echo -e "${RED}Error: .env file not found at /opt/analabit/.env${NC}"
    exit 1
fi

# Update .env file with Grafana credentials
echo -e "${YELLOW}Updating .env file with secure Grafana credentials...${NC}"
if grep -q "GRAFANA_USER" /opt/analabit/.env; then
    # Update existing entries
    sed -i "s/^GRAFANA_USER=.*$/GRAFANA_USER=$GRAFANA_USER/" /opt/analabit/.env
    sed -i "s/^GRAFANA_PASSWORD=.*$/GRAFANA_PASSWORD=$GRAFANA_PASSWORD/" /opt/analabit/.env
else
    # Add new entries
    echo "" >> /opt/analabit/.env
    echo "# Grafana credentials" >> /opt/analabit/.env
    echo "GRAFANA_USER=$GRAFANA_USER" >> /opt/analabit/.env
    echo "GRAFANA_PASSWORD=$GRAFANA_PASSWORD" >> /opt/analabit/.env
fi

# Create htpasswd file for nginx basic auth
echo -e "${YELLOW}Creating htpasswd file for nginx basic auth...${NC}"
if command -v htpasswd &> /dev/null; then
    echo -e "${GREEN}Using htpasswd command...${NC}"
    # Install apache2-utils if not available
    if ! command -v htpasswd &> /dev/null; then
        echo -e "${YELLOW}Installing apache2-utils to use htpasswd...${NC}"
        sudo apt-get update
        sudo apt-get install -y apache2-utils
    fi
    
    # Create htpasswd file
    sudo htpasswd -bc /etc/nginx/.htpasswd $PROMETHEUS_USER $PROMETHEUS_PASSWORD
else
    echo -e "${YELLOW}htpasswd not available, using openssl method...${NC}"
    # Create htpasswd file using openssl
    HASHED_PASSWORD=$(openssl passwd -apr1 $PROMETHEUS_PASSWORD)
    echo "$PROMETHEUS_USER:$HASHED_PASSWORD" | sudo tee /etc/nginx/.htpasswd > /dev/null
fi

# Set proper permissions
sudo chown www-data:www-data /etc/nginx/.htpasswd
sudo chmod 600 /etc/nginx/.htpasswd

echo -e "${GREEN}Generated secure credentials:${NC}"
echo -e "${YELLOW}Grafana:${NC}"
echo -e "  Username: ${GREEN}$GRAFANA_USER${NC}"
echo -e "  Password: ${GREEN}$GRAFANA_PASSWORD${NC}"
echo -e "${YELLOW}Nginx Basic Auth:${NC}"
echo -e "  Username: ${GREEN}$PROMETHEUS_USER${NC}"
echo -e "  Password: ${GREEN}$PROMETHEUS_PASSWORD${NC}"
echo -e "${GREEN}Authentication setup complete!${NC}"
echo -e "${YELLOW}IMPORTANT: Save these credentials in a secure password manager!${NC}"

# Record credentials in a secure file
cat << EOF > /tmp/monitoring_credentials.txt
Grafana:
  Username: $GRAFANA_USER
  Password: $GRAFANA_PASSWORD

Nginx Basic Auth (for Prometheus & Grafana):
  Username: $PROMETHEUS_USER
  Password: $PROMETHEUS_PASSWORD
EOF

echo -e "${GREEN}Credentials saved to /tmp/monitoring_credentials.txt${NC}"
echo -e "${RED}WARNING: Delete this file after saving the credentials securely!${NC}"
echo -e "${YELLOW}Run: shred -u /tmp/monitoring_credentials.txt${NC}"
