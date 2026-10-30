<template>
  <div class="bg-white dark:bg-slate-900 rounded-3xl border border-slate-100 dark:border-slate-800 shadow-xs overflow-hidden">
    <!-- Top Horizontal KPI Strip (matching reference image) -->
    <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 divide-y sm:divide-y-0 sm:divide-x divide-slate-100 dark:divide-slate-800 p-4 sm:p-5">
      <!-- KPI 1: Gadai Emas -->
      <div class="p-3 sm:px-4 sm:py-1">
        <div class="text-xs font-medium text-slate-500 dark:text-slate-400 mb-1">
          Gadai Emas
        </div>
        <div class="flex items-baseline gap-1.5 flex-wrap">
          <span class="text-lg sm:text-xl font-bold text-slate-900 dark:text-white font-mono tracking-tight">
            Rp {{ formatCompact(stats.emas) }}
          </span>
          <span class="inline-flex items-center text-[11px] font-semibold text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-950/40 px-1.5 py-0.5 rounded">
            +{{ gramasiEmas }}g
          </span>
        </div>
        <div class="text-[10px] text-slate-400 mt-0.5">
          Rp {{ formatNum(stats.emas) }}
        </div>
      </div>

      <!-- KPI 2: Gadai Elektronik -->
      <div class="p-3 sm:px-4 sm:py-1">
        <div class="text-xs font-medium text-slate-500 dark:text-slate-400 mb-1">
          Elektronik
        </div>
        <div class="flex items-baseline gap-1.5 flex-wrap">
          <span class="text-lg sm:text-xl font-bold text-slate-900 dark:text-white font-mono tracking-tight">
            Rp {{ formatCompact(stats.nonemas) }}
          </span>
          <span class="inline-flex items-center text-[11px] font-semibold text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-950/40 px-1.5 py-0.5 rounded">
            Target 50M
          </span>
        </div>
        <div class="text-[10px] text-slate-400 mt-0.5">
          Rp {{ formatNum(stats.nonemas) }}
        </div>
      </div>

      <!-- KPI 3: New CIF -->
      <div class="p-3 sm:px-4 sm:py-1">
        <div class="text-xs font-medium text-slate-500 dark:text-slate-400 mb-1">
          New CIF
        </div>
        <div class="flex items-baseline gap-1.5 flex-wrap">
          <span class="text-lg sm:text-xl font-bold text-slate-900 dark:text-white font-mono tracking-tight">
            {{ stats.newcif }}
          </span>
          <span class="inline-flex items-center text-[11px] font-semibold text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-950/40 px-1.5 py-0.5 rounded">
            +Active
          </span>
        </div>
        <div class="text-[10px] text-slate-400 mt-0.5">
          Nasabah Baru Terverifikasi
        </div>
      </div>

      <!-- KPI 4: Total Leads -->
      <div class="p-3 sm:px-4 sm:py-1">
        <div class="text-xs font-medium text-slate-500 dark:text-slate-400 mb-1">
          Leads Masuk
        </div>
        <div class="flex items-baseline gap-1.5 flex-wrap">
          <span class="text-lg sm:text-xl font-bold text-slate-900 dark:text-white font-mono tracking-tight">
            {{ monthlyFunnel.leads }}
          </span>
          <span class="inline-flex items-center text-[11px] font-semibold text-amber-700 dark:text-amber-300 bg-amber-50 dark:bg-amber-950/40 px-1.5 py-0.5 rounded">
            Bulan Ini
          </span>
        </div>
        <div class="text-[10px] text-slate-400 mt-0.5">
          Inquiry & Prospek Aktif
        </div>
      </div>

      <!-- KPI 5: Realisasi SBG / Win Rate -->
      <div class="p-3 sm:px-4 sm:py-1 col-span-2 sm:col-span-1">
        <div class="text-xs font-medium text-slate-500 dark:text-slate-400 mb-1">
          Realisasi DO / SBG
        </div>
        <div class="flex items-baseline gap-1.5 flex-wrap">
          <span class="text-lg sm:text-xl font-bold text-slate-900 dark:text-white font-mono tracking-tight">
            {{ monthlyFunnel.sbg }}
          </span>
          <span class="inline-flex items-center text-[11px] font-semibold text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-950/40 px-1.5 py-0.5 rounded">
            {{ winRate }}% Win
          </span>
        </div>
        <div class="text-[10px] text-slate-400 mt-0.5">
          Konversi Pencairan Berhasil
        </div>
      </div>
    </div>

    <!-- Divider Line -->
    <div class="border-t border-slate-100 dark:border-slate-800"></div>

    <!-- Trends Section (matching reference image chart layout) -->
    <div class="p-5 sm:p-6 space-y-4">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
        <div class="flex items-center gap-3">
          <div class="flex items-center gap-2">
            <h3 class="font-bold text-slate-900 dark:text-white text-base tracking-tight">
              Trends Pacing
            </h3>
            <span class="text-slate-350 dark:text-slate-600">···</span>
          </div>

          <!-- Legend Pills with old theme colors (Amber Selada, Sky Sopiga) -->
          <div class="flex items-center gap-3 text-xs font-medium">
            <div class="flex items-center gap-1.5 cursor-pointer">
              <span class="w-2.5 h-2.5 rounded-full bg-amber-500"></span>
              <span class="text-slate-600 dark:text-slate-300">Selada:</span>
              <span class="font-bold font-mono text-amber-600 dark:text-amber-400">Rp {{ formatCompact(totalSelada) }}</span>
            </div>

            <div class="flex items-center gap-1.5 cursor-pointer">
              <span class="w-2.5 h-2.5 rounded-full bg-sky-500"></span>
              <span class="text-slate-600 dark:text-slate-300">Sopiga:</span>
              <span class="font-bold font-mono text-sky-600 dark:text-sky-400">Rp {{ formatCompact(totalSopiga) }}</span>
            </div>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <span class="text-xs text-amber-700 dark:text-amber-300 font-medium px-2.5 py-1 bg-amber-50 dark:bg-amber-950/40 rounded-lg border border-amber-200/80 dark:border-amber-800/80">
            {{ monthName }} {{ now.getFullYear() }} ▾
          </span>
        </div>
      </div>

      <!-- Chart Canvas / Dual Column Bar Chart -->
      <div class="h-48 flex items-end justify-between gap-1 pt-6 pb-2 px-1 relative">
        <div
          v-for="d in daysInMonth"
          :key="d"
          class="flex-1 h-full flex flex-col justify-end items-center relative group"
        >
          <!-- Floating Interactive Tooltip (matching reference image) -->
          <div
            class="absolute bottom-full mb-3 bg-white dark:bg-slate-900 text-slate-800 dark:text-slate-100 text-[11px] p-2.5 rounded-xl shadow-xl opacity-0 group-hover:opacity-100 transition-all duration-150 z-50 whitespace-nowrap pointer-events-none border border-slate-200 dark:border-slate-700 space-y-1"
            :class="d <= 4 ? 'left-0' : d >= daysInMonth - 4 ? 'right-0' : 'left-1/2 -translate-x-1/2'"
          >
            <div class="font-bold border-b border-slate-100 dark:border-slate-800 pb-1 flex items-center justify-between gap-3 text-slate-700 dark:text-slate-200">
              <span>Tgl {{ d }} {{ monthName }}</span>
              <span v-if="d === currentDay" class="text-[9px] bg-amber-500 text-white px-1.5 py-0.2 rounded font-medium">Hari Ini</span>
            </div>
            <div class="flex items-center gap-2 font-mono text-amber-600 dark:text-amber-400 font-semibold">
              <span class="w-2 h-2 rounded-full bg-amber-500"></span>
              <span>Selada: Rp {{ formatNum(getPacingVal(d, 'selada')) }}</span>
            </div>
            <div class="flex items-center gap-2 font-mono text-sky-600 dark:text-sky-400 font-semibold">
              <span class="w-2 h-2 rounded-full bg-sky-500"></span>
              <span>Sopiga: Rp {{ formatNum(getPacingVal(d, 'sopiga')) }}</span>
            </div>
          </div>

          <!-- Dual Column Bar -->
          <div class="w-full flex items-end justify-center gap-[2px] h-full px-[1px]">
            <!-- Selada Bar (Amber) -->
            <div
              class="w-1/2 rounded-t-sm transition-all duration-300"
              :class="getPacingVal(d, 'selada') > 0 ? 'bg-amber-500 hover:bg-amber-600' : d <= currentDay ? 'bg-amber-100 dark:bg-amber-950/40' : 'bg-slate-100 dark:bg-slate-800'"
              :style="{
                height: `${getPacingHeight(d, 'selada')}%`,
                opacity: d <= currentDay ? 1 : 0.35
              }"
            ></div>

            <!-- Sopiga Bar (Sky) -->
            <div
              class="w-1/2 rounded-t-sm transition-all duration-300"
              :class="getPacingVal(d, 'sopiga') > 0 ? 'bg-sky-500 hover:bg-sky-600' : d <= currentDay ? 'bg-sky-100 dark:bg-sky-950/40' : 'bg-slate-100 dark:bg-slate-800'"
              :style="{
                height: `${getPacingHeight(d, 'sopiga')}%`,
                opacity: d <= currentDay ? 1 : 0.35
              }"
            ></div>
          </div>

          <!-- Current Day Marker -->
          <div
            v-if="d === currentDay"
            class="w-1.5 h-1.5 rounded-full bg-amber-500 mt-1 shadow-xs"
            title="Hari Ini"
          ></div>
        </div>
      </div>

      <!-- X-Axis Labels -->
      <div class="flex justify-between items-center text-[11px] text-slate-400 dark:text-slate-500 pt-2 border-t border-slate-100 dark:border-slate-800 font-mono">
        <span>Tgl 1</span>
        <span>Tgl 5</span>
        <span>Tgl 10</span>
        <span>Tgl 15</span>
        <span>Tgl 20</span>
        <span>Tgl 25</span>
        <span>Tgl {{ daysInMonth }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  stats: {
    emas?: number
    nonemas?: number
    newcif?: number
    [key: string]: any
  }
  gramasiEmas?: number
  monthlyFunnel: {
    leads: number
    prospect: number
    sbg: number
  }
  pacingSelada?: Record<number, number>
  pacingSopiga?: Record<number, number>
}

const props = withDefaults(defineProps<Props>(), {
  gramasiEmas: 0,
  pacingSelada: () => ({}),
  pacingSopiga: () => ({})
})

const now = new Date()
const currentDay = now.getDate()
const daysInMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0).getDate()
const monthName = now.toLocaleDateString('id-ID', { month: 'long' })

const winRate = computed(() => {
  const leads = props.monthlyFunnel?.leads || 0
  const sbg = props.monthlyFunnel?.sbg || 0
  if (!leads) return 0
  return Math.min(100, Math.round((sbg / leads) * 100))
})

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
  if (day > currentDay) return 5
  if (val === 0) return 8
  return Math.min(100, Math.max(10, (val / maxVal) * 90 + 10))
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
