<script setup lang="ts">
import PageHeader from '@/components/PageHeader.vue'
import ThreadList from '@/components/ThreadList.vue'
</script>
<template>
  <div>
    <PageHeader title="Mesajlar" subtitle="Sorumlu olduğunuz öğrenci/veli mesajları" />
    <ThreadList :show-open-button="true" />
  </div>
</template>
