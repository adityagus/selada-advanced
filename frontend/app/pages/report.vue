<template>
  <div class="space-y-6 pb-12 font-sans">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-semibold text-slate-800 tracking-tight">Laporan Rekapitulasi & Performance</h1>
        <p class="text-slate-500 text-sm">Monitoring realisasi aktivitas sales, tingkat konversi, dan rekapitulasi per periode bulan & tahun.</p>
      </div>
      <div class="flex items-center gap-3">
        <button @click="printReport" class="inline-flex items-center gap-2 bg-slate-100 hover:bg-slate-200 text-slate-700 font-medium px-4 py-2.5 rounded-xl text-sm transition-all border border-slate-300 cursor-pointer">
          <svg class="w-4 h-4 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z"/></svg>
          Cetak Laporan
        </button>
      </div>
    </div>

    <!-- Summary Performance Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-xs">
        <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">Total Booking Periodik</span>
        <div class="text-2xl font-bold text-slate-800 mt-1">{{ summaryMetrics.totalBookings }}</div>
      </div>

      <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-xs">
        <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">Pencairan SBG (Deal)</span>
        <div class="text-2xl font-bold text-emerald-600 mt-1">{{ summaryMetrics.totalSbg }}</div>
      </div>

      <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-xs">
        <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">Conversion Rate</span>
        <div class="text-2xl font-bold text-amber-600 mt-1">{{ summaryMetrics.conversionRate }}%</div>
      </div>

      <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-xs">
        <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">Estimasi UP Terkumpul</span>
        <div class="text-xl font-bold text-blue-600 mt-1">Rp {{ formatNum(summaryMetrics.totalEstimatedUp) }}</div>
      </div>
    </div>

    <!-- Main Report Card -->
    <div class="bg-white rounded-2xl border border-slate-200/80 p-6 shadow-xs space-y-6">
      <!-- Tabs -->
      <div class="flex gap-6 border-b border-slate-200">
        <button
          v-for="t in tabs"
          :key="t"
          @click="activeTab = t"
          class="pb-3 text-sm font-medium transition-all cursor-pointer relative"
          :class="activeTab === t ? 'text-amber-600 border-b-2 border-amber-500 font-semibold' : 'text-slate-500 hover:text-slate-800'"
        >
          {{ t }}
        </button>
      </div>

      <!-- Filter Controls Bar -->
      <div class="flex flex-col sm:flex-row items-end gap-4 bg-slate-50 p-4 rounded-xl border border-slate-200">
        <div class="w-full sm:w-48">
          <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1" for="month-select">Bulan</label>
          <select id="month-select" v-model="inputMonth" class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-800">
            <option v-for="(m, idx) in months" :key="idx" :value="idx + 1">
              {{ m }}
            </option>
          </select>
        </div>

        <div class="w-full sm:w-36">
          <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1" for="year-select">Tahun</label>
          <select id="year-select" v-model="inputYear" class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-800">
            <option v-for="y in years" :key="y" :value="y">
              {{ y }}
            </option>
          </select>
        </div>

        <div class="w-full sm:w-48">
          <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">Status Lead</label>
          <select v-model="inputStatus" class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-800">
            <option value="all">Semua Status</option>
            <option value="Leads">Leads</option>
            <option value="Prospect">Prospect</option>
            <option value="Hot Prospect">Hot Prospect</option>
            <option value="SBG">SBG (Deal)</option>
            <option value="Batal">Batal</option>
          </select>
        </div>

        <button @click="applyFilter" class="px-5 py-2.5 bg-amber-500 hover:bg-amber-600 text-white font-medium rounded-xl text-sm transition-all shadow-sm flex items-center gap-1.5 cursor-pointer">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
          Apply Filter
        </button>
      </div>

      <!-- TAB CONTENT 1: DAFTAR BOOKING -->
      <div v-if="activeTab === 'Daftar Booking'">
        <div v-if="showTable" class="overflow-x-auto border border-slate-200/80 rounded-xl">
          <table class="w-full text-left border-collapse text-sm">
            <thead>
              <tr class="bg-slate-50/80 border-b border-slate-200 text-xs font-semibold text-slate-500 uppercase tracking-wider">
                <th class="px-5 py-3.5">No</th>
                <th class="px-5 py-3.5">Nama Customer</th>
                <th class="px-5 py-3.5">Tanggal Rencana / Booking</th>
                <th class="px-5 py-3.5">Sumber Leads</th>
                <th class="px-5 py-3.5">Status Akhir</th>
                <th class="px-5 py-3.5 text-center">Aksi Detail</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-slate-700">
              <tr v-if="isLoading">
                <td colspan="6" class="text-center py-8">
                  <div class="w-6 h-6 border-2 border-amber-500 border-t-transparent rounded-full animate-spin mx-auto"></div>
                </td>
              </tr>
              <tr v-else-if="filteredBookings.length === 0">
                <td colspan="6" class="text-center py-8 text-slate-400 font-medium">
                  Tidak ada data booking/rencana untuk periode yang dipilih.
                </td>
              </tr>
              <tr v-for="(row, idx) in filteredBookings" :key="idx" class="hover:bg-amber-50/30 transition-all">
                <td class="px-5 py-3.5 text-slate-500">{{ idx + 1 }}</td>
                <td class="px-5 py-3.5 font-semibold text-slate-800">{{ row.nama }}</td>
                <td class="px-5 py-3.5 text-slate-600 whitespace-nowrap">{{ formatDate(row.tgl_rencana) }}</td>
                <td class="px-5 py-3.5">
                  <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-slate-100 text-slate-600">
                    {{ row.source || 'Walk In' }}
                  </span>
                </td>
                <td class="px-5 py-3.5">
                  <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold uppercase tracking-wide" :class="getStatusBadgeClass(row.status)">
                    {{ row.status }}
                  </span>
                </td>
                <td class="px-5 py-3.5 text-center">
                  <button @click="showDetail(row)" class="px-3 py-1.5 bg-slate-100 hover:bg-slate-200 text-slate-700 border border-slate-200 rounded-lg text-xs font-medium transition-all inline-flex items-center gap-1 cursor-pointer">
                    Detail <svg width="12" height="12" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/></svg>
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-else class="text-center py-12 border-2 border-dashed border-slate-200 rounded-xl text-slate-400 font-medium">
          Silakan pilih periode Bulan & Tahun lalu klik "Apply Filter" untuk melihat data laporan.
        </div>
      </div>

      <!-- TAB CONTENT 2: RINGKASAN PERFORMANCE -->
      <div v-else class="space-y-6">
        <h3 class="text-base font-semibold text-slate-800">Breakdown Performa Leads Berdasarkan Sumber</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="bg-slate-50 p-4 rounded-xl border border-slate-200 space-y-3">
            <h4 class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Sumber Leads Breakdown</h4>
            <div class="space-y-2">
              <div class="flex justify-between text-xs font-medium text-slate-700">
                <span>Walk In (Kunjungan Outlet)</span>
                <span>{{ getSourceCount('Walk In') }} Leads</span>
              </div>
              <div class="w-full bg-slate-200 rounded-full h-2">
                <div class="bg-amber-500 h-2 rounded-full" :style="{ width: `${getSourcePct('Walk In')}%` }"></div>
              </div>

              <div class="flex justify-between text-xs font-medium text-slate-700 pt-1">
                <span>Media Sosial (FB / IG / WA)</span>
                <span>{{ getSourceCount('Media Sosial') }} Leads</span>
              </div>
              <div class="w-full bg-slate-200 rounded-full h-2">
                <div class="bg-blue-500 h-2 rounded-full" :style="{ width: `${getSourcePct('Media Sosial')}%` }"></div>
              </div>

              <div class="flex justify-between text-xs font-medium text-slate-700 pt-1">
                <span>Referensi Pelanggan</span>
                <span>{{ getSourceCount('Referensi') }} Leads</span>
              </div>
              <div class="w-full bg-slate-200 rounded-full h-2">
                <div class="bg-emerald-500 h-2 rounded-full" :style="{ width: `${getSourcePct('Referensi')}%` }"></div>
              </div>
            </div>
          </div>

          <div class="bg-slate-50 p-4 rounded-xl border border-slate-200 space-y-3">
            <h4 class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Status Konversi</h4>
            <div class="space-y-2 text-xs">
              <div class="flex justify-between text-slate-700"><span>Leads Terkumpul:</span> <strong class="font-semibold text-slate-800">{{ filteredBookings.length }}</strong></div>
              <div class="flex justify-between text-slate-700"><span>Prospek Berminat:</span> <strong class="font-semibold text-blue-600">{{ getStatusCount('Prospect') }}</strong></div>
              <div class="flex justify-between text-slate-700"><span>Hot Prospect:</span> <strong class="font-semibold text-rose-600">{{ getStatusCount('Hot Prospect') }}</strong></div>
              <div class="flex justify-between text-slate-700"><span>SBG Deal (Pencairan):</span> <strong class="font-semibold text-emerald-600">{{ getStatusCount('SBG') }}</strong></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- DETAIL MODAL -->
    <div v-if="selectedItem" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs z-50 flex items-center justify-center p-4">
      <div class="bg-white rounded-2xl max-w-md w-full p-6 shadow-2xl space-y-4">
        <div class="flex items-center justify-between border-b border-slate-100 pb-3">
          <h3 class="font-semibold text-slate-800 text-base">Detail Laporan Booking</h3>
          <button @click="selectedItem = null" class="text-slate-400 hover:text-slate-600 p-1">✕</button>
        </div>
        <div class="space-y-3 text-sm">
          <div class="flex justify-between py-1 border-b border-slate-100"><span class="text-slate-500 font-medium">Customer:</span> <span class="font-semibold text-slate-800">{{ selectedItem.nama }}</span></div>
          <div class="flex justify-between py-1 border-b border-slate-100"><span class="text-slate-500 font-medium">Status:</span> <span class="font-semibold" :class="getStatusBadgeClass(selectedItem.status)">{{ selectedItem.status }}</span></div>
          <div class="flex justify-between py-1 border-b border-slate-100"><span class="text-slate-500 font-medium">Sumber:</span> <span class="text-slate-700">{{ selectedItem.source || 'Walk In' }}</span></div>
          <div class="flex justify-between py-1 border-b border-slate-100"><span class="text-slate-500 font-medium">Tanggal:</span> <span class="text-slate-700">{{ formatDate(selectedItem.tgl_rencana) }}</span></div>
          <div class="py-1">
            <span class="text-slate-500 font-medium block mb-1">Catatan Keterangan:</span>
            <p class="p-3 bg-slate-50 border border-slate-200 rounded-xl text-slate-700 text-xs leading-relaxed whitespace-pre-wrap">{{ selectedItem.keterangan || '-' }}</p>
          </div>
        </div>
        <div class="flex justify-end pt-2">
          <button @click="selectedItem = null" class="px-5 py-2 border border-slate-300 text-slate-600 rounded-xl text-sm font-medium">Tutup</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'

