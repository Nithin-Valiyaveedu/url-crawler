package server

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"

	"url-crawler/internal/config"
	"url-crawler/internal/database"
	"url-crawler/internal/services"
)

type Server struct {
	port           int
	db             database.Service
	crawlerService services.Crawler
	crawlStorage   *database.CrawlStorage
}

func NewServer() *http.Server {
	// Load configuration
	cfg := config.Load()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Log configuration (without sensitive data)
	cfg.LogConfig()

	// Initialize database service with configuration
	dbService := database.New(cfg.Database)

	// Get the underlying database connection using the new GetDB method
	db := dbService.GetDB()
	crawlStorage := database.NewCrawlStorage(db)

	// Initialize Firecrawl crawler service with configuration
	crawlerService := services.NewFirecrawlService(cfg.Crawler)
	if crawlerService == nil {
		log.Fatal("Failed to initialize Firecrawl service. Please ensure FIRECRAWL_API_KEY is set.")
	}
	log.Printf("Initialized Firecrawl crawler service")

	newServer := &Server{
		port:           cfg.Server.Port,
		db:             dbService,
		crawlerService: crawlerService,
		crawlStorage:   crawlStorage,
	}

	// Declare Server config with proper configuration values
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      newServer.RegisterRoutes(cfg.Auth),
		IdleTimeout:  cfg.Server.IdleTimeout,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return server
}
