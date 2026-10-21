import { ref } from 'vue'
import { masterSourceService } from '../services/masterSourceService'
import type { MasterSumberCust, MasterApplicationSource } from '../types/masterSource.types'

const sumberCustList = ref<MasterSumberCust[]>([])
const applicationSourceList = ref<MasterApplicationSource[]>([])
const isFetchedSumberCust = ref(false)
const isFetchedAppSource = ref(false)

export const useMasterSource = () => {
  const fetchSumberCust = async (force = false) => {
    if (isFetchedSumberCust.value && !force) return
    try {
      sumberCustList.value = await masterSourceService.getSumberCust()
      isFetchedSumberCust.value = true
    } catch (err) {
      console.error('Failed to load sumber cust:', err)
    }
  }

  const fetchSumberSource = async (force = false) => {
    if (isFetchedAppSource.value && !force) return
    try {
      applicationSourceList.value = await masterSourceService.getSumberSource()
      isFetchedAppSource.value = true
    } catch (err) {
      console.error('Failed to load application source:', err)
    }
  }

  return {
    sumberCustList,
    applicationSourceList,
    fetchSumberCust,
    fetchSumberSource
  }
}
