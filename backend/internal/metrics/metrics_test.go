package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

// HTTP Metrics Tests

func TestRecordHTTPRequest(t *testing.T) {
	before := testutil.ToFloat64(httpRequestsTotal.WithLabelValues("GET", "/test", "200"))
	RecordHTTPRequest("GET", "/test", "200")
	after := testutil.ToFloat64(httpRequestsTotal.WithLabelValues("GET", "/test", "200"))

	if after != before+1 {
		t.Errorf("Expected counter to increment by 1, got %f", after-before)
	}
}

func TestRecordHTTPDuration(t *testing.T) {
	// Record a duration observation
	RecordHTTPDuration("POST", "/api/test", 1.5)

	// Verify the histogram recorded the observation by checking the count
	count := testutil.CollectAndCount(httpRequestDuration)
	if count == 0 {
		t.Error("Expected histogram to have metrics registered")
	}
}

func TestRecordHTTPRequestSize(t *testing.T) {
	// Record a request size observation
	RecordHTTPRequestSize("PUT", "/upload", 1024.0)

	// Verify the histogram recorded the observation
	count := testutil.CollectAndCount(httpRequestSize)
	if count == 0 {
		t.Error("Expected histogram to have metrics registered")
	}
}

func TestRecordHTTPResponseSize(t *testing.T) {
	// Record a response size observation
	RecordHTTPResponseSize("GET", "/download", 2048.0)

	// Verify the histogram recorded the observation
	count := testutil.CollectAndCount(httpResponseSize)
	if count == 0 {
		t.Error("Expected histogram to have metrics registered")
	}
}

func TestIncrementActiveRequests(t *testing.T) {
	before := testutil.ToFloat64(httpActiveRequests)
	IncrementActiveRequests()
	after := testutil.ToFloat64(httpActiveRequests)

	if after != before+1 {
		t.Errorf("Expected gauge to increment by 1, got %f", after-before)
	}
}

func TestDecrementActiveRequests(t *testing.T) {
	before := testutil.ToFloat64(httpActiveRequests)
	DecrementActiveRequests()
	after := testutil.ToFloat64(httpActiveRequests)

	if after != before-1 {
		t.Errorf("Expected gauge to decrement by 1, got %f", after-before)
	}
}

// Job Queue Metrics Tests

func TestIncrementJobsEnqueued(t *testing.T) {
	before := testutil.ToFloat64(jobsEnqueuedTotal)
	IncrementJobsEnqueued()
	after := testutil.ToFloat64(jobsEnqueuedTotal)

	if after != before+1 {
		t.Errorf("Expected counter to increment by 1, got %f", after-before)
	}
}

func TestIncrementJobsProcessing(t *testing.T) {
	before := testutil.ToFloat64(jobsProcessingCurrent)
	IncrementJobsProcessing()
	after := testutil.ToFloat64(jobsProcessingCurrent)

	if after != before+1 {
		t.Errorf("Expected gauge to increment by 1, got %f", after-before)
	}
}

func TestDecrementJobsProcessing(t *testing.T) {
	before := testutil.ToFloat64(jobsProcessingCurrent)
	DecrementJobsProcessing()
	after := testutil.ToFloat64(jobsProcessingCurrent)

	if after != before-1 {
		t.Errorf("Expected gauge to decrement by 1, got %f", after-before)
	}
}

func TestIncrementJobsCompleted(t *testing.T) {
	before := testutil.ToFloat64(jobsCompletedTotal)
	IncrementJobsCompleted()
	after := testutil.ToFloat64(jobsCompletedTotal)

	if after != before+1 {
		t.Errorf("Expected counter to increment by 1, got %f", after-before)
	}
}

func TestIncrementJobsFailed(t *testing.T) {
	before := testutil.ToFloat64(jobsFailedTotal.WithLabelValues("timeout"))
	IncrementJobsFailed("timeout")
	after := testutil.ToFloat64(jobsFailedTotal.WithLabelValues("timeout"))

	if after != before+1 {
		t.Errorf("Expected counter to increment by 1, got %f", after-before)
	}
}

