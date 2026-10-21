import axiosClient from '~/utils/axiosClient'
import type { MasterKelurahan } from '../types/masterLocation.types'

export const masterLocationService = {
  async searchKelurahan(query?: string): Promise<MasterKelurahan[]> {
    const response = await axiosClient.get<MasterKelurahan[]>('/kelurahan', {
      params: query ? { q: query } : {}
    })
    return response.data
  }
}
