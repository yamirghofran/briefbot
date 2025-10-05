# 🐳 BriefBot Docker Documentation

Complete guide to running BriefBot using Docker for development, testing, and demonstration purposes.

## Table of Contents

- [Quick Start](#quick-start)
- [Prerequisites](#prerequisites)
- [Architecture](#architecture)
- [Services](#services)
- [Environment Configuration](#environment-configuration)
- [Common Tasks](#common-tasks)
- [Troubleshooting](#troubleshooting)
- [Development](#development)

## Quick Start

**For Professor/Evaluator:**

See [SETUP.md](SETUP.md) for step-by-step instructions.

**TL;DR:**

```bash
# 1. Ensure .env file is in project root
# 2. Run everything
docker-compose up --build

# 3. Access application
# Frontend: http://localhost:3000
# Backend: http://localhost:8080
# Swagger: http://localhost:8080/swagger/index.html
# Go Docs: http://localhost:8081
```

## Prerequisites

### Required Software

1. **Docker Desktop** (or Docker Engine + Docker Compose)

   - Download: https://www.docker.com/products/docker-desktop/
   - Minimum version: Docker 20.10+, Docker Compose v2.0+
   - Ensure Docker is running before starting

2. **Git** (to clone repository)
   - Usually pre-installed on macOS/Linux
   - Windows: https://git-scm.com/download/win

### System Requirements

- **CPU**: 2+ cores recommended
- **RAM**: 4GB minimum, 8GB recommended
- **Disk**: 2GB free space for images and builds
- **OS**: macOS 10.15+, Windows 10+, or Linux with Docker support

### Not Required

The following are **NOT** needed (Docker handles them):

- ❌ Go installation
- ❌ Bun/Node.js installation
- ❌ PostgreSQL installation
- ❌ Any Go packages or npm modules

## Architecture

### Service Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     Docker Network                           │
│                                                              │
│  ┌──────────────┐      ┌──────────────┐      ┌───────────┐ │
│  │   Frontend   │─────▶│   Backend    │─────▶│ PostgreSQL│ │
│  │  React+Vite  │      │   Go+Gin     │      │  Database │ │
│  │   :3000      │      │    :8080     │      │   :5432   │ │
│  └──────────────┘      └──────────────┘      └───────────┘ │
│                              │                              │
│                              │                              │
│                        ┌─────▼──────┐                       │
│                        │  pkgsite   │                       │
│                        │  Go Docs   │                       │
│                        │   :8081    │                       │
│                        └────────────┘                       │
│                                                              │
│  ┌──────────────┐      ┌──────────────┐                    │
│  │  Migrations  │      │     Seed     │                    │
│  │  (one-time)  │─────▶│  (one-time)  │                    │
│  └──────────────┘      └──────────────┘                    │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### Service Dependencies

```
postgres (healthy)
    ↓
migrations (runs once)
    ↓
seed (runs once)
    ↓
backend (runs continuously)
    ↓
frontend (runs continuously)

pkgsite (independent)
```

## Services

### 1. PostgreSQL Database

**Image**: `postgres:latest`
**Port**: `5432:5432`
**Purpose**: Main application database

**Configuration**:

- Database: `briefbot`
- Username: `briefbot`
- Password: `briefbot`
- Persistent volume: `postgres_data`

**Health Check**: Runs `pg_isready` every 5 seconds

**Access**:

```bash
# Via docker-compose
docker-compose exec postgres psql -U briefbot -d briefbot

# Via local psql (if installed)
psql -h localhost -U briefbot -d briefbot
```

### 2. Migrations Service

**Build**: `./backend/Dockerfile`
**Purpose**: Run Goose migrations on startup
**Runs**: Once (on-failure restart policy)

**What it does**:

1. Waits for PostgreSQL to be ready
2. Runs all migrations in `backend/sql/migrations/`
3. Exits when complete

**Logs**:

```bash
docker-compose logs migrations
```

### 3. Seed Service

**Image**: `postgres:latest`
**Purpose**: Populate database with test data
**Runs**: Once (after migrations complete)

**What it seeds**:

- Professor user (ID: 1, email: professor@university.edu)
- 10 sample articles (5 unread, 5 read)
- Sample podcast (if table exists)

**Logs**:

```bash
docker-compose logs seed
```

### 4. Backend Service

**Build**: `./backend/Dockerfile`
**Port**: `8080:8080`
**Purpose**: Go REST API server

**Features**:

- RESTful API endpoints
- Swagger documentation at `/swagger/index.html`
- Server-Sent Events (SSE) for real-time updates
- Background workers for async processing

**Environment**: Loads from `.env` file

**Logs**:

```bash
docker-compose logs -f backend
```

**Restart**:

```bash
docker-compose restart backend
```

### 5. Frontend Service

**Build**: `./frontend/Dockerfile`
**Port**: `3000:3000`
**Purpose**: React web application

**Tech Stack**:

- React 19
- Vite dev server
- Bun package manager
- TanStack Router & Query

**Logs**:

```bash
docker-compose logs -f frontend
```

### 6. pkgsite Service

**Build**: `./backend/Dockerfile.pkgsite`
**Port**: `8081:8081`
**Purpose**: Go package documentation server

**Access**: http://localhost:8081

**Features**:

- Browse Go package documentation
- View function signatures
- Explore code structure

## Environment Configuration

### .env File Structure

The `.env` file should be placed in the **project root** (same directory as `docker-compose.yml`).

**Template**: See `.env.example` for all available variables.

### Required Variables

```env
# Database (Docker internal networking)
DATABASE_URL=postgres://briefbot:briefbot@postgres:5432/briefbot?sslmode=disable
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://briefbot:briefbot@postgres:5432/briefbot?sslmode=disable
GOOSE_MIGRATION_DIR=sql/migrations

# Application
PORT=8080
FRONTEND_BASE_URL=http://localhost:3000
```

### Optional Variables (for full functionality)

```env
# AI Services
GROQ_API_KEY=your_groq_key          # For AI summarization
FAL_API_KEY=your_fal_key            # For text-to-speech

# Email (AWS SES)
AWS_ACCESS_KEY_ID=your_key
AWS_SECRET_ACCESS_KEY=your_secret
AWS_REGION=us-east-1
SES_FROM_EMAIL=your@email.com

# Storage (Cloudflare R2)
R2_ACCESS_KEY_ID=your_key
R2_SECRET_ACCESS_KEY=your_secret
R2_ACCOUNT_ID=your_account
R2_BUCKET_NAME=briefbot
R2_PUBLIC_HOST=https://your-bucket.r2.dev

# Feature Flags
DIGEST_PODCAST_ENABLED=true
MAX_CONCURRENT_AUDIO_REQUESTS=5
```

### What Works Without API Keys

✅ **Core Features** (no API keys needed):

- User management
- Item CRUD operations
- Mark items as read/unread
- View and organize articles
- Real-time UI updates
- Database operations
- API documentation

❌ **Advanced Features** (require API keys):

- AI-powered summarization
- Podcast generation
- Email notifications
- Audio storage

## Common Tasks

### Starting the Application

```bash
# Start all services (build if needed)
docker-compose up --build

# Start in background (detached mode)
docker-compose up -d --build

# Start specific service
docker-compose up backend
```

### Stopping the Application

```bash
# Stop all services (preserve data)
docker-compose down

# Stop and remove volumes (fresh start)
docker-compose down -v

# Stop specific service
docker-compose stop backend
```

### Viewing Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f postgres

# Last 100 lines
docker-compose logs --tail=100 backend
```

### Rebuilding Services

```bash
# Rebuild all services
docker-compose build

# Rebuild specific service
docker-compose build backend

# Rebuild and restart
docker-compose up --build -d
```

### Database Operations

```bash
# Connect to database
docker-compose exec postgres psql -U briefbot -d briefbot

# Run SQL file
docker-compose exec -T postgres psql -U briefbot -d briefbot < script.sql

# Backup database
docker-compose exec postgres pg_dump -U briefbot briefbot > backup.sql

# Restore database
docker-compose exec -T postgres psql -U briefbot -d briefbot < backup.sql

# Reset database (WARNING: deletes all data)
docker-compose down -v
docker-compose up -d postgres
docker-compose up migrations seed
```

### Running Commands in Containers

```bash
# Backend shell
docker-compose exec backend sh

# Run Go tests
docker-compose exec backend go test ./...

# Check Go version
docker-compose exec backend go version

# Frontend shell
docker-compose exec frontend sh
```

### Checking Service Status

```bash
# List running services
docker-compose ps

# Check service health
docker-compose ps postgres

# View resource usage
docker stats
```

## Troubleshooting

### Services Won't Start

**Problem**: Containers exit immediately or fail to start

**Solutions**:

```bash
# Check logs for errors
docker-compose logs

# Verify .env file exists
ls -la .env

# Check Docker is running
docker info

# Clean and rebuild
docker-compose down -v
docker system prune -a
docker-compose up --build
```

### Port Conflicts

**Problem**: "Port already in use" error

**Solutions**:

```bash
# Find what's using the port
lsof -i :3000          # macOS/Linux
netstat -ano | findstr :3000    # Windows

# Option 1: Stop conflicting service
# Option 2: Change port in docker-compose.yml
# Change "3000:80" to "3001:80"
```

### Database Connection Errors

**Problem**: Backend can't connect to database

**Solutions**:

```bash
# Check postgres is healthy
docker-compose ps postgres

# View postgres logs
docker-compose logs postgres

# Verify DATABASE_URL uses 'postgres' as host (not 'localhost')
# Correct: postgres://briefbot:briefbot@postgres:5432/...
# Wrong:   postgres://briefbot:briefbot@localhost:5432/...

# Restart services in order
docker-compose down
docker-compose up postgres
# Wait for healthy
docker-compose up backend
```

### Migration Failures

**Problem**: Migrations service fails or exits with error

**Solutions**:

```bash
# View migration logs
docker-compose logs migrations

# Run migrations manually
docker-compose exec postgres psql -U briefbot -d briefbot
# Then check what migrations ran:
# SELECT * FROM goose_db_version;

# Reset and retry
docker-compose down -v
docker-compose up --build
```

### Frontend Build Errors

**Problem**: Frontend fails to build or shows blank page

**Solutions**:

```bash
# View frontend logs
docker-compose logs frontend

# Rebuild frontend
docker-compose build --no-cache frontend
docker-compose up -d frontend

# Check frontend is running
docker-compose exec frontend bun --version
```

### Slow Build Times

**Problem**: Docker builds take too long

**Solutions**:

```bash
# Use BuildKit (faster builds)
export DOCKER_BUILDKIT=1
docker-compose build

# Build in parallel
docker-compose build --parallel

# Use cached layers (don't use --no-cache unless necessary)
docker-compose build
```

### Out of Disk Space

**Problem**: Docker runs out of space

**Solutions**:

```bash
# Check disk usage
docker system df

# Clean up unused resources
docker system prune -a

# Remove all volumes (WARNING: deletes data)
docker volume prune
```

## Development

### Hot Reload (Optional)

To enable hot reload for development, modify `docker-compose.yml`:

```yaml
backend:
  volumes:
    - ./backend:/app
  command: ["air"] # Requires air to be installed in Dockerfile

frontend:
  volumes:
    - ./frontend:/app
    - /app/node_modules
  command: ["bun", "run", "dev"]
```

### Running Tests

```bash
# Backend tests
docker-compose exec backend go test ./...

# Backend tests with coverage
docker-compose exec backend go test -coverprofile=coverage.out ./...

# Frontend tests (if configured)
docker-compose exec frontend bun test
```

### Adding New Migrations

```bash
# Create new migration
docker-compose exec backend goose -dir sql/migrations create migration_name sql

# Run new migrations
docker-compose restart migrations
```

### Debugging

```bash
# Enable verbose logging
docker-compose up --build --verbose

# Inspect container
docker-compose exec backend sh
docker inspect briefbot-backend

# Check environment variables
docker-compose exec backend env
```

## Performance Optimization

### Build Optimization

- Multi-stage builds minimize final image size
- Layer caching speeds up rebuilds
- `.dockerignore` files reduce build context

### Runtime Optimization

- PostgreSQL connection pooling (configured in backend)
- Vite dev server with HMR (Hot Module Replacement)
- Health checks prevent premature requests

### Resource Limits (Optional)

Add to `docker-compose.yml`:

```yaml
services:
  backend:
    deploy:
      resources:
        limits:
          cpus: "1"
          memory: 512M
```

## Security Notes

### For Production

⚠️ **This Docker setup is for development/demonstration only.**

For production:

- Change default database credentials
- Use secrets management (Docker secrets, vault)
- Enable SSL/TLS
- Configure proper firewall rules
- Use non-root users in containers
- Scan images for vulnerabilities
- Keep images updated

### .env File Security

- Never commit `.env` to version control
- Use `.env.example` as template
- Rotate API keys regularly
- Limit permissions: `chmod 600 .env`

## Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [PostgreSQL Docker Hub](https://hub.docker.com/_/postgres)
- [Go Docker Best Practices](https://docs.docker.com/language/golang/)

## Support

For issues:

1. Check logs: `docker-compose logs -f`
2. Review troubleshooting section above
3. Try fresh start: `docker-compose down -v && docker-compose up --build`
4. Check Docker Desktop is running and healthy
