package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port           string
	RateLimitRPS   float64
	DatastoreType  string
	DatastoreFile  string
	CacheEnabled   bool
	CacheTTL       time.Duration
	CacheMaxSize   int
}

func Load() (*Config, error) {
	config := &Config{
		Port:           getEnv("PORT", "8080"),
		DatastoreType:  getEnv("DATASTORE_TYPE", "csv"),
		DatastoreFile:  getEnv("DATASTORE_FILE", "testdata/sample_ips.csv"),
		CacheEnabled:   getEnvBool("CACHE_ENABLED", true),
		CacheMaxSize:   getEnvInt("CACHE_MAX_SIZE", 10000),
	}

	var err error
	config.RateLimitRPS, err = getEnvFloat("RATE_LIMIT_RPS", 10.0)
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT_RPS: %w", err)
	}

	cacheTTLMinutes, err := getEnvFloat("CACHE_TTL_MINUTES", 60.0)
	if err != nil {
		return nil, fmt.Errorf("invalid CACHE_TTL_MINUTES: %w", err)
	}
	config.CacheTTL = time.Duration(cacheTTLMinutes) * time.Minute

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.RateLimitRPS <= 0 {
		return fmt.Errorf("RATE_LIMIT_RPS must be positive, got: %f", c.RateLimitRPS)
	}
	if c.DatastoreType != "csv" {
		return fmt.Errorf("unsupported DATASTORE_TYPE: %s", c.DatastoreType)
	}
	if c.CacheMaxSize <= 0 {
		return fmt.Errorf("CACHE_MAX_SIZE must be positive, got: %d", c.CacheMaxSize)
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) (float64, error) {
	if value := os.Getenv(key); value != "" {
		return strconv.ParseFloat(value, 64)
	}
	return defaultValue, nil
}