const tabs = ['Daftar Booking', 'Ringkasan Performance']
const activeTab = ref('Daftar Booking')
const isLoading = ref(false)
const selectedItem = ref<any>(null)

const months = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
]
const currentYear = new Date().getFullYear()
const years = Array.from({ length: 5 }, (_, i) => currentYear - i)

const inputMonth = ref(new Date().getMonth() + 1)
const inputYear = ref(currentYear)
const inputStatus = ref('all')

const showTable = ref(true)
const allBookings = ref<any[]>([])
const filteredBookings = ref<any[]>([])

const summaryMetrics = reactive({
  totalBookings: 0,
  totalSbg: 0,
  conversionRate: 0,
  totalEstimatedUp: 0
})

onMounted(() => {
  applyFilter()
})

const applyFilter = async () => {
  showTable.value = true
  isLoading.value = true

  const { data } = await useApi('/inquiry/filter', {
    params: { start: 0, length: 300 }
  })
  isLoading.value = false

  const rawList = Array.isArray(data) ? data : (data?.data || [])

  if (rawList) {
    const list = rawList.map((row: any) => ({
      ...row,
      nama: row.nama || row.nm_customer || row.name || 'Customer',
      status: row.nama_hslaktiv || row.status || 'Leads',
      source: row.nm_ref_cust || row.id_sumbercust || row.source || 'Walk In',
      keterangan: row.ket_rencana || row.ket_aktivitas || row.keterangan || '-'
    }))
    allBookings.value = list
    filteredBookings.value = list.filter((row: any) => {
      if (!row.tgl_rencana) return false
      const d = new Date(row.tgl_rencana)
      const matchDate = (d.getMonth() + 1 === inputMonth.value && d.getFullYear() === inputYear.value)
      const matchStatus = inputStatus.value === 'all' || (row.status || '').toLowerCase() === inputStatus.value.toLowerCase()
      return matchDate && matchStatus
    })

    // Calculate Summary Metrics
    summaryMetrics.totalBookings = filteredBookings.value.length
    summaryMetrics.totalSbg = filteredBookings.value.filter(r => (r.status || '').toLowerCase().includes('sbg')).length
    summaryMetrics.conversionRate = summaryMetrics.totalBookings > 0 ? Math.round((summaryMetrics.totalSbg / summaryMetrics.totalBookings) * 100) : 0
    summaryMetrics.totalEstimatedUp = summaryMetrics.totalSbg * 5500000 // Estimated average UP
  } else {
    filteredBookings.value = []
    summaryMetrics.totalBookings = 0
    summaryMetrics.totalSbg = 0
    summaryMetrics.conversionRate = 0
    summaryMetrics.totalEstimatedUp = 0
  }
}

