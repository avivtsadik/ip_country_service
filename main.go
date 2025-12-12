package main

import (
	"fmt"
	"log"

	"ip_country_project/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Printf("IP Country Service starting on port %s\n", cfg.Port)
	fmt.Printf("Rate limit: %.1f RPS\n", cfg.RateLimitRPS)
	fmt.Printf("Datastore: %s (%s)\n", cfg.DatastoreType, cfg.DatastoreFile)
	fmt.Printf("Cache: enabled=%t, TTL=%v, max_size=%d\n", 
		cfg.CacheEnabled, cfg.CacheTTL, cfg.CacheMaxSize)

	// TODO: Initialize datastore, rate limiter, and HTTP server
	fmt.Println("Configuration loaded successfully!")
}
