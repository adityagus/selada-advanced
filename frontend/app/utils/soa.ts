export interface MasterItem {
  id: number
  code: string
  name: string
}

// Default fallback items from db_gadaimulia
export const DEFAULT_SOA_ITEMS: MasterItem[] = [
  { id: 6, code: '02', name: 'MARKETING' },
  { id: 8, code: '04', name: 'WALK IN' },
  { id: 11, code: 'CRM', name: 'CRM' },
  { id: 12, code: 'MITRA-BISNIS', name: 'MITRA BISNIS' },
  { id: 13, code: 'MITRA-PERORAANGAN', name: 'MITRA PERORANGAN' },
  { id: 10, code: 'CGC', name: 'CGC' }
]

export const DEFAULT_AKTIFITAS_ITEMS: MasterItem[] = [
  { id: 1, code: '2', name: 'MEDIA SOSIAL' },
  { id: 3, code: '7', name: 'PRESENTASI' },
  { id: 13, code: 'FACEBOOK', name: 'FACEBOOK' },
  { id: 6, code: '11', name: 'IVENT INTERNAL' },
  { id: 4, code: '9', name: 'CRM' },
  { id: 5, code: '10', name: 'KANVASING/DIRECT SELLING' }
]

export const DEFAULT_JENIS_APLIKASI_ITEMS: MasterItem[] = [
  { id: 9, code: 'TO', name: 'TAKE OVER' },
  { id: 7, code: '01', name: 'PICKUP' },
  { id: 8, code: '02', name: 'WALK IN' }
]

export const SOA_LIST = DEFAULT_SOA_ITEMS.map(i => i.name)

export function parseSoaText(val: any, dynamicList?: MasterItem[]): string {
  if (val === undefined || val === null || val === '') return '-'
  const str = String(val).trim()
  const list = dynamicList && dynamicList.length > 0 ? dynamicList : DEFAULT_SOA_ITEMS
  const match = list.find(item => String(item.id) === str || item.code.toLowerCase() === str.toLowerCase() || item.name.toLowerCase() === str.toLowerCase())
  return match ? match.name : str
}

export function parseAktifitasText(val: any, dynamicList?: MasterItem[]): string {
  if (val === undefined || val === null || val === '' || val === '-') return '-'
  const str = String(val).trim()
  const list = dynamicList && dynamicList.length > 0 ? dynamicList : DEFAULT_AKTIFITAS_ITEMS
  const match = list.find(item => String(item.id) === str || item.code.toLowerCase() === str.toLowerCase() || item.name.toLowerCase() === str.toLowerCase())
  return match ? match.name : str
}

export function parseJenisAplikasiText(val: any, dynamicList?: MasterItem[]): string {
  if (val === undefined || val === null || val === '') return '-'
  const str = String(val).trim()
  const list = dynamicList && dynamicList.length > 0 ? dynamicList : DEFAULT_JENIS_APLIKASI_ITEMS
  const match = list.find(item => String(item.id) === str || item.code.toLowerCase() === str.toLowerCase() || item.name.toLowerCase() === str.toLowerCase())
  return match ? match.name : str
}

export function getAktifitasOptions(soa?: string, dynamicList?: MasterItem[]): MasterItem[] {
  return dynamicList && dynamicList.length > 0 ? dynamicList : DEFAULT_AKTIFITAS_ITEMS
}

export function getJenisAplikasiOptions(soa?: string, dynamicList?: MasterItem[]): MasterItem[] {
  return dynamicList && dynamicList.length > 0 ? dynamicList : DEFAULT_JENIS_APLIKASI_ITEMS
}
