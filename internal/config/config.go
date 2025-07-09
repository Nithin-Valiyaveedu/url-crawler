package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Crawler  CrawlerConfig
	Queue    QueueConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Port         int
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	MaxOpen  int
	MaxIdle  int
	MaxLife  time.Duration
}

type CrawlerConfig struct {
	Timeout          time.Duration
	UserAgent        string
	MaxRedirects     int
	MaxLinksToCheck  int
	RequestDelay     time.Duration
	MaxContentSize   int64
	AllowedDomains   []string
	BlockedDomains   []string
	RespectRobotsTxt bool

	// Firecrawl configuration
	FirecrawlAPIKey string
	FirecrawlAPIURL string
}

type QueueConfig struct {
	Workers    int
	BufferSize int
	MaxRetries int
	RetryDelay time.Duration
}

type AuthConfig struct {
	APIKeys           map[string]string
	RequireAuth       bool
	RateLimitEnabled  bool
	RequestsPerMinute int
	RateLimitWindow   time.Duration
}

func Load() *Config {
	return &Config{
		Server:   loadServerConfig(),
		Database: loadDatabaseConfig(),
		Crawler:  loadCrawlerConfig(),
		Queue:    loadQueueConfig(),
		Auth:     loadAuthConfig(),
	}
}

// loadServerConfig loads server configuration from environment
func loadServerConfig() ServerConfig {
	port, _ := strconv.Atoi(getEnv("PORT", "8080"))
	readTimeout, _ := time.ParseDuration(getEnv("SERVER_READ_TIMEOUT", "10s"))
	writeTimeout, _ := time.ParseDuration(getEnv("SERVER_WRITE_TIMEOUT", "30s"))
	idleTimeout, _ := time.ParseDuration(getEnv("SERVER_IDLE_TIMEOUT", "60s"))

	return ServerConfig{
		Port:         port,
		Host:         getEnv("HOST", ""),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}
}

func loadDatabaseConfig() DatabaseConfig {
	maxOpen, _ := strconv.Atoi(getEnv("DB_MAX_OPEN", "50"))
	maxIdle, _ := strconv.Atoi(getEnv("DB_MAX_IDLE", "50"))
	maxLife, _ := time.ParseDuration(getEnv("DB_MAX_LIFE", "0"))

	return DatabaseConfig{
		Host:     getEnv("URL_CRAWLER_DB_HOST", "localhost"),
		Port:     getEnv("URL_CRAWLER_DB_PORT", "3306"),
		Username: getEnv("URL_CRAWLER_DB_USERNAME", "root"),
		Password: getEnv("URL_CRAWLER_DB_PASSWORD", ""),
		Database: getEnv("URL_CRAWLER_DB_DATABASE", "url_crawler"),
		MaxOpen:  maxOpen,
		MaxIdle:  maxIdle,
		MaxLife:  maxLife,
	}
}

func loadCrawlerConfig() CrawlerConfig {
	timeout, _ := time.ParseDuration(getEnv("CRAWLER_TIMEOUT", "30s"))
	maxRedirects, _ := strconv.Atoi(getEnv("CRAWLER_MAX_REDIRECTS", "10"))
	maxLinksToCheck, _ := strconv.Atoi(getEnv("CRAWLER_MAX_LINKS_CHECK", "10"))
	requestDelay, _ := time.ParseDuration(getEnv("CRAWLER_REQUEST_DELAY", "100ms"))
	maxContentSize, _ := strconv.ParseInt(getEnv("CRAWLER_MAX_CONTENT_SIZE", "10485760"), 10, 64) // 10MB
	respectRobots, _ := strconv.ParseBool(getEnv("CRAWLER_RESPECT_ROBOTS", "true"))

	allowedDomains := strings.Split(getEnv("CRAWLER_ALLOWED_DOMAINS", ""), ",")
	blockedDomains := strings.Split(getEnv("CRAWLER_BLOCKED_DOMAINS", ""), ",")

	// Clean up empty strings from domain lists
	allowedDomains = filterEmptyStrings(allowedDomains)
	blockedDomains = filterEmptyStrings(blockedDomains)

	return CrawlerConfig{
		Timeout:          timeout,
		UserAgent:        getEnv("CRAWLER_USER_AGENT", "URL-Crawler-Bot/1.0"),
		MaxRedirects:     maxRedirects,
		MaxLinksToCheck:  maxLinksToCheck,
		RequestDelay:     requestDelay,
		MaxContentSize:   maxContentSize,
		AllowedDomains:   allowedDomains,
		BlockedDomains:   blockedDomains,
		RespectRobotsTxt: respectRobots,

		// Firecrawl configuration
		FirecrawlAPIKey: getEnv("FIRECRAWL_API_KEY", ""),
		FirecrawlAPIURL: getEnv("FIRECRAWL_API_URL", ""),
	}
}

