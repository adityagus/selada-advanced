<template>
  <div class="p-12 text-center space-y-3 font-sans">
    <AppIcon name="loader" class="w-8 h-8  text-emerald-500   animate-spin mx-auto" />
    <p class="text-slate-500 text-sm font-medium">Mengarahkan ke Halaman Input SBG (Deal)...</p>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

onMounted(() => {
  const id = route.query.id || route.query.id_rencana || ''
  router.replace(`/input-rencana?id=${id}&step=SBG`)
})
</script>
