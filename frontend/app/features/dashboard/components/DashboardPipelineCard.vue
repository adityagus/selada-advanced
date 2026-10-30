<template>
  <div class="bg-white dark:bg-slate-900 rounded-3xl border border-slate-100 dark:border-slate-800 shadow-xs p-5 sm:p-6 space-y-5">
    <!-- Header with Tabs (matching reference image) -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-b border-slate-100 dark:border-slate-800 pb-4">
      <!-- Tabs (Lead Quality style in reference image) -->
      <div class="flex items-center gap-1.5 p-1 bg-slate-100 dark:bg-slate-800 rounded-xl overflow-x-auto">
        <button
          type="button"
          @click="activeTab = 'pipeline'"
          :class="activeTab === 'pipeline' ? 'bg-amber-100 dark:bg-amber-950/60 text-amber-800 dark:text-amber-300 border border-amber-300/60 dark:border-amber-700/60 shadow-2xs font-bold' : 'text-slate-500 hover:text-amber-700 dark:hover:text-amber-300 font-medium'"
          class="px-3.5 py-1.5 text-xs rounded-lg transition cursor-pointer whitespace-nowrap"
        >
          Pipeline Konversi
        </button>
        <button
          type="button"
          @click="activeTab = 'stages'"
          :class="activeTab === 'stages' ? 'bg-amber-100 dark:bg-amber-950/60 text-amber-800 dark:text-amber-300 border border-amber-300/60 dark:border-amber-700/60 shadow-2xs font-bold' : 'text-slate-500 hover:text-amber-700 dark:hover:text-amber-300 font-medium'"
          class="px-3.5 py-1.5 text-xs rounded-lg transition cursor-pointer whitespace-nowrap"
        >
          Detail Funnel
        </button>
        <button
          type="button"
          @click="activeTab = 'deals'"
          :class="activeTab === 'deals' ? 'bg-amber-100 dark:bg-amber-950/60 text-amber-800 dark:text-amber-300 border border-amber-300/60 dark:border-amber-700/60 shadow-2xs font-bold' : 'text-slate-500 hover:text-amber-700 dark:hover:text-amber-300 font-medium'"
          class="px-3.5 py-1.5 text-xs rounded-lg transition cursor-pointer whitespace-nowrap"
        >
          Statistik Deals
        </button>
      </div>

      <!-- Scope Switcher: Bulan Ini vs Total -->
      <div class="flex items-center gap-2">
        <div class="inline-flex p-0.5 bg-slate-100 dark:bg-slate-800 rounded-lg text-xs font-medium border border-slate-200/60 dark:border-slate-700/60">
          <button
            type="button"
            @click="activeScope = 'monthly'"
            :class="activeScope === 'monthly' ? 'bg-amber-500 text-white font-bold shadow-2xs' : 'text-slate-500 hover:text-amber-700 dark:hover:text-amber-300'"
            class="px-2.5 py-1 rounded-md transition cursor-pointer"
          >
            Bulan Ini
          </button>
          <button
            type="button"
            @click="activeScope = 'total'"
            :class="activeScope === 'total' ? 'bg-amber-500 text-white font-bold shadow-2xs' : 'text-slate-500 hover:text-amber-700 dark:hover:text-amber-300'"
            class="px-2.5 py-1 rounded-md transition cursor-pointer"
          >
            Total Akumulasi
          </button>
        </div>

        <button
          @click="$emit('refresh')"
          :disabled="isFetching"
          class="p-1.5 text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 rounded-lg border border-slate-200 dark:border-slate-800 hover:bg-slate-50 transition cursor-pointer disabled:opacity-50"
          title="Refresh pipeline"
        >
          <AppIcon name="refresh" class="w-3.5 h-3.5" :class="{ 'animate-spin': isFetching }" />
        </button>
      </div>
    </div>

    <!-- 3 Stages Cards Strip -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
      <!-- Stage 1: Leads Masuk -->
      <div class="bg-slate-50/80 dark:bg-slate-800/40 rounded-2xl p-4 border border-slate-100 dark:border-slate-800 space-y-2">
        <div class="flex items-center justify-between text-xs">
          <span class="font-semibold text-slate-600 dark:text-slate-300">1. Leads Masuk</span>
          <span class="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-bold bg-amber-50 text-amber-700 dark:bg-amber-950/40 dark:text-amber-400 border border-amber-200/60">
            Awal (100%)
          </span>
        </div>
        <div class="text-2xl font-bold text-slate-900 dark:text-white font-mono">
          {{ formatNum(currentData.leads) }}
        </div>
        <div class="w-full bg-slate-200 dark:bg-slate-700 rounded-full h-1.5 overflow-hidden">
          <div class="bg-amber-500 h-full rounded-full w-full"></div>
        </div>
        <p class="text-[11px] text-slate-400">Kontak & nasabah baru terdaftar</p>
      </div>

      <!-- Stage 2: Prospek -->
      <div class="bg-slate-50/80 dark:bg-slate-800/40 rounded-2xl p-4 border border-slate-100 dark:border-slate-800 space-y-2">
        <div class="flex items-center justify-between text-xs">
          <span class="font-semibold text-slate-600 dark:text-slate-300">2. Prospek</span>
          <span class="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-bold bg-sky-50 text-sky-700 dark:bg-sky-950/40 dark:text-sky-400 border border-sky-200/60">
            {{ getRate(currentData.prospect, currentData.leads) }}% Rate
          </span>
        </div>
        <div class="text-2xl font-bold text-slate-900 dark:text-white font-mono">
          {{ formatNum(currentData.prospect) }}
        </div>
        <div class="w-full bg-slate-200 dark:bg-slate-700 rounded-full h-1.5 overflow-hidden">
          <div
            class="bg-sky-500 h-full rounded-full transition-all duration-700"
            :style="{ width: `${getRate(currentData.prospect, currentData.leads)}%` }"
          ></div>
        </div>
        <p class="text-[11px] text-slate-400">Minat & taksiran terkonfirmasi</p>
      </div>

      <!-- Stage 3: Realisasi SBG / DO -->
      <div class="bg-slate-50/80 dark:bg-slate-800/40 rounded-2xl p-4 border border-slate-100 dark:border-slate-800 space-y-2">
        <div class="flex items-center justify-between text-xs">
          <span class="font-semibold text-slate-600 dark:text-slate-300">3. Realisasi SBG</span>
          <span class="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-bold bg-emerald-50 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-400 border border-emerald-200/60">
            {{ getRate(currentData.sbg, currentData.leads) }}% Win
          </span>
        </div>
        <div class="text-2xl font-bold text-slate-900 dark:text-white font-mono">
          {{ formatNum(currentData.sbg) }}
        </div>
        <div class="w-full bg-slate-200 dark:bg-slate-700 rounded-full h-1.5 overflow-hidden">
          <div
            class="bg-emerald-500 h-full rounded-full transition-all duration-700"
            :style="{ width: `${getRate(currentData.sbg, currentData.leads)}%` }"
          ></div>
        </div>
        <p class="text-[11px] text-slate-400">Surat Bukti Gadai diterbitkan</p>
      </div>
    </div>

    <!-- Expanded Detail Card with Circular Gauge Meter (matching reference image) -->
    <div class="rounded-2xl border border-slate-200 dark:border-slate-800 p-4 sm:p-5 bg-gradient-to-r from-amber-50/20 via-white to-slate-50/30 dark:from-slate-800/40 dark:to-slate-800/20 flex flex-col md:flex-row items-center justify-between gap-6">
      <!-- Left side: Summary details & action buttons -->
      <div class="space-y-3.5 flex-1 w-full">
        <div class="flex items-center justify-between gap-2 flex-wrap">
          <div class="flex items-center gap-2">
            <span class="w-2 h-2 rounded-full bg-emerald-500"></span>
            <h4 class="font-bold text-slate-800 dark:text-white text-sm">
              Ringkasan Efektivitas Konversi Sales
            </h4>
          </div>
          <span class="text-xs font-mono text-slate-400">
            Periode: {{ activeScope === 'monthly' ? 'Bulan Berjalan' : 'Akumulasi' }}
          </span>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-3 gap-3 text-xs">
          <div class="p-2.5 rounded-xl bg-white dark:bg-slate-800 border border-slate-100 dark:border-slate-700">
            <span class="text-slate-400 text-[11px] block">Leads &rarr; Prospek</span>
            <span class="font-bold text-slate-800 dark:text-slate-100 font-mono text-sm">
              {{ getRate(currentData.prospect, currentData.leads) }}%
            </span>
          </div>

          <div class="p-2.5 rounded-xl bg-white dark:bg-slate-800 border border-slate-100 dark:border-slate-700">
            <span class="text-slate-400 text-[11px] block">Prospek &rarr; Realisasi SBG</span>
            <span class="font-bold text-slate-800 dark:text-slate-100 font-mono text-sm">
              {{ getRate(currentData.sbg, currentData.prospect) }}%
            </span>
          </div>

          <div class="p-2.5 rounded-xl bg-white dark:bg-slate-800 border border-slate-100 dark:border-slate-700 col-span-2 sm:col-span-1">
            <span class="text-slate-400 text-[11px] block">Total SBG Terbit</span>
            <span class="font-bold text-emerald-600 dark:text-emerald-400 font-mono text-sm">
              {{ currentData.sbg }} Dokumen
            </span>
          </div>
        </div>

        <!-- Action Buttons (matching View Campaign button in reference image) -->
        <div class="flex items-center gap-2.5 pt-1 flex-wrap">
          <NuxtLink
            to="/inquiry"
            class="inline-flex items-center gap-2 px-4 py-2 bg-amber-500 hover:bg-amber-600 active:scale-95 text-white text-xs font-bold rounded-xl shadow-xs transition cursor-pointer"
          >
            <AppIcon name="plus" class="w-3.5 h-3.5 text-white" />
            <span>Buat Leads Baru</span>
          </NuxtLink>

          <NuxtLink
            to="/inquiry"
            class="inline-flex items-center gap-2 px-3.5 py-2 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 hover:bg-slate-50 text-slate-700 dark:text-slate-200 text-xs font-medium rounded-xl transition cursor-pointer"
          >
            <AppIcon name="document" class="w-3.5 h-3.5 text-slate-400" />
            <span>Buka Data Inquiry</span>
          </NuxtLink>
        </div>
      </div>

      <!-- Right side: Circular Progress Gauge Meter (matching 85 AI Score in reference image) -->
      <div class="flex flex-col items-center justify-center p-3 shrink-0">
        <div class="relative w-24 h-24 flex items-center justify-center">
          <svg class="w-full h-full -rotate-90 transform" viewBox="0 0 100 100">
            <!-- Background track circle -->
            <circle
              cx="50"
              cy="50"
              r="40"
              class="text-slate-100 dark:text-slate-800 stroke-current"
              stroke-width="8"
              fill="transparent"
            />
            <!-- Animated progress circle in amber/emerald -->
            <circle
              cx="50"
              cy="50"
              r="40"
              class="text-amber-500 stroke-current transition-all duration-1000 ease-out"
              stroke-width="8"
              stroke-linecap="round"
              fill="transparent"
              :stroke-dasharray="251.2"
              :stroke-dashoffset="251.2 - (251.2 * winRate) / 100"
            />
          </svg>

          <!-- Score display inside circle -->
          <div class="absolute flex flex-col items-center justify-center text-center">
            <span class="text-xl font-black text-slate-900 dark:text-white font-mono leading-none">
              {{ winRate }}
            </span>
            <span class="text-[9px] font-bold text-amber-600 dark:text-amber-400 uppercase mt-0.5">
              % Win
            </span>
          </div>
        </div>
        <span class="text-[11px] font-bold text-slate-700 dark:text-slate-300 mt-1">
          Conversion Rate
        </span>
        <span class="text-[10px] text-slate-400">
          Leads ke Realisasi
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

