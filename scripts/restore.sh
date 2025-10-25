#!/bin/bash

# Database restore script for Survey System
# Usage: ./scripts/restore.sh <backup_file>

set -e

# Check if backup file is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <backup_file>"
    echo "Example: $0 ./backups/survey_system_20240101_120000.sql.gz"
    exit 1
fi

BACKUP_FILE="$1"

# Check if backup file exists
if [ ! -f "$BACKUP_FILE" ]; then
    echo "Error: Backup file not found: $BACKUP_FILE"
    exit 1
fi

# Database configuration from environment or defaults
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-3306}"
DB_USERNAME="${DB_USERNAME:-survey_user}"
DB_PASSWORD="${DB_PASSWORD:-survey_password}"
DB_DATABASE="${DB_DATABASE:-survey_system}"

echo "WARNING: This will restore the database from backup."
echo "Database: $DB_DATABASE"
echo "Host: $DB_HOST:$DB_PORT"
echo "Backup file: $BACKUP_FILE"
echo ""
read -p "Are you sure you want to continue? (yes/no): " CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    echo "Restore cancelled."
    exit 0
fi

echo "Starting database restore..."

# Restore backup
if command -v mysql &> /dev/null; then
    if [[ "$BACKUP_FILE" == *.gz ]]; then
        gunzip < "$BACKUP_FILE" | mysql -h "$DB_HOST" \
                                         -P "$DB_PORT" \
                                         -u "$DB_USERNAME" \
                                         -p"$DB_PASSWORD" \
                                         "$DB_DATABASE"
    else
        mysql -h "$DB_HOST" \
              -P "$DB_PORT" \
              -u "$DB_USERNAME" \
              -p"$DB_PASSWORD" \
              "$DB_DATABASE" < "$BACKUP_FILE"
    fi
    
    echo "Restore completed successfully!"
else
    echo "Error: mysql command not found"
    exit 1
fi
