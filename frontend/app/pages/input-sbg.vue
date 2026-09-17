<template>
  <div class="space-y-6 pb-12 font-sans max-w-5xl mx-auto">
    <!-- Page Header (Matches v_inputdo.php) -->
    <div class="border-b border-slate-200 pb-4">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 class="text-2xl font-bold text-slate-800 tracking-tight">Input SBG</h1>
          <p class="text-slate-500 text-sm mt-0.5">Input realisasi DO / SBG dari transaksi resmi Sopiga Pro.</p>
        </div>
        <div>
          <NuxtLink
            to="/inquiry"
            class="inline-flex items-center gap-2 bg-slate-600 hover:bg-slate-700 active:scale-95 text-white font-medium px-4 py-2 rounded-xl text-sm transition-all shadow-xs"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
            Kembali
          </NuxtLink>
        </div>
      </div>
    </div>

    <!-- Customer Context Notice (if opened from inquiry) -->
    <div v-if="inqContext.nama" class="bg-blue-50 border border-blue-200 rounded-xl p-4 flex flex-col sm:flex-row sm:items-center justify-between gap-3 text-sm">
      <div class="flex items-center gap-3">
        <div class="w-9 h-9 rounded-lg bg-blue-600 text-white flex items-center justify-center font-bold text-sm shrink-0">
          👤
        </div>
        <div>
          <div class="text-xs text-blue-600 font-semibold uppercase tracking-wider">Customer Inquiry</div>
          <div class="font-bold text-slate-800 text-base">{{ inqContext.nama }} <span v-if="inqContext.hp" class="text-slate-500 font-normal text-sm">({{ inqContext.hp }})</span></div>
          <div v-if="inqContext.ket_rencana" class="text-xs text-slate-600 mt-0.5">Rencana: {{ inqContext.ket_rencana }}</div>
        </div>
      </div>
      <div class="text-xs font-semibold px-3 py-1 rounded-full bg-blue-100 text-blue-800 border border-blue-200 self-start sm:self-center">
        {{ inqContext.nama_hslaktiv || 'Hot Prospect' }}
      </div>
    </div>

    <!-- Main Card (Matches col-md-6 Card in v_inputdo.php) -->
    <div class="max-w-2xl">
      <div class="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden">
        <!-- Card Header (bg-primary / text-white) -->
        <div class="bg-blue-600 text-white px-6 py-4 flex items-center gap-2">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          <h2 class="font-bold text-base">Data SBG</h2>
        </div>

        <!-- Card Body -->
        <div class="p-6 space-y-5">
          <!-- 1. No SBG (Select2 Searchable Dropdown Parity) -->
          <div class="space-y-1.5">
            <label class="block text-sm font-bold text-slate-800">
              No SBG <span class="text-rose-500">*</span>
            </label>

            <!-- Searchable Select UI -->
            <div class="relative">
              <div
                @click="toggleDropdown"
                class="w-full min-h-[42px] px-3.5 py-2 bg-white border border-slate-300 rounded-xl text-sm flex items-center justify-between cursor-pointer hover:border-blue-400 focus-within:ring-2 focus-within:ring-blue-100 transition-all"
                :class="{ 'border-blue-500 ring-2 ring-blue-100': isDropdownOpen }"
              >
                <div class="flex-1 truncate">
                  <span v-if="selectedSbg" class="font-medium text-slate-800">
                    {{ selectedSbg.text }}
                  </span>
                  <span v-else class="text-slate-400">
                    -- Pilih No SBG --
                  </span>
                </div>

                <div class="flex items-center gap-1.5 ml-2">
                  <button
                    v-if="selectedSbg"
                    @click.stop="clearSbg"
                    type="button"
                    class="text-slate-400 hover:text-slate-600 p-0.5 rounded"
                    title="Clear"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
                  </button>
                  <svg class="w-4 h-4 text-slate-400 transition-transform" :class="{ 'rotate-180': isDropdownOpen }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                  </svg>
                </div>
              </div>

              <!-- Dropdown Panel -->
              <div
                v-if="isDropdownOpen"
                class="absolute left-0 right-0 top-full mt-1.5 bg-white border border-slate-200 rounded-xl shadow-xl z-50 overflow-hidden"
              >
                <!-- Search box within dropdown -->
                <div class="p-2 border-b border-slate-100 bg-slate-50">
                  <div class="relative">
                    <input
                      ref="searchInputRef"
                      v-model="searchTerm"
                      @input="onSearchInput"
                      type="text"
                      placeholder="Ketik No SBG atau Nama Customer..."
                      class="w-full pl-8 pr-3 py-1.5 bg-white border border-slate-300 rounded-lg text-xs outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-200"
                    />
                    <svg class="w-3.5 h-3.5 text-slate-400 absolute left-2.5 top-2.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
                  </div>
                </div>

                <!-- Dropdown items list -->
                <div class="max-h-60 overflow-y-auto divide-y divide-slate-50 text-xs">
                  <div v-if="isLoadingSbg" class="py-4 text-center text-slate-400">
                    <div class="w-4 h-4 border-2 border-blue-500 border-t-transparent rounded-full animate-spin mx-auto mb-1"></div>
                    Memuat data SBG...
                  </div>
                  <div v-else-if="sbgList.length === 0" class="py-4 px-3 text-center text-slate-400">
                    Tidak ada SBG tersedia atau semua sudah terdaftar.
                  </div>
                  <button
                    v-for="item in sbgList"
                    :key="item.id"
                    @click="selectSbgItem(item)"
                    type="button"
                    class="w-full text-left px-3.5 py-2.5 hover:bg-blue-50/70 transition-colors flex flex-col gap-0.5 cursor-pointer"
                    :class="{ 'bg-blue-50 font-semibold text-blue-700': selectedSbg?.id === item.id }"
                  >
                    <div class="flex items-center justify-between">
                      <span class="font-bold text-slate-800 font-mono">{{ item.id }}</span>
                      <span class="text-[10px] px-2 py-0.5 rounded font-bold uppercase tracking-wider" :class="item.schema_type === 'CICILAN' ? 'bg-amber-100 text-amber-800' : 'bg-blue-100 text-blue-800'">
                        {{ item.schema_type || 'GADAI' }}
                      </span>
                    </div>
                    <div class="text-slate-600 font-medium">
                      {{ item.nm_customer || '-' }}
                    </div>
                    <div v-if="item.nm_produk" class="text-[11px] text-slate-400">
                      {{ item.nm_produk }}
                    </div>
                  </button>
                </div>
              </div>
            </div>

            <small class="text-slate-400 text-xs flex items-center gap-1">
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
              SBG yang sudah terdaftar tidak akan muncul kembali.
            </small>
          </div>

          <!-- 2. Jenis Transaksi (tipe3 in v_inputdo.php) -->
          <div class="space-y-1.5">
            <label class="block text-sm font-bold text-slate-800">
              Jenis Transaksi <span class="text-rose-500">*</span>
            </label>
            <select
              id="tipe3"
              v-model="form.jnstrx"
              @change="onJnsTrxChange"
              class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm font-medium text-slate-800 focus:outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-100 transition-all"
            >
              <option value="0">-- Pilih Jenis Transaksi --</option>
              <option value="1">Gadai</option>
              <option value="2">Cicilan</option>
            </select>
          </div>

          <!-- 3. Nama Produk (warna3 in v_inputdo.php) -->
          <div class="space-y-1.5">
            <label class="block text-sm font-bold text-slate-800">
              Nama Produk <span class="text-rose-500">*</span>
            </label>
            <select
              id="warna3"
              v-model="form.produk"
              class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm font-medium text-slate-800 focus:outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-100 transition-all"
            >
              <option value="0">-- Pilih Produk --</option>
              <option v-for="prod in produkList" :key="prod.kd_produk" :value="prod.kd_produk">
                {{ prod.nm_produk }}
              </option>
            </select>
          </div>

          <!-- Hidden Context Inputs matching v_inputdo.php -->
          <input type="hidden" :value="form.idrendet" />
          <input type="hidden" :value="form.ket_rencana" />
          <input type="hidden" :value="form.ket_aktivitas" />

          <!-- Alert Banner (sbg-alert in v_inputdo.php) -->
          <div
            v-if="alertState.show"
            class="p-3 rounded-xl text-xs font-semibold flex items-center justify-between gap-2"
            :class="{
              'bg-emerald-50 text-emerald-800 border border-emerald-200': alertState.type === 'success',
              'bg-rose-50 text-rose-800 border border-rose-200': alertState.type === 'danger',
              'bg-amber-50 text-amber-800 border border-amber-200': alertState.type === 'warning'
            }"
          >
            <span>{{ alertState.message }}</span>
            <button @click="alertState.show = false" type="button" class="text-slate-400 hover:text-slate-600 font-bold">&times;</button>
          </div>
        </div>

        <!-- Card Footer -->
        <div class="px-6 py-4 bg-slate-50 border-t border-slate-100 flex justify-end">
          <button
            type="button"
            id="btn-save-sbg"
            @click="onSaveClick"
            class="px-6 py-2.5 bg-blue-600 hover:bg-blue-700 active:scale-95 text-white font-bold rounded-xl text-sm transition-all shadow-xs flex items-center gap-2 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
            :disabled="!isReady || isSubmitting"
          >
            <span v-if="isSubmitting" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
            <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4" />
            </svg>
            <span>Simpan</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Modal Konfirmasi Tanggal (modalConfirmTgl in v_inputdo.php) -->
    <div
      v-if="showConfirmDateModal"
      class="fixed inset-0 bg-slate-900/50 backdrop-blur-xs z-50 flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-2xl max-w-md w-full shadow-2xl border border-slate-100 overflow-hidden animate-in fade-in zoom-in-95 duration-150">
        <!-- Modal Header (bg-warning text-white) -->
        <div class="bg-amber-500 text-white px-6 py-4 flex items-center justify-between">
          <h3 class="font-bold text-base flex items-center gap-2">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            Konfirmasi Tanggal
          </h3>
          <button @click="showConfirmDateModal = false" type="button" class="text-white hover:text-amber-100">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
          </button>
        </div>

        <!-- Modal Body -->
        <div class="p-6 space-y-2">
          <p class="text-sm font-semibold text-slate-800">{{ confirmMsg }}</p>
          <p class="text-xs text-slate-500">Apakah Anda yakin ingin tetap menyimpan data ini?</p>
        </div>

        <!-- Modal Footer -->
        <div class="px-6 py-4 bg-slate-50 border-t border-slate-100 flex justify-end gap-3">
          <button
            @click="showConfirmDateModal = false"
            type="button"
            class="px-4 py-2 border border-slate-300 text-slate-700 font-semibold rounded-xl text-xs hover:bg-slate-100 transition-colors cursor-pointer"
          >
            Batal
          </button>
          <button
            @click="onConfirmSave"
            type="button"
            id="btn-confirm-save"
            class="px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white font-bold rounded-xl text-xs transition-colors shadow-xs cursor-pointer"
          >
            Tetap Simpan
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useApi } from '~/composables/useApi'

