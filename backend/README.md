# BriefBot - AI-Powered Content Curation & Podcast Generation

A Go-based service that curates content from URLs, generates AI-powered summaries, and creates personalized podcast digests with text-to-speech conversion.

## 🚀 Features

- **Content Curation**: Submit URLs and get AI-generated summaries
- **Podcast Generation**: Convert multiple articles into audio podcasts with AI-generated scripts
- **Daily Digests**: Automated email digests with optional podcast audio
- **Asynchronous Processing**: Non-blocking API for long-running operations
- **Background Workers**: Concurrent processing of items and podcasts
- **Cloud Storage**: Audio files stored in Cloudflare R2 (S3-compatible)

## 🏗️ Architecture

### Core Components

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   API Layer     │    │  Service Layer  │    │   Data Layer    │
│  (Gin Router)   │───▶│   (Business)    │───▶│  (PostgreSQL)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

**API Layer**: HTTP endpoints for content submission and retrieval  
**Service Layer**: Business logic for AI processing, email, podcast generation  
**Data Layer**: PostgreSQL with connection pooling for data persistence  

### Background Processing

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Job Queue   │───▶│   Workers   │───▶│  Services   │
│  (Database)   │    │ (Concurrent)│    │(AI/Podcast) │
└─────────────┘    └─────────────┘    └─────────────┘
```

**Job Queue**: PostgreSQL-based queue for pending items/podcasts  
**Workers**: Concurrent goroutines that process jobs  
**Services**: AI content extraction, podcast script generation, text-to-speech  

## 🔧 API Endpoints

### Content Management

#### Submit Content for Processing
```bash
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "title": "Interesting Article",
    "url": "https://example.com/article"
  }'
```

#### Get User's Items
```bash
curl http://localhost:8080/items/user/1
```

#### Get Unread Items
```bash
curl http://localhost:8080/items/user/1/unread
```

#### Mark Item as Read
```bash
curl -X PATCH http://localhost:8080/items/123/read
```

### Podcast Generation

#### Create Podcast from Items
```bash
curl -X POST http://localhost:8080/podcasts \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "title": "My Daily Digest",
    "description": "Daily curated content",
    "item_ids": [1, 2, 3]
  }'
```

#### Create Podcast from Single Item
```bash
curl -X POST http://localhost:8080/podcasts/from-item \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "item_id": 123
  }'
```

#### Get User's Podcasts
```bash
curl http://localhost:8080/podcasts/user/1
```

#### Get Podcast Audio
```bash
curl http://localhost:8080/podcasts/123/audio
```

### Digest System (Unified Service)

#### Regular Daily Digest (No Podcast)
```bash
# Send to all users
curl -X POST http://localhost:8080/digest/trigger

# Send to specific user
curl -X POST http://localhost:8080/digest/trigger/user/1
```

#### Integrated Digest (With Podcast) - **ASYNCHRONOUS**
```bash
# Send to all users (returns immediately, processes in background)
curl -X POST http://localhost:8080/digest/trigger/integrated
# Response: 202 Accepted {"message": "Integrated digest processing started in background"}

# Send to specific user (returns immediately, processes in background)
curl -X POST http://localhost:8080/digest/trigger/integrated/user/1
# Response: 202 Accepted {"message": "Integrated digest processing started for user", "userID": 1}
```

## 🔄 How It Works

### 1. Content Processing Flow

```
URL Submission → Item Creation → Background Processing → AI Analysis → Email Digest
     ↓              ↓                ↓                ↓              ↓
  Pending       Queued          Worker Picks    Groq AI API    Daily Email
  Status        for Work        Up Job          Processing     with Summaries
```

**Step 1**: Submit URL → Creates pending item in database  
**Step 2**: Background worker picks up item → Scrapes content  
**Step 3**: AI service analyzes content → Generates title, summary, metadata  
**Step 4**: Item marked as completed → Available for digest generation  

### 2. Podcast Generation Flow

```
Podcast Request → Script Generation → Audio Generation → File Upload → Ready
      ↓               ↓                ↓               ↓           ↓
   Items Listed    AI Dialogue     Text-to-Speech    R2 Storage   Audio URL
   for Podcast     Scripting       (Fal.ai)         (MP3 File)   Available
