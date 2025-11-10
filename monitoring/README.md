# BriefBot Monitoring Implementation Plan

## 📋 Overview

This directory contains a comprehensive plan for implementing Prometheus and Grafana monitoring for the BriefBot application. The monitoring setup will provide real-time visibility into application performance, job processing, external service integrations, and system health.

## 📚 Documentation Structure

### 1. **plan.md** - High-Level Implementation Plan
- Architecture overview
- Metrics strategy
- Implementation phases
- File structure
- Benefits and considerations

**Start here** to understand the overall approach and architecture.

### 2. **MONITORING_SPEC.md** - Technical Specification
- Detailed metrics definitions
- Code examples for instrumentation
- Middleware implementation
- Service-by-service instrumentation guide
- Dashboard specifications with PromQL queries

**Use this** for detailed technical implementation details and code examples.

### 3. **MONITORING_QUICKSTART.md** - Quick Start Guide
- Step-by-step implementation (45 minutes)
- Minimal viable monitoring setup
- Testing instructions
- Troubleshooting guide

**Use this** to get monitoring up and running quickly.

### 4. **DASHBOARD_EXAMPLES.md** - Dashboard Configurations
- Complete Grafana dashboard JSON
- Panel configurations
- Alert examples
- Import instructions

**Use this** for ready-to-use dashboard configurations.

## 🚀 Quick Start (45 Minutes)

Follow these steps to implement basic monitoring:

### Step 1: Infrastructure (5 min)
```bash
# Add Prometheus, Grafana, and postgres-exporter to docker-compose.yml
# See MONITORING_QUICKSTART.md section 1
```

### Step 2: Configuration (5 min)
```bash
# Create monitoring directory structure
mkdir -p monitoring/grafana/provisioning/{datasources,dashboards}
mkdir -p monitoring/grafana/dashboards

# Create Prometheus config (monitoring/prometheus.yml)
# Create Grafana provisioning configs
# See MONITORING_QUICKSTART.md section 1
```

### Step 3: Go Dependencies (2 min)
```bash
cd backend
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promauto
go get github.com/prometheus/client_golang/prometheus/promhttp
```

### Step 4: Create Metrics Package (10 min)
```bash
# Create backend/internal/metrics/metrics.go
# See MONITORING_QUICKSTART.md section 3
```

### Step 5: Add Middleware (5 min)
```bash
# Create backend/internal/middleware/prometheus.go
# See MONITORING_QUICKSTART.md section 4
```

### Step 6: Update Main Server (5 min)
```bash
# Update backend/cmd/server/main.go
# Add middleware and /metrics endpoint
# See MONITORING_QUICKSTART.md section 5
```

### Step 7: Instrument Services (10 min)
```bash
# Update backend/internal/services/jobqueue.go
# Add metrics calls
# See MONITORING_QUICKSTART.md section 6
```

### Step 8: Test (3 min)
```bash
# Start services
docker-compose up --build

# Verify metrics
curl http://localhost:8080/metrics | grep briefbot

# Access dashboards
# Prometheus: http://localhost:9090
# Grafana: http://localhost:3001 (admin/admin)
```

## 📊 What Gets Monitored

### Application Metrics
- ✅ HTTP requests (rate, duration, errors)
- ✅ Job queue (enqueued, processing, completed, failed)
- ✅ Worker performance
- ✅ Database connections
- ✅ Podcast generation
- ✅ AI API calls
- ✅ Email delivery
- ✅ External services (R2, Speech API, Scraping)
- ✅ SSE connections

### System Metrics
- ✅ PostgreSQL performance
- ✅ Connection pool usage
- ✅ Query performance

## 🎯 Key Metrics to Watch

### Critical (Red Flags)
1. **HTTP 5xx error rate > 1%** - Application errors
2. **Job failure rate > 5%** - Processing issues
3. **DB connections > 90%** - Connection pool exhaustion
4. **Request P95 latency > 2s** - Performance degradation
5. **Queue backlog > 500** - Processing bottleneck

### Performance
1. Request throughput (req/sec)
2. Request latency (P50, P95, P99)
3. Job processing throughput
4. External API latency

### Resources
1. Memory usage
2. CPU usage
3. Database connections
4. Active workers

## 📈 Dashboards

### 1. Overview Dashboard
- HTTP metrics (rate, errors, latency)
- Job processing status
- Database health
- Worker status

### 2. Job Processing Dashboard
- Queue depth by status
- Processing duration
- Success/failure rates
- Worker performance

### 3. External Services Dashboard
- AI API performance
- Speech API metrics
- R2 operations
- Email delivery
- Scraping success rate

### 4. Database Dashboard
- Connection pool usage
- Query performance
- Transaction rate
- Table sizes

## 🔧 Implementation Phases

### Phase 1: Core Infrastructure ✅ (Day 1)
- Add Prometheus, Grafana to docker-compose
- Create basic configuration
- Verify setup works

### Phase 2: Basic Metrics ✅ (Day 2)
- HTTP metrics via middleware
- Database connection metrics
- Basic dashboard

### Phase 3: Service Instrumentation (Day 3-4)
- Job queue metrics
- Worker metrics
- External service metrics
- Podcast generation metrics

### Phase 4: Dashboards (Day 5)
- Create comprehensive dashboards
- Add visualizations
- Configure alerts

