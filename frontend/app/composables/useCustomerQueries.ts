import { useQuery } from '@tanstack/vue-query'
import { customerService } from '~/services/customerService'
import type { Ref, ComputedRef } from 'vue'

export function useSelectCustomersQuery(queryRef?: Ref<string> | ComputedRef<string>) {
  return useQuery({
    queryKey: computed(() => ['select-customers', queryRef?.value || '']),
    queryFn: () => customerService.getSelectCustomers(queryRef?.value),
    staleTime: 1000 * 60 * 10
  })
}
