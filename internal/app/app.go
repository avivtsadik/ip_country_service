package app

import (
	"context"
	"fmt"
	"net/http"

	"ip_country_project/internal/config"
	"ip_country_project/internal/datastores"
	"ip_country_project/internal/handlers"
	"ip_country_project/internal/middleware"
	"ip_country_project/internal/services"
)

// Application holds the application dependencies
type Application struct {
	Handler     http.Handler
	Config      *config.Config
	DataStore   datastores.DataStore
	Service     *services.LocationService
	HTTPHandler *handlers.LocationHandler
}

// New creates a new Application with all dependencies initialized
func New(cfg *config.Config) (*Application, error) {
	// Initialize datastore based on type
	var datastore datastores.DataStore
	switch cfg.DatastoreType {
	case "csv":
		datastore = datastores.NewCSVDataStore(cfg.DatastoreFile)
	case "json":
		datastore = datastores.NewJSONDataStore(cfg.DatastoreFile)
	default:
		return nil, fmt.Errorf("unsupported datastore type: %s", cfg.DatastoreType)
	}

	if err := datastore.Load(context.Background()); err != nil {
		return nil, err
	}

	// Initialize service layer
	service := services.NewLocationService(datastore)

	// Initialize HTTP handler
	httpHandler := handlers.NewLocationHandler(service)

	// Initialize rate limiter
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimitRPS)

	// Setup routes
	mux := http.NewServeMux()
	
	// API v1 endpoints
	mux.Handle("/v1/find-country", rateLimiter.Middleware(http.HandlerFunc(httpHandler.FindCountry)))
	
	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	app := &Application{
		Handler:     mux,
		Config:      cfg,
		DataStore:   datastore,
		Service:     service,
		HTTPHandler: httpHandler,
	}

	return app, nil
}

// NewWithTestConfig creates an application with test-specific configuration
func NewWithTestConfig(datastoreFile string, rateLimitRPS float64) (*Application, error) {
	cfg := &config.Config{
		Port:          "8080", // Default test port
		RateLimitRPS:  rateLimitRPS,
		DatastoreType: "csv",
		DatastoreFile: datastoreFile,
	}

	return New(cfg)
}