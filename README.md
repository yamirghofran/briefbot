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

## Tech Stack

- React (Tanstack)
- HTML/CSS/Javascript for the browser extension
- ShadcnUI
- Go
- Gin
- Colly
- PostgreSQL
- Qdrant (Vector Database)
- FAL.ai (text-to-speech)
- Groq API (LLM)
- Cloudflare R2 (Object Storage)
