package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP Metrics
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

	// Job Queue Metrics
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

	// Worker Metrics
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

	// Database Metrics
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

	// Podcast Metrics
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

	// AI Service Metrics
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

	// External Service Metrics
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

	// SSE Metrics
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

// Database Metrics helpers
func UpdateDBConnectionStats(active, idle, waiting int32) {
	dbConnectionsActive.Set(float64(active))
	dbConnectionsIdle.Set(float64(idle))
	dbConnectionsWaiting.Set(float64(waiting))
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

// External Service Metrics helpers
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
