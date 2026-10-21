export interface SoaRule {
  aktifitasList: string[]
  jenisAplikasiList: string[]
}

export const SOA_OPTIONS_MAP: Record<string, SoaRule> = {
  Marketing: {
    aktifitasList: ['FB', 'Sosmed', 'Event Internal', 'Direct Selling', 'Komunitas'],
    jenisAplikasiList: ['Take Over', 'Pick Up', 'Walk In']
  },
  POS: {
    aktifitasList: ['FB', 'Sosmed', 'Event Internal', 'Direct Selling', 'Banner & Ads cetak'],
    jenisAplikasiList: ['Walk In']
  },
  Branch: {
    aktifitasList: [
      'CRM',
      'Website',
      'FB',
      'Sosmed',
      'Digital Marketing Agency',
      'Event Internal',
      'Direct Selling',
      'Komunitas',
      'Banner & Ads cetak'
    ],
    jenisAplikasiList: ['Take Over', 'Pick Up', 'Walk In']
  }
}

export const SOA_LIST = Object.keys(SOA_OPTIONS_MAP)

export function getAktifitasOptions(soa: string): string[] {
  if (!soa) return []
  return SOA_OPTIONS_MAP[soa]?.aktifitasList || []
}

export function getJenisAplikasiOptions(soa: string): string[] {
  if (!soa) return []
  return SOA_OPTIONS_MAP[soa]?.jenisAplikasiList || []
}
