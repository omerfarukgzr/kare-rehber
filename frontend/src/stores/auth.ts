import { defineStore } from 'pinia'
import type { Role, User } from '@/types/api'

interface AuthState {
  token: string | null
  user: User | null
}

const STORAGE_KEY = 'kare_auth'

function loadFromStorage(): AuthState {
  if (typeof window === 'undefined') return { token: null, user: null }
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return { token: null, user: null }
    return JSON.parse(raw) as AuthState
  } catch {
    return { token: null, user: null }
  }
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => loadFromStorage(),
  getters: {
    isAuthenticated: (s) => Boolean(s.token && s.user),
    role: (s): Role | null => s.user?.role ?? null,
  },
  actions: {
    setSession(token: string, user: User) {
      this.token = token
      this.user = user
      localStorage.setItem(STORAGE_KEY, JSON.stringify({ token, user }))
    },
    logout() {
      this.token = null
      this.user = null
      localStorage.removeItem(STORAGE_KEY)
    },
  },
})
