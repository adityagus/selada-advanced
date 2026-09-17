<template>
  <aside
    class="hidden md:flex flex-col bg-white dark:bg-[#0f172a] border-r border-slate-200/80 dark:border-slate-800/80 fixed top-0 left-0 h-screen z-40 transition-all duration-300 ease-in-out shadow-xs select-none"
    :class="[isExpanded ? 'w-64' : 'w-[72px]']"
  >
    <!-- Brand & Toggle Section -->
    <div class="h-16 px-4 flex items-center justify-between border-b border-slate-100 dark:border-slate-800/80">
      <NuxtLink to="/dashboard" class="flex items-center gap-3 overflow-hidden">
        <div class="brand-logo flex-shrink-0">
          <NuxtImg src="/assets/img/logo-gadaimulia.png" alt="Logo Optimasi" style="width: 36px;" />
        </div>
        <span v-show="isExpanded" class="font-extrabold text-lg text-slate-900 dark:text-amber-400 tracking-tight whitespace-nowrap">
          SELADA
        </span>
      </NuxtLink>

      <!-- Sidebar Collapse Toggle Button -->
      <button
        @click="$emit('toggle-sidebar')"
        class="w-7 h-7 rounded-lg bg-slate-50 dark:bg-slate-800/60 hover:bg-slate-100 dark:hover:bg-slate-700 border border-slate-200 dark:border-slate-700 text-slate-500 dark:text-slate-400 hover:text-slate-800 dark:hover:text-white flex items-center justify-center transition-all flex-shrink-0 cursor-pointer"
        :title="isExpanded ? 'Collapse Sidebar' : 'Expand Sidebar'"
        v-show="isExpanded"
      >
        <svg class="w-4 h-4 transition-transform duration-300" :class="{ 'rotate-180': !isExpanded }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 19l-7-7 7-7m8 14l-7-7 7-7" />
        </svg>
      </button>
    </div>

    <!-- Navigation Menu -->
    <nav class="flex-1 px-3 py-4 space-y-1 overflow-y-auto custom-scrollbar">
      <!-- Dashboard Link -->
      <NuxtLink
        to="/dashboard"
        class="sidebar-nav-item"
        :class="{ 'active': route.path === '/dashboard' }"
      >
        <div class="flex items-center gap-3 min-w-0">
          <svg class="w-5 h-5 flex-shrink-0" :class="route.path === '/dashboard' ? 'text-amber-600 dark:text-amber-400' : 'text-slate-500 dark:text-slate-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"/>
          </svg>
          <span v-show="isExpanded" class="font-medium text-sm truncate">Dashboard</span>
        </div>
      </NuxtLink>

      <!-- Inquiry Link -->
      <NuxtLink
        to="/inquiry"
        class="sidebar-nav-item"
        :class="{ 'active': route.path === '/inquiry' }"
      >
        <div class="flex items-center gap-3 min-w-0">
          <svg class="w-5 h-5 flex-shrink-0" :class="route.path === '/inquiry' ? 'text-amber-600 dark:text-amber-400' : 'text-slate-500 dark:text-slate-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"/>
          </svg>
          <span v-show="isExpanded" class="font-medium text-sm truncate">Inquiry & Rencana</span>
        </div>
      </NuxtLink>


      <!-- Data Customer -->
      <NuxtLink
        to="/data-customer"
        class="sidebar-nav-item"
        :class="{ 'active': route.path === '/data-customer' }"
      >
        <div class="flex items-center gap-3 min-w-0">
          <svg class="w-5 h-5 flex-shrink-0" :class="route.path === '/data-customer' ? 'text-amber-600 dark:text-amber-400' : 'text-slate-500 dark:text-slate-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"/>
          </svg>
          <span v-show="isExpanded" class="font-medium text-sm truncate">Data Customer</span>
        </div>
      </NuxtLink>

      <!-- CeTar / Cek Taksiran -->
      <NuxtLink
        to="/cek-taksiran"
        class="sidebar-nav-item"
        :class="{ 'active': route.path === '/cek-taksiran' }"
      >
        <div class="flex items-center gap-3 min-w-0">
          <svg class="w-5 h-5 flex-shrink-0" :class="route.path === '/cek-taksiran' ? 'text-amber-600 dark:text-amber-400' : 'text-slate-500 dark:text-slate-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z"/>
          </svg>
          <span v-show="isExpanded" class="font-medium text-sm truncate">CeTar</span>
        </div>
      </NuxtLink>

      <!-- Collapsible Sub-menu: Report -->
      <div class="pt-2 border-t border-slate-100 dark:border-slate-800/80 my-2">
        <button
          @click="isReportGroupOpen = !isReportGroupOpen"
          class="w-full sidebar-nav-item justify-between group cursor-pointer"
          :class="{ 'text-slate-900 dark:text-white font-semibold': isReportGroupOpen || route.path.startsWith('/report') }"
        >
          <div class="flex items-center gap-3 min-w-0">
            <svg class="w-5 h-5 flex-shrink-0" :class="route.path.startsWith('/report') ? 'text-amber-600 dark:text-amber-400' : 'text-slate-500 dark:text-slate-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
            </svg>
            <span v-show="isExpanded" class="font-medium text-sm truncate">Report</span>
          </div>
          <svg
            v-show="isExpanded"
            class="w-4 h-4 text-slate-400 dark:text-slate-500 group-hover:text-slate-600 dark:group-hover:text-slate-300 transition-transform duration-200"
            :class="{ 'rotate-180': isReportGroupOpen }"
            fill="none" stroke="currentColor" viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
          </svg>
        </button>

        <!-- Indented Sub-items with vertical line -->
        <div v-show="isReportGroupOpen && isExpanded" class="ml-5 pl-3 border-l border-slate-200 dark:border-slate-800 space-y-1 mt-1 transition-all">
          <NuxtLink
            to="/report"
            class="subnav-item"
            :class="{ 'subnav-active': route.path === '/report' && (!route.query.tab || route.query.tab === 'Laporan Inquiry & Rencana') }"
          >
            Laporan Inquiry
          </NuxtLink>
          <NuxtLink
            to="/report?tab=Laporan Pencairan SBG"
            class="subnav-item"
            :class="{ 'subnav-active': route.path === '/report' && route.query.tab === 'Laporan Pencairan SBG' }"
          >
            Pencairan SBG
          </NuxtLink>
          <NuxtLink
            to="/report?tab=Laporan Sales Performance"
            class="subnav-item"
            :class="{ 'subnav-active': route.path === '/report' && route.query.tab === 'Laporan Sales Performance' }"
          >
            Performance Sales
          </NuxtLink>
        </div>
      </div>
    </nav>

    <!-- Sidebar User Footer -->
    <div class="p-3 border-t border-slate-100 dark:border-slate-800/80 flex items-center justify-between gap-2">
      <div class="flex items-center gap-2.5 overflow-hidden">
        <img
          :src="`https://api.dicebear.com/7.x/avataaars/svg?seed=${userName || 'user'}`"
          alt="User Avatar"
          class="w-9 h-9 rounded-full bg-slate-100 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 flex-shrink-0 object-cover"
        />
        <div v-show="isExpanded" class="flex flex-col min-w-0">
          <span class="text-xs font-bold text-slate-800 dark:text-slate-200 truncate">{{ userName || 'User' }}</span>
          <span class="text-[11px] text-slate-400 dark:text-slate-500 truncate">Cabang {{ userBranch || '-' }}</span>
        </div>
      </div>
      <button
        @click="$emit('logout')"
        class="w-8 h-8 rounded-lg bg-slate-50 dark:bg-slate-800/60 hover:bg-red-50 dark:hover:bg-red-950/40 text-slate-400 hover:text-red-600 dark:hover:text-red-400 border border-slate-200 dark:border-slate-700 hover:border-red-200 dark:hover:border-red-800 flex items-center justify-center transition-all flex-shrink-0 cursor-pointer"
        :title="isExpanded ? '' : 'Logout'"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/>
        </svg>
      </button>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from '#app'

