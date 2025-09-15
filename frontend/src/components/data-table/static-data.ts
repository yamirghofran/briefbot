export type ItemType = "article" | "video" | "podcast" | "document" | "other"

export interface Item {
  id: string
  user_id: number
  url: string
  is_read: boolean
  text_content: string
  summary: string
  title: string
  type: ItemType
  tags: string[]
  platform: string
  authors: string[]
  created_at: Date
  modified_at: Date
}

export const itemTypes: Array<{
  value: ItemType
  label: string
  icon: string
}> = [
  {
    value: "article",
    label: "Article",
    icon: "📄",
  },
  {
    value: "video",
    label: "Video",
    icon: "🎥",
  },
  {
    value: "podcast",
    label: "Podcast",
    icon: "🎙️",
  },
  {
    value: "document",
    label: "Document",
    icon: "📑",
  },
  {
    value: "other",
    label: "Other",
    icon: "📁",
  },
]

export const platforms = [
  "Medium",
  "YouTube",
  "Spotify",
  "GitHub",
  "Substack",
  "LinkedIn",
  "Twitter",
  "ArXiv",
  "Other",
]

// Sample data
export const sampleItems: Item[] = [
  {
    id: "item-1",
    user_id: 1,
    url: "https://medium.com/@author/building-modern-web-apps",
    is_read: false,
    text_content: "A comprehensive guide to building modern web applications with React and TypeScript...",
    summary: "This article covers the fundamentals of modern web development using React and TypeScript, including best practices for component architecture and state management.",
    title: "Building Modern Web Applications",
    type: "article",
    tags: ["react", "typescript", "web-development", "frontend"],
    platform: "Medium",
    authors: ["Sarah Johnson"],
    created_at: new Date("2024-01-15T10:30:00Z"),
    modified_at: new Date("2024-01-15T10:30:00Z"),
  },
  {
    id: "item-2",
    user_id: 1,
    url: "https://youtube.com/watch?v=ai-intro-machine-learning",
    is_read: true,
    text_content: "An introduction to machine learning concepts and practical applications in real-world scenarios...",
    summary: "This video tutorial provides a beginner-friendly introduction to machine learning, covering basic algorithms and their practical applications.",
    title: "Introduction to Machine Learning",
    type: "video",
    tags: ["machine-learning", "ai", "tutorial"],
    platform: "YouTube",
    authors: ["Dr. Michael Chen"],
    created_at: new Date("2024-01-12T14:15:00Z"),
    modified_at: new Date("2024-01-12T14:15:00Z"),
  },
  {
    id: "item-3",
    user_id: 1,
    url: "https://spotify.com/episode/startup-founder-interview",
    is_read: false,
    text_content: "An in-depth interview with successful startup founders discussing their journey and lessons learned...",
    summary: "This podcast episode features interviews with three successful startup founders who share their experiences, challenges, and key insights from their entrepreneurial journey.",
    title: "Startup Founder Interview",
    type: "podcast",
    tags: ["startup", "entrepreneurship", "interview"],
    platform: "Spotify",
    authors: ["Alex Rivera", "Maria Garcia"],
    created_at: new Date("2024-01-18T09:45:00Z"),
    modified_at: new Date("2024-01-18T09:45:00Z"),
  },
  {
    id: "item-4",
    user_id: 1,
    url: "https://github.com/user/advanced-react-patterns",
    is_read: true,
    text_content: "A detailed documentation of advanced React patterns and best practices for building scalable applications...",
    summary: "This GitHub repository contains comprehensive documentation on advanced React patterns including hooks, context, and performance optimization techniques.",
    title: "Advanced React Patterns",
    type: "document",
    tags: ["react", "patterns", "documentation", "advanced"],
    platform: "GitHub",
    authors: ["David Kim"],
    created_at: new Date("2024-01-16T16:20:00Z"),
    modified_at: new Date("2024-01-16T16:20:00Z"),
  },
  {
    id: "item-5",
    user_id: 1,
    url: "https://linkedin.com/posts/ai-future-technology",
    is_read: false,
    text_content: "Exploring the future of AI technology and its impact on various industries...",
    summary: "This LinkedIn post discusses emerging AI technologies and their potential impact on different industries, including healthcare, finance, and education.",
    title: "Future of AI Technology",
    type: "article",
    tags: ["ai", "future", "technology", "industry"],
    platform: "LinkedIn",
    authors: ["Jennifer Wu"],
    created_at: new Date("2024-01-14T11:30:00Z"),
    modified_at: new Date("2024-01-14T11:30:00Z"),
  },
  {
    id: "item-6",
    user_id: 1,
    url: "https://arxiv.org/paper/quantum-computing-applications",
    is_read: true,
    text_content: "A research paper on quantum computing applications in cryptography and optimization problems...",
    summary: "This academic paper explores the current state and future applications of quantum computing, particularly focusing on its potential in cryptography and solving complex optimization problems.",
    title: "Quantum Computing Applications",
    type: "document",
    tags: ["quantum-computing", "research", "cryptography", "academic"],
    platform: "ArXiv",
    authors: ["Dr. Robert Lee", "Dr. Emily Zhang"],
    created_at: new Date("2024-01-20T08:00:00Z"),
    modified_at: new Date("2024-01-20T08:00:00Z"),
  },
  {
    id: "item-7",
    user_id: 1,
    url: "https://substack.com/newsletter/product-management-tips",
    is_read: false,
    text_content: "Weekly newsletter with practical product management tips and strategies...",
    summary: "This newsletter issue provides actionable product management advice, including user research techniques, roadmap planning, and stakeholder communication strategies.",
    title: "Product Management Tips",
    type: "article",
    tags: ["product-management", "newsletter", "tips"],
    platform: "Substack",
    authors: ["Tom Anderson"],
    created_at: new Date("2024-01-17T07:15:00Z"),
    modified_at: new Date("2024-01-17T07:15:00Z"),
  },
  {
    id: "item-8",
    user_id: 1,
    url: "https://twitter.com/tech-lead-thread/scaling-systems",
    is_read: true,
    text_content: "A comprehensive Twitter thread on scaling systems and architecture decisions...",
    summary: "This Twitter thread from an experienced tech lead covers real-world lessons learned from scaling large systems, including database sharding, caching strategies, and microservices architecture.",
    title: "Scaling Systems Architecture",
    type: "article",
    tags: ["scaling", "architecture", "twitter-thread", "systems"],
    platform: "Twitter",
    authors: ["Sarah Patel"],
    created_at: new Date("2024-01-19T13:45:00Z"),
    modified_at: new Date("2024-01-19T13:45:00Z"),
  },
]