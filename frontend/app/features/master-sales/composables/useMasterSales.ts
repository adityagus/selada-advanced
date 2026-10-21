import { ref, computed } from 'vue'
import { masterSalesService } from '../services/masterSalesService'
import type { MasterSalesUser } from '../types/masterSales.types'

const salesList = ref<MasterSalesUser[]>([])
const isFetched = ref(false)
const isLoading = ref(false)

export const useMasterSales = () => {
  const fetchSalesList = async (force = false) => {
    if (isFetched.value && !force) return
    isLoading.value = true
    try {
      salesList.value = await masterSalesService.getSalesList()
      isFetched.value = true
    } catch (err) {
      console.error('Failed to load sales list:', err)
    } finally {
      isLoading.value = false
    }
  }

  // Sorted by omzet descending by default
  const sortedSalesList = computed(() => {
    return [...salesList.value].sort((a, b) => {
      const omzetA = Number(a.total_omzet ?? 0)
      const omzetB = Number(b.total_omzet ?? 0)
      if (omzetB !== omzetA) return omzetB - omzetA
      return (a.nama || '').localeCompare(b.nama || '')
    })
  })

  return {
    salesList,
    sortedSalesList,
    isLoading,
    fetchSalesList
  }
}
