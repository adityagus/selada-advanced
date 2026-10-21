import { ref, computed } from 'vue'
import { masterBranchService } from '../services/masterBranchService'
import type { MasterBranch } from '../types/masterBranch.types'

const branches = ref<MasterBranch[]>([])
const isFetched = ref(false)
const isLoading = ref(false)
const currentBranchName = ref('')

export const useMasterBranch = () => {
  const fetchBranches = async (forceRefresh = false) => {
    if (isFetched.value && !forceRefresh) return
    isLoading.value = true
    try {
      branches.value = await masterBranchService.getBranchList()
      isFetched.value = true
    } catch (err) {
      console.error('Failed to load branches:', err)
    } finally {
      isLoading.value = false
    }
  }

  const fetchCurrentBranch = async () => {
    try {
      const res = await masterBranchService.getCurrentBranch()
      currentBranchName.value = res.branch_name
      return res
    } catch (err) {
      console.error('Failed to load current branch:', err)
      return null
    }
  }

  const branchMap = computed(() => {
    const map = new Map<number, string>()
    for (const b of branches.value) {
      map.set(b.id_cabang, b.nm_cabang)
    }
    return map
  })

  const getBranchName = (id: number): string => {
    return branchMap.value.get(id) || `Cabang #${id}`
  }

  return {
    branches,
    isLoading,
    currentBranchName,
    fetchBranches,
    fetchCurrentBranch,
    getBranchName
  }
}
