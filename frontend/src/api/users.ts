import { api } from './client'
import type { Page, Role, User } from '@/types/api'

export interface CreateUserPayload {
  role: Role
  fullName: string
  phone: string
  email?: string
  city?: string
  password?: string
}

export interface UpdateUserPayload {
  fullName?: string
  email?: string | null
  city?: string | null
  isActive?: boolean
}

export interface CreateUserResponse {
  user: User
  generatedPassword?: string
}

export interface ListUsersParams {
  page?: number
  pageSize?: number
  role?: Role
  isActive?: 'true' | 'false'
  city?: string
  search?: string
}

export const usersApi = {
  list(params: ListUsersParams = {}) {
    return api.get<Page<User>>('/admin/users', { params }).then(r => r.data)
  },
  get(id: string) {
    return api.get<User>(`/admin/users/${id}`).then(r => r.data)
  },
  create(payload: CreateUserPayload) {
    return api.post<CreateUserResponse>('/admin/users', payload).then(r => r.data)
  },
  update(id: string, payload: UpdateUserPayload) {
    return api.patch<User>(`/admin/users/${id}`, payload).then(r => r.data)
  },
  resetPassword(id: string) {
    return api.post<{ generatedPassword: string }>(`/admin/users/${id}/reset-password`).then(r => r.data)
  },
}
