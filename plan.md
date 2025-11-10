# Prometheus and Grafana Monitoring Setup Plan

## Overview
This plan outlines the implementation of comprehensive monitoring for BriefBot using Prometheus for metrics collection and Grafana for visualization. The setup will track application performance, worker job processing, database health, API endpoints, and external service integrations.

## Architecture

### Components
1. **Prometheus** - Metrics collection and storage
2. **Grafana** - Metrics visualization and dashboards
3. **Prometheus Go Client** - Application instrumentation
4. **PostgreSQL Exporter** - Database metrics
5. **Node Exporter** (optional) - System-level metrics

### Metrics Strategy

#### 1. Application Metrics (Go Backend)
- **HTTP Metrics**
  - Request count by endpoint, method, status code
  - Request duration histogram
  - Request size histogram
  - Response size histogram
  - Active requests gauge

- **Job Queue Metrics**
  - Items enqueued counter
  - Items processing gauge
  - Items completed counter
  - Items failed counter
  - Job processing duration histogram
  - Queue depth by status (pending, processing, completed, failed)
  - Retry count histogram

- **Worker Metrics**
  - Active workers gauge
  - Worker utilization percentage
  - Jobs processed per worker counter
  - Worker errors counter
  - Batch processing duration histogram

- **Podcast Generation Metrics**
  - Podcasts generated counter
  - Podcast generation duration histogram
  - Audio generation requests (concurrent gauge)
  - Audio generation failures counter
  - Speech API latency histogram

- **AI Service Metrics**
  - AI API calls counter (by operation: summarize, extract, generate)
  - AI API latency histogram
  - AI API errors counter
  - Token usage counter (if available)

- **Email Service Metrics**
  - Emails sent counter
  - Email failures counter
  - Digest generation duration histogram

- **Database Metrics**
  - Connection pool stats (active, idle, waiting)
  - Query duration histogram
  - Database errors counter

- **External Service Metrics**
  - Scraping service requests counter
  - Scraping failures counter
  - R2 upload/download operations counter
  - R2 operation duration histogram
  - FAL API requests counter
  - Groq API requests counter

- **SSE Metrics**
  - Active SSE connections gauge
  - SSE events sent counter
  - SSE connection duration histogram

#### 2. System Metrics
- CPU usage
- Memory usage
- Disk I/O
- Network I/O

#### 3. Database Metrics (PostgreSQL)
- Connection count
- Transaction rate
- Query performance
- Cache hit ratio
- Table sizes
- Index usage

## Implementation Plan

### Phase 1: Infrastructure Setup

#### 1.1 Update docker-compose.yml
**File**: `/docker-compose.yml`

Add new services:
```yaml
  # Prometheus for metrics collection
  prometheus:
    image: prom/prometheus:latest
    container_name: briefbot-prometheus
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/usr/share/prometheus/console_libraries'
      - '--web.console.templates=/usr/share/prometheus/consoles'
      - '--storage.tsdb.retention.time=30d'
    ports:
      - "9090:9090"
    networks:
      - briefbot-network
    restart: unless-stopped

  # Grafana for visualization
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
      - GF_SERVER_ROOT_URL=http://localhost:3001
    ports:
      - "3001:3000"
    networks:
      - briefbot-network
    depends_on:
      - prometheus
    restart: unless-stopped

  # PostgreSQL Exporter for database metrics
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

#### 1.2 Create Prometheus Configuration
**File**: `/monitoring/prometheus.yml`

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s
  external_labels:
    cluster: 'briefbot'
    environment: 'development'

scrape_configs:
  # Backend Go application
  - job_name: 'briefbot-backend'
    static_configs:
      - targets: ['backend:8080']
    metrics_path: '/metrics'

  # PostgreSQL database
  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres-exporter:9187']

  # Prometheus itself
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
```

#### 1.3 Create Grafana Provisioning
**File**: `/monitoring/grafana/provisioning/datasources/prometheus.yml`

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

**File**: `/monitoring/grafana/provisioning/dashboards/dashboard.yml`

```yaml
apiVersion: 1

providers:
  - name: 'BriefBot Dashboards'
    orgId: 1
    folder: ''
    type: file
    disableDeletion: false
    updateIntervalSeconds: 10
    allowUiUpdates: true
    options:
      path: /var/lib/grafana/dashboards
```

### Phase 2: Backend Instrumentation

#### 2.1 Add Prometheus Dependencies
**File**: `/backend/go.mod`

Add dependency:
```
github.com/prometheus/client_golang v1.19.0
```

Run:
```bash
cd backend && go get github.com/prometheus/client_golang/prometheus
cd backend && go get github.com/prometheus/client_golang/prometheus/promauto
cd backend && go get github.com/prometheus/client_golang/prometheus/promhttp
```

#### 2.2 Create Metrics Package
**File**: `/backend/internal/metrics/metrics.go`

