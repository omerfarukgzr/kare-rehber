<script setup lang="ts">
import PageHeader from '@/components/PageHeader.vue'
import ThreadList from '@/components/ThreadList.vue'
</script>
<template>
  <div>
    <PageHeader title="Mesajlar" subtitle="Yönetim ve koordinatör ile iletişim" />
    <ThreadList :show-open-button="true" />
  </div>
</template>
