<template>
  <div class="taksiran-container">
    <!-- Header -->
    <div class="taksiran-header mb-8">
      <h1>Cek Taksiran Gadai</h1>
      <p class="text-muted">Silakan isi formulir di bawah ini untuk mensimulasikan nilai taksiran barang gadai Anda.</p>
    </div>

    <!-- Main Calculator Form -->
    <div class="taksiran-card card">
      <form @submit.prevent="handleSubmit" class="calculator-form">
        <!-- Jenis Barang Toggle -->
        <div class="form-group">
          <label class="form-label">Jenis Barang</label>
          <div class="radio-group">
            <label class="radio-label">
              <input type="radio" v-model="form.jenisBarang" value="umum" @change="resetFields" /> Umum (Kendaraan / Elektronik)
            </label>
            <label class="radio-label">
              <input type="radio" v-model="form.jenisBarang" value="emas" @change="resetFields" /> Emas / Perhiasan
            </label>
          </div>
        </div>

        <!-- Category & Grade -->
        <div class="form-row">
          <div class="form-group flex-1">
            <label class="form-label" for="kategori-select">Kategori Barang</label>
            <select id="kategori-select" v-model="form.kategori" class="form-select" required>
              <option value="">- Pilih Kategori -</option>
              <option v-if="form.jenisBarang === 'umum'" value="Kendaraan">Kendaraan</option>
              <option v-if="form.jenisBarang === 'umum'" value="Elektronik">Elektronik</option>
              <option v-if="form.jenisBarang === 'emas'" value="Perhiasan">Perhiasan Emas</option>
              <option v-if="form.jenisBarang === 'emas'" value="Lainnya">Lain-lain</option>
            </select>
          </div>

          <div class="form-group flex-1">
            <label class="form-label" for="grade-select">Grade Kondisi</label>
            <select id="grade-select" v-model="form.grade" @change="calculateLoan" class="form-select" required>
              <option value="A">Grade A (90% Taksiran)</option>
              <option value="B">Grade B (85% Taksiran)</option>
              <option value="C">Grade C (80% Taksiran)</option>
              <option value="D">Grade D (75% Taksiran)</option>
            </select>
          </div>
        </div>

        <!-- Kode Barang & Search Button -->
        <div class="form-row">
          <div class="form-group flex-1 search-group">
            <label class="form-label" for="kode-input">Kode Barang / Karat</label>
            <div class="input-search-wrapper">
              <input
                id="kode-input"
                v-model="form.kodeBarang"
                type="text"
                placeholder="Ketik kode barang (cth: Beat, iPhone, atau Karat 24)"
                class="form-input"
                required
              />
              <button type="button" @click="searchBarang" class="btn btn-secondary search-action-btn" :disabled="isSearching">
                <span v-if="isSearching" class="spinner-small"></span>
                <span v-else>🔍 Cari</span>
              </button>
            </div>
            <span v-if="searchError" class="error-hint">⚠️ {{ searchError }}</span>
          </div>

          <div class="form-group flex-1">
            <label class="form-label" for="nama-input">Nama Deskripsi Barang</label>
            <input
              id="nama-input"
              v-model="form.namaBarang"
              type="text"
              class="form-input disabled-field"
              readonly
              disabled
            />
          </div>
        </div>

        <!-- Nilai Taksir & Nilai Pinjaman -->
        <div class="form-row">
          <div class="form-group flex-grow">
            <label class="form-label" for="taksir-input">Nilai Taksir Barang</label>
            <input
              id="taksir-input"
              v-model="displayTaksir"
              type="text"
              class="form-input disabled-field highlight-field"
              readonly
              disabled
            />
          </div>

          <div class="form-group flex-grow">
            <label class="form-label" for="pinjam-input">Nilai Maksimal Pinjaman</label>
            <input
              id="pinjam-input"
              v-model="displayPinjaman"
              type="text"
              class="form-input disabled-field highlight-field"
              readonly
              disabled
            />
          </div>
        </div>

        <!-- Sewa Modal, Admin, Estimasi Net Pencairan -->
        <div class="form-row">
          <div class="form-group flex-1">
            <label class="form-label">Estimasi Sewa Modal (Per 15 Hari)</label>
            <input
              v-model="displaySewaModal"
              type="text"
              class="form-input disabled-field"
              readonly
              disabled
            />
          </div>
          <div class="form-group flex-1">
            <label class="form-label">Biaya Administrasi</label>
            <input
              v-model="displayBiayaAdmin"
              type="text"
              class="form-input disabled-field"
              readonly
              disabled
            />
          </div>
          <div class="form-group flex-grow">
            <label class="form-label">Estimasi Pencairan Net</label>
            <input
              v-model="displayNetPencairan"
              type="text"
              class="form-input disabled-field highlight-success"
              readonly
              disabled
            />
          </div>
        </div>

        <!-- Submit & Save Estimation -->
        <div class="form-actions-box">
          <button type="submit" class="btn btn-primary submit-estimate-btn" :disabled="isSubmitting || form.nilaiTaksir <= 0">
            <span v-if="isSubmitting" class="spinner-small"></span>
            <span v-else>Simpan Transaksi Taksiran 💾</span>
          </button>
        </div>

        <!-- Success/Error alert banner -->
        <div v-if="statusMessage" class="status-banner" :class="statusClass">
          {{ statusMessage }}
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useApi } from '~/composables/useApi'

