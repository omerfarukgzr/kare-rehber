<script setup lang="ts">
import PageHeader from '@/components/PageHeader.vue'
import ThreadList from '@/components/ThreadList.vue'
</script>
<template>
  <div>
    <PageHeader title="Mesajlar & Geri Bildirimler" subtitle="Öğrenci/veli başvuruları, geri bildirim ve şikayetler" />
    <ThreadList :show-open-button="true" />
  </div>
</template>
