import { api } from './client'

export interface EvaluationWeek {
  id: string
  weekNo: number
  label: string
  startDate: string
  endDate: string
  isOpen: boolean
  openedAt: string | null
  closedAt: string | null
}

export const weeksApi = {
  listAll() {
    return api.get<{ items: EvaluationWeek[] }>('/admin/weeks').then(r => r.data.items)
  },
  listOpen() {
    return api.get<{ items: EvaluationWeek[] }>('/weeks/open').then(r => r.data.items)
  },
  create(payload: { weekNo: number; label?: string; startDate: string; endDate: string }) {
    return api.post<EvaluationWeek>('/admin/weeks', payload).then(r => r.data)
  },
  update(id: string, payload: { label?: string; startDate?: string; endDate?: string }) {
    return api.patch<EvaluationWeek>(`/admin/weeks/${id}`, payload).then(r => r.data)
  },
  open(id: string) {
    return api.post<EvaluationWeek>(`/admin/weeks/${id}/open`).then(r => r.data)
  },
  close(id: string) {
    return api.post<EvaluationWeek>(`/admin/weeks/${id}/close`).then(r => r.data)
  },
  generate(payload: { startDate: string; weekCount: number }) {
    return api.post<{ items: EvaluationWeek[]; created: number }>('/admin/weeks/generate', payload).then(r => r.data)
  },
}
