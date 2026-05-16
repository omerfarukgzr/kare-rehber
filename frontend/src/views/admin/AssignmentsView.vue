<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import PageHeader from '@/components/PageHeader.vue'
import { assignmentsApi, type Assignment } from '@/api/assignments'
import { usersApi } from '@/api/users'
import type { Page, User } from '@/types/api'

const loading = ref(false)
const data = ref<Page<Assignment>>({ items: [], totalCount: 0, page: 1, pageSize: 20 })

const filters = reactive<{ coachId?: string; city?: string; search?: string; page: number; pageSize: number }>({
  page: 1,
  pageSize: 20,
})

const coaches = ref<User[]>([])
const coordinators = ref<User[]>([])

async function loadCoaches() {
  const r = await usersApi.list({ role: 'coach', isActive: 'true', pageSize: 200 })
  coaches.value = r.items
}
async function loadCoordinators() {
  const r = await usersApi.list({ role: 'coordinator', isActive: 'true', pageSize: 200 })
  coordinators.value = r.items
}

async function fetchData() {
  loading.value = true
  try {
    data.value = await assignmentsApi.list({
      page: filters.page,
      pageSize: filters.pageSize,
      coachId: filters.coachId,
      city: filters.city,
      search: filters.search,
    })
  } finally {
    loading.value = false
  }
}
onMounted(() => { loadCoaches(); loadCoordinators(); fetchData() })

const assignDialog = ref(false)
const assignCoachId = ref<string>('')
const assignCity = ref('')
const assignSearch = ref('')
const studentsWithoutCoach = ref<User[]>([])
const selectedStudentIds = ref<string[]>([])
const assignLoading = ref(false)

async function openAssign() {
  assignCoachId.value = ''
  assignCity.value = ''
  assignSearch.value = ''
  selectedStudentIds.value = []
  await loadFreeStudents()
  assignDialog.value = true
}
async function loadFreeStudents() {
  const r = await assignmentsApi.studentsWithoutCoach({ city: assignCity.value, search: assignSearch.value })
  studentsWithoutCoach.value = r.items
}

async function submitAssign() {
  if (!assignCoachId.value) {
    ElMessage.warning('Koç seçin')
    return
  }
  if (selectedStudentIds.value.length === 0) {
    ElMessage.warning('En az bir öğrenci seçin')
    return
  }
  assignLoading.value = true
  try {
    const res = await assignmentsApi.assign(assignCoachId.value, selectedStudentIds.value)
    ElMessage.success(`${res.created} eşleştirme oluşturuldu${res.failed ? `, ${res.failed} başarısız` : ''}`)
    assignDialog.value = false
    fetchData()
  } finally {
    assignLoading.value = false
  }
}

async function unassign(row: Assignment) {
  const action = await ElMessageBox.confirm(
    `${row.studentName} öğrencisinin koç eşleştirmesini sonlandırmak istiyor musunuz?`,
    'Eşleştirmeyi sonlandır',
    { type: 'warning' },
  ).catch(() => 'cancel')
  if (action === 'cancel') return
  await assignmentsApi.unassign(row.studentId)
  ElMessage.success('Sonlandırıldı')
  fetchData()
}

const coordDialog = ref(false)
const coordSelectedIds = ref<string[]>([])
const coordValue = ref<string | null>(null)
function openCoordinator(rows: Assignment[]) {
  coordSelectedIds.value = rows.map(r => r.studentId)
  coordValue.value = null
  coordDialog.value = true
}
async function submitCoord() {
  const res = await assignmentsApi.setCoordinator(coordSelectedIds.value, coordValue.value)
  ElMessage.success(`${res.updated} öğrenci güncellendi`)
  coordDialog.value = false
}

const tableSelection = ref<Assignment[]>([])
function onSelectionChange(rows: Assignment[]) { tableSelection.value = rows }
</script>

