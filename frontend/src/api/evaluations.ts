import { api } from './client'
import type { Page, User } from '@/types/api'

export type EvaluationStatus = 'pending' | 'approved' | 'edited_by_admin'

export interface Evaluation {
  id: string
  assignmentId: string
  weekId: string
  weekNo: number
  weekLabel: string
  weekIsOpen: boolean
  coachId: string
  coachName: string
  studentId: string
  studentName: string
  studentCity: string | null
  courseStatus: string | null
  homeworkDone: boolean | null
  motivation: number | null
  behavior: number | null
  generalNote: string | null
  adminOnlyNote?: string | null
  status: EvaluationStatus
  submittedBy: string
  submittedAt: string
  lastEditedBy?: string | null
  lastEditedAt?: string | null
  approvedBy?: string | null
  approvedAt?: string | null
}

export interface EvaluationVersion {
  id: string
  versionNo: number
  snapshot: Record<string, unknown>
  editedBy: string | null
  editedAt: string
  changeReason: string | null
}

export interface CreateEvaluationPayload {
  studentId: string
  weekId: string
  courseStatus?: string | null
  homeworkDone?: boolean | null
  motivation?: number | null
  behavior?: number | null
  generalNote?: string | null
  adminOnlyNote?: string | null
}

export interface UpdateEvaluationPayload {
  courseStatus?: string | null
  homeworkDone?: boolean | null
  motivation?: number | null
  behavior?: number | null
  generalNote?: string | null
  adminOnlyNote?: string | null
  changeReason?: string | null
}

export const evaluationsApi = {
  list(params: { weekId?: string; coachId?: string; studentId?: string; status?: EvaluationStatus; search?: string; page?: number; pageSize?: number } = {}) {
    return api.get<Page<Evaluation>>('/evaluations', { params }).then(r => r.data)
  },
  get(id: string) {
    return api.get<Evaluation>(`/evaluations/${id}`).then(r => r.data)
  },
  coachCreate(payload: CreateEvaluationPayload) {
    return api.post<Evaluation>('/coach/evaluations', payload).then(r => r.data)
  },
  coachUpdate(id: string, payload: UpdateEvaluationPayload) {
    return api.patch<Evaluation>(`/coach/evaluations/${id}`, payload).then(r => r.data)
  },
  coachStudents() {
    return api.get<{ items: User[] }>('/coach/students').then(r => r.data.items)
  },
  adminUpdate(id: string, payload: UpdateEvaluationPayload) {
    return api.patch<Evaluation>(`/admin/evaluations/${id}`, payload).then(r => r.data)
  },
  approve(id: string) {
    return api.post<Evaluation>(`/admin/evaluations/${id}/approve`).then(r => r.data)
  },
  revoke(id: string) {
    return api.post<Evaluation>(`/admin/evaluations/${id}/revoke`).then(r => r.data)
  },
  versions(id: string) {
    return api.get<{ items: EvaluationVersion[] }>(`/admin/evaluations/${id}/versions`).then(r => r.data.items)
  },
}
