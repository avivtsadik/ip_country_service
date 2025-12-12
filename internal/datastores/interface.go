package datastores

import "ip_country_project/internal/models"

type DataStore interface {
	FindLocation(ip string) (*models.Location, error)
	Load() error
	// Close Close() included for future HTTP clients and Redis connections
	// CSV implementation returns nil - no persistent resources to clean
	Close() error
}