<template>
  <div>
    <PageHeader title="Eşleştirmeler" subtitle="Aktif koç-öğrenci eşleştirmeleri ve koordinatör atamaları">
      <template #actions>
        <el-button type="primary" @click="openAssign">+ Yeni Eşleştirme</el-button>
        <el-button :disabled="tableSelection.length === 0" @click="openCoordinator(tableSelection)">Koordinatör Ata</el-button>
      </template>
    </PageHeader>

    <el-card class="filters">
      <el-form :inline="true">
        <el-form-item label="Koç">
          <el-select v-model="filters.coachId" clearable placeholder="Hepsi" style="width: 240px" filterable>
            <el-option v-for="c in coaches" :key="c.id" :label="`${c.fullName} (${c.phone})`" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="Şehir">
          <el-input v-model="filters.city" clearable style="width: 160px" />
        </el-form-item>
        <el-form-item label="Ara">
          <el-input v-model="filters.search" clearable style="width: 220px" @keyup.enter="() => { filters.page = 1; fetchData() }" />
        </el-form-item>
        <el-button type="primary" @click="() => { filters.page = 1; fetchData() }">Filtrele</el-button>
      </el-form>
    </el-card>

    <el-card>
      <el-table v-loading="loading" :data="data.items" stripe @selection-change="onSelectionChange">
        <el-table-column type="selection" width="48" />
        <el-table-column label="Koç" min-width="200">
          <template #default="{ row }">
            <strong>{{ row.coachName }}</strong>
            <div class="dim">{{ row.coachPhone }}</div>
          </template>
        </el-table-column>
        <el-table-column label="Öğrenci" min-width="200">
          <template #default="{ row }">
            <strong>{{ row.studentName }}</strong>
            <div class="dim">{{ row.studentPhone }}</div>
          </template>
        </el-table-column>
        <el-table-column label="Şehir" prop="studentCity" width="120" />
        <el-table-column label="Başlangıç" width="160">
          <template #default="{ row }">{{ new Date(row.startedAt).toLocaleDateString('tr-TR') }}</template>
        </el-table-column>
        <el-table-column label="İşlem" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="danger" @click="unassign(row)">Sonlandır</el-button>
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

    <el-dialog v-model="assignDialog" title="Yeni Eşleştirme" width="780px">
      <el-form>
        <el-form-item label="Koç" required>
          <el-select v-model="assignCoachId" filterable placeholder="Koç seçin" style="width: 100%">
            <el-option v-for="c in coaches" :key="c.id" :label="`${c.fullName} (${c.phone}, ${c.city ?? '—'})`" :value="c.id" />
          </el-select>
        </el-form-item>
      </el-form>

      <el-divider>Atanmamış Öğrenciler</el-divider>

      <el-form :inline="true">
        <el-form-item label="Şehir">
          <el-input v-model="assignCity" clearable style="width: 140px" />
        </el-form-item>
        <el-form-item label="Ara">
          <el-input v-model="assignSearch" clearable style="width: 200px" @keyup.enter="loadFreeStudents" />
        </el-form-item>
        <el-button @click="loadFreeStudents">Filtrele</el-button>
      </el-form>

      <el-table :data="studentsWithoutCoach" max-height="320" @selection-change="(rows: User[]) => selectedStudentIds = rows.map(r => r.id)">
        <el-table-column type="selection" width="48" />
        <el-table-column label="Ad Soyad" prop="fullName" />
        <el-table-column label="Telefon" prop="phone" width="140" />
        <el-table-column label="Şehir" prop="city" width="120" />
      </el-table>
      <p class="hint">Seçili: {{ selectedStudentIds.length }} öğrenci</p>

      <template #footer>
        <el-button @click="assignDialog = false">İptal</el-button>
        <el-button type="primary" :loading="assignLoading" @click="submitAssign">Eşleştir</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="coordDialog" title="Koordinatör Ata" width="480px">
      <p>{{ coordSelectedIds.length }} öğrenci için koordinatör ata.</p>
      <el-select v-model="coordValue" clearable placeholder="Koordinatör seç (boş bırakırsanız kaldırılır)" style="width: 100%" filterable>
        <el-option v-for="c in coordinators" :key="c.id" :label="`${c.fullName} (${c.city ?? '—'})`" :value="c.id" />
      </el-select>
      <template #footer>
        <el-button @click="coordDialog = false">İptal</el-button>
        <el-button type="primary" @click="submitCoord">Kaydet</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.filters { margin-bottom: 16px; }
.pager { margin-top: 16px; display: flex; justify-content: flex-end; }
.dim { color: #6b7280; font-size: 12px; }
.hint { color: #6b7280; font-size: 13px; margin: 8px 0; }
</style>
