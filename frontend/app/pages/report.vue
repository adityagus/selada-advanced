<template>
  <div class="report-container">
    <!-- Header -->
    <div class="report-header mb-8">
      <h1>Report</h1>
      <p class="text-muted">Laporan data aktivitas booking berdasarkan bulan dan tahun.</p>
    </div>

    <!-- Main Report Card -->
    <div class="report-card card">
      <!-- Tabs -->
      <div class="tabs-bar mb-6">
        <button
          v-for="t in tabs"
          :key="t"
          @click="activeTab = t"
          class="tab-btn"
          :class="{ active: activeTab === t }"
        >
          {{ t }}
        </button>
      </div>

      <!-- Filter Controls -->
      <div class="filter-controls mb-6">
        <div class="filter-group">
          <label class="form-label" for="month-select">Bulan</label>
          <select id="month-select" v-model="inputMonth" class="form-select select-inline">
            <option v-for="(m, idx) in months" :key="idx" :value="idx + 1">
              {{ m }}
            </option>
          </select>
        </div>

        <div class="filter-group">
          <label class="form-label" for="year-select">Tahun</label>
          <select id="year-select" v-model="inputYear" class="form-select select-inline">
            <option v-for="y in years" :key="y" :value="y">
              {{ y }}
            </option>
          </select>
        </div>

        <button @click="applyFilter" class="btn btn-primary filter-btn">
          🔎 Apply Filter
        </button>
      </div>

      <!-- Tab Content: Booking List -->
      <div v-if="activeTab === 'Daftar Booking'">
        <div v-if="showTable" class="table-responsive">
          <table class="report-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Nama Customer</th>
                <th>Tanggal Rencana / Booking</th>
                <th>Status</th>
                <th class="text-center">Action</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="isLoading">
                <td colspan="5" class="text-center py-6">
                  <div class="spinner-small centered"></div>
                </td>
              </tr>
              <tr v-else-if="filteredBookings.length === 0">
                <td colspan="5" class="text-center py-6 text-muted">
                  Tidak ada data booking/rencana untuk periode yang dipilih.
                </td>
              </tr>
              <tr v-for="(row, idx) in filteredBookings" :key="idx">
                <td>{{ idx + 1 }}</td>
                <td><strong>{{ row.nama }}</strong></td>
                <td>{{ formatDate(row.tgl_rencana) }}</td>
                <td>
                  <span class="badge" :class="getStatusBadgeClass(row.status)">
                    {{ row.status }}
                  </span>
                </td>
                <td>
                  <div class="action-buttons-cell">
                    <button @click="showDetail(row)" class="btn btn-secondary btn-small">Detail 👁️</button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-else class="empty-state text-muted py-6">
          Silakan pilih periode Bulan & Tahun lalu klik "Apply Filter" untuk melihat data laporan.
        </div>
      </div>

      <!-- Tab Content: Lainnya -->
      <div v-else class="empty-state text-muted py-6">
        Laporan statistik lainnya akan segera hadir.
      </div>
    </div>

    <!-- Detail Modal -->
    <div v-if="selectedItem" class="modal-backdrop">
      <div class="modal-content card detail-modal">
        <div class="modal-header">
          <h3>Detail Laporan Booking</h3>
          <button @click="selectedItem = null" class="close-btn">✕</button>
        </div>
        <div class="detail-grid">
          <div class="detail-row"><span class="detail-lbl">Customer:</span> <span class="detail-val">{{ selectedItem.nama }}</span></div>
          <div class="detail-row"><span class="detail-lbl">Status:</span> <span class="badge" :class="getStatusBadgeClass(selectedItem.status)">{{ selectedItem.status }}</span></div>
          <div class="detail-row"><span class="detail-lbl">Source:</span> <span class="detail-val">{{ selectedItem.source || 'Walk In' }}</span></div>
          <div class="detail-row"><span class="detail-lbl">Tanggal:</span> <span class="detail-val">{{ formatDate(selectedItem.tgl_rencana) }}</span></div>
          <div class="detail-row"><span class="detail-lbl">Keterangan:</span> <span class="detail-val remarks">{{ selectedItem.keterangan || '-' }}</span></div>
        </div>
        <div class="modal-actions">
          <button @click="selectedItem = null" class="btn btn-secondary">Tutup</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useApi } from '~/composables/useApi'

const tabs = ['Daftar Booking', 'Statistik Lainnya']
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

const selectedMonth = ref<number | null>(null)
const selectedYear = ref<number | null>(null)
const showTable = ref(false)

