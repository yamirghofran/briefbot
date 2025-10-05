# BriefBot

## Overview

An AI-enabled platform for managing links and extracting knowledge faster.

- Browser Plugin
- Web UI
- Daily Digest Email Notifications

## 🐳 Quick Start with Docker

**For evaluators/professors**: The easiest way to run BriefBot is with Docker.

```bash
# 1. Ensure you have the .env file in the project root
# 2. Run everything with one command
docker-compose up --build

# 3. Access the application
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
# Swagger Docs: http://localhost:8080/swagger/index.html
# Go Docs: http://localhost:8081
```

**📖 Full Docker Instructions**: See [SETUP.md](SETUP.md) for detailed setup instructions.

**📚 Docker Documentation**: See [README.docker.md](README.docker.md) for comprehensive Docker documentation.

### Prerequisites for Docker Setup

- Docker Desktop installed ([Download here](https://www.docker.com/products/docker-desktop/))
- `.env` file in project root (provided separately)
- No other dependencies needed!

## 📚 Built-in Documentation

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

## Tech Stack

- React (Tanstack)
- ShadcnUI
- Go
- Gin
- Colly
- PostgreSQL
- Qdrant (Vector Database)
- FAL.ai (text-to-speech)
- Groq API (LLM)
- Cloudflare R2 (Object Storage)
