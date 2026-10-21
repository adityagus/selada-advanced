<template>
  <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
    <!-- Card 1: Gadai Emas -->
    <div class="bg-white p-5 rounded-3xl border border-slate-100 hover:border-slate-200/80 shadow-xs hover:shadow-sm transition-all duration-300 flex flex-col justify-between group">
      <div>
        <div class="flex items-center justify-between mb-4">
          <h4 class="font-bold text-slate-800 text-sm md:text-base tracking-tight">
            Emas
          </h4>
          <div class="w-10 h-10 rounded-2xl bg-amber-50 text-amber-600 border border-amber-100 flex items-center justify-center flex-shrink-0 shadow-2xs group-hover:scale-105 transition-transform">
            <AppIcon name="coin" class="w-5 h-5" />
          </div>
        </div>

        <div class="flex items-baseline justify-between gap-2 mb-3 flex-wrap">
          <span class="text-xl lg:text-2xl font-bold text-slate-900 tracking-tight">
            Rp {{ formatNum(stats.emas) }}
          </span>
          <span class="inline-flex items-center gap-1 bg-emerald-50 text-emerald-600 border border-emerald-200/60 font-normal text-xs px-2.5 py-1 rounded-full">
            <AppIcon name="arrow-up" class="w-3 h-3 text-emerald-600" />
            {{ gramasiEmas }}g
          </span>
        </div>
      </div>

      <div class="pt-2 space-y-1.5">
        <div class="flex items-center justify-between text-xs">
          <span class="font-medium text-slate-500">Target {{ emasCard.target_emas_g }}g</span>
          <span class="font-bold text-amber-600">{{ emasCard.pct_emas }}%</span>
        </div>
        <div class="w-full bg-slate-100 rounded-full h-2 overflow-hidden">
          <div
            class="h-full rounded-full transition-all duration-1000"
            :class="emasCard.is_emas_bonus_done ? 'bg-emerald-500' : 'bg-amber-500'"
            :style="{ width: `${emasCard.pct_emas}%` }"
          ></div>
        </div>
        <p class="text-[11px] text-slate-400 font-medium truncate">{{ emasCard.bonus_emas_title }}</p>
      </div>
    </div>

    <!-- Card 2: Gadai Non Emas / Elektronik -->
    <div class="bg-white p-5 rounded-3xl border border-slate-100 hover:border-slate-200/80 shadow-xs hover:shadow-sm transition-all duration-300 flex flex-col justify-between group">
      <div>
        <div class="flex items-center justify-between mb-4">
          <h4 class="font-bold text-slate-800 text-sm md:text-base tracking-tight">
            Elektronik
          </h4>
          <div class="w-10 h-10 rounded-2xl bg-amber-50 text-amber-600 border border-amber-100 flex items-center justify-center flex-shrink-0 shadow-2xs group-hover:scale-105 transition-transform">
            <AppIcon name="monitor" class="w-5 h-5" />
          </div>
        </div>

        <div class="flex items-baseline justify-between gap-2 mb-3 flex-wrap">
          <span class="text-xl lg:text-2xl font-bold text-slate-900 tracking-tight">
            Rp {{ formatNum(stats.nonemas) }}
          </span>
          <span class="inline-flex items-center gap-1 bg-amber-50 text-amber-700 border border-amber-200/60 font-medium text-xs px-2.5 py-1 rounded-full">
            Target 50M
          </span>
        </div>
      </div>

      <div class="pt-2 space-y-1.5">
        <div class="flex items-center justify-between text-xs">
          <span class="font-medium text-slate-500">Pencapaian UP</span>
          <span class="font-bold text-amber-600">{{ nonEmasCard.pct_nonemas }}%</span>
        </div>
        <div class="w-full bg-slate-100 rounded-full h-2 overflow-hidden">
          <div
            class="h-full rounded-full transition-all duration-1000"
            :class="nonEmasCard.is_nonemas_bonus_done ? 'bg-emerald-500' : 'bg-amber-500'"
            :style="{ width: `${nonEmasCard.pct_nonemas}%` }"
          ></div>
        </div>
        <p class="text-[11px] text-slate-400 font-medium truncate">{{ nonEmasCard.bonus_title }}</p>
      </div>
    </div>

    <!-- Card 3: New CIF (Nasabah Baru) -->
    <div class="bg-white p-5 rounded-3xl border border-slate-100 hover:border-slate-200/80 shadow-xs hover:shadow-sm transition-all duration-300 flex flex-col justify-between group">
      <div>
        <div class="flex items-center justify-between mb-4">
          <h4 class="font-bold text-slate-800 text-sm md:text-base tracking-tight">
            New CIF (Nasabah Baru)
          </h4>
          <div class="w-10 h-10 rounded-2xl bg-slate-900 text-amber-400 border border-slate-800 flex items-center justify-center flex-shrink-0 shadow-2xs group-hover:scale-105 transition-transform">
            <AppIcon name="user-plus" class="w-5 h-5" />
          </div>
        </div>

        <div class="flex items-baseline justify-between gap-2 mb-3">
          <span class="text-2xl lg:text-3xl font-extrabold text-slate-900 tracking-tight">
            {{ formatNum(stats.newcif) }}
          </span>
          <span class="inline-flex items-center gap-1 bg-amber-50 text-amber-700 border border-amber-200/60 font-bold text-xs px-2.5 py-1 rounded-full">
            <AppIcon name="arrow-up" class="w-3 h-3 text-amber-600" />
            Active
          </span>
        </div>
      </div>

      <div class="pt-2 flex items-center justify-between text-xs text-slate-400 font-medium">
        <span>vs. <span class="font-semibold text-slate-500"> {{ formatNum(lastMonthCif) }} Bulan Lalu</span></span>
      </div>
    </div>

    <!-- Card 4: New Leads -->
    <div class="bg-white p-5 rounded-3xl border border-slate-100 hover:border-slate-200/80 shadow-xs hover:shadow-sm transition-all duration-300 flex flex-col justify-between group">
      <div>
        <div class="flex items-center justify-between mb-4">
          <h4 class="font-bold text-slate-800 text-sm md:text-base tracking-tight">
            New Leads
          </h4>
          <div class="w-10 h-10 rounded-2xl bg-amber-50 text-amber-600 border border-amber-100 flex items-center justify-center flex-shrink-0 shadow-2xs group-hover:scale-105 transition-transform">
            <AppIcon name="document" class="w-5 h-5" />
          </div>
        </div>

        <div class="flex items-baseline justify-between gap-2 mb-3">
          <span class="text-2xl lg:text-3xl font-extrabold text-slate-900 tracking-tight">
            {{ formatNum(monthlyFunnel.leads) }}
          </span>
          <span class="inline-flex items-center gap-1 bg-amber-50 text-amber-700 border border-amber-200/60 font-bold text-xs px-2.5 py-1 rounded-full">
            <AppIcon name="arrow-up" class="w-3 h-3 text-amber-600" />
            Leads
          </span>
        </div>
      </div>

      <div class="pt-2 flex items-center justify-between text-xs text-slate-400 font-medium">
        <span>vs. <span class="font-semibold text-slate-500"> {{ formatNum(lastMonthLeads) }} Bulan Lalu</span></span>
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
