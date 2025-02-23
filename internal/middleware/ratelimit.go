package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"api-gateway/config"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Define a struct to hold client's rate limiter and last seen time
type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Create a map to store client rate limiters and configuration
var (
	clients = make(map[string]*client)
	mu      sync.Mutex
	cfg     *config.RateLimitConfig
)

// InitRateLimit initializes the rate limiter with configuration
func InitRateLimit(config *config.RateLimitConfig) {
	cfg = config
	// Start cleanup routine
	go cleanupRoutine()
}

// cleanupRoutine periodically cleans up old rate limiters
func cleanupRoutine() {
	for {
		fmt.Println("cleanupRoutine")
		time.Sleep(time.Duration(cfg.CleanupInterval) * time.Minute)
		mu.Lock()
		for ip, client := range clients {
			if time.Since(client.lastSeen) > time.Duration(cfg.CleanupInterval)*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

// getClientLimiter returns a rate limiter for the provided IP address
func getClientLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if c, exists := clients[ip]; exists {
		c.lastSeen = time.Now()
		return c.limiter
	}

	// Create a new rate limiter using configured values
	limiter := rate.NewLimiter(rate.Every(time.Minute/time.Duration(cfg.RequestsPerMinute)), cfg.BurstSize)
	clients[ip] = &client{
		limiter:  limiter,
		lastSeen: time.Now(),
	}

	return limiter
}

// RateLimit middleware for gin
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if rate limit is configured
		if cfg == nil {
			c.Next()
			return
		}

		ip := c.ClientIP()
		limiter := getClientLimiter(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
