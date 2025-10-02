import axios from 'axios'
import type {
  User,
  Item,
  CreateUserRequest,
  UpdateUserRequest,
  CreateItemRequest,
  UpdateItemRequest,
  SubmitUrlRequest,
  Podcast,
  CreatePodcastRequest,
  PodcastResponse,
} from '@/types'

const API_BASE_URL = 'http://localhost:8080'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// User API functions
export const userApi = {
  createUser: async (data: CreateUserRequest): Promise<User> => {
    const response = await api.post<User>('/users', data)
    return response.data
  },

  getUser: async (id: number): Promise<User> => {
    const response = await api.get<User>(`/users/${id}`)
    return response.data
  },

  getUserByEmail: async (email: string): Promise<User> => {
    const response = await api.get<User>(`/users/email/${email}`)
    return response.data
  },

  listUsers: async (): Promise<User[]> => {
    const response = await api.get<User[]>('/users')
    return response.data
  },

  updateUser: async (id: number, data: UpdateUserRequest): Promise<User> => {
    const response = await api.put<User>(`/users/${id}`, data)
    return response.data
  },

  deleteUser: async (id: number): Promise<void> => {
    await api.delete(`/users/${id}`)
  },
}

// Item API functions
export const itemApi = {
  createItem: async (data: CreateItemRequest): Promise<Item> => {
    const response = await api.post<Item>('/items', data)
    return response.data
  },

  submitUrl: async (data: SubmitUrlRequest): Promise<Item> => {
    const requestData: CreateItemRequest = {
      user_id: data.user_id,
      url: data.url,
      text_content: '', // Will be populated by backend scraping
    }
    const response = await api.post<Item>('/items', requestData)
    return response.data
  },

  getItem: async (id: number): Promise<Item> => {
    const response = await api.get<Item>(`/items/${id}`)
    return response.data
  },

  getItemById: async (id: number, userId: number): Promise<Item> => {
    const response = await api.get<Item>(`/items/${id}?user_id=${userId}`)
    return response.data
  },

  getItemsByUser: async (userId: number): Promise<Item[]> => {
    const response = await api.get<Item[]>(`/items/user/${userId}`)
    return response.data
  },

  getUnreadItemsByUser: async (userId: number): Promise<Item[]> => {
    const response = await api.get<Item[]>(`/items/user/${userId}/unread`)
    return response.data
  },

  updateItem: async (id: number, data: UpdateItemRequest): Promise<void> => {
    await api.put(`/items/${id}`, data)
  },

  patchItem: async (id: number, data: Partial<{ title: string; summary: string; tags: string[]; authors: string[] }>): Promise<Item> => {
    const response = await api.patch<Item>(`/items/${id}`, data)
    return response.data
  },

  markItemAsRead: async (id: number): Promise<Item> => {
    const response = await api.patch<Item>(`/items/${id}/read`)
    return response.data
  },

  toggleItemReadStatus: async (id: number): Promise<Item> => {
    const response = await api.patch<Item>(`/items/${id}/toggle-read`)
    return response.data
  },

  deleteItem: async (id: number): Promise<void> => {
    await api.delete(`/items/${id}`)
  },
}

// Digest API functions
export const digestApi = {
  triggerIntegratedDigest: async (): Promise<void> => {
    await api.post('/digest/trigger/integrated')
  },

  triggerIntegratedDigestForUser: async (userId: number): Promise<void> => {
    await api.post(`/digest/trigger/integrated/user/${userId}`)
  },

  triggerDailyDigest: async (): Promise<void> => {
    await api.post('/digest/trigger')
  },

  triggerDailyDigestForUser: async (userId: number): Promise<void> => {
    await api.post(`/digest/trigger/user/${userId}`)
  },
}

// Podcast API functions
export const podcastApi = {
  createPodcast: async (data: CreatePodcastRequest): Promise<PodcastResponse> => {
    const response = await api.post<PodcastResponse>('/podcasts', data)
    return response.data
  },

  getPodcast: async (id: number): Promise<Podcast> => {
    const response = await api.get<Podcast>(`/podcasts/${id}`)
    return response.data
  },

  getPodcastsByUser: async (userId: number): Promise<{ podcasts: Podcast[]; count: number }> => {
    const response = await api.get<{ podcasts: Podcast[]; count: number }>(`/podcasts/user/${userId}`)
    return response.data
  },

  getPodcastProcessingStatus: async (id: number): Promise<{
    podcast_id: number
    status: string
    is_pending: boolean
    is_writing: boolean
    is_generating: boolean
    is_processing: boolean
    is_completed: boolean
    is_failed: boolean
    audio_url?: string | null
  }> => {
    const response = await api.get(`/podcasts/${id}/status`)
    return response.data
  },

  deletePodcast: async (id: number): Promise<void> => {
    await api.delete(`/podcasts/${id}`)
  },
}