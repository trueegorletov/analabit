#!/bin/bash
set -e

# Configuration
NGINX_CONF="/etc/nginx/sites-available/analabit"
BACKUP_FILE="${NGINX_CONF}.bak.$(date +%Y%m%d%H%M%S)"
PROM_USER="$1"
PROM_PASS="$2"

# Ensure parameters are provided
if [ -z "$PROM_USER" ] || [ -z "$PROM_PASS" ]; then
    echo "Usage: $0 <prometheus_user> <prometheus_password>"
    exit 1
fi

# Check if Nginx config exists
if [ ! -f "$NGINX_CONF" ]; then
    echo "Error: Nginx config file not found at $NGINX_CONF"
    exit 1
fi

# Create backup
echo "Creating backup of Nginx config: $BACKUP_FILE"
sudo cp "$NGINX_CONF" "$BACKUP_FILE"

# Check if monitoring locations already exist
if grep -q "location /prometheus/" "$NGINX_CONF"; then
    echo "Monitoring endpoints already configured in Nginx"
    exit 0
fi 

echo "Adding monitoring endpoints to Nginx configuration..."

# Create temporary files for the monitoring locations
cat > /tmp/monitoring.locations << EOF
    # Prometheus endpoint (restricted to authenticated users)
    location /prometheus/ {
        auth_basic "Prometheus";
        auth_basic_user_file /etc/nginx/.prometheus_htpasswd;
        proxy_pass http://localhost:9090/;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    # Grafana endpoint
    location /grafana/ {
        proxy_pass http://localhost:3500/;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
EOF

# Try to find HTTPS server block
if grep -q "listen 443 ssl" "$NGINX_CONF"; then
    echo "Found HTTPS server block, inserting monitoring locations..."
    
    # A simpler approach: find the last closing brace of the HTTPS server block
    # and insert our locations before it
    SERVER_BLOCK_START=$(grep -n -B 1 "server {" "$NGINX_CONF" | grep -B 1000 "listen 443 ssl" | tail -n 1 | cut -d: -f1)
    
    # Extract the file up to the start of the HTTPS server block
    sudo head -n "$SERVER_BLOCK_START" "$NGINX_CONF" > /tmp/nginx_start
    
    # Extract the server block to a temporary file
    sudo sed -n "$SERVER_BLOCK_START,\$p" "$NGINX_CONF" > /tmp/server_block
    
    # Find the first occurrence of a closing brace in the server block
    # This will be the end of the HTTPS server block
    SERVER_BLOCK_END=$(grep -n "}" /tmp/server_block | head -n 1 | cut -d: -f1)
    
    # Edit the server block to insert our locations before the closing brace
    head -n $(( SERVER_BLOCK_END - 1 )) /tmp/server_block > /tmp/server_start
    cat /tmp/monitoring.locations >> /tmp/server_start
    echo "}" >> /tmp/server_start
    
    # Get the rest of the file after the server block
    tail -n +$(( SERVER_BLOCK_START + SERVER_BLOCK_END )) "$NGINX_CONF" > /tmp/nginx_end
    
    # Combine everything back together
    cat /tmp/nginx_start /tmp/server_start /tmp/nginx_end > /tmp/nginx_new
    sudo mv /tmp/nginx_new "$NGINX_CONF"
    sudo chmod 644 "$NGINX_CONF"
    
    # Verify the changes were applied
    if grep -q "location /prometheus/" "$NGINX_CONF"; then
        echo "Successfully added monitoring locations to HTTPS server block"
    else
        echo "Modification failed, falling back to append method"
        # Fallback to appending a new server block
        cat > /tmp/new_server_block << EOF

# Monitoring server block for Prometheus and Grafana
server {
    listen 443 ssl;
    server_name analabit.ru;
    
# Include our monitoring locations
EOF
        cat /tmp/monitoring.locations >> /tmp/new_server_block
        cat >> /tmp/new_server_block << EOF
    
    ssl_certificate /etc/letsencrypt/live/analabit.ru/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/analabit.ru/privkey.pem;
    include /etc/letsencrypt/options-ssl-nginx.conf;
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem;
}
EOF
        sudo bash -c "cat /tmp/new_server_block >> '$NGINX_CONF'"
        echo "Added new server block for monitoring as a fallback"
    fi
else
    echo "No HTTPS server block found, adding new server block for monitoring..."
    
    # Create a new server block
    cat > /tmp/new_server_block << EOF

# Monitoring server block for Prometheus and Grafana
server {
    listen 443 ssl;
    server_name analabit.ru;
    
EOF
    cat /tmp/monitoring.locations >> /tmp/new_server_block
    cat >> /tmp/new_server_block << EOF
    
    ssl_certificate /etc/letsencrypt/live/analabit.ru/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/analabit.ru/privkey.pem;
    include /etc/letsencrypt/options-ssl-nginx.conf;
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem;
}
EOF
    sudo bash -c "cat /tmp/new_server_block >> '$NGINX_CONF'"
    echo "Added new server block for monitoring"
fi

# Clean up temporary files
rm -f /tmp/nginx_start /tmp/server_block /tmp/server_start /tmp/nginx_end /tmp/monitoring.locations /tmp/new_server_block

# Create htpasswd file for Prometheus authentication
echo "Setting up Prometheus authentication..."
sudo htpasswd -bc /etc/nginx/.prometheus_htpasswd "$PROM_USER" "$PROM_PASS"

# Test and reload Nginx
echo "Testing Nginx configuration..."
sudo nginx -t && sudo systemctl reload nginx
echo "Nginx configuration updated and reloaded successfully"
