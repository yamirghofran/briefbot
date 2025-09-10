export interface User {
  id: number
  name?: string | null
  email?: string | null
  auth_provider?: string | null
  oauth_id?: string | null
  password_hash?: string | null
  created_at?: string | null
  updated_at?: string | null
}

export interface Item {
  id: number
  user_id?: number | null
  url?: string | null
  is_read?: boolean | null
  file_key?: string | null
  text_content?: string | null
  summary?: string | null
  created_at?: string | null
  modified_at?: string | null
}

export interface CreateUserRequest {
  name: string
  email: string
}

export interface UpdateUserRequest {
  name?: string
  email?: string
}

export interface CreateItemRequest {
  user_id: number
  url: string
  text_content: string
}

export interface UpdateItemRequest {
  user_id?: number
  url?: string
  text_content?: string
  summary?: string
  is_read?: boolean
}