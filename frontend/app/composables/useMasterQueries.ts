import { useQuery } from '@tanstack/vue-query'
import { masterSourceService } from '~/features/master-source/services/masterSourceService'
import { masterSalesService } from '~/features/master-sales/services/masterSalesService'
import { masterLocationService } from '~/features/master-location/services/masterLocationService'
import type { Ref, ComputedRef } from 'vue'

export function useSumberCustQuery() {
  return useQuery({
    queryKey: ['master-sumbercust'],
    queryFn: () => masterSourceService.getSumberCust(),
    staleTime: 1000 * 60 * 30
  })
}

export function useSumberSourceQuery() {
  return useQuery({
    queryKey: ['master-sumber-source'],
    queryFn: () => masterSourceService.getSumberSource(),
    staleTime: 1000 * 60 * 30
  })
}

export function useSalesListQuery() {
  return useQuery({
    queryKey: ['master-sales-list'],
    queryFn: () => masterSalesService.getSalesList(),
    staleTime: 1000 * 60 * 10
  })
}

export function useKelurahanQuery(queryRef: Ref<string> | ComputedRef<string>) {
  return useQuery({
    queryKey: computed(() => ['master-kelurahan', queryRef.value]),
    queryFn: () => masterLocationService.searchKelurahan(queryRef.value),
    enabled: computed(() => !!queryRef.value && queryRef.value.length >= 2),
    staleTime: 1000 * 60 * 15
  })
}
