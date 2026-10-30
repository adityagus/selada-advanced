<template>
  <div class="bg-white dark:bg-slate-900 rounded-3xl border border-slate-100 dark:border-slate-800 shadow-xs p-5 sm:p-6 space-y-4">
    <!-- Header with 3 dots & Badge (matching Instagram card in reference image) -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <h3 class="font-bold text-slate-800 dark:text-white text-base tracking-tight">
          Insentif Sales
        </h3>
        <span class="text-slate-350 dark:text-slate-600">···</span>
      </div>

      <div class="flex items-center gap-1.5">
        <span
          class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-semibold"
          :class="isQualified ? 'bg-emerald-50 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-400 border border-emerald-200/80' : 'bg-amber-50 text-amber-700 dark:bg-amber-950/40 dark:text-amber-400 border border-amber-200/80'"
        >
          <span class="w-1.5 h-1.5 rounded-full" :class="isQualified ? 'bg-emerald-500' : 'bg-amber-500'"></span>
          {{ isQualified ? "Qualified" : "Belum Qualified" }}
        </span>
      </div>
    </div>

    <!-- Review Sales Controls (if reviewer) -->
    <div v-if="isReview && sortedSalesList.length > 0" class="pt-1">
      <div class="flex items-center gap-2">
        <div class="relative flex-1">
          <select
            :value="selectedSalesUser"
            @change="$emit('update:selectedSalesUser', ($event.target as HTMLSelectElement).value)"
            class="w-full pl-3 pr-8 py-1.5 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium text-slate-800 dark:text-slate-200 focus:outline-none focus:border-amber-500 transition appearance-none cursor-pointer"
          >
            <option v-for="s in sortedSalesList" :key="s.username" :value="s.username">
              {{ s.nama }} — {{ s.nm_jabatan || "Sales" }}
            </option>
          </select>
          <AppIcon name="chevron-down" class="w-3.5 h-3.5 text-slate-400 absolute right-2.5 top-2.5 pointer-events-none" />
        </div>
        <button
          @click="$emit('review-change')"
          :disabled="isFetching"
          class="px-2.5 py-1.5 bg-amber-500 hover:bg-amber-600 text-white text-xs font-bold rounded-xl transition cursor-pointer disabled:opacity-50"
        >
          {{ isFetching ? '...' : 'Pilih' }}
        </button>
      </div>
    </div>

    <!-- Main Incentive Amount (matching Followers in reference image) -->
    <div class="space-y-1">
      <div class="flex items-center justify-between">
        <span class="text-xs text-slate-400 font-medium">
          {{ isQualified ? "Estimasi Insentif Bulan Ini" : "Potensi Insentif" }}
        </span>
        <button
          @click="isPrivate = !isPrivate"
          type="button"
          class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 p-1 rounded-lg hover:bg-slate-50 transition cursor-pointer"
          :title="isPrivate ? 'Tampilkan Nominal' : 'Sembunyikan Nominal'"
        >
          <AppIcon name="eye-off" v-if="isPrivate" class="w-4 h-4" />
          <AppIcon name="eye" v-else class="w-4 h-4" />
        </button>
      </div>

      <div class="flex items-baseline gap-2">
        <span v-if="!isPrivate" class="text-lg font-mono text-slate-400">Rp</span>
        <h2 class="text-3xl sm:text-4xl font-black tracking-tight text-slate-900 dark:text-white font-mono select-none">
          <template v-if="isPrivate">••••••••</template>
          <template v-else>{{ formatNum(estIncentive) }}</template>
        </h2>
      </div>
    </div>

    <!-- Target Poin Box with Circular Donut Gauge (matching Interactions row in reference image) -->
    <div class="p-4 rounded-2xl bg-slate-50/80 dark:bg-slate-800/40 border border-slate-100 dark:border-slate-800 flex items-center justify-between gap-4">
      <div class="space-y-1">
        <span class="text-xs font-medium text-slate-500 dark:text-slate-400">
          Target Poin Performa
        </span>
        <div class="text-xl font-bold text-slate-900 dark:text-white font-mono">
          {{ pointTotal }} <span class="text-xs font-normal text-slate-400">/ {{ targetGoalPoints }} Poin</span>
        </div>
        <p class="text-[11px] text-slate-500 dark:text-slate-400 line-clamp-1">
          {{ targetDescription }}
        </p>
      </div>

      <!-- Donut Circle Gauge (matching reference image) -->
      <div class="relative w-14 h-14 flex items-center justify-center shrink-0">
        <svg class="w-full h-full -rotate-90 transform" viewBox="0 0 60 60">
          <circle
            cx="30"
            cy="30"
            r="24"
            class="text-slate-200 dark:text-slate-700 stroke-current"
            stroke-width="5"
            fill="transparent"
          />
          <circle
            cx="30"
            cy="30"
            r="24"
            class="text-amber-500 stroke-current transition-all duration-1000 ease-out"
            stroke-width="5"
            stroke-linecap="round"
            fill="transparent"
            :stroke-dasharray="150.8"
            :stroke-dashoffset="150.8 - (150.8 * Math.min(100, Math.max(0, targetProgressPct))) / 100"
          />
        </svg>
        <span class="absolute text-[11px] font-bold text-slate-800 dark:text-slate-100 font-mono">
          {{ targetProgressPct }}%
        </span>
      </div>
    </div>

    <!-- Sub-metrics Grid (matching Likes/Comments/Saves/Shares in reference image) -->
    <div class="grid grid-cols-2 gap-2 text-xs pt-1 border-t border-slate-100 dark:border-slate-800">
      <div class="flex justify-between items-center py-1">
        <span class="text-slate-500 dark:text-slate-400">Gadai Emas</span>
        <span class="font-bold text-slate-800 dark:text-white font-mono">{{ gramasiEmas }}g</span>
      </div>

      <div class="flex justify-between items-center py-1">
        <span class="text-slate-500 dark:text-slate-400">Grade Sales</span>
        <span class="font-bold text-amber-600 dark:text-amber-400 font-mono">{{ userGrade }}</span>
      </div>

      <div class="flex justify-between items-center py-1">
        <span class="text-slate-500 dark:text-slate-400">Target Non-Emas</span>
        <span class="font-bold text-slate-800 dark:text-white font-mono">{{ nonEmasCard.pct_nonemas }}%</span>
      </div>

      <div class="flex justify-between items-center py-1">
        <span class="text-slate-500 dark:text-slate-400">Rate Komisi</span>
        <span class="font-bold text-emerald-600 dark:text-emerald-400 font-mono">{{ appliedRate }}%</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { NonEmasCardInfo } from '../types/dashboard.types'

interface Props {
  isReview?: boolean
  selectedSalesUser?: string
  sortedSalesList?: any[]
  isFetching?: boolean
  estIncentive?: number
  isQualified?: boolean
  userGrade?: string
  appliedRate?: number
  pointTotal?: number
  targetGoalPoints?: number
  targetProgressPct?: number
  targetDescription?: string
  gramasiEmas?: number
  nonEmasCard: NonEmasCardInfo
}

withDefaults(defineProps<Props>(), {
  isReview: false,
  selectedSalesUser: '',
  sortedSalesList: () => [],
  isFetching: false,
  estIncentive: 0,
  isQualified: false,
  userGrade: 'JUNIOR',
  appliedRate: 25,
  pointTotal: 0,
  targetGoalPoints: 25,
  targetProgressPct: 0,
  targetDescription: 'Target Poin Qualified',
  gramasiEmas: 0
})

defineEmits<{
  (e: 'update:selectedSalesUser', val: string): void
  (e: 'review-change'): void
  (e: 'review-reset'): void
}>()

const isPrivate = ref(true)

const formatNum = (val?: number) => {
  if (!val && val !== 0) return '0'
  return new Intl.NumberFormat('id-ID').format(val)
}
</script>
