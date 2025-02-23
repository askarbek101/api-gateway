package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server           ServerConfig           `mapstructure:"server"`
	Redis            RedisConfig            `mapstructure:"redis"`
	JWT              JWTConfig              `mapstructure:"jwt"`
	Cache            CacheConfig            `mapstructure:"cache"`
	RateLimit        RateLimitConfig        `mapstructure:"rate_limit"`
	ExternalServices ExternalServicesConfig `mapstructure:"external_services"`
}

type ServerConfig struct {
	Port         string   `mapstructure:"port"`
	Mode         string   `mapstructure:"mode"`
	TrustedProxy string   `mapstructure:"trusted_proxy"` // CIDR format for trusted proxies
	AllowOrigins []string `mapstructure:"allow_origins"` // CORS allowed origins
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	DB   int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

type CacheConfig struct {
	Duration int `mapstructure:"duration"` // Default cache duration in seconds
}

type RateLimitConfig struct {
	RequestsPerMinute int `mapstructure:"requests_per_minute"` // Number of requests allowed per minute
	BurstSize         int `mapstructure:"burst_size"`          // Maximum burst size
	CleanupInterval   int `mapstructure:"cleanup_interval"`    // Cleanup interval in minutes
}

type ExternalServicesConfig struct {
	UserService map[string]string `mapstructure:"user_service"`
}

// LoadConfig reads configuration from environment variables or config file
func LoadConfig() (*Config, error) {
	// Server defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.trusted_proxy", "127.0.0.1/32")

	// Default CORS origins - ensure at least one origin is allowed
	viper.SetDefault("server.allow_origins", []string{
		"http://localhost:8080",
		"http://127.0.0.1:8080",
	})

	// Redis defaults
	viper.SetDefault("redis.host", "gqetIorpvdYUyjDtOYwIzDhaiiQUoEOx@tramway.proxy.rlwy.net")
	viper.SetDefault("redis.port", "11897")
	viper.SetDefault("redis.db", 0)

	// JWT defaults
	viper.SetDefault("jwt.secret", "your-secret-key")

	// Cache defaults
	viper.SetDefault("cache.duration", 60) // 1 minute

	// Rate limit defaults
	viper.SetDefault("rate_limit.requests_per_minute", 100)
	viper.SetDefault("rate_limit.burst_size", 100)
	viper.SetDefault("rate_limit.cleanup_interval", 5)

	// External Services defaults
	viper.SetDefault("external_services.user_service", map[string]string{
		"base_url": "http://localhost:8081",
		"timeout":  "30s",
	})

	// Enable environment variable binding
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	// Replace dots with underscores in env variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Ensure we have at least one allowed origin
	if len(config.Server.AllowOrigins) == 0 {
		config.Server.AllowOrigins = []string{"http://localhost:8080"}
	}

	return &config, nil
}
