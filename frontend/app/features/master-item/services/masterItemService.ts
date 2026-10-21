import axiosClient from '~/utils/axiosClient'
import type { MasterKatalogItem, MasterMotorcycleSpec, SalesPointResponse } from '../types/masterItem.types'

export const masterItemService = {
  async getKatalogList(category: 'sport' | 'trail' | 'moped'): Promise<MasterKatalogItem[]> {
    const response = await axiosClient.get<MasterKatalogItem[]>(`/katalog/list/${category}`)
    return response.data
  },

  async getKatalogSpec(id: string): Promise<MasterMotorcycleSpec> {
    const response = await axiosClient.get<MasterMotorcycleSpec>(`/katalog/spec/${id}`)
    return response.data
  },

  async getSalesPoint(): Promise<SalesPointResponse> {
    const response = await axiosClient.get<SalesPointResponse>('/katalog/point')
    return response.data
  }
}