const form = ref({
  jenisBarang: 'umum',
  kategori: '',
  grade: 'A',
  kodeBarang: '',
  namaBarang: '',
  nilaiTaksir: 0,
  nilaiPinjaman: 0,
  sewaModal: 0,
  biayaAdmin: 0,
  estimasiPencairan: 0
})

const isSearching = ref(false)
const isSubmitting = ref(false)
const searchError = ref('')
const statusMessage = ref('')
const statusClass = ref('')

const displayTaksir = computed(() => formatRupiah(form.value.nilaiTaksir))
const displayPinjaman = computed(() => formatRupiah(form.value.nilaiPinjaman))
const displaySewaModal = computed(() => formatRupiah(form.value.sewaModal))
const displayBiayaAdmin = computed(() => formatRupiah(form.value.biayaAdmin))
const displayNetPencairan = computed(() => formatRupiah(form.value.estimasiPencairan))

const resetFields = () => {
  form.value.kategori = ''
  form.value.kodeBarang = ''
  form.value.namaBarang = ''
  form.value.nilaiTaksir = 0
  form.value.nilaiPinjaman = 0
  form.value.sewaModal = 0
  form.value.biayaAdmin = 0
  form.value.estimasiPencairan = 0
  searchError.value = ''
  statusMessage.value = ''
}

const searchBarang = async () => {
  searchError.value = ''
  isSearching.value = true
  statusMessage.value = ''

  const searchPath = form.value.jenisBarang === 'umum' ? '/taksiran/barang-umum' : '/taksiran/barang-emas'
  
  const { data, error } = await useApi(searchPath, {
    params: {
      search: form.value.kodeBarang,
      start: 0,
      length: 1
    }
  })

  isSearching.value = false

  if (error || !data || !data.data || data.data.length === 0) {
    searchError.value = 'Barang tidak ditemukan. Coba keyword lain.'
    form.value.namaBarang = 'Barang Tidak Terdaftar'
    form.value.nilaiTaksir = 0
    calculateLoan()
    return
  }

  // Found item
  const item = data.data[0]
  form.value.namaBarang = item.nama_barang || item.tipe_barang || 'Perhiasan Emas Acuan'
  
  // Gold estimation vs general estimation
  if (form.value.jenisBarang === 'emas') {
    // Emas taksiran depends on karat/gold rate (dummy model formula)
    form.value.nilaiTaksir = Number(item.harga_pasar) || 850000
  } else {
    form.value.nilaiTaksir = Number(item.harga_taksir) || Number(item.harga_pasar) || 2000000
  }

  calculateLoan()
}

