.PHONY: help build run test clean docker-build docker-up docker-down migrate backup restore

# Default target
help:
	@echo "Survey System - Available commands:"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-up     - Start all services with Docker Compose"
	@echo "  make docker-down   - Stop all services"
	@echo "  make migrate       - Run database migrations"
	@echo "  make backup        - Backup database"
	@echo "  make restore       - Restore database from backup"
	@echo "  make generate-key  - Generate encryption key"
	@echo "  make hash-password - Generate password hash"
	@echo "  make lint          - Run linters"
	@echo "  make fmt           - Format code"

# Build the application
build:
	@echo "Building application..."
	@go build -ldflags="-s -w" -o bin/survey-system ./cmd/server
	@echo "Build complete: bin/survey-system"

# Run the application
run:
	@echo "Starting application..."
	@go run ./cmd/server

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf dist/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t survey-system:latest .
	@echo "Docker image built: survey-system:latest"

# Start all services with Docker Compose
docker-up:
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d
	@echo "Services started. Check status with: docker-compose ps"

# Stop all services
docker-down:
	@echo "Stopping services..."
	@docker-compose down
	@echo "Services stopped"

# View logs
docker-logs:
	@docker-compose logs -f app

# Run database migrations
migrate:
	@echo "Running database migrations..."
	@mysql -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USERNAME) -p$(DB_PASSWORD) $(DB_DATABASE) < migrations/001_create_tables.sql
	@mysql -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USERNAME) -p$(DB_PASSWORD) $(DB_DATABASE) < migrations/002_seed_data.sql
	@echo "Migrations complete"

# Backup database
backup:
	@echo "Backing up database..."
	@./scripts/backup.sh

# Restore database
restore:
	@echo "Restoring database..."
	@./scripts/restore.sh $(BACKUP_FILE)

# Generate encryption key
generate-key:
	@echo "Generating encryption key..."
	@go run scripts/generate_key.go

# Generate password hash
hash-password:
	@echo "Generating password hash..."
	@go run scripts/hash_password.go $(PASSWORD)

# Run linters
lint:
	@echo "Running linters..."
	@go vet ./...
	@golangci-lint run || echo "golangci-lint not installed, skipping..."

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted"

# Install development dependencies
dev-deps:
	@echo "Installing development dependencies..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Development dependencies installed"

# Build for all platforms
build-all:
	@echo "Building for all platforms..."
	@mkdir -p dist
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/survey-system-linux-amd64 ./cmd/server
	@GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/survey-system-linux-arm64 ./cmd/server
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/survey-system-darwin-amd64 ./cmd/server
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/survey-system-darwin-arm64 ./cmd/server
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/survey-system-windows-amd64.exe ./cmd/server
	@GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o dist/survey-system-windows-arm64.exe ./cmd/server
	@echo "Build complete for all platforms"