const route = useRoute()
const router = useRouter()

// Inquiry Context
const inqContext = reactive({
  id_rencana: 0,
  id_rencana_det: 0,
  nama: '',
  hp: '',
  ket_rencana: '',
  ket_aktivitas: '',
  nama_hslaktiv: ''
})

// Form State (Matches v_inputdo.php)
const form = reactive({
  nosbg: '',
  jnstrx: '0', // 0: None, 1: Gadai, 2: Cicilan
  produk: '0',
  idrendet: '',
  ket_rencana: '',
  ket_aktivitas: '',
  nm_customer: '',
  tgl_cair: ''
})

// SBG List & Searchable Dropdown State
const sbgList = ref<any[]>([])
const selectedSbg = ref<any>(null)
const isDropdownOpen = ref(false)
const searchTerm = ref('')
const isLoadingSbg = ref(false)
const searchInputRef = ref<HTMLInputElement | null>(null)

// Products List
const produkList = ref<any[]>([])

// Alerts & Modal
const isSubmitting = ref(false)
const showConfirmDateModal = ref(false)
const confirmMsg = ref('')
const calculatedDiffDays = ref(0)

const alertState = reactive({
  show: false,
  message: '',
  type: 'info' as 'success' | 'danger' | 'warning' | 'info'
})

