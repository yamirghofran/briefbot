# BriefBot

## Overview

An AI-enabled platform for managing links and extracting knowledge faster.

When you save a link with Briefbot, it extracts the metadata such as a proper title, authors, platform, and tags. It also generates a summary that covers the key topics for when you want to skim over it.
You also have the option to select multiple items and generate an engaging NotebookLM style podcast about those items.
Finally, when you click the "Trigger Digest" button, it sends you the summaries and a podcast about the items you saved yesterday but didn't read.
You can filter you items based on type, author, platform or search over them.

## 🐳 Quick Start with Docker

**For evaluators/professors**: The easiest way to run BriefBot is with Docker.

1. Make sure you have Docker Desktop installed and running.
2. Clone this repository and navigate to it.
3. Create a `.env` file in project root similar to the example.
4. Run `docker-compose up --build`
5. Navigate to http://localhost:3000 to use the app.

```bash
# 1. Ensure you have the .env file in the project root
# 2. Run everything with one command
docker-compose up --build
```

You can access:

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Swagger Docs: http://localhost:8080/swagger/index.html
- Go Docs: http://localhost:8081/github.com/yamirghofran/briefbot
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3001 (admin/admin)
- **Metrics Endpoint**: http://localhost:8080/metrics

**📖 Full Docker Instructions**: See [SETUP.md](SETUP.md) for detailed setup instructions.

**📚 Docker Documentation**: See [README.docker.md](README.docker.md) for comprehensive Docker documentation.

### Prerequisites for Docker Setup

- Docker Desktop installed ([Download here](https://www.docker.com/products/docker-desktop/))
- `.env` file in project root (provided separately)
- No other dependencies needed!

### Setting Up the Extension

1. Navigate to you browsers Extensions -> Manage Extensions settings
2. Turn on "Developer Mode"
3. Click on "Load Unpacked"
4. Select the "browser-plugin" folder from this repo.
5. Pin the extension.

You can then use the extension on the webpage you want to save without having to go to the briefbot website.

## Built-in Documentation

BriefBot includes comprehensive documentation tools that run automatically with Docker:

### Swagger API Documentation (Port 8080)

Interactive REST API documentation generated from code annotations.

**URL**: http://localhost:8080/swagger/index.html

**Features**:

- **Complete API Reference**: All endpoints with descriptions, parameters, and response schemas
- **Try It Out**: Test API endpoints directly in the browser without curl or Postman
- **Request/Response Examples**: See example payloads for every endpoint
- **Schema Explorer**: Browse data models and their fields
- **Organized by Tags**: Endpoints grouped by feature (users, items, podcasts, digest)

**Example Endpoints**:

- `GET /users` - List all users
- `POST /items` - Create new item
- `GET /items/user/:userID` - Get user's items
- `POST /daily-digest/trigger` - Send digest emails

### Go Package Documentation - pkgsite (Port 8081)

Official Go documentation server (same as pkg.go.dev) running locally for BriefBot's codebase.

**URL**: http://localhost:8081

**Features**:

- **Package Explorer**: Browse all Go packages in the project
- **Source Code Navigation**: Jump to function definitions and implementations
- **Function Signatures**: View all exported functions, types, and constants
- **Code Examples**: See usage examples from comments
- **Cross-References**: Navigate between related packages and types
- **Package Dependencies**: Understand how packages relate to each other

**Key Packages to Explore**:

- `internal/handlers` - HTTP request handlers
- `internal/services` - Business logic layer
- `internal/db` - Database queries and models
- `cmd/server` - Application entry point

### Why Both?

**Swagger** focuses on the **HTTP API interface** - what external clients see and use.

**pkgsite** focuses on the **internal Go code** - how the application is structured and implemented.

Together, they provide complete documentation from API consumer perspective (Swagger) and developer perspective (pkgsite).

### DeepWiki

Extensive documentation is also available at [DeepWiki](https://deepwiki.com/yamirghofran/briefbot)

## 📊 Monitoring & Observability

BriefBot includes comprehensive monitoring with **Prometheus** and **Grafana** for real-time visibility into application performance, job processing, and system health.

### Monitoring Stack

- **Prometheus** (Port 9090) - Metrics collection and storage
- **Grafana** (Port 3001) - Metrics visualization and dashboards
- **PostgreSQL Exporter** (Port 9187) - Database metrics
- **Metrics Endpoint** (Port 8080/metrics) - Application metrics

### Key Metrics Tracked

**Application Metrics:**
- HTTP request rate, latency (P50, P95, P99), and error rates
- Active HTTP requests and response sizes
- Job queue depth by status (pending, processing, completed, failed)
- Job processing duration and throughput
- Worker performance and utilization
- Database connection pool usage

**External Services:**
- AI API calls, latency, and errors
- Podcast generation metrics
- Email delivery success/failure rates
- Scraping service performance
- SSE connection tracking

### Quick Start

1. Start all services with `docker-compose up --build`
2. Access Grafana at http://localhost:3001 (login: admin/admin)
3. View the "BriefBot Overview" dashboard
4. Explore metrics in Prometheus at http://localhost:9090

### Useful Prometheus Queries

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

### Monitoring Documentation

For detailed monitoring documentation, see:
- **[MONITORING_README.md](MONITORING_README.md)** - Overview and quick start
- **[MONITORING_QUICKSTART.md](MONITORING_QUICKSTART.md)** - Step-by-step setup guide
- **[MONITORING_SPEC.md](MONITORING_SPEC.md)** - Technical specifications
- **[DASHBOARD_EXAMPLES.md](DASHBOARD_EXAMPLES.md)** - Dashboard configurations

## Tech Stack

**Frontend:**
- React (Tanstack)
- HTML/CSS/Javascript for the browser extension
- ShadcnUI

**Backend:**
- Go
- Gin (HTTP framework)
- Colly (web scraping)
- PostgreSQL (database)

**AI & External Services:**
- Groq API (LLM for summarization)
- FAL.ai (text-to-speech for podcasts)
- Cloudflare R2 (object storage)
- AWS SES (email delivery)

**Monitoring & Observability:**
- Prometheus (metrics collection)
- Grafana (visualization)
- PostgreSQL Exporter (database metrics)
