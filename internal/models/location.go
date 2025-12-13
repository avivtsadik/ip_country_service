package models

type Location struct {
	IP      string `json:"ip,omitempty"`
	Country string `json:"country"`
	City    string `json:"city"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}