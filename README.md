# BriefBot

## Overview

An AI-enabled platform for managing links and extracting knowledge faster.

- Browser Plugin
- Telegram Bot Integration
- Web UI

## Tech Stack

- React (Tanstack)
- ShadcnUI
- Go
- Gin
- PostgreSQL
- Qdrant (Vector Database)
- Redis (Background jobs and pipeline management)
- Eleven Labs API (text-to-speech)
- Groq API (LLM)
- Cloudflare Browser Rendering
- Cloudflare R2 (Object Storage)

## Environment Variables:

• DATABASE_URL: PostgreSQL connection string
• PORT: Server port (default: 8080)

## API Endpoints Created:

### User Endpoints:

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

## Quick Test Sequence:

Run these commands in order to test the full flow:

# 1. Create user

curl -X POST http://localhost:8081/users -H "Content-Type: application/json" -d '{"name":"Test User","email":"test@example.com"}'

# 2. Create item (use user ID from response above)

curl -X POST http://localhost:8081/items -H "Content-Type: application/json" -d '{"user_id":1,"url":"https://test.com",
"text_content":"Test content"}'

# 3. Get user

curl -X GET http://localhost:8081/users/1

# 4. Get items for user

curl -X GET http://localhost:8081/items/user/1

# 5. Mark item as read

curl -X PATCH http://localhost:8081/items/1/read

# 6. Check unread items (should be empty now)

curl -X GET http://localhost:8081/items/user/1/unread