func TestRecordJobProcessingDuration(t *testing.T) {
	// Record a job processing duration
	RecordJobProcessingDuration(45.5)

	// Verify the histogram recorded the observation
	count := testutil.CollectAndCount(jobProcessingDuration)
	if count == 0 {
		t.Error("Expected histogram to have metrics registered")
	}
}

func TestSetJobQueueDepth(t *testing.T) {
	SetJobQueueDepth("pending", 10.0)
	value := testutil.ToFloat64(jobQueueDepth.WithLabelValues("pending"))

	if value != 10.0 {
		t.Errorf("Expected gauge to be set to 10.0, got %f", value)
	}

	SetJobQueueDepth("pending", 5.0)
	value = testutil.ToFloat64(jobQueueDepth.WithLabelValues("pending"))

	if value != 5.0 {
		t.Errorf("Expected gauge to be set to 5.0, got %f", value)
	}
}

// Worker Metrics Tests

func TestSetWorkersActive(t *testing.T) {
	SetWorkersActive(5.0)
	value := testutil.ToFloat64(workersActive)

	if value != 5.0 {
		t.Errorf("Expected gauge to be set to 5.0, got %f", value)
	}

	SetWorkersActive(3.0)
	value = testutil.ToFloat64(workersActive)

	if value != 3.0 {
		t.Errorf("Expected gauge to be set to 3.0, got %f", value)
	}
}

func TestIncrementWorkerJobsProcessed(t *testing.T) {
	before := testutil.ToFloat64(workerJobsProcessedTotal.WithLabelValues("worker-1"))
	IncrementWorkerJobsProcessed("worker-1")
	after := testutil.ToFloat64(workerJobsProcessedTotal.WithLabelValues("worker-1"))

	if after != before+1 {
		t.Errorf("Expected counter to increment by 1, got %f", after-before)
	}
}

func TestIncrementWorkerErrors(t *testing.T) {
	before := testutil.ToFloat64(workerErrorsTotal.WithLabelValues("worker-2", "network_error"))
	IncrementWorkerErrors("worker-2", "network_error")
	after := testutil.ToFloat64(workerErrorsTotal.WithLabelValues("worker-2", "network_error"))

	if after != before+1 {
		t.Errorf("Expected counter to increment by 1, got %f", after-before)
	}
}

func TestRecordWorkerBatchDuration(t *testing.T) {
	// Record a worker batch duration
	RecordWorkerBatchDuration(30.5)

	// Verify the histogram recorded the observation
	count := testutil.CollectAndCount(workerBatchDuration)
	if count == 0 {
		t.Error("Expected histogram to have metrics registered")
	}
}

// Database Metrics Tests

func TestUpdateDBConnectionStats(t *testing.T) {
	UpdateDBConnectionStats(10, 5, 2)

	active := testutil.ToFloat64(dbConnectionsActive)
	idle := testutil.ToFloat64(dbConnectionsIdle)
	waiting := testutil.ToFloat64(dbConnectionsWaiting)

	if active != 10.0 {
		t.Errorf("Expected active connections to be 10.0, got %f", active)
	}
	if idle != 5.0 {
		t.Errorf("Expected idle connections to be 5.0, got %f", idle)
	}
	if waiting != 2.0 {
		t.Errorf("Expected waiting connections to be 2.0, got %f", waiting)
	}
}

// Podcast Metrics Tests

func TestIncrementPodcastsGenerated(t *testing.T) {
	before := testutil.ToFloat64(podcastsGeneratedTotal)
	IncrementPodcastsGenerated()
	after := testutil.ToFloat64(podcastsGeneratedTotal)

	if after != before+1 {
		t.Errorf("Expected counter to increment by 1, got %f", after-before)
	}
}

func TestRecordPodcastGenerationDuration(t *testing.T) {
	// Record a podcast generation duration
	RecordPodcastGenerationDuration(120.5)

	// Verify the histogram recorded the observation
	count := testutil.CollectAndCount(podcastGenerationDuration)
	if count == 0 {
		t.Error("Expected histogram to have metrics registered")
	}
}

