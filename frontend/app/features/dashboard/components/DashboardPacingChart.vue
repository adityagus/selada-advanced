<template>
  <div class="grid grid-cols-1 lg:grid-cols-12 gap-4">
    <!-- Daily Pacing Chart (8 cols) -->
    <div
      class="lg:col-span-8 bg-white dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 p-5 sm:p-6 flex flex-col justify-between space-y-6"
    >
      <!-- Header with Totals & Legend -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 pb-4 border-b border-slate-100 dark:border-slate-800">
        <div>
          <h3 class="font-semibold text-slate-900 dark:text-slate-100 text-base tracking-tight">
            Pacing Realisasi Harian
          </h3>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">
            Perbandingan omzet harian sistem Selada Gadai vs Sopiga Pro
          </p>
        </div>

        <!-- Legend & Totals Capsule -->
        <div class="flex items-center gap-2 text-xs flex-wrap">
          <div class="flex items-center gap-1.5 px-2.5 py-1 rounded-md bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-700 dark:text-slate-300">
            <span class="w-2 h-2 rounded-xs bg-slate-800 dark:bg-slate-200"></span>
            <span>Selada:</span>
            <strong class="font-mono font-medium">Rp {{ formatCompact(totalSelada) }}</strong>
          </div>

          <div class="flex items-center gap-1.5 px-2.5 py-1 rounded-md bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-700 dark:text-slate-300">
            <span class="w-2 h-2 rounded-xs bg-amber-500"></span>
            <span>Sopiga:</span>
            <strong class="font-mono font-medium">Rp {{ formatCompact(totalSopiga) }}</strong>
          </div>
        </div>
      </div>

      <!-- Pacing Chart Bars Container -->
      <div class="h-44 flex items-end justify-between gap-1 pt-4 pb-2 px-1 relative">
        <div
          v-for="d in daysInMonth"
          :key="d"
          class="flex-1 h-full flex flex-col justify-end items-center relative group"
        >
          <!-- Minimal Floating Tooltip -->
          <div
            class="absolute bottom-full mb-2 bg-slate-900 text-white text-[11px] p-2.5 rounded-lg shadow-lg opacity-0 group-hover:opacity-100 transition-opacity duration-150 z-50 whitespace-nowrap pointer-events-none space-y-1 border border-slate-800"
            :class="d <= 4 ? 'left-0' : d >= daysInMonth - 4 ? 'right-0' : 'left-1/2 -translate-x-1/2'"
          >
            <div class="font-medium text-slate-300 border-b border-slate-800 pb-1 flex items-center justify-between gap-3">
              <span>Tgl {{ d }} {{ monthName }}</span>
              <span v-if="d === currentDay" class="text-[9px] text-amber-400 font-normal">Hari Ini</span>
            </div>
            <div class="text-slate-200 font-mono">
              Selada: Rp {{ formatNum(getPacingVal(d, 'selada')) }}
            </div>
            <div class="text-amber-400 font-mono">
              Sopiga: Rp {{ formatNum(getPacingVal(d, 'sopiga')) }}
            </div>
          </div>

          <!-- Dual Column Bar -->
          <div class="w-full flex items-end justify-center gap-[1px] h-full px-[1px]">
            <!-- Selada Bar -->
            <div
              class="w-1/2 rounded-t-xs transition-all duration-300"
              :class="getPacingVal(d, 'selada') > 0 ? 'bg-slate-800 dark:bg-slate-200' : d <= currentDay ? 'bg-slate-200 dark:bg-slate-800' : 'bg-slate-100 dark:bg-slate-900'"
              :style="{
                height: `${getPacingHeight(d, 'selada')}%`,
                opacity: d <= currentDay ? 1 : 0.3
              }"
            ></div>

            <!-- Sopiga Bar -->
            <div
              class="w-1/2 rounded-t-xs transition-all duration-300"
              :class="getPacingVal(d, 'sopiga') > 0 ? 'bg-amber-500' : d <= currentDay ? 'bg-amber-200/50 dark:bg-amber-950/40' : 'bg-slate-100 dark:bg-slate-900'"
              :style="{
                height: `${getPacingHeight(d, 'sopiga')}%`,
                opacity: d <= currentDay ? 1 : 0.3
              }"
            ></div>
          </div>

          <!-- Current Day Dot -->
          <div
            v-if="d === currentDay"
            class="w-1 h-1 rounded-full bg-slate-900 dark:bg-white mt-1"
            title="Hari Ini"
          ></div>
        </div>
      </div>

      <!-- Timeline Axis Markers -->
      <div class="flex justify-between items-center text-[10px] text-slate-400 dark:text-slate-500 pt-2 border-t border-slate-100 dark:border-slate-800 font-mono">
        <span>Tgl 1</span>
        <span>Tgl 5</span>
        <span>Tgl 10</span>
        <span>Tgl 15</span>
        <span>Tgl 20</span>
        <span>Tgl 25</span>
        <span>Tgl {{ daysInMonth }}</span>
      </div>
    </div>

    <!-- Live Activity Card (4 cols) -->
    <div
      class="lg:col-span-4 bg-white dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 p-5 sm:p-6 flex flex-col justify-between space-y-4"
    >
      <div class="flex items-center justify-between pb-3 border-b border-slate-100 dark:border-slate-800">
        <div class="flex items-center gap-2">
          <h3 class="font-semibold text-slate-900 dark:text-slate-100 text-base tracking-tight">
            Aktivitas Insentif
          </h3>
          <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
        </div>
        <span class="text-xs text-slate-400 font-mono">
          {{ liveFeeds.length }} Transaksi
        </span>
      </div>

      <!-- Activity Feed List -->
      <div class="flex-1 space-y-2 overflow-y-auto max-h-[320px] pr-0.5">
        <div
          v-if="liveFeeds.length === 0"
          class="py-12 text-center text-slate-400 dark:text-slate-500 text-xs font-normal space-y-1"
        >
          <p>Belum ada transaksi insentif hari ini.</p>
        </div>

        <div
          v-for="(feed, idx) in liveFeeds"
          :key="idx"
          class="p-2.5 rounded-lg bg-slate-50 dark:bg-slate-800/40 border border-slate-100 dark:border-slate-800 flex items-center justify-between"
        >
          <div class="flex items-center gap-2.5 min-w-0">
            <div class="w-7 h-7 rounded-md bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-300 flex items-center justify-center shrink-0">
              <AppIcon name="coin" v-if="getFeedCategory(feed) === 'emas'" class="w-3.5 h-3.5" />
              <AppIcon name="monitor" v-else-if="getFeedCategory(feed) === 'nonemas'" class="w-3.5 h-3.5" />
              <AppIcon name="user-plus" v-else-if="getFeedCategory(feed) === 'cif'" class="w-3.5 h-3.5" />
              <AppIcon name="star" v-else class="w-3.5 h-3.5" />
            </div>

            <div class="min-w-0">
              <p class="text-xs font-medium text-slate-800 dark:text-slate-200 truncate">
                {{ feed.type }}
              </p>
              <p class="text-[10px] text-slate-400 truncate">
                {{ feed.code }} <span v-if="feed.customer">• {{ feed.customer }}</span>
              </p>
            </div>
          </div>

          <div class="text-right shrink-0 pl-2">
            <p class="text-xs font-semibold text-emerald-600 dark:text-emerald-400 font-mono">
              +Rp {{ formatNum(feed.amount) }}
            </p>
            <p v-if="feed.sub" class="text-[9px] text-slate-400 font-mono">
              {{ feed.sub }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { LiveFeedItem } from '../types/dashboard.types'

interface Props {
  pacingSelada?: Record<number, number>
  pacingSopiga?: Record<number, number>
  liveFeeds?: LiveFeedItem[]
}

const props = withDefaults(defineProps<Props>(), {
  pacingSelada: () => ({}),
  pacingSopiga: () => ({}),
  liveFeeds: () => []
})

const now = new Date()
const currentDay = now.getDate()
const daysInMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0).getDate()
const monthName = now.toLocaleDateString('id-ID', { month: 'short' })

