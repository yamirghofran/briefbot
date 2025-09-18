# BriefBot

## Overview

An AI-enabled platform for managing links and extracting knowledge faster.

- Browser Plugin
- Telegram Bot Integration
- Web UI
- Daily Digest Email Notifications

## Tech Stack

- React (Tanstack)
- ShadcnUI
- Go
- Gin
- Colly
- PostgreSQL
- Qdrant (Vector Database)
- Redis (Background jobs and pipeline management)
- Eleven Labs API (text-to-speech)
- Groq API (LLM)
- Cloudflare Browser Rendering
- Cloudflare R2 (Object Storage)

## Environment Variables:

### Core Variables:
• DATABASE_URL: PostgreSQL connection string
• PORT: Server port (default: 8080)

### Email Service Variables (Required for Daily Digest):
• AWS_ACCESS_KEY_ID: AWS access key for SES
• AWS_SECRET_ACCESS_KEY: AWS secret key for SES
• AWS_REGION: AWS region for SES
• SES_FROM_EMAIL: Sender email address
• SES_FROM_NAME: Sender name (optional, defaults to "BriefBot")
• SES_REPLY_TO_EMAIL: Reply-to email address (optional, defaults to SES_FROM_EMAIL)

### Daily Digest Variables (Optional):
• DAILY_DIGEST_SUBJECT: Email subject template (optional, defaults to "Your Daily BriefBot Digest - %s")

## API Endpoints Created:

### User Endpoints:

- GET /users - List all users
- POST /users - Create user
- GET /users/:id - Get user by ID
- GET /users/email/:email - Get user by email
- PUT /users/:id - Update user
- DELETE /users/:id - Delete user

### Item Endpoints:

- POST /items - Create item
- GET /items/:id - Get item by ID
- GET /items/user/:userID - Get all items for a user
- GET /items/user/:userID/unread - Get unread items for a user
- PUT /items/:id - Update item
- PATCH /items/:id/read - Mark item as read
- DELETE /items/:id - Delete item

### Daily Digest Endpoints:

- POST /daily-digest/trigger - Send daily digest to all users with unread items from yesterday
- POST /daily-digest/trigger/user/:userID - Send daily digest to specific user with unread items from yesterday

## Quick Test Sequence:

Run these commands in order to test the full flow:

# 1. Create user

```
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name":"Test User","email":"test@example.com","auth_provider":"local","oauth_id":null,"password_hash":"hashedpassword"}'
```

# 2. Create item (use user ID from response above)

```
curl -X POST http://localhost:8080/items -H "Content-Type: application/json" -d '{"title":"Test Title","user_id":1,"url":"https://test.com","text_content":"Test content","summary":"Test summary","type":"article","platform":"web","tags":["test"],"authors":["Test Author"]}'
```

# 3. Get user

```
curl -X GET http://localhost:8080/users/1
```

# 4. Get items for user

```
curl -X GET http://localhost:8080/items/user/1
```

# 5. Mark item as read

```
curl -X PATCH http://localhost:8080/items/1/read
```

# 6. Test daily digest (requires email service configuration)

```
# Send daily digest to all users
curl -X POST http://localhost:8080/daily-digest/trigger

# Send daily digest to specific user (user ID 1)
curl -X POST http://localhost:8080/daily-digest/trigger/user/1
```

## Daily Digest Feature

The daily digest automatically compiles unread items from the previous day and sends them via email to users.

### Daily Digest Curl Examples:

# Send daily digest to all users (requires email service to be configured)

```
curl -X POST http://localhost:8080/daily-digest/trigger
```

# Send daily digest to specific user (e.g., user ID 1)

```
curl -X POST http://localhost:8080/daily-digest/trigger/user/1
```

### Daily Digest Configuration:

- `DAILY_DIGEST_SUBJECT`: Email subject template (optional, defaults to "Your Daily BriefBot Digest - %s")

### Daily Digest Behavior:

- Only sends to users with valid email addresses
- Only includes unread items from the previous day (24-48 hours ago)
- Only includes items with `processing_status = 'completed'`
- Sends individual emails to each user with their specific unread items
- Gracefully skips users with no unread items from the previous day
- Continues processing other users if one user fails
