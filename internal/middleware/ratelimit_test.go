package middleware

import (
	"testing"
	"time"
)

func TestRateLimiter_Allow(t *testing.T) {
	// Test basic allow functionality
	rl := NewRateLimiter(2) // 2 RPS

	// Should allow first 2 requests immediately
	if !rl.Allow() {
		t.Fatal("first request should be allowed")
	}
	if !rl.Allow() {
		t.Fatal("second request should be allowed")
	}

	// Third request should be denied (bucket empty)
	if rl.Allow() {
		t.Fatal("third request should be denied - rate limited")
	}
}

func TestRateLimiter_TokenRefill(t *testing.T) {
	rl := NewRateLimiter(2) // 2 RPS

	// Exhaust tokens
	rl.Allow()
	rl.Allow()
	
	// Should be rate limited
	if rl.Allow() {
		t.Fatal("should be rate limited")
	}

	// Wait for 1 second for token refill
	time.Sleep(1100 * time.Millisecond) // Slightly over 1 second

	// Should allow requests again after refill
	if !rl.Allow() {
		t.Fatal("should allow request after token refill")
	}
}

func TestRateLimiter_BurstHandling(t *testing.T) {
	rl := NewRateLimiter(5) // 5 RPS
	
	// Should handle initial burst up to capacity
	allowed := 0
	for i := 0; i < 10; i++ {
		if rl.Allow() {
			allowed++
		}
	}

	if allowed != 5 {
		t.Fatalf("expected 5 requests allowed in burst, got %d", allowed)
	}
}

func TestRateLimiter_ZeroRPS(t *testing.T) {
	// Edge case: 0 RPS should deny all requests
	rl := NewRateLimiter(0)
	
	if rl.Allow() {
		t.Fatal("zero RPS should deny all requests")
	}
}

func TestRateLimiter_HighRPS(t *testing.T) {
	// Test with high RPS
	rl := NewRateLimiter(100)
	
	// Should allow many requests initially
	for i := 0; i < 100; i++ {
		if !rl.Allow() {
			t.Fatalf("request %d should be allowed with high RPS", i+1)
		}
	}
	
	// 101st request should be denied
	if rl.Allow() {
		t.Fatal("request beyond capacity should be denied")
	}
}

func TestRateLimiter_ConcurrentAccess(t *testing.T) {
	rl := NewRateLimiter(10)
	
	// Test concurrent access doesn't cause race conditions
	done := make(chan bool, 20)
	
	for i := 0; i < 20; i++ {
		go func() {
			rl.Allow() // Just call it, don't care about result
			done <- true
		}()
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 20; i++ {
		<-done
	}
	
	// If we get here without race detector issues, test passes
}