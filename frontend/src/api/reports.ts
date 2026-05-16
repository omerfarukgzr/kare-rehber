import { api } from './client'

export interface SummaryReport {
  activeStudents: number
  activeCoaches: number
  activeCoordinators: number
  activeAssignments: number
  totalEvaluations: number
  pendingEvaluations: number
  approvedEvaluations: number
  openThreads: number
  pendingRegistrations: number
}

export interface CoachPerformanceRow {
  coachId: string
  coachName: string
  studentCount: number
  evaluationCount: number
  approvedCount: number
  avgMotivation: number
  avgBehavior: number
}

export interface CityDistributionRow { city: string | null; count: number }

export interface WeekStatsRow {
  weekId: string
  weekNo: number
  label: string
  startDate: string
  totalEvaluations: number
  approvedCount: number
  activeAssignments: number
}

export interface StudentTrendRow {
  weekNo: number
  weekLabel: string
  motivation: number | null
  behavior: number | null
  status: string
}

export const reportsApi = {
  summary() {
    return api.get<SummaryReport>('/admin/reports/summary').then(r => r.data)
  },
  coachPerformance() {
    return api.get<{ items: CoachPerformanceRow[] }>('/admin/reports/coach-performance').then(r => r.data.items)
  },
  cityDistribution(role: string) {
    return api.get<{ items: CityDistributionRow[] }>('/admin/reports/city-distribution', { params: { role } }).then(r => r.data.items)
  },
  weekStats() {
    return api.get<{ items: WeekStatsRow[] }>('/admin/reports/week-stats').then(r => r.data.items)
  },
  studentTrend(studentId: string) {
    return api.get<{ items: StudentTrendRow[] }>('/reports/student-trend', { params: { studentId } }).then(r => r.data.items)
  },
}