const getPacingVal = (day: number, type: 'selada' | 'sopiga') => {
  const data = type === 'selada' ? props.pacingSelada : props.pacingSopiga
  return data[day] || 0
}

const totalSelada = computed(() => {
  return Object.values(props.pacingSelada || {}).reduce((acc, v) => acc + (Number(v) || 0), 0)
})

const totalSopiga = computed(() => {
  return Object.values(props.pacingSopiga || {}).reduce((acc, v) => acc + (Number(v) || 0), 0)
})

const getPacingHeight = (day: number, type: 'selada' | 'sopiga') => {
  const val = getPacingVal(day, type)
  const maxVal = 5000000
  if (day > currentDay) return 4
  if (val === 0) return 6
  return Math.min(100, Math.max(10, (val / maxVal) * 90 + 10))
}

const getFeedCategory = (feed: any) => {
  if (feed.isBonus || feed.is_bonus || (feed.type && feed.type.toLowerCase().includes('bonus'))) {
    return 'bonus'
  }
  const typeStr = (feed.type || '').toLowerCase()
  const iconStr = (feed.iconText || '').toLowerCase()
  if (typeStr.includes('cif') || typeStr.includes('nasabah') || iconStr === 'cif') {
    return 'cif'
  }
  if (typeStr.includes('elektronik') || typeStr.includes('laptop') || typeStr.includes('hp') || iconStr === 'el' || typeStr.includes('non emas')) {
    return 'nonemas'
  }
  return 'emas'
}

const formatNum = (val?: number) => {
  if (!val && val !== 0) return '0'
  return new Intl.NumberFormat('id-ID').format(val)
}

const formatCompact = (val?: number) => {
  if (!val) return '0'
  if (val >= 1000000000) return (val / 1000000000).toFixed(1) + ' M'
  if (val >= 1000000) return (val / 1000000).toFixed(1) + ' Jt'
  if (val >= 1000) return (val / 1000).toFixed(0) + ' Rb'
  return String(val)
}
</script>
