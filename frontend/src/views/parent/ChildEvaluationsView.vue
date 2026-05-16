<script setup lang="ts">
import { onMounted, ref } from 'vue'
import PageHeader from '@/components/PageHeader.vue'
import EvaluationDisplay from '@/components/EvaluationDisplay.vue'
import { evaluationsApi, type Evaluation } from '@/api/evaluations'

const loading = ref(false)
const evaluations = ref<Evaluation[]>([])

async function fetchData() {
  loading.value = true
  try {
    const r = await evaluationsApi.list({ pageSize: 200 })
    evaluations.value = r.items
  } finally {
    loading.value = false
  }
}
onMounted(fetchData)

const detail = ref<Evaluation | null>(null)
const open = ref(false)
function openDetail(ev: Evaluation) { detail.value = ev; open.value = true }
</script>

<template>
  <div>
    <PageHeader title="Çocuğumun Değerlendirmeleri" subtitle="Yönetici onayından geçen haftalık raporlar" />

    <el-card v-loading="loading">
      <el-empty v-if="!loading && evaluations.length === 0" description="Henüz onaylanmış bir değerlendirme yok" />
      <el-timeline v-else>
        <el-timeline-item v-for="ev in evaluations" :key="ev.id" type="success" :timestamp="`${ev.weekLabel} • ${ev.coachName}`" placement="top">
          <el-card class="card" @click="openDetail(ev)">
            <p v-if="ev.generalNote" class="snippet">{{ ev.generalNote }}</p>
            <p v-else class="dim">Detayları görmek için tıklayın</p>
            <div class="ratings">
              <span v-if="ev.motivation !== null">Motivasyon: <strong>{{ ev.motivation }}/5</strong></span>
              <span v-if="ev.behavior !== null">Davranış: <strong>{{ ev.behavior }}/5</strong></span>
              <span v-if="ev.homeworkDone !== null">Ödev: <strong>{{ ev.homeworkDone ? 'Evet' : 'Hayır' }}</strong></span>
            </div>
          </el-card>
        </el-timeline-item>
      </el-timeline>
    </el-card>

    <el-dialog v-model="open" title="Haftalık Değerlendirme" width="640px">
      <EvaluationDisplay v-if="detail" :ev="detail" :show-admin-note="false" />
      <template #footer><el-button @click="open = false">Kapat</el-button></template>
    </el-dialog>
  </div>
</template>

<style scoped>
.card { cursor: pointer; }
.snippet { margin: 0 0 8px; line-height: 1.5; }
.dim { color: #6b7280; margin: 0 0 8px; }
.ratings { display: flex; gap: 16px; color: #4b5563; font-size: 13px; }
</style>
