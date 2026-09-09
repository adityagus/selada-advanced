export interface RencanaFilterParams {
  start?: number
  length?: number
  search?: string
  status?: string
  sales_user?: string
}

export interface InputRencanaPayload {
  cust?: string
  custlama?: number
  nm?: string
  nama_customer?: string
  hp?: string
  tgl_rencana?: string
  jam_rencana?: string
  status?: string
  hslaktiv?: number
  katcust?: string
  ket?: string
  ket_rencana?: string
  ket_aktivitas?: string
  sumbercust?: string
  id_sumbercust?: number
  id_rencana?: number
  alamat?: string
  rt?: string
  rw?: string
  kelurahan?: string
  kecamatan?: string
  kota?: string
  aktiv?: number
}

export interface InquiryItem {
  id_rencana: number
  id_customer?: number
  nama: string
  hp: string
  tgl_rencana?: string
  status: string
  nama_hslaktiv?: string
  id_hslaktiv?: number
  ket_rencana?: string
  ket_aktivitas?: string
  id_sumbercust?: number
  nm_sumbercust?: string
}
