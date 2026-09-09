import axiosClient from '~/utils/axiosClient'
import type { DashboardMetricsResponse } from '~/types/dashboard'

export const dashboardService = {
  async getMetrics(salesUser?: string): Promise<DashboardMetricsResponse> {
    const params: Record<string, string> = {}
    if (salesUser) {
      params.sales_user = salesUser
    }
    const response = await axiosClient.get<DashboardMetricsResponse>('/dashboard/metrics', { params })
    return response.data
  },

  async getInquiryCounts(salesUser?: string): Promise<Record<string, number>> {
    const params: Record<string, string> = {}
    if (salesUser) {
      params.sales_user = salesUser
    }
    const response = await axiosClient.get<Record<string, number>>('/inquiry/counts', { params })
    return response.data
  }
}
