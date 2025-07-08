#!/bin/bash
set -euo pipefail

# Create systemd service for frontend
sudo tee /etc/systemd/system/analabit-webui.service > /dev/null << 'EOF'
[Unit]
Description=Analabit Web UI
After=network.target

[Service]
Type=simple
User=analabit
Group=analabit
WorkingDirectory=/opt/analabit-webui
ExecStart=/usr/bin/npm start
Restart=always
RestartSec=10
Environment=NODE_ENV=production
Environment=PORT=3000
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

# Enable the service
sudo systemctl daemon-reload
sudo systemctl enable analabit-webui.service

echo "Frontend systemd service created and enabled"
