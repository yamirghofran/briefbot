# BriefBot Monitoring Implementation Specification

## Table of Contents
1. [Metrics Package Design](#metrics-package-design)
2. [Middleware Implementation](#middleware-implementation)
3. [Service Instrumentation](#service-instrumentation)
4. [Dashboard Specifications](#dashboard-specifications)
5. [Code Examples](#code-examples)

## Metrics Package Design

### Package Structure
```
backend/internal/metrics/
├── metrics.go          # Main metrics definitions
├── http.go            # HTTP-specific metrics
├── jobs.go            # Job queue metrics
├── workers.go         # Worker metrics
├── external.go        # External service metrics
└── database.go        # Database metrics
```

### Metric Naming Convention
Follow Prometheus best practices:
- Use snake_case
- Include application prefix: `briefbot_`
- Counter suffix: `_total`
- Histogram/Summary suffix: `_duration_seconds`, `_bytes`
- Gauge: no suffix

### Core Metrics Definitions

#### HTTP Metrics (`http.go`)
```go
var (
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "briefbot_http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )

    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "briefbot_http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
        },
        []string{"method", "endpoint"},
    )

    httpRequestSize = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "briefbot_http_request_size_bytes",
            Help:    "HTTP request size in bytes",
            Buckets: prometheus.ExponentialBuckets(100, 10, 8),
        },
        []string{"method", "endpoint"},
    )

    httpResponseSize = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "briefbot_http_response_size_bytes",
            Help:    "HTTP response size in bytes",
            Buckets: prometheus.ExponentialBuckets(100, 10, 8),
        },
        []string{"method", "endpoint"},
    )

    httpActiveRequests = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "briefbot_http_active_requests",
            Help: "Number of active HTTP requests",
        },
    )
)
```

#### Job Queue Metrics (`jobs.go`)
```go
var (
    jobsEnqueuedTotal = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "briefbot_jobs_enqueued_total",
            Help: "Total number of jobs enqueued",
        },
    )

    jobsProcessingCurrent = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "briefbot_jobs_processing_current",
            Help: "Number of jobs currently being processed",
        },
    )

    jobsCompletedTotal = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "briefbot_jobs_completed_total",
            Help: "Total number of jobs completed successfully",
        },
    )

    jobsFailedTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "briefbot_jobs_failed_total",
            Help: "Total number of jobs that failed",
        },
        []string{"error_type"},
    )

    jobProcessingDuration = promauto.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "briefbot_job_processing_duration_seconds",
            Help:    "Job processing duration in seconds",
            Buckets: []float64{1, 5, 10, 30, 60, 120, 300, 600},
        },
    )

    jobQueueDepth = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "briefbot_job_queue_depth",
            Help: "Number of jobs in queue by status",
        },
        []string{"status"},
    )

    jobRetriesTotal = promauto.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "briefbot_job_retries_total",
            Help:    "Distribution of job retry counts",
            Buckets: []float64{0, 1, 2, 3, 4, 5},
        },
    )
)
```

#### Worker Metrics (`workers.go`)
```go
var (
    workersActive = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "briefbot_workers_active",
            Help: "Number of active workers",
        },
    )

    workerJobsProcessedTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "briefbot_worker_jobs_processed_total",
            Help: "Total number of jobs processed by each worker",
        },
        []string{"worker_id"},
    )

    workerErrorsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "briefbot_worker_errors_total",
            Help: "Total number of worker errors",
        },
        []string{"worker_id", "error_type"},
    )

    workerBatchDuration = promauto.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "briefbot_worker_batch_duration_seconds",
            Help:    "Worker batch processing duration in seconds",
            Buckets: []float64{1, 5, 10, 30, 60, 120, 300},
        },
    )

    workerUtilization = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "briefbot_worker_utilization",
            Help: "Worker utilization percentage (0-100)",
        },
        []string{"worker_id"},
    )
)
```

#### Podcast Metrics (`external.go`)
```go
var (
    podcastsGeneratedTotal = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "briefbot_podcasts_generated_total",
            Help: "Total number of podcasts generated",
        },
    )

    podcastGenerationDuration = promauto.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "briefbot_podcast_generation_duration_seconds",
            Help:    "Podcast generation duration in seconds",
            Buckets: []float64{10, 30, 60, 120, 300, 600, 1200},
        },
    )

    podcastAudioRequestsCurrent = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "briefbot_podcast_audio_requests_current",
            Help: "Number of concurrent audio generation requests",
        },
    )

    podcastGenerationFailuresTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "briefbot_podcast_generation_failures_total",
            Help: "Total number of podcast generation failures",
        },
        []string{"stage"},
    )
)
```

#### AI Service Metrics (`external.go`)
```go
var (
    aiAPICallsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "briefbot_ai_api_calls_total",
            Help: "Total number of AI API calls",
        },
        []string{"operation", "provider"},
    )

    aiAPILatency = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "briefbot_ai_api_latency_seconds",
            Help:    "AI API call latency in seconds",
            Buckets: []float64{.1, .5, 1, 2, 5, 10, 30},
        },
        []string{"operation", "provider"},
    )

    aiAPIErrorsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "briefbot_ai_api_errors_total",
            Help: "Total number of AI API errors",
        },
        []string{"operation", "provider", "error_type"},
    )

    aiTokensUsedTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "briefbot_ai_tokens_used_total",
            Help: "Total number of AI tokens used",
        },
        []string{"operation", "provider"},
    )
)
```

#### Database Metrics (`database.go`)
```go
var (
    dbConnectionsActive = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "briefbot_db_connections_active",
            Help: "Number of active database connections",
        },
    )

    dbConnectionsIdle = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "briefbot_db_connections_idle",
            Help: "Number of idle database connections",
        },
    )

    dbConnectionsWaiting = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "briefbot_db_connections_waiting",
            Help: "Number of connections waiting for a database connection",
        },
    )

    dbQueryDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "briefbot_db_query_duration_seconds",
            Help:    "Database query duration in seconds",
            Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
        },
        []string{"query_type"},
    )

    dbErrorsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "briefbot_db_errors_total",
            Help: "Total number of database errors",
        },
        []string{"error_type"},
    )
)
```

#### External Service Metrics (`external.go`)
```go
var (
    // R2 Storage
    r2OperationsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "briefbot_r2_operations_total",
            Help: "Total number of R2 operations",
        },
        []string{"operation"},
    )

    r2OperationDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "briefbot_r2_operation_duration_seconds",
            Help:    "R2 operation duration in seconds",
            Buckets: []float64{.1, .5, 1, 2, 5, 10, 30},
        },
        []string{"operation"},
    )

    // Speech Service (FAL)
    speechAPIRequestsTotal = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "briefbot_speech_api_requests_total",
            Help: "Total number of speech API requests",
        },
    )

    speechAPILatency = promauto.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "briefbot_speech_api_latency_seconds",
            Help:    "Speech API latency in seconds",
            Buckets: []float64{1, 5, 10, 30, 60, 120, 300},
        },
    )

    // Email Service
    emailsSentTotal = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "briefbot_emails_sent_total",
            Help: "Total number of emails sent",
        },
    )

    emailFailuresTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "briefbot_email_failures_total",
            Help: "Total number of email failures",
        },
        []string{"error_type"},
    )

    // Scraping Service
    scrapingRequestsTotal = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "briefbot_scraping_requests_total",
            Help: "Total number of scraping requests",
        },
    )

    scrapingFailuresTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "briefbot_scraping_failures_total",
            Help: "Total number of scraping failures",
        },
        []string{"error_type"},
    )
)
```

#### SSE Metrics (`external.go`)
```go
var (
    sseConnectionsActive = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "briefbot_sse_connections_active",
            Help: "Number of active SSE connections",
        },
    )

    sseEventsSentTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "briefbot_sse_events_sent_total",
            Help: "Total number of SSE events sent",
        },
        []string{"event_type"},
    )

    sseConnectionDuration = promauto.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "briefbot_sse_connection_duration_seconds",
            Help:    "SSE connection duration in seconds",
            Buckets: []float64{10, 30, 60, 300, 600, 1800, 3600},
        },
    )
)
```

## Middleware Implementation

### Prometheus Middleware (`backend/internal/middleware/prometheus.go`)

```go
package middleware

import (
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/yamirghofran/briefbot/internal/metrics"
)

// PrometheusMiddleware returns a Gin middleware that records HTTP metrics
func PrometheusMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Skip metrics endpoint to avoid recursion
        if c.Request.URL.Path == "/metrics" {
            c.Next()
            return
        }

        start := time.Now()
        
        // Increment active requests
        metrics.IncrementActiveRequests()
        defer metrics.DecrementActiveRequests()

        // Record request size
        requestSize := computeApproximateRequestSize(c.Request)
        
        // Process request
        c.Next()

        // Calculate duration
        duration := time.Since(start).Seconds()
        
        // Get status code
        status := strconv.Itoa(c.Writer.Status())
        
        // Get endpoint (use route pattern, not actual path)
        endpoint := c.FullPath()
        if endpoint == "" {
            endpoint = "unknown"
        }
        
        // Record metrics
        metrics.RecordHTTPRequest(c.Request.Method, endpoint, status)
        metrics.RecordHTTPDuration(c.Request.Method, endpoint, duration)
        metrics.RecordHTTPRequestSize(c.Request.Method, endpoint, float64(requestSize))
        metrics.RecordHTTPResponseSize(c.Request.Method, endpoint, float64(c.Writer.Size()))
    }
}

func computeApproximateRequestSize(r *http.Request) int {
    size := 0
    if r.URL != nil {
        size += len(r.URL.String())
    }
    size += len(r.Method)
    size += len(r.Proto)
    for name, values := range r.Header {
        size += len(name)
        for _, value := range values {
            size += len(value)
        }
    }
    size += len(r.Host)
    if r.ContentLength != -1 {
        size += int(r.ContentLength)
    }
    return size
}
```

### Metrics Helper Functions (`backend/internal/metrics/metrics.go`)

```go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

// HTTP Metrics helpers
func RecordHTTPRequest(method, endpoint, status string) {
    httpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
}

func RecordHTTPDuration(method, endpoint string, duration float64) {
    httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
}

func RecordHTTPRequestSize(method, endpoint string, size float64) {
    httpRequestSize.WithLabelValues(method, endpoint).Observe(size)
}

func RecordHTTPResponseSize(method, endpoint string, size float64) {
    httpResponseSize.WithLabelValues(method, endpoint).Observe(size)
}

func IncrementActiveRequests() {
    httpActiveRequests.Inc()
}

func DecrementActiveRequests() {
    httpActiveRequests.Dec()
}

// Job Queue Metrics helpers
func IncrementJobsEnqueued() {
    jobsEnqueuedTotal.Inc()
}

func IncrementJobsProcessing() {
    jobsProcessingCurrent.Inc()
}

func DecrementJobsProcessing() {
    jobsProcessingCurrent.Dec()
}

func IncrementJobsCompleted() {
    jobsCompletedTotal.Inc()
}

func IncrementJobsFailed(errorType string) {
    jobsFailedTotal.WithLabelValues(errorType).Inc()
}

func RecordJobProcessingDuration(duration float64) {
    jobProcessingDuration.Observe(duration)
}

func SetJobQueueDepth(status string, count float64) {
    jobQueueDepth.WithLabelValues(status).Set(count)
}

// Worker Metrics helpers
func SetWorkersActive(count float64) {
    workersActive.Set(count)
}

func IncrementWorkerJobsProcessed(workerID string) {
    workerJobsProcessedTotal.WithLabelValues(workerID).Inc()
}

func IncrementWorkerErrors(workerID, errorType string) {
    workerErrorsTotal.WithLabelValues(workerID, errorType).Inc()
}

func RecordWorkerBatchDuration(duration float64) {
    workerBatchDuration.Observe(duration)
}

// Podcast Metrics helpers
func IncrementPodcastsGenerated() {
    podcastsGeneratedTotal.Inc()
}

func RecordPodcastGenerationDuration(duration float64) {
    podcastGenerationDuration.Observe(duration)
}

func IncrementPodcastAudioRequests() {
    podcastAudioRequestsCurrent.Inc()
}

func DecrementPodcastAudioRequests() {
    podcastAudioRequestsCurrent.Dec()
}

func IncrementPodcastGenerationFailures(stage string) {
    podcastGenerationFailuresTotal.WithLabelValues(stage).Inc()
}

// AI Service Metrics helpers
func RecordAIAPICall(operation, provider string, duration float64) {
    aiAPICallsTotal.WithLabelValues(operation, provider).Inc()
    aiAPILatency.WithLabelValues(operation, provider).Observe(duration)
}

func IncrementAIAPIErrors(operation, provider, errorType string) {
    aiAPIErrorsTotal.WithLabelValues(operation, provider, errorType).Inc()
}

// Database Metrics helpers
func UpdateDBConnectionStats(active, idle, waiting int32) {
    dbConnectionsActive.Set(float64(active))
    dbConnectionsIdle.Set(float64(idle))
    dbConnectionsWaiting.Set(float64(waiting))
}

// External Service Metrics helpers
func RecordR2Operation(operation string, duration float64) {
    r2OperationsTotal.WithLabelValues(operation).Inc()
    r2OperationDuration.WithLabelValues(operation).Observe(duration)
}

func RecordSpeechAPIRequest(duration float64) {
    speechAPIRequestsTotal.Inc()
    speechAPILatency.Observe(duration)
}

func IncrementEmailsSent() {
    emailsSentTotal.Inc()
}

func IncrementEmailFailures(errorType string) {
    emailFailuresTotal.WithLabelValues(errorType).Inc()
}

func IncrementScrapingRequests() {
    scrapingRequestsTotal.Inc()
}

func IncrementScrapingFailures(errorType string) {
    scrapingFailuresTotal.WithLabelValues(errorType).Inc()
}

// SSE Metrics helpers
func IncrementSSEConnections() {
    sseConnectionsActive.Inc()
}

func DecrementSSEConnections() {
    sseConnectionsActive.Dec()
}

func IncrementSSEEventsSent(eventType string) {
    sseEventsSentTotal.WithLabelValues(eventType).Inc()
}

func RecordSSEConnectionDuration(duration float64) {
    sseConnectionDuration.Observe(duration)
}
```

## Service Instrumentation Examples

### Job Queue Service (`jobqueue.go`)

```go
// In EnqueueItem method
func (s *jobQueueService) EnqueueItem(ctx context.Context, userID int32, title string, url string) (*db.Item, error) {
    // ... existing code ...
    
    item, err := s.querier.CreatePendingItem(ctx, params)
    if err != nil {
        return nil, fmt.Errorf("failed to enqueue item: %w", err)
    }

    // Add metrics
    metrics.IncrementJobsEnqueued()

    // ... rest of existing code ...
}

// In MarkItemAsProcessing method
func (s *jobQueueService) MarkItemAsProcessing(ctx context.Context, itemID int32) error {
    // ... existing code ...
    
    err = s.querier.UpdateItemAsProcessing(ctx, itemID)
    if err != nil {
        return fmt.Errorf("failed to mark item as processing: %w", err)
    }

    // Add metrics
    metrics.IncrementJobsProcessing()

    // ... rest of existing code ...
}

// In CompleteItem method
func (s *jobQueueService) CompleteItem(ctx context.Context, itemID int32, ...) error {
    start := time.Now()
    defer func() {
        duration := time.Since(start).Seconds()
        metrics.RecordJobProcessingDuration(duration)
    }()

    // ... existing code ...
    
    if err != nil {
        return fmt.Errorf("failed to complete item: %w", err)
    }

    // Add metrics
    metrics.DecrementJobsProcessing()
    metrics.IncrementJobsCompleted()

    // ... rest of existing code ...
}

// In FailItem method
func (s *jobQueueService) FailItem(ctx context.Context, itemID int32, errorMsg string) error {
    // ... existing code ...
    
    // Add metrics
    metrics.DecrementJobsProcessing()
    metrics.IncrementJobsFailed(categorizeError(errorMsg))

    // ... rest of existing code ...
}

// Helper function to categorize errors
func categorizeError(errorMsg string) string {
    if strings.Contains(errorMsg, "network") {
        return "network"
    } else if strings.Contains(errorMsg, "timeout") {
        return "timeout"
    } else if strings.Contains(errorMsg, "parse") {
        return "parse"
    }
    return "unknown"
}
```

### Worker Service (`worker.go`)

```go
// In Start method
func (s *workerService) Start(ctx context.Context) error {
    // ... existing code ...
    
    // Start worker goroutines
    for i := 0; i < s.workerCount; i++ {
        s.wg.Add(1)
        go s.worker(i + 1)
    }

    // Update metrics
    metrics.SetWorkersActive(float64(s.workerCount))

    // ... rest of existing code ...
}

// In worker method
func (s *workerService) worker(workerID int) {
    defer s.wg.Done()
    
    workerIDStr := fmt.Sprintf("worker-%d", workerID)
    log.Printf("[Worker %d] Starting", workerID)

    for {
        select {
        case <-s.ctx.Done():
            log.Printf("[Worker %d] Shutting down", workerID)
            return
        case <-time.After(s.pollInterval):
            // Process batch
            batchStart := time.Now()
            
            items, err := s.jobQueueService.DequeuePendingItems(s.ctx, s.batchSize)
            if err != nil {
                log.Printf("[Worker %d] Error dequeuing items: %v", workerID, err)
                metrics.IncrementWorkerErrors(workerIDStr, "dequeue")
                continue
            }

            if len(items) == 0 {
                continue
            }

            // Process each item
            for _, item := range items {
                if err := s.processItem(s.ctx, item); err != nil {
                    log.Printf("[Worker %d] Error processing item %d: %v", workerID, item.ID, err)
                    metrics.IncrementWorkerErrors(workerIDStr, "process")
                } else {
                    metrics.IncrementWorkerJobsProcessed(workerIDStr)
                }
            }

            // Record batch duration
            batchDuration := time.Since(batchStart).Seconds()
            metrics.RecordWorkerBatchDuration(batchDuration)
        }
    }
}
```

### Podcast Service (`podcast.go`)

```go
// In GeneratePodcast method
func (s *podcastService) GeneratePodcast(ctx context.Context, userID int32, itemIDs []int32) (*db.Podcast, error) {
    start := time.Now()
    defer func() {
        duration := time.Since(start).Seconds()
        metrics.RecordPodcastGenerationDuration(duration)
    }()

    // ... existing code ...

    // Generate audio for each segment
    for i, segment := range dialogSegments {
        metrics.IncrementPodcastAudioRequests()
        defer metrics.DecrementPodcastAudioRequests()

        audioData, err := s.speechService.GenerateSpeech(ctx, segment.Text, segment.Voice)
        if err != nil {
            metrics.IncrementPodcastGenerationFailures("audio_generation")
            return nil, fmt.Errorf("failed to generate audio for segment %d: %w", i, err)
        }
        // ... rest of code ...
    }

    // ... existing code ...

    // Upload to R2
    if s.r2Service != nil {
        audioURL, err := s.r2Service.UploadFile(ctx, ...)
        if err != nil {
            metrics.IncrementPodcastGenerationFailures("upload")
            return nil, fmt.Errorf("failed to upload podcast: %w", err)
        }
    }

    metrics.IncrementPodcastsGenerated()
    return podcast, nil
}
```

### AI Service (`ai.go`)

```go
// In Summarize method
func (s *aiService) Summarize(ctx context.Context, content string) (string, error) {
    start := time.Now()
    defer func() {
        duration := time.Since(start).Seconds()
        metrics.RecordAIAPICall("summarize", "groq", duration)
    }()

    // ... existing code ...

    if err != nil {
        metrics.IncrementAIAPIErrors("summarize", "groq", categorizeAIError(err))
        return "", fmt.Errorf("failed to summarize: %w", err)
    }

    // ... rest of existing code ...
}

// Similar instrumentation for ExtractMetadata, GeneratePodcastDialog, etc.
```

### Database Pool Monitoring (`main.go`)

```go
// In main function, after creating pool
func main() {
    // ... existing code ...

    // Create connection pool
    pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
    if err != nil {
        log.Fatalf("Unable to create connection pool: %v", err)
    }
    defer pool.Close()

    // Start database metrics collector
    go func() {
        ticker := time.NewTicker(15 * time.Second)
        defer ticker.Stop()
        
        for range ticker.C {
            stats := pool.Stat()
            metrics.UpdateDBConnectionStats(
                stats.AcquiredConns(),
                stats.IdleConns(),
                stats.EmptyAcquireCount(),
            )
        }
    }()

    // ... rest of existing code ...
}
```

## Dashboard Specifications

### Overview Dashboard Structure

**Panels (Top to Bottom, Left to Right):**

1. **Row 1: Key Metrics**
   - Request Rate (Graph): `rate(briefbot_http_requests_total[5m])`
   - Error Rate (Graph): `rate(briefbot_http_requests_total{status=~"5.."}[5m])`
   - P95 Latency (Graph): `histogram_quantile(0.95, rate(briefbot_http_request_duration_seconds_bucket[5m]))`
   - Active Requests (Stat): `briefbot_http_active_requests`

2. **Row 2: Job Processing**
   - Jobs Enqueued (Graph): `rate(briefbot_jobs_enqueued_total[5m])`
   - Jobs Completed (Graph): `rate(briefbot_jobs_completed_total[5m])`
   - Jobs Failed (Graph): `rate(briefbot_jobs_failed_total[5m])`
   - Queue Depth (Graph): `briefbot_job_queue_depth`

3. **Row 3: Workers**
   - Active Workers (Stat): `briefbot_workers_active`
   - Worker Utilization (Gauge): `briefbot_worker_utilization`
   - Jobs per Worker (Graph): `rate(briefbot_worker_jobs_processed_total[5m])`

4. **Row 4: Database**
   - DB Connections (Graph): `briefbot_db_connections_active`, `briefbot_db_connections_idle`
   - DB Query Duration (Heatmap): `briefbot_db_query_duration_seconds`

5. **Row 5: External Services**
   - AI API Calls (Graph): `rate(briefbot_ai_api_calls_total[5m])`
   - R2 Operations (Graph): `rate(briefbot_r2_operations_total[5m])`
   - Email Sent (Graph): `rate(briefbot_emails_sent_total[5m])`

### Query Examples

```promql
# HTTP request rate by endpoint
sum(rate(briefbot_http_requests_total[5m])) by (endpoint)

# Error rate percentage
sum(rate(briefbot_http_requests_total{status=~"5.."}[5m])) / sum(rate(briefbot_http_requests_total[5m])) * 100

# P50, P95, P99 latency
histogram_quantile(0.50, rate(briefbot_http_request_duration_seconds_bucket[5m]))
histogram_quantile(0.95, rate(briefbot_http_request_duration_seconds_bucket[5m]))
histogram_quantile(0.99, rate(briefbot_http_request_duration_seconds_bucket[5m]))

# Job processing throughput
rate(briefbot_jobs_completed_total[5m])

# Job failure rate
rate(briefbot_jobs_failed_total[5m]) / rate(briefbot_jobs_enqueued_total[5m]) * 100

# Database connection pool utilization
briefbot_db_connections_active / (briefbot_db_connections_active + briefbot_db_connections_idle) * 100

# AI API success rate
(rate(briefbot_ai_api_calls_total[5m]) - rate(briefbot_ai_api_errors_total[5m])) / rate(briefbot_ai_api_calls_total[5m]) * 100
```

## Testing Checklist

- [ ] Metrics endpoint accessible at `/metrics`
- [ ] HTTP metrics recorded for all endpoints
- [ ] Job queue metrics update correctly
- [ ] Worker metrics track active workers
- [ ] Podcast generation metrics recorded
- [ ] AI service metrics captured
- [ ] Database pool metrics updated
- [ ] External service metrics tracked
- [ ] SSE metrics recorded
- [ ] Prometheus scrapes metrics successfully
- [ ] Grafana displays all dashboards
- [ ] Load test shows accurate metrics
- [ ] No performance degradation

## Performance Considerations

1. **Metric Cardinality**: Keep label values bounded
   - Use `endpoint` (route pattern) not `path` (actual URL)
   - Categorize error types instead of unique messages
   - Limit worker IDs to actual worker count

2. **Sampling**: For high-frequency operations, consider sampling
   - Example: Sample 10% of requests for detailed tracing

3. **Aggregation**: Pre-aggregate where possible
   - Use counters and let Prometheus calculate rates

4. **Memory**: Monitor Prometheus memory usage
   - Default retention: 30 days
   - Adjust based on metric volume

