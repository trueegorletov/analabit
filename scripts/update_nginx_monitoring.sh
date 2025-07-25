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
        proxy_pass http://localhost:9090/prometheus/;
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
        proxy_set_header X-Forwarded-Host \$host;
        proxy_set_header X-Forwarded-Server \$host;
    }
EOF

# Try to find HTTPS server block
if grep -q "listen 443 ssl" "$NGINX_CONF"; then
    echo "Found HTTPS server block, inserting monitoring locations..."
    
    # Use awk to find the HTTPS server block and insert monitoring locations
    # This approach is more reliable than complex grep chains
    awk '
    BEGIN { 
        in_https_server = 0
        brace_count = 0
        server_start = 0
        monitoring_inserted = 0
    }
    
    # Look for server block start
    /^[[:space:]]*server[[:space:]]*{/ {
        if (!in_https_server) {
            server_start = NR
            brace_count = 1
            in_https_server = 1
            print $0
            next
        }
    }
    
    # If we are in a server block, check for HTTPS
    in_https_server && /listen[[:space:]]+443[[:space:]]+ssl/ {
        in_https_server = 2  # Mark as HTTPS server block
        print $0
        next
    }
    
    # Count braces when in server block
    in_https_server {
        # Count opening braces
        gsub(/{/, "&")
        brace_count += gsub(/{/, "&")
        
        # Count closing braces
        gsub(/}/, "&")
        close_braces = gsub(/}/, "&")
        brace_count -= close_braces
        
        # If this is the closing brace of the HTTPS server block
        if (in_https_server == 2 && brace_count == 0 && close_braces > 0) {
            # Insert monitoring locations before the closing brace
            if (!monitoring_inserted) {
                while ((getline line < "/tmp/monitoring.locations") > 0) {
                    print line
                }
                close("/tmp/monitoring.locations")
                monitoring_inserted = 1
            }
            print $0
            in_https_server = 0
            next
        }
        
        # If brace count reaches 0 but this was not HTTPS, reset
        if (brace_count == 0) {
            in_https_server = 0
        }
    }
    
    # Print all other lines as-is
    { print $0 }
    ' "$NGINX_CONF" > /tmp/nginx_new
    
    # Check if monitoring locations were inserted
    if grep -q "location /prometheus/" /tmp/nginx_new; then
        sudo mv /tmp/nginx_new "$NGINX_CONF"
        sudo chmod 644 "$NGINX_CONF"
        echo "Successfully added monitoring locations to HTTPS server block"
    else
        echo "AWK method failed, falling back to append method"
        rm -f /tmp/nginx_new
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
rm -f /tmp/monitoring.locations /tmp/new_server_block /tmp/nginx_new

# Create htpasswd file for Prometheus authentication
echo "Setting up Prometheus authentication..."
sudo htpasswd -bc /etc/nginx/.prometheus_htpasswd "$PROM_USER" "$PROM_PASS"

# Test and reload Nginx
echo "Testing Nginx configuration..."
sudo nginx -t && sudo systemctl reload nginx
echo "Nginx configuration updated and reloaded successfully"
