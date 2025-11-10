package middleware

import (
	"net/http"
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
