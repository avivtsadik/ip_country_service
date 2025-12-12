package tests

import (
	"errors"
	"os"
	"testing"

	"ip_country_project/internal/datastores"
	appErrors "ip_country_project/internal/errors"
)

// Helper function to create test CSV file and datastore
func setupTestDatastore(t *testing.T, data string) *datastores.CSVDataStore {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "test_*.csv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	_, _ = tmpFile.WriteString(data)
	_ = tmpFile.Close()

	ds := datastores.NewCSVDataStore(tmpFile.Name())

	// Use t.Cleanup instead of returning cleanup function
	t.Cleanup(func() {
		_ = os.Remove(tmpFile.Name())
	})

	return ds
}

func TestCSVDataStore_Load(t *testing.T) {
	testData := "8.8.8.8,Mountain View,United States\n1.1.1.1,San Francisco,United States\n"
	ds := setupTestDatastore(t, testData)

	err := ds.Load()
	if err != nil {
		t.Fatalf("failed to load CSV: %v", err)
	}
}

func TestCSVDataStore_FindLocation_Success(t *testing.T) {
	testData := "8.8.8.8,Mountain View,United States\n1.1.1.1,San Francisco,United States\n"
	ds := setupTestDatastore(t, testData)

	if err := ds.Load(); err != nil {
		t.Fatalf("failed to load CSV: %v", err)
	}

	location, err := ds.FindLocation("8.8.8.8")
	if err != nil {
		t.Fatalf("expected successful lookup, got error: %v", err)
	}

	if location.Country != "United States" {
		t.Errorf("expected country 'United States', got '%s'", location.Country)
	}
	if location.City != "Mountain View" {
		t.Errorf("expected city 'Mountain View', got '%s'", location.City)
	}
}

func TestCSVDataStore_FindLocation_NotFound(t *testing.T) {
	testData := "8.8.8.8,Mountain View,United States\n"
	ds := setupTestDatastore(t, testData)

	if err := ds.Load(); err != nil {
		t.Fatalf("failed to load CSV: %v", err)
	}

	_, err := ds.FindLocation("1.2.3.4")
	if !errors.Is(err, appErrors.ErrIPNotFound) {
		t.Errorf("expected ErrIPNotFound, got %v", err)
	}
}

func TestCSVDataStore_FindLocation_InvalidIP(t *testing.T) {
	testData := "8.8.8.8,Mountain View,United States\n"
	ds := setupTestDatastore(t, testData)

	if err := ds.Load(); err != nil {
		t.Fatalf("failed to load CSV: %v", err)
	}

	_, err := ds.FindLocation("invalid-ip")
	if !errors.Is(err, appErrors.ErrInvalidIP) {
		t.Errorf("expected ErrInvalidIP, got %v", err)
	}
}

func TestCSVDataStore_Load_InvalidFile(t *testing.T) {
	ds := datastores.NewCSVDataStore("nonexistent.csv")
	err := ds.Load()
	if err == nil {
		t.Error("expected error loading nonexistent file")
	}
}

func TestCSVDataStore_Load_MalformedCSV(t *testing.T) {
	testData := "8.8.8.8,Mountain View\n1.1.1.1,San Francisco,United States,Extra\n"
	ds := setupTestDatastore(t, testData)

	err := ds.Load()
	if err == nil {
		t.Error("expected error loading malformed CSV")
	}
}

func TestCSVDataStore_Load_InvalidIPInCSV(t *testing.T) {
	testData := "invalid-ip,Mountain View,United States\n"
	ds := setupTestDatastore(t, testData)

	err := ds.Load()
	if err == nil {
		t.Error("expected error loading CSV with invalid IP")
	}
}
