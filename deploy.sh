#!/bin/bash

# ==========================================
# CONFIGURATION
# ==========================================
SERVER_IP="160.191.242.73"
SSH_USER="root"
PORT="22"
SSH_TARGET="$SSH_USER@$SERVER_IP"

REMOTE_BASE_DIR="service-go"
SERVICE_DIR="worker-service"
REMOTE_PATH="/root/$REMOTE_BASE_DIR/$SERVICE_DIR"
BINARY_NAME="worker-service"
SERVICE_FILE="/etc/systemd/system/$BINARY_NAME.service"

echo "================================================="
echo " START DEPLOYING: $BINARY_NAME"
echo " TARGET: $SSH_TARGET"
echo "================================================="

# 1. Build the binary
echo "[1/5] Building Go binary for Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -o ../$BINARY_NAME cmd/main.go

if [ $? -ne 0 ]; then
  echo "Error: Build failed!"
  exit 1
fi
echo "==> Build successful."
cd ..

# 2. Create remote directory if not exists
echo "[2/5] Ensuring remote directory exists..."
ssh -p $PORT $SSH_TARGET "mkdir -p $REMOTE_PATH"

# 3. Create systemd service file if not exists
echo "[3/5] Setting up systemd service..."
ssh -p $PORT $SSH_TARGET "
if [ ! -f $SERVICE_FILE ]; then
  echo 'Creating systemd service file...'
  cat > $SERVICE_FILE << 'EOF'
[Unit]
Description=Worker Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/root/service-go/worker-service
ExecStart=/root/service-go/worker-service/worker-service
Restart=on-failure
RestartSec=5s
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF
  systemctl daemon-reload
  systemctl enable $BINARY_NAME
  echo 'Service registered with systemd.'
else
  echo 'Service file already exists, skipping.'
fi
"

# 4. Stop service, upload binary, set permissions
echo "[4/5] Stopping service, uploading binary..."
ssh -p $PORT $SSH_TARGET "systemctl stop $BINARY_NAME 2>/dev/null || true"

scp -P $PORT $BINARY_NAME $SSH_TARGET:$REMOTE_PATH/
if [ $? -ne 0 ]; then
  echo "Error: Failed to upload binary!"
  exit 1
fi

ssh -p $PORT $SSH_TARGET "chmod +x $REMOTE_PATH/$BINARY_NAME"

# 5. Start service & verify
echo "[5/5] Starting service..."
ssh -p $PORT $SSH_TARGET "systemctl start $BINARY_NAME"

sleep 2

ssh -p $PORT $SSH_TARGET "systemctl status $BINARY_NAME --no-pager"

echo "================================================="
echo " DEPLOYMENT COMPLETED!"
echo " Binary: $REMOTE_PATH/$BINARY_NAME"
echo ""
echo " Useful commands:"
echo "   View logs : journalctl -u $BINARY_NAME -f"
echo "   Status    : systemctl status $BINARY_NAME"
echo "   Restart   : systemctl restart $BINARY_NAME"
echo "================================================="