### Phase 5: Documentation & Testing (Day 6)
- Update README
- Load testing
- Verify all metrics

### Phase 6: Alerting (Optional - Day 7)
- Define alert rules
- Configure notifications
- Test alerts

## 📁 File Structure After Implementation

```
briefbot/
├── docker-compose.yml (updated)
├── monitoring/
│   ├── README.md (this file)
│   ├── prometheus.yml
│   ├── alerts.yml (optional)
│   └── grafana/
│       ├── provisioning/
│       │   ├── datasources/
│       │   │   └── prometheus.yml
│       │   └── dashboards/
│       │       └── dashboard.yml
│       └── dashboards/
│           ├── overview.json
│           ├── jobs.json
│           └── external.json
├── backend/
│   ├── go.mod (updated)
│   ├── cmd/server/main.go (updated)
│   └── internal/
│       ├── metrics/
│       │   └── metrics.go (new)
│       ├── middleware/
│       │   └── prometheus.go (new)
│       └── services/ (updated)
└── README.md (updated)
```

## 🌐 Access URLs

After starting with `docker-compose up`:

- **Application**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Metrics Endpoint**: http://localhost:8080/metrics
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3001 (admin/admin)
- **Swagger API Docs**: http://localhost:8080/swagger/index.html
- **PostgreSQL**: localhost:5432

## 🧪 Testing

### Generate Test Traffic

```bash
# HTTP requests
for i in {1..100}; do curl http://localhost:8080/users; sleep 0.1; done

# Create jobs
for i in {1..10}; do
    curl -X POST http://localhost:8080/items \
        -H "Content-Type: application/json" \
        -d "{\"user_id\":1,\"url\":\"https://example.com/article-$i\"}"
done
```

### Verify Metrics

```bash
# Check metrics endpoint
curl http://localhost:8080/metrics | grep briefbot

# Check Prometheus targets
open http://localhost:9090/targets

# Check Grafana dashboards
open http://localhost:3001
```

## 🔍 Troubleshooting

### Metrics Not Showing
1. Check metrics endpoint: `curl http://localhost:8080/metrics`
2. Verify Prometheus scraping: http://localhost:9090/targets
3. Check backend logs: `docker-compose logs backend`

### Grafana Can't Connect
1. Verify datasource URL: `http://prometheus:9090`
2. Test connection in Grafana settings
3. Check Prometheus is running: `docker-compose ps prometheus`

### No Data in Dashboards
1. Generate test traffic
2. Wait 1-2 minutes for data collection
3. Verify queries in Prometheus: http://localhost:9090/graph
4. Check time range in Grafana (last 15 minutes)

## 📖 Useful Prometheus Queries

```promql
# Total request rate
sum(rate(briefbot_http_requests_total[5m]))

# Error rate percentage
sum(rate(briefbot_http_requests_total{status=~"5.."}[5m])) / sum(rate(briefbot_http_requests_total[5m])) * 100

# P95 latency
histogram_quantile(0.95, rate(briefbot_http_request_duration_seconds_bucket[5m]))

# Job success rate
rate(briefbot_jobs_completed_total[5m]) / rate(briefbot_jobs_enqueued_total[5m]) * 100

# Database connection utilization
briefbot_db_connections_active / (briefbot_db_connections_active + briefbot_db_connections_idle) * 100
```

## 🎓 Learning Resources

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Prometheus Best Practices](https://prometheus.io/docs/practices/naming/)
- [PromQL Tutorial](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [Grafana Dashboard Best Practices](https://grafana.com/docs/grafana/latest/dashboards/build-dashboards/best-practices/)

## 🚦 Next Steps

1. **Start with Quick Start**: Follow MONITORING_QUICKSTART.md
2. **Add More Metrics**: Use MONITORING_SPEC.md for detailed instrumentation
3. **Create Dashboards**: Import examples from DASHBOARD_EXAMPLES.md
4. **Set Up Alerts**: Configure critical alerts for production
5. **Monitor & Iterate**: Refine metrics and dashboards based on usage

## 💡 Best Practices

1. **Keep cardinality low**: Avoid high-cardinality labels (user IDs, timestamps)
2. **Use meaningful names**: Follow Prometheus naming conventions
3. **Monitor what matters**: Focus on business-critical metrics
4. **Set up alerts**: Don't wait for users to report issues
5. **Regular reviews**: Review dashboards weekly, refine as needed
6. **Document anomalies**: Keep notes on unusual patterns
7. **Test in dev**: Verify metrics work before production

## 🤝 Contributing

When adding new features to BriefBot:

1. Add relevant metrics in `internal/metrics/`
2. Instrument the code with metric calls
3. Update dashboards to include new metrics
4. Document new metrics in this README
5. Test metrics with load testing

## 📝 License

Same as BriefBot project (MIT)

## 🙋 Support

For questions or issues:
1. Check troubleshooting section above
2. Review Prometheus/Grafana logs
3. Consult official documentation
4. Open an issue in the project repository

---

**Ready to get started?** → Open [MONITORING_QUICKSTART.md](MONITORING_QUICKSTART.md)

**Need technical details?** → See [MONITORING_SPEC.md](MONITORING_SPEC.md)

**Want dashboards?** → Check [DASHBOARD_EXAMPLES.md](DASHBOARD_EXAMPLES.md)

**Understanding the plan?** → Read [plan.md](plan.md)