func loadQueueConfig() QueueConfig {
	workers, _ := strconv.Atoi(getEnv("QUEUE_WORKERS", "3"))
	bufferSize, _ := strconv.Atoi(getEnv("QUEUE_BUFFER_SIZE", "100"))
	maxRetries, _ := strconv.Atoi(getEnv("QUEUE_MAX_RETRIES", "3"))
	retryDelay, _ := time.ParseDuration(getEnv("QUEUE_RETRY_DELAY", "5s"))

	return QueueConfig{
		Workers:    workers,
		BufferSize: bufferSize,
		MaxRetries: maxRetries,
		RetryDelay: retryDelay,
	}
}

func loadAuthConfig() AuthConfig {
	requireAuth, _ := strconv.ParseBool(getEnv("AUTH_REQUIRED", "true"))
	rateLimitEnabled, _ := strconv.ParseBool(getEnv("RATE_LIMIT_ENABLED", "true"))
	requestsPerMinute, _ := strconv.Atoi(getEnv("RATE_LIMIT_REQUESTS_PER_MINUTE", "60"))
	rateLimitWindow, _ := time.ParseDuration(getEnv("RATE_LIMIT_WINDOW", "1m"))

	// Load API keys from environment
	apiKeys := make(map[string]string)

	// Support multiple API keys via environment variables
	// Format: API_KEY_<NAME>=<key>
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 && strings.HasPrefix(parts[0], "API_KEY_") {
			name := strings.TrimPrefix(parts[0], "API_KEY_")
			name = strings.ToLower(name)
			apiKeys[parts[1]] = name
		}
	}

	// // Add default development key if no keys are configured
	// if len(apiKeys) == 0 && !requireAuth {
	// 	apiKeys["dev-api-key-12345"] = "development"
	// }

	return AuthConfig{
		APIKeys:           apiKeys,
		RequireAuth:       requireAuth,
		RateLimitEnabled:  rateLimitEnabled,
		RequestsPerMinute: requestsPerMinute,
		RateLimitWindow:   rateLimitWindow,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// filterEmptyStrings removes empty strings from a slice
func filterEmptyStrings(slice []string) []string {
	var filtered []string
	for _, str := range slice {
		if trimmed := strings.TrimSpace(str); trimmed != "" {
			filtered = append(filtered, trimmed)
		}
	}
	return filtered
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return ErrInvalidPort
	}

	if c.Database.Host == "" {
		return ErrMissingDBHost
	}

	if c.Database.Username == "" {
		return ErrMissingDBUsername
	}

	if c.Queue.Workers <= 0 {
		return ErrInvalidWorkerCount
	}

	if c.Queue.BufferSize <= 0 {
		return ErrInvalidBufferSize
	}

	return nil
}

// Configuration errors
var (
	ErrInvalidPort        = fmt.Errorf("invalid port number")
	ErrMissingDBHost      = fmt.Errorf("database host is required")
	ErrMissingDBUsername  = fmt.Errorf("database username is required")
	ErrInvalidWorkerCount = fmt.Errorf("worker count must be greater than 0")
	ErrInvalidBufferSize  = fmt.Errorf("buffer size must be greater than 0")
)

// LogConfig logs the current configuration (without sensitive data)
func (c *Config) LogConfig() {
	log.Println("=== URL Crawler Configuration ===")
	log.Printf("Server: %s:%d", c.Server.Host, c.Server.Port)
	log.Printf("Database: %s:%s@%s:%s/%s", c.Database.Username, "***", c.Database.Host, c.Database.Port, c.Database.Database)
	log.Printf("Queue Workers: %d", c.Queue.Workers)
	log.Printf("Auth Required: %t", c.Auth.RequireAuth)
	log.Printf("Rate Limiting: %t", c.Auth.RateLimitEnabled)
	log.Printf("Crawler Timeout: %s", c.Crawler.Timeout)
	log.Printf("Crawler User Agent: %s", c.Crawler.UserAgent)
	log.Println("=================================")
}
