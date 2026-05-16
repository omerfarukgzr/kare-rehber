<script setup lang="ts">
import type { Evaluation } from '@/api/evaluations'

defineProps<{ ev: Evaluation; showAdminNote?: boolean }>()

const motivationColor = (n: number | null) => {
  if (n == null) return ''
  if (n >= 4) return 'success'
  if (n >= 3) return 'warning'
  return 'danger'
}
</script>

<template>
  <el-descriptions :column="2" border>
    <el-descriptions-item label="Hafta">{{ ev.weekLabel }}</el-descriptions-item>
    <el-descriptions-item label="Durum">
      <el-tag v-if="ev.status === 'approved'" type="success">Onaylandı</el-tag>
      <el-tag v-else-if="ev.status === 'edited_by_admin'" type="warning">Yönetici Düzenledi</el-tag>
      <el-tag v-else type="info">Beklemede</el-tag>
    </el-descriptions-item>
    <el-descriptions-item label="Koç">{{ ev.coachName }}</el-descriptions-item>
    <el-descriptions-item label="Öğrenci">{{ ev.studentName }}</el-descriptions-item>
    <el-descriptions-item label="Ders Durumu">{{ ev.courseStatus || '—' }}</el-descriptions-item>
    <el-descriptions-item label="Ödev Yapıldı">
      <el-tag v-if="ev.homeworkDone === true" type="success">Evet</el-tag>
      <el-tag v-else-if="ev.homeworkDone === false" type="danger">Hayır</el-tag>
      <span v-else>—</span>
    </el-descriptions-item>
    <el-descriptions-item label="Motivasyon">
      <el-tag v-if="ev.motivation" :type="motivationColor(ev.motivation)">{{ ev.motivation }} / 5</el-tag>
      <span v-else>—</span>
    </el-descriptions-item>
    <el-descriptions-item label="Davranış">
      <el-tag v-if="ev.behavior" :type="motivationColor(ev.behavior)">{{ ev.behavior }} / 5</el-tag>
      <span v-else>—</span>
    </el-descriptions-item>
    <el-descriptions-item label="Genel Yorum" :span="2">
      <div class="note">{{ ev.generalNote || '—' }}</div>
    </el-descriptions-item>
    <el-descriptions-item v-if="showAdminNote && ev.adminOnlyNote" label="Yönetici Notu (sadece admin görür)" :span="2">
      <div class="note admin-note">{{ ev.adminOnlyNote }}</div>
    </el-descriptions-item>
  </el-descriptions>
</template>

<style scoped>
.note { white-space: pre-wrap; line-height: 1.5; }
.admin-note { background: #fef3c7; padding: 8px; border-radius: 6px; }
</style>
