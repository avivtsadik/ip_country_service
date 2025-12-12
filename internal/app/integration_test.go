package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"ip_country_project/internal/datastores"
	"ip_country_project/internal/handlers"
	"ip_country_project/internal/middleware"
	"ip_country_project/internal/models"
	"ip_country_project/internal/services"
)

func setupTestHandler(t *testing.T) http.Handler {
	t.Helper()

	// Create test CSV data
	testData := "8.8.8.8,Mountain View,United States\n1.1.1.1,San Francisco,United States\n"

	tmpFile, err := os.CreateTemp("", "test_handlers_*.csv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	_, _ = tmpFile.WriteString(testData)
	_ = tmpFile.Close()

	// Use t.Cleanup instead of returning cleanup function
	t.Cleanup(func() {
		_ = os.Remove(tmpFile.Name())
	})

	// Use the app setup with test config
	application, err := NewWithTestConfig(tmpFile.Name(), 10) // 10 RPS for tests
	if err != nil {
		t.Fatalf("failed to create test application: %v", err)
	}

	return application.Handler
}

func TestIntegration_FindCountry_Success(t *testing.T) {
	handler := setupTestHandler(t)

	req := httptest.NewRequest("GET", "/v1/find-country?ip=8.8.8.8", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var location models.Location
	err := json.Unmarshal(rr.Body.Bytes(), &location)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if location.Country != "United States" {
		t.Errorf("expected country 'United States', got '%s'", location.Country)
	}
	if location.City != "Mountain View" {
		t.Errorf("expected city 'Mountain View', got '%s'", location.City)
	}
}

func TestIntegration_FindCountry_MissingIP(t *testing.T) {
	handler := setupTestHandler(t)

	req := httptest.NewRequest("GET", "/v1/find-country", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}

	var errorResp models.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &errorResp)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if errorResp.Error != "missing ip parameter" {
		t.Errorf("unexpected error message: %s", errorResp.Error)
	}
}

func TestIntegration_FindCountry_InvalidIP(t *testing.T) {
	handler := setupTestHandler(t)

	req := httptest.NewRequest("GET", "/v1/find-country?ip=invalid", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}

	var errorResp models.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &errorResp)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if errorResp.Error != "invalid IP address format" {
		t.Errorf("unexpected error message: %s", errorResp.Error)
	}
}

func TestIntegration_FindCountry_IPNotFound(t *testing.T) {
	handler := setupTestHandler(t)

	req := httptest.NewRequest("GET", "/v1/find-country?ip=1.2.3.4", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rr.Code)
	}

	var errorResp models.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &errorResp)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if errorResp.Error != "IP address not found" {
		t.Errorf("unexpected error message: %s", errorResp.Error)
	}
}

func TestIntegration_FindCountry_MethodNotAllowed(t *testing.T) {
	handler := setupTestHandler(t)

	req := httptest.NewRequest("POST", "/v1/find-country?ip=8.8.8.8", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rr.Code)
	}

	// Check Allow header
	allowHeader := rr.Header().Get("Allow")
	if allowHeader != "GET" {
		t.Errorf("expected Allow header 'GET', got '%s'", allowHeader)
	}
}

func TestIntegration_RateLimiting(t *testing.T) {
	// Create handler with very low rate limit for testing
	testData := "8.8.8.8,Mountain View,United States\n"
	tmpFile, err := os.CreateTemp("", "test_rate_limit_*.csv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	_, _ = tmpFile.WriteString(testData)
	_ = tmpFile.Close()
	t.Cleanup(func() {
		_ = os.Remove(tmpFile.Name())
	})

	datastore := datastores.NewCSVDataStore(tmpFile.Name())
	if loadErr := datastore.Load(); loadErr != nil {
		t.Fatalf("failed to load CSV: %v", loadErr)
	}

	service := services.NewLocationService(datastore)
	handlerObj := handlers.NewLocationHandler(service)
	rateLimiter := middleware.NewRateLimiter(1) // 1 RPS - very restrictive

	mux := http.NewServeMux()
	mux.Handle("/v1/find-country", rateLimiter.Middleware(http.HandlerFunc(handlerObj.FindCountry)))

	// First request should succeed
	req1 := httptest.NewRequest("GET", "/v1/find-country?ip=8.8.8.8", nil)
	rr1 := httptest.NewRecorder()
	mux.ServeHTTP(rr1, req1)

	if rr1.Code != http.StatusOK {
		t.Errorf("first request should succeed, got status %d", rr1.Code)
	}

	// Second immediate request should be rate limited
	req2 := httptest.NewRequest("GET", "/v1/find-country?ip=8.8.8.8", nil)
	rr2 := httptest.NewRecorder()
	mux.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusTooManyRequests {
		t.Errorf("second request should be rate limited, got status %d", rr2.Code)
	}

	var errorResp models.ErrorResponse
	if err := json.Unmarshal(rr2.Body.Bytes(), &errorResp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if errorResp.Error != "rate limit exceeded" {
		t.Errorf("unexpected error message: %s", errorResp.Error)
	}
}

func TestIntegration_ContentType(t *testing.T) {
	handler := setupTestHandler(t)

	req := httptest.NewRequest("GET", "/v1/find-country?ip=8.8.8.8", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", contentType)
	}
}
