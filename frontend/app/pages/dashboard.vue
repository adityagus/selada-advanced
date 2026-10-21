<template>
  <div class="space-y-6 pb-12 font-sans">
    <!-- Header Section -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-slate-900 tracking-tight">
          Dashboard
        </h1>
        <p v-if="isReview && selectedSalesUser" class="text-xs text-slate-500 mt-1 flex items-center gap-2">
          <AppIcon name="user" class="w-4 h-4 text-amber-500" />
          <span>Meninjau Dashboard Sales:</span>
          <span class="font-semibold text-amber-800 bg-amber-50 px-2.5 py-0.5 rounded-lg border border-amber-200/60">{{ activeSalesName }}</span>
        </p>
      </div>

      <div class="flex items-center gap-3">
        <NuxtLink
          to="/inquiry"
          class="inline-flex items-center gap-2 bg-amber-500 hover:bg-amber-600 active:scale-95 text-white font-medium px-4 py-2.5 rounded-xl text-sm transition-all shadow-xs"
        >
          <AppIcon name="plus" class="w-4 h-4" />
          Leads Baru
        </NuxtLink>
      </div>
    </div>

    <!-- Section 1: Hero & Target & Leaderboard Summary -->
    <DashboardHero
      :is-review="isReview"
      v-model:selected-sales-user="selectedSalesUser"
      :sorted-sales-list="sortedSalesList"
      :is-fetching="isFetching"
      :est-incentive="estIncentive"
      :is-qualified="isQualified"
      :user-grade="userGrade"
      :applied-rate="appliedRate"
      :point-total="pointTotal"
      :target-goal-points="targetGoalPoints"
      :target-progress-pct="targetProgressPct"
      :target-description="targetDescription"
      :motivation-text="motivationText"
      :leaderboard="leaderboard"
      @review-change="handleReviewSales"
      @review-reset="resetReview"
      @open-leaderboard="showLeaderboardModal = true"
    />

    <!-- Section 2: 4 Metrics Cards -->
    <DashboardMetrics
      :stats="stats"
      :emas-card="emasCard"
      :non-emas-card="nonEmasCard"
      :gramasi-emas="gramasiEmas"
      :last-month-cif="lastMonthCif"
      :last-month-leads="lastMonthLeads"
      :monthly-funnel="monthlyFunnel"
    />

    <!-- Section 3: Sales Funnel -->
    <DashboardFunnel
      :stats="stats"
      :monthly-funnel="monthlyFunnel"
      :is-fetching="isFetching"
      @refresh="refetch()"
    />

    <!-- Section 4: Daily Pacing Chart & Live Feeds -->
    <DashboardPacingChart
      :pacing-selada="pacingSelada"
      :pacing-sopiga="pacingSopiga"
      :live-feeds="liveFeeds"
    />

    <!-- Section 5: Full Leaderboard Modal -->
    <DashboardLeaderboardModal
      v-model:show="showLeaderboardModal"
      :leaderboard="leaderboard"
      :selected-sales-user="selectedSalesUser"
      :sales-name="salesName"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { useDashboardMetricsQuery } from '~/composables/useDashboardQueries'
import type { LeaderboardInfo } from '~/features/dashboard/types/dashboard.types'

const showLeaderboardModal = ref(false)

const stats = ref({
  all: 0,
  leads: 0,
  prospect: 0,
  hotprospect: 0,
  sbg: 0,
  batal: 0,
  newcif: 0,
  emas: 0,
  nonemas: 0,
  leads_monthly: 0
})

const isReview = ref(false)
const salesList = ref<any[]>([])
const selectedSalesUser = ref('')

const sortedSalesList = computed(() => {
  return [...salesList.value].sort((a, b) => {
    const omzetA = Number(a.total_omzet ?? 0)
    const omzetB = Number(b.total_omzet ?? 0)
    if (omzetB !== omzetA) return omzetB - omzetA
    return (a.nama || '').localeCompare(b.nama || '')
  })
})

const monthlyFunnel = reactive({
  leads: 0,
  prospect: 0,
  hotprospect: 0,
  sbg: 0
})

const estIncentive = ref(0)
const isQualified = ref(false)
const userGrade = ref('JUNIOR')
const appliedRate = ref(25)
const pointTotal = ref(0)
const targetGoalPoints = ref(25)
const dailyTargetProgress = ref(0)
const targetDescription = ref('Mengejar Qualified 1 (Target: 25 Poin)')
const motivationText = ref('Semangat! Kejar Poin untuk Bonus Insentif!')

const salesName = ref('Sales')
const leaderboard = ref<LeaderboardInfo>({
  rank: 1,
  total_participants: 1,
  gap_to_next: 0,
  top_performers: []
})

const gramasiEmas = ref(0)
const lastMonthCif = ref(0)
const lastMonthLeads = ref(0)

const emasCard = ref({
  target_emas_g: 300,
  pct_emas: 0,
  bonus_emas_title: 'Tier Bonus: 300g = 1 Jt',
  sub_emas_text: 'Target berikutnya: 300g (Bonus 1 Jt)',
  is_emas_bonus_done: false
})

