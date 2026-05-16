<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import {
  HomeFilled,
  UserFilled,
  CircleCheck,
  Connection,
  Calendar,
  Document,
  ChatLineRound,
  Message,
  DataAnalysis,
  SwitchButton,
} from '@element-plus/icons-vue'

const router = useRouter()
const auth = useAuthStore()

const menuItems = computed(() => [
  { path: '/admin/dashboard', label: 'Panel', icon: HomeFilled },
  { path: '/admin/users', label: 'Kullanıcılar', icon: UserFilled },
  { path: '/admin/approvals', label: 'Onaylar', icon: CircleCheck },
  { path: '/admin/assignments', label: 'Eşleştirmeler', icon: Connection },
  { path: '/admin/weeks', label: 'Hafta Yönetimi', icon: Calendar },
  { path: '/admin/evaluations', label: 'Değerlendirmeler', icon: Document },
  { path: '/admin/threads', label: 'Mesajlar', icon: ChatLineRound },
  { path: '/admin/sms', label: 'SMS', icon: Message },
  { path: '/admin/reports', label: 'Raporlar', icon: DataAnalysis },
])

function logout() {
  auth.logout()
  router.replace('/login')
}
</script>

<template>
  <el-container class="admin-layout">
    <el-aside width="240px" class="admin-aside">
      <div class="brand">
        <div class="brand-title">KARE-REHBER</div>
        <div class="brand-sub">Yönetim Paneli</div>
      </div>
      <el-menu
        :default-active="$route.path"
        router
        class="admin-menu"
        background-color="#1f2937"
        text-color="#e5e7eb"
        active-text-color="#fbbf24"
      >
        <el-menu-item v-for="m in menuItems" :key="m.path" :index="m.path">
          <el-icon><component :is="m.icon" /></el-icon>
          <span>{{ m.label }}</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="admin-header">
        <div class="header-spacer"></div>
        <el-dropdown>
          <span class="user-trigger">
            <el-avatar :size="32">{{ auth.user?.fullName?.[0] ?? 'A' }}</el-avatar>
            <span class="user-name">{{ auth.user?.fullName ?? 'Yönetici' }}</span>
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
      </el-header>

      <el-main class="admin-main">
        <RouterView />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.admin-layout {
  height: 100vh;
}
.admin-aside {
  background: #1f2937;
  display: flex;
  flex-direction: column;
}
.brand {
  padding: 20px 16px;
  color: #fff;
  border-bottom: 1px solid #374151;
}
.brand-title {
  font-size: 18px;
  font-weight: 700;
  letter-spacing: 0.5px;
}
.brand-sub {
  font-size: 12px;
  color: #9ca3af;
  margin-top: 2px;
}
.admin-menu {
  flex: 1;
  border-right: none !important;
}
.admin-header {
  background: #fff;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  align-items: center;
  padding: 0 24px;
}
.header-spacer { flex: 1; }
.user-trigger {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}
.user-name { color: #374151; font-weight: 500; }
.admin-main {
  background: #f5f7fa;
  padding: 24px;
}
</style>
