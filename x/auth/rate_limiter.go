package auth

import (
	"errors"
	"time"
)

// RateLimiter controls the rate of operations.
type RateLimiter struct {
	capacity     int
	tokens       int
	rate         time.Duration
	lastRefill   time.Time
}

// NewRateLimiter creates a new RateLimiter with the specified capacity and refill rate.
func NewRateLimiter(capacity int, rate time.Duration) *RateLimiter {
	return &RateLimiter{
		capacity:   capacity,
		tokens:     capacity,
		rate:       rate,
		lastRefill: time.Now(),
	}
}

// Allow checks if an operation is allowed under the current rate limit.
func (rl *RateLimiter) Allow() error {
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)

	// Refill tokens based on elapsed time
	numTokens := int(elapsed / rl.rate)
	if numTokens > 0 {
		rl.tokens = min(rl.capacity, rl.tokens+numTokens)
		rl.lastRefill = now
	}

	if rl.tokens > 0 {
		rl.tokens--
		return nil
	}

	return errors.New("rate limit exceeded")
}
