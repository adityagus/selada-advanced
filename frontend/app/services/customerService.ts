import axiosClient from '~/utils/axiosClient'
import type { CustomerSelectItem } from '~/types/customer'

export const customerService = {
  async getSelectCustomers(query?: string): Promise<CustomerSelectItem[]> {
    const params: Record<string, string> = {}
    if (query) {
      params.q = query
    }
    const response = await axiosClient.get<CustomerSelectItem[]>('/customer/select', { params })
    return response.data
  }
}