Create a centralized metrics package that defines all application metrics:
- HTTP metrics (request count, duration, size)
- Job queue metrics (enqueued, processing, completed, failed)
- Worker metrics (active workers, jobs processed)
- Podcast metrics (generated, duration)
- AI service metrics (calls, latency, errors)
- Email metrics (sent, failed)
- Database metrics (connections, query duration)
- External service metrics (R2, FAL, Groq)
- SSE metrics (connections, events)

#### 2.3 Create Middleware for HTTP Metrics
**File**: `/backend/internal/middleware/prometheus.go`

Create Gin middleware to automatically track:
- HTTP request count by endpoint, method, status
- Request duration
- Request/response sizes
- Active requests

#### 2.4 Instrument Services

**Files to modify**:
1. `/backend/internal/services/jobqueue.go`
   - Add metrics for enqueue, dequeue, complete, fail operations
   - Track queue depth by status
   - Track processing duration

2. `/backend/internal/services/worker.go`
   - Track active workers
   - Track jobs processed per worker
   - Track worker errors
   - Track batch processing duration

3. `/backend/internal/services/podcast.go`
   - Track podcast generation count
   - Track generation duration
   - Track concurrent audio requests
   - Track failures

4. `/backend/internal/services/ai.go`
   - Track API calls by operation type
   - Track latency
   - Track errors

5. `/backend/internal/services/email.go`
   - Track emails sent
   - Track failures
   - Track digest generation duration

6. `/backend/internal/services/speech.go`
   - Track FAL API requests
   - Track latency
   - Track failures

7. `/backend/internal/services/r2.go`
   - Track upload/download operations
   - Track operation duration
   - Track failures

8. `/backend/internal/services/scraping.go`
   - Track scraping requests
   - Track failures

9. `/backend/internal/services/sse.go`
   - Track active connections
   - Track events sent
   - Track connection duration

#### 2.5 Add Database Pool Metrics
**File**: `/backend/cmd/server/main.go`

Add metrics collector for pgxpool stats:
- Active connections
- Idle connections
- Waiting connections
- Total connections

#### 2.6 Add Metrics Endpoint
**File**: `/backend/internal/handlers/routes.go`

