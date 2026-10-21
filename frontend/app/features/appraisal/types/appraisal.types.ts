export interface ElekCategory {
  kd_jenis_barang: string
  nm_jenis_barang: string
}

export interface ElekWilayah {
  region_id: string
  region_name: string
}

export interface ElekScoringOption {
  option_id: string | number
  option_name: string
  score_deduction?: number
  score_percent?: number
  is_reject?: boolean
}

export interface ElekScoringItem {
  item_id: string | number
  item_name: string
  selected: number | null
}

export interface ElekScoringSection {
  section_id: string | number
  section_name: string
  section_type: 'reject' | 'positif' | 'negatif'
  options: ElekScoringOption[]
  items?: ElekScoringItem[]
  selectedOption?: number | null
}

export interface ElekCalculationResult {
  basePrice: number
  deductions: number
  maxLoan: number
  requestedLoan: number
  adminFee: number
}

export interface EmasItem {
  id?: string | number
  collateral_name: string
  karat: number
  total: number
  gross_weight: number
  net_weight: number
  price_per_gram: number
  total_taksir: number
  notes?: string
  is_lm: boolean
}

export interface HistoryTaksiranItem {
  id?: number | string
  nama_sales?: string
  jenis_barang: string
  nama_barang: string
  nilai_taksir: number
  nilai_pinjaman: number
  items_emas?: any[]
  created_at: string
}
