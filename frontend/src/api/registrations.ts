import { api } from './client'
import type { Page } from '@/types/api'

export type RegistrationKind = 'student' | 'coach'
export type RegistrationStatus = 'pending' | 'approved' | 'rejected'

export interface Registration {
  id: string
  kind: RegistrationKind
  status: RegistrationStatus
  fullName: string
  phone: string
  email: string | null
  city: string | null
  payload: Record<string, unknown>
  reviewedAt: string | null
  reviewNote: string | null
  createdAt: string
}

export interface DecisionResponse {
  registration: Registration
  userId?: string
  generatedPassword?: string
  smsLogId?: string
}

export interface StudentRegistrationPayload {
  fullName: string
  phone: string
  email?: string
  city?: string
  school?: string
  grade?: string
  parentFullName?: string
  parentPhone?: string
  notes?: string
}

export interface CoachRegistrationPayload {
  fullName: string
  phone: string
  email?: string
  city?: string
  bio?: string
  experience?: string
}

export const registrationsApi = {
  applyStudent(payload: StudentRegistrationPayload) {
    return api.post<Registration>('/registrations/student', payload).then(r => r.data)
  },
  applyCoach(payload: CoachRegistrationPayload) {
    return api.post<Registration>('/registrations/coach', payload).then(r => r.data)
  },
  list(params: { kind?: RegistrationKind; status?: RegistrationStatus; search?: string; page?: number; pageSize?: number } = {}) {
    return api.get<Page<Registration>>('/admin/registrations', { params }).then(r => r.data)
  },
  decide(id: string, approve: boolean, note: string) {
    return api.post<DecisionResponse>(`/admin/registrations/${id}/decision`, { approve, note }).then(r => r.data)
  },
}
