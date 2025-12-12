package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ip_country_project/internal/app"
	"ip_country_project/internal/config"
	safe "ip_country_project/internal/utils"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Printf("IP Country Service starting on port %s\n", cfg.Port)
	fmt.Printf("Rate limit: %.1f RPS\n", cfg.RateLimitRPS)
	fmt.Printf("Datastore: %s (%s)\n", cfg.DatastoreType, cfg.DatastoreFile)

	// Initialize application
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Create server with production timeouts
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      application.Handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("Server starting on port :%s\n", cfg.Port)
	fmt.Println("Endpoints available:")
	fmt.Println("  GET /v1/find-country?ip=8.8.8.8")
	fmt.Println("  GET /health")

	// Start server in safe goroutine
	safe.Go(func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed to start: %v", err)
		}
	})

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
