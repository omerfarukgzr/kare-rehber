<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import PageHeader from '@/components/PageHeader.vue'
import EvaluationDisplay from '@/components/EvaluationDisplay.vue'
import { evaluationsApi, type Evaluation, type EvaluationStatus, type EvaluationVersion } from '@/api/evaluations'
import { weeksApi, type EvaluationWeek } from '@/api/weeks'
import { usersApi } from '@/api/users'
import type { Page, User } from '@/types/api'

const loading = ref(false)
const data = ref<Page<Evaluation>>({ items: [], totalCount: 0, page: 1, pageSize: 20 })
const weeks = ref<EvaluationWeek[]>([])
const coaches = ref<User[]>([])

const filters = reactive<{ weekId?: string; coachId?: string; status?: EvaluationStatus; search?: string; page: number; pageSize: number }>({
  page: 1, pageSize: 20,
})

async function fetchData() {
  loading.value = true
  try {
    data.value = await evaluationsApi.list({
      page: filters.page,
      pageSize: filters.pageSize,
      weekId: filters.weekId,
      coachId: filters.coachId,
      status: filters.status,
      search: filters.search,
    })
  } finally {
    loading.value = false
  }
}

async function loadFilters() {
  weeks.value = await weeksApi.listAll()
  const r = await usersApi.list({ role: 'coach', isActive: 'true', pageSize: 200 })
  coaches.value = r.items
}

onMounted(() => { loadFilters(); fetchData() })

const detailOpen = ref(false)
const current = ref<Evaluation | null>(null)
const editMode = ref(false)
const editForm = ref<{ courseStatus: string; homeworkDone: boolean | null; motivation: number; behavior: number; generalNote: string; adminOnlyNote: string; changeReason: string }>(
  { courseStatus: '', homeworkDone: null, motivation: 0, behavior: 0, generalNote: '', adminOnlyNote: '', changeReason: '' },
)
const versions = ref<EvaluationVersion[]>([])

async function openDetail(ev: Evaluation) {
  current.value = ev
  editMode.value = false
  editForm.value = {
    courseStatus: ev.courseStatus ?? '',
    homeworkDone: ev.homeworkDone,
    motivation: ev.motivation ?? 0,
    behavior: ev.behavior ?? 0,
    generalNote: ev.generalNote ?? '',
    adminOnlyNote: ev.adminOnlyNote ?? '',
    changeReason: '',
  }
  versions.value = await evaluationsApi.versions(ev.id).catch(() => [])
  detailOpen.value = true
}

async function saveEdit() {
  if (!current.value) return
  await evaluationsApi.adminUpdate(current.value.id, {
    courseStatus: editForm.value.courseStatus,
    homeworkDone: editForm.value.homeworkDone,
    motivation: editForm.value.motivation || null,
    behavior: editForm.value.behavior || null,
    generalNote: editForm.value.generalNote,
    adminOnlyNote: editForm.value.adminOnlyNote,
    changeReason: editForm.value.changeReason || null,
  })
  ElMessage.success('Düzenleme kaydedildi (eski sürüm logta)')
  detailOpen.value = false
  fetchData()
}

async function approve() {
  if (!current.value) return
  await evaluationsApi.approve(current.value.id)
  ElMessage.success('Onaylandı')
  detailOpen.value = false
  fetchData()
}

async function revoke() {
  if (!current.value) return
  const a = await ElMessageBox.confirm('Onayı geri çekmek istiyor musunuz? Değerlendirme tekrar beklemede gözükecek.', 'Onayı Geri Çek', { type: 'warning' }).catch(() => 'cancel')
  if (a === 'cancel') return
  await evaluationsApi.revoke(current.value.id)
  ElMessage.success('Onay geri çekildi')
  detailOpen.value = false
  fetchData()
}

const statusType: Record<EvaluationStatus, '' | 'success' | 'warning' | 'info' | 'danger'> = {
  pending: 'info',
  edited_by_admin: 'warning',
  approved: 'success',
}
const statusLabel: Record<EvaluationStatus, string> = {
  pending: 'Beklemede',
  edited_by_admin: 'Düzenlendi',
  approved: 'Onaylandı',
}
</script>

