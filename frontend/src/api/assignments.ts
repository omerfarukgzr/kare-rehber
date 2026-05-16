import { api } from './client'
import type { Page, User } from '@/types/api'

export interface Assignment {
  id: string
  coachId: string
  coachName: string
  coachPhone: string
  studentId: string
  studentName: string
  studentPhone: string
  studentCity: string | null
  startedAt: string
  endedAt: string | null
  isActive: boolean
}

export const assignmentsApi = {
  list(params: { coachId?: string; studentId?: string; city?: string; search?: string; page?: number; pageSize?: number } = {}) {
    return api.get<Page<Assignment>>('/admin/assignments', { params }).then(r => r.data)
  },
  assign(coachId: string, studentIds: string[]) {
    return api.post<{ created: number; failed: number; errors?: string[] }>('/admin/assignments', { coachId, studentIds }).then(r => r.data)
  },
  unassign(studentId: string) {
    return api.post('/admin/assignments/unassign', { studentId }).then(r => r.data)
  },
  studentsWithoutCoach(params: { city?: string; search?: string } = {}) {
    return api.get<{ items: User[] }>('/admin/assignments/students-without-coach', { params }).then(r => r.data)
  },
  setCoordinator(studentIds: string[], coordinatorId: string | null) {
    return api.post<{ updated: number }>('/admin/assignments/set-coordinator', { studentIds, coordinatorId }).then(r => r.data)
  },
  setParent(studentId: string, parentId: string | null) {
    return api.post('/admin/assignments/set-parent', { studentId, parentId }).then(r => r.data)
  },
}
