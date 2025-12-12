package main

import (
	"fmt"
	"log"
	"net/http"

	"ip_country_project/internal/config"
	"ip_country_project/internal/datastores"
	"ip_country_project/internal/handlers"
	"ip_country_project/internal/middleware"
	"ip_country_project/internal/services"
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

	// Initialize datastore
	datastore := datastores.NewCSVDataStore(cfg.DatastoreFile)
	if err := datastore.Load(); err != nil {
		log.Fatalf("Failed to load datastore: %v", err)
	}

	// Initialize service layer
	service := services.NewLocationService(datastore)

	// Initialize HTTP components
	handler := handlers.NewLocationHandler(service)
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimitRPS)

	// Setup routes
	mux := http.NewServeMux()
	mux.Handle("/v1/find-country", rateLimiter.Middleware(http.HandlerFunc(handler.FindCountry)))
	
	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	fmt.Printf("Server starting on port :%s\n", cfg.Port)
	fmt.Println("Endpoints available:")
	fmt.Println("  GET /v1/find-country?ip=8.8.8.8")
	fmt.Println("  GET /health")
	
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
