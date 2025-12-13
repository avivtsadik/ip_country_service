package datastores

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"sync"

	"ip_country_project/internal/errors"
	"ip_country_project/internal/models"
	"ip_country_project/internal/utils"
)

type CSVDataStore struct {
	filePath string
	data     map[string]*models.Location
	mutex    sync.RWMutex
}

func NewCSVDataStore(filePath string) *CSVDataStore {
	return &CSVDataStore{
		filePath: filePath,
		data:     make(map[string]*models.Location),
	}
}

func (c *CSVDataStore) Load(ctx context.Context) error {
	file, err := os.Open(c.filePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// CSV is loaded fully into memory at startup.
	// This is acceptable for the exercise and small datasets.
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i, record := range records {
		if len(record) != 3 {
			return fmt.Errorf("invalid CSV format at line %d: expected 3 fields, got %d", i+1, len(record))
		}

		ip, city, country := record[0], record[1], record[2]

		if !utils.IsValidIP(ip) {
			return fmt.Errorf("invalid IP address at line %d: %s", i+1, ip)
		}

		c.data[ip] = &models.Location{
			IP:      ip,
			City:    city,
			Country: country,
		}
	}

	return nil
}

func (c *CSVDataStore) FindLocation(ctx context.Context, ip string) (*models.Location, error) {
	if !utils.IsValidIP(ip) {
		return nil, errors.ErrInvalidIP
	}

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	location, exists := c.data[ip]
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

func (c *CSVDataStore) Close() error {
	return nil
}
