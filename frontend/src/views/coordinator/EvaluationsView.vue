<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import PageHeader from '@/components/PageHeader.vue'
import EvaluationDisplay from '@/components/EvaluationDisplay.vue'
import { evaluationsApi, type Evaluation, type EvaluationStatus } from '@/api/evaluations'
import { weeksApi, type EvaluationWeek } from '@/api/weeks'
import type { Page } from '@/types/api'

const loading = ref(false)
const data = ref<Page<Evaluation>>({ items: [], totalCount: 0, page: 1, pageSize: 20 })
const weeks = ref<EvaluationWeek[]>([])

const filters = reactive<{ weekId?: string; page: number; pageSize: number }>({ page: 1, pageSize: 20 })

async function fetchData() {
  loading.value = true
  try {
    data.value = await evaluationsApi.list({ page: filters.page, pageSize: filters.pageSize, weekId: filters.weekId })
  } finally {
    loading.value = false
  }
}
onMounted(async () => {
  weeks.value = await weeksApi.listAll()
  fetchData()
})

const open = ref(false)
const current = ref<Evaluation | null>(null)
function openDetail(ev: Evaluation) { current.value = ev; open.value = true }

const statusType: Record<EvaluationStatus, '' | 'success' | 'warning' | 'info'> = {
  pending: 'info', edited_by_admin: 'warning', approved: 'success',
}
const statusLabel: Record<EvaluationStatus, string> = { pending: 'Beklemede', edited_by_admin: 'Düzenlendi', approved: 'Onaylandı' }
</script>

<template>
  <div>
    <PageHeader title="Değerlendirmeler" subtitle="Sorumlu öğrencilerinizin haftalık değerlendirmeleri (aktif sürüm)" />

    <el-card class="filters">
      <el-form :inline="true">
        <el-form-item label="Hafta">
          <el-select v-model="filters.weekId" clearable placeholder="Hepsi" style="width: 220px">
            <el-option v-for="w in weeks" :key="w.id" :label="w.label" :value="w.id" />
          </el-select>
        </el-form-item>
        <el-button type="primary" @click="() => { filters.page = 1; fetchData() }">Filtrele</el-button>
      </el-form>
    </el-card>

    <el-card>
      <el-table v-loading="loading" :data="data.items" stripe>
        <el-table-column label="Hafta" prop="weekLabel" min-width="180" />
        <el-table-column label="Öğrenci" prop="studentName" min-width="140" />
        <el-table-column label="Koç" prop="coachName" min-width="140" />
        <el-table-column label="Durum" width="140">
          <template #default="{ row }">
            <el-tag :type="statusType[row.status as EvaluationStatus]">{{ statusLabel[row.status as EvaluationStatus] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Motivasyon" width="110"><template #default="{ row }">{{ row.motivation ?? '—' }}</template></el-table-column>
        <el-table-column label="Davranış" width="100"><template #default="{ row }">{{ row.behavior ?? '—' }}</template></el-table-column>
        <el-table-column label="İşlem" width="100" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDetail(row)">Detay</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="pager"
        background layout="total, prev, pager, next"
        :total="data.totalCount"
        v-model:current-page="filters.page"
        v-model:page-size="filters.pageSize"
        @current-change="fetchData"
      />
    </el-card>

    <el-dialog v-model="open" title="Değerlendirme Detayı" width="700px">
      <EvaluationDisplay v-if="current" :ev="current" :show-admin-note="false" />
      <template #footer><el-button @click="open = false">Kapat</el-button></template>
    </el-dialog>
  </div>
</template>

<style scoped>
.filters { margin-bottom: 16px; }
.pager { margin-top: 16px; display: flex; justify-content: flex-end; }
</style>
