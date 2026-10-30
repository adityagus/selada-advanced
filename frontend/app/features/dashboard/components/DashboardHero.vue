<template>
  <div class="space-y-4">
    <!-- Review Sales Controls (When Active) -->
    <div
      v-if="isReview && sortedSalesList.length > 0"
      class="bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl px-4 py-3 flex flex-wrap items-center justify-between gap-3"
    >
      <div class="flex items-center gap-2 text-xs text-slate-600 dark:text-slate-400 font-medium">
        <span class="w-2 h-2 rounded-full bg-amber-500"></span>
        <span>Review Performa Sales:</span>
      </div>

      <div class="flex items-center gap-2 flex-wrap">
        <div class="relative">
          <select
            :value="selectedSalesUser"
            @change="$emit('update:selectedSalesUser', ($event.target as HTMLSelectElement).value)"
            class="pl-3 pr-8 py-1.5 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-xs font-medium text-slate-800 dark:text-slate-200 focus:outline-none focus:ring-1 focus:ring-slate-400 dark:focus:ring-slate-600 transition min-w-[220px] cursor-pointer appearance-none"
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
          class="px-3 py-1.5 bg-amber-500 hover:bg-amber-600 active:scale-95 text-white text-xs font-semibold rounded-lg shadow-2xs transition disabled:opacity-50 cursor-pointer"
        >
          {{ isFetching ? 'Memuat...' : 'Tampilkan' }}
        </button>

        <button
          v-if="selectedSalesUser"
          @click="$emit('review-reset')"
          class="text-xs text-slate-500 hover:text-amber-700 dark:hover:text-amber-400 font-medium cursor-pointer px-2.5 py-1.5 rounded-lg hover:bg-amber-50 dark:hover:bg-amber-950/40 transition"
        >
          Reset ke Saya
        </button>
      </div>
    </div>

    <!-- Main Overview Row -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-4">
      <!-- Left: Incentive & Performance Status (8 cols) -->
      <div
        class="lg:col-span-8 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-5 sm:p-6 flex flex-col justify-between space-y-6"
      >
        <!-- Top Row: Label, Status, Privacy -->
        <div class="flex items-start justify-between gap-4">
          <div class="space-y-1">
            <div class="flex items-center gap-2">
              <span class="text-xs font-medium text-slate-500 dark:text-slate-400">
                {{ isQualified ? "Estimasi Insentif Berjalan" : "Potensi Insentif" }}
              </span>
              <span
                class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded text-[11px] font-medium"
                :class="isQualified ? 'bg-emerald-50 dark:bg-emerald-950/40 text-emerald-700 dark:text-emerald-400 border border-emerald-200 dark:border-emerald-800' : 'bg-amber-50 dark:bg-amber-950/40 text-amber-800 dark:text-amber-300 border border-amber-200 dark:border-amber-800'"
              >
                <span class="w-1.5 h-1.5 rounded-full" :class="isQualified ? 'bg-emerald-500' : 'bg-amber-500'"></span>
                {{ isQualified ? 'Qualified' : 'Belum Qualified' }}
              </span>
            </div>

            <!-- Big Metric Number -->
            <div class="flex items-baseline gap-2 pt-1">
              <span v-if="!isPrivate" class="text-lg font-medium text-slate-400 dark:text-slate-500">Rp</span>
              <span class="text-3xl sm:text-4xl font-semibold tracking-tight text-slate-900 dark:text-slate-100 font-mono">
                <template v-if="isPrivate">••••••••</template>
                <template v-else>{{ formatNum(estIncentive) }}</template>
              </span>
            </div>
          </div>

          <!-- Privacy Toggle Button -->
          <button
            @click="isPrivate = !isPrivate"
            type="button"
            class="p-2 text-slate-400 hover:text-amber-600 dark:hover:text-amber-400 rounded-lg border border-slate-200 dark:border-slate-800 hover:bg-amber-50/50 dark:hover:bg-slate-800/60 transition cursor-pointer"
            :title="isPrivate ? 'Tampilkan Nominal' : 'Sembunyikan Nominal'"
          >
            <AppIcon name="eye-off" v-if="isPrivate" class="w-4 h-4" />
            <AppIcon name="eye" v-else class="w-4 h-4" />
          </button>
        </div>

        <!-- Middle: Target Progress (Linear Style with Sidebar Amber Color) -->
        <div class="pt-4 border-t border-slate-100 dark:border-slate-800 space-y-2.5">
          <div class="flex items-center justify-between text-xs">
            <div class="flex items-center gap-2">
              <span class="font-medium text-slate-700 dark:text-slate-300">Target Poin</span>
              <span class="text-slate-400">•</span>
              <span class="text-slate-500 dark:text-slate-400">{{ targetDescription }}</span>
            </div>
            <div class="font-mono text-slate-800 dark:text-slate-200 font-medium">
              <span class="text-amber-600 dark:text-amber-400 font-bold">{{ pointTotal }}</span>
              <span class="text-slate-400"> / {{ targetGoalPoints }} pts</span>
              <span class="ml-2 text-slate-500">({{ targetProgressPct }}%)</span>
            </div>
          </div>

          <!-- Clean Minimal Track in Sidebar Amber Primary -->
          <div class="w-full bg-slate-100 dark:bg-slate-800 rounded-full h-2 overflow-hidden">
            <div
              class="h-full rounded-full transition-all duration-500"
              :class="targetProgressPct >= 100 ? 'bg-emerald-500' : 'bg-amber-500'"
              :style="{ width: `${Math.min(100, Math.max(0, targetProgressPct))}%` }"
            ></div>
          </div>

          <div class="flex items-center justify-between text-[11px] text-slate-400 dark:text-slate-500">
            <span>Sisa kebutuhan: {{ Math.max(0, (targetGoalPoints || 25) - (pointTotal || 0)) }} Poin</span>
            <div class="flex items-center gap-3">
              <span>Grade: <strong class="text-slate-700 dark:text-slate-300 font-semibold">{{ userGrade }}</strong></span>
              <span>Rate: <strong class="text-slate-700 dark:text-slate-300 font-semibold">{{ appliedRate }}%</strong></span>
            </div>
          </div>
        </div>
      </div>

      <!-- Right: Leaderboard Position Summary (4 cols) -->
      <div
        class="lg:col-span-4 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-5 sm:p-6 flex flex-col justify-between space-y-4"
      >
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <span class="text-xs font-medium text-slate-500 dark:text-slate-400">Peringkat Sales</span>
            <span class="text-[11px] text-slate-400">Nasional</span>
          </div>

          <div class="flex items-baseline gap-2">
            <span class="text-3xl font-semibold tracking-tight text-slate-900 dark:text-slate-100 font-mono">
              #{{ leaderboard.rank }}
            </span>
            <span class="text-xs text-slate-400">
              dari {{ leaderboard.total_participants }} sales
            </span>
          </div>

          <p class="text-xs text-slate-500 dark:text-slate-400 leading-relaxed">
            <template v-if="leaderboard.rank === 1">
              Peringkat tertinggi nasional saat ini.
            </template>
            <template v-else>
              Gap ke peringkat #{{ leaderboard.rank - 1 }}:
              <strong class="font-mono text-slate-800 dark:text-slate-200 font-semibold">
                Rp {{ formatNum(leaderboard.gap_to_next) }}
              </strong>
            </template>
          </p>
        </div>

        <button
          @click="$emit('open-leaderboard')"
          type="button"
          class="w-full py-2 px-3 text-xs font-semibold text-amber-700 dark:text-amber-300 bg-amber-50 dark:bg-amber-950/40 hover:bg-amber-100 dark:hover:bg-amber-900/40 border border-amber-200/80 dark:border-amber-800/80 rounded-lg transition flex items-center justify-center gap-1.5 cursor-pointer"
        >
          <AppIcon name="chart" class="w-3.5 h-3.5 text-amber-600 dark:text-amber-400" />
          <span>Buka Leaderboard</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { LeaderboardInfo } from '../types/dashboard.types'

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
  motivationText?: string
  leaderboard?: LeaderboardInfo
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
  motivationText: '',
  leaderboard: () => ({ rank: 1, total_participants: 1, gap_to_next: 0, top_performers: [] })
})

defineEmits<{
  (e: 'update:selectedSalesUser', val: string): void
  (e: 'review-change'): void
  (e: 'review-reset'): void
  (e: 'open-leaderboard'): void
}>()

const isPrivate = ref(true)

const formatNum = (val?: number) => {
  if (!val && val !== 0) return '0'
  return new Intl.NumberFormat('id-ID').format(val)
}
</script>
