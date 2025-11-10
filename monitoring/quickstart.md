# Monitoring Quick Start Guide

## 1. Add Infrastructure (5 minutes)

### Update docker-compose.yml
Add three new services at the end:

```yaml
  prometheus:
    image: prom/prometheus:latest
    container_name: briefbot-prometheus
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--storage.tsdb.retention.time=30d'
    ports:
      - "9090:9090"
    networks:
      - briefbot-network
    restart: unless-stopped

  grafana:
    image: grafana/grafana:latest
    container_name: briefbot-grafana
    volumes:
      - grafana_data:/var/lib/grafana
      - ./monitoring/grafana/provisioning:/etc/grafana/provisioning
      - ./monitoring/grafana/dashboards:/var/lib/grafana/dashboards
    environment:
      - GF_SECURITY_ADMIN_USER=admin
      - GF_SECURITY_ADMIN_PASSWORD=admin
      - GF_USERS_ALLOW_SIGN_UP=false
    ports:
      - "3001:3000"
    networks:
      - briefbot-network
    depends_on:
      - prometheus
    restart: unless-stopped

  postgres-exporter:
    image: prometheuscommunity/postgres-exporter:latest
    container_name: briefbot-postgres-exporter
    environment:
      DATA_SOURCE_NAME: "postgresql://briefbot:briefbot@postgres:5432/briefbot?sslmode=disable"
    ports:
      - "9187:9187"
    networks:
      - briefbot-network
    depends_on:
      - postgres
    restart: unless-stopped
```

Add volumes:
```yaml
volumes:
  postgres_data:
  prometheus_data:
  grafana_data:
```

### Create monitoring directory structure
```bash
mkdir -p monitoring/grafana/provisioning/datasources
mkdir -p monitoring/grafana/provisioning/dashboards
mkdir -p monitoring/grafana/dashboards
```

### Create Prometheus config
**File**: `monitoring/prometheus.yml`
```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'briefbot-backend'
    static_configs:
      - targets: ['backend:8080']
    metrics_path: '/metrics'

  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres-exporter:9187']

  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
```

### Create Grafana datasource config
**File**: `monitoring/grafana/provisioning/datasources/prometheus.yml`
```yaml
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
    editable: true
```

### Create Grafana dashboard config
**File**: `monitoring/grafana/provisioning/dashboards/dashboard.yml`
```yaml
apiVersion: 1

providers:
  - name: 'BriefBot'
    orgId: 1
    folder: ''
    type: file
    disableDeletion: false
    updateIntervalSeconds: 10
    allowUiUpdates: true
    options:
      path: /var/lib/grafana/dashboards
```

## 2. Add Go Dependencies (2 minutes)

```bash
cd backend
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promauto
go get github.com/prometheus/client_golang/prometheus/promhttp
```

## 3. Create Metrics Package (10 minutes)

**File**: `backend/internal/metrics/metrics.go`

```go
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
            Help: "Total HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )

    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "briefbot_http_request_duration_seconds",
            Help:    "HTTP request duration",
            Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
        },
        []string{"method", "endpoint"},
    )

    httpActiveRequests = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "briefbot_http_active_requests",
            Help: "Active HTTP requests",
        },
    )

    // Job Metrics
    jobsEnqueuedTotal = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "briefbot_jobs_enqueued_total",
            Help: "Total jobs enqueued",
        },
    )

    jobsProcessingCurrent = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "briefbot_jobs_processing_current",
            Help: "Jobs currently processing",
        },
    )

    jobsCompletedTotal = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "briefbot_jobs_completed_total",
            Help: "Total jobs completed",
        },
    )

    jobsFailedTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "briefbot_jobs_failed_total",
            Help: "Total jobs failed",
        },
        []string{"error_type"},
    )

    jobProcessingDuration = promauto.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "briefbot_job_processing_duration_seconds",
            Help:    "Job processing duration",
            Buckets: []float64{1, 5, 10, 30, 60, 120, 300, 600},
        },
    )

    // Database Metrics
    dbConnectionsActive = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "briefbot_db_connections_active",
            Help: "Active database connections",
        },
    )

    dbConnectionsIdle = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "briefbot_db_connections_idle",
            Help: "Idle database connections",
        },
    )
)

// Helper functions
func RecordHTTPRequest(method, endpoint, status string) {
    httpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
}

func RecordHTTPDuration(method, endpoint string, duration float64) {
    httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
}

func IncrementActiveRequests() {
    httpActiveRequests.Inc()
}

func DecrementActiveRequests() {
    httpActiveRequests.Dec()
}

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

func UpdateDBConnectionStats(active, idle int32) {
    dbConnectionsActive.Set(float64(active))
    dbConnectionsIdle.Set(float64(idle))
}
```

