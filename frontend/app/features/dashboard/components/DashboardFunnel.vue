<template>
  <div class="bg-white dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 p-5 sm:p-6 space-y-5">
    <!-- Header with View Toggle -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 pb-4 border-b border-slate-100 dark:border-slate-800">
      <div>
        <h3 class="font-semibold text-slate-900 dark:text-slate-100 text-base tracking-tight">
          Pipeline Konversi
        </h3>
        <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">
          Tingkat konversi nasabah dari leads hingga realisasi SBG
        </p>
      </div>

      <div class="flex items-center gap-2">
        <div class="inline-flex p-0.5 bg-slate-100 dark:bg-slate-800 rounded-lg border border-slate-200/60 dark:border-slate-700/60">
          <button
            type="button"
            @click="activeScope = 'monthly'"
            :class="activeScope === 'monthly' ? 'bg-white dark:bg-slate-700 text-slate-900 dark:text-white font-medium shadow-2xs' : 'text-slate-500 hover:text-slate-800 dark:hover:text-slate-200 font-normal'"
            class="px-2.5 py-1 text-xs rounded-md transition cursor-pointer"
          >
            Bulan Ini
          </button>
          <button
            type="button"
            @click="activeScope = 'total'"
            :class="activeScope === 'total' ? 'bg-white dark:bg-slate-700 text-slate-900 dark:text-white font-medium shadow-2xs' : 'text-slate-500 hover:text-slate-800 dark:hover:text-slate-200 font-normal'"
            class="px-2.5 py-1 text-xs rounded-md transition cursor-pointer"
          >
            Total Akumulasi
          </button>
        </div>

        <button
          @click="$emit('refresh')"
          :disabled="isFetching"
          type="button"
          class="p-1.5 text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 rounded-lg border border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800 transition cursor-pointer disabled:opacity-50"
          title="Refresh pipeline"
        >
          <AppIcon name="refresh" class="w-3.5 h-3.5" :class="{ 'animate-spin': isFetching }" />
        </button>
      </div>
    </div>

    <!-- 3-Stage Minimal Pipeline -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
      <!-- Stage 1: Leads -->
      <div class="bg-slate-50 dark:bg-slate-800/40 rounded-lg p-4 border border-slate-100 dark:border-slate-800 space-y-3">
        <div class="flex items-center justify-between text-xs">
          <span class="font-medium text-slate-500 dark:text-slate-400">1. Leads Masuk</span>
          <span class="font-mono text-slate-400">100%</span>
        </div>

        <div>
          <div class="text-2xl font-semibold text-slate-900 dark:text-slate-100 font-mono tracking-tight">
            {{ formatNum(activeFunnelData.leads) }}
          </div>
          <p class="text-xs text-slate-400 mt-0.5">Kontak & inquiry baru</p>
        </div>

        <div class="w-full bg-slate-200 dark:bg-slate-700 rounded-full h-1">
          <div class="bg-slate-700 dark:bg-slate-300 h-full rounded-full w-full"></div>
        </div>
      </div>

      <!-- Stage 2: Prospek -->
      <div class="bg-slate-50 dark:bg-slate-800/40 rounded-lg p-4 border border-slate-100 dark:border-slate-800 space-y-3">
        <div class="flex items-center justify-between text-xs">
          <span class="font-medium text-slate-500 dark:text-slate-400">2. Prospek</span>
          <span class="font-mono text-slate-700 dark:text-slate-300 font-medium">
            {{ getRate(activeFunnelData.prospect, activeFunnelData.leads) }}%
          </span>
        </div>

        <div>
          <div class="text-2xl font-semibold text-slate-900 dark:text-slate-100 font-mono tracking-tight">
            {{ formatNum(activeFunnelData.prospect) }}
          </div>
          <p class="text-xs text-slate-400 mt-0.5">Minat terkonfirmasi</p>
        </div>

        <div class="w-full bg-slate-200 dark:bg-slate-700 rounded-full h-1 overflow-hidden">
          <div
            class="bg-slate-700 dark:bg-slate-300 h-full rounded-full transition-all duration-500"
            :style="{ width: `${getRate(activeFunnelData.prospect, activeFunnelData.leads)}%` }"
          ></div>
        </div>
      </div>

      <!-- Stage 3: Realisasi SBG -->
      <div class="bg-slate-50 dark:bg-slate-800/40 rounded-lg p-4 border border-slate-100 dark:border-slate-800 space-y-3">
        <div class="flex items-center justify-between text-xs">
          <span class="font-medium text-slate-500 dark:text-slate-400">3. Realisasi SBG / DO</span>
          <span class="font-mono text-emerald-600 dark:text-emerald-400 font-medium">
            {{ getRate(activeFunnelData.sbg, activeFunnelData.leads) }}%
          </span>
        </div>

        <div>
          <div class="text-2xl font-semibold text-slate-900 dark:text-slate-100 font-mono tracking-tight">
            {{ formatNum(activeFunnelData.sbg) }}
          </div>
          <p class="text-xs text-slate-400 mt-0.5">Pencairan selesai</p>
        </div>

        <div class="w-full bg-slate-200 dark:bg-slate-700 rounded-full h-1 overflow-hidden">
          <div
            class="bg-emerald-500 h-full rounded-full transition-all duration-500"
            :style="{ width: `${getRate(activeFunnelData.sbg, activeFunnelData.leads)}%` }"
          ></div>
        </div>
      </div>
    </div>

    <!-- Summary Metrics Strip -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-2 pt-1 text-xs">
      <div class="flex items-center justify-between py-2 px-3 rounded-lg bg-slate-50 dark:bg-slate-800/60 border border-slate-100 dark:border-slate-800">
        <span class="text-slate-500 dark:text-slate-400">Konversi Total (Leads ke SBG):</span>
        <span class="font-medium font-mono text-emerald-600 dark:text-emerald-400">
          {{ getRate(activeFunnelData.sbg, activeFunnelData.leads) }}%
        </span>
      </div>
      <div class="flex items-center justify-between py-2 px-3 rounded-lg bg-slate-50 dark:bg-slate-800/60 border border-slate-100 dark:border-slate-800">
        <span class="text-slate-500 dark:text-slate-400">Leads ke Prospek:</span>
        <span class="font-medium font-mono text-slate-700 dark:text-slate-300">
          {{ getRate(activeFunnelData.prospect, activeFunnelData.leads) }}%
        </span>
      </div>
      <div class="flex items-center justify-between py-2 px-3 rounded-lg bg-slate-50 dark:bg-slate-800/60 border border-slate-100 dark:border-slate-800">
        <span class="text-slate-500 dark:text-slate-400">Prospek ke Realisasi:</span>
        <span class="font-medium font-mono text-slate-700 dark:text-slate-300">
          {{ getRate(activeFunnelData.sbg, activeFunnelData.prospect) }}%
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
    hotprospect: number
    sbg: number
  }
  isFetching?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  isFetching: false
})

defineEmits<{
  (e: 'refresh'): void
}>()

const activeScope = ref<'monthly' | 'total'>('monthly')

const activeFunnelData = computed(() => {
  if (activeScope.value === 'total') {
    return {
      leads: props.stats?.leads || 0,
      prospect: props.stats?.prospect || 0,
      hotprospect: props.stats?.hotprospect || 0,
      sbg: props.stats?.sbg || 0
    }
  }
  return {
    leads: props.monthlyFunnel?.leads || 0,
    prospect: props.monthlyFunnel?.prospect || 0,
    hotprospect: props.monthlyFunnel?.hotprospect || 0,
    sbg: props.monthlyFunnel?.sbg || 0
  }
})

const formatNum = (val?: number) => {
  if (!val && val !== 0) return '0'
  return new Intl.NumberFormat('id-ID').format(val)
}

const getRate = (current?: number, total?: number) => {
  if (!total || total <= 0) return 0
  if (!current || current <= 0) return 0
  return Math.min(100, Math.round((current / total) * 100))
}
</script>
