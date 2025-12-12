package services

import (
	"errors"
	"testing"

	appErrors "ip_country_project/internal/errors"
	"ip_country_project/internal/models"
)

// mockDataStore implements datastore's.DataStore interface for testing
type mockDataStore struct {
	findLocationFunc func(ip string) (*models.Location, error)
	loadFunc         func() error
	closeFunc        func() error
}

func (m *mockDataStore) FindLocation(ip string) (*models.Location, error) {
	if m.findLocationFunc != nil {
		return m.findLocationFunc(ip)
	}
	return nil, appErrors.ErrIPNotFound
}

func (m *mockDataStore) Load() error {
	if m.loadFunc != nil {
		return m.loadFunc()
	}
	return nil
}

func (m *mockDataStore) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}

func TestLocationService_FindCountry_Success(t *testing.T) {
	expectedLocation := &models.Location{
		IP:      "8.8.8.8",
		Country: "United States",
		City:    "Mountain View",
	}

	mockDS := &mockDataStore{
		findLocationFunc: func(ip string) (*models.Location, error) {
			if ip == "8.8.8.8" {
				return expectedLocation, nil
			}
			return nil, appErrors.ErrIPNotFound
		},
	}

	service := NewLocationService(mockDS)

	location, err := service.FindCountry("8.8.8.8")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if location.Country != expectedLocation.Country {
		t.Errorf("expected country '%s', got '%s'", expectedLocation.Country, location.Country)
	}
	if location.City != expectedLocation.City {
		t.Errorf("expected city '%s', got '%s'", expectedLocation.City, location.City)
	}
}

func TestLocationService_FindCountry_InvalidIP(t *testing.T) {
	mockDS := &mockDataStore{}
	service := NewLocationService(mockDS)

	testCases := []string{
		"",
		"invalid-ip",
		"256.256.256.256",
		"abc.def.ghi.jkl",
		"192.168.1",
		"192.168.1.1.1",
	}

	for _, invalidIP := range testCases {
		t.Run("invalid_ip_"+invalidIP, func(t *testing.T) {
			_, err := service.FindCountry(invalidIP)
			if !errors.Is(err, appErrors.ErrInvalidIP) {
				t.Errorf("expected ErrInvalidIP for '%s', got: %v", invalidIP, err)
			}
		})
	}
}

func TestLocationService_FindCountry_IPNotFound(t *testing.T) {
	mockDS := &mockDataStore{
		findLocationFunc: func(ip string) (*models.Location, error) {
			return nil, appErrors.ErrIPNotFound
		},
	}

	service := NewLocationService(mockDS)

	_, err := service.FindCountry("1.2.3.4")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// The service wraps datastore errors, so we need to unwrap
	if !errors.Is(err, appErrors.ErrIPNotFound) {
		t.Errorf("expected wrapped ErrIPNotFound, got: %v", err)
	}
}

func TestLocationService_FindCountry_DatastoreError(t *testing.T) {
	datastoreErr := errors.New("database connection failed")

	mockDS := &mockDataStore{
		findLocationFunc: func(ip string) (*models.Location, error) {
			return nil, datastoreErr
		},
	}

	service := NewLocationService(mockDS)

	_, err := service.FindCountry("8.8.8.8")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Should wrap the original error
	if !errors.Is(err, datastoreErr) {
		t.Errorf("expected wrapped datastore error, got: %v", err)
	}
}

func TestLocationService_FindCountry_IPNormalization(t *testing.T) {
	var receivedIP string

	mockDS := &mockDataStore{
		findLocationFunc: func(ip string) (*models.Location, error) {
			receivedIP = ip
			return &models.Location{
				IP:      ip,
				Country: "Test",
				City:    "Test",
			}, nil
		},
	}

	service := NewLocationService(mockDS)

	// Test whitespace trimming
	_, err := service.FindCountry("  8.8.8.8  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedIP != "8.8.8.8" {
		t.Errorf("expected normalized IP '8.8.8.8', but datastore received '%s'", receivedIP)
	}
}

func TestNewLocationService(t *testing.T) {
	mockDS := &mockDataStore{}
	service := NewLocationService(mockDS)

	if service == nil {
		t.Fatal("NewLocationService returned nil")
	}

	if service.datastore != mockDS {
		t.Error("NewLocationService didn't set datastore correctly")
	}
}
