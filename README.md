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

## API Endpoints Created:

### User Endpoints:

• POST /users - Create user
• GET /users/:id - Get user by ID
• GET /users/email/:email - Get user by email
• PUT /users/:id - Update user
• DELETE /users/:id - Delete user

### Item Endpoints:

• POST /items - Create item
• GET /items/:id - Get item by ID
• GET /items/user/:userID - Get all items for a user
• GET /items/user/:userID/unread - Get unread items for a user
• PUT /items/:id - Update item
• PATCH /items/:id/read - Mark item as read
• DELETE /items/:id - Delete item
