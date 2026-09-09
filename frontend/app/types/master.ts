export interface SalesUserItem {
  username: string
  nama: string
  id_jabatan: number
  nm_jabatan: string
  fkuser: string
}

export interface SumberCustItem {
  id_sumbercust: number
  nm_sumbercust: string
}

export interface ApplicationSourceItem {
  code: string
  name: string
}

export interface KelurahanItem {
  id_kelurahan?: number
  nm_kelurahan: string
  nm_kecamatan?: string
  nm_kota?: string
}
