package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"url-crawler/internal/config"
	customMiddleware "url-crawler/internal/middleware"
)

func (s *Server) RegisterRoutes(authCfg config.AuthConfig) http.Handler {
	e := echo.New()

	//Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(customMiddleware.RequestIDMiddleware())
	e.Use(customMiddleware.SecurityHeadersMiddleware())

	// CORS configuration
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Setup authentication using provided configuration
	authConfig := customMiddleware.NewAuthConfig(authCfg)

	// Rate limiting using configuration
	rateLimitConfig := customMiddleware.NewRateLimitConfig(authCfg)

	// Apply authentication and rate limiting middleware
	e.Use(customMiddleware.AuthMiddleware(authConfig))
	e.Use(customMiddleware.RateLimitMiddleware(rateLimitConfig))

	// Basic routes (no auth required)
	e.GET("/", s.APIInfoHandler)
	e.GET("/health", s.healthHandler)

	// API group
	api := e.Group("/api")

	// Health check endpoint
	api.GET("/health", s.crawlHandler.HealthCheck)

	// Crawl endpoints
	crawlGroup := api.Group("/crawl")
	{
		// Create new crawl request
		crawlGroup.POST("", s.crawlHandler.CreateCrawlRequest)

		// Get all crawl results (with pagination, filtering, sorting)
		crawlGroup.GET("", s.crawlHandler.GetCrawlResults)

		// Get crawl statistics
		crawlGroup.GET("/stats", s.crawlHandler.GetCrawlStats)

		// Bulk operations
		crawlGroup.POST("/rerun", s.crawlHandler.RerunCrawlResults)
		crawlGroup.DELETE("", s.crawlHandler.DeleteCrawlResults)

		// Individual crawl result operations
		crawlGroup.GET("/:id", s.crawlHandler.GetCrawlResult)
		crawlGroup.GET("/:id/status", s.crawlHandler.GetCrawlStatus)
	}

	return e
}

func (s *Server) APIInfoHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "URL Crawler API",
		"version": "1.0.0",
		"status":  "healthy",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c echo.Context) error {
	dbHealth := s.db.Health()

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": c.Get("request_id"),
		"database":  dbHealth,
		"services": map[string]string{
			"crawler": "healthy",
			"queue":   "healthy",
		},
	}

	return c.JSON(http.StatusOK, response)
}
