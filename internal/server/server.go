package server

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"

	"url-crawler/internal/config"
	"url-crawler/internal/database"
	"url-crawler/internal/handlers"
	"url-crawler/internal/services"
)

type Server struct {
	port int

	//database
	db database.Service

	// Services
	crawlerService services.Crawler
	queueService   *services.QueueService
	crawlStorage   *database.CrawlStorage

	// Handlers
	crawlHandler *handlers.CrawlHandler
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

	// Initialize queue service with configuration
	queueService := services.NewQueueService(cfg.Queue.Workers, crawlerService, crawlStorage)

	// Initialize handlers
	crawlHandler := handlers.NewCrawlHandler(queueService, crawlStorage)

	newServer := &Server{
		port:           cfg.Server.Port,
		db:             dbService,
		crawlerService: crawlerService,
		queueService:   queueService,
		crawlStorage:   crawlStorage,
		crawlHandler:   crawlHandler,
	}

	// Start the queue service
	queueService.Start()

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
