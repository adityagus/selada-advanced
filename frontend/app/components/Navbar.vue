<template>
  <header class="sticky top-0 z-30 bg-white/90 dark:bg-[#0f172a]/90 backdrop-blur-md border-b border-slate-200/80 dark:border-slate-800 px-4 md:px-8 py-3 flex items-center justify-between gap-4 shadow-xs">
    <!-- Expand Button (Desktop collapsed state) -->
    <div class="hidden md:block">
      <button
        @click="$emit('toggle-sidebar')"
        class="w-7 h-7 rounded-lg bg-slate-50 dark:bg-slate-800 hover:bg-slate-100 dark:hover:bg-slate-700 border border-slate-200 dark:border-slate-700 text-slate-500 dark:text-slate-400 hover:text-slate-800 dark:hover:text-white flex items-center justify-center transition-all flex-shrink-0 cursor-pointer"
        :title="isExpanded ? 'Collapse Sidebar' : 'Expand Sidebar'"
        v-show="!isExpanded"
      >
        <svg class="w-4 h-4 transition-transform duration-300" :class="{ 'rotate-180': !isExpanded }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 19l-7-7 7-7m8 14l-7-7 7-7" />
        </svg>
      </button>
    </div>

    <!-- Mobile Logo & Title -->
    <div class="flex items-center gap-2 sm:hidden">
      <NuxtImg src="/assets/img/logo-gadaimulia.png" alt="Logo Optimasi" style="width: 32px;" />
      <span class="font-extrabold text-base text-slate-900 dark:text-amber-400 tracking-tight">SELADA</span>
    </div>

    <!-- Right Header Actions (Theme, Notifications, Avatar) -->
    <div class="flex items-center gap-3">
      <!-- Light / Dark Toggle (Sun / Moon Icon) -->
      <button
        @click="$emit('toggle-theme')"
        class="w-9 h-9 rounded-full border border-slate-200 dark:border-slate-700/80 text-amber-600 dark:text-amber-400 bg-amber-50/50 dark:bg-slate-800/80 hover:bg-amber-100 dark:hover:bg-slate-700 flex items-center justify-center transition-all cursor-pointer"
        :title="isDark ? 'Switch to Light Gold Theme' : 'Switch to Dark Gold & Obsidian Theme'"
      >
        <!-- Sun Icon (in Dark mode) -->
        <svg v-if="isDark" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
        </svg>
        <!-- Moon Icon (in Light mode) -->
        <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
        </svg>
      </button>

      <!-- Notification Bell with badge dot -->
      <button
        class="relative w-9 h-9 rounded-full border border-slate-200 dark:border-slate-700/80 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 flex items-center justify-center transition-all cursor-pointer"
        title="Notifications"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
        </svg>
        <span class="absolute top-2 right-2 w-2 h-2 rounded-full bg-amber-500 ring-2 ring-white dark:ring-slate-900"></span>
      </button>

      <!-- User Avatar Profile Dropdown -->
      <div class="relative">
        <button
          @click="isProfileOpen = !isProfileOpen"
          class="flex items-center gap-2 rounded-full p-0.5 border border-slate-200 dark:border-slate-700 hover:border-amber-400 transition-all cursor-pointer"
        >
          <img
            :src="`https://api.dicebear.com/7.x/avataaars/svg?seed=${userName || 'user'}`"
            alt="User Avatar"
            class="w-8 h-8 rounded-full bg-slate-100 dark:bg-slate-800 object-cover"
          />
        </button>

        <!-- Profile Dropdown Menu -->
        <div
          v-show="isProfileOpen"
          class="absolute right-0 mt-2 w-48 bg-white dark:bg-[#111827] border border-slate-200 dark:border-slate-800 rounded-2xl shadow-xl py-2 z-50 animate-fade-in"
        >
          <div class="px-4 py-2 border-b border-slate-100 dark:border-slate-800">
            <p class="text-xs font-bold text-slate-900 dark:text-slate-100 truncate">{{ userName || 'User' }}</p>
            <p class="text-[11px] text-slate-500 dark:text-slate-400">Cabang {{ userBranch || '-' }}</p>
          </div>
          <button
            @click="$emit('logout')"
            class="w-full text-left px-4 py-2 text-xs text-red-600 dark:text-red-400 font-semibold hover:bg-red-50 dark:hover:bg-red-950/40 flex items-center gap-2 transition-colors cursor-pointer"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/>
            </svg>
            Logout
          </button>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from '#app'

defineProps<{
  isExpanded: boolean
  isDark: boolean
  userName?: string
  userBranch?: string
}>()

defineEmits<{
  (e: 'toggle-sidebar'): void
  (e: 'toggle-theme'): void
  (e: 'logout'): void
}>()

const route = useRoute()
const isProfileOpen = ref(false)

watch(() => route.path, () => {
  isProfileOpen.value = false
})
</script>

<style scoped>
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-6px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fade-in {
  animation: fadeIn 0.2s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}
</style>
