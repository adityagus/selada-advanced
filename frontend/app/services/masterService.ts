import axiosClient from '~/utils/axiosClient'
import type { SalesUserItem, SumberCustItem, ApplicationSourceItem } from '~/types/master'

export const masterService = {
  async getSumberCust(): Promise<SumberCustItem[]> {
    const response = await axiosClient.get<SumberCustItem[]>('/sumbercust')
    return response.data
  },

  async getSumberSource(): Promise<ApplicationSourceItem[]> {
    const response = await axiosClient.get<ApplicationSourceItem[]>('/master/sumber_source')
    return response.data
  },

  async getSalesList(): Promise<SalesUserItem[]> {
    const response = await axiosClient.get<SalesUserItem[]>('/sales/list')
    return response.data
  },

  async getKelurahan(query?: string): Promise<any[]> {
    const response = await axiosClient.get('/kelurahan', {
      params: query ? { q: query } : {}
    })
    return response.data
  }
}
