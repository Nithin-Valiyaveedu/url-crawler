package services

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"url-crawler/internal/config"
	"url-crawler/internal/models"

	"github.com/google/uuid"
	"github.com/mendableai/firecrawl-go"
)

// FirecrawlService implements the crawler interface using Firecrawl SDK
type FirecrawlService struct {
	app *firecrawl.FirecrawlApp
}

// NewFirecrawlServiceWithConfig creates a new Firecrawl-based crawler service using configuration
func NewFirecrawlService(cfg config.CrawlerConfig) *FirecrawlService {
	// Use configuration values
	apiKey := cfg.FirecrawlAPIKey
	apiUrl := cfg.FirecrawlAPIURL

	if apiKey == "" {
		log.Printf("Warning: FIRECRAWL_API_KEY not configured")
		return nil
	}

	if apiUrl == "" {
		apiUrl = "https://api.firecrawl.dev"
	}

	// Initialize the FirecrawlApp
	app, err := firecrawl.NewFirecrawlApp(apiKey, apiUrl)
	if err != nil {
		log.Printf("Warning: Failed to initialize FirecrawlApp: %v", err)
		return nil
	}

	log.Printf("Firecrawl service initialized with API URL: %s (using config)", apiUrl)
	return &FirecrawlService{
		app: app,
	}
}

// AnalyzeURL performs comprehensive analysis using Firecrawl
func (fs *FirecrawlService) AnalyzeURL(targetURL string) (*models.CrawlResult, error) {
	if fs.app == nil {
		return nil, fmt.Errorf("firecrawl service not properly initialized")
	}

	// Initialize result
	result := &models.CrawlResult{
		ID:        uuid.New().String(),
		URL:       targetURL,
		Status:    models.CrawlStatusRunning,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		HeadingCounts: models.HeadingCounts{
			H1: 0, H2: 0, H3: 0, H4: 0, H5: 0, H6: 0,
		},
		BrokenLinks: models.BrokenLinks{},
	}

	log.Printf("Starting Firecrawl analysis for URL: %s", targetURL)

	// Use ScrapeURL for single page analysis
	waitFor := 3000
	scrapeParams := &firecrawl.ScrapeParams{
		Formats:     []string{"markdown", "html"},
		IncludeTags: []string{"title", "h1", "h2", "h3", "h4", "h5", "h6", "form", "input", "a", "link"},
		WaitFor:     &waitFor,
	}
	scrapeResponse, err := fs.app.ScrapeURL(targetURL, scrapeParams)
	if err != nil {
		result.Status = models.CrawlStatusError
		errorMsg := fmt.Sprintf("Firecrawl scrape failed: %v", err)
		result.ErrorMessage = &errorMsg
		return result, fmt.Errorf("failed to scrape URL with Firecrawl: %w", err)
	}

	log.Printf("Firecrawl successfully scraped URL: %s", targetURL)

	// Extract data from Firecrawl response
	if err := fs.extractDataFromFirecrawlDocument(scrapeResponse, result); err != nil {
		log.Printf("Warning: Failed to extract some data from response: %v", err)
		// Don't fail the entire operation, just log the warning
	}

	// Set completion status
	result.Status = models.CrawlStatusCompleted
	result.UpdatedAt = time.Now()

	log.Printf("Firecrawl analysis completed for URL: %s", targetURL)
	return result, nil
}

// extractDataFromFirecrawlDocument extracts relevant data from Firecrawl document
func (fs *FirecrawlService) extractDataFromFirecrawlDocument(doc *firecrawl.FirecrawlDocument, result *models.CrawlResult) error {
	// Extract title from metadata
	if doc.Metadata != nil && doc.Metadata.Title != nil {
		result.Title = strings.TrimSpace(*doc.Metadata.Title)
	}

	// Extract HTML content
	if doc.HTML != "" {
		fs.analyzeHTMLContent(doc.HTML, result)
	}

	// Extract markdown content
	if doc.Markdown != "" {
		fs.analyzeMarkdownContent(doc.Markdown, result)
	}

	// Extract metadata if available
	if doc.Metadata != nil {
		fs.extractFirecrawlMetadata(doc.Metadata, result)
	}

	return nil
}

// analyzeHTMLContent analyzes HTML content for various elements
func (fs *FirecrawlService) analyzeHTMLContent(html string, result *models.CrawlResult) {
	// Detect login forms
	result.HasLoginForm = fs.detectLoginForm(html)

	// Count headings
	fs.countHeadings(html, result)

	// Analyze links
	fs.analyzeLinks(html, result)

	// Detect HTML version
	result.HTMLVersion = fs.detectHTMLVersion(html)
}

// analyzeMarkdownContent analyzes markdown content for additional insights
func (fs *FirecrawlService) analyzeMarkdownContent(markdown string, result *models.CrawlResult) {
	// Count headings in markdown (as backup/validation)
	headingRegex := regexp.MustCompile(`(?m)^#+\s+`)
	headingMatches := headingRegex.FindAllString(markdown, -1)

	// Validate heading counts against markdown
	markdownHeadingCount := len(headingMatches)
	if markdownHeadingCount > 0 {
		log.Printf("Markdown validation: Found %d headings in markdown content", markdownHeadingCount)
	}
}

