Database Schema

- Record

  - URL
  - Platform (e.g. Youtube, Github, Arxiv, Wallstreet Journal)
  - Type (e.g. Github Repo, Research Paper, Article, Youtube video)
  - Authors (e.g. Wallstreet Jourrnal, Greg Pauloski)
  - Tags (e.g. AI, Computer Science, Politics, News, Databases, Go)
  - isRead (boolean)
  - createdAt
  - modifiedAt
  - File URL (Cloudflare R2)

- Embedding Pipeline
- Text -> Chunking -> Gemini Embedding -> -> Store in Database -> Store in Qdrant
- File -> Cloudflare Markdown -> Gemini Embedding -> Store in Database -> Store in Qdrant

- Summarization Pipeline

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
