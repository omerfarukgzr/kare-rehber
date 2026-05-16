export type Role = 'admin' | 'coach' | 'student' | 'parent' | 'coordinator'

export interface User {
  id: string
  role: Role
  fullName: string
  phone: string
  email: string | null
  city: string | null
  isActive: boolean
  createdAt: string
}

export interface LoginRequest {
  identifier: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface ApiError {
  error: {
    code: string
    message: string
    details?: unknown
  }
}

export interface Page<T> {
  items: T[]
  totalCount: number
  page: number
  pageSize: number
}
