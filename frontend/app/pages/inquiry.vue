<template>
  <div class="space-y-6 pb-12 font-sans">
    <!-- Header Section -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-semibold text-slate-800 tracking-tight">Inquiry & Rencana Aktivitas</h1>
        <p class="text-slate-500 text-sm">Kelola pipeline sales, update status prospek, dan catat janji temu customer.</p>
      </div>
      <div class="flex items-center gap-3">
        <button @click="openModal = true" class="inline-flex items-center gap-2 bg-amber-500 hover:bg-amber-600 active:scale-95 text-white font-medium px-5 py-2.5 rounded-xl text-sm transition-all shadow-md shadow-amber-200 cursor-pointer">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
          Add Leads Baru
        </button>
      </div>
    </div>

    <!-- Summary Metrics Bar -->
    <div class="grid grid-cols-2 sm:grid-cols-3 gap-4">
      <div @click="selectedStatus = 'Leads'; fetchInquiries()" class="bg-white p-4 rounded-xl border border-slate-200/80 shadow-xs hover:border-amber-400 transition-all cursor-pointer">
        <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">Leads</span>
        <div class="text-2xl font-bold text-slate-800 mt-1">{{ summaryCounts.leads }}</div>
      </div>
      <div @click="selectedStatus = 'Prospect'; fetchInquiries()" class="bg-white p-4 rounded-xl border border-slate-200/80 shadow-xs hover:border-blue-400 transition-all cursor-pointer">
        <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">Prospect</span>
        <div class="text-2xl font-bold text-blue-600 mt-1">{{ summaryCounts.prospect }}</div>
      </div>
      <div @click="selectedStatus = 'Hot Prospect'; fetchInquiries()" class="bg-white p-4 rounded-xl border border-slate-200/80 shadow-xs hover:border-rose-400 transition-all cursor-pointer">
        <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">Hot Prospect</span>
        <div class="text-2xl font-bold text-rose-600 mt-1">{{ summaryCounts.hotProspect }}</div>
      </div>
      <!-- <div @click="selectedStatus = 'SBG'; fetchInquiries()" class="bg-white p-4 rounded-xl border border-slate-200/80 shadow-xs hover:border-emerald-400 transition-all cursor-pointer">
        <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">SBG (Deal)</span>
        <div class="text-2xl font-bold text-emerald-600 mt-1">{{ summaryCounts.sbg }}</div>
      </div> -->
    </div>

    <!-- Search & Filter Card -->
    <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-xs flex flex-col md:flex-row gap-4 justify-between items-end">
      <div class="flex-1 w-full md:w-auto">
        <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Cari Customer / Kontak</label>
        <div class="relative">
          <input
            v-model="searchQuery"
            @input="debounceFetch"
            type="text"
            placeholder="Ketik nama atau nomor handphone..."
            class="w-full pl-10 pr-4 py-2.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-800 focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all"
          />
          <svg class="w-4 h-4 text-slate-400 absolute left-3.5 top-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
        </div>
      </div>

      <div class="w-full md:w-56">
        <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Filter Status</label>
        <select v-model="selectedStatus" @change="fetchInquiries" class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-800 focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all">
          <option value="all">Semua Status Aktif</option>
          <option value="Leads">Leads</option>
          <option value="Prospect">Prospect</option>
          <option value="Hot Prospect">Hot Prospect</option>
        </select>
      </div>

      <div class="w-full md:w-48">
        <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Sumber Leads</label>
        <select v-model="selectedSource" @change="fetchInquiries" class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all">
          <option value="all">Semua Sumber</option>
          <option v-for="s in sourcesList" :key="s" :value="s">{{ s }}</option>
        </select>
      </div>
    </div>

    <!-- Data Table Card -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-sm">
          <thead>
            <tr class="bg-slate-50/80 border-b border-slate-200 text-xs font-semibold text-slate-500 uppercase tracking-wider">
              <th class="px-5 py-3.5">Tanggal</th>
              <th class="px-5 py-3.5">Nama Customer</th>
              <th class="px-5 py-3.5 text-center">Sumber</th>
              <th class="px-5 py-3.5 text-center">Status</th>
              <th class="px-5 py-3.5">Aktivitas Terakhir</th>
              <th class="px-5 py-3.5">Keterangan Barang</th>
              <th class="px-5 py-3.5 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-slate-700">
            <tr v-if="isLoadingTable">
              <td colspan="6" class="text-center py-8">
                <div class="w-6 h-6 border-2 border-amber-500 border-t-transparent rounded-full animate-spin mx-auto"></div>
              </td>
            </tr>
            <tr v-else-if="inquiries.length === 0">
              <td colspan="6" class="text-center py-8 text-slate-400 font-medium">
                Tidak ada data inquiry ditemukan.
              </td>
            </tr>
            <!-- <pre>{{ inquiries }}</pre> -->
            <tr v-for="item in inquiries" :key="item.id_rencana" class="hover:bg-amber-50/30 transition-all">
              <td class="px-5 py-3.5 whitespace-nowrap text-slate-500">{{ formatDate(item.tgl_rencana) }}</td>
              <td class="px-5 py-3.5 font-semibold text-slate-800">
                {{ item.nama }}
                <!-- <span v-if="item.hp" class="block text-xs font-normal text-slate-400">{{ item.hp }}</span> -->
              </td>
              <td class="px-5 py-3.5">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-slate-100 text-slate-600">
                  {{ item.source }}
                </span>
              </td>
              <td class="px-5 py-3.5 flex justify-center">
                <span class="inline-flex items-center px-5 py-0.5 rounded-full text-xs font-semibold uppercase tracking-wide" :class="getStatusBadgeClass(item.status)">
                  {{ item.status }}
                </span>
              </td>
              <td class="px-5 py-3.5">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold uppercase tracking-wide">
                  {{ item.aktivitas }}
                </span>
              </td>
              <td class="px-5 py-3.5">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold uppercase tracking-wide">
                  {{ item.keterangan_barang }}
                </span>
              </td>
              <td class="px-5 py-3.5 text-center">
                <div class="inline-flex items-center gap-1.5" >
                  <NuxtLink to="/input-do" class="p-1.5 text-green-600 hover-green-500 rounded-lg border border-green-200 transition-all" title="Submit SBG" v-if="item.id_hslaktiv == 3">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
  <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
