import axios from 'axios'
import type {
  User,
  Item,
  CreateUserRequest,
  UpdateUserRequest,
  CreateItemRequest,
  UpdateItemRequest,
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

  getItem: async (id: number): Promise<Item> => {
    const response = await api.get<Item>(`/items/${id}`)
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

  deleteItem: async (id: number): Promise<void> => {
    await api.delete(`/items/${id}`)
  },
}