```

**Step 1**: Submit items for podcast → Creates pending podcast  
**Step 2**: AI generates dialogue script → Conversational format  
**Step 3**: Text-to-speech conversion → Multiple audio segments  
**Step 4**: Audio stitching & upload → Single MP3 file in R2  
**Step 5**: Podcast ready → Audio URL available for download  

### 3. Integrated Digest Flow (NEW - ASYNCHRONOUS)

```
Digest Trigger → Fetch Items → Generate Podcast → Create Email → Send Email
      ↓            ↓             ↓               ↓           ↓
  202 Accepted   Get Unread   Background TTS   Include     Delivered
  (Immediate)    Items        & Upload         Audio Link  to User
```

**Step 1**: API call → Returns 202 immediately  
**Step 2**: Background processing → Fetches user's unread items  
**Step 3**: Podcast generation → AI script + TTS + R2 upload  
**Step 4**: Email creation → Includes podcast download link at top  
**Step 5**: Email delivery → User receives digest with audio option  

## ⚙️ Configuration

### Environment Variables

```bash
# Database
DATABASE_URL=postgres://user:pass@localhost:5432/briefbot?sslmode=disable

# AI Service (Groq/OpenAI Compatible)
GROQ_API_KEY=your_groq_api_key

# Text-to-Speech (Fal.ai)
FAL_API_KEY=your_fal_api_key

# Email Service (AWS SES)
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_REGION=us-east-1
SES_FROM_EMAIL=noreply@yourdomain.com
SES_FROM_NAME=BriefBot

# Cloud Storage (Cloudflare R2)
R2_ACCESS_KEY_ID=your_r2_access_key
R2_SECRET_ACCESS_KEY=your_r2_secret_key
R2_ACCOUNT_ID=your_account_id
R2_BUCKET_NAME=briefbot-audio
R2_PUBLIC_HOST=https://audio.yourdomain.com

# Optional Settings
MAX_CONCURRENT_AUDIO_REQUESTS=5  # Default: 5
DIGEST_PODCAST_ENABLED=true        # Enable podcast generation in digests
```

### Worker Configuration

```go
WorkerConfig{
    WorkerCount:    2,               // Number of concurrent workers
    PollInterval:   5 * time.Second, // How often to check for jobs
    MaxRetries:     3,               // Max retries for failed jobs
    BatchSize:      10,              // Jobs per batch
    EnablePodcasts: true,            // Process podcasts
}
```

## 🚀 Quick Start

### 1. Clone and Setup
```bash
git clone <repository>
cd briefbot/backend
cp .env.example .env  # Configure your environment variables
```

### 2. Database Setup
```bash
# Run migrations
goose -dir sql/migrations postgres "$DATABASE_URL" up

# Generate SQLC code (if modifying queries)
sqlc generate
```

### 3. Run the Server
```bash
go run cmd/server/main.go
# Server starts on port 8080 (or PORT env var)
```

### 4. Test the Flow
```bash
# 1. Submit content
curl -X POST localhost:8080/items -d '{"user_id":1,"url":"https://example.com/article"}'

# 2. Wait for processing (check status)
curl localhost:8080/items/user/1/unread

# 3. Trigger integrated digest (async)
curl -X POST localhost:8080/digest/trigger/integrated/user/1
# Returns immediately: 202 Accepted

# 4. Check logs for completion
# Background processing will complete in 2-5 minutes
```

## 📊 Performance Characteristics

- **Content Processing**: ~30-60 seconds per item (scraping + AI analysis)
- **Podcast Generation**: ~2-5 minutes (script writing + TTS + audio stitching)
- **Email Delivery**: ~1-2 seconds per user
- **Concurrent Audio**: Up to 5 simultaneous TTS requests (configurable)
- **Database**: Connection pooled (25 max connections)

## 🔍 Monitoring

The application provides comprehensive logging:
- Job processing status and timing
- AI service calls and responses
- Podcast generation progress
- Email delivery success/failure
- Background worker activity

## 🛠️ Development

### Testing
```bash
go test ./internal/services -v  # Run service tests
go test ./internal/db -v      # Run database tests
```

### Code Generation
```bash
sqlc generate  # Regenerate database models from SQL queries
```

### Linting
```bash
golangci-lint run  # Run linter
```

## 📚 API Documentation

The API uses Swagger/OpenAPI annotations. When running the server, documentation is available at:
- Swagger UI: `http://localhost:8080/swagger/index.html`
- OpenAPI JSON: `http://localhost:8080/swagger/doc.json`

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## 📄 License

[Your License Here]

---

**BriefBot** - Your AI-powered content curator and podcast generator 🎧