</svg>
                  </NuxtLink>
                  <NuxtLink :to="`/input-rencana?id=${item.id_rencana || item.id_customer}`" class="p-1.5 text-amber-600 hover:bg-amber-50 rounded-lg border border-amber-200 transition-all" title="Buka Halaman Progress & Form SBG" v-if="item.id_hslaktiv != 3">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg>
                  </NuxtLink>
                  <button @click="showDetail(item)" class="p-1.5 text-slate-600 hover:bg-slate-100 rounded-lg border border-slate-200 transition-all" title="Detail">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/></svg>
                  </button>
                  <button @click="handleDelete(item.id_rencana)" class="p-1.5 text-rose-600 hover:bg-rose-50 rounded-lg border border-rose-200 transition-all" title="Hapus">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- MODAL: BUAT RENCANA LEADS BARU -->
    <div v-if="openModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs z-50 flex items-center justify-center p-4">
      <div class="bg-white rounded-2xl max-w-lg w-full p-6 shadow-2xl space-y-4 max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between border-b border-slate-100 pb-3">
          <h3 class="font-semibold text-slate-800 text-base">Buat Rencana Leads Baru</h3>
          <button @click="closeFormModal" class="text-slate-400 hover:text-slate-600 p-1">✕</button>
        </div>

        <form @submit.prevent="handleSubmitRencana" class="space-y-4">
          <div>
            <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Tipe Customer</label>
            <div class="grid grid-cols-2 gap-3">
              <label class="flex items-center justify-center gap-2 p-2.5 border rounded-xl cursor-pointer text-sm font-medium" :class="form.cust === 'Baru' ? 'bg-amber-50 border-amber-400 text-amber-700' : 'bg-slate-50 border-slate-200 text-slate-600'">
                <input type="radio" v-model="form.cust" value="Baru" class="hidden" /> Customer Baru
              </label>
              <label class="flex items-center justify-center gap-2 p-2.5 border rounded-xl cursor-pointer text-sm font-medium" :class="form.cust === 'Lama' ? 'bg-amber-50 border-amber-400 text-amber-700' : 'bg-slate-50 border-slate-200 text-slate-600'">
                <input type="radio" v-model="form.cust" value="Lama" class="hidden" /> Pelanggan Lama
              </label>
            </div>
          </div>

          <template v-if="form.cust === 'Baru'">
            <div>
              <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">Nama Lengkap</label>
              <input v-model="form.nm" type="text" placeholder="Masukkan nama lengkap customer" class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm" required />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">No. Handphone / WA</label>
              <input v-model="form.hp" type="text" placeholder="08xxxxxxxxxx" class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm" required />
            </div>
          </template>

          <template v-else>
            <div>
              <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">Pilih Customer Existing</label>
              <select v-model="form.custlama" class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm" required>
                <option value="">- Pilih Customer -</option>
                <option v-for="c in existingCustomers" :key="c.id_customer" :value="c.id_customer">
                  {{ c.nama }} ({{ c.hp }})
                </option>
              </select>
            </div>
          </template>

          <div>
            <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">Sumber Leads</label>
            <select v-model="form.sumbercust" class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm">
              <option v-for="s in sourcesList" :key="s" :value="s">{{ s }}</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">Catatan Rencana Aktivitas</label>
            <textarea v-model="form.ket" rows="3" placeholder="Contoh: Berminat gadai laptop Asus, janji ketemuan tgl 5..." class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm"></textarea>
          </div>

          <div v-if="modalError" class="p-3 bg-rose-50 border border-rose-200 rounded-xl text-xs text-rose-600">
            {{ modalError }}
          </div>

          <div class="flex justify-end gap-3 pt-2">
            <button type="button" @click="closeFormModal" class="px-4 py-2 border border-slate-300 text-slate-600 rounded-xl text-sm font-medium">Batal</button>
            <button type="submit" class="px-5 py-2 bg-amber-500 hover:bg-amber-600 text-white rounded-xl text-sm font-medium shadow-sm" :disabled="isSubmitting">
              <span v-if="isSubmitting">Menyimpan...</span>
              <span v-else>Simpan Rencana</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL: DETAIL INQUIRY -->
    <div v-if="selectedItem" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs z-50 flex items-center justify-center p-4">
      <div class="bg-white rounded-2xl max-w-md w-full p-6 shadow-2xl space-y-4">
        <div class="flex items-center justify-between border-b border-slate-100 pb-3">
          <h3 class="font-semibold text-slate-800 text-base">Detail Aktivitas</h3>
          <button @click="selectedItem = null" class="text-slate-400 hover:text-slate-600 p-1">✕</button>
        </div>
        <div class="space-y-3 text-sm">
          <div class="flex justify-between py-1 border-b border-slate-100"><span class="text-slate-500 font-medium">Nama:</span> <span class="font-semibold text-slate-800">{{ selectedItem.nama }}</span></div>
          <div class="flex justify-between py-1 border-b border-slate-100"><span class="text-slate-500 font-medium">Status:</span> <span class="font-semibold" :class="getStatusBadgeClass(selectedItem.status)">{{ selectedItem.status }}</span></div>
          <div class="flex justify-between py-1 border-b border-slate-100"><span class="text-slate-500 font-medium">Sumber:</span> <span class="text-slate-700">{{ selectedItem.source || 'Walk In' }}</span></div>
          <div class="flex justify-between py-1 border-b border-slate-100"><span class="text-slate-500 font-medium">Tanggal Rencana:</span> <span class="text-slate-700">{{ formatDate(selectedItem.tgl_rencana) }}</span></div>
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

