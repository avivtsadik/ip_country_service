package utils

import (
	"net"
	"strings"
)

// NormalizeIP cleans and validates IP address input
func NormalizeIP(ip string) string {
	trimmed := strings.TrimSpace(ip)
	parsed := net.ParseIP(trimmed)
	if parsed == nil {
		return ""
	}
	return parsed.String()
}

// IsValidIP checks if the IP address is valid
func IsValidIP(ip string) bool {
	return NormalizeIP(ip) != ""
}