<template>
  <div class="inquiry-container">
    <!-- Header -->
    <div class="inquiry-header mb-8">
      <div class="header-left-box">
        <h1>Inquiry</h1>
        <p class="text-muted">Kelola dan lacak semua inquiry aktivitas sales di satu tempat</p>
      </div>
      <div class="header-right-box">
        <button @click="openModal = true" class="btn btn-primary">
          <span>➕</span> Add Leads
        </button>
      </div>
    </div>

    <!-- Filter & Search Controls -->
    <div class="filter-card card mb-8">
      <div class="filter-inputs">
        <div class="input-col search-col">
          <label class="form-label" for="search-cust">Cari Customer</label>
          <input
            id="search-cust"
            v-model="searchQuery"
            @input="debounceFetch"
            type="text"
            placeholder="Ketik nama customer..."
            class="form-input"
          />
        </div>
        <div class="input-col status-col">
          <label class="form-label" for="status-filter">Status Customer</label>
          <select id="status-filter" v-model="selectedStatus" @change="fetchInquiries" class="form-select">
            <option value="all">Semua Status Aktif</option>
            <option value="Leads">Leads</option>
            <option value="Prospect">Prospect</option>
            <option value="Hot Prospect">Hot Prospect</option>
            <option value="Batal">Batal</option>
            <option value="SBG">SBG (Deal)</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Data Table Card -->
    <div class="table-card card">
      <div class="table-responsive">
        <table class="inquiry-table">
          <thead>
            <tr>
              <th>Tanggal</th>
              <th>Nama Customer</th>
              <th>Source</th>
              <th>Status</th>
              <th>Keterangan</th>
              <th class="text-center">Action</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="isLoadingTable">
              <td colspan="6" class="text-center py-6">
                <div class="spinner-small centered"></div>
              </td>
            </tr>
            <tr v-else-if="inquiries.length === 0">
              <td colspan="6" class="text-center py-6 text-muted">
                Tidak ada data inquiry ditemukan.
              </td>
            </tr>
            <tr v-for="item in inquiries" :key="item.id_rencana">
              <td>{{ formatDate(item.tgl_rencana) }}</td>
              <td><strong>{{ item.nama }}</strong></td>
              <td><span class="badge badge-info">{{ item.source || 'Walk In' }}</span></td>
              <td>
                <span class="badge" :class="getStatusBadgeClass(item.status)">
                  {{ item.status }}
                </span>
              </td>
              <td class="text-truncate">{{ item.keterangan || '-' }}</td>
              <td>
                <div class="action-buttons-cell">
                  <button @click="showDetail(item)" class="btn-icon btn-detail" title="Detail">👁️</button>
                  <button @click="handleDelete(item.id_rencana)" class="btn-icon btn-delete" title="Hapus">🗑️</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Form Wizard: Buat Rencana Baru -->
    <div v-if="openModal" class="modal-backdrop">
      <div class="modal-content card">
        <div class="modal-header">
          <h3>Buat Rencana Aktivitas Baru</h3>
          <button @click="closeFormModal" class="close-btn">✕</button>
        </div>

        <form @submit.prevent="handleSubmitRencana" class="modal-form">
          <!-- Step 1: Customer Type -->
          <div class="form-group">
            <label class="form-label">Tipe Customer</label>
            <div class="radio-group">
              <label class="radio-label">
                <input type="radio" v-model="form.cust" value="Baru" /> Baru
              </label>
              <label class="radio-label">
                <input type="radio" v-model="form.cust" value="Lama" /> Pelanggan Lama
              </label>
            </div>
          </div>

          <!-- Step 2: Customer Fields -->
          <div v-if="form.cust === 'Baru'" class="form-row">
            <div class="form-group flex-1">
              <label class="form-label" for="nm">Nama Lengkap</label>
              <input id="nm" v-model="form.nm" type="text" placeholder="Nama customer baru" class="form-input" required />
            </div>
            <div class="form-group flex-1">
              <label class="form-label" for="hp">No. Handphone</label>
              <input id="hp" v-model="form.hp" type="text" placeholder="08xxxxxxxxx" class="form-input" required />
            </div>
          </div>

          <div v-else class="form-group">
            <label class="form-label" for="custlama">Pilih Customer</label>
            <select id="custlama" v-model="form.custlama" class="form-select" required>
              <option value="">- Pilih Customer Existing -</option>
              <option v-for="c in existingCustomers" :key="c.id_customer" :value="c.id_customer">
                {{ c.nama }} ({{ c.hp }})
              </option>
            </select>
          </div>

          <div class="form-row">
            <div class="form-group flex-1">
              <label class="form-label" for="sumbercust">Sumber leads</label>
              <select id="sumbercust" v-model="form.sumbercust" class="form-select">
                <option value="Walk In">Walk In</option>
                <option value="Media Sosial">Media Sosial</option>
                <option value="Event">Event</option>
                <option value="Referensi">Referensi</option>
              </select>
            </div>
            <div class="form-group flex-1">
              <label class="form-label" for="hslaktiv">Target Status Awal</label>
              <select id="hslaktiv" v-model.number="form.hslaktiv" class="form-select" required>
                <option :value="1">Leads</option>
                <option :value="2">Prospect</option>
                <option :value="3">Hot Prospect</option>
              </select>
            </div>
          </div>

          <div class="form-group">
            <label class="form-label" for="ket">Catatan Rencana Aktivitas</label>
            <textarea id="ket" v-model="form.ket" rows="3" placeholder="Ketik catatan aktivitas atau janji temu..." class="form-input form-textarea"></textarea>
          </div>

          <div v-if="modalError" class="modal-error">
            <span>⚠️</span> {{ modalError }}
          </div>

          <div class="modal-actions">
            <button type="button" @click="closeFormModal" class="btn btn-secondary">Batal</button>
            <button type="submit" class="btn btn-primary" :disabled="isSubmitting">
              <span v-if="isSubmitting" class="spinner-small"></span>
              <span v-else>Simpan Rencana 💾</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Detail Modal -->
    <div v-if="selectedItem" class="modal-backdrop">
      <div class="modal-content card detail-modal">
        <div class="modal-header">
          <h3>Detail Aktivitas Customer</h3>
          <button @click="selectedItem = null" class="close-btn">✕</button>
        </div>
        <div class="detail-grid">
          <div class="detail-row"><span class="detail-lbl">Nama:</span> <span class="detail-val">{{ selectedItem.nama }}</span></div>
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
import { ref, reactive, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'

const inquiries = ref<any[]>([])
const searchQuery = ref('')
const selectedStatus = ref('all')
const isLoadingTable = ref(false)
const openModal = ref(false)
const isSubmitting = ref(false)
const modalError = ref('')
const existingCustomers = ref<any[]>([])
const selectedItem = ref<any>(null)

const form = reactive({
  cust: 'Baru',
  custlama: '',
  nm: '',
  hp: '',
  sumbercust: 'Walk In',
  hslaktiv: 1,
  katcust: '1',
  ket: '',
  aktiv: 1
})

onMounted(() => {
  fetchInquiries()
  fetchExistingCustomers()
})

const fetchInquiries = async () => {
  isLoadingTable.value = true
  const { data } = await useApi('/inquiry/filter', {
    params: {
      status: selectedStatus.value,
      search: searchQuery.value,
      start: 0,
      length: 100
    }
  })
  isLoadingTable.value = false
  if (data && data.data) {
    inquiries.value = data.data
  }
}

const fetchExistingCustomers = async () => {
  const { data } = await useApi('/inquiry/filter', {
    params: { start: 0, length: 100 }
  })
  if (data && data.data) {
    const map = new Map()
    for (const item of data.data) {
      map.set(item.nama, { id_customer: item.id_rencana, nama: item.nama, hp: item.source || '0812' })
    }
    existingCustomers.value = Array.from(map.values())
  }
}

let debounceTimer: any = null
const debounceFetch = () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    fetchInquiries()
  }, 400)
}

