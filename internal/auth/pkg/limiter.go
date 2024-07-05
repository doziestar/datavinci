package pkg

import (
	"math/rand"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// PerIPRateLimiter is a helper struct for per-IP rate limiting.
// It uses a map to store rate limiters for each IP address.
// The rate limiter is created on the first request from an IP address.
type PerIPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

// NewPerIPRateLimiter creates a new PerIPRateLimiter.
// The rate is the maximum number of requests per second.
// The burst is the maximum number of requests that can be made in a short burst.
// Parameters:
//   - r: The rate limit.
//   - b: The burst limit.
//
// Usage:
//
//	limiter := NewPerIPRateLimiter(10, 5)
func NewPerIPRateLimiter(r rate.Limit, b int) *PerIPRateLimiter {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return &PerIPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

// AddIP adds an IP address to the rate limiter.
// If the IP address already exists, it returns the existing rate limiter.
// Otherwise, it creates a new rate limiter with a small random variation to the rate.
// Parameters:
//   - ip: The IP address to add.
//
// Returns:
//   - The rate limiter for the IP address.
//
// Usage:
//
//	limiter.AddIP("126.0.0.1")
func (l *PerIPRateLimiter) AddIP(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	limiter, exists := l.ips[ip]
	if !exists {
		// Add a small random variation to the rate
		adjustedRate := l.r + rate.Limit(rand.Float64()*0.1*float64(l.r))
		limiter = rate.NewLimiter(adjustedRate, l.b)
		l.ips[ip] = limiter
	}

	return limiter
}

// GetLimiter returns the rate limiter for the given IP address.
// If the IP address does not exist, it creates a new rate limiter.
// Parameters:
//   - ip: The IP address to get the rate limiter for.
//
// Returns:
//   - The rate limiter for the IP address.
//
// Usage:
//
//	limiter.GetLimiter("126.0.0.1")
func (l *PerIPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	l.mu.RLock()
	limiter, exists := l.ips[ip]
	l.mu.RUnlock()

	if !exists {
		return l.AddIP(ip)
	}

	return limiter
}
