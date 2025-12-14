package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"ip_country_project/internal/errors"
	"ip_country_project/internal/models"
)

// RateLimiter implements a token bucket algorithm for rate limiting
type RateLimiter struct {
	capacity   float64   // maximum tokens in bucket (equals RPS)
	tokens     float64   // current available tokens
	refillRate float64   // tokens added per second (equals RPS)
	lastRefill time.Time // last time tokens were refilled
	mutex      sync.Mutex // protects concurrent access to token state
}

// NewRateLimiter creates a new token bucket rate limiter with the specified requests per second
func NewRateLimiter(rps float64) *RateLimiter {
	return &RateLimiter{
		capacity:   rps,            // bucket capacity equals max RPS
		tokens:     rps,            // start with full bucket
		refillRate: rps,            // refill at RPS rate
		lastRefill: time.Now(),     // track when bucket was last refilled
	}
}

// Allow checks if a request should be allowed based on token bucket algorithm
func (l *RateLimiter) Allow() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(l.lastRefill).Seconds()

	// Refill tokens based on elapsed time
	if elapsed > 0 {
		refill := elapsed * l.refillRate
		l.tokens += refill
		if l.tokens > l.capacity {
			l.tokens = l.capacity
		}
		l.lastRefill = now
	}

	// Check if we have tokens available
	if l.tokens >= 1 {
		l.tokens -= 1
		return true
	}

	return false
}

// Middleware returns an HTTP middleware that enforces rate limiting
func (l *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !l.Allow() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			if err := json.NewEncoder(w).Encode(models.ErrorResponse{Error: errors.ErrRateLimited.Error()}); err != nil {
				w.Write([]byte(errors.ErrRateLimited.Error()))
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}
