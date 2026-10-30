<template>
  <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
    <!-- Card 1: Gadai Emas -->
    <div
      class="bg-white dark:bg-slate-900 p-5 rounded-xl border border-slate-200 dark:border-slate-800 flex flex-col justify-between space-y-4"
    >
      <div class="space-y-3">
        <div class="flex items-center justify-between">
          <span class="text-xs font-medium text-slate-500 dark:text-slate-400">
            Gadai Emas
          </span>
          <div class="w-8 h-8 rounded-lg bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 flex items-center justify-center">
            <AppIcon name="coin" class="w-4 h-4" />
          </div>
        </div>

        <div>
          <div class="text-2xl font-semibold text-slate-900 dark:text-slate-100 tracking-tight font-mono">
            Rp {{ formatNum(stats.emas) }}
          </div>
          <div class="flex items-center gap-1.5 mt-1 text-xs text-slate-500 dark:text-slate-400 font-medium">
            <span class="text-emerald-600 dark:text-emerald-400 font-semibold">{{ gramasiEmas }} gram</span>
            <span class="text-slate-300 dark:text-slate-700">•</span>
            <span>Akumulasi emas</span>
          </div>
        </div>
      </div>

      <div class="pt-3 border-t border-slate-100 dark:border-slate-800 space-y-1.5">
        <div class="flex items-center justify-between text-xs">
          <span class="text-slate-500 dark:text-slate-400">Target {{ emasCard.target_emas_g }}g</span>
          <span class="font-medium text-slate-700 dark:text-slate-300 font-mono">{{ emasCard.pct_emas }}%</span>
        </div>
        <div class="w-full bg-slate-100 dark:bg-slate-800 rounded-full h-1.5 overflow-hidden">
          <div
            class="h-full rounded-full transition-all duration-500"
            :class="emasCard.is_emas_bonus_done ? 'bg-emerald-500' : 'bg-amber-500'"
            :style="{ width: `${Math.min(100, Math.max(0, emasCard.pct_emas || 0))}%` }"
          ></div>
        </div>
        <p class="text-[11px] text-slate-400 truncate" :title="emasCard.bonus_emas_title">
          {{ emasCard.bonus_emas_title }}
        </p>
      </div>
    </div>

    <!-- Card 2: Gadai Non Emas / Elektronik -->
    <div
      class="bg-white dark:bg-slate-900 p-5 rounded-xl border border-slate-200 dark:border-slate-800 flex flex-col justify-between space-y-4"
    >
      <div class="space-y-3">
        <div class="flex items-center justify-between">
          <span class="text-xs font-medium text-slate-500 dark:text-slate-400">
            Gadai Elektronik
          </span>
          <div class="w-8 h-8 rounded-lg bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 flex items-center justify-center">
            <AppIcon name="monitor" class="w-4 h-4" />
          </div>
        </div>

        <div>
          <div class="text-2xl font-semibold text-slate-900 dark:text-slate-100 tracking-tight font-mono">
            Rp {{ formatNum(stats.nonemas) }}
          </div>
          <div class="flex items-center gap-1.5 mt-1 text-xs text-slate-500 dark:text-slate-400 font-medium">
            <span class="text-slate-700 dark:text-slate-300 font-semibold">Target UP: 50 Jt</span>
            <span class="text-slate-300 dark:text-slate-700">•</span>
            <span>Omzet aktif</span>
          </div>
        </div>
      </div>

      <div class="pt-3 border-t border-slate-100 dark:border-slate-800 space-y-1.5">
        <div class="flex items-center justify-between text-xs">
          <span class="text-slate-500 dark:text-slate-400">Pencapaian UP</span>
          <span class="font-medium text-slate-700 dark:text-slate-300 font-mono">{{ nonEmasCard.pct_nonemas }}%</span>
        </div>
        <div class="w-full bg-slate-100 dark:bg-slate-800 rounded-full h-1.5 overflow-hidden">
          <div
            class="h-full rounded-full transition-all duration-500"
            :class="nonEmasCard.is_nonemas_bonus_done ? 'bg-emerald-500' : 'bg-slate-700 dark:bg-slate-300'"
            :style="{ width: `${Math.min(100, Math.max(0, nonEmasCard.pct_nonemas || 0))}%` }"
          ></div>
        </div>
        <p class="text-[11px] text-slate-400 truncate" :title="nonEmasCard.bonus_title">
          {{ nonEmasCard.bonus_title }}
        </p>
      </div>
    </div>

    <!-- Card 3: New CIF (Nasabah Baru) -->
    <div
      class="bg-white dark:bg-slate-900 p-5 rounded-xl border border-slate-200 dark:border-slate-800 flex flex-col justify-between space-y-4"
    >
      <div class="space-y-3">
        <div class="flex items-center justify-between">
          <span class="text-xs font-medium text-slate-500 dark:text-slate-400">
            Nasabah Baru (New CIF)
          </span>
          <div class="w-8 h-8 rounded-lg bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 flex items-center justify-center">
            <AppIcon name="user-plus" class="w-4 h-4" />
          </div>
        </div>

        <div>
          <div class="text-2xl font-semibold text-slate-900 dark:text-slate-100 tracking-tight font-mono">
            {{ formatNum(stats.newcif) }}
            <span class="text-xs font-normal text-slate-400">CIF</span>
          </div>
          <div class="flex items-center gap-1.5 mt-1 text-xs text-slate-500 dark:text-slate-400 font-medium">
            <span>Akuisisi bulan berjalan</span>
          </div>
        </div>
      </div>

      <div class="pt-3 border-t border-slate-100 dark:border-slate-800 flex items-center justify-between text-xs text-slate-500 dark:text-slate-400">
        <span>Bulan Lalu</span>
        <span class="font-medium font-mono text-slate-700 dark:text-slate-300">{{ formatNum(lastMonthCif) }} CIF</span>
      </div>
    </div>

    <!-- Card 4: New Leads -->
    <div
      class="bg-white dark:bg-slate-900 p-5 rounded-xl border border-slate-200 dark:border-slate-800 flex flex-col justify-between space-y-4"
    >
      <div class="space-y-3">
        <div class="flex items-center justify-between">
          <span class="text-xs font-medium text-slate-500 dark:text-slate-400">
            Total Leads Masuk
          </span>
          <div class="w-8 h-8 rounded-lg bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 flex items-center justify-center">
            <AppIcon name="document" class="w-4 h-4" />
          </div>
        </div>

        <div>
          <div class="text-2xl font-semibold text-slate-900 dark:text-slate-100 tracking-tight font-mono">
            {{ formatNum(monthlyFunnel.leads) }}
            <span class="text-xs font-normal text-slate-400">Leads</span>
          </div>
          <div class="flex items-center gap-1.5 mt-1 text-xs text-slate-500 dark:text-slate-400 font-medium">
            <span>Inquiry aktif tercatat</span>
          </div>
        </div>
      </div>

      <div class="pt-3 border-t border-slate-100 dark:border-slate-800 flex items-center justify-between text-xs text-slate-500 dark:text-slate-400">
        <span>Bulan Lalu</span>
        <span class="font-medium font-mono text-slate-700 dark:text-slate-300">{{ formatNum(lastMonthLeads) }} Leads</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { EmasCardInfo, NonEmasCardInfo } from '../types/dashboard.types'

interface Props {
  stats: {
    emas?: number
    nonemas?: number
    newcif?: number
    [key: string]: any
  }
  emasCard: EmasCardInfo
  nonEmasCard: NonEmasCardInfo
  gramasiEmas?: number
  lastMonthCif?: number
  lastMonthLeads?: number
  monthlyFunnel: {
    leads: number
    [key: string]: any
  }
}

withDefaults(defineProps<Props>(), {
  gramasiEmas: 0,
  lastMonthCif: 0,
  lastMonthLeads: 0
})

const formatNum = (val?: number) => {
  if (!val && val !== 0) return '0'
  return new Intl.NumberFormat('id-ID').format(val)
}
</script>
