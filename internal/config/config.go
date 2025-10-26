package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Encryption EncryptionConfig `mapstructure:"encryption"`
	CORS       CORSConfig       `mapstructure:"cors"`
	OneLink    OneLinkConfig    `mapstructure:"onelink"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	Expiration time.Duration `mapstructure:"expiration"`
}

// EncryptionConfig holds encryption configuration
type EncryptionConfig struct {
	Key string `mapstructure:"key"`
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
}

// OneLinkConfig holds one-time link configuration
type OneLinkConfig struct {
	BaseURL           string        `mapstructure:"base_url"`
	DefaultExpiration time.Duration `mapstructure:"default_expiration"`
	MaxExpiration     time.Duration `mapstructure:"max_expiration"`
}

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set config file path
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath("./config")
		v.AddConfigPath(".")
	}

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Enable environment variable override
	v.AutomaticEnv()
	v.SetEnvPrefix("SURVEY")

	// Map environment variables to config fields
	// Database
	v.BindEnv("database.host", "DB_HOST")
	v.BindEnv("database.port", "DB_PORT")
	v.BindEnv("database.username", "DB_USERNAME")
	v.BindEnv("database.password", "DB_PASSWORD")
	v.BindEnv("database.database", "DB_DATABASE")

	// Redis
	v.BindEnv("redis.host", "REDIS_HOST")
	v.BindEnv("redis.port", "REDIS_PORT")
	v.BindEnv("redis.password", "REDIS_PASSWORD")

	// JWT
	v.BindEnv("jwt.secret", "JWT_SECRET")

	// Encryption
	v.BindEnv("encryption.key", "ENCRYPTION_KEY")

	// Server
	v.BindEnv("server.port", "SERVER_PORT")
	v.BindEnv("server.mode", "SERVER_MODE")

	// Unmarshal config into struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := validate(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// validate validates the configuration
func validate(config *Config) error {
	// Validate encryption key length (must be 32 bytes for AES-256)
	if len(config.Encryption.Key) != 32 {
		return fmt.Errorf("encryption key must be exactly 32 bytes, got %d bytes", len(config.Encryption.Key))
	}

	// Validate JWT secret is not empty
	if config.JWT.Secret == "" {
		return fmt.Errorf("JWT secret cannot be empty")
	}

	// Validate database configuration
	if config.Database.Host == "" {
		return fmt.Errorf("database host cannot be empty")
	}
	if config.Database.Database == "" {
		return fmt.Errorf("database name cannot be empty")
	}

	// Validate Redis configuration
	if config.Redis.Host == "" {
		return fmt.Errorf("redis host cannot be empty")
	}

	// Validate server port
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	return nil
}
