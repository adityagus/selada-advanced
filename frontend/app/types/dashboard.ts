export interface SalesFunnelCounts {
  all: number
  leads: number
  prospect: number
  hotprospect: number
  sbg: number
  batal: number
  leads_monthly: number
  prospect_monthly: number
  hotprospect_monthly: number
  sbg_monthly: number
}

export interface LiveFeedItem {
  id?: string
  code: string
  date?: string
  type: string
  customer: string
  amount: number
  poin?: number
  sub: string
  loan_amount?: number
  iconText: string
  isBonus?: boolean
  is_bonus?: boolean
}

export interface EmasCardInfo {
  target_emas_g: number
  pct_emas: number
  bonus_emas_title: string
  sub_emas_text: string
  is_emas_bonus_done: boolean
}

export interface NonEmasCardInfo {
  pct_nonemas: number
  gap_nonemas: number
  bonus_title: string
  sub_text: string
  is_nonemas_bonus_done: boolean
}

export interface LeaderboardInfo {
  rank: number
  total_participants: number
  gap_to_next: number
}

export interface DashboardMetricsResponse {
  counts?: SalesFunnelCounts
  newcif?: number
  last_month_cif?: number
  leads?: number
  leads_monthly?: number
  prospect_monthly?: number
  hotprospect_monthly?: number
  sbg_monthly?: number
  last_month_leads?: number
  emas?: number
  nonemas?: number
  gramasi_emas?: number
  total_trx_count?: number
  est_incentive?: number
  main_incentive_amount?: number
  total_take_home?: number
  is_qualified?: boolean
  user_grade?: string
  applied_rate?: number
  poin_total?: number
  target_goal_points?: number
  daily_target_progress?: number
  target_description?: string
  motivation_text?: string
  is_review?: boolean
  sales_list?: Array<{
    username: string
    nama: string
    id_jabatan: number
    nm_jabatan: string
    fkuser: string
  }>
  selected_sales_user?: string
  label_point?: { icon: string; text: string }
  label_leaderboard?: { icon: string; text: string }
  emas_card?: EmasCardInfo
  nonemas_card?: NonEmasCardInfo
  leaderboard?: LeaderboardInfo
  pacing_daily?: {
    selada?: Record<number, number>
    sopiga?: Record<number, number>
  }
  live_feeds?: LiveFeedItem[]
  username?: string
}
