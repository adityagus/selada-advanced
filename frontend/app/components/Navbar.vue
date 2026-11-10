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
        <AppIcon name="chevrons-left" class="w-4 h-4 transition-transform duration-300" :class="{ 'rotate-180': !isExpanded }" />
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
        <AppIcon name="sun" v-if="isDark" class="w-4 h-4" />
        <AppIcon name="moon" v-else class="w-4 h-4" />
      </button>

      <!-- Interactive Notification Bell & Popover -->
      <div class="relative" ref="notifDropdownRef">
        <button
          @click="toggleNotifications"
          class="relative w-9 h-9 rounded-full border border-slate-200 dark:border-slate-700/80 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 flex items-center justify-center transition-all cursor-pointer"
          :class="{ 'bg-amber-50 dark:bg-slate-800 border-amber-300 dark:border-amber-500/50 text-amber-600 dark:text-amber-400': isNotifOpen }"
          title="Notifikasi & Reminder Sales"
          aria-label="Lihat Notifikasi"
        >
          <AppIcon name="bell" class="w-4 h-4" />
          <!-- Animated Unread Badge Dot -->
          <span
            v-if="unreadCount > 0"
            class="absolute -top-1 -right-1 flex h-4 min-w-4 px-1 items-center justify-center rounded-full bg-rose-500 text-[10px] font-bold text-white shadow-xs ring-2 ring-white dark:ring-slate-900 animate-pulse"
          >
            {{ unreadCount > 9 ? '9+' : unreadCount }}
          </span>
        </button>

        <!-- Notification Popover Menu -->
        <div
          v-show="isNotifOpen"
          class="absolute right-0 mt-2 w-80 sm:w-96 bg-white dark:bg-[#111827] border border-slate-200 dark:border-slate-800 rounded-2xl shadow-2xl py-3 z-50 animate-fade-in"
        >
          <!-- Header -->
          <div class="px-4 pb-3 border-b border-slate-100 dark:border-slate-800/80 flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="font-bold text-sm text-slate-800 dark:text-slate-100">Notifikasi</span>
              <span
                v-if="unreadCount > 0"
                class="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-300"
              >
                {{ unreadCount }} baru
              </span>
            </div>
            <button
              v-if="unreadCount > 0"
              @click="markAllAsRead"
              class="text-xs text-amber-600 dark:text-amber-400 hover:underline font-medium cursor-pointer"
            >
              Tandai dibaca
            </button>
          </div>

          <!-- Filter Tabs -->
          <div class="px-4 pt-2.5 pb-1 flex gap-2">
            <button
              @click="activeNotifTab = 'all'"
              class="px-3 py-1 rounded-lg text-xs font-semibold transition-all cursor-pointer"
              :class="activeNotifTab === 'all' ? 'bg-amber-500 text-white shadow-xs' : 'bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400 hover:bg-slate-200'"
            >
              Semua ({{ notifications.length }})
            </button>
            <button
              @click="activeNotifTab = 'priority'"
              class="px-3 py-1 rounded-lg text-xs font-semibold transition-all cursor-pointer"
              :class="activeNotifTab === 'priority' ? 'bg-amber-500 text-white shadow-xs' : 'bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400 hover:bg-slate-200'"
            >
              Prioritas ({{ priorityCount }})
            </button>
          </div>

          <!-- Notification Items List -->
          <div class="max-h-80 overflow-y-auto divide-y divide-slate-100 dark:divide-slate-800/60 mt-1">
            <div v-if="isLoadingNotifs" class="py-8 text-center text-slate-400">
              <AppIcon name="loader" class="w-5 h-5 text-amber-500 animate-spin mx-auto mb-1" />
              <span class="text-xs">Memuat notifikasi...</span>
            </div>

            <div v-else-if="filteredNotifications.length === 0" class="py-8 text-center text-slate-400 dark:text-slate-500 px-4">
              <AppIcon name="check-circle" class="w-7 h-7 text-emerald-500 mx-auto mb-2 opacity-80" />
              <p class="text-xs font-medium">Tidak ada notifikasi baru</p>
              <p class="text-[11px] text-slate-400 mt-0.5">Semua leads dan reminder telah terpantau rapi.</p>
            </div>

            <div
              v-for="item in filteredNotifications"
              :key="item.id"
              @click="handleNotificationClick(item)"
              class="p-3.5 hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors flex items-start gap-3 cursor-pointer group"
              :class="{ 'bg-amber-50/40 dark:bg-amber-950/20': !item.isRead }"
            >
              <!-- Icon Indicator -->
              <div
                class="w-8 h-8 rounded-xl flex items-center justify-center flex-shrink-0 mt-0.5"
                :class="getNotifIconStyle(item.type)"
              >
                <AppIcon :name="item.icon" class="w-4 h-4" />
              </div>

              <!-- Content -->
              <div class="flex-1 min-w-0">
                <div class="flex items-center justify-between gap-1">
                  <p class="text-xs font-semibold text-slate-800 dark:text-slate-100 truncate group-hover:text-amber-600 dark:group-hover:text-amber-400 transition-colors">
                    {{ item.title }}
                  </p>
                  <span v-if="!item.isRead" class="w-2 h-2 rounded-full bg-rose-500 flex-shrink-0"></span>
                </div>
                <p class="text-[11px] text-slate-600 dark:text-slate-400 line-clamp-2 mt-0.5 leading-relaxed">
                  {{ item.message }}
                </p>
                <div class="flex items-center gap-2 mt-1 text-[10px] text-slate-400 dark:text-slate-500 font-medium">
                  <span class="flex items-center gap-1">
                    <AppIcon name="clock" class="w-3 h-3" />
                    {{ item.time }}
                  </span>
                  <span v-if="item.customer" class="px-1.5 py-0.2 bg-slate-100 dark:bg-slate-800 rounded text-slate-600 dark:text-slate-300">
                    {{ item.customer }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- Footer / Link to Inquiry -->
          <div class="px-4 pt-2.5 pb-1 border-t border-slate-100 dark:border-slate-800 flex items-center justify-between">
            <span class="text-[11px] text-slate-400">Total {{ notifications.length }} aktivitas</span>
            <NuxtLink
              to="/inquiry"
              @click="isNotifOpen = false"
              class="text-xs font-semibold text-amber-600 dark:text-amber-400 hover:text-amber-700 flex items-center gap-1 cursor-pointer"
            >
              Lihat di Inquiry
              <AppIcon name="arrow-right" class="w-3.5 h-3.5" />
            </NuxtLink>
          </div>
        </div>
      </div>

      <!-- User Avatar Profile Dropdown -->
      <div class="relative" ref="profileDropdownRef">
        <button
          @click="isProfileOpen = !isProfileOpen"
          class="flex items-center justify-center w-8 h-8 rounded-full border border-slate-300 dark:border-slate-700 hover:border-amber-500 dark:hover:border-amber-400 text-slate-600 dark:text-slate-300 bg-slate-50 dark:bg-slate-800 transition-all cursor-pointer"
          title="Profile"
        >
          <AppIcon name="user" class="w-4 h-4" />
        </button>

        <!-- Profile Dropdown Menu -->
        <div
          v-show="isProfileOpen"
          class="absolute right-0 mt-2 w-48 bg-white dark:bg-[#111827] border border-slate-200 dark:border-slate-800 rounded-2xl shadow-xl py-2 z-50 animate-fade-in"
        >
          <div class="px-4 py-2 border-b border-slate-100 dark:border-slate-800 flex items-center gap-2.5">
            <div class="w-7 h-7 rounded-full border border-slate-300 dark:border-slate-700 flex items-center justify-center text-slate-600 dark:text-slate-300 bg-slate-50 dark:bg-slate-800/80 flex-shrink-0">
              <AppIcon name="user" class="w-3.5 h-3.5" />
            </div>
            <div class="min-w-0">
              <p class="text-xs font-bold text-slate-900 dark:text-slate-100 truncate">{{ userName || 'User' }}</p>
              <p class="text-[11px] text-slate-500 dark:text-slate-400 truncate">{{ userBranch || '-' }}</p>
            </div>
          </div>
          <button
            @click="$emit('logout')"
            class="w-full text-left px-4 py-2 text-xs text-red-600 dark:text-red-400 font-semibold hover:bg-red-50 dark:hover:bg-red-950/40 flex items-center gap-2 transition-colors cursor-pointer"
          >
            <AppIcon name="logout" class="w-4 h-4" />
            Logout
          </button>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from '#app'
import { useApi } from '~/composables/useApi'

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

interface NotificationItem {
  id: string
  title: string
  message: string
  time: string
  customer?: string
  type: 'urgent' | 'warning' | 'info' | 'success'
  icon: 'flame' | 'clock' | 'warning' | 'check' | 'star'
  isRead: boolean
  targetLink?: string
}

const route = useRoute()
const router = useRouter()
const isProfileOpen = ref(false)
const isNotifOpen = ref(false)
const activeNotifTab = ref<'all' | 'priority'>('all')
const isLoadingNotifs = ref(false)
const notifications = ref<NotificationItem[]>([])

const notifDropdownRef = ref<HTMLElement | null>(null)
const profileDropdownRef = ref<HTMLElement | null>(null)

const unreadCount = computed(() => notifications.value.filter(n => !n.isRead).length)
const priorityCount = computed(() => notifications.value.filter(n => n.type === 'urgent' || n.type === 'warning').length)

const filteredNotifications = computed(() => {
  if (activeNotifTab.value === 'priority') {
    return notifications.value.filter(n => n.type === 'urgent' || n.type === 'warning')
  }
  return notifications.value
})

const getNotifIconStyle = (type: string) => {
  switch (type) {
    case 'urgent':
      return 'bg-rose-100 text-rose-600 dark:bg-rose-950/50 dark:text-rose-400'
    case 'warning':
      return 'bg-amber-100 text-amber-600 dark:bg-amber-950/50 dark:text-amber-400'
    case 'success':
      return 'bg-emerald-100 text-emerald-600 dark:bg-emerald-950/50 dark:text-emerald-400'
    default:
      return 'bg-blue-100 text-blue-600 dark:bg-blue-950/50 dark:text-blue-400'
  }
}

const toggleNotifications = () => {
  isNotifOpen.value = !isNotifOpen.value
  if (isNotifOpen.value && notifications.value.length === 0) {
    fetchNotifications()
  }
}

const markAllAsRead = () => {
  notifications.value.forEach(n => { n.isRead = true })
  saveReadState()
}

const handleNotificationClick = (item: NotificationItem) => {
  item.isRead = true
  saveReadState()
  isNotifOpen.value = false
  if (item.targetLink) {
    router.push(item.targetLink)
  }
}

const saveReadState = () => {
  if (typeof localStorage !== 'undefined') {
    const readIds = notifications.value.filter(n => n.isRead).map(n => n.id)
    localStorage.setItem('selada_read_notifs', JSON.stringify(readIds))
  }
}

const loadReadState = (): string[] => {
  if (typeof localStorage !== 'undefined') {
    try {
      const saved = localStorage.getItem('selada_read_notifs')
      return saved ? JSON.parse(saved) : []
    } catch {
      return []
    }
  }
  return []
}

const fetchNotifications = async () => {
  isLoadingNotifs.value = true
  try {
    const { data } = await useApi<any>('/inquiry/filter', {
      params: { start: 0, length: 25 }
    })

    const rawList = Array.isArray(data) ? data : (data?.data || [])
    const readIds = loadReadState()
    const generated: NotificationItem[] = []

    if (rawList && rawList.length > 0) {
      for (let i = 0; i < rawList.length; i++) {
        const item = rawList[i]
        const custName = item.nama || item.nm_customer || 'Customer'
        const status = (item.nama_hslaktiv || item.status || '').toLowerCase()
        const id = `notif-${item.id_rencana_det || item.id_rencana || i}`

        if (status.includes('hot')) {
          generated.push({
            id,
            title: '🔥 Hot Prospect Butuh Follow-up Segera',
            message: `Prospek ${custName} memiliki minat tinggi. Lakukan pendekatan closing untuk pencairan SBG!`,
            time: item.tgl_rencana ? new Date(item.tgl_rencana).toLocaleDateString('id-ID') : 'Hari ini',
            customer: custName,
            type: 'urgent',
            icon: 'flame',
            isRead: readIds.includes(id),
            targetLink: '/inquiry'
          })
        } else if (status.includes('prospect')) {
          generated.push({
            id,
            title: '🎯 Prospek Aktif Menunggu Konfirmasi',
            message: `Customer ${custName} telah dihubungi dan tertarik. Jadwalkan taksiran jaminan.`,
            time: item.tgl_rencana ? new Date(item.tgl_rencana).toLocaleDateString('id-ID') : 'Jadwal Aktif',
            customer: custName,
            type: 'warning',
            icon: 'clock',
            isRead: readIds.includes(id),
            targetLink: '/inquiry'
          })
        } else if (status.includes('lead')) {
          generated.push({
            id,
            title: '📋 Leads Baru Perlu Dihubungi',
            message: `Leads masuk atas nama ${custName}. Segera lakukan panggilan pertama.`,
            time: 'Baru Masuk',
            customer: custName,
            type: 'info',
            icon: 'star',
            isRead: readIds.includes(id),
            targetLink: '/inquiry'
          })
        }
      }
    }

    // Default system guidance notification if list is short
    if (generated.length === 0) {
      generated.push({
        id: 'sys-1',
        title: '✨ Selamat Datang di Selada CRM',
        message: 'Pantau target omzet sales dan realisasi SBG harian melalui menu Inquiry & Laporan.',
        time: 'Info',
        type: 'info',
        icon: 'star',
        isRead: readIds.includes('sys-1'),
        targetLink: '/inquiry'
      })
    }

    notifications.value = generated
  } catch {
    // Graceful fallback
    notifications.value = [
      {
        id: 'fallback-1',
        title: '🎯 Pantau Aktivitas Harian',
        message: 'Periksa daftar rencana inquiry dan follow-up hari ini.',
        time: 'Hari ini',
        type: 'info',
        icon: 'star',
        isRead: false,
        targetLink: '/inquiry'
      }
    ]
  } finally {
    isLoadingNotifs.value = false
  }
}

// Close popovers on click outside
const handleClickOutside = (event: MouseEvent) => {
  const target = event.target as Node
  if (notifDropdownRef.value && !notifDropdownRef.value.contains(target)) {
    isNotifOpen.value = false
  }
  if (profileDropdownRef.value && !profileDropdownRef.value.contains(target)) {
    isProfileOpen.value = false
  }
}

onMounted(() => {
  fetchNotifications()
  if (typeof document !== 'undefined') {
    document.addEventListener('click', handleClickOutside)
  }
})

onUnmounted(() => {
  if (typeof document !== 'undefined') {
    document.removeEventListener('click', handleClickOutside)
  }
})

watch(() => route.path, () => {
  isProfileOpen.value = false
  isNotifOpen.value = false
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
