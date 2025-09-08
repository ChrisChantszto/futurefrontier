package service

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/onetakesolutions/onetake-corpsite-backend/internal/config"
	"github.com/onetakesolutions/onetake-corpsite-backend/internal/models"
)

// LoggerService handles async logging with queue and backup mechanisms
type LoggerService struct {
	config       config.LoggingConfig
	esService    *ElasticsearchService
	mongoService *MongoBackupService
	logger       *zap.Logger

	// Channels for async processing
	apiLogQueue   chan *models.APILogEntry
	errorLogQueue chan *models.ErrorLogEntry

	// Worker management
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Metrics
	mu              sync.RWMutex
	processedLogs   int64
	failedLogs      int64
	queuedLogs      int64
}

// NewLoggerService creates a new logger service with async processing
func NewLoggerService(cfg config.LoggingConfig, database *mongo.Database, logger *zap.Logger) (*LoggerService, error) {
	// Initialize Elasticsearch service
	esService, err := NewElasticsearchService(cfg, logger)
	if err != nil {
		logger.Warn("Failed to initialize Elasticsearch service, will use MongoDB backup only", zap.Error(err))
	}

	// Initialize MongoDB backup service
	mongoService, err := NewMongoBackupService(database, logger)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	service := &LoggerService{
		config:        cfg,
		esService:     esService,
		mongoService:  mongoService,
		logger:        logger,
		apiLogQueue:   make(chan *models.APILogEntry, cfg.QueueSize),
		errorLogQueue: make(chan *models.ErrorLogEntry, cfg.QueueSize),
		ctx:           ctx,
		cancel:        cancel,
	}

	// Start worker goroutines
	service.startWorkers()

	// Start recovery goroutine
	service.startRecoveryWorker()

	logger.Info("Logger service started", 
		zap.Int("workers", cfg.WorkerCount),
		zap.Int("queue_size", cfg.QueueSize),
		zap.Bool("local_mode", cfg.LocalMode))

	return service, nil
}

// startWorkers starts the background worker goroutines
func (ls *LoggerService) startWorkers() {
	for i := 0; i < ls.config.WorkerCount; i++ {
		ls.wg.Add(2) // One for API logs, one for error logs

		// API log worker
		go func(workerID int) {
			defer ls.wg.Done()
			ls.logger.Debug("Starting API log worker", zap.Int("worker_id", workerID))

			for {
				select {
				case <-ls.ctx.Done():
					return
				case logEntry := <-ls.apiLogQueue:
					ls.processAPILog(logEntry)
				}
			}
		}(i)

		// Error log worker
		go func(workerID int) {
			defer ls.wg.Done()
			ls.logger.Debug("Starting error log worker", zap.Int("worker_id", workerID))

			for {
				select {
				case <-ls.ctx.Done():
					return
				case logEntry := <-ls.errorLogQueue:
					ls.processErrorLog(logEntry)
				}
			}
		}(i)
	}
}

// startRecoveryWorker starts the recovery worker that retries failed logs from MongoDB
func (ls *LoggerService) startRecoveryWorker() {
	ls.wg.Add(1)
	go func() {
		defer ls.wg.Done()
		ticker := time.NewTicker(time.Duration(ls.config.RetryInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ls.ctx.Done():
				return
			case <-ticker.C:
				ls.recoverFailedLogs()
			}
		}
	}()
}

// LogAPI queues an API log entry for async processing
func (ls *LoggerService) LogAPI(logEntry *models.APILogEntry) {
	logEntry.CreatedAt = time.Now()
	logEntry.ESWritten = false

	select {
	case ls.apiLogQueue <- logEntry:
		ls.incrementQueuedLogs()
	default:
		// Queue is full, log directly to MongoDB as backup
		ls.logger.Warn("API log queue is full, writing directly to MongoDB backup")
		if err := ls.mongoService.SaveAPILog(ls.ctx, logEntry); err != nil {
			ls.logger.Error("Failed to save API log to MongoDB backup", zap.Error(err))
		}
		ls.incrementFailedLogs()
	}
}

// LogError queues an error log entry for async processing
func (ls *LoggerService) LogError(logEntry *models.ErrorLogEntry) {
	logEntry.CreatedAt = time.Now()
	logEntry.ESWritten = false

	select {
	case ls.errorLogQueue <- logEntry:
		ls.incrementQueuedLogs()
	default:
		// Queue is full, log directly to MongoDB as backup
		ls.logger.Warn("Error log queue is full, writing directly to MongoDB backup")
		if err := ls.mongoService.SaveErrorLog(ls.ctx, logEntry); err != nil {
			ls.logger.Error("Failed to save error log to MongoDB backup", zap.Error(err))
		}
		ls.incrementFailedLogs()
	}
}

