<template>
  <div class="space-y-6">
    <!-- Review Sales Controls Card -->
    <div v-if="isReview && sortedSalesList.length > 0" class="bg-white border border-slate-100 rounded-2xl p-4 flex flex-wrap items-center justify-between gap-4 shadow-xs">
      <div class="flex items-center gap-2.5 font-medium text-xs text-slate-600">
        <span class="w-2 h-2 rounded-full bg-amber-500 animate-pulse"></span>
        <AppIcon name="user" class="w-4 h-4 text-amber-500" />
        <span>Pilih Sales Person:</span>
      </div>
      <div class="flex items-center gap-3 flex-wrap">
        <select
          :value="selectedSalesUser"
          @change="$emit('update:selectedSalesUser', ($event.target as HTMLSelectElement).value)"
          class="px-3.5 py-2 bg-slate-50 border border-slate-200 rounded-xl text-xs sm:text-sm font-medium text-slate-800 focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all min-w-[240px]"
        >
          <option v-for="s in sortedSalesList" :key="s.username" :value="s.username">
            {{ s.nama }} - {{ s.nm_jabatan || "SALES" }}
          </option>
        </select>
        <button
          @click="$emit('review-change')"
          :disabled="isFetching"
          class="inline-flex items-center gap-2 px-4 py-2 bg-amber-500 hover:bg-amber-600 active:scale-95 text-white font-medium text-xs sm:text-sm rounded-xl shadow-xs transition-all cursor-pointer disabled:opacity-50"
        >
          <AppIcon name="refresh" class="w-4 h-4" :class="{ 'animate-spin': isFetching }" />
          <span>Ganti Sales</span>
        </button>
        <button
          v-if="selectedSalesUser"
          @click="$emit('review-reset')"
          class="text-xs text-slate-500 hover:text-slate-800 font-medium cursor-pointer px-2.5 py-1.5 rounded-lg hover:bg-slate-100 transition-colors"
        >
          Reset
        </button>
      </div>
    </div>

    <!-- HERO & LEADERBOARD ROW -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Hero Card (2 cols) -->
      <div
        class="lg:col-span-2 bg-[#1e293b] text-white rounded-3xl p-6 md:p-8 shadow-xl relative overflow-hidden flex flex-col min-h-[220px] justify-center"
      >
        <!-- Background Decorative Glow -->
        <div class="absolute -top-24 -right-24 w-72 h-72 bg-amber-500/10 rounded-full blur-3xl pointer-events-none"></div>
        <div class="absolute -bottom-24 -left-24 w-72 h-72 bg-emerald-500/10 rounded-full blur-3xl pointer-events-none"></div>

        <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-6 relative z-10">
          <div class="space-y-3 flex-1">
            <div class="flex items-center gap-3">
              <span class="text-base sm:text-lg font-semibold text-amber-400">
                {{ isQualified ? "Estimasi Insentif Bulan Ini" : "Potensi Insentif Terkumpul" }}
              </span>

              <!-- Privacy Toggle Button -->
              <button
                @click="isPrivate = !isPrivate"
                type="button"
                :class="isPrivate ? 'bg-rose-500/20 text-rose-300 border-rose-500/30' : 'bg-white/10 text-slate-300 border-white/20 hover:bg-white/20'"
                class="p-1.5 rounded-lg border text-xs flex items-center gap-1.5 transition-all cursor-pointer"
                :title="isPrivate ? 'Tampilkan Nilai' : 'Sembunyikan Nilai'"
              >
                <AppIcon name="eye-off" v-if="isPrivate" class="w-4 h-4" />
                <AppIcon name="eye" v-else class="w-4 h-4" />
              </button>
            </div>

            <!-- Incentive Display -->
            <div class="flex items-baseline gap-2">
              <span v-if="!isPrivate" class="text-2xl text-slate-400 font-light">Rp</span>
              <h2 class="text-4xl md:text-5xl font-bold tracking-tight text-white select-none">
                <template v-if="isPrivate">******</template>
                <template v-else>{{ formatNum(estIncentive) }}</template>
              </h2>
            </div>

            <!-- Qualification Status Badge -->
            <div class="pt-1">
              <span v-if="isPrivate" class="text-sm font-normal text-slate-400">*****</span>
              <div
                v-else
                class="inline-flex items-center gap-2 px-3 py-1.5 rounded-xl text-xs font-medium"
                :class="isQualified ? 'bg-emerald-500/10 border border-emerald-500/20 text-emerald-400' : 'bg-slate-800/80 border border-slate-700 text-slate-300'"
              >
                <AppIcon name="check-circle" v-if="isQualified" class="w-4 h-4 text-emerald-400" />
                <AppIcon name="lock" v-else class="w-4 h-4 text-amber-400" />
                <span>Status: {{ isQualified ? "Qualified" : "Potensi Terkunci" }} | Grade: {{ userGrade }} (Rate: {{ isQualified ? appliedRate : (appliedRate || 25) }}%)</span>
              </div>
            </div>
          </div>

          <!-- Target Progress Card (Right Side inside Hero) -->
          <div class="w-full md:w-96 lg:w-[380px] shrink-0 bg-slate-800/50 backdrop-blur-md rounded-2xl p-5 border border-white/10 space-y-3.5">
            <div class="flex justify-between items-center text-xs">
              <span class="text-slate-400 font-medium">Target Progress</span>
              <span class="font-semibold text-white tabular-nums">{{ pointTotal }} / {{ targetGoalPoints }} Poin</span>
            </div>

            <!-- Progress Bar -->
            <div class="w-full bg-slate-700/60 rounded-full h-2.5 overflow-hidden p-0.5">
              <div
                class="bg-gradient-to-r from-amber-400 to-amber-500 h-full rounded-full transition-all duration-1000 ease-out"
                :style="{ width: `${Math.min(100, Math.max(0, targetProgressPct))}%` }"
              ></div>
            </div>

            <!-- Percentage Badge & Target Description -->
            <div class="space-y-1.5 pt-0.5">
              <div class="flex items-center gap-2">
                <span class="inline-flex items-center gap-1.5 py-0.5 rounded-full text-xs font-bold text-emerald-400 tabular-nums">
                  {{ targetProgressPct }}%
                </span>
              </div>
              <p class="text-xs text-slate-200 font-medium leading-relaxed">
                {{ targetDescription }}
              </p>
            </div>

            <!-- Motivation Footer -->
            <div class="pt-1 flex items-center gap-2 text-xs text-slate-300">
              <AppIcon name="bolt" class="w-4 h-4 text-amber-400 shrink-0" />
              <span class="leading-snug">{{ motivationText }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Leaderboard / Ranking Card (1 col) -->
      <div
        class="bg-white rounded-3xl p-6 border border-slate-100 shadow-xs flex flex-col items-center justify-between text-center relative group overflow-hidden"
      >
        <div class="absolute top-0 right-0 p-4 opacity-5 group-hover:scale-110 transition-transform duration-500">
          <AppIcon name="sparkles" class="w-24 h-24 text-slate-900" />
        </div>

        <div>
          <div
            class="w-14 h-14 rounded-2xl bg-amber-50 text-amber-500 flex items-center justify-center mx-auto mb-3 shadow-inner border border-amber-100"
          >
            <AppIcon name="sparkles" class="w-7 h-7" />
          </div>
          <h3 class="text-slate-400 text-xs font-semibold mb-1">
            Posisi Leaderboard
          </h3>
          <div class="text-3xl font-bold text-slate-800 tracking-tight mb-2">
            #{{ leaderboard.rank }}
            <span class="text-sm font-normal text-slate-500">dari {{ leaderboard.total_participants }} Sales</span>
          </div>
          <p class="text-xs text-slate-500 max-w-[240px] leading-relaxed mx-auto">
            <template v-if="leaderboard.rank === 1">
              <strong class="text-emerald-600 font-semibold">Luar biasa!</strong> Anda memimpin peringkat #1!
            </template>
            <template v-else>
              Butuh
              <strong class="text-emerald-600 font-semibold">Rp {{ formatNum(leaderboard.gap_to_next) }}</strong>
              omzet lagi untuk menyalip ke peringkat #{{ leaderboard.rank - 1 }}.
            </template>
          </p>
        </div>

        <button
          @click="$emit('open-leaderboard')"
          class="mt-4 px-4 py-2.5 text-xs font-semibold bg-amber-500 hover:bg-amber-600 active:scale-95 text-white rounded-xl shadow-xs transition-all flex items-center justify-center gap-2 w-full cursor-pointer"
        >
          <AppIcon name="chart" class="w-4 h-4" />
          <span>Lihat Leaderboard Lengkap</span>
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
  targetDescription: 'Mengejar Qualified 1 (Target: 25 Poin)',
  motivationText: 'Semangat! Kejar Poin untuk Bonus Insentif!',
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
