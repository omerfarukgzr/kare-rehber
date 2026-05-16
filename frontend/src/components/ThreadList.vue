<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { threadsApi, type Thread, type ThreadKind, type ThreadMessage, type ThreadStatus } from '@/api/threads'
import type { Page } from '@/types/api'
import { useAuthStore } from '@/stores/auth'

const props = defineProps<{
  showOpenButton?: boolean
}>()
defineEmits<{ refresh: [] }>()

const auth = useAuthStore()

const loading = ref(false)
const data = ref<Page<Thread>>({ items: [], totalCount: 0, page: 1, pageSize: 20 })
const filters = reactive<{ kind?: ThreadKind; status?: ThreadStatus; search?: string; page: number; pageSize: number }>({ status: 'open', page: 1, pageSize: 20 })

async function fetchData() {
  loading.value = true
  try {
    data.value = await threadsApi.list({ ...filters })
  } finally {
    loading.value = false
  }
}
onMounted(fetchData)

const open = ref(false)
const current = ref<Thread | null>(null)
const messages = ref<ThreadMessage[]>([])
const reply = ref('')
const sending = ref(false)

async function openThread(t: Thread) {
  current.value = t
  messages.value = []
  open.value = true
  messages.value = await threadsApi.messages(t.id)
  fetchData()
}

async function send() {
  if (!current.value || !reply.value.trim()) return
  sending.value = true
  try {
    const m = await threadsApi.send(current.value.id, reply.value.trim())
    messages.value.push(m)
    reply.value = ''
  } finally {
    sending.value = false
  }
}

async function closeThread() {
  if (!current.value) return
  await threadsApi.close(current.value.id)
  ElMessage.success('Konuşma kapatıldı')
  open.value = false
  fetchData()
}

const createOpen = ref(false)
const createForm = ref<{ kind: ThreadKind; subject: string; body: string }>({ kind: 'general', subject: '', body: '' })
async function submitCreate() {
  if (!createForm.value.subject.trim() || !createForm.value.body.trim()) {
    ElMessage.warning('Konu ve mesaj gerekli')
    return
  }
  await threadsApi.create(createForm.value)
  ElMessage.success('Konuşma açıldı')
  createOpen.value = false
  createForm.value = { kind: 'general', subject: '', body: '' }
  fetchData()
}

const kindLabel: Record<ThreadKind, string> = { general: 'Genel', feedback: 'Geri Bildirim', complaint: 'Şikayet' }
const kindColor: Record<ThreadKind, '' | 'info' | 'warning' | 'danger'> = { general: 'info', feedback: 'warning', complaint: 'danger' }

watch(() => filters.kind, () => { filters.page = 1; fetchData() })
watch(() => filters.status, () => { filters.page = 1; fetchData() })
</script>

