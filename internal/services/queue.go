package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"url-crawler/internal/config"
	"url-crawler/internal/models"

	"github.com/google/uuid"
)

// QueueService manages background crawling tasks
type QueueService struct {
	queue       chan *CrawlTask
	workers     int
	bufferSize  int
	maxRetries  int
	retryDelay  time.Duration
	crawler     Crawler
	storage     CrawlStorage
	running     bool
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	mu          sync.RWMutex
	activeTasks map[string]*CrawlTask
}

// CrawlTask represents a crawling task
type CrawlTask struct {
	ID        string
	URL       string
	CreatedAt time.Time
	Status    models.CrawlStatus
}

// CrawlStorage interface for persisting crawl results
type CrawlStorage interface {
	SaveCrawlResult(result *models.CrawlResult) error
	UpdateCrawlStatus(id string, status models.CrawlStatus, errorMsg *string) error
	GetCrawlResult(id string) (*models.CrawlResult, error)
}

// NewQueueService creates a new queue service (backward compatibility)
func NewQueueService(workers int, crawler Crawler, storage CrawlStorage) *QueueService {
	// Create default config
	defaultConfig := config.QueueConfig{
		Workers:    workers,
		BufferSize: 100,
		MaxRetries: 3,
		RetryDelay: 5 * time.Second,
	}
	return NewQueueServiceWithConfig(defaultConfig, crawler, storage)
}

// NewQueueServiceWithConfig creates a new queue service using configuration
func NewQueueServiceWithConfig(cfg config.QueueConfig, crawler Crawler, storage CrawlStorage) *QueueService {
	ctx, cancel := context.WithCancel(context.Background())

	return &QueueService{
		queue:       make(chan *CrawlTask, cfg.BufferSize),
		workers:     cfg.Workers,
		bufferSize:  cfg.BufferSize,
		maxRetries:  cfg.MaxRetries,
		retryDelay:  cfg.RetryDelay,
		crawler:     crawler,
		storage:     storage,
		ctx:         ctx,
		cancel:      cancel,
		activeTasks: make(map[string]*CrawlTask),
	}
}

// Start begins processing crawl tasks
func (q *QueueService) Start() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.running {
		return
	}

	q.running = true

	// Start worker goroutines
	for i := 0; i < q.workers; i++ {
		q.wg.Add(1)
		go q.worker(i)
	}

	log.Printf("Queue service started with %d workers", q.workers)
}

// Stop gracefully stops the queue service
func (q *QueueService) Stop() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.running {
		return
	}

	q.running = false
	q.cancel()
	close(q.queue)

	log.Println("Waiting for workers to finish...")
	q.wg.Wait()
	log.Println("Queue service stopped")
}

// EnqueueURL adds a URL to the crawling queue
func (q *QueueService) EnqueueURL(url string) (*models.CrawlResult, error) {
	// Validate URL
	if err := q.crawler.ValidateURL(url); err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Create crawl result record
	result := &models.CrawlResult{
		ID:            uuid.New().String(),
		URL:           url,
		Status:        models.CrawlStatusQueued,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		HeadingCounts: models.HeadingCounts{},
		BrokenLinks:   models.BrokenLinks{},
	}

	// Save initial record to database
	if err := q.storage.SaveCrawlResult(result); err != nil {
		return nil, fmt.Errorf("failed to save crawl result: %w", err)
	}

	// Create task
	task := &CrawlTask{
		ID:        result.ID,
		URL:       url,
		CreatedAt: time.Now(),
		Status:    models.CrawlStatusQueued,
	}

	// Add to active tasks
	q.mu.Lock()
	q.activeTasks[task.ID] = task
	q.mu.Unlock()

	// Try to enqueue (non-blocking)
	select {
	case q.queue <- task:
		log.Printf("Enqueued crawl task for URL: %s (ID: %s)", url, result.ID)
	default:
		// Queue is full
		q.mu.Lock()
		delete(q.activeTasks, task.ID)
		q.mu.Unlock()

		// Update status to error
		errorMsg := "Queue is full"
		q.storage.UpdateCrawlStatus(result.ID, models.CrawlStatusError, &errorMsg)

		return nil, fmt.Errorf("queue is full, please try again later")
	}

	return result, nil
}

