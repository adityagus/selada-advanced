import { useQuery } from '@tanstack/vue-query'
import { rencanaService } from '~/services/rencanaService'
import type { RencanaFilterParams } from '~/types/rencana'
import type { Ref, ComputedRef } from 'vue'

export function useInquiriesQuery(paramsRef?: Ref<RencanaFilterParams> | ComputedRef<RencanaFilterParams>) {
  return useQuery({
    queryKey: computed(() => ['inquiries', paramsRef?.value || {}]),
    queryFn: () => rencanaService.getInquiries(paramsRef?.value),
    staleTime: 1000 * 60 * 3
  })
}

export function useInquiryDetailQuery(idRef: Ref<string | null> | ComputedRef<string | null>) {
  return useQuery({
    queryKey: computed(() => ['inquiry-detail', idRef.value]),
    queryFn: async () => {
      const id = idRef.value
      if (!id) return null
      const res = await rencanaService.getInquiries({ search: id, start: 0, length: 100 })
      let rawList = Array.isArray(res) ? res : (res?.data || [])
      if (!rawList || rawList.length === 0) {
        const fallback = await rencanaService.getInquiries({ start: 0, length: 500 })
        rawList = Array.isArray(fallback) ? fallback : (fallback?.data || [])
      }
      return rawList
    },
    enabled: computed(() => !!idRef.value),
    staleTime: 1000 * 60 * 5
  })
}
