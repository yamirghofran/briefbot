#!/bin/bash
set -e

echo "🧪 Testing Docker builds..."
echo ""

# Test backend build
echo "📦 Building backend..."
docker build -t briefbot-backend:test -f backend/Dockerfile backend/
echo "✅ Backend build successful"
echo ""

# Test frontend build
echo "📦 Building frontend..."
docker build \
  --build-arg VITE_API_URL=http://localhost:8080 \
  -t briefbot-frontend:test \
  -f frontend/Dockerfile.prod \
  frontend/
echo "✅ Frontend build successful"
echo ""

echo "🎉 All builds completed successfully!"
echo ""
echo "To run the containers:"
echo "  docker run -d -p 8080:8080 briefbot-backend:test"
echo "  docker run -d -p 3000:80 briefbot-frontend:test"

