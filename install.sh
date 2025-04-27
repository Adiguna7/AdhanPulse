#!/bin/bash

REPO_URL="https://github.com/Adiguna7/AdhanPulse/releases/download/v1.0.0/adhan_pulse-linux-amd64"
INSTALL_DIR="/usr/local/bin/adhan_pulse"
SERVICE_NAME="adhan_pulse.service"
SERVICE_FILE="/etc/systemd/system/$SERVICE_NAME"
WORKING_DIR="/var/lib/adhan_pulse"

if [ "$(id -u)" -ne 0 ]; then
    echo "Please run as root (sudo)"
    exit 1
fi

if [ ! -d "$WORKING_DIR" ]; then
    echo "Creating working directory $WORKING_DIR..."
    mkdir -p "$WORKING_DIR"
    chown nobody:nogroup "$WORKING_DIR"
fi

echo "Downloading the Go binary..."
curl -L "$REPO_URL" -o "$INSTALL_DIR"

if [ $? -ne 0 ]; then
    echo "Failed to download the binary. Please check the URL."
    exit 1
fi

chmod +x "$INSTALL_DIR"

echo "Creating systemd service file..."

cat <<EOF > "$SERVICE_FILE"
[Unit]
Description=Adhan Pulse is a lightweight service that calculates daily prayer times and sends timely desktop notifications.
After=network.target

[Service]
ExecStart=$INSTALL_DIR
Restart=always
User=nobody
Group=nogroup
WorkingDirectory=$WORKING_DIR

[Install]
WantedBy=multi-user.target
EOF

echo "Reloading systemd..."
systemctl daemon-reload

echo "Enabling and starting the service..."
systemctl enable "$SERVICE_NAME"
systemctl start "$SERVICE_NAME"

echo "Installation complete. The service is now running."
