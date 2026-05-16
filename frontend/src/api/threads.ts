import { api } from './client'
import type { Page } from '@/types/api'

export type ThreadKind = 'general' | 'feedback' | 'complaint'
export type ThreadStatus = 'open' | 'closed'

export interface Thread {
  id: string
  kind: ThreadKind
  subject: string
  openedByUserId: string
  openedByName: string
  openedByRole: string
  status: ThreadStatus
  unreadCount: number
  lastMessageAt: string | null
  lastMessageBody: string | null
  createdAt: string
}

export interface ThreadMessage {
  id: string
  threadId: string
  fromUserId: string
  fromName: string
  fromRole: string
  body: string
  readAt: string | null
  createdAt: string
}

export const threadsApi = {
  list(params: { kind?: ThreadKind; status?: ThreadStatus; search?: string; page?: number; pageSize?: number } = {}) {
    return api.get<Page<Thread>>('/threads', { params }).then(r => r.data)
  },
  create(payload: { kind: ThreadKind; subject: string; body: string }) {
    return api.post<Thread>('/threads', payload).then(r => r.data)
  },
  messages(id: string) {
    return api.get<{ items: ThreadMessage[] }>(`/threads/${id}/messages`).then(r => r.data.items)
  },
  send(id: string, body: string) {
    return api.post<ThreadMessage>(`/threads/${id}/messages`, { body }).then(r => r.data)
  },
  close(id: string) {
    return api.post(`/threads/${id}/close`).then(r => r.data)
  },
}
