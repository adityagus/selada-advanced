import { ref } from 'vue'
import { masterLocationService } from '../services/masterLocationService'
import type { MasterKelurahan } from '../types/masterLocation.types'

export const useMasterLocation = () => {
  const kelurahanResults = ref<MasterKelurahan[]>([])
  const isSearching = ref(false)

  const searchKelurahan = async (query: string) => {
    if (!query || query.trim().length < 2) {
      kelurahanResults.value = []
      return
    }
    isSearching.value = true
    try {
      kelurahanResults.value = await masterLocationService.searchKelurahan(query)
    } catch (err) {
      console.error('Failed searching kelurahan:', err)
      kelurahanResults.value = []
    } finally {
      isSearching.value = false
    }
  }

  return {
    kelurahanResults,
    isSearching,
    searchKelurahan
  }
}
