import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { rencanaService } from '~/services/rencanaService'
import type { InputRencanaPayload } from '~/types/rencana'

export function useInputRencanaMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (payload: InputRencanaPayload) => rencanaService.inputRencana(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inquiries'] })
      queryClient.invalidateQueries({ queryKey: ['inquiry-detail'] })
      queryClient.invalidateQueries({ queryKey: ['dashboard-metrics'] })
      queryClient.invalidateQueries({ queryKey: ['inquiry-counts'] })
    }
  })
}

export function useDeleteRencanaMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => rencanaService.deleteRencana(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inquiries'] })
      queryClient.invalidateQueries({ queryKey: ['dashboard-metrics'] })
      queryClient.invalidateQueries({ queryKey: ['inquiry-counts'] })
    }
  })
}
