<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/PageHeader.vue'
import { smsApi, type MissingCoachRow, type SmsLog, type SmsTemplate } from '@/api/sms'
import { weeksApi, type EvaluationWeek } from '@/api/weeks'
import { usersApi } from '@/api/users'
import type { Page, User } from '@/types/api'

const tab = ref('send')

const templates = ref<SmsTemplate[]>([])
const allUsers = ref<User[]>([])

async function loadInit() {
  templates.value = await smsApi.templates()
  const r = await usersApi.list({ pageSize: 500 })
  allUsers.value = r.items
}
onMounted(loadInit)

const sendForm = ref<{ recipients: string[]; phones: string; body: string; templateKey: string }>({
  recipients: [],
  phones: '',
  body: '',
  templateKey: '',
})
const sending = ref(false)

watch(() => sendForm.value.templateKey, (k) => {
  const tpl = templates.value.find(t => t.key === k)
  if (tpl) sendForm.value.body = tpl.body
})

async function send() {
  if (!sendForm.value.body.trim()) {
    ElMessage.warning('Mesaj gerekli')
    return
  }
  const phones = sendForm.value.phones.split(/[,;\n\s]+/).map(s => s.trim()).filter(Boolean)
  if (sendForm.value.recipients.length === 0 && phones.length === 0) {
    ElMessage.warning('En az bir alıcı veya telefon numarası gerekli')
    return
  }
  sending.value = true
  try {
    const res = await smsApi.send({
      userIds: sendForm.value.recipients,
      phones,
      body: sendForm.value.body,
      templateKey: sendForm.value.templateKey || undefined,
    })
    ElMessage.success(`${res.sent} SMS gönderildi (mock)${res.failed ? `, ${res.failed} başarısız` : ''}`)
    sendForm.value.recipients = []
    sendForm.value.phones = ''
    sendForm.value.body = ''
    sendForm.value.templateKey = ''
  } finally {
    sending.value = false
  }
}

const logsLoading = ref(false)
const logs = ref<Page<SmsLog>>({ items: [], totalCount: 0, page: 1, pageSize: 20 })
const logFilters = reactive<{ search?: string; page: number; pageSize: number }>({ page: 1, pageSize: 20 })
async function loadLogs() {
  logsLoading.value = true
  try {
    logs.value = await smsApi.logs({ ...logFilters })
  } finally {
    logsLoading.value = false
  }
}
watch(tab, (t) => { if (t === 'logs') loadLogs(); if (t === 'missing') loadWeeks() })

const weeks = ref<EvaluationWeek[]>([])
const selectedWeekId = ref<string>('')
const missing = ref<MissingCoachRow[]>([])
const missingLoading = ref(false)
async function loadWeeks() {
  if (weeks.value.length === 0) {
    weeks.value = await weeksApi.listAll()
    if (!selectedWeekId.value) {
      const open = weeks.value.find(w => w.isOpen)
      selectedWeekId.value = open?.id ?? weeks.value[0]?.id
    }
  }
  if (selectedWeekId.value) loadMissing()
}
async function loadMissing() {
  if (!selectedWeekId.value) return
  missingLoading.value = true
  try {
    missing.value = await smsApi.missingCoaches(selectedWeekId.value)
  } finally {
    missingLoading.value = false
  }
}
const missingSelected = ref<MissingCoachRow[]>([])
function pickMissingForSend() {
  const ids = Array.from(new Set((missingSelected.value.length ? missingSelected.value : missing.value).map(r => r.coachId)))
  sendForm.value.recipients = ids
  const week = weeks.value.find(w => w.id === selectedWeekId.value)
  sendForm.value.templateKey = 'missing_evaluation_reminder'
  sendForm.value.body = `Sayın koçumuz, ${week?.label ?? 'aktif hafta'} için öğrenci değerlendirmenizi henüz girmediniz. En kısa sürede tamamlamanızı rica ederiz.`
  tab.value = 'send'
}
</script>