const inquiries = ref<any[]>([])
const searchQuery = ref('')
const selectedStatus = ref('all')
const selectedSource = ref('all')
const isLoadingTable = ref(false)
const openModal = ref(false)
const isSubmitting = ref(false)
const modalError = ref('')
const existingCustomers = ref<any[]>([])
const selectedItem = ref<any>(null)

const summaryCounts = reactive({
  leads: 0,
  prospect: 0,
  hotProspect: 0,
  sbg: 0
})

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

const sourcesList = ref<string[]>(['Walk In', 'Media Sosial', 'Event', 'Referensi'])

onMounted(() => {
  fetchSummaryCounts()
  fetchInquiries()
  fetchExistingCustomers()
  fetchSources()
})

const fetchSummaryCounts = async () => {
  try {
    const { data } = await useApi('/inquiry/counts')
    if (data) {
      summaryCounts.leads = data.leads || 0
      summaryCounts.prospect = data.prospect || 0
      summaryCounts.hotProspect = data.hotprospect || 0
      summaryCounts.sbg = data.sbg || 0
    }
  } catch (e) {
    console.error('Failed to fetch inquiry counts:', e)
  }
}

const fetchSources = async () => {
  const { data } = await useApi('/sumbercust')
  if (data && Array.isArray(data) && data.length > 0) {
    sourcesList.value = data.map((s: any) => s.nm_ref_cust || s.kd_ref_cust)
    if (sourcesList.value.length > 0) {
      form.sumbercust = sourcesList.value[0]
    }
  }
}

