-- Database schema

-- Create crawl_results table
CREATE TABLE IF NOT EXISTS crawl_results (
    id VARCHAR(36) PRIMARY KEY,
    url TEXT NOT NULL,
    title TEXT,
    html_version VARCHAR(50),
    internal_links_count INT DEFAULT 0,
    external_links_count INT DEFAULT 0,
    inaccessible_links_count INT DEFAULT 0,
    has_login_form BOOLEAN DEFAULT FALSE,
    heading_counts JSON,
    broken_links JSON,
    external_links JSON,
    status ENUM('queued', 'running', 'completed', 'error') NOT NULL DEFAULT 'queued',
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Add indexes directly in table creation
    INDEX idx_crawl_status (status),
    INDEX idx_crawl_created_at (created_at),
    INDEX idx_crawl_updated_at (updated_at),
    INDEX idx_crawl_url_hash (url(255)),
    INDEX idx_crawl_status_updated (status, updated_at),
    FULLTEXT KEY idx_url_title_fulltext (url, title)
);

SHOW TABLES; 