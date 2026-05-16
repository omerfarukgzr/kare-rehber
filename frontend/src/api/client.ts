import axios, { type AxiosInstance } from 'axios'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'

// Production'da Railway tarafından inject edilen VITE_API_URL kullanılır.
// Dev'de boş bırakılır → vite.config.ts proxy'si /api/* taleplerini backend'e yönlendirir.
const apiBase = (import.meta.env.VITE_API_URL ?? '').replace(/\/$/, '')

export const api: AxiosInstance = axios.create({
  baseURL: `${apiBase}/api/v1`,
  timeout: 20000,
})

api.interceptors.request.use((config) => {
  const auth = useAuthStore()
  if (auth.token) {
    config.headers = config.headers ?? {}
    config.headers.Authorization = `Bearer ${auth.token}`
  }
  return config
})

api.interceptors.response.use(
  (resp) => resp,
  (error) => {
    const status = error?.response?.status
    const apiErr = error?.response?.data?.error
    const message: string = apiErr?.message || error?.message || 'Bir hata oluştu'

    if (status === 401) {
      const auth = useAuthStore()
      auth.logout()
      if (typeof window !== 'undefined' && window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    } else if (status && status >= 400) {
      ElMessage.error(message)
    }

    return Promise.reject(error)
  },
)
