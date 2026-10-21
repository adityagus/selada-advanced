<template>
  <div class="bg-white rounded-3xl border border-slate-100 p-6 shadow-xs space-y-6">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2.5">
        <div class="w-7 h-7 rounded-lg bg-slate-100 flex items-center justify-center text-slate-500">
          <AppIcon name="user" class="w-4 h-4" />
        </div>
        <h3 class="font-semibold text-slate-800 text-lg">Sales Funnel</h3>
      </div>
      <button
        @click="$emit('refresh')"
        :disabled="isFetching"
        class="flex items-center gap-1.5 text-xs font-semibold px-3 py-1.5 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-xl transition-all disabled:opacity-50 cursor-pointer"
      >
        <AppIcon name="refresh" class="w-3.5 h-3.5" :class="{ 'animate-spin': isFetching }" />
        <span>{{ isFetching ? "Refetching..." : "Refresh Data" }}</span>
      </button>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- Column 1: Akumulasi Keseluruhan -->
      <div class="bg-slate-50/80 rounded-3xl p-6 border border-slate-100/90 space-y-5">
        <h4 class="text-sm font-bold text-slate-800 text-center">
          Akumulasi Keseluruhan
        </h4>

        <div class="flex">
          <div class="sales-funnel" style="width: 100%; max-width: none">
            <!-- Stage 1: Leads -->
            <div class="relative flex items-center justify-between gap-3">
              <div class="flex-1 flex justify-center">
                <div
                  class="bg-[#E50000] text-white text-xs font-bold py-2.5 px-4 rounded-t-md rounded-b-3xl text-center shadow-2xs transition-all duration-300"
                  :style="{ width: getFunnelWidth(1, stats.leads, stats.leads) }"
                >
                  Leads - {{ stats.leads }}
                </div>
              </div>
            </div>

            <!-- Stage 2: Prospek -->
            <div class="relative flex items-center justify-between gap-3">
              <div class="flex-1 flex justify-center">
                <div
                  class="bg-[#8B0000] text-white text-xs font-bold py-2.5 px-4 rounded-b-3xl text-center shadow-2xs transition-all duration-300"
                  :style="{ width: getFunnelWidth(2, stats.prospect, stats.leads) }"
                >
                  Prospek - {{ stats.prospect }}
                </div>
              </div>
            </div>

            <!-- Stage 3: Hot Prospek -->
            <div class="relative flex items-center justify-between gap-3">
              <div class="flex-1 flex justify-center">
                <div
                  class="bg-[#F59E0B] text-white text-xs font-bold py-2.5 px-4 rounded-b-3xl text-center shadow-2xs transition-all duration-300"
                  :style="{ width: getFunnelWidth(3, stats.hotprospect, stats.prospect || stats.leads) }"
                >
                  Hot Prospek - {{ stats.hotprospect }}
                </div>
              </div>
            </div>

            <!-- Stage 4: DO -->
            <div class="relative flex items-center justify-between gap-3">
              <div class="flex-1 flex justify-center">
                <div
                  class="bg-[#00C853] text-white text-xs font-bold py-2.5 px-4 rounded-b-3xl text-center shadow-2xs transition-all duration-300"
                  :style="{ width: getFunnelWidth(4, stats.sbg, stats.hotprospect || stats.leads) }"
                >
                  DO - {{ stats.sbg }}
                </div>
              </div>
            </div>
          </div>
          <div class="sales-percent flex flex-col justify-around">
            <div class="w-16 text-right font-semibold text-xs text-red-500 flex items-center justify-end gap-0.5 whitespace-nowrap">
              {{ calcDrop(stats.prospect, stats.leads) }}
            </div>
            <div class="w-16 text-right font-semibold text-xs text-amber-500 flex items-center justify-end gap-0.5 whitespace-nowrap">
              {{ calcDrop(stats.hotprospect, stats.prospect) }}
            </div>
            <div class="w-16 text-right font-semibold text-xs text-green-500 flex items-center justify-end gap-0.5 whitespace-nowrap">
              {{ calcDrop(stats.sbg, stats.hotprospect) }}
            </div>
          </div>
        </div>
      </div>

      <!-- Column 2: Bulan ini -->
      <div class="bg-slate-50/80 rounded-3xl p-6 border border-slate-100/90 space-y-5">
        <h4 class="text-sm font-bold text-slate-800 text-center">
          Bulan ini
        </h4>

        <div class="flex">
          <div class="sales-funnel" style="width: 100%; max-width: none">
            <!-- Stage 1: Leads -->
            <div class="relative flex items-center justify-between gap-3">
              <div class="flex-1 flex justify-center">
                <div
                  class="bg-[#E50000] text-white text-xs font-bold py-2.5 px-4 rounded-t-md rounded-b-3xl text-center shadow-2xs transition-all duration-300"
                  :style="{ width: getFunnelWidth(1, monthlyFunnel.leads, monthlyFunnel.leads) }"
                >
                  Leads - {{ monthlyFunnel.leads }}
                </div>
              </div>
            </div>

            <!-- Stage 2: Prospek -->
            <div class="relative flex items-center justify-between gap-3">
              <div class="flex-1 flex justify-center">
                <div
                  class="bg-[#8B0000] text-white text-xs font-bold py-2.5 px-4 rounded-b-3xl text-center shadow-2xs transition-all duration-300"
                  :style="{ width: getFunnelWidth(2, monthlyFunnel.prospect, monthlyFunnel.leads) }"
                >
                  Prospek - {{ monthlyFunnel.prospect }}
                </div>
              </div>
            </div>

            <!-- Stage 3: Hot Prospek -->
            <div class="relative flex items-center justify-between gap-3">
              <div class="flex-1 flex justify-center">
                <div
                  class="bg-[#F59E0B] text-white text-xs font-bold py-2.5 px-4 rounded-b-3xl text-center shadow-2xs transition-all duration-300"
                  :style="{ width: getFunnelWidth(3, monthlyFunnel.hotprospect, monthlyFunnel.prospect || monthlyFunnel.leads) }"
                >
                  Hot Prospek - {{ monthlyFunnel.hotprospect }}
                </div>
              </div>
            </div>

            <!-- Stage 4: DO -->
            <div class="relative flex items-center justify-between gap-3">
              <div class="flex-1 flex justify-center">
                <div
                  class="bg-[#00C853] text-white text-xs font-bold py-2.5 px-4 rounded-b-3xl text-center shadow-2xs transition-all duration-300"
                  :style="{ width: getFunnelWidth(4, monthlyFunnel.sbg, monthlyFunnel.hotprospect || monthlyFunnel.leads) }"
                >
                  DO - {{ monthlyFunnel.sbg }}
                </div>
              </div>
            </div>
          </div>
          <div class="sales-percent flex flex-col justify-around">
            <div class="w-16 text-right font-semibold text-xs text-red-500 flex items-center justify-end gap-0.5 whitespace-nowrap">
              {{ calcDrop(monthlyFunnel.prospect, monthlyFunnel.leads) }}
            </div>
            <div class="w-16 text-right font-semibold text-xs text-amber-500 flex items-center justify-end gap-0.5 whitespace-nowrap">
              {{ calcDrop(monthlyFunnel.hotprospect, monthlyFunnel.prospect) }}
            </div>
            <div class="w-16 text-right font-semibold text-xs text-green-500 flex items-center justify-end gap-0.5 whitespace-nowrap">
              {{ calcDrop(monthlyFunnel.sbg, monthlyFunnel.hotprospect) }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
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

withDefaults(defineProps<Props>(), {
  isFetching: false
})

defineEmits<{
  (e: 'refresh'): void
}>()

const calcDrop = (curr?: number, prev?: number) => {
  if (!prev || prev <= 0) return '0 %'
  if (!curr || curr >= prev) return '0 %'
  const drop = Math.round(100 - (((prev - curr) / prev) * 100))
  return `${drop} %`
}

const getFunnelWidth = (stage: number, count?: number, base?: number) => {
  const stageBases: Record<number, { min: number; default: number; max: number }> = {
    1: { min: 100, default: 100, max: 100 },
    2: { min: 72, default: 80, max: 85 },
    3: { min: 52, default: 60, max: 65 },
    4: { min: 32, default: 40, max: 45 }
  }

  const config = stageBases[stage] || { min: 40, default: 50, max: 60 }
  if (!base || base <= 0 || count === undefined || count === null) {
    return `${config.default}%`
  }

  const ratio = Math.min(1, Math.max(0, count / base))
  const calculatedWidth = config.min + ratio * (config.max - config.min)
  return `${Math.round(calculatedWidth)}%`
}
</script>

<style scoped>
.sales-funnel {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.sales-percent {
  display: flex;
  flex-direction: column;
  justify-content: space-around;
  padding-left: 0.5rem;
}
</style>
