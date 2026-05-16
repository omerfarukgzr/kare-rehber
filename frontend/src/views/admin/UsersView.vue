<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import PageHeader from '@/components/PageHeader.vue'
import { usersApi, type CreateUserPayload, type UpdateUserPayload } from '@/api/users'
import type { Page, Role, User } from '@/types/api'

const loading = ref(false)
const data = ref<Page<User>>({ items: [], totalCount: 0, page: 1, pageSize: 20 })

const filters = reactive<{ role?: Role; isActive?: 'true' | 'false'; city?: string; search?: string; page: number; pageSize: number }>({
  page: 1,
  pageSize: 20,
})

const roleLabel: Record<Role, string> = {
  admin: 'Yönetici',
  coach: 'Koç',
  student: 'Öğrenci',
  parent: 'Veli',
  coordinator: 'Koordinatör',
}
const roleType: Record<Role, '' | 'success' | 'warning' | 'info' | 'primary' | 'danger'> = {
  admin: 'danger',
  coach: 'primary',
  student: 'success',
  parent: 'info',
  coordinator: 'warning',
}

async function fetchData() {
  loading.value = true
  try {
    data.value = await usersApi.list({
      page: filters.page,
      pageSize: filters.pageSize,
      role: filters.role,
      isActive: filters.isActive,
      city: filters.city,
      search: filters.search,
    })
  } finally {
    loading.value = false
  }
}

onMounted(fetchData)

const createDialog = ref(false)
const createForm = ref<CreateUserPayload>({ role: 'coach', fullName: '', phone: '', email: '', city: '', password: '' })
const createLoading = ref(false)

function openCreate() {
  createForm.value = { role: 'coach', fullName: '', phone: '', email: '', city: '', password: '' }
  createDialog.value = true
}
async function submitCreate() {
  if (!createForm.value.fullName || !createForm.value.phone) {
    ElMessage.warning('Ad soyad ve telefon zorunlu')
    return
  }
  createLoading.value = true
  try {
    const res = await usersApi.create(createForm.value)
    if (res.generatedPassword) {
      ElMessageBox.alert(
        `Üretilen şifre: ${res.generatedPassword}\n\nBu şifre kullanıcının telefonuna SMS ile (mock) gönderildi.`,
        'Kullanıcı oluşturuldu',
        { confirmButtonText: 'Tamam' },
      )
    } else {
      ElMessage.success('Kullanıcı oluşturuldu')
    }
    createDialog.value = false
    fetchData()
  } finally {
    createLoading.value = false
  }
}

const editDialog = ref(false)
const editForm = ref<UpdateUserPayload & { id: string }>({ id: '', fullName: '', email: '', city: '', isActive: true })
function openEdit(u: User) {
  editForm.value = { id: u.id, fullName: u.fullName, email: u.email ?? '', city: u.city ?? '', isActive: u.isActive }
  editDialog.value = true
}
const editLoading = ref(false)
async function submitEdit() {
  editLoading.value = true
  try {
    await usersApi.update(editForm.value.id, {
      fullName: editForm.value.fullName ?? undefined,
      email: editForm.value.email ?? null,
      city: editForm.value.city ?? null,
      isActive: editForm.value.isActive,
    })
    ElMessage.success('Güncellendi')
    editDialog.value = false
    fetchData()
  } finally {
    editLoading.value = false
  }
}

async function toggleActive(u: User) {
  await usersApi.update(u.id, { isActive: !u.isActive })
  ElMessage.success(u.isActive ? 'Pasifleştirildi' : 'Aktifleştirildi')
  fetchData()
}

async function resetPwd(u: User) {
  const action = await ElMessageBox.confirm(
    `${u.fullName} kullanıcısının şifresini sıfırlamak istediğinize emin misiniz? Yeni şifre SMS (mock) ile gönderilir.`,
    'Şifre sıfırla',
    { confirmButtonText: 'Sıfırla', cancelButtonText: 'İptal', type: 'warning' },
  ).catch(() => 'cancel')
  if (action === 'cancel') return
  const res = await usersApi.resetPassword(u.id)
  ElMessageBox.alert(`Yeni şifre: ${res.generatedPassword}`, 'Sıfırlandı', { confirmButtonText: 'Tamam' })
}

function applyFilters() {
  filters.page = 1
  fetchData()
}
function clearFilters() {
  filters.role = undefined
  filters.isActive = undefined
  filters.city = undefined
  filters.search = undefined
  applyFilters()
}
</script>

