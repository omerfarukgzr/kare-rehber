<script setup lang="ts">
defineProps<{ title: string; subtitle?: string }>()
</script>

<template>
  <div class="page-header">
    <div>
      <h2>{{ title }}</h2>
      <p v-if="subtitle">{{ subtitle }}</p>
    </div>
    <div class="actions">
      <slot name="actions" />
    </div>
  </div>
</template>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
  gap: 12px;
}
.page-header h2 { margin: 0 0 4px; font-size: 22px; color: #1f2937; }
.page-header p { margin: 0; color: #6b7280; font-size: 14px; }
.actions { display: flex; gap: 8px; flex-wrap: wrap; }
</style>
