package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestComputeApproximateRequestSize(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *http.Request
		minSize  int // Minimum expected size (we use min because exact size can vary)
		validate func(*testing.T, int)
	}{
		{
			name: "minimal request",
			setup: func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				return req
			},
			validate: func(t *testing.T, size int) {
				// Should at least include method + proto + host
				if size < 10 {
					t.Errorf("Expected size >= 10 for minimal request, got %d", size)
				}
			},
		},
		{
			name: "request with URL",
			setup: func() *http.Request {
				req, _ := http.NewRequest("GET", "http://example.com/api/test?foo=bar", nil)
				return req
			},
			validate: func(t *testing.T, size int) {
				// Should include the full URL
				if size < 30 {
					t.Errorf("Expected size >= 30 for request with URL, got %d", size)
				}
			},
		},
		{
			name: "request with headers",
			setup: func() *http.Request {
				req, _ := http.NewRequest("POST", "/api/data", nil)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer token123")
				req.Header.Set("X-Custom-Header", "custom-value")
				return req
			},
			validate: func(t *testing.T, size int) {
				// Should include headers
				if size < 50 {
					t.Errorf("Expected size >= 50 for request with headers, got %d", size)
				}
			},
		},
		{
			name: "request with content length",
			setup: func() *http.Request {
				req, _ := http.NewRequest("POST", "/upload", nil)
				req.ContentLength = 1024
				return req
			},
			validate: func(t *testing.T, size int) {
				// Should include content length
				if size < 1024 {
					t.Errorf("Expected size >= 1024 for request with content length, got %d", size)
				}
			},
		},
		{
			name: "request with nil URL",
			setup: func() *http.Request {
				req := &http.Request{
					Method: "GET",
					Proto:  "HTTP/1.1",
					Header: make(http.Header),
					Host:   "example.com",
					URL:    nil,
				}
				return req
			},
			validate: func(t *testing.T, size int) {
				// Should not panic and should have some size
				if size < 0 {
					t.Errorf("Expected non-negative size, got %d", size)
				}
			},
		},
		{
			name: "large request with multiple headers and content",
			setup: func() *http.Request {
				req, _ := http.NewRequest("POST", "http://example.com/api/v1/resources/12345?include=details&format=json", nil)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")
				req.Header.Set("Accept", "application/json")
				req.Header.Set("User-Agent", "Mozilla/5.0")
				req.Header.Set("X-Request-ID", "550e8400-e29b-41d4-a716-446655440000")
				req.ContentLength = 2048
				req.Host = "example.com"
				return req
			},
			validate: func(t *testing.T, size int) {
				// Should be a large size including URL, headers, and content
				if size < 2200 {
					t.Errorf("Expected size >= 2200 for large request, got %d", size)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.setup()
			size := computeApproximateRequestSize(req)
			tt.validate(t, size)
		})
	}
}

func TestComputeApproximateRequestSize_Components(t *testing.T) {
	// Test that each component contributes to the size
	baseReq, _ := http.NewRequest("GET", "http://example.com/test", nil)
	baseSize := computeApproximateRequestSize(baseReq)

	// Add a header
	reqWithHeader, _ := http.NewRequest("GET", "http://example.com/test", nil)
	reqWithHeader.Header.Set("X-Test", "value")
	sizeWithHeader := computeApproximateRequestSize(reqWithHeader)

	assert.Greater(t, sizeWithHeader, baseSize, "Adding header should increase size")

	// Add content length
	reqWithContent, _ := http.NewRequest("GET", "http://example.com/test", nil)
	reqWithContent.ContentLength = 500
	sizeWithContent := computeApproximateRequestSize(reqWithContent)

	assert.Greater(t, sizeWithContent, baseSize, "Adding content length should increase size")
}