Add Prometheus metrics endpoint:
```go
router.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

### Phase 3: Dashboard Creation

#### 3.1 Create Main Application Dashboard
**File**: `/monitoring/grafana/dashboards/briefbot-overview.json`

Panels:
- HTTP request rate (requests/sec by endpoint)
- HTTP error rate (4xx, 5xx)
- HTTP request duration (p50, p95, p99)
- Active HTTP requests
- Job queue depth (by status)
- Job processing rate
- Job failure rate
- Active workers
- Worker utilization
- Database connection pool usage
- Database query duration
- Memory usage
- CPU usage

#### 3.2 Create Job Processing Dashboard
**File**: `/monitoring/grafana/dashboards/briefbot-jobs.json`

Panels:
- Items enqueued over time
- Items by status (stacked area chart)
- Job processing duration histogram
- Job failure rate by type
- Retry count distribution
- Worker performance (jobs/worker)
- Queue backlog trend
- Processing throughput

#### 3.3 Create External Services Dashboard
**File**: `/monitoring/grafana/dashboards/briefbot-external.json`

Panels:
- AI API calls (by operation)
- AI API latency (p50, p95, p99)
- AI API error rate
- Speech API requests
- Speech API latency
- R2 operations (upload/download)
- R2 operation duration
- Scraping requests
- Scraping failures
- Email delivery rate

#### 3.4 Create Database Dashboard
**File**: `/monitoring/grafana/dashboards/briefbot-database.json`

Panels:
- Connection pool usage
- Active vs idle connections
- Query rate
- Query duration
- Transaction rate
- Cache hit ratio
- Database size
- Table sizes
- Slow queries

#### 3.5 Create Podcast Generation Dashboard
**File**: `/monitoring/grafana/dashboards/briefbot-podcasts.json`

Panels:
- Podcasts generated over time
- Podcast generation duration
- Concurrent audio requests
- Audio generation failures
- Speech API performance
- R2 upload performance for podcasts
- Average podcast size

### Phase 4: Alerting (Optional but Recommended)

#### 4.1 Create Alert Rules
**File**: `/monitoring/prometheus/alerts.yml`

Define alert rules for:
- High error rate (>5% 5xx errors)
- High request latency (p95 > 1s)
- Job queue backlog (>100 pending items)
- High job failure rate (>10%)
- Database connection pool exhaustion (>90% used)
- Worker failures
- External service failures
- Low disk space

#### 4.2 Configure Alertmanager (Optional)
**File**: `/monitoring/alertmanager.yml`

Configure notification channels:
- Email
- Slack
- PagerDuty
- Webhook

### Phase 5: Documentation

#### 5.1 Update README.md
Add section about monitoring:
- How to access Prometheus (http://localhost:9090)
- How to access Grafana (http://localhost:3001)
- Default credentials
- Available dashboards
- Key metrics to watch

#### 5.2 Create Monitoring Guide
**File**: `/monitoring/README.md`

Document:
- Architecture overview
- Metrics catalog
- Dashboard guide
- Alert guide
- Troubleshooting
- Best practices

#### 5.3 Update .env.example
Add monitoring-related environment variables if needed:
```
# Monitoring
METRICS_ENABLED=true
METRICS_PATH=/metrics
```

### Phase 6: Testing

#### 6.1 Verify Metrics Collection
- Start all services with docker-compose
- Generate test traffic
- Verify metrics appear in Prometheus
- Verify dashboards display data in Grafana

#### 6.2 Load Testing
- Use tool like `hey` or `wrk` to generate load
- Verify metrics accurately reflect load
- Check for any performance impact

#### 6.3 Alert Testing
- Trigger alert conditions
- Verify alerts fire correctly
- Verify notifications work

## File Structure

```
briefbot/
├── docker-compose.yml (updated)
├── monitoring/
│   ├── README.md
│   ├── prometheus.yml
│   ├── alerts.yml (optional)
│   ├── alertmanager.yml (optional)
│   └── grafana/
│       ├── provisioning/
│       │   ├── datasources/
│       │   │   └── prometheus.yml
│       │   └── dashboards/
│       │       └── dashboard.yml
│       └── dashboards/
│           ├── briefbot-overview.json
│           ├── briefbot-jobs.json
│           ├── briefbot-external.json
│           ├── briefbot-database.json
│           └── briefbot-podcasts.json
├── backend/
│   ├── go.mod (updated)
│   ├── cmd/server/main.go (updated)
│   └── internal/
│       ├── metrics/
│       │   └── metrics.go (new)
│       ├── middleware/
│       │   └── prometheus.go (new)
│       ├── handlers/
│       │   └── routes.go (updated)
│       └── services/
│           ├── jobqueue.go (updated)
│           ├── worker.go (updated)
│           ├── podcast.go (updated)
│           ├── ai.go (updated)
│           ├── email.go (updated)
│           ├── speech.go (updated)
│           ├── r2.go (updated)
│           ├── scraping.go (updated)
│           └── sse.go (updated)
└── README.md (updated)
```

## Implementation Order

1. **Infrastructure** (Day 1)
   - Add Prometheus, Grafana, postgres-exporter to docker-compose
   - Create Prometheus config
   - Create Grafana provisioning configs
   - Test basic setup

2. **Basic Metrics** (Day 2)
   - Add Prometheus Go client dependency
   - Create metrics package
   - Add HTTP middleware
   - Add /metrics endpoint
   - Verify basic metrics collection

3. **Service Instrumentation** (Day 3-4)
   - Instrument job queue service
   - Instrument worker service
   - Instrument podcast service
   - Instrument AI service
   - Instrument email service
   - Instrument external services (R2, speech, scraping)
   - Instrument SSE service
   - Add database pool metrics

4. **Dashboards** (Day 5)
   - Create overview dashboard
   - Create job processing dashboard
   - Create external services dashboard
   - Create database dashboard
   - Create podcast dashboard

5. **Documentation & Testing** (Day 6)
   - Update README
   - Create monitoring guide
   - Test all metrics
   - Load testing
   - Final verification

6. **Alerting** (Optional - Day 7)
   - Define alert rules
   - Configure Alertmanager
   - Test alerts

## Key Metrics to Track

### Critical Metrics (Red Flags)
1. HTTP 5xx error rate > 1%
2. Job failure rate > 5%
3. Database connection pool > 90% utilized
4. Request p95 latency > 2s
5. Queue backlog > 500 items
6. Worker errors > 10/hour

### Performance Metrics
1. Request throughput (req/sec)
2. Request latency (p50, p95, p99)
3. Job processing throughput
4. Database query duration
5. External API latency

### Resource Metrics
1. Memory usage
2. CPU usage
3. Database connections
4. Active workers
5. SSE connections

## Benefits

1. **Visibility**: Real-time insight into application performance
2. **Debugging**: Quickly identify bottlenecks and issues
3. **Capacity Planning**: Understand resource usage patterns
4. **SLA Monitoring**: Track uptime and performance SLAs
5. **Alerting**: Proactive notification of issues
6. **Historical Analysis**: Trend analysis and capacity planning
7. **External Service Monitoring**: Track third-party API performance

## Considerations

1. **Performance Impact**: Minimal (<1% overhead with proper implementation)
2. **Storage**: Prometheus data retention (30 days default, configurable)
3. **Security**: Metrics endpoint should be protected in production
4. **Cardinality**: Avoid high-cardinality labels (e.g., user IDs, timestamps)
5. **Naming**: Follow Prometheus naming conventions (snake_case, _total suffix for counters)

## Next Steps After Implementation

1. Set up production monitoring infrastructure
2. Configure production alerts
3. Set up notification channels (Slack, PagerDuty)
4. Create runbooks for common alerts
5. Train team on dashboard usage
6. Establish SLOs and SLIs
7. Regular dashboard reviews and refinements

## Resources

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Prometheus Go Client](https://github.com/prometheus/client_golang)
- [Prometheus Best Practices](https://prometheus.io/docs/practices/naming/)
- [Grafana Dashboard Best Practices](https://grafana.com/docs/grafana/latest/dashboards/build-dashboards/best-practices/)