const filteredBookings = ref<any[]>([])

const applyFilter = async () => {
  selectedMonth.value = inputMonth.value
  selectedYear.value = inputYear.value
  showTable.value = true
  isLoading.value = true

  // Fetch from backend
  const { data } = await useApi('/inquiry/filter', {
    params: { start: 0, length: 150 }
  })
  isLoading.value = false

  if (data && data.data) {
    // Filter local records by chosen month and year
    filteredBookings.value = data.data.filter((row: any) => {
      if (!row.tgl_rencana) return false
      const d = new Date(row.tgl_rencana)
      return (
        d.getMonth() + 1 === selectedMonth.value &&
        d.getFullYear() === selectedYear.value
      )
    })
  } else {
    filteredBookings.value = []
  }
}

const showDetail = (item: any) => {
  selectedItem.value = item
}

const getStatusBadgeClass = (status: string) => {
  switch (status) {
    case 'Leads': return 'badge-warning'
    case 'Prospect': return 'badge-info'
    case 'Hot Prospect': return 'badge-danger'
    case 'SBG': return 'badge-success'
    case 'Batal': return 'badge-danger'
    default: return 'badge-info'
  }
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', { day: '2-digit', month: '2-digit', year: 'numeric' })
}
</script>

<style scoped>
.report-container {
  max-width: 1200px;
  margin: 0 auto;
}

.report-header h1 {
  font-size: 1.75rem;
  font-weight: 700;
}

.tabs-bar {
  display: flex;
  gap: 1.5rem;
  border-bottom: 1px solid var(--border-color);
  margin-bottom: 2rem;
}

.tab-btn {
  background: none;
  border: none;
  padding: 0.75rem 0.5rem;
  color: var(--text-muted);
  font-weight: 600;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  transition: var(--transition-smooth);
  font-family: var(--font-body);
}

.tab-btn:hover, .tab-btn.active {
  color: var(--selada-400);
}

.tab-btn.active {
  border-bottom-color: var(--selada-400);
}

.filter-controls {
  display: flex;
  gap: 1.5rem;
  align-items: flex-end;
  flex-wrap: wrap;
  background: #f9fafb;
  padding: 1.25rem;
  border-radius: 12px;
  border: 1px solid var(--border-color);
}

.filter-group {
  display: flex;
  flex-direction: column;
}

.select-inline {
  min-width: 160px;
}

.filter-btn {
  margin-bottom: 0.25rem;
}

.table-responsive {
  width: 100%;
  overflow-x: auto;
  margin-top: 1.5rem;
}

.report-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.report-table th, .report-table td {
  padding: 0.85rem 1.25rem;
  border-bottom: 1px solid var(--border-color);
}

.report-table th {
  background: #f9fafb;
  color: var(--text-muted);
  font-weight: 600;
  text-transform: uppercase;
  font-size: 0.75rem;
}

.report-table tr:hover td {
  background: #fffbeb;
}

.text-center { text-align: center; }
.py-6 { padding-top: 1.5rem; padding-bottom: 1.5rem; }

.empty-state {
  text-align: center;
  padding: 3rem 1.5rem;
  border: 1px dashed var(--border-color);
  border-radius: 12px;
  background: #fafafa;
}

.action-buttons-cell {
  display: flex;
  justify-content: center;
}

.btn-small {
  padding: 0.35rem 0.75rem;
  font-size: 0.8rem;
}

/* Modal */
.modal-backdrop {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  padding: 1rem;
}

.modal-content {
  width: 100%;
  max-width: 560px;
  max-height: 90vh;
  overflow-y: auto;
  animation: fadeIn 0.3s ease-out;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.modal-header h3 {
  font-size: 1.15rem;
  font-weight: 700;
}

.close-btn {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: #f3f4f6;
  border: none;
  cursor: pointer;
  font-size: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: var(--transition-smooth);
  color: var(--text-muted);
}

.close-btn:hover {
  background: #e5e7eb;
  color: var(--text-main);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 1rem;
}

.detail-modal {
  max-width: 480px;
}

.detail-grid {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-bottom: 1.5rem;
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 1rem;
}

.detail-row {
  display: flex;
}

.detail-lbl {
  width: 100px;
  font-weight: 600;
  color: var(--text-muted);
}

.detail-val {
  flex-grow: 1;
  color: var(--text-main);
}

.remarks {
  white-space: pre-wrap;
  background: #f9fafb;
  padding: 0.75rem;
  border-radius: 6px;
  border: 1px solid var(--border-color);
}

.centered {
  margin: 0 auto;
}
</style>