const showAlert = (msg: string, type: 'success' | 'danger' | 'warning' = 'danger') => {
  alertState.message = msg
  alertState.type = type
  alertState.show = true
  if (type === 'success') {
    setTimeout(() => {
      alertState.show = false
    }, 6000)
  }
}

// Check Ready State (matches checkReady() in v_inputdo.php)
const isReady = computed(() => {
  return (
    form.nosbg !== '' &&
    form.jnstrx !== '0' &&
    form.produk !== '0'
  )
})

onMounted(async () => {
  const queryId = (route.query.id || route.query.idrendet || route.query.id_rencana || '') as string
  if (queryId) {
    form.idrendet = queryId
    await fetchInquiryContext(queryId)
  }

  // Load initial SBGs from Sopiga
  await fetchSbgList('')
})

// Fetch Inquiry Customer context to prefill idrendet and match customer
const fetchInquiryContext = async (id: string) => {
  try {
    const { data } = await useApi('/sopiga/inquiry-context', {
      params: { id }
    })
    if (data) {
      inqContext.id_rencana = data.id_rencana || 0
      inqContext.id_rencana_det = data.id_rencana_det || 0
      inqContext.nama = data.nama || ''
      inqContext.hp = data.hp || ''
      inqContext.ket_rencana = data.ket_rencana || ''
      inqContext.ket_aktivitas = data.ket_aktivitas || ''
      inqContext.nama_hslaktiv = data.nama_hslaktiv || ''

      form.idrendet = String(data.id_rencana_det || data.id_rencana || id)
      form.ket_rencana = data.ket_rencana || ''
      form.ket_aktivitas = data.ket_aktivitas || ''
      form.nm_customer = data.nama || ''

      // Pre-filter SBGs using the customer's name
      if (data.nama) {
        searchTerm.value = data.nama
        await fetchSbgList(data.nama)
      }
    }
  } catch (e) {
    console.warn('Inquiry context not found, proceeding in standalone mode:', e)
  }
}

