<template>
  <div class="bg-white dark:bg-slate-900 rounded-3xl border border-slate-100 dark:border-slate-800 shadow-xs p-5 sm:p-6 space-y-4">
    <!-- Header with 3 dots & Badge (matching X.com card in reference image) -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <h3 class="font-bold text-slate-800 dark:text-white text-base tracking-tight">
          Leaderboard Sales
        </h3>
        <span class="text-slate-350 dark:text-slate-600">···</span>
      </div>

      <span class="text-[11px] font-semibold text-slate-500 bg-slate-100 dark:bg-slate-800 px-2 py-0.5 rounded-full">
        Nasional
      </span>
    </div>

    <!-- Big Ranking Display (matching Followers in reference image) -->
    <div>
      <div class="text-xs text-slate-400 font-medium">Posisi Peringkat</div>
      <div class="flex items-baseline gap-2 mt-0.5">
        <span class="text-3xl font-black text-slate-900 dark:text-white font-mono tracking-tight">
          #{{ leaderboard.rank }}
        </span>
        <span class="text-xs font-semibold text-emerald-600 dark:text-emerald-400">
          dari {{ leaderboard.total_participants }} Sales
        </span>
      </div>
      <p class="text-[11px] text-slate-500 dark:text-slate-400 mt-1">
        <template v-if="leaderboard.rank === 1">
          Juara 1 Nasional! Pertahankan performa terbaik.
        </template>
        <template v-else>
          Butuh <strong class="text-amber-600 dark:text-amber-400 font-mono">Rp {{ formatCompact(leaderboard.gap_to_next) }}</strong> omzet lagi untuk menyalip ke #{{ leaderboard.rank - 1 }}
        </template>
      </p>
    </div>

    <!-- Capsule / Pill Bar Chart (matching vertical pill bars in reference image) -->
    <div class="p-3.5 rounded-2xl bg-slate-50/80 dark:bg-slate-800/40 border border-slate-100 dark:border-slate-800 space-y-2">
      <div class="flex items-center justify-between text-[11px] text-slate-400">
        <span>Distribusi Top Sales</span>
        <span class="font-mono">Poin Performa</span>
      </div>

      <!-- Rounded Vertical Pill Bars -->
      <div class="h-16 flex items-end justify-between gap-2 px-1 pt-2">
        <div
          v-for="(item, idx) in pillBars"
          :key="idx"
          class="flex-1 flex flex-col items-center gap-1 h-full justify-end group relative"
        >
          <!-- Tooltip on hover -->
          <div class="absolute bottom-full mb-1.5 opacity-0 group-hover:opacity-100 transition bg-slate-900 text-white text-[10px] py-1 px-2 rounded-md shadow-lg pointer-events-none whitespace-nowrap z-30">
            {{ item.name }}: {{ item.points }} pts
          </div>

          <div
            class="w-full max-w-[14px] rounded-full transition-all duration-500"
            :class="item.isMe ? 'bg-amber-500 shadow-xs' : 'bg-slate-300 dark:bg-slate-700'"
            :style="{ height: `${item.pct}%` }"
          ></div>
          <span class="text-[9px] font-mono text-slate-400" :class="{ 'font-bold text-amber-600': item.isMe }">
            {{ item.label }}
          </span>
        </div>
      </div>
    </div>

    <!-- Compact Live Feed Activity Stream -->
    <div class="space-y-2 pt-1">
      <div class="flex items-center justify-between text-xs">
        <span class="font-bold text-slate-700 dark:text-slate-300">Aktivitas Insentif Terbaru</span>
        <span class="text-[10px] text-slate-400 font-mono">{{ liveFeeds.length }} Transaksi</span>
      </div>

      <div class="space-y-1.5 max-h-36 overflow-y-auto pr-0.5">
        <div
          v-if="liveFeeds.length === 0"
          class="py-4 text-center text-slate-400 text-xs"
        >
          Belum ada transaksi hari ini
        </div>

        <div
          v-for="(feed, idx) in liveFeeds.slice(0, 3)"
          :key="idx"
          class="flex items-center justify-between p-2 rounded-xl bg-slate-50 dark:bg-slate-800/40 text-xs border border-slate-100 dark:border-slate-800"
        >
          <div class="flex items-center gap-2 min-w-0">
            <div class="w-6 h-6 rounded-lg bg-amber-50 text-amber-600 flex items-center justify-center shrink-0">
              <AppIcon name="coin" class="w-3.5 h-3.5" />
            </div>
            <div class="min-w-0">
              <p class="font-bold text-slate-800 dark:text-white truncate text-[11px]">{{ feed.type }}</p>
              <p class="text-[9px] text-slate-400 truncate">{{ feed.code }}</p>
            </div>
          </div>
          <div class="text-right shrink-0 font-mono text-xs font-bold text-emerald-600 dark:text-emerald-400">
            +Rp {{ formatCompact(feed.amount) }}
          </div>
        </div>
      </div>
    </div>

    <!-- Full Leaderboard CTA Button -->
    <button
      @click="$emit('open-leaderboard')"
      type="button"
      class="w-full py-2.5 px-4 text-xs font-bold text-white bg-amber-500 hover:bg-amber-600 active:scale-95 rounded-xl shadow-xs transition cursor-pointer flex items-center justify-center gap-2"
    >
      <AppIcon name="chart" class="w-4 h-4 text-white" />
      <span>Lihat Leaderboard Lengkap</span>
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { LeaderboardInfo, LiveFeedItem } from '../types/dashboard.types'

interface Props {
  leaderboard?: LeaderboardInfo
  liveFeeds?: LiveFeedItem[]
}

const props = withDefaults(defineProps<Props>(), {
  leaderboard: () => ({ rank: 1, total_participants: 1, gap_to_next: 0, top_performers: [] }),
  liveFeeds: () => []
})

defineEmits<{
  (e: 'open-leaderboard'): void
}>()

const pillBars = computed(() => {
  const top = props.leaderboard?.top_performers || []
  if (top.length === 0) {
    return [
      { label: '#1', name: 'Top 1', points: 30, pct: 100, isMe: props.leaderboard.rank === 1 },
      { label: '#2', name: 'Top 2', points: 24, pct: 80, isMe: props.leaderboard.rank === 2 },
      { label: '#3', name: 'Top 3', points: 18, pct: 60, isMe: props.leaderboard.rank === 3 },
      { label: '#4', name: 'Top 4', points: 12, pct: 40, isMe: props.leaderboard.rank === 4 },
      { label: '#5', name: 'Top 5', points: 8, pct: 25, isMe: props.leaderboard.rank >= 5 }
    ]
  }

  const maxPoints = Math.max(...top.map(t => t.points || 1), 1)
  return top.slice(0, 6).map((t, idx) => ({
    label: `#${t.rank || idx + 1}`,
    name: t.name,
    points: t.points,
    pct: Math.min(100, Math.max(15, Math.round(((t.points || 0) / maxPoints) * 100))),
    isMe: t.rank === props.leaderboard.rank
  }))
})

const formatCompact = (val?: number) => {
  if (!val) return '0'
  if (val >= 1000000000) return (val / 1000000000).toFixed(1) + ' M'
  if (val >= 1000000) return (val / 1000000).toFixed(1) + ' Jt'
  if (val >= 1000) return (val / 1000).toFixed(0) + ' Rb'
  return String(val)
}
</script>
