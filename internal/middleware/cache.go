package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"api-gateway/config"
)

var (
	redisClient *redis.Client
	cacheCfg    *config.CacheConfig
)

// InitRedis initializes the Redis client with the given address and cache config
func InitRedis(addr string, cfg *config.CacheConfig) error {
	cacheCfg = cfg
	redisClient = redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// Test the connection
	_, err := redisClient.Ping(context.Background()).Result()
	return err
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Cache middleware caches GET requests using Redis
func Cache() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip caching for non-GET requests
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// Skip if Redis client is not initialized
		if redisClient == nil {
			c.Next()
			return
		}

		key := c.Request.URL.String()

		// Try to get from cache
		if val, err := redisClient.Get(context.Background(), key).Result(); err == nil {
			var cached struct {
				Status int         `json:"status"`
				Header http.Header `json:"header"`
				Data   string      `json:"data"`
			}

			if err := json.Unmarshal([]byte(val), &cached); err == nil {
				// Set headers from cache
				for k, v := range cached.Header {
					for _, vv := range v {
						c.Writer.Header().Add(k, vv)
					}
				}
				c.Writer.WriteHeader(cached.Status)
				c.Writer.Write([]byte(cached.Data))
				c.Abort()
				return
			}
		}

		// Create a custom ResponseWriter to capture the response
		w := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = w

		c.Next()

		// Only cache successful responses
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 && !c.IsAborted() {
			cached := struct {
				Status int         `json:"status"`
				Header http.Header `json:"header"`
				Data   string      `json:"data"`
			}{
				Status: c.Writer.Status(),
				Header: w.Header(),
				Data:   w.body.String(),
			}

			// Attempt to cache but don't block on errors
			if data, err := json.Marshal(cached); err == nil {
				ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
				defer cancel()

				redisClient.Set(
					ctx,
					key,
					data,
					time.Duration(cacheCfg.Duration)*time.Second,
				)
			}
		}
	}
}