// Fetch SBG list from Sopiga (parity with C_sistem/getsbglist)
const fetchSbgList = async (term: string) => {
  isLoadingSbg.value = true
  try {
    const { data } = await useApi('/sopiga/getsbglist', {
      params: {
        term: term || '',
        jnstrx: form.jnstrx !== '0' ? form.jnstrx : ''
      }
    })
    const list = Array.isArray(data) ? data : []
    sbgList.value = list
  } catch (e) {
    console.error('Failed to load SBG list:', e)
    sbgList.value = []
  } finally {
    isLoadingSbg.value = false
  }
}

let searchDebounce: any = null
const onSearchInput = () => {
  clearTimeout(searchDebounce)
  searchDebounce = setTimeout(() => {
    fetchSbgList(searchTerm.value)
  }, 300)
}

const toggleDropdown = () => {
  isDropdownOpen.value = !isDropdownOpen.value
  if (isDropdownOpen.value) {
    nextTick(() => {
      searchInputRef.value?.focus()
    })
  }
}

// On selecting an SBG (matches select2:select event in v_inputdo.php)
const selectSbgItem = (item: any) => {
  selectedSbg.value = item
  form.nosbg = item.id
  form.nm_customer = item.nm_customer || ''
  form.tgl_cair = item.tgl_cair || ''
  isDropdownOpen.value = false

  const jnstrxVal = item.jnstrx_val
  const kdProduk = item.kd_produk

  if (jnstrxVal && jnstrxVal !== '0') {
    form.jnstrx = String(jnstrxVal)
    loadProduk(jnstrxVal, kdProduk)
  }
}

const clearSbg = () => {
  selectedSbg.value = null
  form.nosbg = ''
  form.tgl_cair = ''
  form.jnstrx = '0'
  form.produk = '0'
  produkList.value = []
}

