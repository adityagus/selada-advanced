export interface MasterKatalogItem {
  id_motor: string
  nama_motor: string
  gambar_color: string
}

export interface ColorSpec {
  color_name: string
  gambar_color: string
}

export interface MasterMotorcycleSpec {
  id_motor: string
  nama_motor: string
  harga: number
  detail: string
  brosur: string
  colors: ColorSpec[]
  fiturs: string[]
}

export interface SalesPointResponse {
  point: number
}
