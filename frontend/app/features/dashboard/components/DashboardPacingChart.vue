<template>
  <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
    <!-- Daily Pacing Chart (2 cols) -->
    <div
      class="lg:col-span-2 bg-white rounded-3xl border border-slate-100 p-6 shadow-xs flex flex-col justify-between"
    >
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 mb-4">
        <div>
          <h3 class="font-semibold text-slate-800 text-lg">
            Grafik Pacing Harian (Booking Selada vs Sopiga Pro)
          </h3>
          <p class="text-xs text-slate-500">
            Monitoring realisasi omzet harian bulan ini
          </p>
        </div>
        <div class="flex items-center gap-4 text-xs font-medium">
          <div class="flex items-center gap-1.5">
            <span class="w-3 h-3 rounded-full bg-emerald-500"></span>
            <span class="text-slate-600">Selada</span>
          </div>
          <div class="flex items-center gap-1.5">
            <span class="w-3 h-3 rounded-full bg-sky-500"></span>
            <span class="text-slate-600">Sopiga Pro</span>
          </div>
        </div>
      </div>

      <!-- Pacing Chart Bars -->
      <div class="h-44 flex items-end justify-between gap-1 pt-6 pb-2 px-1 relative">
        <div
          v-for="d in daysInMonth"
          :key="d"
          class="flex-1 h-full flex flex-col justify-end items-center relative group"
          :class="d === currentDay ? 'ring-2 ring-amber-400 ring-offset-1 z-10 rounded-t-sm' : ''"
        >
          <!-- Tooltip -->
          <div
            class="absolute bottom-full mb-2 bg-slate-800 text-white text-[10px] p-2.5 rounded-xl shadow-2xl opacity-0 group-hover:opacity-100 transition-opacity z-50 whitespace-nowrap pointer-events-none space-y-1"
            :class="d <= 4 ? 'left-0' : d >= daysInMonth - 4 ? 'right-0' : 'left-1/2 -translate-x-1/2'"
          >
            <div class="font-semibold border-b border-slate-700 pb-1 text-amber-400">
              Tgl {{ d }} {{ d === currentDay ? "(Hari Ini)" : "" }}
            </div>
            <div class="flex items-center gap-1.5 text-emerald-400 font-medium">
              <span class="w-2 h-2 rounded-full bg-emerald-400"></span>
              <span>Selada: Rp {{ formatNum(getPacingVal(d, "selada")) }}</span>
            </div>
            <div class="flex items-center gap-1.5 text-sky-400 font-medium">
              <span class="w-2 h-2 rounded-full bg-sky-400"></span>
              <span>Sopiga Pro: Rp {{ formatNum(getPacingVal(d, "sopiga")) }}</span>
            </div>
          </div>

          <!-- Dual Bar -->
          <div class="w-full flex items-end justify-center gap-[2px] h-full px-[1px]">
            <div
              class="w-1/2 rounded-t-sm transition-all duration-500"
              :class="getPacingVal(d, 'selada') > 0 ? 'bg-emerald-500 shadow-sm' : d <= currentDay ? 'bg-emerald-200/60' : 'bg-slate-100'"
              :style="{
                height: `${getPacingHeight(d, 'selada')}%`,
                opacity: d <= currentDay ? 1 : 0.4
              }"
            ></div>
            <div
              class="w-1/2 rounded-t-sm transition-all duration-500"
              :class="getPacingVal(d, 'sopiga') > 0 ? 'bg-sky-500 shadow-sm' : d <= currentDay ? 'bg-sky-200/60' : 'bg-slate-100'"
              :style="{
                height: `${getPacingHeight(d, 'sopiga')}%`,
                opacity: d <= currentDay ? 1 : 0.4
              }"
            ></div>
          </div>
        </div>
      </div>
      <div class="flex justify-between items-center text-[10px] text-slate-400 pt-2">
        <span>Tgl 1</span>
        <span>Tgl 15</span>
        <span>Tgl {{ daysInMonth }}</span>
      </div>
    </div>

    <!-- Live Feed (1 col) -->
    <div class="bg-white rounded-3xl border border-slate-100 p-6 shadow-xs flex flex-col">
      <div class="flex items-center justify-between mb-4">
        <h3 class="font-semibold text-slate-800 text-lg flex items-center gap-2">
          Live Insentif
          <span class="flex h-2 w-2 relative">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
          </span>
        </h3>
      </div>

      <div class="flex-1 space-y-3 overflow-y-auto max-h-[350px] pr-1">
        <div
          v-for="(feed, idx) in liveFeeds"
          :key="idx"
          class="p-3 rounded-2xl bg-slate-50 border border-slate-100 flex items-center justify-between group hover:bg-white hover:shadow-md transition-all"
        >
          <div class="flex items-center gap-3">
            <div
              class="w-9 h-9 rounded-xl flex items-center justify-center shadow-xs border transition-all"
              :class="getFeedBadgeClass(feed)"
            >
              <AppIcon name="coin" v-if="getFeedCategory(feed) === 'emas'" class="w-4 h-4" />
              <AppIcon name="monitor" v-else-if="getFeedCategory(feed) === 'nonemas'" class="w-4 h-4" />
              <AppIcon name="user-plus" v-else-if="getFeedCategory(feed) === 'cif'" class="w-4 h-4" />
              <AppIcon name="star" v-else class="w-4 h-4" />
            </div>
            <div>
              <p class="text-xs font-semibold text-slate-800">
                {{ feed.type }}
              </p>
              <p class="text-[10px] text-slate-500 uppercase tracking-tight font-medium">
                {{ feed.code }} • {{ feed.customer }}
              </p>
            </div>
          </div>
          <div class="text-right">
            <p class="text-xs font-bold text-emerald-600">
              +Rp {{ formatNum(feed.amount) }}
            </p>
            <p class="text-[9px] text-amber-600 font-medium">
              {{ feed.sub }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
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

const getPacingVal = (day: number, type: 'selada' | 'sopiga') => {
  const data = type === 'selada' ? props.pacingSelada : props.pacingSopiga
  return data[day] || 0
}

const getPacingHeight = (day: number, type: 'selada' | 'sopiga') => {
  const val = getPacingVal(day, type)
  const maxVal = 5000000
  if (day > currentDay) return 8
  if (val === 0) return 12
  return Math.min(100, Math.max(15, (val / maxVal) * 85 + 15))
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

const getFeedBadgeClass = (feed: any) => {
  const cat = getFeedCategory(feed)
  if (cat === 'bonus') return 'bg-amber-50 text-amber-600 border-amber-200/80'
  if (cat === 'cif') return 'bg-blue-50 text-blue-600 border-blue-200/80'
  if (cat === 'nonemas') return 'bg-sky-50 text-sky-600 border-sky-200/80'
  return 'bg-emerald-50 text-emerald-600 border-emerald-200/80'
}

const formatNum = (val?: number) => {
  if (!val && val !== 0) return '0'
  return new Intl.NumberFormat('id-ID').format(val)
}
</script>