<template>
  <div>
    <PageHeader title="SMS" subtitle="Tekli ve toplu SMS gönderimi (mock provider — DB'ye log)" />

    <el-tabs v-model="tab">
      <el-tab-pane label="Gönder" name="send">
        <el-card>
          <el-form label-position="top">
            <el-form-item label="Alıcılar (kullanıcılar)">
              <el-select v-model="sendForm.recipients" multiple filterable placeholder="Kullanıcı seç" style="width: 100%">
                <el-option v-for="u in allUsers" :key="u.id" :label="`${u.fullName} (${u.role}, ${u.phone})`" :value="u.id" />
              </el-select>
            </el-form-item>
            <el-form-item label="Ek telefon numaraları (virgül, satır veya boşlukla ayır)">
              <el-input v-model="sendForm.phones" type="textarea" :rows="2" placeholder="05551112233, 05552223344..." />
            </el-form-item>
            <el-form-item label="Şablon">
              <el-select v-model="sendForm.templateKey" clearable placeholder="Şablon seç (opsiyonel)" style="width: 360px">
                <el-option v-for="t in templates" :key="t.key" :label="t.name" :value="t.key" />
              </el-select>
            </el-form-item>
            <el-form-item label="Mesaj">
              <el-input v-model="sendForm.body" type="textarea" :rows="6" placeholder="Mesaj içeriği..." />
            </el-form-item>
            <el-button type="primary" :loading="sending" @click="send">Gönder (Mock)</el-button>
          </el-form>
        </el-card>
      </el-tab-pane>

      <el-tab-pane label="Eksik Koçlar" name="missing">
        <el-card>
          <el-form :inline="true">
            <el-form-item label="Hafta">
              <el-select v-model="selectedWeekId" placeholder="Hafta" style="width: 280px" @change="loadMissing">
                <el-option v-for="w in weeks" :key="w.id" :label="w.label" :value="w.id" />
              </el-select>
            </el-form-item>
            <el-button @click="loadMissing">Yenile</el-button>
            <el-button type="primary" :disabled="missing.length === 0" @click="pickMissingForSend">
              Hatırlatma SMS Hazırla ({{ Array.from(new Set((missingSelected.length ? missingSelected : missing).map(r => r.coachId))).length }} koç)
            </el-button>
          </el-form>
          <el-table v-loading="missingLoading" :data="missing" stripe @selection-change="(rows: MissingCoachRow[]) => missingSelected = rows">
            <el-table-column type="selection" width="48" />
            <el-table-column label="Koç" prop="coachName" />
            <el-table-column label="Telefon" prop="coachPhone" width="160" />
            <el-table-column label="Öğrenci" prop="studentName" />
          </el-table>
          <p v-if="missing.length === 0 && !missingLoading" class="hint">Bu hafta için tüm değerlendirmeler girilmiş 🎉</p>
        </el-card>
      </el-tab-pane>

      <el-tab-pane label="Loglar" name="logs">
        <el-card>
          <el-form :inline="true">
            <el-form-item label="Ara">
              <el-input v-model="logFilters.search" clearable placeholder="telefon veya içerik" style="width: 240px" @keyup.enter="() => { logFilters.page = 1; loadLogs() }" />
            </el-form-item>
            <el-button @click="() => { logFilters.page = 1; loadLogs() }">Filtrele</el-button>
          </el-form>
          <el-table v-loading="logsLoading" :data="logs.items" stripe>
            <el-table-column label="Tarih" width="180">
              <template #default="{ row }">{{ new Date(row.sentAt).toLocaleString('tr-TR') }}</template>
            </el-table-column>
            <el-table-column label="Telefon" prop="toPhone" width="140" />
            <el-table-column label="Şablon" prop="templateKey" width="220" />
            <el-table-column label="İçerik" prop="body" />
            <el-table-column label="Durum" prop="status" width="110">
              <template #default="{ row }">
                <el-tag :type="row.status === 'mock_sent' ? 'info' : row.status === 'sent' ? 'success' : 'danger'" size="small">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            class="pager"
            background layout="total, prev, pager, next"
            :total="logs.totalCount"
            v-model:current-page="logFilters.page"
            v-model:page-size="logFilters.pageSize"
            @current-change="loadLogs"
          />
        </el-card>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<style scoped>
.hint { color: #6b7280; font-size: 13px; }
.pager { margin-top: 16px; display: flex; justify-content: flex-end; }
</style>