const calculateLoan = async () => {
  if (form.value.nilaiTaksir <= 0) {
    form.value.nilaiPinjaman = 0
    form.value.sewaModal = 0
    form.value.biayaAdmin = 0
    form.value.estimasiPencairan = 0
    return
  }

  // Calculate loan based on condition grade
  let multiplier = 0.90
  if (form.value.grade === 'B') multiplier = 0.85
  if (form.value.grade === 'C') multiplier = 0.80
  if (form.value.grade === 'D') multiplier = 0.75

  form.value.nilaiPinjaman = Math.round(form.value.nilaiTaksir * multiplier)

  // Fetch admin fee dynamically from Go API
  const { data } = await useApi('/taksiran/biaya-admin', {
    params: { pinjaman: form.value.nilaiPinjaman }
  })
  
  form.value.biayaAdmin = data?.biaya_admin || 15000 // Fallback

  // Sewa modal is typically 1.2% for general, 1.5% for gold per 15-day cycle
  const sewaRate = form.value.jenisBarang === 'emas' ? 0.015 : 0.012
  form.value.sewaModal = Math.round(form.value.nilaiPinjaman * sewaRate)

  // Net payout estimation
  form.value.estimasiPencairan = form.value.nilaiPinjaman - form.value.biayaAdmin
}

const handleSubmit = async () => {
  isSubmitting.value = true
  statusMessage.value = ''

  const { error } = await useApi('/taksiran/simpan', {
    method: 'POST',
    body: {
      jenis_barang: form.value.jenisBarang,
      kategori: form.value.kategori,
      grade: form.value.grade,
      kode_barang: form.value.kodeBarang,
      nama_barang: form.value.namaBarang,
      nilai_taksir: form.value.nilaiTaksir,
      nilai_pinjaman: form.value.nilaiPinjaman,
      sewa_modal: form.value.sewaModal,
      biaya_admin: form.value.biayaAdmin,
      estimasi_pencairan: form.value.estimasiPencairan
    }
  })

  isSubmitting.value = false

  if (error) {
    statusMessage.value = 'Gagal menyimpan transaksi taksiran: ' + error
    statusClass.value = 'banner-error'
  } else {
    statusMessage.value = 'Simulasi taksiran gadai berhasil disimpan sebagai rencana transaksi! 🎉'
    statusClass.value = 'banner-success'
    resetFields()
  }
}

const formatRupiah = (num: number) => {
  if (!num) return 'Rp 0'
  return 'Rp ' + num.toLocaleString('id-ID')
}
</script>

<style scoped>
.taksiran-container {
  max-width: 1000px;
  margin: 0 auto;
}

.taksiran-header h1 {
  font-size: 1.75rem;
  font-weight: 700;
}

.taksiran-card {
  padding: 2rem;
}

.radio-group {
  display: flex;
  gap: 2rem;
  margin-top: 0.5rem;
}

.radio-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  color: var(--text-main);
}

.form-row {
  display: flex;
  gap: 1.5rem;
}

.flex-1 { flex: 1; }
.flex-grow { flex-grow: 1; }

@media (max-width: 768px) {
  .form-row {
    flex-direction: column;
    gap: 0;
  }
}

.input-search-wrapper {
  display: flex;
  gap: 0.5rem;
}

.search-action-btn {
  flex-shrink: 0;
}

.disabled-field {
  background: #f9fafb;
  color: var(--text-muted);
  cursor: not-allowed;
}

.highlight-field {
  color: var(--selada-400) !important;
  font-weight: 700;
}

.highlight-success {
  color: var(--color-success) !important;
  font-weight: 700;
  background: #ecfdf5;
  border-color: #a7f3d0;
}

.error-hint {
  font-size: 0.75rem;
  color: #dc2626;
  margin-top: 0.25rem;
  display: block;
}

.form-actions-box {
  display: flex;
  justify-content: flex-end;
  margin-top: 1.5rem;
  border-top: 1px solid var(--border-color);
  padding-top: 1.5rem;
}

.submit-estimate-btn {
  padding: 0.85rem 2rem;
}

.status-banner {
  margin-top: 1.5rem;
  padding: 1rem;
  border-radius: 8px;
  text-align: center;
  font-weight: 500;
  font-size: 0.95rem;
}

.banner-success {
  background: #ecfdf5;
  border: 1px solid #a7f3d0;
  color: #059669;
}

.banner-error {
  background: #fef2f2;
  border: 1px solid #fecaca;
  color: #dc2626;
}
</style>
