package services

import (
	"fmt"

	"ip_country_project/internal/datastores"
	"ip_country_project/internal/errors"
	"ip_country_project/internal/models"
	"ip_country_project/internal/utils"
)

// LocationService provides business logic for IP location lookups
type LocationService struct {
	datastore datastores.DataStore
}

func NewLocationService(datastore datastores.DataStore) *LocationService {
	return &LocationService{
		datastore: datastore,
	}
}

func (s *LocationService) FindCountry(ip string) (*models.Location, error) {
	// Normalize and validate IP format
	normalizedIP := utils.NormalizeIP(ip)
	if normalizedIP == "" {
		return nil, errors.ErrInvalidIP
	}

	// Delegate to datastore
	location, err := s.datastore.FindLocation(normalizedIP)
	if err != nil {
		// Wrap datastore errors with context for better logging
		return nil, fmt.Errorf("datastore lookup failed: %w", err)
	}

	return location, nil
}