func TestIncrementPodcastAudioRequests(t *testing.T) {
	before := testutil.ToFloat64(podcastAudioRequestsCurrent)
	IncrementPodcastAudioRequests()
	after := testutil.ToFloat64(podcastAudioRequestsCurrent)

	if after != before+1 {
		t.Errorf("Expected gauge to increment by 1, got %f", after-before)
	}
}

func TestDecrementPodcastAudioRequests(t *testing.T) {
	before := testutil.ToFloat64(podcastAudioRequestsCurrent)
	DecrementPodcastAudioRequests()
	after := testutil.ToFloat64(podcastAudioRequestsCurrent)

	if after != before-1 {
		t.Errorf("Expected gauge to decrement by 1, got %f", after-before)
	}
}

func TestIncrementPodcastGenerationFailures(t *testing.T) {
	before := testutil.ToFloat64(podcastGenerationFailuresTotal.WithLabelValues("audio_generation"))
	IncrementPodcastGenerationFailures("audio_generation")
	after := testutil.ToFloat64(podcastGenerationFailuresTotal.WithLabelValues("audio_generation"))

	if after != before+1 {
		t.Errorf("Expected counter to increment by 1, got %f", after-before)
	}
}

// AI Service Metrics Tests

func TestRecordAIAPICall(t *testing.T) {
	beforeCalls := testutil.ToFloat64(aiAPICallsTotal.WithLabelValues("generate_script", "openai"))

	RecordAIAPICall("generate_script", "openai", 2.5)

	afterCalls := testutil.ToFloat64(aiAPICallsTotal.WithLabelValues("generate_script", "openai"))

	if afterCalls != beforeCalls+1 {
		t.Errorf("Expected calls counter to increment by 1, got %f", afterCalls-beforeCalls)
	}

	// Verify latency histogram recorded the observation
	count := testutil.CollectAndCount(aiAPILatency)
	if count == 0 {
		t.Error("Expected latency histogram to have metrics registered")
	}
}

func TestIncrementAIAPIErrors(t *testing.T) {
	before := testutil.ToFloat64(aiAPIErrorsTotal.WithLabelValues("summarize", "anthropic", "rate_limit"))
	IncrementAIAPIErrors("summarize", "anthropic", "rate_limit")
	after := testutil.ToFloat64(aiAPIErrorsTotal.WithLabelValues("summarize", "anthropic", "rate_limit"))

	if after != before+1 {
		t.Errorf("Expected counter to increment by 1, got %f", after-before)
	}
}

// External Service Metrics Tests

func TestIncrementEmailsSent(t *testing.T) {
	before := testutil.ToFloat64(emailsSentTotal)
	IncrementEmailsSent()
	after := testutil.ToFloat64(emailsSentTotal)

	if after != before+1 {
		t.Errorf("Expected counter to increment by 1, got %f", after-before)
	}
}

func TestIncrementEmailFailures(t *testing.T) {
	before := testutil.ToFloat64(emailFailuresTotal.WithLabelValues("smtp_error"))
	IncrementEmailFailures("smtp_error")
	after := testutil.ToFloat64(emailFailuresTotal.WithLabelValues("smtp_error"))

	if after != before+1 {
		t.Errorf("Expected counter to increment by 1, got %f", after-before)
	}
}

func TestIncrementScrapingRequests(t *testing.T) {
	before := testutil.ToFloat64(scrapingRequestsTotal)
	IncrementScrapingRequests()
	after := testutil.ToFloat64(scrapingRequestsTotal)

	if after != before+1 {
		t.Errorf("Expected counter to increment by 1, got %f", after-before)
	}
}

func TestIncrementScrapingFailures(t *testing.T) {
	before := testutil.ToFloat64(scrapingFailuresTotal.WithLabelValues("timeout"))
	IncrementScrapingFailures("timeout")
	after := testutil.ToFloat64(scrapingFailuresTotal.WithLabelValues("timeout"))

	if after != before+1 {
		t.Errorf("Expected counter to increment by 1, got %f", after-before)
	}
}

