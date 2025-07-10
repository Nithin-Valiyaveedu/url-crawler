package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"url-crawler/internal/models"
)

// CrawlStorage implements the storage interface for crawl results
type CrawlStorage struct {
	db *sql.DB
}

// NewCrawlStorage creates a new crawl storage instance
func NewCrawlStorage(db *sql.DB) *CrawlStorage {
	return &CrawlStorage{db: db}
}

// SaveCrawlResult saves or updates a crawl result in the database
func (cs *CrawlStorage) SaveCrawlResult(result *models.CrawlResult) error {
	query := `
		INSERT INTO crawl_results (
			id, url, title, html_version, internal_links_count, external_links_count,
			inaccessible_links_count, has_login_form, heading_counts, broken_links,
			external_links, status, error_message, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			title = VALUES(title),
			html_version = VALUES(html_version),
			internal_links_count = VALUES(internal_links_count),
			external_links_count = VALUES(external_links_count),
			inaccessible_links_count = VALUES(inaccessible_links_count),
			has_login_form = VALUES(has_login_form),
			heading_counts = VALUES(heading_counts),
			broken_links = VALUES(broken_links),
			external_links = VALUES(external_links),
			status = VALUES(status),
			error_message = VALUES(error_message),
			updated_at = VALUES(updated_at)
	`

	_, err := cs.db.Exec(query,
		result.ID,
		result.URL,
		result.Title,
		result.HTMLVersion,
		result.InternalLinksCount,
		result.ExternalLinksCount,
		result.InaccessibleLinksCount,
		result.HasLoginForm,
		result.HeadingCounts,
		result.BrokenLinks,
		result.ExternalLinks,
		result.Status,
		result.ErrorMessage,
		result.CreatedAt,
		result.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save crawl result: %w", err)
	}

	return nil
}

// UpdateCrawlStatus updates only the status and error message of a crawl result
func (cs *CrawlStorage) UpdateCrawlStatus(id string, status models.CrawlStatus, errorMsg *string) error {
	query := `
		UPDATE crawl_results 
		SET status = ?, error_message = ?, updated_at = ? 
		WHERE id = ?
	`

	_, err := cs.db.Exec(query, status, errorMsg, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update crawl status: %w", err)
	}

	return nil
}

// GetCrawlResult retrieves a single crawl result by ID
func (cs *CrawlStorage) GetCrawlResult(id string) (*models.CrawlResult, error) {
	query := `
		SELECT id, url, title, html_version, internal_links_count, external_links_count,
			   inaccessible_links_count, has_login_form, heading_counts, broken_links,
			   external_links, status, error_message, created_at, updated_at
		FROM crawl_results 
		WHERE id = ?
	`

	row := cs.db.QueryRow(query, id)

	result := &models.CrawlResult{}

	err := row.Scan(
		&result.ID,
		&result.URL,
		&result.Title,
		&result.HTMLVersion,
		&result.InternalLinksCount,
		&result.ExternalLinksCount,
		&result.InaccessibleLinksCount,
		&result.HasLoginForm,
		&result.HeadingCounts,
		&result.BrokenLinks,
		&result.ExternalLinks,
		&result.Status,
		&result.ErrorMessage,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("crawl result not found")
		}
		return nil, fmt.Errorf("failed to get crawl result: %w", err)
	}

	return result, nil
}

