# BriefBot Setup Instructions

## Prerequisites

**Required:**

- **Docker Desktop** installed and running
  - macOS/Windows: [Download Docker Desktop](https://www.docker.com/products/docker-desktop/)
  - Linux: Install Docker Engine + Docker Compose plugin
  - Minimum version: Docker 20.10+ and Docker Compose v2.0+

**Not Required:**

- ❌ Go (handled by Docker)
- ❌ Bun/Node.js (handled by Docker)
- ❌ PostgreSQL (runs in Docker container)
- ❌ Any other dependencies

## Quick Start

### Step 1: Clone the Repository

```bash
git clone <repository-url>
cd briefbot
```

### Step 2: Add Environment File

**You should have received a `.env` file via email. Place it in the project root directory:**

```
briefbot/
├── .env          ← PUT THE FILE HERE (same level as docker-compose.yml)
├── docker-compose.yml
├── README.md
└── ...
```

#### On macOS/Linux:

```bash
# If you downloaded the .env file to Downloads folder:
mv ~/Downloads/.env .env

# Or create it manually:
nano .env
# (paste contents from email, save with Ctrl+X, then Y)
```

#### On Windows:

```bash
# If you downloaded to Downloads folder:
move %USERPROFILE%\Downloads\.env .env

# Or create it manually:
notepad .env
# (paste contents from email, save)
```

### Step 3: Verify File Location

```bash
# Check that .env exists in the right place
ls -la .env          # macOS/Linux
dir .env             # Windows
```

You should see the file listed. The filename must be exactly `.env` (not `.env.txt`).

### Step 4: Start the Application

```bash
docker-compose up --build
```

**What happens:**

1. Downloads required Docker images (first time only, ~5 minutes)
2. Builds the Go backend (~2-3 minutes)
3. Builds the React frontend (~3-4 minutes)
4. Starts PostgreSQL database
5. Runs database migrations
6. Seeds test data (professor user + sample items)
7. Starts all services

**Wait for this message:**

```
briefbot-backend    | Server starting on port 8080
briefbot-frontend   | ... (nginx started)
```

### Step 5: Access the Application

Once all services are running (30-60 seconds after startup):

- **Frontend UI**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Swagger API Docs**: http://localhost:8080/swagger/index.html
- **Go Documentation**: http://localhost:8081
- **Database**:
  - Host: localhost
  - Port: 5432
  - Database: briefbot
  - Username: briefbot
  - Password: briefbot

### Step 6: Test the Application

**Pre-seeded data:**

- User ID 1: Professor Demo (professor@university.edu)
- 10 sample articles (5 unread, 5 read)

**Try these actions:**

1. View all items in the UI
2. Mark items as read/unread
3. Add a new URL
4. View API documentation at http://localhost:8080/swagger/index.html
5. Explore Go package docs at http://localhost:8081

## Stopping the Application

```bash
# Stop all services (data is preserved)
docker-compose down

# Stop and remove all data (fresh start)
docker-compose down -v
```

## Troubleshooting

### Port Already in Use

If you see "port already in use" errors:

```bash
# Option 1: Stop conflicting services
# Find what's using the port
lsof -i :3000    # macOS/Linux
netstat -ano | findstr :3000    # Windows

# Option 2: Change ports in docker-compose.yml
# Edit docker-compose.yml and change port mappings
# For example: "3001:80" instead of "3000:80"
```

### Services Won't Start

```bash
# View logs to see what's wrong
docker-compose logs -f

# View logs for specific service
docker-compose logs -f backend
docker-compose logs -f postgres

# Restart fresh
docker-compose down -v
docker-compose up --build
```

### ".env file not found" Error

- Ensure `.env` is in the same directory as `docker-compose.yml`
- Check the filename is exactly `.env` (not `.env.txt`)
- On Windows, disable "Hide extensions for known file types" to verify

### Database Connection Issues

```bash
# Check if postgres is healthy
docker-compose ps

# View postgres logs
docker-compose logs postgres

# Connect to database manually
docker-compose exec postgres psql -U briefbot -d briefbot
```

### Build Errors

```bash
# Clean Docker cache and rebuild
docker-compose down -v
docker system prune -a
docker-compose up --build
```

### Frontend Not Loading

```bash
# Check if frontend container is running
docker-compose ps

# View frontend logs
docker-compose logs frontend

# Rebuild frontend only
docker-compose up --build frontend
```

## Development Commands

```bash
# View logs for all services
docker-compose logs -f

# View logs for specific service
docker-compose logs -f backend

# Execute command in backend container
docker-compose exec backend sh

# Execute command in database
docker-compose exec postgres psql -U briefbot -d briefbot

# Restart specific service
docker-compose restart backend

# Rebuild and restart specific service
docker-compose up --build -d backend
```

## System Requirements

- **RAM**: 4GB minimum, 8GB recommended
- **Disk**: 2GB free space for images and builds
- **OS**: macOS 10.15+, Windows 10+, or Linux with Docker support
- **Internet**: Required for initial download of Docker images

## Architecture Overview

```
┌─────────────────┐
│   Frontend      │  React + Vite (Bun)
│   Port: 3000    │  Served by nginx
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Backend       │  Go + Gin Framework
│   Port: 8080    │  REST API + Swagger
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   PostgreSQL    │  Database
│   Port: 5432    │  Latest version
└─────────────────┘

┌─────────────────┐
│   pkgsite       │  Go Documentation
│   Port: 8081    │  Package explorer
└─────────────────┘
```

## What's Included

✅ **Full Application Stack**

- React frontend with modern UI components
- Go backend with RESTful API
- PostgreSQL database with migrations
- Real-time updates via Server-Sent Events (SSE)

✅ **Development Tools**

- Swagger API documentation
- Go package documentation (pkgsite)
- Database access via psql

✅ **Test Data**

- Professor user account
- 10 sample articles
- Various item states (read/unread)

## Features Available

### With Real API Keys (from .env):

✅ AI-powered article summarization (GROQ)
✅ Podcast generation with text-to-speech (FAL)
✅ Email digest notifications (AWS SES)
✅ Audio storage (Cloudflare R2)

### Without API Keys:

✅ User management
✅ Item CRUD operations
✅ Mark items as read/unread
✅ View and organize articles
✅ Real-time UI updates
✅ API documentation

## Support

If you encounter any issues:

1. Check the troubleshooting section above
2. View service logs: `docker-compose logs -f`
3. Ensure Docker Desktop is running
4. Verify `.env` file is in the correct location
5. Try a fresh start: `docker-compose down -v && docker-compose up --build`

## Next Steps

After successful setup:

1. Explore the frontend UI at http://localhost:3000
2. Review API documentation at http://localhost:8080/swagger/index.html
3. Browse Go package docs at http://localhost:8081
4. Test API endpoints using Swagger UI
5. Add your own URLs and test the workflow