// GetActiveTask returns an active task by ID
func (q *QueueService) GetActiveTask(id string) (*CrawlTask, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	task, exists := q.activeTasks[id]
	return task, exists
}

// GetQueueStats returns statistics about the queue
func (q *QueueService) GetQueueStats() map[string]interface{} {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return map[string]interface{}{
		"queue_length": len(q.queue),
		"active_tasks": len(q.activeTasks),
		"workers":      q.workers,
		"running":      q.running,
	}
}

// worker processes tasks from the queue
func (q *QueueService) worker(id int) {
	defer q.wg.Done()

	log.Printf("Worker %d started", id)

	for {
		select {
		case task, ok := <-q.queue:
			if !ok {
				log.Printf("Worker %d: Queue closed, exiting", id)
				return
			}

			q.processTask(task, id)

		case <-q.ctx.Done():
			log.Printf("Worker %d: Context cancelled, exiting", id)
			return
		}
	}
}

// processTask handles the actual crawling of a URL
func (q *QueueService) processTask(task *CrawlTask, workerID int) {
	log.Printf("Worker %d: Processing task %s for URL: %s", workerID, task.ID, task.URL)

	// Update task status to running
	task.Status = models.CrawlStatusRunning
	if err := q.storage.UpdateCrawlStatus(task.ID, models.CrawlStatusRunning, nil); err != nil {
		log.Printf("Worker %d: Failed to update task status to running: %v", workerID, err)
	}

	// Perform the actual crawling
	result, err := q.crawler.AnalyzeURL(task.URL)
	if err != nil {
		log.Printf("Worker %d: Failed to crawl URL %s: %v", workerID, task.URL, err)

		// Update status to error
		errorMsg := err.Error()
		if updateErr := q.storage.UpdateCrawlStatus(task.ID, models.CrawlStatusError, &errorMsg); updateErr != nil {
			log.Printf("Worker %d: Failed to update task status to error: %v", workerID, updateErr)
		}
	} else {
		// Update the result with the correct ID and save
		result.ID = task.ID
		result.Status = models.CrawlStatusCompleted
		result.UpdatedAt = time.Now()

		if err := q.storage.SaveCrawlResult(result); err != nil {
			log.Printf("Worker %d: Failed to save crawl result: %v", workerID, err)

			// Update status to error
			errorMsg := "Failed to save crawl result"
			q.storage.UpdateCrawlStatus(task.ID, models.CrawlStatusError, &errorMsg)
		} else {
			log.Printf("Worker %d: Successfully completed crawl for URL: %s", workerID, task.URL)
		}
	}

	// Remove from active tasks
	q.mu.Lock()
	delete(q.activeTasks, task.ID)
	q.mu.Unlock()
}

// RequeueTask re-adds a task to the queue (for re-running analysis)
func (q *QueueService) RequeueTask(id string) error {
	// Get the existing crawl result
	result, err := q.storage.GetCrawlResult(id)
	if err != nil {
		return fmt.Errorf("failed to get crawl result: %w", err)
	}

	// Create new task
	task := &CrawlTask{
		ID:        id,
		URL:       result.URL,
		CreatedAt: time.Now(),
		Status:    models.CrawlStatusQueued,
	}

	// Update status to queued
	if err := q.storage.UpdateCrawlStatus(id, models.CrawlStatusQueued, nil); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	// Add to active tasks
	q.mu.Lock()
	q.activeTasks[task.ID] = task
	q.mu.Unlock()

	select {
	case q.queue <- task:
		log.Printf("Re-queued crawl task for URL: %s (ID: %s)", result.URL, id)
		return nil
	default:
		// Queue is full
		q.mu.Lock()
		delete(q.activeTasks, task.ID)
		q.mu.Unlock()

		errorMsg := "Queue is full"
		q.storage.UpdateCrawlStatus(id, models.CrawlStatusError, &errorMsg)

		return fmt.Errorf("queue is full, please try again later")
	}
}