// GetCrawlResults retrieves crawl results with filtering, sorting, and pagination
func (cs *CrawlStorage) GetCrawlResults(filters models.CrawlFilters) (*models.PaginatedCrawlResults, error) {
	// Validate filters
	if err := filters.Validate(); err != nil {
		return nil, fmt.Errorf("invalid filters: %w", err)
	}

	// Build WHERE clause
	var whereConditions []string
	var args []interface{}

	if filters.Status != nil {
		whereConditions = append(whereConditions, "status = ?")
		args = append(args, *filters.Status)
	}

	if filters.Search != "" {
		whereConditions = append(whereConditions, "(url LIKE ? OR title LIKE ?)")
		searchTerm := "%" + filters.Search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count total results
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM crawl_results 
		%s
	`, whereClause)

	var total int
	err := cs.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count results: %w", err)
	}

	// Calculate pagination
	offset := (filters.Page - 1) * filters.PageSize
	totalPages := (total + filters.PageSize - 1) / filters.PageSize

	// Build main query
	query := fmt.Sprintf(`
		SELECT id, url, title, html_version, internal_links_count, external_links_count,
			   inaccessible_links_count, has_login_form, heading_counts, broken_links,
			   external_links, status, error_message, created_at, updated_at
		FROM crawl_results 
		%s
		ORDER BY %s %s
		LIMIT ? OFFSET ?
	`, whereClause, filters.SortBy, filters.SortDir)

	// Add pagination parameters
	args = append(args, filters.PageSize, offset)

	rows, err := cs.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query crawl results: %w", err)
	}
	defer rows.Close()

	var results []models.CrawlResult
	for rows.Next() {
		result := models.CrawlResult{}

		err := rows.Scan(
			&result.ID,
			&result.URL,
			&result.Title,
			&result.HTMLVersion,
			&result.InternalLinksCount,
			&result.ExternalLinksCount,
			&result.InaccessibleLinksCount,
			&result.HasLoginForm,
			&result.HeadingCounts,
			&result.BrokenLinks,
			&result.ExternalLinks,
			&result.Status,
			&result.ErrorMessage,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan crawl result: %w", err)
		}

		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating results: %w", err)
	}

	return &models.PaginatedCrawlResults{
		Results:    results,
		Total:      total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}, nil
}

// DeleteCrawlResults deletes multiple crawl results by their IDs
func (cs *CrawlStorage) DeleteCrawlResults(ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	// Create placeholders for the IN clause
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
		DELETE FROM crawl_results 
		WHERE id IN (%s)
	`, strings.Join(placeholders, ","))

	result, err := cs.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete crawl results: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no crawl results were deleted")
	}

	return nil
}

// GetCrawlStats returns statistics about crawl results
func (cs *CrawlStorage) GetCrawlStats() (*models.CrawlStats, error) {
	query := `
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN status = 'queued' THEN 1 ELSE 0 END) as queued,
			SUM(CASE WHEN status = 'running' THEN 1 ELSE 0 END) as running,
			SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed,
			SUM(CASE WHEN status = 'error' THEN 1 ELSE 0 END) as error
		FROM crawl_results
	`

	stats := &models.CrawlStats{}
	err := cs.db.QueryRow(query).Scan(
		&stats.Total,
		&stats.Queued,
		&stats.Running,
		&stats.Completed,
		&stats.Error,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get crawl stats: %w", err)
	}

	return stats, nil
}

// UpdateCrawlResultsBulkStatus updates the status of multiple crawl results
func (cs *CrawlStorage) UpdateCrawlResultsBulkStatus(ids []string, status models.CrawlStatus) error {
	if len(ids) == 0 {
		return nil
	}

	// Create placeholders for the IN clause
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids)+2) // +2 for status and updated_at
	for i, id := range ids {
		placeholders[i] = "?"
		args[i+2] = id // Start from index 2
	}

	// Set status and updated_at as first arguments
	args[0] = status
	args[1] = time.Now()

	query := fmt.Sprintf(`
		UPDATE crawl_results 
		SET status = ?, updated_at = ?, error_message = NULL
		WHERE id IN (%s)
	`, strings.Join(placeholders, ","))

	result, err := cs.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update crawl results status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no crawl results were updated")
	}

	return nil
}

// CleanupOldCrawlResults removes crawl results older than the specified duration
func (cs *CrawlStorage) CleanupOldCrawlResults(olderThan time.Duration) (int64, error) {
	cutoffTime := time.Now().Add(-olderThan)

	query := `
		DELETE FROM crawl_results 
		WHERE created_at < ? AND status IN ('completed', 'error')
	`

	result, err := cs.db.Exec(query, cutoffTime)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup old crawl results: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get affected rows: %w", err)
	}

	return rowsAffected, nil
}