<template>
  <div>
    <el-card class="filters">
      <el-form :inline="true">
        <el-form-item label="Tip">
          <el-select v-model="filters.kind" clearable placeholder="Hepsi" style="width: 160px">
            <el-option label="Genel" value="general" />
            <el-option label="Geri Bildirim" value="feedback" />
            <el-option label="Şikayet" value="complaint" />
          </el-select>
        </el-form-item>
        <el-form-item label="Durum">
          <el-select v-model="filters.status" clearable placeholder="Hepsi" style="width: 140px">
            <el-option label="Açık" value="open" />
            <el-option label="Kapalı" value="closed" />
          </el-select>
        </el-form-item>
        <el-form-item label="Ara">
          <el-input v-model="filters.search" clearable style="width: 200px" @keyup.enter="() => { filters.page = 1; fetchData() }" />
        </el-form-item>
        <el-button type="primary" @click="() => { filters.page = 1; fetchData() }">Filtrele</el-button>
        <el-button v-if="props.showOpenButton" type="success" @click="createOpen = true">+ Yeni Konuşma</el-button>
      </el-form>
    </el-card>

    <el-card v-loading="loading">
      <el-empty v-if="!loading && data.items.length === 0" description="Konuşma yok" />
      <div v-else class="thread-list">
        <div v-for="t in data.items" :key="t.id" class="thread" :class="{ unread: t.unreadCount > 0 }" @click="openThread(t)">
          <div class="row">
            <el-tag size="small" :type="kindColor[t.kind]">{{ kindLabel[t.kind] }}</el-tag>
            <strong>{{ t.subject }}</strong>
            <el-badge v-if="t.unreadCount > 0" :value="t.unreadCount" class="badge" />
            <el-tag v-if="t.status === 'closed'" type="info" size="small">Kapalı</el-tag>
          </div>
          <div class="meta">
            <span>{{ t.openedByName }} ({{ t.openedByRole }})</span>
            <span v-if="t.lastMessageAt">• Son: {{ new Date(t.lastMessageAt).toLocaleString('tr-TR') }}</span>
          </div>
          <div v-if="t.lastMessageBody" class="snippet">{{ t.lastMessageBody.slice(0, 140) }}</div>
        </div>
      </div>
      <el-pagination
        class="pager"
        background layout="total, prev, pager, next"
        :total="data.totalCount"
        v-model:current-page="filters.page"
        v-model:page-size="filters.pageSize"
        @current-change="fetchData"
      />
    </el-card>

    <el-dialog v-model="open" :title="current?.subject ?? 'Konuşma'" width="640px" top="5vh">
      <div v-if="current">
        <el-tag :type="kindColor[current.kind]">{{ kindLabel[current.kind] }}</el-tag>
        <span v-if="current.status === 'closed'" class="closed-tag">Kapalı</span>
        <p class="dim">Açan: {{ current.openedByName }} ({{ current.openedByRole }})</p>

        <div class="msgs">
          <div v-for="m in messages" :key="m.id" class="msg" :class="{ mine: m.fromUserId === auth.user?.id }">
            <div class="bubble">
              <div class="from">{{ m.fromName }} <span class="role">({{ m.fromRole }})</span></div>
              <div class="body">{{ m.body }}</div>
              <div class="time">{{ new Date(m.createdAt).toLocaleString('tr-TR') }}</div>
            </div>
          </div>
        </div>

        <el-input v-if="current.status === 'open'" v-model="reply" type="textarea" :rows="3" placeholder="Yanıt yazın..." />
        <div v-if="current.status === 'open'" class="actions">
          <el-button :loading="sending" type="primary" @click="send">Gönder</el-button>
          <el-button type="warning" @click="closeThread">Konuşmayı Kapat</el-button>
        </div>
      </div>
    </el-dialog>

    <el-dialog v-model="createOpen" title="Yeni Konuşma" width="540px">
      <el-form label-position="top" :model="createForm">
        <el-form-item label="Tip" required>
          <el-select v-model="createForm.kind" style="width: 100%">
            <el-option label="Genel" value="general" />
            <el-option label="Geri Bildirim" value="feedback" />
            <el-option label="Şikayet" value="complaint" />
          </el-select>
        </el-form-item>
        <el-form-item label="Konu" required><el-input v-model="createForm.subject" placeholder="Kısa başlık" /></el-form-item>
        <el-form-item label="Mesaj" required><el-input v-model="createForm.body" type="textarea" :rows="5" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createOpen = false">İptal</el-button>
        <el-button type="primary" @click="submitCreate">Aç</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.filters { margin-bottom: 16px; }
.thread-list { display: flex; flex-direction: column; gap: 1px; background: #e5e7eb; }
.thread { background: #fff; padding: 12px 16px; cursor: pointer; transition: background 0.15s; }
.thread:hover { background: #f9fafb; }
.thread.unread { background: #fffbeb; }
.thread .row { display: flex; align-items: center; gap: 8px; margin-bottom: 4px; }
.thread .badge { margin-left: auto; }
.thread .meta { font-size: 12px; color: #6b7280; }
.thread .snippet { color: #4b5563; font-size: 13px; margin-top: 6px; }
.pager { margin-top: 16px; display: flex; justify-content: flex-end; }
.dim { color: #6b7280; font-size: 13px; }
.closed-tag { margin-left: 8px; color: #6b7280; }
.msgs { max-height: 360px; overflow-y: auto; padding: 12px 0; display: flex; flex-direction: column; gap: 12px; }
.msg { display: flex; }
.msg.mine { justify-content: flex-end; }
.msg .bubble { max-width: 70%; background: #f3f4f6; padding: 10px 14px; border-radius: 12px; }
.msg.mine .bubble { background: #dbeafe; }
.msg .from { font-size: 12px; font-weight: 600; color: #1f2937; margin-bottom: 4px; }
.msg .role { color: #6b7280; font-weight: 400; }
.msg .body { white-space: pre-wrap; line-height: 1.4; }
.msg .time { font-size: 11px; color: #6b7280; margin-top: 4px; }
.actions { display: flex; justify-content: flex-end; gap: 8px; margin-top: 12px; }
</style>