interface Props {
  stats: {
    leads?: number
    prospect?: number
    hotprospect?: number
    sbg?: number
    [key: string]: any
  }
  monthlyFunnel: {
    leads: number
    prospect: number
    sbg: number
    [key: string]: any
  }
  isFetching?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  isFetching: false
})

defineEmits<{
  (e: 'refresh'): void
}>()

const activeTab = ref<'pipeline' | 'stages' | 'deals'>('pipeline')
const activeScope = ref<'monthly' | 'total'>('monthly')

const currentData = computed(() => {
  if (activeScope.value === 'total') {
    return {
      leads: props.stats?.leads || 0,
      prospect: props.stats?.prospect || 0,
      sbg: props.stats?.sbg || 0
    }
  }
  return {
    leads: props.monthlyFunnel?.leads || 0,
    prospect: props.monthlyFunnel?.prospect || 0,
    sbg: props.monthlyFunnel?.sbg || 0
  }
})

const winRate = computed(() => {
  const leads = currentData.value.leads || 0
  const sbg = currentData.value.sbg || 0
  if (!leads) return 0
  return Math.min(100, Math.round((sbg / leads) * 100))
})

const getRate = (current?: number, total?: number) => {
  if (!total || total <= 0) return 0
  if (!current || current <= 0) return 0
  return Math.min(100, Math.round((current / total) * 100))
}

const formatNum = (val?: number) => {
  if (!val && val !== 0) return '0'
  return new Intl.NumberFormat('id-ID').format(val)
}
</script>