const handleSubmitRencana = async () => {
  isSubmitting.value = true
  modalError.value = ''

  const payload: any = {
    cust: form.cust,
    sumbercust: form.sumbercust,
    hslaktiv: form.hslaktiv,
    katcust: form.katcust,
    ket: form.ket,
    aktiv: form.aktiv
  }

  if (form.cust === 'Lama') {
    payload.custlama = Number(form.custlama)
  } else {
    payload.nm = form.nm
    payload.hp = form.hp
  }

  const { error } = await useApi('/rencana/input', {
    method: 'POST',
    body: payload
  })

  isSubmitting.value = false

  if (error) {
    modalError.value = error
    return
  }

  fetchInquiries()
  closeFormModal()
}

const handleDelete = async (id: number) => {
  if (confirm('Apakah Anda yakin ingin menghapus rencana aktivitas ini?')) {
    const { error } = await useApi(`/rencana/${id}`, { method: 'DELETE' })
    if (!error) {
      fetchInquiries()
    }
  }
}

const showDetail = (item: any) => {
  selectedItem.value = item
}

const closeFormModal = () => {
  openModal.value = false
  modalError.value = ''
  form.cust = 'Baru'
  form.custlama = ''
  form.nm = ''
  form.hp = ''
  form.ket = ''
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
.inquiry-container {
  max-width: 1200px;
  margin: 0 auto;
}

.inquiry-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 1rem;
}

