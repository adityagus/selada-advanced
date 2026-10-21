import axiosClient from '~/utils/axiosClient'
import type { MasterSumberCust, MasterApplicationSource } from '../types/masterSource.types'

export const masterSourceService = {
  async getSumberCust(): Promise<MasterSumberCust[]> {
    const response = await axiosClient.get<MasterSumberCust[]>('/sumbercust')
    return response.data
  },

  async getSumberSource(): Promise<MasterApplicationSource[]> {
    const response = await axiosClient.get<MasterApplicationSource[]>('/master/sumber_source')
    return response.data
  }
}
