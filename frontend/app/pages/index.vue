<template>
  <div class="loading-screen">
    <div class="spinner"></div>
    <p>Loading Selada V2...</p>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useCookie, navigateTo } from '#app'

definePageMeta({
  layout: false
})

const token = useCookie('auth_token')

onMounted(async () => {
  if (token.value) {
    await navigateTo('/dashboard')
  } else {
    await navigateTo('/login')
  }
})
</script>

<style scoped>
.loading-screen {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #fef3c7 0%, #ffffff 50%, #fef3c7 100%);
  color: var(--text-muted);
  gap: 1rem;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #fde68a;
  border-radius: 50%;
  border-top-color: var(--selada-400);
  animation: spin 1s ease-in-out infinite;
}
</style>
