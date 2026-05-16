<script setup lang="ts">
import PageHeader from '@/components/PageHeader.vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
</script>
<template>
  <div>
    <PageHeader :title="`Hoş geldiniz, ${auth.user?.fullName ?? 'Koç'}`" subtitle="Bu hafta yapılacaklar ve öğrencileriniz" />
    <el-row :gutter="20">
      <el-col :span="8"><el-card><div class="stat"><div class="num">—</div><div class="label">Atanmış Öğrenci</div></div></el-card></el-col>
      <el-col :span="8"><el-card><div class="stat"><div class="num">—</div><div class="label">Bu Hafta Girilen</div></div></el-card></el-col>
      <el-col :span="8"><el-card><div class="stat"><div class="num">—</div><div class="label">Eksik Değerlendirme</div></div></el-card></el-col>
    </el-row>
  </div>
</template>
<style scoped>
.stat { text-align: center; padding: 8px 0; }
.stat .num { font-size: 32px; font-weight: 700; color: #1f2937; }
.stat .label { color: #6b7280; font-size: 13px; margin-top: 4px; }
</style>
