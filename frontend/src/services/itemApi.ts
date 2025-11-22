import axios from 'axios'
import type {
  Item,
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