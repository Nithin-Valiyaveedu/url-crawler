package middleware

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"
	"time"
	"url-crawler/internal/config"

	"github.com/labstack/echo/v4"
)

// AuthConfig holds authentication configuration
type AuthConfig struct {
	APIKeys   map[string]string // key hash -> name mapping
	SkipPaths []string          // paths that don't require authentication
}

// NewAuthConfigFromConfig creates an auth configuration from the main config
func NewAuthConfig(cfg config.AuthConfig) *AuthConfig {
	authConfig := &AuthConfig{
		APIKeys: make(map[string]string),
		SkipPaths: []string{
			"/health",
			"/api/health",
			"/", // Allow root path for basic health check
		},
	}

	// Add all configured API keys
	for key, name := range cfg.APIKeys {
		authConfig.AddAPIKey(key, name)
	}

	return authConfig
}

// AddAPIKey adds an API key to the configuration
func (ac *AuthConfig) AddAPIKey(key, name string) {
	hash := sha256.Sum256([]byte(key))
	ac.APIKeys[fmt.Sprintf("%x", hash)] = name
}

// shouldSkipAuth checks if the path should skip authentication
func (ac *AuthConfig) shouldSkipAuth(path string) bool {
	for _, skipPath := range ac.SkipPaths {
		if path == skipPath || strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// AuthMiddleware creates an authentication middleware
func AuthMiddleware(config *AuthConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip authentication for certain paths
			if config.shouldSkipAuth(c.Request().URL.Path) {
				return next(c)
			}

			// Get authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Missing authorization header",
				})
			}

			// Extract API key from Bearer token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid authorization format. Use 'Bearer <api-key>'",
				})
			}

			apiKey := parts[1]
			if apiKey == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Missing API key",
				})
			}

			// Hash the provided API key
			hash := sha256.Sum256([]byte(apiKey))
			keyHash := fmt.Sprintf("%x", hash)

			// Check if the API key exists
			if name, exists := config.APIKeys[keyHash]; exists {
				// Set user context for logging/auditing
				c.Set("api_key_name", name)
				c.Set("api_key_hash", keyHash)
				return next(c)
			}

			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid API key",
			})
		}
	}
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int                         // Max requests per minute
	WindowSize        time.Duration               // Time window for rate limiting
	KeyGenerator      func(c echo.Context) string // Function to generate rate limit key
}

// DefaultRateLimitConfig creates a default rate limit configuration
func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		RequestsPerMinute: 60, // 60 requests per minute
		WindowSize:        time.Minute,
		KeyGenerator: func(c echo.Context) string {
			// Use API key if available, otherwise use IP
			if keyName := c.Get("api_key_name"); keyName != nil {
				return fmt.Sprintf("api:%s", keyName)
			}
			return fmt.Sprintf("ip:%s", c.RealIP())
		},
	}
}

// NewRateLimitConfigFromConfig creates a rate limit configuration from the main config
func NewRateLimitConfig(cfg config.AuthConfig) *RateLimitConfig {
	if !cfg.RateLimitEnabled {
		// Return a config that allows unlimited requests
		return &RateLimitConfig{
			RequestsPerMinute: 999999, // Effectively unlimited
			WindowSize:        time.Minute,
			KeyGenerator: func(c echo.Context) string {
				return "unlimited"
			},
		}
	}

	return &RateLimitConfig{
		RequestsPerMinute: cfg.RequestsPerMinute,
		WindowSize:        cfg.RateLimitWindow,
		KeyGenerator: func(c echo.Context) string {
			// Use API key if available, otherwise use IP
			if keyName := c.Get("api_key_name"); keyName != nil {
				return fmt.Sprintf("api:%s", keyName)
			}
			return fmt.Sprintf("ip:%s", c.RealIP())
		},
	}
}

// Simple in-memory rate limiter
type rateLimiter struct {
	requests map[string][]time.Time
}

var globalRateLimiter = &rateLimiter{
	requests: make(map[string][]time.Time),
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(config *RateLimitConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := config.KeyGenerator(c)
			now := time.Now()

			// Clean up old requests
			if requests, exists := globalRateLimiter.requests[key]; exists {
				var validRequests []time.Time
				cutoff := now.Add(-config.WindowSize)

				for _, reqTime := range requests {
					if reqTime.After(cutoff) {
						validRequests = append(validRequests, reqTime)
					}
				}

				globalRateLimiter.requests[key] = validRequests
			}

			// Check if limit exceeded
			if len(globalRateLimiter.requests[key]) >= config.RequestsPerMinute {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "Rate limit exceeded",
				})
			}

			// Add current request
			globalRateLimiter.requests[key] = append(globalRateLimiter.requests[key], now)

			return next(c)
		}
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			rid := req.Header.Get(echo.HeaderXRequestID)
			if rid == "" {
				rid = generateRequestID()
			}

			res.Header().Set(echo.HeaderXRequestID, rid)
			c.Set("request_id", rid)

			return next(c)
		}
	}
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// SecurityHeadersMiddleware adds security headers to responses
func SecurityHeadersMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			res := c.Response()

			// Security headers
			res.Header().Set("X-Content-Type-Options", "nosniff")
			res.Header().Set("X-Frame-Options", "DENY")
			res.Header().Set("X-XSS-Protection", "1; mode=block")
			res.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			res.Header().Set("Content-Security-Policy", "default-src 'self'")

			return next(c)
		}
	}
}