const showDetail = (item: any) => {
  selectedItem.value = item
}

const printReport = () => {
  if (typeof window !== 'undefined') {
    window.print()
  }
}

const getSourceCount = (src: string) => {
  return filteredBookings.value.filter(b => (b.source || 'Walk In').toLowerCase() === src.toLowerCase()).length
}

const getSourcePct = (src: string) => {
  if (!filteredBookings.value.length) return 0
  return Math.round((getSourceCount(src) / filteredBookings.value.length) * 100)
}

const getStatusCount = (st: string) => {
  return filteredBookings.value.filter(b => (b.status || '').toLowerCase().includes(st.toLowerCase())).length
}

const getStatusBadgeClass = (status: string) => {
  switch (status) {
    case 'Leads': return 'bg-amber-100 text-amber-800 border border-amber-200'
    case 'Prospect': return 'bg-blue-100 text-blue-800 border border-blue-200'
    case 'Hot Prospect': return 'bg-rose-100 text-rose-800 border border-rose-200'
    case 'SBG': return 'bg-emerald-100 text-emerald-800 border border-emerald-200'
    case 'Batal': return 'bg-slate-100 text-slate-600 border border-slate-200'
    default: return 'bg-slate-100 text-slate-600'
  }
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

const formatNum = (val: number) => {
  if (!val && val !== 0) return '0'
  return new Intl.NumberFormat('id-ID').format(val)
}
</script>