// detectLoginForm analyzes HTML for login form patterns
func (fs *FirecrawlService) detectLoginForm(html string) bool {
	htmlLower := strings.ToLower(html)

	// Look for password fields (most reliable indicator)
	passwordPatterns := []string{
		`type="password"`,
		`type='password'`,
		`input[type="password"]`,
		`input[type='password']`,
	}

	hasPasswordField := false
	for _, pattern := range passwordPatterns {
		if strings.Contains(htmlLower, pattern) {
			hasPasswordField = true
			break
		}
	}

	if !hasPasswordField {
		return false
	}

	// Look for additional login indicators
	loginIndicators := []string{
		"login", "signin", "sign-in", "log-in", "auth", "authentication",
		"username", "email", "user", "account",
		"password", "pwd", "pass",
		"submit", "button",
		"loginform", "authform", "signupform",
	}

	indicatorCount := 0
	for _, indicator := range loginIndicators {
		if strings.Contains(htmlLower, indicator) {
			indicatorCount++
		}
	}

	// If we have a password field and multiple login indicators, it's likely a login form
	return indicatorCount >= 2
}

// countHeadings counts H1-H6 headings in HTML
func (fs *FirecrawlService) countHeadings(html string, result *models.CrawlResult) {
	// Count each heading level
	for i := 1; i <= 6; i++ {
		pattern := fmt.Sprintf(`(?i)<h%d[^>]*>`, i)
		regex := regexp.MustCompile(pattern)
		matches := regex.FindAllString(html, -1)
		count := len(matches)

		switch i {
		case 1:
			result.HeadingCounts.H1 = count
		case 2:
			result.HeadingCounts.H2 = count
		case 3:
			result.HeadingCounts.H3 = count
		case 4:
			result.HeadingCounts.H4 = count
		case 5:
			result.HeadingCounts.H5 = count
		case 6:
			result.HeadingCounts.H6 = count
		}
	}

	log.Printf("Heading counts: H1=%d, H2=%d, H3=%d, H4=%d, H5=%d, H6=%d",
		result.HeadingCounts.H1, result.HeadingCounts.H2, result.HeadingCounts.H3,
		result.HeadingCounts.H4, result.HeadingCounts.H5, result.HeadingCounts.H6)
}

// analyzeLinks analyzes links in the HTML content
func (fs *FirecrawlService) analyzeLinks(html string, result *models.CrawlResult) {
	// Simple link counting for now
	// In a production environment, you'd want more sophisticated link analysis

	// Count internal vs external links
	linkRegex := regexp.MustCompile(`(?i)href=["']([^"']+)["']`)
	matches := linkRegex.FindAllStringSubmatch(html, -1)

	internalCount := 0
	externalCount := 0

	for _, match := range matches {
		if len(match) > 1 {
			href := match[1]
			if strings.HasPrefix(href, "http") {
				externalCount++
			} else if strings.HasPrefix(href, "/") || strings.HasPrefix(href, "#") {
				internalCount++
			}
		}
	}

	result.InternalLinksCount = internalCount
	result.ExternalLinksCount = externalCount

	log.Printf("Link analysis: Internal=%d, External=%d", internalCount, externalCount)
}

// detectHTMLVersion detects HTML version from DOCTYPE or content
func (fs *FirecrawlService) detectHTMLVersion(html string) string {
	htmlUpper := strings.ToUpper(html)

	if strings.Contains(htmlUpper, "<!DOCTYPE HTML>") {
		return "HTML5"
	}
	if strings.Contains(htmlUpper, "HTML 4.01") {
		return "HTML 4.01"
	}
	if strings.Contains(htmlUpper, "XHTML 1.0") {
		return "XHTML 1.0"
	}
	if strings.Contains(htmlUpper, "XHTML 1.1") {
		return "XHTML 1.1"
	}

	// Default assumption for modern websites
	return "HTML5"
}

// extractFirecrawlMetadata extracts additional metadata from Firecrawl document metadata
func (fs *FirecrawlService) extractFirecrawlMetadata(metadata *firecrawl.FirecrawlDocumentMetadata, result *models.CrawlResult) {
	if metadata.StatusCode != nil && *metadata.StatusCode >= 400 {
		log.Printf("Warning: HTTP status code %d for URL: %s", *metadata.StatusCode, result.URL)
	}

	if metadata.Description != nil && *metadata.Description != "" {
		log.Printf("Page description: %s", *metadata.Description)
	}
}

// ValidateURL validates the URL format and content
func (fs *FirecrawlService) ValidateURL(targetURL string) error {
	if targetURL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		return fmt.Errorf("URL must start with http:// or https://")
	}

	// Basic validation for malicious patterns
	maliciousPatterns := []string{
		"javascript:",
		"data:",
		"file:",
		"ftp:",
	}

	lowerURL := strings.ToLower(targetURL)
	for _, pattern := range maliciousPatterns {
		if strings.Contains(lowerURL, pattern) {
			return fmt.Errorf("potentially malicious URL pattern detected")
		}
	}

	return nil
}
