<template>
  <div class="min-h-screen bg-[#f3f4f6] dark:bg-[#090d16] text-slate-800 dark:text-slate-100 font-sans flex flex-col md:flex-row overflow-x-hidden w-full transition-colors duration-200">
    <!-- Desktop Sidebar Component -->
    <Sidebar
      :is-expanded="isExpanded"
      :user-name="userName"
      :user-branch="userBranch"
      @toggle-sidebar="toggleSidebar"
      @logout="handleLogout"
    />

    <!-- Main Content Section (Navbar + Slot) -->
    <div
      class="flex-1 flex flex-col min-w-0 transition-all duration-300 ml-0"
      :class="[isExpanded ? 'md:ml-64' : 'md:ml-[72px]']"
    >
      <!-- Top Navbar Component -->
      <Navbar
        :is-expanded="isExpanded"
        :is-dark="isDark"
        :user-name="userName"
        :user-branch="userBranch"
        @toggle-sidebar="toggleSidebar"
        @toggle-theme="toggleTheme"
        @logout="handleLogout"
      />

      <!-- Page Content Slot -->
      <main class="flex-1 p-4 md:p-8 pb-20 md:pb-8 overflow-y-auto">
        <slot />
      </main>
    </div>

    <!-- Mobile Bottom Navigation -->
    <BottomNav />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useCookie, navigateTo } from '#app'
import { useApi } from '~/composables/useApi'
import { isTokenExpired, clearAuthCookies } from '~/utils/auth'
import Sidebar from '~/components/Sidebar.vue'
import Navbar from '~/components/Navbar.vue'

const token = useCookie('auth_token')
const tokenExpiry = useCookie('token_expiry')
const userName = useCookie('user_name')
const userBranch = useCookie('user_branch')

const isExpanded = ref(true)
const isDark = ref(false)

const toggleSidebar = () => {
  isExpanded.value = !isExpanded.value
}

const toggleTheme = () => {
  isDark.value = !isDark.value
  if (typeof document !== 'undefined') {
    if (isDark.value) {
      document.documentElement.classList.add('dark')
      localStorage.setItem('theme', 'dark')
    } else {
      document.documentElement.classList.remove('dark')
      localStorage.setItem('theme', 'light')
    }
  }
}

onMounted(() => {
  if (isTokenExpired(token.value, tokenExpiry.value)) {
    clearAuthCookies()
    navigateTo('/login')
  }

  // Check saved theme
  if (typeof window !== 'undefined') {
    const savedTheme = localStorage.getItem('theme')
    if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
      isDark.value = true
      document.documentElement.classList.add('dark')
    } else {
      isDark.value = false
      document.documentElement.classList.remove('dark')
    }
  }
})

const handleLogout = async () => {
  try {
    await useApi('/logout', { method: 'POST' })
  } catch (e) {
    // Ignore logout failures
  } finally {
    clearAuthCookies()
    navigateTo('/login')
  }
}
</script>
