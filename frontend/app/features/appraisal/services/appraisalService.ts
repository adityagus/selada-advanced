import axiosClient from '~/utils/axiosClient'
import type { ElekCategory, ElekWilayah, HistoryTaksiranItem } from '../types/appraisal.types'

export const appraisalService = {
  async getCategories(): Promise<ElekCategory[]> {
    const res = await axiosClient.get<ElekCategory[]>('/taksiran/jenis-barang')
    return res.data
  },

  async getWilayah(): Promise<ElekWilayah[]> {
    const res = await axiosClient.get<ElekWilayah[]>('/taksiran/wilayah')
    return res.data
  },

  async getElectronicScoring(categoryId: string, gradeId?: string) {
    const res = await axiosClient.get('/taksiran/electronic-scoring', {
      params: { category_id: categoryId, grade_id: gradeId }
    })
    return res.data
  },

  async getNilaiTaksiran(params: { category_id: string; region_id: string; item_id: string }) {
    const res = await axiosClient.get('/taksiran/nilai-taksir', { params })
    return res.data
  },

  async searchBarang(query: string, categoryId?: string) {
    const res = await axiosClient.get('/taksiran/barang-umum', {
      params: { q: query, category_id: categoryId }
    })
    return res.data
  },

  async getBarangEmas() {
    const res = await axiosClient.get('/taksiran/barang-emas')
    return res.data
  },

  async getGoldStle(regionId?: string) {
    const res = await axiosClient.get('/taksiran/gold-stle', {
      params: { region_id: regionId }
    })
    return res.data
  },

  async getBiayaAdminEmas() {
    const res = await axiosClient.get('/taksiran/biaya-admin-emas')
    return res.data
  },

  async simpanTaksiran(payload: {
    jenis_barang: string
    nama_barang: string
    nilai_taksir: number
    nilai_pinjaman: number
    items_emas?: any[]
  }) {
    const res = await axiosClient.post('/taksiran/simpan', payload)
    return res.data
  },

  async getHistory(params?: { sales_user?: string; search?: string }): Promise<HistoryTaksiranItem[]> {
    const res = await axiosClient.get<{ data: HistoryTaksiranItem[] }>('/taksiran/history', { params })
    return res.data?.data || []
  }
}