.inquiry-header h1 {
  font-size: 1.75rem;
  font-weight: 700;
}

.filter-card {
  padding: 1.25rem 1.5rem;
}

.filter-inputs {
  display: flex;
  gap: 1.5rem;
  flex-wrap: wrap;
}

.input-col {
  display: flex;
  flex-direction: column;
}

.search-col {
  flex: 2;
  min-width: 250px;
}

.status-col {
  flex: 1;
  min-width: 180px;
}

.table-card {
  padding: 0;
  overflow: hidden;
  border-radius: 12px;
}

.table-responsive {
  width: 100%;
  overflow-x: auto;
}

.inquiry-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.inquiry-table th, .inquiry-table td {
  padding: 0.85rem 1.25rem;
  border-bottom: 1px solid var(--border-color);
}

.inquiry-table th {
  background: #f9fafb;
  color: var(--text-muted);
  font-weight: 600;
  text-transform: uppercase;
  font-size: 0.75rem;
  letter-spacing: 0.5px;
}

.inquiry-table tr:last-child td {
  border-bottom: none;
}

.inquiry-table tr:hover td {
  background: #fffbeb;
}

.text-truncate {
  max-width: 200px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.text-center { text-align: center; }
.py-6 { padding-top: 1.5rem; padding-bottom: 1.5rem; }

.action-buttons-cell {
  display: flex;
  gap: 0.5rem;
  justify-content: center;
}

.btn-icon {
  background: #f3f4f6;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: var(--transition-smooth);
  font-size: 0.85rem;
}

.btn-icon:hover {
  background: #e5e7eb;
  transform: scale(1.05);
}

.btn-delete:hover {
  background: #fee2e2;
  border-color: #fecaca;
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

.radio-group {
  display: flex;
  gap: 1.5rem;
  margin-top: 0.5rem;
}

.radio-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
}

.form-row {
  display: flex;
  gap: 1rem;
}

.flex-1 { flex: 1; }

@media (max-width: 576px) {
  .form-row {
    flex-direction: column;
    gap: 0;
  }
}

.modal-error {
  background: #fef2f2;
  border: 1px solid #fecaca;
  color: #dc2626;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  margin-bottom: 1.25rem;
  font-size: 0.85rem;
}

/* Detail */
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
