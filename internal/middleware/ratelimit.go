package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"ip_country_project/internal/models"
)

type RateLimiter struct {
	capacity   float64   // max tokens = RPS
	tokens     float64   // current tokens
	refillRate float64   // tokens per second
	lastRefill time.Time // last refill timestamp
	mutex      sync.Mutex
}

func NewRateLimiter(rps float64) *RateLimiter {
	return &RateLimiter{
		capacity:   rps,
		tokens:     rps, // start full
		refillRate: rps,
		lastRefill: time.Now(),
	}
}

func (l *RateLimiter) Allow() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(l.lastRefill).Seconds()

	if elapsed > 0 {
		refill := elapsed * l.refillRate
		l.tokens += refill
		if l.tokens > l.capacity {
			l.tokens = l.capacity
		}
		l.lastRefill = now
	}

	if l.tokens >= 1 {
		l.tokens -= 1
		return true
	}

	return false
}

func (l *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !l.Allow() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "rate limit exceeded"})
			return
		}
		next.ServeHTTP(w, r)
	})
}