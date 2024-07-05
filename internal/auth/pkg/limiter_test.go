package pkg_test

import (
	"auth/pkg"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestPerIPRateLimiter(t *testing.T) {
	t.Run("NewPerIPRateLimiter", func(t *testing.T) {
		limiter := pkg.NewPerIPRateLimiter(1, 5)
		assert.NotNil(t, limiter, "NewPerIPRateLimiter should return a non-nil limiter")
	})

	t.Run("AddIP", func(t *testing.T) {
		limiter := pkg.NewPerIPRateLimiter(1, 5)
		ip := "192.168.1.1"
		rateLimiter := limiter.AddIP(ip)
		assert.NotNil(t, rateLimiter, "AddIP should return a non-nil rate.Limiter")
	})

	t.Run("GetLimiter", func(t *testing.T) {
		limiter := pkg.NewPerIPRateLimiter(1, 5)
		ip := "192.168.1.1"

		// First call should add the IP
		rateLimiter1 := limiter.GetLimiter(ip)
		assert.NotNil(t, rateLimiter1, "GetLimiter should return a non-nil rate.Limiter")

		// Second call should return the same limiter
		rateLimiter2 := limiter.GetLimiter(ip)
		assert.Equal(t, rateLimiter1, rateLimiter2, "GetLimiter should return the same rate.Limiter for the same IP")
	})

	t.Run("RateLimitingBehavior", func(t *testing.T) {
		limiter := pkg.NewPerIPRateLimiter(10, 5) // 10 requests per second, burst of 5
		ip := "192.168.1.1"

		rateLimiter := limiter.GetLimiter(ip)

		// Should allow burst
		for i := 0; i < 5; i++ {
			assert.True(t, rateLimiter.Allow(), "Should allow burst of 5 requests")
		}

		// Next request should be rate limited
		assert.False(t, rateLimiter.Allow(), "Should not allow 6th request immediately")

		// Wait for 100ms, should allow one more request
		time.Sleep(100 * time.Millisecond)
		assert.True(t, rateLimiter.Allow(), "Should allow request after waiting")
	})

	t.Run("ConcurrentAccess", func(t *testing.T) {
		limiter := pkg.NewPerIPRateLimiter(100, 10) // High limit to avoid rate limiting in this test
		ips := []string{"192.168.1.1", "192.168.1.2", "192.168.1.3", "192.168.1.4", "192.168.1.5"}

		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for _, ip := range ips {
					rateLimiter := limiter.GetLimiter(ip)
					assert.NotNil(t, rateLimiter, "GetLimiter should return a non-nil rate.Limiter")
				}
			}()
		}

		wg.Wait()
	})

	t.Run("DifferentIPsDifferentLimiters", func(t *testing.T) {
		limiter := pkg.NewPerIPRateLimiter(1, 5)
		ip1 := "192.168.1.1"
		ip2 := "192.168.1.2"

		rateLimiter1 := limiter.GetLimiter(ip1)
		rateLimiter2 := limiter.GetLimiter(ip2)

		assert.NotEqual(t, rateLimiter1, rateLimiter2, "Different IPs should have different rate limiters")

		// Additional check to ensure the limiters are truly different
		rateLimiter1.Allow() // Use up one token from the first limiter
		assert.True(t, rateLimiter2.Allow(), "Second limiter should still allow a request")
	})

	t.Run("RespectsBurstAndRate", func(t *testing.T) {
		r := rate.Limit(2) // 2 requests per second
		b := 3             // burst of 3
		limiter := pkg.NewPerIPRateLimiter(r, b)
		ip := "192.168.1.1"
		rateLimiter := limiter.GetLimiter(ip)

		// Should allow burst
		for i := 0; i < b; i++ {
			assert.True(t, rateLimiter.Allow(), "Should allow burst of 3 requests")
		}

		// Next request should be rate limited
		assert.False(t, rateLimiter.Allow(), "Should not allow 4th request immediately")

		// Wait for 500ms, should allow one more request
		time.Sleep(500 * time.Millisecond)
		assert.True(t, rateLimiter.Allow(), "Should allow request after waiting")

		// Next request should be rate limited again
		assert.False(t, rateLimiter.Allow(), "Should not allow another request immediately")
	})
}
