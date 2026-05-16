<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/PageHeader.vue'
import { evaluationsApi, type CreateEvaluationPayload, type Evaluation } from '@/api/evaluations'
import { weeksApi, type EvaluationWeek } from '@/api/weeks'
import type { User } from '@/types/api'

const router = useRouter()

const students = ref<User[]>([])
const weeks = ref<EvaluationWeek[]>([])
const loading = ref(false)
const submitting = ref(false)

const form = ref<CreateEvaluationPayload>({
  studentId: '',
  weekId: '',
  courseStatus: '',
  homeworkDone: null,
  motivation: null,
  behavior: null,
  generalNote: '',
  adminOnlyNote: '',
})

const existing = ref<Evaluation | null>(null)

async function loadInit() {
  loading.value = true
  try {
    const [s, w] = await Promise.all([
      evaluationsApi.coachStudents(),
      weeksApi.listOpen(),
    ])
    students.value = s
    weeks.value = w
    if (w.length) form.value.weekId = w[0].id
  } finally {
    loading.value = false
  }
}
onMounted(loadInit)

async function checkExisting() {
  existing.value = null
  if (!form.value.studentId || !form.value.weekId) return
  const data = await evaluationsApi.list({ studentId: form.value.studentId, weekId: form.value.weekId, pageSize: 1 })
  if (data.items.length) {
    existing.value = data.items[0]
    form.value.courseStatus = existing.value.courseStatus ?? ''
    form.value.homeworkDone = existing.value.homeworkDone
    form.value.motivation = existing.value.motivation
    form.value.behavior = existing.value.behavior
    form.value.generalNote = existing.value.generalNote ?? ''
    form.value.adminOnlyNote = existing.value.adminOnlyNote ?? ''
  } else {
    form.value.courseStatus = ''
    form.value.homeworkDone = null
    form.value.motivation = null
    form.value.behavior = null
    form.value.generalNote = ''
    form.value.adminOnlyNote = ''
  }
}
watch(() => [form.value.studentId, form.value.weekId], checkExisting)

async function submit() {
  if (!form.value.studentId || !form.value.weekId) {
    ElMessage.warning('Öğrenci ve hafta seçin')
    return
  }
  submitting.value = true
  try {
    if (existing.value) {
      await evaluationsApi.coachUpdate(existing.value.id, {
        courseStatus: form.value.courseStatus,
        homeworkDone: form.value.homeworkDone,
        motivation: form.value.motivation,
        behavior: form.value.behavior,
        generalNote: form.value.generalNote,
        adminOnlyNote: form.value.adminOnlyNote,
      })
      ElMessage.success('Değerlendirme güncellendi')
    } else {
      await evaluationsApi.coachCreate(form.value)
      ElMessage.success('Değerlendirme kaydedildi (yönetici onayına gönderildi)')
    }
    router.push('/coach/students')
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div>
    <PageHeader title="Haftalık Değerlendirme" subtitle="Öğrenciniz için haftalık takip notu girin" />

    <el-card v-loading="loading">
      <el-form label-position="top" :model="form">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Öğrenci" required>
              <el-select v-model="form.studentId" placeholder="Öğrenci seç" filterable style="width: 100%">
                <el-option v-for="s in students" :key="s.id" :label="s.fullName" :value="s.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Hafta" required>
              <el-select v-model="form.weekId" placeholder="Hafta seç" style="width: 100%">
                <el-option v-for="w in weeks" :key="w.id" :label="w.label" :value="w.id" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-alert v-if="existing && existing.status !== 'pending'" type="warning" show-icon :closable="false" class="mt">
          Bu değerlendirme yönetici tarafından düzenlenmiş veya onaylanmıştır. Koç tarafından düzenlenemez.
        </el-alert>
        <el-alert v-else-if="existing" type="info" show-icon :closable="false" class="mt">
          Bu hafta için zaten bir değerlendirme girmiştiniz. Düzenleyebilirsiniz.
        </el-alert>

        <el-divider />

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Ders Durumu">
              <el-input v-model="form.courseStatus as string" placeholder="örn. İyi / Orta / Düşük" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Ödev Yapıldı mı?">
              <el-radio-group v-model="form.homeworkDone">
                <el-radio :value="true">Evet</el-radio>
                <el-radio :value="false">Hayır</el-radio>
                <el-radio :value="null">Belirtme</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Motivasyon (1-5)">
              <el-rate v-model="form.motivation as number" :max="5" allow-half="false" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Davranış (1-5)">
              <el-rate v-model="form.behavior as number" :max="5" allow-half="false" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="Genel Yorum (veliye de gösterilir)">
          <el-input v-model="form.generalNote as string" type="textarea" :rows="4" placeholder="Bu hafta ile ilgili genel değerlendirmeniz..." />
        </el-form-item>

        <el-form-item label="Yöneticiye Özel Not (veliye gösterilmez)">
          <el-input v-model="form.adminOnlyNote as string" type="textarea" :rows="3" placeholder="Sadece yönetimin görmesi gereken bilgiler..." />
        </el-form-item>

        <div class="actions">
          <el-button @click="router.push('/coach/students')">İptal</el-button>
          <el-button type="primary" :loading="submitting" :disabled="(existing?.status ?? 'pending') !== 'pending'" @click="submit">
            {{ existing ? 'Güncelle' : 'Kaydet' }}
          </el-button>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<style scoped>
.actions { display: flex; justify-content: flex-end; gap: 8px; margin-top: 16px; }
.mt { margin-top: 16px; }
</style>
