import { useQuery } from '@tanstack/vue-query'
import { masterService } from '~/services/masterService'
import type { Ref, ComputedRef } from 'vue'

export function useSumberCustQuery() {
  return useQuery({
    queryKey: ['master-sumbercust'],
    queryFn: () => masterService.getSumberCust(),
    staleTime: 1000 * 60 * 30
  })
}

export function useSumberSourceQuery() {
  return useQuery({
    queryKey: ['master-sumber-source'],
    queryFn: () => masterService.getSumberSource(),
    staleTime: 1000 * 60 * 30
  })
}

export function useSalesListQuery() {
  return useQuery({
    queryKey: ['master-sales-list'],
    queryFn: () => masterService.getSalesList(),
    staleTime: 1000 * 60 * 10
  })
}

export function useKelurahanQuery(queryRef: Ref<string> | ComputedRef<string>) {
  return useQuery({
    queryKey: computed(() => ['master-kelurahan', queryRef.value]),
    queryFn: () => masterService.getKelurahan(queryRef.value),
    enabled: computed(() => !!queryRef.value && queryRef.value.length >= 2),
    staleTime: 1000 * 60 * 15
  })
}
