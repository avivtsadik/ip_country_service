package datastores

import (
	"context"
	"errors"
	"os"
	"testing"

	appErrors "ip_country_project/internal/errors"
)

// Helper function to create test JSON file and datastore
func setupTestJSONDatastore(t *testing.T, data string) *JSONDataStore {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "test_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	_, _ = tmpFile.WriteString(data)
	_ = tmpFile.Close()

	ds := NewJSONDataStore(tmpFile.Name())

	// Use t.Cleanup instead of returning cleanup function
	t.Cleanup(func() {
		_ = os.Remove(tmpFile.Name())
	})

	return ds
}

func TestJSONDataStore_Load(t *testing.T) {
	testData := `[
		{"ip": "8.8.8.8", "city": "Mountain View", "country": "United States"},
		{"ip": "1.1.1.1", "city": "San Francisco", "country": "United States"}
	]`
	ds := setupTestJSONDatastore(t, testData)

	err := ds.Load(context.Background())
	if err != nil {
		t.Fatalf("failed to load JSON: %v", err)
	}
}

func TestJSONDataStore_FindLocation_Success(t *testing.T) {
	testData := `[
		{"ip": "8.8.8.8", "city": "Mountain View", "country": "United States"},
		{"ip": "1.1.1.1", "city": "San Francisco", "country": "United States"}
	]`
	ds := setupTestJSONDatastore(t, testData)

	if err := ds.Load(context.Background()); err != nil {
		t.Fatalf("failed to load JSON: %v", err)
	}

	location, err := ds.FindLocation(context.Background(), "8.8.8.8")
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

func TestJSONDataStore_FindLocation_NotFound(t *testing.T) {
	testData := `[
		{"ip": "8.8.8.8", "city": "Mountain View", "country": "United States"}
	]`
	ds := setupTestJSONDatastore(t, testData)

	if err := ds.Load(context.Background()); err != nil {
		t.Fatalf("failed to load JSON: %v", err)
	}

	_, err := ds.FindLocation(context.Background(), "1.2.3.4")
	if !errors.Is(err, appErrors.ErrIPNotFound) {
		t.Errorf("expected ErrIPNotFound, got %v", err)
	}
}

func TestJSONDataStore_FindLocation_InvalidIP(t *testing.T) {
	testData := `[
		{"ip": "8.8.8.8", "city": "Mountain View", "country": "United States"}
	]`
	ds := setupTestJSONDatastore(t, testData)

	if err := ds.Load(context.Background()); err != nil {
		t.Fatalf("failed to load JSON: %v", err)
	}

	_, err := ds.FindLocation(context.Background(), "invalid-ip")
	if !errors.Is(err, appErrors.ErrInvalidIP) {
		t.Errorf("expected ErrInvalidIP, got %v", err)
	}
}

func TestJSONDataStore_Load_InvalidFile(t *testing.T) {
	ds := NewJSONDataStore("nonexistent.json")
	err := ds.Load(context.Background())
	if err == nil {
		t.Error("expected error loading nonexistent file")
	}
}

func TestJSONDataStore_Load_MalformedJSON(t *testing.T) {
	testData := `[
		{"ip": "8.8.8.8", "city": "Mountain View", "country": "United States"},
		{"ip": "1.1.1.1", "city": "San Francisco"
	]`
	ds := setupTestJSONDatastore(t, testData)

	err := ds.Load(context.Background())
	if err == nil {
		t.Error("expected error loading malformed JSON")
	}
}

func TestJSONDataStore_Load_InvalidIPInJSON(t *testing.T) {
	testData := `[
		{"ip": "invalid-ip", "city": "Mountain View", "country": "United States"}
	]`
	ds := setupTestJSONDatastore(t, testData)

	err := ds.Load(context.Background())
	if err == nil {
		t.Error("expected error loading JSON with invalid IP")
	}
}

func TestJSONDataStore_Load_MissingFields(t *testing.T) {
	testData := `[
		{"ip": "8.8.8.8", "city": "Mountain View"}
	]`
	ds := setupTestJSONDatastore(t, testData)

	err := ds.Load(context.Background())
	if err != nil {
		t.Fatalf("failed to load JSON with missing field: %v", err)
	}

	location, err := ds.FindLocation(context.Background(), "8.8.8.8")
	if err != nil {
		t.Fatalf("expected successful lookup, got error: %v", err)
	}

	if location.Country != "" {
		t.Errorf("expected empty country, got '%s'", location.Country)
	}
	if location.City != "Mountain View" {
		t.Errorf("expected city 'Mountain View', got '%s'", location.City)
	}
}