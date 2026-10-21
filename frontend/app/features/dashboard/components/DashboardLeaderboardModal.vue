<template>
  <Teleport to="body">
    <div
      v-if="show"
      class="fixed inset-0 z-[999] bg-slate-900/60 backdrop-blur-sm flex items-center justify-center p-3 sm:p-5 md:p-6 animate-in fade-in duration-200"
      @click.self="$emit('update:show', false)"
    >
      <div
        class="bg-white rounded-3xl max-w-4xl xl:max-w-5xl w-full max-h-[92vh] flex flex-col shadow-2xl border border-slate-100 overflow-hidden animate-in zoom-in-95 duration-200"
      >
        <!-- Modal Header -->
        <div class="px-6 py-5 border-b border-slate-100 flex items-center justify-between bg-gradient-to-r from-amber-50/50 via-white to-slate-50/50">
          <div class="flex items-center gap-3.5">
            <div class="w-11 h-11 rounded-2xl bg-gradient-to-tr from-amber-500 to-amber-400 text-white flex items-center justify-center font-bold shadow-md shadow-amber-500/20">
              <AppIcon name="sparkles" class="w-6 h-6 text-white" />
            </div>
            <div>
              <div class="flex items-center gap-2">
                <h3 class="font-bold text-slate-800 text-lg md:text-xl">Leaderboard Performa Sales</h3>
                <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-semibold bg-amber-100 text-amber-800 border border-amber-200">
                  <AppIcon name="chart" class="w-3 h-3 text-amber-600" />
                  Bulan Berjalan
                </span>
              </div>
              <p class="text-xs text-slate-500 mt-0.5">Peringkat performa pencairan omzet & akumulasi poin per sales</p>
            </div>
          </div>
          <button
            @click="$emit('update:show', false)"
            class="p-2 text-slate-400 hover:text-slate-700 hover:bg-slate-100 rounded-xl transition-colors cursor-pointer"
            title="Tutup (Esc)"
          >
            <AppIcon name="close" class="w-5 h-5" />
          </button>
        </div>

        <!-- Scrollable Modal Body -->
        <div class="p-6 overflow-y-auto space-y-5">
          <!-- Top 3 Podium Highlights -->
          <div v-if="topThreePerformers.length > 0" class="grid grid-cols-1 md:grid-cols-3 gap-3.5">
            <!-- Rank 1: Gold -->
            <div
              v-if="topThreePerformers[0]"
              class="relative overflow-hidden rounded-2xl p-4 bg-gradient-to-b from-amber-50 via-amber-50/40 to-white border-2 border-amber-300 shadow-xs flex flex-col justify-between"
            >
              <div class="absolute -right-3 -top-3 w-16 h-16 bg-amber-200/40 rounded-full blur-xl pointer-events-none"></div>
              <div class="flex items-start justify-between gap-2 mb-2">
                <div class="flex items-center gap-2">
                  <span class="w-8 h-8 rounded-xl bg-gradient-to-tr from-amber-500 to-yellow-400 text-white font-extrabold text-sm flex items-center justify-center shadow-sm">
                    1
                  </span>
                  <div class="min-w-0">
                    <span class="text-[10px] font-bold text-amber-700 tracking-wider uppercase">Juara 1</span>
                    <h4 class="font-bold text-slate-900 text-sm truncate" :title="topThreePerformers[0].name">
                      {{ topThreePerformers[0].name }}
                    </h4>
                  </div>
                </div>
                <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-[11px] font-bold bg-amber-100 text-amber-900 border border-amber-200 shrink-0">
                  <AppIcon name="star" class="w-3 h-3 text-amber-600" />
                  {{ topThreePerformers[0].points }} pts
                </span>
              </div>
              <div class="mt-2 pt-1 flex items-baseline justify-between text-xs">
                <span class="text-slate-500 text-[11px]">Total Omzet:</span>
                <span class="font-mono font-bold text-amber-900">
                  <span class="text-[10px] text-amber-600 font-sans mr-0.5">Rp</span>{{ formatNum(topThreePerformers[0].total_omzet) }}
                </span>
              </div>
            </div>

            <!-- Rank 2: Silver -->
            <div
              v-if="topThreePerformers[1]"
              class="relative overflow-hidden rounded-2xl p-4 bg-gradient-to-b from-slate-50 via-slate-50/50 to-white border border-slate-200 shadow-xs flex flex-col justify-between"
            >
              <div class="flex items-start justify-between gap-2 mb-2">
                <div class="flex items-center gap-2">
                  <span class="w-8 h-8 rounded-xl bg-gradient-to-tr from-slate-400 to-slate-300 text-white font-extrabold text-sm flex items-center justify-center shadow-sm">
                    2
                  </span>
                  <div class="min-w-0">
                    <span class="text-[10px] font-bold text-slate-500 tracking-wider uppercase">Peringkat 2</span>
                    <h4 class="font-bold text-slate-800 text-sm truncate" :title="topThreePerformers[1].name">
                      {{ topThreePerformers[1].name }}
                    </h4>
                  </div>
                </div>
                <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-[11px] font-semibold bg-slate-100 text-slate-700 border border-slate-200 shrink-0">
                  {{ topThreePerformers[1].points }} pts
                </span>
              </div>
              <div class="mt-2 pt-1 flex items-baseline justify-between text-xs">
                <span class="text-slate-500 text-[11px]">Total Omzet:</span>
                <span class="font-mono font-bold text-slate-800">
                  <span class="text-[10px] text-slate-400 font-sans mr-0.5">Rp</span>{{ formatNum(topThreePerformers[1].total_omzet) }}
                </span>
              </div>
            </div>

            <!-- Rank 3: Bronze -->
            <div
              v-if="topThreePerformers[2]"
              class="relative overflow-hidden rounded-2xl p-4 bg-gradient-to-b from-orange-50/60 via-orange-50/20 to-white border border-orange-200 shadow-xs flex flex-col justify-between"
            >
              <div class="flex items-start justify-between gap-2 mb-2">
                <div class="flex items-center gap-2">
                  <span class="w-8 h-8 rounded-xl bg-gradient-to-tr from-amber-700 to-orange-500 text-white font-extrabold text-sm flex items-center justify-center shadow-sm">
                    3
                  </span>
                  <div class="min-w-0">
                    <span class="text-[10px] font-bold text-amber-800 tracking-wider uppercase">Peringkat 3</span>
                    <h4 class="font-bold text-slate-800 text-sm truncate" :title="topThreePerformers[2].name">
                      {{ topThreePerformers[2].name }}
                    </h4>
                  </div>
                </div>
                <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-[11px] font-semibold bg-orange-100 text-orange-900 border border-orange-200 shrink-0">
                  {{ topThreePerformers[2].points }} pts
                </span>
              </div>
              <div class="mt-2 pt-1 flex items-baseline justify-between text-xs">
                <span class="text-slate-500 text-[11px]">Total Omzet:</span>
                <span class="font-mono font-bold text-orange-950">
                  <span class="text-[10px] text-orange-700 font-sans mr-0.5">Rp</span>{{ formatNum(topThreePerformers[2].total_omzet) }}
                </span>
              </div>
            </div>
          </div>

          <!-- Controls Bar: Filter Pills & Search -->
          <div class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3 pt-1">
            <!-- Filter Tabs -->
            <div class="flex items-center gap-1.5 p-1 bg-slate-100 rounded-xl">
              <button
                v-for="f in filterOptions"
                :key="f.key"
                @click="leaderboardActiveFilter = f.key"
                class="px-3 py-1.5 rounded-lg text-xs font-semibold transition-all cursor-pointer"
                :class="leaderboardActiveFilter === f.key ? 'bg-white text-slate-900 shadow-xs' : 'text-slate-500 hover:text-slate-800'"
              >
                {{ f.label }}
              </button>
            </div>

            <!-- Search Field -->
            <div class="relative w-full sm:w-64">
              <AppIcon name="search" class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
              <input
                v-model="leaderboardSearch"
                type="text"
                placeholder="Cari nama sales..."
                class="w-full pl-9 pr-8 py-2 bg-slate-50 border border-slate-200 rounded-xl text-xs font-medium text-slate-800 focus:outline-none focus:border-amber-500 focus:bg-white transition-all"
              />
              <button
                v-if="leaderboardSearch"
                @click="leaderboardSearch = ''"
                class="absolute right-2.5 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 p-0.5 cursor-pointer"
              >
                <AppIcon name="close" class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>

          <!-- Table Container -->
          <div class="border border-slate-100 rounded-2xl overflow-hidden shadow-2xs bg-white">
            <table class="w-full text-xs text-left">
              <thead>
                <tr class="bg-slate-50 border-b border-slate-100 text-slate-500 uppercase tracking-wider font-semibold text-[11px]">
                  <th class="w-14 px-3.5 py-3 text-center">Rank</th>
                  <th class="px-4 py-3">Nama Sales</th>
                  <th class="px-3.5 py-3 text-right">Gadai Emas</th>
                  <th class="px-3.5 py-3 text-right">Non Emas</th>
                  <th class="px-4 py-3 text-right">Total Omzet</th>
                  <th class="w-28 px-3.5 py-3 text-right">Poin</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr
                  v-for="p in filteredPerformers"
                  :key="p.employee_id || p.name"
                  class="hover:bg-slate-50/80 transition-colors"
                  :class="isCurrentUser(p) ? 'bg-amber-50/50 font-medium' : ''"
                >
                  <!-- Rank Column -->
                  <td class="w-14 px-3.5 py-3.5 text-center font-bold">
                    <span
                      v-if="p.rank === 1"
                      class="inline-flex items-center justify-center w-6 h-6 rounded-lg bg-amber-100 text-amber-800 font-extrabold text-[11px]"
                    >
                      1
                    </span>
                    <span
                      v-else-if="p.rank === 2"
                      class="inline-flex items-center justify-center w-6 h-6 rounded-lg bg-slate-200 text-slate-700 font-extrabold text-[11px]"
                    >
                      2
                    </span>
                    <span
                      v-else-if="p.rank === 3"
                      class="inline-flex items-center justify-center w-6 h-6 rounded-lg bg-orange-100 text-orange-800 font-extrabold text-[11px]"
                    >
                      3
                    </span>
                    <span v-else class="text-slate-400 font-mono text-xs">
                      #{{ p.rank }}
                    </span>
                  </td>

                  <!-- Nama Sales Column -->
                  <td class="min-w-[200px] px-4 py-3.5">
                    <div class="flex items-center gap-2.5">
                      <div
                        :class="[
                          p.rank === 1 ? 'bg-amber-100 text-amber-800 border-amber-200' :
                          p.rank === 2 ? 'bg-slate-200 text-slate-700 border-slate-300' :
                          p.rank === 3 ? 'bg-orange-100 text-orange-800 border-orange-200' :
                          'bg-slate-100 text-slate-600 border-slate-200'
                        ]"
                        class="w-7 h-7 rounded-full border flex items-center justify-center text-[10px] font-bold shrink-0 uppercase"
                      >
                        {{ getInitials(p.name) }}
                      </div>
                      <div class="min-w-0">
                        <div class="flex items-center gap-1.5 flex-wrap">
                          <span
                            :class="isCurrentUser(p) ? 'text-amber-950 font-bold' : 'text-slate-800 font-semibold'"
                            class="truncate"
                          >
                            {{ p.name }}
                          </span>
                          <span
                            v-if="isCurrentUser(p)"
                            class="px-1.5 py-0.5 text-[9px] font-extrabold uppercase tracking-wider rounded bg-amber-500 text-white shrink-0"
                          >
                            Anda
                          </span>
                        </div>
                      </div>
                    </div>
                  </td>

                  <!-- Gadai Emas Column -->
                  <td class="w-36 min-w-[130px] px-3.5 py-3.5 text-right whitespace-nowrap">
                    <div class="inline-flex items-baseline justify-end gap-1 font-mono tabular-nums">
                      <span class="text-[10px] text-slate-400 font-sans font-medium">Rp</span>
                      <span class="text-xs font-medium text-slate-700">{{ formatNum(p.total_emas ?? 0) }}</span>
                    </div>
                  </td>

                  <!-- Non Emas Column -->
                  <td class="w-36 min-w-[130px] px-3.5 py-3.5 text-right whitespace-nowrap">
                    <div class="inline-flex items-baseline justify-end gap-1 font-mono tabular-nums">
                      <span class="text-[10px] text-slate-400 font-sans font-medium">Rp</span>
                      <span class="text-xs font-medium text-slate-700">{{ formatNum(p.total_non_emas ?? 0) }}</span>
                    </div>
                  </td>

                  <!-- Total Omzet Column -->
                  <td class="w-44 min-w-[150px] px-4 py-3.5 text-right whitespace-nowrap">
                    <div class="inline-flex items-baseline justify-end gap-1 font-mono tabular-nums">
                      <span class="text-[10px] text-emerald-600 font-sans font-bold">Rp</span>
                      <span class="text-xs font-bold text-slate-900">{{ formatNum(p.total_omzet ?? 0) }}</span>
                    </div>
                  </td>

                  <!-- Poin Column -->
                  <td class="w-28 min-w-[95px] px-3.5 py-3.5 text-right whitespace-nowrap">
                    <div class="inline-flex items-center justify-end">
                      <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-md font-mono text-xs font-bold bg-amber-50 text-amber-800 border border-amber-200/70 tabular-nums">
                        {{ p.points }}
                        <span class="text-[9px] font-sans font-medium text-amber-600">pts</span>
                      </span>
                    </div>
                  </td>
                </tr>

                <!-- Empty State (No Results) -->
                <tr v-if="filteredPerformers.length === 0">
                  <td colspan="6" class="py-12 text-center">
                    <div class="max-w-xs mx-auto flex flex-col items-center">
                      <div class="w-12 h-12 rounded-2xl bg-slate-100 text-slate-400 flex items-center justify-center mb-3">
                        <AppIcon name="search" class="w-6 h-6" />
                      </div>
                      <p class="text-xs font-semibold text-slate-700">Tidak ada data sales ditemukan</p>
                      <p class="text-[11px] text-slate-400 mt-1">
                        {{ leaderboardSearch ? `Tidak ada hasil pencarian untuk "${leaderboardSearch}"` : 'Belum ada data performa leaderboard periode ini.' }}
                      </p>
                      <button
                        v-if="leaderboardSearch"
                        @click="leaderboardSearch = ''"
                        class="mt-3 px-3 py-1.5 text-xs font-medium text-amber-600 hover:text-amber-700 bg-amber-50 rounded-lg cursor-pointer"
                      >
                        Reset Pencarian
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="px-6 py-4 bg-slate-50 border-t border-slate-100 flex flex-col sm:flex-row items-center justify-between gap-3 text-xs">
          <div class="flex items-center gap-3 text-slate-500">
            <span>
              Menampilkan <strong class="font-semibold text-slate-700">{{ filteredPerformers.length }}</strong> dari <strong class="font-semibold text-slate-700">{{ leaderboard.top_performers?.length || 0 }}</strong> Sales
            </span>
            <span class="w-1 h-1 rounded-full bg-slate-300"></span>
            <span>
              Total Omzet: <strong class="font-mono font-bold text-slate-800">Rp {{ formatNum(totalLeaderboardOmzet) }}</strong>
            </span>
          </div>
          <button
            @click="$emit('update:show', false)"
            class="w-full sm:w-auto px-5 py-2 bg-slate-800 hover:bg-slate-900 active:scale-95 text-white rounded-xl text-xs font-semibold transition-all shadow-xs cursor-pointer"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import type { LeaderboardInfo } from '../types/dashboard.types'