const nonEmasCard = ref({
  pct_nonemas: 0,
  gap_nonemas: 0,
  bonus_title: 'Bonus UP 50 Jt = 500 Rb',
  sub_text: 'Target UP: 50 Jt',
  is_nonemas_bonus_done: false
})

const targetProgressPct = computed(() => {
  if (dailyTargetProgress.value !== undefined && dailyTargetProgress.value !== null) {
    return Math.min(100, Math.max(0, Math.round(dailyTargetProgress.value)))
  }
  if (!targetGoalPoints.value) return 0
  return Math.min(100, Math.max(0, Math.round((pointTotal.value / targetGoalPoints.value) * 100)))
})

const pacingSelada = ref<Record<number, number>>({})
const pacingSopiga = ref<Record<number, number>>({})
const liveFeeds = ref<any[]>([])

// TanStack Query integration
const {
  data: metricsData,
  isFetching,
  refetch
} = useDashboardMetricsQuery(selectedSalesUser)

const handleReviewSales = () => {
  refetch()
}

const resetReview = () => {
  selectedSalesUser.value = ''
  refetch()
}

const activeSalesName = computed(() => {
  if (!selectedSalesUser.value) return salesName.value
  const targetUser = selectedSalesUser.value.toLowerCase().trim()
  const s = sortedSalesList.value.find(
    (item) => item.username && item.username.toLowerCase().trim() === targetUser
  )
  if (!s) return selectedSalesUser.value.toUpperCase()
  return s.nm_jabatan ? `${s.nama} (${s.nm_jabatan})` : s.nama
})

watch(
  metricsData,
  (data) => {
    if (!data) return
    if (data.counts) {
      stats.value = { ...stats.value, ...data.counts }
    }
    if (data.newcif !== undefined) stats.value.newcif = data.newcif
    if (data.last_month_cif !== undefined) lastMonthCif.value = data.last_month_cif
    if (data.leads_monthly !== undefined) stats.value.leads_monthly = data.leads_monthly
    if (data.last_month_leads !== undefined) lastMonthLeads.value = data.last_month_leads
    if (data.emas !== undefined) stats.value.emas = data.emas
    if (data.nonemas !== undefined) stats.value.nonemas = data.nonemas
    if (data.gramasi_emas !== undefined) gramasiEmas.value = data.gramasi_emas

    if (data.est_incentive !== undefined) estIncentive.value = data.est_incentive
    if (data.is_qualified !== undefined) isQualified.value = data.is_qualified
    if (data.user_grade !== undefined) userGrade.value = data.user_grade
    if (data.applied_rate !== undefined) appliedRate.value = data.applied_rate
    if (data.poin_total !== undefined) pointTotal.value = data.poin_total
    if (data.target_goal_points !== undefined) targetGoalPoints.value = data.target_goal_points
    if (data.daily_target_progress !== undefined) dailyTargetProgress.value = data.daily_target_progress
    if (data.target_description !== undefined) targetDescription.value = data.target_description
    if (data.motivation_text !== undefined) motivationText.value = data.motivation_text

    if (data.is_review !== undefined) isReview.value = data.is_review
    if (data.sales_list && Array.isArray(data.sales_list)) salesList.value = data.sales_list
    if (data.selected_sales_user && !selectedSalesUser.value) {
      selectedSalesUser.value = data.selected_sales_user
    } else if (data.is_review && !selectedSalesUser.value && sortedSalesList.value.length > 0) {
      selectedSalesUser.value = sortedSalesList.value[0].username
    }

    if (data.emas_card) emasCard.value = { ...emasCard.value, ...data.emas_card }
    if (data.nonemas_card) nonEmasCard.value = { ...nonEmasCard.value, ...data.nonemas_card }

    if (data.leaderboard) leaderboard.value = { ...leaderboard.value, ...data.leaderboard }
    if (data.live_feeds && Array.isArray(data.live_feeds)) liveFeeds.value = data.live_feeds
    if (data.username) salesName.value = data.username

    if (data.branch_name) {
      const branchNameCookie = useCookie('user_branch_name')
      if (!branchNameCookie.value || /^\d+$/.test(branchNameCookie.value)) {
        branchNameCookie.value = data.branch_name
      }
    }

    if (data.pacing_daily) {
      if (data.pacing_daily.selada) {
        const sMap: Record<number, number> = {}
        for (const [k, v] of Object.entries(data.pacing_daily.selada)) {
          sMap[Number(k)] = Number(v)
        }
        pacingSelada.value = sMap
      }
      if (data.pacing_daily.sopiga) {
        const spMap: Record<number, number> = {}
        for (const [k, v] of Object.entries(data.pacing_daily.sopiga)) {
          spMap[Number(k)] = Number(v)
        }
        pacingSopiga.value = spMap
      }
    }

    monthlyFunnel.leads =
      data.leads_monthly ??
      data.counts?.leads_monthly ??
      data.counts?.leads ??
      stats.value.leads ??
      0
    monthlyFunnel.prospect = data.prospect_monthly ?? data.counts?.prospect_monthly ?? 0
    monthlyFunnel.hotprospect = data.hotprospect_monthly ?? data.counts?.hotprospect_monthly ?? 0
    monthlyFunnel.sbg = data.sbg_monthly ?? data.counts?.sbg_monthly ?? 0
  },
  { immediate: true }
)
</script>
