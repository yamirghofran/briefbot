import axios from 'axios'
import type {
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