defineProps<{
  isExpanded: boolean
  userName?: string
  userBranch?: string
}>()

defineEmits<{
  (e: 'toggle-sidebar'): void
  (e: 'logout'): void
}>()

const route = useRoute()
const isReportGroupOpen = ref(true)
</script>

<style scoped>
.sidebar-nav-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.6rem 0.75rem;
  border-radius: 0.75rem;
  color: #475569;
  font-weight: 500;
  font-size: 0.875rem;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
  border: 1px solid transparent;
}

:deep(.dark) .sidebar-nav-item {
  color: #94a3b8;
}

.sidebar-nav-item:hover {
  background-color: #f8fafc;
  color: #0f172a;
}

:deep(.dark) .sidebar-nav-item:hover {
  background-color: rgba(30, 41, 59, 0.6);
  color: #f8fafc;
}

.sidebar-nav-item.active {
  background-color: #fef3c7;
  color: #d97706;
  font-weight: 600;
  border-color: rgba(245, 158, 11, 0.25);
}

:deep(.dark) .sidebar-nav-item.active {
  background-color: rgba(217, 119, 6, 0.18);
  color: #fbbf24;
  font-weight: 600;
  border-color: rgba(251, 191, 36, 0.3);
}

.subnav-item {
  display: block;
  padding: 0.4rem 0.75rem;
  font-size: 0.775rem;
  color: #64748b;
  border-radius: 0.5rem;
  transition: all 0.15s ease;
  font-weight: 500;
}

:deep(.dark) .subnav-item {
  color: #94a3b8;
}

.subnav-item:hover {
  color: #0f172a;
  background-color: #f8fafc;
}

:deep(.dark) .subnav-item:hover {
  color: #f8fafc;
  background-color: rgba(30, 41, 59, 0.6);
}

.subnav-active {
  color: #d97706;
  font-weight: 700;
  background-color: #fef3c7;
}

:deep(.dark) .subnav-active {
  color: #fbbf24;
  background-color: rgba(217, 119, 6, 0.2);
}

/* Custom subtle scrollbar */
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #e2e8f0;
  border-radius: 2px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #cbd5e1;
}
</style>
