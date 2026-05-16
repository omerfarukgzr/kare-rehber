<script setup lang="ts">
import { onMounted, ref } from 'vue'
import PageHeader from '@/components/PageHeader.vue'
import { reportsApi, type SummaryReport } from '@/api/reports'

const summary = ref<SummaryReport | null>(null)
const loading = ref(false)

async function fetchData() {
  loading.value = true
  try {
    summary.value = await reportsApi.summary()
  } finally {
    loading.value = false
  }
}
onMounted(fetchData)
</script>

<template>
  <div>
    <PageHeader title="Yönetim Paneli" subtitle="Sistemin genel durumu ve hızlı erişimler" />
    <el-row v-loading="loading" :gutter="20">
      <el-col :span="6"><el-card><div class="stat"><div class="num">{{ summary?.activeStudents ?? '—' }}</div><div class="label">Aktif Öğrenci</div></div></el-card></el-col>
      <el-col :span="6"><el-card><div class="stat"><div class="num">{{ summary?.activeCoaches ?? '—' }}</div><div class="label">Aktif Koç</div></div></el-card></el-col>
      <el-col :span="6"><el-card><div class="stat"><div class="num">{{ summary?.activeAssignments ?? '—' }}</div><div class="label">Aktif Eşleştirme</div></div></el-card></el-col>
      <el-col :span="6"><el-card><div class="stat"><div class="num">{{ summary?.pendingRegistrations ?? '—' }}</div><div class="label">Bekleyen Başvuru</div></div></el-card></el-col>
    </el-row>
    <el-row v-loading="loading" :gutter="20" class="mt">
      <el-col :span="6"><el-card><div class="stat"><div class="num">{{ summary?.pendingEvaluations ?? '—' }}</div><div class="label">Bekleyen Değerlendirme</div></div></el-card></el-col>
      <el-col :span="6"><el-card><div class="stat"><div class="num">{{ summary?.approvedEvaluations ?? '—' }}</div><div class="label">Onaylı Değerlendirme</div></div></el-card></el-col>
      <el-col :span="6"><el-card><div class="stat"><div class="num">{{ summary?.openThreads ?? '—' }}</div><div class="label">Açık Mesaj</div></div></el-card></el-col>
      <el-col :span="6"><el-card><div class="stat"><div class="num">{{ summary?.totalEvaluations ?? '—' }}</div><div class="label">Toplam Değerlendirme</div></div></el-card></el-col>
    </el-row>
  </div>
</template>

<style scoped>
.stat { text-align: center; padding: 8px 0; }
.stat .num { font-size: 32px; font-weight: 700; color: #1f2937; }
.stat .label { color: #6b7280; font-size: 13px; margin-top: 4px; }
.mt { margin-top: 16px; }
</style>
