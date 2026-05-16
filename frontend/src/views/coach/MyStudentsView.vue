<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import PageHeader from '@/components/PageHeader.vue'
import { evaluationsApi, type Evaluation } from '@/api/evaluations'
import { weeksApi, type EvaluationWeek } from '@/api/weeks'
import type { User } from '@/types/api'

const router = useRouter()

const loading = ref(false)
const students = ref<User[]>([])
const openWeeks = ref<EvaluationWeek[]>([])
const evaluationsByStudent = ref<Record<string, Evaluation[]>>({})

async function loadAll() {
  loading.value = true
  try {
    const [s, w] = await Promise.all([
      evaluationsApi.coachStudents(),
      weeksApi.listOpen(),
    ])
    students.value = s
    openWeeks.value = w

    const evals = await evaluationsApi.list({ pageSize: 200 })
    const grouped: Record<string, Evaluation[]> = {}
    for (const ev of evals.items) {
      if (!grouped[ev.studentId]) grouped[ev.studentId] = []
      grouped[ev.studentId].push(ev)
    }
    evaluationsByStudent.value = grouped
  } finally {
    loading.value = false
  }
}
onMounted(loadAll)

function statusOfStudent(studentId: string): { label: string; type: 'success' | 'warning' | 'info' | 'danger' } {
  if (!openWeeks.value.length) return { label: 'Açık hafta yok', type: 'info' }
  const activeWeek = openWeeks.value[openWeeks.value.length - 1]
  const evs = evaluationsByStudent.value[studentId] ?? []
  const hasActive = evs.some(e => e.weekId === activeWeek.id)
  if (hasActive) return { label: 'Bu hafta girildi', type: 'success' }
  return { label: 'Bu hafta eksik', type: 'warning' }
}

function startEvaluation(studentId: string) {
  router.push({ name: 'coach-evaluation-new', query: { studentId } })
}
</script>

<template>
  <div>
    <PageHeader title="Öğrencilerim" subtitle="Atanmış öğrencileriniz ve bu haftaki durumları">
      <template #actions>
        <el-button type="primary" @click="router.push('/coach/evaluations/new')">Yeni Değerlendirme</el-button>
      </template>
    </PageHeader>

    <el-card v-loading="loading">
      <el-empty v-if="!loading && students.length === 0" description="Henüz size atanmış bir öğrenci yok" />
      <el-table v-else :data="students" stripe>
        <el-table-column label="Ad Soyad" prop="fullName" />
        <el-table-column label="Telefon" prop="phone" width="160" />
        <el-table-column label="Şehir" prop="city" width="140" />
        <el-table-column label="Bu Hafta" width="180">
          <template #default="{ row }">
            <el-tag :type="statusOfStudent(row.id).type">{{ statusOfStudent(row.id).label }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="İşlem" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="startEvaluation(row.id)">Değerlendir</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>
