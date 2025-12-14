package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Host          string
	Port          string
	RateLimitRPS  float64
	DatastoreType string
	DatastoreFile string
}

func Load() (*Config, error) {
	config := &Config{
		Host:          getEnv("HOST", "localhost"),
		Port:          getEnv("PORT", "8080"),
		DatastoreType: getEnv("DATASTORE_TYPE", "csv"),
		DatastoreFile: getEnv("DATASTORE_FILE", "testdata/sample_ips.csv"),
	}

	var err error
	config.RateLimitRPS, err = getEnvFloat("RATE_LIMIT_RPS", 10.0)
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT_RPS: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.RateLimitRPS <= 0 {
		return fmt.Errorf("RATE_LIMIT_RPS must be positive, got: %f", c.RateLimitRPS)
	}
	if c.DatastoreType != "csv" && c.DatastoreType != "json" {
		return fmt.Errorf("unsupported DATASTORE_TYPE: %s (supported: csv, json)", c.DatastoreType)
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) (float64, error) {
	if value := os.Getenv(key); value != "" {
		return strconv.ParseFloat(value, 64)
	}
	return defaultValue, nil
}
