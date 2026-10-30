// Older records stored source metadata and progress notes in ket_rencana.
// New records keep the item description in ket_rencana and notes in ket_aktivitas.
const LEGACY_NOTE = /^(\[SOA:|Update progress status:|\[Hot Prospect\]|\[Leads\]|\[Prospect\]|\[SBG\])/

export function getRencanaDetails(item: Record<string, any>) {
  const note = String(item.ket_aktivitas || item.ket_rencana || '')
  const header = note.match(/^\[SOA: ([^\]]+)\]/)?.[1] || ''
  const metadata = Object.fromEntries(header.split(' | ').map((part) => {
    const separator = part.indexOf(': ')
    return separator < 0 ? ['SOA', part] : [part.slice(0, separator), part.slice(separator + 2)]
  }))
  const description = String(item.ket_rencana || '').trim()
  // Keterangan barang is pure user input: never show SOA/Aktifitas/Jenis metadata
  // or auto-generated progress notes that older records stored in ket_rencana.
  const isLegacyNote = LEGACY_NOTE.test(description)
  return {
    keterangan_aktivitas: String(item.ket_aktivitas || '').replace(/^\[SOA: [^\]]+\]\s*/, ''),
    soa: item.soa || (['Marketing', 'POS', 'Branch'].includes(item.id_sumbercust) ? item.id_sumbercust : ''),
    aktifitas: item.aktifitas || '',
    jenis_aplikasi: item.jenis_aplikasi || '',
    rt: item.rt,
    rw: item.rw,
    keterangan_barang: isLegacyNote ? '' : String(item.keterangan_barang ?? description),
  }
}
