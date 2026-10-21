<script setup lang="ts">
defineProps<{
  name: 'dashboard' | 'inquiry' | 'users' | 'user' | 'user-plus' | 'calculator' | 'chart' | 'logout' | 'chevrons-left' | 'chevron-down' | 'sun' | 'moon' | 'bell' | 'clock' | 'monitor' | 'coin' | 'search' | 'plus' | 'close' | 'trash' | 'refresh' | 'eye' | 'eye-off' | 'check-circle' | 'check' | 'lock' | 'bolt' | 'sparkles' | 'arrow-up' | 'arrow-left' | 'arrow-right' | 'document' | 'star' | 'trending-up' | 'message' | 'phone' | 'external-link' | 'warning' | 'info' | 'pin' | 'flame' | 'camera' | 'download' | 'edit' | 'printer' | 'diamond' | 'card' | 'ring' | 'package' | 'loader'
  title?: string
}>()
const baseURL = useRuntimeConfig().app.baseURL
</script>

<template>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="20"
    height="20"
    viewBox="0 0 24 24"
    fill="none"
    :aria-hidden="title ? undefined : true"
    :role="title ? 'img' : undefined"
    :aria-label="title"
    focusable="false"
  >
    <title v-if="title">{{ title }}</title>
    <use :href="`${baseURL}icons/sprite.svg#${name}`" />
  </svg>
</template>