interface Props {
  show?: boolean
  leaderboard?: LeaderboardInfo
  selectedSalesUser?: string
  salesName?: string
}

const props = withDefaults(defineProps<Props>(), {
  show: false,
  leaderboard: () => ({ rank: 1, total_participants: 1, gap_to_next: 0, top_performers: [] }),
  selectedSalesUser: '',
  salesName: 'Sales'
})

const emit = defineEmits<{
  (e: 'update:show', val: boolean): void
}>()

const leaderboardSearch = ref('')
const leaderboardActiveFilter = ref<'all' | 'top10' | 'emas' | 'nonemas'>('all')

const filterOptions = [
  { key: 'all' as const, label: 'Semua' },
  { key: 'top10' as const, label: 'Top 10' },
  { key: 'emas' as const, label: 'Emas Terbanyak' },
  { key: 'nonemas' as const, label: 'Elektronik' }
]

const isCurrentUser = (p: { name?: string; employee_id?: number }) => {
  if (!p) return false
  if (props.selectedSalesUser && p.name && p.name.toLowerCase().includes(props.selectedSalesUser.toLowerCase())) {
    return true
  }
  if (props.salesName && p.name && p.name.toLowerCase() === props.salesName.toLowerCase()) {
    return true
  }
  return false
}

const getInitials = (name: string) => {
  if (!name) return 'S'
  const parts = name.trim().split(/\s+/)
  if (parts.length >= 2 && parts[0] && parts[1]) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return name.slice(0, 2).toUpperCase()
}

