package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"url-crawler/internal/database"
	"url-crawler/internal/models"
	"url-crawler/internal/services"
)

// CrawlHandler handles all crawl-related HTTP requests
type CrawlHandler struct {
	queue     *services.QueueService
	storage   *database.CrawlStorage
	validator *validator.Validate
}

// NewCrawlHandler creates a new crawl handler
func NewCrawlHandler(queue *services.QueueService, storage *database.CrawlStorage) *CrawlHandler {
	return &CrawlHandler{
		queue:     queue,
		storage:   storage,
		validator: validator.New(),
	}
}

// CreateCrawlRequest handles POST /api/crawl requests
func (h *CrawlHandler) CreateCrawlRequest(c echo.Context) error {
	var req models.CrawlRequest

	// Bind request body
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request data: " + err.Error(),
		})
	}

	// Enqueue the URL for crawling
	result, err := h.queue.EnqueueURL(req.URL)
	if err != nil {
		if strings.Contains(err.Error(), "queue is full") {
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Return response
	response := models.CrawlRequestResponse{
		ID:      result.ID,
		URL:     result.URL,
		Status:  result.Status,
		Message: "URL added to crawl queue successfully",
	}

	return c.JSON(http.StatusCreated, response)
}

// GetCrawlResults handles GET /api/crawl requests
func (h *CrawlHandler) GetCrawlResults(c echo.Context) error {
	// Parse query parameters
	filters := models.DefaultFilters()

	// Parse page
	if pageStr := c.QueryParam("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filters.Page = page
		}
	}

	// Parse page size
	if pageSizeStr := c.QueryParam("pageSize"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 && pageSize <= 100 {
			filters.PageSize = pageSize
		}
	}

	// Parse search
	if search := c.QueryParam("search"); search != "" {
		filters.Search = search
	}

	// Parse status filter
	if statusStr := c.QueryParam("status"); statusStr != "" {
		status := models.CrawlStatus(statusStr)
		if status.IsValid() {
			filters.Status = &status
		}
	}

	// Parse sort parameters
	if sortBy := c.QueryParam("sortBy"); sortBy != "" {
		// Validate allowed sort fields
		allowedFields := map[string]bool{
			"url": true, "title": true, "status": true,
			"created_at": true, "updated_at": true,
			"internal_links_count": true, "external_links_count": true,
		}
		if allowedFields[sortBy] {
			filters.SortBy = sortBy
		}
	}

	if sortDir := c.QueryParam("sortDir"); sortDir == "asc" || sortDir == "desc" {
		filters.SortDir = sortDir
	}

	// Get results from storage
	results, err := h.storage.GetCrawlResults(filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve crawl results",
		})
	}

	return c.JSON(http.StatusOK, results)
}

// GetCrawlResult handles GET /api/crawl/:id requests
func (h *CrawlHandler) GetCrawlResult(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing crawl result ID",
		})
	}

	result, err := h.storage.GetCrawlResult(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Crawl result not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve crawl result",
		})
	}

	return c.JSON(http.StatusOK, result)
}

// DeleteCrawlResults handles DELETE /api/crawl requests
func (h *CrawlHandler) DeleteCrawlResults(c echo.Context) error {
	var req struct {
		IDs []string `json:"ids" validate:"required,min=1"`
	}

	// Bind request body
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request data: " + err.Error(),
		})
	}

	// Delete from storage
	err := h.storage.DeleteCrawlResults(req.IDs)
	if err != nil {
		if strings.Contains(err.Error(), "no crawl results were deleted") {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "No crawl results found for the provided IDs",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete crawl results",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":       "Crawl results deleted successfully",
		"deleted_count": len(req.IDs),
	})
}

// RerunCrawlResults handles POST /api/crawl/rerun requests
func (h *CrawlHandler) RerunCrawlResults(c echo.Context) error {
	var req struct {
		IDs []string `json:"ids" validate:"required,min=1"`
	}

	// Bind request body
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request data: " + err.Error(),
		})
	}

	// Update status to queued first
	err := h.storage.UpdateCrawlResultsBulkStatus(req.IDs, models.CrawlStatusQueued)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update crawl results status",
		})
	}

	// Requeue each task
	var successCount int
	var errors []string

	for _, id := range req.IDs {
		if err := h.queue.RequeueTask(id); err != nil {
			errors = append(errors, "Failed to requeue "+id+": "+err.Error())
		} else {
			successCount++
		}
	}

	response := map[string]interface{}{
		"message":         "Rerun operation completed",
		"success_count":   successCount,
		"total_requested": len(req.IDs),
	}

	if len(errors) > 0 {
		response["errors"] = errors
		return c.JSON(http.StatusPartialContent, response)
	}

	return c.JSON(http.StatusOK, response)
}

// GetCrawlStats handles GET /api/crawl/stats requests
func (h *CrawlHandler) GetCrawlStats(c echo.Context) error {
	// Get database stats
	dbStats, err := h.storage.GetCrawlStats()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve crawl statistics",
		})
	}

	// Get queue stats
	queueStats := h.queue.GetQueueStats()

	// Combine stats
	response := map[string]interface{}{
		"database":  dbStats,
		"queue":     queueStats,
		"timestamp": h.getCurrentTimestamp(),
	}

	return c.JSON(http.StatusOK, response)
}

// GetCrawlStatus handles GET /api/crawl/:id/status requests
func (h *CrawlHandler) GetCrawlStatus(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing crawl result ID",
		})
	}

	// First check if it's in the active queue
	if task, exists := h.queue.GetActiveTask(id); exists {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"id":        task.ID,
			"status":    task.Status,
			"url":       task.URL,
			"queued_at": task.CreatedAt,
		})
	}

	// If not in queue, get from database
	result, err := h.storage.GetCrawlResult(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Crawl result not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve crawl status",
		})
	}

	response := map[string]interface{}{
		"id":         result.ID,
		"status":     result.Status,
		"url":        result.URL,
		"created_at": result.CreatedAt,
		"updated_at": result.UpdatedAt,
	}

	if result.ErrorMessage != nil {
		response["error_message"] = *result.ErrorMessage
	}

	return c.JSON(http.StatusOK, response)
}

// HealthCheck handles GET /api/health requests
func (h *CrawlHandler) HealthCheck(c echo.Context) error {
	queueStats := h.queue.GetQueueStats()

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": h.getCurrentTimestamp(),
		"queue":     queueStats,
	}

	return c.JSON(http.StatusOK, response)
}

// Helper method to get current timestamp
func (h *CrawlHandler) getCurrentTimestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}

// validateCrawlIDs validates that all provided IDs exist in the database
func (h *CrawlHandler) validateCrawlIDs(ids []string) error {
	for _, id := range ids {
		_, err := h.storage.GetCrawlResult(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return fmt.Errorf("crawl result with ID %s not found", id)
			}
			return err
		}
	}
	return nil
}
