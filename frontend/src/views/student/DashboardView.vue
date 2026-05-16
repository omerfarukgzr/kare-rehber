<script setup lang="ts">
import PageHeader from '@/components/PageHeader.vue'
import { useAuthStore } from '@/stores/auth'
const auth = useAuthStore()
</script>
<template>
  <div>
    <PageHeader :title="`Hoş geldin, ${auth.user?.fullName ?? 'Öğrenci'}`" subtitle="Mesajlaş, geri bildirim gönder" />
    <el-card><el-empty description="Mesajlar bölümünden yönetimle iletişime geçebilirsiniz" /></el-card>
  </div>
</template>
