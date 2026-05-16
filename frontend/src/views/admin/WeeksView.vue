<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import PageHeader from '@/components/PageHeader.vue'
import { weeksApi, type EvaluationWeek } from '@/api/weeks'

const loading = ref(false)
const weeks = ref<EvaluationWeek[]>([])

async function fetchData() {
  loading.value = true
  try {
    weeks.value = await weeksApi.listAll()
  } finally {
    loading.value = false
  }
}
onMounted(fetchData)

const createOpen = ref(false)
const createForm = ref({ weekNo: 1, label: '', startDate: '', endDate: '' })
function openCreate() {
  const next = (weeks.value[weeks.value.length - 1]?.weekNo ?? 0) + 1
  createForm.value = { weekNo: next, label: '', startDate: '', endDate: '' }
  createOpen.value = true
}
async function submitCreate() {
  if (!createForm.value.startDate || !createForm.value.endDate) {
    ElMessage.warning('Tarihler zorunlu')
    return
  }
  await weeksApi.create(createForm.value)
  ElMessage.success('Hafta oluşturuldu')
  createOpen.value = false
  fetchData()
}

const generateOpen = ref(false)
const generateForm = ref({ startDate: '', weekCount: 36 })
async function submitGenerate() {
  if (!generateForm.value.startDate) {
    ElMessage.warning('Başlangıç tarihi gerekli')
    return
  }
  const r = await weeksApi.generate(generateForm.value)
  ElMessage.success(`${r.created} hafta oluşturuldu`)
  generateOpen.value = false
  fetchData()
}

async function toggleOpen(w: EvaluationWeek) {
  if (w.isOpen) {
    const a = await ElMessageBox.confirm(
      `${w.label} kapatılacak. Bu hafta için yeni değerlendirme girilemeyecek. Emin misiniz?`,
      'Kapat', { type: 'warning' },
    ).catch(() => 'cancel')
    if (a === 'cancel') return
    await weeksApi.close(w.id)
    ElMessage.success('Kapatıldı')
  } else {
    await weeksApi.open(w.id)
    ElMessage.success('Açıldı')
  }
  fetchData()
}

const editOpen = ref(false)
const editForm = ref<{ id: string; label: string; startDate: string; endDate: string }>({ id: '', label: '', startDate: '', endDate: '' })
function openEdit(w: EvaluationWeek) {
  editForm.value = { id: w.id, label: w.label, startDate: w.startDate, endDate: w.endDate }
  editOpen.value = true
}
async function submitEdit() {
  await weeksApi.update(editForm.value.id, {
    label: editForm.value.label,
    startDate: editForm.value.startDate,
    endDate: editForm.value.endDate,
  })
  ElMessage.success('Güncellendi')
  editOpen.value = false
  fetchData()
}
</script>

<template>
  <div>
    <PageHeader title="Hafta Yönetimi" subtitle="Akademik yıl haftalarını tanımla, aç ve kapat">
      <template #actions>
        <el-button @click="generateOpen = true">Toplu Üret (36 Hafta)</el-button>
        <el-button type="primary" @click="openCreate">+ Yeni Hafta</el-button>
      </template>
    </PageHeader>

    <el-card>
      <el-table v-loading="loading" :data="weeks" stripe>
        <el-table-column label="No" prop="weekNo" width="80" />
        <el-table-column label="Etiket" prop="label" min-width="240" />
        <el-table-column label="Başlangıç" prop="startDate" width="120" />
        <el-table-column label="Bitiş" prop="endDate" width="120" />
        <el-table-column label="Durum" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.isOpen" type="success">Açık</el-tag>
            <el-tag v-else type="info">Kapalı</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="İşlemler" width="280" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openEdit(row)">Düzenle</el-button>
            <el-button link :type="row.isOpen ? 'warning' : 'success'" @click="toggleOpen(row)">
              {{ row.isOpen ? 'Kapat' : 'Aç' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="createOpen" title="Yeni Hafta" width="500px">
      <el-form label-position="top">
        <el-form-item label="Hafta No"><el-input-number v-model="createForm.weekNo" :min="1" /></el-form-item>
        <el-form-item label="Etiket (boş bırakılırsa otomatik üretilir)"><el-input v-model="createForm.label" /></el-form-item>
        <el-form-item label="Başlangıç">
          <el-date-picker v-model="createForm.startDate" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item label="Bitiş">
          <el-date-picker v-model="createForm.endDate" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createOpen = false">İptal</el-button>
        <el-button type="primary" @click="submitCreate">Oluştur</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="generateOpen" title="Toplu Hafta Üretimi" width="480px">
      <p class="hint">Belirtilen başlangıç tarihinden itibaren n adet 7'şer günlük hafta otomatik oluşturulur.</p>
      <el-form label-position="top">
        <el-form-item label="Başlangıç (Pazartesi tercih edilir)">
          <el-date-picker v-model="generateForm.startDate" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item label="Hafta sayısı"><el-input-number v-model="generateForm.weekCount" :min="1" :max="52" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="generateOpen = false">İptal</el-button>
        <el-button type="primary" @click="submitGenerate">Üret</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="editOpen" title="Hafta Düzenle" width="500px">
      <el-form label-position="top">
        <el-form-item label="Etiket"><el-input v-model="editForm.label" /></el-form-item>
        <el-form-item label="Başlangıç"><el-date-picker v-model="editForm.startDate" type="date" value-format="YYYY-MM-DD" style="width: 100%" /></el-form-item>
        <el-form-item label="Bitiş"><el-date-picker v-model="editForm.endDate" type="date" value-format="YYYY-MM-DD" style="width: 100%" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editOpen = false">İptal</el-button>
        <el-button type="primary" @click="submitEdit">Kaydet</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.hint { color: #6b7280; font-size: 13px; }
</style>
