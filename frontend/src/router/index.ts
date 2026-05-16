import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import type { Role } from '@/types/api'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/auth/LoginView.vue'),
    meta: { public: true },
  },
  {
    path: '/register/student',
    name: 'student-registration',
    component: () => import('@/views/auth/StudentRegistrationView.vue'),
    meta: { public: true },
  },
  {
    path: '/register/coach',
    name: 'coach-registration',
    component: () => import('@/views/auth/CoachRegistrationView.vue'),
    meta: { public: true },
  },

  {
    path: '/admin',
    component: () => import('@/layouts/AdminLayout.vue'),
    meta: { roles: ['admin'] satisfies Role[] },
    children: [
      { path: '', redirect: '/admin/dashboard' },
      { path: 'dashboard', name: 'admin-dashboard', component: () => import('@/views/admin/DashboardView.vue') },
      { path: 'users', name: 'admin-users', component: () => import('@/views/admin/UsersView.vue') },
      { path: 'approvals', name: 'admin-approvals', component: () => import('@/views/admin/ApprovalsView.vue') },
      { path: 'assignments', name: 'admin-assignments', component: () => import('@/views/admin/AssignmentsView.vue') },
      { path: 'weeks', name: 'admin-weeks', component: () => import('@/views/admin/WeeksView.vue') },
      { path: 'evaluations', name: 'admin-evaluations', component: () => import('@/views/admin/EvaluationsView.vue') },
      { path: 'threads', name: 'admin-threads', component: () => import('@/views/admin/ThreadsView.vue') },
      { path: 'sms', name: 'admin-sms', component: () => import('@/views/admin/SmsView.vue') },
      { path: 'reports', name: 'admin-reports', component: () => import('@/views/admin/ReportsView.vue') },
    ],
  },

  {
    path: '/coach',
    component: () => import('@/layouts/RoleLayout.vue'),
    meta: { roles: ['coach'] satisfies Role[] },
    children: [
      { path: '', redirect: '/coach/dashboard' },
      { path: 'dashboard', name: 'coach-dashboard', component: () => import('@/views/coach/DashboardView.vue') },
      { path: 'students', name: 'coach-students', component: () => import('@/views/coach/MyStudentsView.vue') },
      { path: 'evaluations/new', name: 'coach-evaluation-new', component: () => import('@/views/coach/EvaluationFormView.vue') },
    ],
  },

  {
    path: '/coordinator',
    component: () => import('@/layouts/RoleLayout.vue'),
    meta: { roles: ['coordinator'] satisfies Role[] },
    children: [
      { path: '', redirect: '/coordinator/dashboard' },
      { path: 'dashboard', name: 'coordinator-dashboard', component: () => import('@/views/coordinator/DashboardView.vue') },
      { path: 'students', name: 'coordinator-students', component: () => import('@/views/coordinator/MyStudentsView.vue') },
      { path: 'evaluations', name: 'coordinator-evaluations', component: () => import('@/views/coordinator/EvaluationsView.vue') },
      { path: 'threads', name: 'coordinator-threads', component: () => import('@/views/coordinator/ThreadsView.vue') },
    ],
  },

  {
    path: '/parent',
    component: () => import('@/layouts/RoleLayout.vue'),
    meta: { roles: ['parent'] satisfies Role[] },
    children: [
      { path: '', redirect: '/parent/dashboard' },
      { path: 'dashboard', name: 'parent-dashboard', component: () => import('@/views/parent/DashboardView.vue') },
      { path: 'evaluations', name: 'parent-evaluations', component: () => import('@/views/parent/ChildEvaluationsView.vue') },
      { path: 'threads', name: 'parent-threads', component: () => import('@/views/parent/ThreadsView.vue') },
    ],
  },

  {
    path: '/student',
    component: () => import('@/layouts/RoleLayout.vue'),
    meta: { roles: ['student'] satisfies Role[] },
    children: [
      { path: '', redirect: '/student/dashboard' },
      { path: 'dashboard', name: 'student-dashboard', component: () => import('@/views/student/DashboardView.vue') },
      { path: 'threads', name: 'student-threads', component: () => import('@/views/student/ThreadsView.vue') },
    ],
  },

  { path: '/', redirect: () => homeForRole(null) },
  { path: '/:pathMatch(.*)*', redirect: '/login' },
]

function homeForRole(role: Role | null): string {
  switch (role) {
    case 'admin':
      return '/admin/dashboard'
    case 'coach':
      return '/coach/dashboard'
    case 'coordinator':
      return '/coordinator/dashboard'
    case 'parent':
      return '/parent/dashboard'
    case 'student':
      return '/student/dashboard'
    default:
      return '/login'
  }
}

export const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  const isPublic = to.meta?.public === true

  if (isPublic) {
    if (auth.isAuthenticated && to.name === 'login') {
      return homeForRole(auth.role)
    }
    return true
  }

  if (!auth.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }

  const allowed = (to.meta?.roles as Role[] | undefined) ?? null
  if (allowed && auth.role && !allowed.includes(auth.role)) {
    return homeForRole(auth.role)
  }

  return true
})

export { homeForRole }
