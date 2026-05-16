<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { registrationsApi } from '@/api/registrations'

const router = useRouter()
const submitting = ref(false)

const form = ref({
  fullName: '',
  phone: '',
  email: '',
  city: '',
  bio: '',
  experience: '',
})

async function submit() {
  if (!form.value.fullName || !form.value.phone) {
    ElMessage.warning('Ad soyad ve telefon zorunlu')
    return
  }
  submitting.value = true
  try {
    await registrationsApi.applyCoach(form.value)
    ElMessage.success('Başvurunuz alındı. Yönetim incelemesinden sonra koç havuzuna eklenecek ve giriş bilgileriniz SMS ile gönderilecektir.')
    router.push('/login')
  } catch {
    /* interceptor */
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="register-page">
    <div class="register-card">
      <h2>Koç Başvuru Formu</h2>
      <p class="hint">Yönetim incelemesi ardından koç havuzuna alınırsınız.</p>

      <el-form :model="form" label-position="top">
        <el-row :gutter="16">
          <el-col :span="12"><el-form-item label="Ad Soyad" required><el-input v-model="form.fullName" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="Telefon" required><el-input v-model="form.phone" /></el-form-item></el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12"><el-form-item label="E-posta"><el-input v-model="form.email" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="Şehir"><el-input v-model="form.city" /></el-form-item></el-col>
        </el-row>
        <el-form-item label="Kısa Özgeçmiş"><el-input v-model="form.bio" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="Eğitim/Mentörlük Deneyimi"><el-input v-model="form.experience" type="textarea" :rows="3" /></el-form-item>
        <div class="actions">
          <el-button @click="router.push('/login')">Geri</el-button>
          <el-button type="primary" :loading="submitting" @click="submit">Başvuruyu Gönder</el-button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<style scoped>
.register-page { min-height: 100vh; padding: 40px 20px; background: #f5f7fa; display: flex; justify-content: center; }
.register-card { width: 100%; max-width: 720px; background: #fff; padding: 32px; border-radius: 12px; box-shadow: 0 4px 16px rgba(0,0,0,0.08); }
.hint { color: #6b7280; margin-top: 0; }
.actions { display: flex; gap: 8px; justify-content: flex-end; }
</style>
