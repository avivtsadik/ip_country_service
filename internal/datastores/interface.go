package datastores

import (
	"context"
	"ip_country_project/internal/models"
)

type DataStore interface {
	FindLocation(ctx context.Context, ip string) (*models.Location, error)
	Load(ctx context.Context) error
	// Close Close() included for future HTTP clients and Redis connections
	// CSV implementation returns nil - no persistent resources to clean
	Close() error
}
