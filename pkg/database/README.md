# Database Package

This package provides database connection management and migration utilities for the Survey System.

## Features

- MySQL database connection using GORM
- Connection pool configuration
- Database health checks
- Automatic migration support
- Manual SQL migration scripts

## Usage

### Initialize Database Connection

```go
import (
    "survey-system/internal/config"
    "survey-system/pkg/database"
)

// Load configuration
cfg, err := config.Load("./config/config.yaml")
if err != nil {
    log.Fatal(err)
}

// Initialize database
db, err := database.InitDB(&cfg.Database)
if err != nil {
    log.Fatal(err)
}
defer database.Close()
```

### Run Auto-Migration

```go
// Run GORM auto-migration for all models
if err := database.AutoMigrate(db); err != nil {
    log.Fatal(err)
}
```

### Health Check

```go
// Check database connection health
if err := database.HealthCheck(); err != nil {
    log.Printf("Database health check failed: %v", err)
}
```

## Configuration

Database configuration is managed through the config package. Required settings:

```yaml
database:
  host: localhost
  port: 3306
  username: root
  password: your_password
  database: survey_system
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 1h
```

## Migration Scripts

Manual SQL migration scripts are located in the `migrations/` directory:

- `001_create_tables.sql` - Creates all database tables with proper indexes and foreign keys

These scripts can be run manually or through a migration tool if needed.

## Models

The following models are supported:

- **User** - System users (admins)
- **Survey** - Survey/questionnaire definitions
- **Question** - Questions within surveys
- **Response** - Survey responses/submissions
- **OneLink** - One-time access links

All models are defined in the `internal/model` package.