const fetchInquiries = async () => {
  isLoadingTable.value = true
  const statusParam = selectedStatus.value === 'all' ? '' : selectedStatus.value


  const  appSource  = await useApi('/sumbercust')
  
  const { data } = await useApi('/inquiry/filter', {
    params: {
      status: statusParam,
      status_filter: statusParam,
      search: searchQuery.value || '',
      start: 0,
      length: 300
    }
  })
  console.log('data', data)
  isLoadingTable.value = false

  const rawList = Array.isArray(data) ? data : (data?.data || [])

  if (rawList && rawList.length >= 0) {
    let list = rawList.map((i: any) => {
      const source = appSource.data?.find((value : any) => { return value.kd_ref_cust === i.id_sumbercust})?.nm_ref_cust || i.id_sumbercust || 'Walk In'
      return {
        ...i,
        id_customer: i.id_customer || i.id_rencana,
        nama: i.nama ,
        hp: i.hp,
        status: i.nama_hslaktiv,
        source: source,
        aktivitas: i.ket_aktivitas,
        keterangan_barang: i.ket_rencana
      }
    })

    // Exclude SBG / Deal status from Inquiry page list
    list = list.filter((i: any) => {
      const st = (i.status || i.nama_hslaktiv || '').toLowerCase()
      const idAktiv = Number(i.id_hslaktiv)
      return idAktiv !== 7 && !st.includes('sbg') && !st.includes('deal') && !st.includes('do')
    })

    if (selectedSource.value !== 'all') {
      list = list.filter((i: any) => (i.source || 'Walk In').toLowerCase().includes(selectedSource.value.toLowerCase()))
    }
    inquiries.value = list
    console.log('inquiries', inquiries.value);

    // Recalculate summary metrics if count endpoint not fetched
    console.log('inquiries', inquiries);
    console.log('rawlist', rawList);
    if (summaryCounts.leads === 0 && summaryCounts.prospect === 0) {
      summaryCounts.leads = rawList.filter((i: any) => {
        const st = (i.nama_hslaktiv || i.status || '').toLowerCase()
        return st.includes('lead') || st === ''
      }).length

      summaryCounts.prospect = rawList.filter((i: any) => {
        const st = (i.nama_hslaktiv || i.status || '').toLowerCase()
        return st === 'prospect'
      }).length

      summaryCounts.hotProspect = rawList.filter((i: any) => {
        const st = (i.nama_hslaktiv || i.status || '').toLowerCase()
        return st.includes('hot')
      }).length

      summaryCounts.sbg = rawList.filter((i: any) => {
        const st = (i.nama_hslaktiv || i.status || '').toLowerCase()
        return st.includes('sbg') || st.includes('do') || st.includes('deal')
      }).length
    }
  } else {
    inquiries.value = []
  }
}

const fetchExistingCustomers = async () => {
  const { data } = await useApi('/customer/select')
  if (data && Array.isArray(data) && data.length > 0) {
    existingCustomers.value = data.map((c: any) => ({
      id_customer: c.id,
      nama: c.text,
      hp: ''
    }))
  } else {
    const { data: inqData } = await useApi('/inquiry/filter', {
      params: { start: 0, length: 100 }
    })
    if (inqData && inqData.data) {
      const map = new Map()
      for (const item of inqData.data) {
        map.set(item.nama, { id_customer: item.id_customer || item.id_rencana, nama: item.nama, hp: item.id_sumbercust || 'Walk In' })
      }
      existingCustomers.value = Array.from(map.values())
    }
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

  if (form.cust === 'Lama' && form.custlama) {
    const custId = Number(form.custlama)
    const existingCount = inquiries.value.filter(
      (i: any) => Number(i.id_customer) === custId
    ).length
    if (existingCount >= 2) {
      modalError.value = 'Batas maksimal rencana / leads per customer hanya boleh 2. Customer ini sudah memiliki 2 rencana.'
      isSubmitting.value = false
      return
    }
  }

  const payload: any = {
    cust: form.cust,
    sumbercust: form.sumbercust,
    hslaktiv: 1, // Default status Leads
    katcust: '1',
    ket: form.ket,
    aktiv: form.aktiv
  }

  console.log('payload', payload);

  if (form.cust === 'Lama') {
    payload.custlama = Number(form.custlama)
  } else {
    payload.nm = form.nm
    payload.hp = form.hp
  }

  const { data, error } = await useApi('/rencana/input', {
    method: 'POST',
    body: payload
  })

  isSubmitting.value = false

  if (error) {
    modalError.value = typeof error === 'string' ? error : (error.message || 'Gagal menyimpan rencana')
    return
  }

  fetchSummaryCounts()
  fetchInquiries()
  closeFormModal()
}

const handleDelete = async (id: number) => {
  if (confirm('Apakah Anda yakin ingin menghapus rencana aktivitas ini?')) {
    const { error } = await useApi(`/rencana/${id}`, { method: 'DELETE' })
    if (!error) {
      fetchSummaryCounts()
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
    case 'Leads': return 'bg-amber-100 text-amber-800 border border-amber-200'
    case 'Prospect': return 'bg-blue-100 text-blue-800 border border-blue-200'
    case 'Hot Prospect': return 'bg-rose-100 text-rose-800 border border-rose-200 text-center'
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
</script>
