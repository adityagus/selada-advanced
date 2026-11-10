// Older records stored source metadata and progress notes in ket_rencana.
// New records keep the item description in ket_rencana and notes in ket_aktivitas.
const LEGACY_NOTE = /^(\[SOA:|Update progress status:|\[Hot Prospect\]|\[Leads\]|\[Prospect\]|\[SBG\])/

export function getRencanaDetails(item: Record<string, any>) {
  if (!item) {
    return {
      keterangan_aktivitas: '',
      soa: '',
      aktifitas: '',
      jenis_aplikasi: '',
      keterangan_barang: '',
    }
  }

  const note = String(item.ket_aktivitas || item.ket_rencana || '')
  const header = note.match(/^\[SOA: ([^\]]+)\]/)?.[1] || ''
  const metadata = Object.fromEntries(header.split(' | ').map((part) => {
    const separator = part.indexOf(': ')
    return separator < 0 ? ['SOA', part] : [part.slice(0, separator), part.slice(separator + 2)]
  }))
  const description = String(item.ket_rencana || '').trim()
  const isLegacyNote = LEGACY_NOTE.test(description)
  const rawKetAktivitas = String(item.ket_aktivitas || '').trim()
  const fallbackNote = isLegacyNote ? description : ''
  return {
    keterangan_aktivitas: (rawKetAktivitas || fallbackNote).replace(/^\[SOA: [^\]]+\]\s*/, ''),
    soa: item.soa || metadata.SOA || (['Marketing', 'POS', 'Branch'].includes(item.id_sumbercust) ? item.id_sumbercust : ''),
    aktifitas: item.aktifitas || metadata.Aktifitas || '',
    jenis_aplikasi: item.jenis_aplikasi || metadata.Jenis || '',
    keterangan_barang: isLegacyNote ? '' : String(item.keterangan_barang ?? description),
  }
}
