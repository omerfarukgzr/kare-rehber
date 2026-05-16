<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import PageHeader from '@/components/PageHeader.vue'
import { registrationsApi, type Registration, type RegistrationKind, type RegistrationStatus } from '@/api/registrations'
import type { Page } from '@/types/api'

const loading = ref(false)
const data = ref<Page<Registration>>({ items: [], totalCount: 0, page: 1, pageSize: 20 })

const filters = reactive<{ kind?: RegistrationKind; status?: RegistrationStatus; search?: string; page: number; pageSize: number }>({
  status: 'pending',
  page: 1,
  pageSize: 20,
})

async function fetchData() {
  loading.value = true
  try {
    data.value = await registrationsApi.list({
      page: filters.page,
      pageSize: filters.pageSize,
      kind: filters.kind,
      status: filters.status,
      search: filters.search,
    })
  } finally {
    loading.value = false
  }
}
onMounted(fetchData)

const detailOpen = ref(false)
const detailItem = ref<Registration | null>(null)
const decisionNote = ref('')

function openDetail(item: Registration) {
  detailItem.value = item
  decisionNote.value = ''
  detailOpen.value = true
}

const decideLoading = ref(false)

async function decide(approve: boolean) {
  if (!detailItem.value) return
  if (!approve) {
    const action = await ElMessageBox.prompt('Reddetme nedenini yazın (opsiyonel)', 'Başvuru Reddi', {
      confirmButtonText: 'Reddet',
      cancelButtonText: 'İptal',
      inputType: 'textarea',
    }).catch(() => null)
    if (!action) return
    decisionNote.value = (action as { value: string }).value || ''
  }
  decideLoading.value = true
  try {
    const res = await registrationsApi.decide(detailItem.value.id, approve, decisionNote.value)
    if (approve && res.generatedPassword) {
      ElMessageBox.alert(
        `Kullanıcı oluşturuldu.\n\nÜretilen şifre: ${res.generatedPassword}\n\nBu şifre kullanıcıya SMS ile (mock) gönderildi.`,
        'Onaylandı',
        { confirmButtonText: 'Tamam' },
      )
    } else {
      ElMessage.success(approve ? 'Onaylandı' : 'Reddedildi')
    }
    detailOpen.value = false
    fetchData()
  } finally {
    decideLoading.value = false
  }
}

const statusType: Record<RegistrationStatus, '' | 'success' | 'warning' | 'info' | 'danger'> = {
  pending: 'warning',
  approved: 'success',
  rejected: 'danger',
}
const kindLabel: Record<RegistrationKind, string> = {
  student: 'Öğrenci',
  coach: 'Koç',
}
</script>

<template>
  <div>
    <PageHeader title="Başvurular" subtitle="Öğrenci ve koç başvurularını inceleyin, onaylayın veya reddedin" />

    <el-card class="filters">
      <el-form :inline="true">
        <el-form-item label="Tip">
          <el-select v-model="filters.kind" clearable placeholder="Hepsi" style="width: 160px">
            <el-option label="Öğrenci" value="student" />
            <el-option label="Koç" value="coach" />
          </el-select>
        </el-form-item>
        <el-form-item label="Durum">
          <el-select v-model="filters.status" clearable placeholder="Hepsi" style="width: 160px">
            <el-option label="Bekleyen" value="pending" />
            <el-option label="Onaylandı" value="approved" />
            <el-option label="Reddedildi" value="rejected" />
          </el-select>
        </el-form-item>
        <el-form-item label="Ara">
          <el-input v-model="filters.search" clearable placeholder="ad, telefon, e-posta" style="width: 220px" @keyup.enter="() => { filters.page = 1; fetchData() }" />
        </el-form-item>
        <el-button type="primary" @click="() => { filters.page = 1; fetchData() }">Filtrele</el-button>
      </el-form>
    </el-card>

    <el-card>
      <el-table v-loading="loading" :data="data.items" stripe>
        <el-table-column label="Tip" width="100">
          <template #default="{ row }">
            <el-tag :type="row.kind === 'coach' ? 'primary' : 'success'">{{ kindLabel[row.kind as RegistrationKind] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Ad Soyad" prop="fullName" min-width="160" />
        <el-table-column label="Telefon" prop="phone" width="140" />
        <el-table-column label="Şehir" prop="city" width="120" />
        <el-table-column label="Durum" width="120">
          <template #default="{ row }">
            <el-tag :type="statusType[row.status as RegistrationStatus]">{{ row.status === 'pending' ? 'Bekliyor' : row.status === 'approved' ? 'Onaylandı' : 'Reddedildi' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Tarih" width="180">
          <template #default="{ row }">{{ new Date(row.createdAt).toLocaleString('tr-TR') }}</template>
        </el-table-column>
        <el-table-column label="İşlem" width="120" fixed="right">
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

    <el-dialog v-model="detailOpen" title="Başvuru Detayı" width="640px">
      <div v-if="detailItem">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="Tip">{{ kindLabel[detailItem.kind] }}</el-descriptions-item>
          <el-descriptions-item label="Durum">
            <el-tag :type="statusType[detailItem.status]">{{ detailItem.status }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="Ad Soyad">{{ detailItem.fullName }}</el-descriptions-item>
          <el-descriptions-item label="Telefon">{{ detailItem.phone }}</el-descriptions-item>
          <el-descriptions-item label="E-posta">{{ detailItem.email || '—' }}</el-descriptions-item>
          <el-descriptions-item label="Şehir">{{ detailItem.city || '—' }}</el-descriptions-item>
        </el-descriptions>

        <el-divider>Form Bilgileri</el-divider>
        <pre class="payload">{{ JSON.stringify(detailItem.payload, null, 2) }}</pre>

        <el-divider v-if="detailItem.status === 'pending'">Onay Notu (opsiyonel)</el-divider>
        <el-input v-if="detailItem.status === 'pending'" v-model="decisionNote" type="textarea" :rows="2" placeholder="Onay/red notu" />
      </div>
      <template #footer>
        <el-button @click="detailOpen = false">Kapat</el-button>
        <template v-if="detailItem?.status === 'pending'">
          <el-button type="danger" :loading="decideLoading" @click="decide(false)">Reddet</el-button>
          <el-button type="success" :loading="decideLoading" @click="decide(true)">Onayla</el-button>
        </template>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.filters { margin-bottom: 16px; }
.pager { margin-top: 16px; display: flex; justify-content: flex-end; }
.payload { background: #f5f7fa; padding: 12px; border-radius: 6px; font-size: 12px; max-height: 200px; overflow: auto; }
</style>
