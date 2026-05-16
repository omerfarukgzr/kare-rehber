<script setup lang="ts">
import PageHeader from './PageHeader.vue'

defineProps<{ title: string; subtitle?: string; phase?: string }>()
</script>

<template>
  <div>
    <PageHeader :title="title" :subtitle="subtitle" />
    <el-card>
      <div class="placeholder">
        <el-empty :description="phase ? `Bu modül ${phase}'de tamamlanacak.` : 'Bu modül henüz aktif değil.'" />
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.placeholder { padding: 24px 0; }
</style>
