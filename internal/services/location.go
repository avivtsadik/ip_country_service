package services

import (
	"context"
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

func (s *LocationService) FindCountry(ctx context.Context, ip string) (*models.Location, error) {
	// Normalize and validate IP format
	normalizedIP := utils.NormalizeIP(ip)
	if normalizedIP == "" {
		return nil, errors.ErrInvalidIP
	}

	// Delegate to datastore with context
	location, err := s.datastore.FindLocation(ctx, normalizedIP)
	if err != nil {
		return nil, err
	}

	return location, nil
}
