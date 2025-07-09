package services

import "url-crawler/internal/models"

// Crawler interface
type Crawler interface {
	// AnalyzeURL performs comprehensive analysis of the given URL
	AnalyzeURL(targetURL string) (*models.CrawlResult, error)

	// ValidateURL validates URL format before crawling
	ValidateURL(targetURL string) error
}