<template>
  <div>
    <PageHeader title="Kullanıcılar" subtitle="Tüm sistem kullanıcıları">
      <template #actions>
        <el-button type="primary" @click="openCreate">+ Yeni Kullanıcı</el-button>
      </template>
    </PageHeader>

    <el-card class="filters">
      <el-form :inline="true">
        <el-form-item label="Rol">
          <el-select v-model="filters.role" placeholder="Hepsi" clearable style="width: 160px">
            <el-option label="Yönetici" value="admin" />
            <el-option label="Koç" value="coach" />
            <el-option label="Öğrenci" value="student" />
            <el-option label="Veli" value="parent" />
            <el-option label="Koordinatör" value="coordinator" />
          </el-select>
        </el-form-item>
        <el-form-item label="Durum">
          <el-select v-model="filters.isActive" placeholder="Hepsi" clearable style="width: 140px">
            <el-option label="Aktif" value="true" />
            <el-option label="Pasif" value="false" />
          </el-select>
        </el-form-item>
        <el-form-item label="Şehir">
          <el-input v-model="filters.city" placeholder="Şehir" clearable style="width: 160px" />
        </el-form-item>
        <el-form-item label="Ara">
          <el-input v-model="filters.search" placeholder="ad, telefon, e-posta" clearable style="width: 220px" @keyup.enter="applyFilters" />
        </el-form-item>
        <el-button type="primary" @click="applyFilters">Filtrele</el-button>
        <el-button @click="clearFilters">Temizle</el-button>
      </el-form>
    </el-card>

    <el-card>
      <el-table v-loading="loading" :data="data.items" stripe>
        <el-table-column label="Ad Soyad" prop="fullName" min-width="180">
          <template #default="{ row }">
            <div>
              <strong>{{ row.fullName }}</strong>
              <el-tag v-if="!row.isActive" type="info" size="small" style="margin-left: 6px">Pasif</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="Rol" width="130">
          <template #default="{ row }">
            <el-tag :type="roleType[row.role as Role]">{{ roleLabel[row.role as Role] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Telefon" prop="phone" width="160" />
        <el-table-column label="E-posta" prop="email" min-width="180" />
        <el-table-column label="Şehir" prop="city" width="120" />
        <el-table-column label="İşlemler" width="280" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openEdit(row)">Düzenle</el-button>
            <el-button link type="warning" @click="resetPwd(row)">Şifre Sıfırla</el-button>
            <el-button link :type="row.isActive ? 'danger' : 'success'" @click="toggleActive(row)">
              {{ row.isActive ? 'Pasifleştir' : 'Aktifleştir' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="pager"
        background
        layout="total, prev, pager, next, sizes"
        :total="data.totalCount"
        v-model:current-page="filters.page"
        v-model:page-size="filters.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        @current-change="fetchData"
        @size-change="fetchData"
      />
    </el-card>

    <el-dialog v-model="createDialog" title="Yeni Kullanıcı" width="540px">
      <el-form :model="createForm" label-position="top">
        <el-form-item label="Rol" required>
          <el-select v-model="createForm.role" style="width: 100%">
            <el-option label="Yönetici" value="admin" />
            <el-option label="Koordinatör" value="coordinator" />
            <el-option label="Koç" value="coach" />
            <el-option label="Öğrenci" value="student" />
            <el-option label="Veli" value="parent" />
          </el-select>
        </el-form-item>
        <el-form-item label="Ad Soyad" required><el-input v-model="createForm.fullName" /></el-form-item>
        <el-form-item label="Telefon" required><el-input v-model="createForm.phone" placeholder="05XX XXX XX XX" /></el-form-item>
        <el-form-item label="E-posta"><el-input v-model="createForm.email" /></el-form-item>
        <el-form-item label="Şehir"><el-input v-model="createForm.city" /></el-form-item>
        <el-form-item label="Şifre (boş bırakılırsa otomatik üretilir ve SMS ile gönderilir)">
          <el-input v-model="createForm.password" type="password" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialog = false">İptal</el-button>
        <el-button type="primary" :loading="createLoading" @click="submitCreate">Oluştur</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="editDialog" title="Kullanıcı Düzenle" width="540px">
      <el-form :model="editForm" label-position="top">
        <el-form-item label="Ad Soyad"><el-input v-model="editForm.fullName" /></el-form-item>
        <el-form-item label="E-posta"><el-input v-model="editForm.email as string" /></el-form-item>
        <el-form-item label="Şehir"><el-input v-model="editForm.city as string" /></el-form-item>
        <el-form-item label="Aktif"><el-switch v-model="editForm.isActive" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialog = false">İptal</el-button>
        <el-button type="primary" :loading="editLoading" @click="submitEdit">Kaydet</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.filters { margin-bottom: 16px; }
.pager { margin-top: 16px; display: flex; justify-content: flex-end; }
</style>
