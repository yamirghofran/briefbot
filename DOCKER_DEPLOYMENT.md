# Docker Deployment Guide

This guide explains how to build, deploy, and configure the Briefbot application using Docker containers.

## 📦 Container Images

The application consists of two main containers:

- **Backend**: Go API server with PostgreSQL migrations
- **Frontend**: React SPA served by nginx

## 🏗️ Building Images

### Backend

```bash
# Development build
docker build -t briefbot-backend:dev -f backend/Dockerfile backend/

# Production build with specific platform
docker build --platform linux/amd64 -t briefbot-backend:prod -f backend/Dockerfile backend/
```

### Frontend

The frontend requires the `VITE_API_URL` build argument to configure the API endpoint:

```bash
# Development build (localhost)
docker build \
  --build-arg VITE_API_URL=http://localhost:8080 \
  -t briefbot-frontend:dev \
  -f frontend/Dockerfile.prod \
  frontend/

# Production build (Azure)
docker build \
  --build-arg VITE_API_URL=https://your-backend.azurewebsites.net \
  -t briefbot-frontend:prod \
  -f frontend/Dockerfile.prod \
  frontend/

# Staging build
docker build \
  --build-arg VITE_API_URL=https://staging-backend.azurewebsites.net \
  -t briefbot-frontend:staging \
  -f frontend/Dockerfile.prod \
  frontend/
```

## 🚀 Running Containers Locally

### Backend

```bash
docker run -d \
  --name briefbot-backend \
  -p 8080:8080 \
  -e DATABASE_URL=postgres://user:pass@host:5432/briefbot \
  -e PORT=8080 \
  briefbot-backend:dev
```

### Frontend

```bash
docker run -d \
  --name briefbot-frontend \
  -p 3000:80 \
  briefbot-frontend:dev
```

## 🌐 Deployment Scenarios

### Azure Container Apps

```bash
# Tag images for Azure Container Registry
docker tag briefbot-backend:prod yourregistry.azurecr.io/briefbot-backend:latest
docker tag briefbot-frontend:prod yourregistry.azurecr.io/briefbot-frontend:latest

# Push to ACR
docker push yourregistry.azurecr.io/briefbot-backend:latest
docker push yourregistry.azurecr.io/briefbot-frontend:latest

# Deploy using Azure CLI
az containerapp create \
  --name briefbot-backend \
  --resource-group briefbot-rg \
  --environment briefbot-env \
  --image yourregistry.azurecr.io/briefbot-backend:latest \
  --target-port 8080 \
  --env-vars DATABASE_URL=secretref:database-url

az containerapp create \
  --name briefbot-frontend \
  --resource-group briefbot-rg \
  --environment briefbot-env \
  --image yourregistry.azurecr.io/briefbot-frontend:latest \
  --target-port 80
```

### Azure App Service (Web App for Containers)

```bash
# Create App Service Plan
az appservice plan create \
  --name briefbot-plan \
  --resource-group briefbot-rg \
  --is-linux \
  --sku B1

# Create backend web app
az webapp create \
  --name briefbot-backend \
  --resource-group briefbot-rg \
  --plan briefbot-plan \
  --deployment-container-image-name yourregistry.azurecr.io/briefbot-backend:latest

# Configure backend environment variables
az webapp config appsettings set \
  --name briefbot-backend \
  --resource-group briefbot-rg \
  --settings DATABASE_URL="postgres://..." PORT=8080

# Create frontend web app
az webapp create \
  --name briefbot-frontend \
  --resource-group briefbot-rg \
  --plan briefbot-plan \
  --deployment-container-image-name yourregistry.azurecr.io/briefbot-frontend:latest
```

### DockerHub Deployment

Images are automatically published to DockerHub via GitHub Actions when pushing to the `main` branch.

