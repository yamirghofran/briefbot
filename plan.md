Database Schema

- Record

  - URL
  - Platform (e.g. Youtube, Github, Arxiv, Wallstreet Journal)
  - Type (e.g. Github Repo, Research Paper, Article, Youtube video)
  - Authors (e.g. Wallstreet Jourrnal, Greg Pauloski)
  - Tags (e.g. AI, Computer Science, Politics, News, Databases, Go)
  - isRead (boolean)
  - File URL (Cloudflare R2)
  - Text content
  - Summary
  - createdAt
  - modifiedAt

- Create Record pipeline

1. Identify item type
2. Fetch text version of content (curl/web scrabing -> Cloudflare Broswer Rendering if doesn't work)
3. Create summary
4. Create/assign author, platform, etc. (fetch authors/platforms -> select the ones that apply -> create & select ones that don't exist.)
5. Create/assign tags (fetch tags -> select the ones that apply -> create & select ones that don't exist.)
6. Index (Generate & store embeddings)

- Embedding Pipeline
- Text -> Chunking -> Gemini Embedding -> -> Store in Database -> Store in Qdrant
- File -> Cloudflare Markdown -> Gemini Embedding -> Store in Database -> Store in Qdrant

- Daily Summarization Pipeline

1. Find links added during yesterday
2. Get the text content (e.g. Youtube Transcript, PDF text, Github Repo)
3. Find key points from each item
4. Create summary

- Podcast Creation

1. Find links added during yesterday
2. Get the text content (e.g. Youtube Transcript, PDF text, Github Repo)
3. Find key points from each item
4. Create a dialogue between 2 people with Kimi
5. Text-to-speech with Eleven Labs
6. Send to user via Telegram

- Entry points for adding links
  - Browser plugin
  - Telegram bot
  - Web UI

## Pipelines

- Record creation
  - Embedding pipeline
- Daily summary
- Daily podcast
  Organize and keep track of pipelines with Redis.

## Task List

- [ ] Ugly Frontend (Basic CRUD functionality)
- [ ] Filters (date, tags, author, platform, type, etc.)
- [ ] AI-assisted filter creator
- [ ] Cloudflare Browser Rendering to fetch page content.
- [ ] Create record by sending link/file to Telegram Bot
