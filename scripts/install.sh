#!/bin/bash

# Installation script for Survey System
# This script installs the Survey System as a systemd service

set -e

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "Please run as root (use sudo)"
    exit 1
fi

# Configuration
INSTALL_DIR="/opt/survey-system"
SERVICE_USER="www-data"
SERVICE_GROUP="www-data"

echo "Survey System Installation Script"
echo "=================================="
echo ""

# Check if binary exists
if [ ! -f "./survey-system" ]; then
    echo "Error: survey-system binary not found in current directory"
    echo "Please build the application first: make build"
    exit 1
fi

# Create installation directory
echo "Creating installation directory: $INSTALL_DIR"
mkdir -p "$INSTALL_DIR"

# Copy binary
echo "Copying binary..."
cp ./survey-system "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/survey-system"

# Copy configuration files
echo "Copying configuration files..."
mkdir -p "$INSTALL_DIR/config"
mkdir -p "$INSTALL_DIR/migrations"
cp -r config/* "$INSTALL_DIR/config/" || true
cp -r migrations/* "$INSTALL_DIR/migrations/" || true

# Copy .env.example if .env doesn't exist
if [ ! -f "$INSTALL_DIR/.env" ]; then
    echo "Creating .env file from template..."
    cp .env.example "$INSTALL_DIR/.env"
    echo "IMPORTANT: Edit $INSTALL_DIR/.env with your configuration"
fi

# Set ownership
echo "Setting ownership..."
chown -R $SERVICE_USER:$SERVICE_GROUP "$INSTALL_DIR"

# Install systemd service
echo "Installing systemd service..."
cp scripts/survey-system.service /etc/systemd/system/
systemctl daemon-reload

echo ""
echo "Installation complete!"
echo ""
echo "Next steps:"
echo "1. Edit configuration: nano $INSTALL_DIR/.env"
echo "2. Run database migrations: cd $INSTALL_DIR && mysql -u user -p database < migrations/001_create_tables.sql"
echo "3. Enable service: systemctl enable survey-system"
echo "4. Start service: systemctl start survey-system"
echo "5. Check status: systemctl status survey-system"
echo "6. View logs: journalctl -u survey-system -f"
echo ""