// SSE Metrics Tests

func TestIncrementSSEConnections(t *testing.T) {
	before := testutil.ToFloat64(sseConnectionsActive)
	IncrementSSEConnections()
	after := testutil.ToFloat64(sseConnectionsActive)

	if after != before+1 {
		t.Errorf("Expected gauge to increment by 1, got %f", after-before)
	}
}

func TestDecrementSSEConnections(t *testing.T) {
	before := testutil.ToFloat64(sseConnectionsActive)
	DecrementSSEConnections()
	after := testutil.ToFloat64(sseConnectionsActive)

	if after != before-1 {
		t.Errorf("Expected gauge to decrement by 1, got %f", after-before)
	}
}

func TestIncrementSSEEventsSent(t *testing.T) {
	before := testutil.ToFloat64(sseEventsSentTotal.WithLabelValues("status_update"))
	IncrementSSEEventsSent("status_update")
	after := testutil.ToFloat64(sseEventsSentTotal.WithLabelValues("status_update"))

	if after != before+1 {
		t.Errorf("Expected counter to increment by 1, got %f", after-before)
	}
}

// Integration tests to verify multiple operations

func TestHTTPMetricsIntegration(t *testing.T) {
	// Simulate a complete HTTP request cycle
	IncrementActiveRequests()
	RecordHTTPRequest("POST", "/api/items", "201")
	RecordHTTPDuration("POST", "/api/items", 0.15)
	RecordHTTPRequestSize("POST", "/api/items", 512.0)
	RecordHTTPResponseSize("POST", "/api/items", 256.0)
	DecrementActiveRequests()

	// Verify all metrics were recorded
	requests := testutil.ToFloat64(httpRequestsTotal.WithLabelValues("POST", "/api/items", "201"))
	if requests == 0 {
		t.Error("Expected HTTP request to be recorded")
	}
}

func TestJobQueueMetricsIntegration(t *testing.T) {
	// Simulate a complete job lifecycle
	IncrementJobsEnqueued()
	SetJobQueueDepth("pending", 5.0)
	IncrementJobsProcessing()
	RecordJobProcessingDuration(25.0)
	IncrementJobsCompleted()
	DecrementJobsProcessing()
	SetJobQueueDepth("pending", 4.0)

	// Verify metrics were recorded
	completed := testutil.ToFloat64(jobsCompletedTotal)
	if completed == 0 {
		t.Error("Expected job completion to be recorded")
	}
}

func TestWorkerMetricsIntegration(t *testing.T) {
	// Simulate worker activity
	SetWorkersActive(3.0)
	IncrementWorkerJobsProcessed("worker-test")
	RecordWorkerBatchDuration(15.5)

	// Verify metrics were recorded
	active := testutil.ToFloat64(workersActive)
	if active != 3.0 {
		t.Errorf("Expected 3 active workers, got %f", active)
	}
}

func TestPodcastGenerationMetricsIntegration(t *testing.T) {
	// Simulate podcast generation
	IncrementPodcastAudioRequests()
	RecordPodcastGenerationDuration(180.0)
	IncrementPodcastsGenerated()
	DecrementPodcastAudioRequests()

	// Verify metrics were recorded
	generated := testutil.ToFloat64(podcastsGeneratedTotal)
	if generated == 0 {
		t.Error("Expected podcast generation to be recorded")
	}
}

func TestAIServiceMetricsIntegration(t *testing.T) {
	// Simulate AI API calls
	RecordAIAPICall("generate_script", "openai", 3.5)
	RecordAIAPICall("generate_script", "openai", 2.8)

	// Verify metrics were recorded
	calls := testutil.ToFloat64(aiAPICallsTotal.WithLabelValues("generate_script", "openai"))
	if calls == 0 {
		t.Error("Expected AI API calls to be recorded")
	}
}