const topThreePerformers = computed(() => {
  const list = props.leaderboard.top_performers || []
  return list.slice(0, 3)
})

const totalLeaderboardOmzet = computed(() => {
  const list = props.leaderboard.top_performers || []
  return list.reduce((acc, curr) => acc + (Number(curr.total_omzet) || 0), 0)
})

const filteredPerformers = computed(() => {
  let list = [...(props.leaderboard.top_performers || [])]

  const query = leaderboardSearch.value.trim().toLowerCase()
  if (query) {
    list = list.filter((p) => (p.name || "").toLowerCase().includes(query))
  }

  if (leaderboardActiveFilter.value === 'top10') {
    list = list.slice(0, 10)
  } else if (leaderboardActiveFilter.value === 'emas') {
    list.sort((a, b) => (Number(b.total_emas) || 0) - (Number(a.total_emas) || 0))
  } else if (leaderboardActiveFilter.value === 'nonemas') {
    list.sort((a, b) => (Number(b.total_non_emas) || 0) - (Number(a.total_non_emas) || 0))
  }

  return list
})

const handleGlobalKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && props.show) {
    emit('update:show', false)
  }
}

watch(() => props.show, (isOpen) => {
  if (typeof document !== 'undefined') {
    document.body.style.overflow = isOpen ? 'hidden' : ''
  }
})

onMounted(() => {
  if (typeof window !== 'undefined') {
    window.addEventListener('keydown', handleGlobalKeydown)
  }
})

onUnmounted(() => {
  if (typeof window !== 'undefined') {
    window.removeEventListener('keydown', handleGlobalKeydown)
  }
  if (typeof document !== 'undefined') {
    document.body.style.overflow = ''
  }
})

const formatNum = (val?: number) => {
  if (!val && val !== 0) return '0'
  return new Intl.NumberFormat('id-ID').format(val)
}
</script>