```bash
# Pull from DockerHub
docker pull <your-dockerhub-username>/briefbot-backend:latest
docker pull <your-dockerhub-username>/briefbot-frontend:latest

# Run with docker-compose
docker-compose -f docker-compose.prod.yml up -d
```

## 🔧 Environment Variables

### Backend Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://user:pass@host:5432/briefbot?sslmode=disable` |
| `PORT` | Server port | `8080` |
| `FRONTEND_BASE_URL` | Frontend URL for CORS | `https://briefbot.com` |
| `R2_ACCESS_KEY_ID` | Cloudflare R2 access key | `your-access-key` |
| `R2_SECRET_ACCESS_KEY` | Cloudflare R2 secret key | `your-secret-key` |
| `R2_ACCOUNT_ID` | Cloudflare R2 account ID | `your-account-id` |
| `R2_BUCKET_NAME` | Cloudflare R2 bucket name | `briefbot-storage` |
| `R2_PUBLIC_HOST` | Cloudflare R2 public host | `https://pub-xxx.r2.dev` |
| `OPENAI_API_KEY` | OpenAI API key | `sk-...` |
| `RESEND_API_KEY` | Resend email API key | `re_...` |

### Frontend Build Arguments

| Variable | Description | Example |
|----------|-------------|---------|
| `VITE_API_URL` | Backend API URL (build-time) | `https://api.briefbot.com` |

## 🔍 Health Checks

Both containers include health check endpoints:

### Backend
```bash
curl http://localhost:8080/health
```

### Frontend
```bash
curl http://localhost:80/health
```

## 📊 CI/CD Pipeline

The GitHub Actions workflow automatically:

1. ✅ Runs tests and linting
2. ✅ Builds Docker images for linux/amd64
3. ✅ Tags images with:
   - `latest` (main branch)
   - `main-<git-sha>` (commit-specific)
4. ✅ Pushes to DockerHub
5. ✅ Generates build summary

### Triggering Builds

```bash
# Push to main branch triggers automatic build
git push origin main

# Create a release tag for versioned builds
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

## 🛠️ Multi-Environment Setup

### Development
```bash
# Use localhost API
docker build --build-arg VITE_API_URL=http://localhost:8080 -t briefbot-frontend:dev frontend/
```

### Staging
```bash
# Use staging API
docker build --build-arg VITE_API_URL=https://staging-api.briefbot.com -t briefbot-frontend:staging frontend/
```

### Production
```bash
# Use production API
docker build --build-arg VITE_API_URL=https://api.briefbot.com -t briefbot-frontend:prod frontend/
```

## 🔐 Security Best Practices

1. **Never commit secrets** - Use environment variables or secret managers
2. **Use multi-stage builds** - Reduces image size and attack surface
3. **Scan images** - Use tools like Trivy or Snyk
4. **Update base images** - Regularly update Alpine and Go versions
5. **Use specific tags** - Avoid `latest` in production
6. **Enable health checks** - Ensure containers are monitored

## 📝 Troubleshooting

### Frontend can't connect to backend

**Problem**: API calls fail with CORS or connection errors

**Solution**: Ensure `VITE_API_URL` matches your backend URL:
```bash
# Check the built config
docker run briefbot-frontend:dev cat /usr/share/nginx/html/assets/index-*.js | grep -o 'http[s]*://[^"]*'
```

### Backend health check fails

**Problem**: Container restarts due to failed health checks

**Solution**: Verify the `/health` endpoint exists and returns 200:
```bash
docker exec briefbot-backend curl -f http://localhost:8080/health
```

### Image size too large

**Problem**: Docker images are consuming too much space

**Solution**: 
- Ensure `.dockerignore` is properly configured
- Use multi-stage builds (already implemented)
- Clean up build cache: `docker builder prune`

## 📚 Additional Resources

- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Azure Container Apps Documentation](https://learn.microsoft.com/en-us/azure/container-apps/)
- [GitHub Actions Docker Build](https://docs.github.com/en/actions/publishing-packages/publishing-docker-images)
