package main

import (
	"fmt"
	"log"

	"ip_country_project/internal/config"
	"ip_country_project/internal/datastores"
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

	// Test datastore loading
	datastore := datastores.NewCSVDataStore(cfg.DatastoreFile)
	if err := datastore.Load(); err != nil {
		log.Fatalf("Failed to load datastore: %v", err)
	}
	
	// Test successful lookup
	location, err := datastore.FindLocation("8.8.8.8")
	if err != nil {
		fmt.Printf("Error finding IP 8.8.8.8: %v\n", err)
	} else {
		fmt.Printf("Test lookup 8.8.8.8: %s, %s\n", location.City, location.Country)
	}
	
	// Test IP not found
	location, err = datastore.FindLocation("1.2.3.4")
	if err != nil {
		fmt.Printf("Error finding IP 1.2.3.4: %v\n", err)
	} else {
		fmt.Printf("Test lookup 1.2.3.4: %s, %s\n", location.City, location.Country)
	}
	
	// Test invalid IP format
	location, err = datastore.FindLocation("invalid-ip")
	if err != nil {
		fmt.Printf("Error finding invalid IP: %v\n", err)
	}
	
	fmt.Println("Configuration and datastore loaded successfully!")
}
