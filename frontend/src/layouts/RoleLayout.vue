<script setup lang="ts">
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { SwitchButton } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

interface MenuItem { path: string; label: string }

const menuByRole: Record<string, MenuItem[]> = {
  coach: [
    { path: '/coach/dashboard', label: 'Panel' },
    { path: '/coach/students', label: 'Öğrencilerim' },
    { path: '/coach/evaluations/new', label: 'Yeni Değerlendirme' },
  ],
  coordinator: [
    { path: '/coordinator/dashboard', label: 'Panel' },
    { path: '/coordinator/students', label: 'Öğrencilerim' },
    { path: '/coordinator/evaluations', label: 'Değerlendirmeler' },
    { path: '/coordinator/threads', label: 'Mesajlar' },
  ],
  parent: [
    { path: '/parent/dashboard', label: 'Panel' },
    { path: '/parent/evaluations', label: 'Çocuğumun Değerlendirmeleri' },
    { path: '/parent/threads', label: 'Mesajlar' },
  ],
  student: [
    { path: '/student/dashboard', label: 'Panel' },
    { path: '/student/threads', label: 'Mesajlar' },
  ],
}

const items = computed<MenuItem[]>(() => menuByRole[auth.role ?? ''] ?? [])

const roleTitle = computed(() => {
  switch (auth.role) {
    case 'coach': return 'Koç Paneli'
    case 'coordinator': return 'Koordinatör Paneli'
    case 'parent': return 'Veli Paneli'
    case 'student': return 'Öğrenci Paneli'
    default: return 'Panel'
  }
})

function logout() {
  auth.logout()
  router.replace('/login')
}
</script>

<template>
  <el-container class="role-layout">
    <el-header class="role-header">
      <div class="brand">
        <div class="title">KARE-REHBER</div>
        <div class="sub">{{ roleTitle }}</div>
      </div>

      <el-menu :default-active="route.path" router mode="horizontal" class="top-menu">
        <el-menu-item v-for="m in items" :key="m.path" :index="m.path">
          {{ m.label }}
        </el-menu-item>
      </el-menu>

      <div class="header-right">
        <el-dropdown>
          <span class="user-trigger">
            <el-avatar :size="30">{{ auth.user?.fullName?.[0] ?? 'U' }}</el-avatar>
            <span class="user-name">{{ auth.user?.fullName }}</span>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="logout">
                <el-icon><SwitchButton /></el-icon>
                Çıkış Yap
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-header>

    <el-main class="role-main">
      <RouterView />
    </el-main>
  </el-container>
</template>

<style scoped>
.role-layout { height: 100vh; }
.role-header {
  display: flex;
  align-items: center;
  background: #fff;
  border-bottom: 1px solid #e5e7eb;
  padding: 0 24px;
  gap: 24px;
  height: 60px;
}
.brand .title { font-weight: 700; font-size: 16px; }
.brand .sub { font-size: 12px; color: #6b7280; }
.top-menu { flex: 1; border-bottom: none !important; }
.header-right { display: flex; align-items: center; }
.user-trigger { display: flex; align-items: center; gap: 8px; cursor: pointer; }
.user-name { font-weight: 500; color: #374151; }
.role-main { background: #f5f7fa; padding: 24px; }
</style>
