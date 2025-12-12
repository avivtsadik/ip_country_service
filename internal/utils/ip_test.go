package utils

import "testing"

func TestNormalizeIP(t *testing.T) {
	// Valid cases
	if result := NormalizeIP("8.8.8.8"); result != "8.8.8.8" {
		t.Errorf("expected '8.8.8.8', got '%s'", result)
	}
	
	// Whitespace trimming
	if result := NormalizeIP("  192.168.1.1  "); result != "192.168.1.1" {
		t.Errorf("expected '192.168.1.1', got '%s'", result)
	}
	
	// Invalid cases should return empty (including leading zeros which Go considers invalid)
	invalid := []string{"", "invalid", "256.256.256.256", "192.168.1", "008.008.008.008"}
	for _, ip := range invalid {
		if result := NormalizeIP(ip); result != "" {
			t.Errorf("expected empty for invalid IP '%s', got '%s'", ip, result)
		}
	}
}

func TestIsValidIP(t *testing.T) {
	// Valid IPs
	valid := []string{"8.8.8.8", "192.168.1.1", "127.0.0.1"}
	for _, ip := range valid {
		if !IsValidIP(ip) {
			t.Errorf("expected '%s' to be valid", ip)
		}
	}
	
	// Invalid IPs
	invalid := []string{"", "invalid", "256.256.256.256", "192.168.1"}
	for _, ip := range invalid {
		if IsValidIP(ip) {
			t.Errorf("expected '%s' to be invalid", ip)
		}
	}
}