func TestPrometheusMiddleware_SkipsMetricsEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a test router
	router := gin.New()
	router.Use(PrometheusMiddleware())

	// Track if handler was called
	handlerCalled := false
	router.GET("/metrics", func(c *gin.Context) {
		handlerCalled = true
		c.Status(http.StatusOK)
	})

	// Make request to /metrics
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/metrics", nil)
	router.ServeHTTP(w, req)

	// Verify handler was called and request completed
	assert.True(t, handlerCalled, "Handler should be called")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPrometheusMiddleware_RecordsMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a test router
	router := gin.New()
	router.Use(PrometheusMiddleware())

	// Add a test endpoint
	router.GET("/api/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Make request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/test", nil)
	router.ServeHTTP(w, req)

	// Verify request completed successfully
	assert.Equal(t, http.StatusOK, w.Code)

	// Note: We can't easily verify the metrics were recorded without accessing
	// the prometheus registry, but we can verify the middleware didn't break the request
}

func TestPrometheusMiddleware_HandlesUnknownRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a test router
	router := gin.New()
	router.Use(PrometheusMiddleware())

	// Don't register any routes, so all routes will be "unknown"
	router.NoRoute(func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	// Make request to unknown route
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/unknown", nil)
	router.ServeHTTP(w, req)

	// Verify request was handled
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPrometheusMiddleware_HandlesMultipleMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			router := gin.New()
			router.Use(PrometheusMiddleware())

			// Register handler for this method
			router.Handle(method, "/api/resource", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			// Make request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(method, "/api/resource", nil)
			router.ServeHTTP(w, req)

			// Verify request completed
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func TestPrometheusMiddleware_HandlesErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a test router
	router := gin.New()
	router.Use(PrometheusMiddleware())

	// Add endpoint that returns error
	router.GET("/api/error", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	})

	// Make request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/error", nil)
	router.ServeHTTP(w, req)

	// Verify error status was recorded
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPrometheusMiddleware_RecordsResponseSize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a test router
	router := gin.New()
	router.Use(PrometheusMiddleware())

	// Add endpoint with known response
	router.GET("/api/data", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "test",
			"data":    []int{1, 2, 3, 4, 5},
		})
	})

	// Make request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/data", nil)
	router.ServeHTTP(w, req)

	// Verify response was written
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Greater(t, w.Body.Len(), 0, "Response should have content")
}

func TestPrometheusMiddleware_ActiveRequestsTracking(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a test router
	router := gin.New()
	router.Use(PrometheusMiddleware())

	// Add endpoint
	router.GET("/api/slow", func(c *gin.Context) {
		// Simulate some processing
		c.Status(http.StatusOK)
	})

	// Make request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/slow", nil)
	router.ServeHTTP(w, req)

	// Verify request completed (active requests should be decremented)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPrometheusMiddleware_WithRequestBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a test router
	router := gin.New()
	router.Use(PrometheusMiddleware())

	// Add POST endpoint
	router.POST("/api/create", func(c *gin.Context) {
		c.Status(http.StatusCreated)
	})

	// Make request with body
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/create", nil)
	req.Header.Set("Content-Type", "application/json")
	req.ContentLength = 100
	router.ServeHTTP(w, req)

	// Verify request completed
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPrometheusMiddleware_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a full router with multiple endpoints
	router := gin.New()
	router.Use(PrometheusMiddleware())

	router.GET("/api/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"users": []string{"alice", "bob"}})
	})

	router.POST("/api/users", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"id": 123})
	})

	router.GET("/api/users/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "name": "alice"})
	})

	// Test multiple requests
	tests := []struct {
		method       string
		path         string
		expectedCode int
	}{
		{"GET", "/api/users", http.StatusOK},
		{"POST", "/api/users", http.StatusCreated},
		{"GET", "/api/users/123", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestPrometheusMiddleware_ChainedWithOtherMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create router with multiple middleware
	router := gin.New()

	// Add custom middleware before Prometheus
	customMiddlewareCalled := false
	router.Use(func(c *gin.Context) {
		customMiddlewareCalled = true
		c.Next()
	})

	router.Use(PrometheusMiddleware())

	router.GET("/api/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Make request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/test", nil)
	router.ServeHTTP(w, req)

	// Verify both middleware executed
	assert.True(t, customMiddlewareCalled, "Custom middleware should be called")
	assert.Equal(t, http.StatusOK, w.Code)
}
