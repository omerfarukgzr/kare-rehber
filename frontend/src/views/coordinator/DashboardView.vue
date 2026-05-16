<script setup lang="ts">
import PageHeader from '@/components/PageHeader.vue'
import { useAuthStore } from '@/stores/auth'
const auth = useAuthStore()
</script>
<template>
  <div>
    <PageHeader :title="`Hoş geldiniz, ${auth.user?.fullName ?? 'Koordinatör'}`" subtitle="Sorumlu olduğunuz öğrencilerin durumu" />
    <el-card><el-empty description="Henüz veri yok" /></el-card>
  </div>
</template>