<template>
  <div>
    <PageHeader title="Değerlendirmeler" subtitle="Tüm haftalık değerlendirmeler — düzenle, onayla, geçmişi gör" />

    <el-card class="filters">
      <el-form :inline="true">
        <el-form-item label="Hafta">
          <el-select v-model="filters.weekId" clearable placeholder="Hepsi" style="width: 220px">
            <el-option v-for="w in weeks" :key="w.id" :label="w.label" :value="w.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="Koç">
          <el-select v-model="filters.coachId" clearable filterable placeholder="Hepsi" style="width: 220px">
            <el-option v-for="c in coaches" :key="c.id" :label="c.fullName" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="Durum">
          <el-select v-model="filters.status" clearable placeholder="Hepsi" style="width: 160px">
            <el-option label="Beklemede" value="pending" />
            <el-option label="Düzenlendi" value="edited_by_admin" />
            <el-option label="Onaylandı" value="approved" />
          </el-select>
        </el-form-item>
        <el-button type="primary" @click="() => { filters.page = 1; fetchData() }">Filtrele</el-button>
      </el-form>
    </el-card>

    <el-card>
      <el-table v-loading="loading" :data="data.items" stripe>
        <el-table-column label="Hafta" prop="weekLabel" min-width="180" />
        <el-table-column label="Koç" prop="coachName" min-width="140" />
        <el-table-column label="Öğrenci" prop="studentName" min-width="140" />
        <el-table-column label="Durum" width="140">
          <template #default="{ row }">
            <el-tag :type="statusType[row.status as EvaluationStatus]">{{ statusLabel[row.status as EvaluationStatus] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Motivasyon" width="120">
          <template #default="{ row }">{{ row.motivation ?? '—' }}</template>
        </el-table-column>
        <el-table-column label="Davranış" width="100">
          <template #default="{ row }">{{ row.behavior ?? '—' }}</template>
        </el-table-column>
        <el-table-column label="Tarih" width="160">
          <template #default="{ row }">{{ new Date(row.submittedAt).toLocaleDateString('tr-TR') }}</template>
        </el-table-column>
        <el-table-column label="İşlem" width="100" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDetail(row)">Detay</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="pager"
        background
        layout="total, prev, pager, next"
        :total="data.totalCount"
        v-model:current-page="filters.page"
        v-model:page-size="filters.pageSize"
        @current-change="fetchData"
      />
    </el-card>

    <el-dialog v-model="detailOpen" title="Değerlendirme Detayı" width="780px" top="5vh">
      <div v-if="current">
        <EvaluationDisplay v-if="!editMode" :ev="current" :show-admin-note="true" />

        <el-form v-else label-position="top">
          <el-row :gutter="16">
            <el-col :span="12"><el-form-item label="Ders Durumu"><el-input v-model="editForm.courseStatus" /></el-form-item></el-col>
            <el-col :span="12">
              <el-form-item label="Ödev Yapıldı mı?">
                <el-radio-group v-model="editForm.homeworkDone">
                  <el-radio :value="true">Evet</el-radio>
                  <el-radio :value="false">Hayır</el-radio>
                  <el-radio :value="null">Belirtme</el-radio>
                </el-radio-group>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="16">
            <el-col :span="12"><el-form-item label="Motivasyon (1-5)"><el-rate v-model="editForm.motivation" :max="5" /></el-form-item></el-col>
            <el-col :span="12"><el-form-item label="Davranış (1-5)"><el-rate v-model="editForm.behavior" :max="5" /></el-form-item></el-col>
          </el-row>
          <el-form-item label="Genel Yorum"><el-input v-model="editForm.generalNote" type="textarea" :rows="3" /></el-form-item>
          <el-form-item label="Yönetici Notu (sadece admin görür)"><el-input v-model="editForm.adminOnlyNote" type="textarea" :rows="3" /></el-form-item>
          <el-form-item label="Değişiklik nedeni (log için)"><el-input v-model="editForm.changeReason" /></el-form-item>
        </el-form>

        <el-divider v-if="versions.length">Geçmiş Sürümler ({{ versions.length }})</el-divider>
        <el-collapse v-if="versions.length">
          <el-collapse-item v-for="v in versions" :key="v.id" :title="`v${v.versionNo} — ${new Date(v.editedAt).toLocaleString('tr-TR')}${v.changeReason ? ' — ' + v.changeReason : ''}`" :name="v.id">
            <pre class="snap">{{ JSON.stringify(v.snapshot, null, 2) }}</pre>
          </el-collapse-item>
        </el-collapse>
      </div>
      <template #footer>
        <el-button @click="detailOpen = false">Kapat</el-button>
        <template v-if="current">
          <el-button v-if="!editMode" @click="editMode = true">Düzenle</el-button>
          <el-button v-else @click="editMode = false">Düzenlemeyi İptal Et</el-button>
          <el-button v-if="editMode" type="primary" @click="saveEdit">Düzenlemeyi Kaydet</el-button>
          <el-button v-if="!editMode && current.status === 'approved'" type="warning" @click="revoke">Onayı Geri Çek</el-button>
          <el-button v-if="!editMode && current.status !== 'approved'" type="success" @click="approve">Onayla</el-button>
        </template>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.filters { margin-bottom: 16px; }
.pager { margin-top: 16px; display: flex; justify-content: flex-end; }
.snap { background: #f5f7fa; padding: 12px; border-radius: 6px; font-size: 12px; max-height: 200px; overflow: auto; }
</style>
