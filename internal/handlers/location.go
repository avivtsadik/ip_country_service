package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	appErrors "ip_country_project/internal/errors"
	"ip_country_project/internal/models"
	"ip_country_project/internal/services"
)

type LocationHandler struct {
	service *services.LocationService
}

func NewLocationHandler(service *services.LocationService) *LocationHandler {
	return &LocationHandler{
		service: service,
	}
}

func (h *LocationHandler) FindCountry(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		h.writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get IP parameter from query string
	ip := r.URL.Query().Get("ip")
	if ip == "" {
		h.writeError(w, "missing ip parameter", http.StatusBadRequest)
		return
	}

	// Call service with request context
	location, err := h.service.FindCountry(r.Context(), ip)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(location)
}

func (h *LocationHandler) handleServiceError(w http.ResponseWriter, err error) {
	if errors.Is(err, appErrors.ErrInvalidIP) {
		h.writeError(w, "invalid IP address format", http.StatusBadRequest)
		return
	}
	
	if errors.Is(err, appErrors.ErrIPNotFound) {
		h.writeError(w, "IP address not found", http.StatusNotFound)
		return
	}

	// All other errors are internal server errors
	h.writeError(w, "internal server error", http.StatusInternalServerError)
}

func (h *LocationHandler) writeError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{Error: message})
}