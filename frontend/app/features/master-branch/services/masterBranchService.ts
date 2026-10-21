import axiosClient from '~/utils/axiosClient'
import type { MasterBranch, CurrentBranchResponse } from '../types/masterBranch.types'

export const masterBranchService = {
  async getBranchList(): Promise<MasterBranch[]> {
    const response = await axiosClient.get<MasterBranch[]>('/master/cabang')
    return response.data
  },

  async getCurrentBranch(): Promise<CurrentBranchResponse> {
    const response = await axiosClient.get<CurrentBranchResponse>('/branch/current')
    return response.data
  }
}
