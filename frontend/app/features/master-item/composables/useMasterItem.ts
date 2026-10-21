import { ref } from 'vue'
import { masterItemService } from '../services/masterItemService'
import type { MasterKatalogItem, MasterMotorcycleSpec } from '../types/masterItem.types'

export const useMasterItem = () => {
  const katalogList = ref<MasterKatalogItem[]>([])
  const currentSpec = ref<MasterMotorcycleSpec | null>(null)
  const salesPoint = ref<number>(0)
  const isLoading = ref(false)

  const fetchKatalog = async (category: 'sport' | 'trail' | 'moped') => {
    isLoading.value = true
    try {
      katalogList.value = await masterItemService.getKatalogList(category)
    } catch (err) {
      console.error('Failed loading katalog:', err)
      katalogList.value = []
    } finally {
      isLoading.value = false
    }
  }

  const fetchSpec = async (id: string) => {
    isLoading.value = true
    try {
      currentSpec.value = await masterItemService.getKatalogSpec(id)
      return currentSpec.value
    } catch (err) {
      console.error('Failed loading spec:', err)
      currentSpec.value = null
      return null
    } finally {
      isLoading.value = false
    }
  }

  const fetchSalesPoint = async () => {
    try {
      const res = await masterItemService.getSalesPoint()
      salesPoint.value = res.point || 0
      return salesPoint.value
    } catch (err) {
      console.error('Failed loading sales point:', err)
      return 0
    }
  }

  return {
    katalogList,
    currentSpec,
    salesPoint,
    isLoading,
    fetchKatalog,
    fetchSpec,
    fetchSalesPoint
  }
}
