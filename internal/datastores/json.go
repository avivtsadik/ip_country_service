package datastores

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"ip_country_project/internal/errors"
	"ip_country_project/internal/models"
	"ip_country_project/internal/utils"
)

type JSONDataStore struct {
	filePath string
	data     map[string]*models.Location
	mutex    sync.RWMutex
}

func NewJSONDataStore(filePath string) *JSONDataStore {
	return &JSONDataStore{
		filePath: filePath,
		data:     make(map[string]*models.Location),
	}
}

func (j *JSONDataStore) Load(ctx context.Context) error {
	file, err := os.Open(j.filePath)
	if err != nil {
		return fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	var locations []models.Location
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&locations); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	j.mutex.Lock()
	defer j.mutex.Unlock()

	for i, location := range locations {
		if !utils.IsValidIP(location.IP) {
			return fmt.Errorf("invalid IP address at index %d: %s", i, location.IP)
		}

		j.data[location.IP] = &models.Location{
			IP:      location.IP,
			City:    location.City,
			Country: location.Country,
		}
	}

	return nil
}

func (j *JSONDataStore) FindLocation(ctx context.Context, ip string) (*models.Location, error) {
	if !utils.IsValidIP(ip) {
		return nil, errors.ErrInvalidIP
	}

	j.mutex.RLock()
	defer j.mutex.RUnlock()

	location, exists := j.data[ip]
	if !exists {
		return nil, errors.ErrIPNotFound
	}

	// Return a copy to avoid race conditions
	return &models.Location{
		IP:      location.IP,
		City:    location.City,
		Country: location.Country,
	}, nil
}

func (j *JSONDataStore) Close() error {
	return nil
}