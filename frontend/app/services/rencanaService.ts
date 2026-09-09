import axiosClient from '~/utils/axiosClient'
import type { RencanaFilterParams, InputRencanaPayload } from '~/types/rencana'

export const rencanaService = {
  async getInquiries(params?: RencanaFilterParams): Promise<any> {
    const response = await axiosClient.get('/inquiry/filter', { params })
    return response.data
  },

  async inputRencana(payload: InputRencanaPayload): Promise<any> {
    const response = await axiosClient.post('/rencana/input', payload)
    console.log('response', response.data)
    return response.data
  },

  async deleteRencana(id: number): Promise<any> {
    const response = await axiosClient.delete(`/rencana/${id}`)
    return response.data
  }
}
