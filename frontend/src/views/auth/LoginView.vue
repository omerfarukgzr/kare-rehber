<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { authApi } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'
import { homeForRole } from '@/router'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  identifier: '',
  password: '',
})

const rules: FormRules = {
  identifier: [{ required: true, message: 'Telefon veya e-posta gerekli', trigger: 'blur' }],
  password: [{ required: true, message: 'Şifre gerekli', trigger: 'blur' }],
}

interface QuickUser {
  label: string
  identifier: string
  password: string
  type: 'primary' | 'success' | 'warning' | 'info' | 'danger'
}

const quickUsers: QuickUser[] = [
  { label: 'Admin',       identifier: '05000000000', password: 'Admin123!', type: 'danger'  },
  { label: 'Koç',         identifier: '05311111111', password: 'Test123!',  type: 'primary' },
  { label: 'Öğrenci',     identifier: '05322222221', password: 'Test123!',  type: 'success' },
  { label: 'Veli',        identifier: '05333333331', password: 'Test123!',  type: 'warning' },
  { label: 'Koordinatör', identifier: '05344444441', password: 'Test123!',  type: 'info'    },
]

async function doLogin(identifier: string, password: string) {
  loading.value = true
  try {
    const res = await authApi.login({ identifier, password })
    auth.setSession(res.token, res.user)
    ElMessage.success('Hoş geldiniz')
    const redirect = (route.query.redirect as string) || homeForRole(res.user.role)
    router.replace(redirect)
  } catch {
    /* axios interceptor toast eder */
  } finally {
    loading.value = false
  }
}

async function submit() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  await doLogin(form.identifier, form.password)
}

async function quickLogin(u: QuickUser) {
  form.identifier = u.identifier
  form.password = u.password
  await doLogin(u.identifier, u.password)
}
</script>

<template>
  <div class="login-page">
    <div class="login-card">
      <div class="brand">
        <div class="logo">KARE</div>
        <h1>KARE-REHBER</h1>
        <p>Koçluk Takip Sistemi</p>
      </div>

      <el-form ref="formRef" :model="form" :rules="rules" label-position="top" @submit.prevent="submit">
        <el-form-item label="Telefon veya E-posta" prop="identifier">
          <el-input v-model="form.identifier" size="large" placeholder="0555 555 55 55 / kullanici@ornek.com" />
        </el-form-item>
        <el-form-item label="Şifre" prop="password">
          <el-input v-model="form.password" type="password" size="large" show-password placeholder="••••••••" />
        </el-form-item>
        <el-button type="primary" size="large" native-type="submit" :loading="loading" class="submit-btn">
          Giriş Yap
        </el-button>
      </el-form>

      <div class="quick-login">
        <div class="quick-divider">
          <span>Test kullanıcıları</span>
        </div>
        <div class="quick-buttons">
          <el-button
            v-for="u in quickUsers"
            :key="u.label"
            :type="u.type"
            size="default"
            :loading="loading"
            @click="quickLogin(u)"
          >
            {{ u.label }}
          </el-button>
        </div>
        <p class="quick-hint">Tek tıkla seçtiğin rolle hızlıca giriş yap</p>
      </div>

      <div class="register-links">
        <span>Sisteme kayıtlı değil misiniz?</span>
        <el-button link type="primary" @click="router.push('/register/student')">Öğrenci Başvurusu</el-button>
        <el-button link type="primary" @click="router.push('/register/coach')">Koç Başvurusu</el-button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1f2937 0%, #0f172a 100%);
  padding: 20px;
}
.login-card {
  width: 100%;
  max-width: 460px;
  background: #fff;
  border-radius: 12px;
  padding: 36px 32px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.3);
}
.brand { text-align: center; margin-bottom: 28px; }
.logo {
  width: 60px; height: 60px; line-height: 60px;
  background: #1f2937; color: #fbbf24;
  border-radius: 12px; margin: 0 auto 16px;
  font-weight: 800; font-size: 18px; letter-spacing: 1px;
}
.brand h1 { font-size: 22px; margin: 0 0 6px; color: #1f2937; }
.brand p { color: #6b7280; margin: 0; font-size: 14px; }
.submit-btn { width: 100%; }

.quick-login {
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px dashed #e5e7eb;
}
.quick-divider {
  text-align: center;
  font-size: 12px;
  color: #9ca3af;
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-bottom: 12px;
}
.quick-buttons {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(80px, 1fr));
  gap: 8px;
}
.quick-buttons :deep(.el-button) {
  margin: 0;
}
.quick-hint {
  text-align: center;
  font-size: 11px;
  color: #9ca3af;
  margin: 10px 0 0;
}

.register-links {
  margin-top: 24px;
  text-align: center;
  font-size: 13px;
  color: #6b7280;
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 4px;
  align-items: center;
}
</style>
