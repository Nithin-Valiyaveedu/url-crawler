package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// CrawlStatus represents the different states of a crawl operation
type CrawlStatus string

const (
	CrawlStatusQueued    CrawlStatus = "queued"
	CrawlStatusRunning   CrawlStatus = "running"
	CrawlStatusCompleted CrawlStatus = "completed"
	CrawlStatusError     CrawlStatus = "error"
)

// HeadingCounts represents the count of each heading level
type HeadingCounts struct {
	H1 int `json:"h1" db:"h1"`
	H2 int `json:"h2" db:"h2"`
	H3 int `json:"h3" db:"h3"`
	H4 int `json:"h4" db:"h4"`
	H5 int `json:"h5" db:"h5"`
	H6 int `json:"h6" db:"h6"`
}

// Value implements the driver.Valuer interface for database storage
func (h HeadingCounts) Value() (driver.Value, error) {
	return json.Marshal(h)
}

// Scan implements the sql.Scanner interface for database retrieval
func (h *HeadingCounts) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, h)
}

// BrokenLink represents a link that returned an error status
type BrokenLink struct {
	URL        string `json:"url" db:"url"`
	StatusCode int    `json:"statusCode" db:"status_code"`
	StatusText string `json:"statusText" db:"status_text"`
}

// BrokenLinks is a slice of BrokenLink that can be stored in database as JSON
type BrokenLinks []BrokenLink

// Value implements the driver.Valuer interface for database storage
func (bl BrokenLinks) Value() (driver.Value, error) {
	if bl == nil {
		return "[]", nil
	}
	return json.Marshal(bl)
}

// Scan implements the sql.Scanner interface for database retrieval
func (bl *BrokenLinks) Scan(value interface{}) error {
	if value == nil {
		*bl = BrokenLinks{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, bl)
}

// ExternalLinks is a slice of strings representing external link URLs
type ExternalLinks []string

// Value implements the driver.Valuer interface for database storage
func (el ExternalLinks) Value() (driver.Value, error) {
	if el == nil {
		return "[]", nil
	}
	return json.Marshal(el)
}

// Scan implements the sql.Scanner interface for database retrieval
func (el *ExternalLinks) Scan(value interface{}) error {
	if value == nil {
		*el = ExternalLinks{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, el)
}

// CrawlResult represents the complete analysis result of a crawled URL
type CrawlResult struct {
	ID                     string        `json:"id" db:"id"`
	URL                    string        `json:"url" db:"url"`
	Title                  string        `json:"title" db:"title"`
	HTMLVersion            string        `json:"htmlVersion" db:"html_version"`
	InternalLinksCount     int           `json:"internalLinksCount" db:"internal_links_count"`
	ExternalLinksCount     int           `json:"externalLinksCount" db:"external_links_count"`
	InaccessibleLinksCount int           `json:"inaccessibleLinksCount" db:"inaccessible_links_count"`
	HasLoginForm           bool          `json:"hasLoginForm" db:"has_login_form"`
	HeadingCounts          HeadingCounts `json:"headingCounts" db:"heading_counts"`
	BrokenLinks            BrokenLinks   `json:"brokenLinks" db:"broken_links"`
	ExternalLinks          ExternalLinks `json:"externalLinks" db:"external_links"`
	Status                 CrawlStatus   `json:"status" db:"status"`
	ErrorMessage           *string       `json:"errorMessage,omitempty" db:"error_message"`
	CreatedAt              time.Time     `json:"createdAt" db:"created_at"`
	UpdatedAt              time.Time     `json:"updatedAt" db:"updated_at"`
}

// CrawlRequest represents a request to crawl a URL
type CrawlRequest struct {
	URL string `json:"url" validate:"required,url"`
}

// CrawlRequestResponse represents the response when a crawl is requested
type CrawlRequestResponse struct {
	ID      string      `json:"id"`
	URL     string      `json:"url"`
	Status  CrawlStatus `json:"status"`
	Message string      `json:"message"`
}

// PaginatedCrawlResults represents paginated crawl results
type PaginatedCrawlResults struct {
	Results    []CrawlResult `json:"results"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"pageSize"`
	TotalPages int           `json:"totalPages"`
}

// CrawlFilters represents filters for querying crawl results
type CrawlFilters struct {
	Status   *CrawlStatus `json:"status,omitempty"`
	Search   string       `json:"search,omitempty"`
	Page     int          `json:"page"`
	PageSize int          `json:"pageSize"`
	SortBy   string       `json:"sortBy,omitempty"`
	SortDir  string       `json:"sortDir,omitempty"`
}

// CrawlStats represents statistics about crawl operations
type CrawlStats struct {
	Total     int `json:"total"`
	Queued    int `json:"queued"`
	Running   int `json:"running"`
	Completed int `json:"completed"`
	Error     int `json:"error"`
}

// ValidateStatus checks if the provided status is valid
func (status CrawlStatus) IsValid() bool {
	switch status {
	case CrawlStatusQueued, CrawlStatusRunning, CrawlStatusCompleted, CrawlStatusError:
		return true
	default:
		return false
	}
}

// String returns the string representation of CrawlStatus
func (status CrawlStatus) String() string {
	return string(status)
}

// DefaultFilters returns default filter values
func DefaultFilters() CrawlFilters {
	return CrawlFilters{
		Page:     1,
		PageSize: 10,
		SortBy:   "updated_at",
		SortDir:  "desc",
	}
}

// Validate validates the crawl filters
func (f *CrawlFilters) Validate() error {
	if f.Page < 1 {
		f.Page = 1
	}

	if f.PageSize < 1 || f.PageSize > 100 {
		f.PageSize = 10
	}

	if f.SortBy == "" {
		f.SortBy = "updated_at"
	}

	if f.SortDir != "asc" && f.SortDir != "desc" {
		f.SortDir = "desc"
	}

	if f.Status != nil && !f.Status.IsValid() {
		return errors.New("invalid status filter")
	}

	return nil
}