// processAPILog processes a single API log entry
func (ls *LoggerService) processAPILog(logEntry *models.APILogEntry) {
	ctx, cancel := context.WithTimeout(ls.ctx, 10*time.Second)
	defer cancel()

	// Try to write to Elasticsearch first
	if ls.esService != nil {
		err := ls.esService.WriteAPILog(ctx, logEntry)
		if err == nil {
			logEntry.ESWritten = true
			ls.incrementProcessedLogs()
			return
		}

		ls.logger.Warn("Failed to write API log to Elasticsearch, saving to MongoDB backup", zap.Error(err))
	}

	// Fallback to MongoDB
	if err := ls.mongoService.SaveAPILog(ctx, logEntry); err != nil {
		ls.logger.Error("Failed to save API log to MongoDB backup", zap.Error(err))
		ls.incrementFailedLogs()
		return
	}

	ls.incrementProcessedLogs()
}

// processErrorLog processes a single error log entry
func (ls *LoggerService) processErrorLog(logEntry *models.ErrorLogEntry) {
	ctx, cancel := context.WithTimeout(ls.ctx, 10*time.Second)
	defer cancel()

	// Try to write to Elasticsearch first
	if ls.esService != nil {
		err := ls.esService.WriteErrorLog(ctx, logEntry)
		if err == nil {
			logEntry.ESWritten = true
			ls.incrementProcessedLogs()
			return
		}

		ls.logger.Warn("Failed to write error log to Elasticsearch, saving to MongoDB backup", zap.Error(err))
	}

	// Fallback to MongoDB
	if err := ls.mongoService.SaveErrorLog(ctx, logEntry); err != nil {
		ls.logger.Error("Failed to save error log to MongoDB backup", zap.Error(err))
		ls.incrementFailedLogs()
		return
	}

	ls.incrementProcessedLogs()
}

// recoverFailedLogs attempts to retry failed logs from MongoDB to Elasticsearch
func (ls *LoggerService) recoverFailedLogs() {
	if ls.esService == nil || !ls.esService.IsHealthy(ls.ctx) {
		return
	}

	ctx, cancel := context.WithTimeout(ls.ctx, 30*time.Second)
	defer cancel()

	// Recover API logs
	apiLogs, err := ls.mongoService.GetFailedAPILogs(ctx, 100) // Process 100 at a time
	if err != nil {
		ls.logger.Error("Failed to get failed API logs from MongoDB", zap.Error(err))
	} else {
		for _, logEntry := range apiLogs {
			if err := ls.esService.WriteAPILog(ctx, &logEntry); err == nil {
				// Successfully written to ES, remove from MongoDB
				if err := ls.mongoService.DeleteAPILog(ctx, logEntry.ID); err != nil {
					ls.logger.Error("Failed to delete recovered API log from MongoDB", zap.Error(err))
				}
			} else {
				// Update retry count
				if err := ls.mongoService.UpdateAPILogRetry(ctx, logEntry.ID); err != nil {
					ls.logger.Error("Failed to update API log retry count", zap.Error(err))
				}
			}
		}
	}

	// Recover error logs
	errorLogs, err := ls.mongoService.GetFailedErrorLogs(ctx, 100)
	if err != nil {
		ls.logger.Error("Failed to get failed error logs from MongoDB", zap.Error(err))
	} else {
		for _, logEntry := range errorLogs {
			if err := ls.esService.WriteErrorLog(ctx, &logEntry); err == nil {
				// Successfully written to ES, remove from MongoDB
				if err := ls.mongoService.DeleteErrorLog(ctx, logEntry.ID); err != nil {
					ls.logger.Error("Failed to delete recovered error log from MongoDB", zap.Error(err))
				}
			} else {
				// Update retry count
				if err := ls.mongoService.UpdateErrorLogRetry(ctx, logEntry.ID); err != nil {
					ls.logger.Error("Failed to update error log retry count", zap.Error(err))
				}
			}
		}
	}
}

// GetStats returns logging statistics
func (ls *LoggerService) GetStats() (processed, failed, queued int64) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()
	return ls.processedLogs, ls.failedLogs, ls.queuedLogs
}

// Shutdown gracefully shuts down the logger service
func (ls *LoggerService) Shutdown() {
	ls.logger.Info("Shutting down logger service...")
	
	ls.cancel()
	
	// Close channels
	close(ls.apiLogQueue)
	close(ls.errorLogQueue)
	
	// Wait for workers to finish
	ls.wg.Wait()
	
	ls.logger.Info("Logger service shutdown complete")
}

// Helper methods for metrics
func (ls *LoggerService) incrementProcessedLogs() {
	ls.mu.Lock()
	ls.processedLogs++
	ls.mu.Unlock()
}

func (ls *LoggerService) incrementFailedLogs() {
	ls.mu.Lock()
	ls.failedLogs++
	ls.mu.Unlock()
}

func (ls *LoggerService) incrementQueuedLogs() {
	ls.mu.Lock()
	ls.queuedLogs++
	ls.mu.Unlock()
}
