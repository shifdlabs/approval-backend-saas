package middleware

import (
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
			for ip, b := range clients {
				if now.After(b.reset) {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		mu.Lock()
		b, ok := clients[ip]
		if !ok || now.After(b.reset) {
			b = &bucket{count: 0, reset: now.Add(window)}
			clients[ip] = b
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
