# Docker Quick Start Guide

## 🚀 Quick Commands

### Build Images Locally

```bash
# Backend (production-ready)
docker build -t briefbot-backend:latest -f backend/Dockerfile backend/

# Frontend (with API URL)
docker build \
  --build-arg VITE_API_URL=http://localhost:8080 \
  -t briefbot-frontend:latest \
  -f frontend/Dockerfile.prod \
  frontend/
```

### Run Locally with Docker

```bash
# Start PostgreSQL
docker run -d \
  --name briefbot-postgres \
  -e POSTGRES_USER=briefbot \
  -e POSTGRES_PASSWORD=briefbot \
  -e POSTGRES_DB=briefbot \
  -p 5432:5432 \
  postgres:latest

# Start Backend
docker run -d \
  --name briefbot-backend \
  -p 8080:8080 \
  -e DATABASE_URL=postgres://briefbot:briefbot@host.docker.internal:5432/briefbot?sslmode=disable \
  briefbot-backend:latest

# Start Frontend
docker run -d \
  --name briefbot-frontend \
  -p 3000:80 \
  briefbot-frontend:latest
```

### Pull from DockerHub

```bash
# Replace <username> with your DockerHub username
docker pull <username>/briefbot-backend:latest
docker pull <username>/briefbot-frontend:latest
```

## 📝 Environment Configuration

### Local Development

Create `frontend/.env.local`:
```bash
VITE_API_URL=http://localhost:8080
```

### Azure Deployment

Set GitHub Actions variable:
```bash
# In GitHub repo settings → Secrets and variables → Actions → Variables
VITE_API_URL=https://your-backend.azurewebsites.net
```

## 🔄 CI/CD Workflow

The GitHub Actions workflow automatically builds and pushes images when you push to `main`:

1. Push changes to main branch
2. GitHub Actions runs tests
3. Builds Docker images for linux/amd64
4. Pushes to DockerHub with tags:
   - `latest`
   - `main-<git-sha>`

## 🧪 Testing Images

```bash
# Test backend health
docker run --rm briefbot-backend:latest /app/server --version

# Test frontend build
docker run --rm -p 8080:80 briefbot-frontend:latest
# Visit http://localhost:8080
```

## 📦 Image Tags

| Tag | Description | Use Case |
|-----|-------------|----------|
| `latest` | Latest main branch build | Development/Testing |
| `main-<sha>` | Specific commit | Rollback/Debugging |
| `v1.0.0` | Semantic version | Production releases |

## 🔍 Troubleshooting

### Check container logs
```bash
docker logs briefbot-backend
docker logs briefbot-frontend
```

### Inspect running container
```bash
docker exec -it briefbot-backend sh
docker exec -it briefbot-frontend sh
```

### Verify API URL in frontend
```bash
docker run --rm briefbot-frontend:latest cat /usr/share/nginx/html/assets/index-*.js | grep -o 'http[s]*://[^"]*' | head -1
```
