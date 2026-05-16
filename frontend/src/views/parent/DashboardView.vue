<script setup lang="ts">
import PageHeader from '@/components/PageHeader.vue'
import { useAuthStore } from '@/stores/auth'
const auth = useAuthStore()
</script>
<template>
  <div>
    <PageHeader :title="`Hoş geldiniz, ${auth.user?.fullName ?? 'Veli'}`" subtitle="Çocuğunuzun gelişim takibi" />
    <el-card><el-empty description="Henüz onaylanmış değerlendirme yok" /></el-card>
  </div>
</template>
