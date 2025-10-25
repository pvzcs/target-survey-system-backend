#!/bin/bash

# Database backup script for Survey System
# Usage: ./scripts/backup.sh

set -e

# Configuration
BACKUP_DIR="${BACKUP_DIR:-./backups}"
DATE=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS="${RETENTION_DAYS:-7}"

# Database configuration from environment or defaults
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-3306}"
DB_USERNAME="${DB_USERNAME:-survey_user}"
DB_PASSWORD="${DB_PASSWORD:-survey_password}"
DB_DATABASE="${DB_DATABASE:-survey_system}"

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

# Backup filename
BACKUP_FILE="$BACKUP_DIR/survey_system_$DATE.sql.gz"

echo "Starting database backup..."
echo "Database: $DB_DATABASE"
echo "Host: $DB_HOST:$DB_PORT"
echo "Backup file: $BACKUP_FILE"

# Perform backup
if command -v mysqldump &> /dev/null; then
    mysqldump -h "$DB_HOST" \
              -P "$DB_PORT" \
              -u "$DB_USERNAME" \
              -p"$DB_PASSWORD" \
              --single-transaction \
              --routines \
              --triggers \
              --events \
              "$DB_DATABASE" | gzip > "$BACKUP_FILE"
    
    echo "Backup completed successfully!"
    echo "Backup size: $(du -h "$BACKUP_FILE" | cut -f1)"
else
    echo "Error: mysqldump command not found"
    exit 1
fi

# Clean up old backups
echo "Cleaning up backups older than $RETENTION_DAYS days..."
find "$BACKUP_DIR" -name "survey_system_*.sql.gz" -mtime +$RETENTION_DAYS -delete

echo "Backup process completed!"
