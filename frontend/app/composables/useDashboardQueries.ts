import { useQuery } from '@tanstack/vue-query'
import { dashboardService } from '~/services/dashboardService'
import type { Ref, ComputedRef } from 'vue'

export function useDashboardMetricsQuery(salesUserRef?: Ref<string | undefined> | ComputedRef<string | undefined>) {
  return useQuery({
    queryKey: computed(() => ['dashboard-metrics', salesUserRef?.value || '']),
    queryFn: () => dashboardService.getMetrics(salesUserRef?.value),
    staleTime: 1000 * 60 * 2
  })
}

export function useInquiryCountsQuery(salesUserRef?: Ref<string | undefined> | ComputedRef<string | undefined>) {
  return useQuery({
    queryKey: computed(() => ['inquiry-counts', salesUserRef?.value || '']),
    queryFn: () => dashboardService.getInquiryCounts(salesUserRef?.value),
    staleTime: 1000 * 60 * 2
  })
}