## 4. Add Prometheus Middleware (5 minutes)

**File**: `backend/internal/middleware/prometheus.go`

```go
package middleware

import (
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/yamirghofran/briefbot/internal/metrics"
)

func PrometheusMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Request.URL.Path == "/metrics" {
            c.Next()
            return
        }

        start := time.Now()
        metrics.IncrementActiveRequests()
        defer metrics.DecrementActiveRequests()

        c.Next()

        duration := time.Since(start).Seconds()
        status := strconv.Itoa(c.Writer.Status())
        endpoint := c.FullPath()
        if endpoint == "" {
            endpoint = "unknown"
        }

        metrics.RecordHTTPRequest(c.Request.Method, endpoint, status)
        metrics.RecordHTTPDuration(c.Request.Method, endpoint, duration)
    }
}
```

## 5. Update Main Server (5 minutes)

**File**: `backend/cmd/server/main.go`

Add imports:
```go
import (
    // ... existing imports ...
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/yamirghofran/briefbot/internal/metrics"
    "github.com/yamirghofran/briefbot/internal/middleware"
)
```

Add middleware and metrics endpoint:
```go
func main() {
    // ... existing code ...

    // Initialize Gin router
    router := gin.Default()

    // Add Prometheus middleware
    router.Use(middleware.PrometheusMiddleware())

    // Add CORS middleware
    router.Use(func(c *gin.Context) {
        // ... existing CORS code ...
    })

    // Setup routes
    handlers.SetupRoutes(router, userService, itemService, digestService, podcastService, sseManager)

    // Metrics endpoint
    router.GET("/metrics", gin.WrapH(promhttp.Handler()))

    // Swagger documentation endpoint
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // ... rest of existing code ...

    // Start database metrics collector
    go func() {
        ticker := time.NewTicker(15 * time.Second)
        defer ticker.Stop()
        
        for range ticker.C {
            stats := pool.Stat()
            metrics.UpdateDBConnectionStats(
                stats.AcquiredConns(),
                stats.IdleConns(),
            )
        }
    }()

    // ... rest of existing code ...
}
```

## 6. Instrument Job Queue (10 minutes)

**File**: `backend/internal/services/jobqueue.go`

Add import:
```go
import (
    // ... existing imports ...
    "time"
    "github.com/yamirghofran/briefbot/internal/metrics"
)
```

Update methods:
```go
func (s *jobQueueService) EnqueueItem(ctx context.Context, userID int32, title string, url string) (*db.Item, error) {
    // ... existing code ...
    
    item, err := s.querier.CreatePendingItem(ctx, params)
    if err != nil {
        return nil, fmt.Errorf("failed to enqueue item: %w", err)
    }

    metrics.IncrementJobsEnqueued()  // ADD THIS

    // ... rest of code ...
}

func (s *jobQueueService) MarkItemAsProcessing(ctx context.Context, itemID int32) error {
    // ... existing code ...
    
    err = s.querier.UpdateItemAsProcessing(ctx, itemID)
    if err != nil {
        return fmt.Errorf("failed to mark item as processing: %w", err)
    }

    metrics.IncrementJobsProcessing()  // ADD THIS

    // ... rest of code ...
}

func (s *jobQueueService) CompleteItem(ctx context.Context, itemID int32, title, textContent, summary, itemType, platform string, tags, authors []string) error {
    start := time.Now()  // ADD THIS
    defer func() {
        duration := time.Since(start).Seconds()
        metrics.RecordJobProcessingDuration(duration)
    }()

    // ... existing code ...
    
    if err != nil {
        return fmt.Errorf("failed to complete item: %w", err)
    }

    metrics.DecrementJobsProcessing()  // ADD THIS
    metrics.IncrementJobsCompleted()   // ADD THIS

    // ... rest of code ...
}

func (s *jobQueueService) FailItem(ctx context.Context, itemID int32, errorMsg string) error {
    // ... existing code ...
    
    metrics.DecrementJobsProcessing()  // ADD THIS
    metrics.IncrementJobsFailed("unknown")  // ADD THIS (can categorize errors later)

    // ... rest of code ...
}
```

