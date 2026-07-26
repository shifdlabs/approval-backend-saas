package middleware

import (
	"Microservice/helper"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimit is a lightweight in-memory, per-client-IP fixed-window rate limiter.
// It exists to blunt scraping/abuse of the unauthenticated public verification
// endpoint (AUDIT SEC-10) without pulling in an external dependency or Redis.
//
// Note: state is per-process, so behind multiple replicas the effective limit is
// per instance. For a single-instance deployment this is sufficient; move to a
// shared store (Redis) if the service is horizontally scaled.
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	return RateLimitBy(func(c *gin.Context) string { return c.ClientIP() }, limit, window)
}

// RateLimitByUser buckets by authenticated user id instead of client IP. Use it
// for authenticated, expensive endpoints (e.g. AI draft generation): a whole
// office shares one public IP, so an IP bucket would let one user exhaust
// everyone else's quota. Falls back to the IP when the request is unauthenticated.
func RateLimitByUser(limit int, window time.Duration) gin.HandlerFunc {
	return RateLimitBy(func(c *gin.Context) string {
		if userID := helper.GetUserID(c); userID != "" {
			return userID
		}
		return c.ClientIP()
	}, limit, window)
}

// RateLimitBy is the shared implementation; keyFn decides what counts as one
// client. Each call gets its own independent bucket store.
func RateLimitBy(keyFn func(*gin.Context) string, limit int, window time.Duration) gin.HandlerFunc {
	type bucket struct {
		count int
		reset time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*bucket)
	)

	// Periodically evict stale buckets so the map does not grow unbounded.
	go func() {
		ticker := time.NewTicker(window)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			mu.Lock()
			for key, b := range clients {
				if now.After(b.reset) {
					delete(clients, key)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		key := keyFn(c)
		now := time.Now()

		mu.Lock()
		b, ok := clients[key]
		if !ok || now.After(b.reset) {
			b = &bucket{count: 0, reset: now.Add(window)}
			clients[key] = b
		}
		b.count++
		exceeded := b.count > limit
		retryAfter := int(time.Until(b.reset).Seconds())
		mu.Unlock()

		if exceeded {
			if retryAfter < 1 {
				retryAfter = 1
			}
			c.Header("Retry-After", itoa(retryAfter))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"code":    http.StatusTooManyRequests,
				"message": "Too many requests. Please try again later.",
			})
			return
		}

		c.Next()
	}
}

// itoa avoids importing strconv for a single small conversion.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
