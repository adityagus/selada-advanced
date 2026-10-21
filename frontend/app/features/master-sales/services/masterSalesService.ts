import axiosClient from '~/utils/axiosClient'
import type { MasterSalesUser } from '../types/masterSales.types'

export const masterSalesService = {
  async getSalesList(): Promise<MasterSalesUser[]> {
    const response = await axiosClient.get<MasterSalesUser[]>('/sales/list')
    return response.data
  }
}