## 7. Test the Setup (5 minutes)

```bash
# Start everything
docker-compose up --build

# Access services
# Prometheus: http://localhost:9090
# Grafana: http://localhost:3001 (admin/admin)
# Metrics endpoint: http://localhost:8080/metrics

# Generate some traffic
curl http://localhost:8080/users
curl -X POST http://localhost:8080/items -H "Content-Type: application/json" -d '{"user_id":1,"url":"https://example.com"}'

# Check metrics
curl http://localhost:8080/metrics | grep briefbot

# Check Prometheus
# Go to http://localhost:9090/targets - should see backend as UP
# Go to http://localhost:9090/graph - query: briefbot_http_requests_total

# Check Grafana
# Go to http://localhost:3001
# Login: admin/admin
# Go to Explore, select Prometheus datasource
# Query: briefbot_http_requests_total
```

## 8. Create Basic Dashboard (10 minutes)

In Grafana (http://localhost:3001):

1. Click "+" → "Dashboard"
2. Click "Add visualization"
3. Select "Prometheus" datasource
4. Add these panels:

**Panel 1: Request Rate**
- Query: `rate(briefbot_http_requests_total[5m])`
- Legend: `{{method}} {{endpoint}}`
- Type: Time series

**Panel 2: Active Requests**
- Query: `briefbot_http_active_requests`
- Type: Stat

**Panel 3: Jobs Processing**
- Query: `briefbot_jobs_processing_current`
- Type: Gauge

**Panel 4: Job Completion Rate**
- Query: `rate(briefbot_jobs_completed_total[5m])`
- Type: Time series

5. Save dashboard as "BriefBot Overview"

## Next Steps

1. Add more metrics to other services (see MONITORING_SPEC.md)
2. Create additional dashboards for specific areas
3. Set up alerts for critical metrics
4. Add more detailed instrumentation

## Useful Prometheus Queries

```promql
# Total request rate
sum(rate(briefbot_http_requests_total[5m]))

# Error rate
sum(rate(briefbot_http_requests_total{status=~"5.."}[5m]))

# P95 latency
histogram_quantile(0.95, rate(briefbot_http_request_duration_seconds_bucket[5m]))

# Jobs in queue
briefbot_jobs_processing_current

# Job completion rate
rate(briefbot_jobs_completed_total[5m])

# Job failure rate
rate(briefbot_jobs_failed_total[5m])

# Database connections
briefbot_db_connections_active
briefbot_db_connections_idle
```

## Troubleshooting

**Prometheus not scraping backend:**
- Check `docker-compose logs prometheus`
- Verify backend is accessible: `docker exec briefbot-prometheus wget -O- http://backend:8080/metrics`
- Check Prometheus targets: http://localhost:9090/targets

**No metrics showing:**
- Verify metrics endpoint: `curl http://localhost:8080/metrics`
- Check for Go errors: `docker-compose logs backend`
- Ensure middleware is registered before routes

**Grafana can't connect to Prometheus:**
- Check datasource config in Grafana
- Verify Prometheus URL: `http://prometheus:9090`
- Test connection in Grafana datasource settings

## Access URLs

- **Application**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Metrics**: http://localhost:8080/metrics
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3001 (admin/admin)
- **Swagger**: http://localhost:8080/swagger/index.html

