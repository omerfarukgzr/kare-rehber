import { api } from './client'
import type { Page } from '@/types/api'

export interface SmsTemplate {
  key: string
  name: string
  body: string
}

export interface SmsLog {
  id: string
  toPhone: string
  toUserId: string | null
  body: string
  sentByUserId: string | null
  status: string
  providerName: string
  templateKey: string | null
  sentAt: string
}

export interface MissingCoachRow {
  coachId: string
  coachName: string
  coachPhone: string
  studentId: string
  studentName: string
}

export const smsApi = {
  templates() {
    return api.get<{ items: SmsTemplate[] }>('/admin/sms/templates').then(r => r.data.items)
  },
  send(payload: { userIds?: string[]; phones?: string[]; body: string; templateKey?: string }) {
    return api.post<{ sent: number; failed: number; errors?: string[] }>('/admin/sms/send', payload).then(r => r.data)
  },
  logs(params: { search?: string; page?: number; pageSize?: number } = {}) {
    return api.get<Page<SmsLog>>('/admin/sms/logs', { params }).then(r => r.data)
  },
  missingCoaches(weekId: string) {
    return api.get<{ items: MissingCoachRow[] }>('/admin/sms/missing-coaches', { params: { weekId } }).then(r => r.data.items)
  },
}