// On Jenis Transaksi change (matches tipe3 change event in v_inputdo.php)
const onJnsTrxChange = () => {
  if (form.jnstrx !== '0') {
    loadProduk(form.jnstrx, null)
  } else {
    produkList.value = []
    form.produk = '0'
  }
}

// Load products based on transaction type (matches loadProduk() in v_inputdo.php)
const loadProduk = async (tipeVal: string, selectVal: string | null) => {
  try {
    const { data } = await useApi('/sopiga/getproduk', {
      params: { id: tipeVal }
    })
    const items = Array.isArray(data) ? data : []
    produkList.value = items

    if (selectVal) {
      const match = items.find((i: any) => String(i.kd_produk) === String(selectVal))
      if (match) {
        form.produk = match.kd_produk
      } else {
        form.produk = items[0]?.kd_produk || '0'
      }
    } else {
      form.produk = items[0]?.kd_produk || '0'
    }
  } catch (e) {
    console.error('Failed to load products:', e)
    produkList.value = []
    form.produk = '0'
  }
}

// Save Click Handler with Date Confirmation Logic (matches v_inputdo.php lines 122-175)
const onSaveClick = () => {
  const nosbg = form.nosbg
  const tglCair = form.tgl_cair
  let diffDays = 0

  // Validasi Tanggal Cair vs Tanggal Submit (Hari ini)
  if (tglCair) {
    const dCair = new Date(tglCair)
    const options: Intl.DateTimeFormatOptions = {
      day: '2-digit',
      month: 'long',
      year: 'numeric'
    }

    const formattedDate = new Intl.DateTimeFormat('id-ID', options).format(dCair)
    const dNow = new Date()
    dCair.setHours(0, 0, 0, 0)
    dNow.setHours(0, 0, 0, 0)

    if (dCair.getTime() !== dNow.getTime()) {
      const diffTime = Math.abs(dCair.getTime() - dNow.getTime())
      diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))

      // Tampilkan Modal Custom bukannya confirm browser
      confirmMsg.value = `Tanggal cair (${formattedDate}) berbeda dengan tanggal hari ini. Selisih ${diffDays} hari.`
      calculatedDiffDays.value = diffDays
      showConfirmDateModal.value = true
      return // Berhenti dulu, tunggu user klik di modal
    }
  }

  executeSave(nosbg, diffDays)
}

// Confirm Save click inside modal
const onConfirmSave = () => {
  showConfirmDateModal.value = false
  executeSave(form.nosbg, calculatedDiffDays.value)
}

// Execute Save (matches executeSave() in v_inputdo.php)
const executeSave = async (nosbg: string, diffDays: number) => {
  const nmCustomer = selectedSbg.value?.nm_customer || form.nm_customer || inqContext.nama
  const jnstrx = form.jnstrx
  const produk = form.produk
  const idrendet = form.idrendet || String(inqContext.id_rencana_det || inqContext.id_rencana)
  const ketRencana = form.ket_rencana || inqContext.ket_rencana
  const ketAktivitas = form.ket_aktivitas || inqContext.ket_aktivitas

  if (!nosbg || jnstrx === '0' || produk === '0') {
    showAlert('Lengkapi semua data terlebih dahulu.', 'warning')
    return
  }

  isSubmitting.value = true

  try {
    const payload = {
      nosbg,
      jnstrx,
      produk,
      idrendet,
      nm_customer: nmCustomer,
      selisih_hari: diffDays,
      ket_rencana: ketRencana,
      ket_aktivitas: ketAktivitas
    }

    const { data, error } = await useApi('/sopiga/inputdo', {
      method: 'POST',
      body: payload
    })

    if (error || (data && data.status === 'error')) {
      const msg = typeof error === 'string' ? error : (data?.message || error?.message || 'Gagal menyimpan data DO.')
      showAlert(msg, 'danger')
      return
    }

    showAlert('Berhasil menyimpan.', 'success')
    setTimeout(() => {
      router.push('/inquiry')
    }, 1200)
  } catch (e: any) {
    showAlert('Gagal menyimpan. Silakan coba lagi: ' + (e.message || e), 'danger')
  } finally {
    isSubmitting.value = false
  }
}
</script>
