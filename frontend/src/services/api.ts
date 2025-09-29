import axios from 'axios'
import type {
  User,
  Item,
  CreateUserRequest,
  UpdateUserRequest,
  CreateItemRequest,
  UpdateItemRequest,
  SubmitUrlRequest,
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