<template>
  <div class="space-y-6 pb-12 font-sans">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-slate-800 dark:text-white tracking-tight">Data Customer & CIF</h1>
        <p class="text-slate-500 dark:text-slate-400 text-sm">Direktori nasabah, data CIF terintegrasi, dan riwayat inquiry dari database.</p>
      </div>
      <div class="flex items-center gap-3">
        <button
          @click="fetchCustomers"
          class="px-3.5 py-2 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200 rounded-xl text-xs font-semibold shadow-xs transition-all flex items-center gap-1.5 cursor-pointer"
        >
          <AppIcon name="refresh" class="w-4 h-4 text-amber-500" :class="{ 'animate-spin': isLoading }" />
          Refresh Data
        </button>
      </div>
    </div>

    <!-- Metric Summary Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <div class="bg-white dark:bg-slate-800 p-5 rounded-2xl border border-slate-200/80 dark:border-slate-700/80 shadow-xs flex items-center justify-between">
        <div>
          <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">Total Customer</span>
          <div class="text-2xl font-extrabold text-slate-800 dark:text-white mt-1">{{ totalCustomersCount }}</div>
        </div>
        <div class="p-3 bg-amber-50 dark:bg-amber-950/40 text-amber-600 dark:text-amber-400 rounded-xl">
          <AppIcon name="users" class="w-6 h-6" />
        </div>
      </div>

      <div class="bg-white dark:bg-slate-800 p-5 rounded-2xl border border-slate-200/80 dark:border-slate-700/80 shadow-xs flex items-center justify-between">
        <div>
          <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">Customer Active SBG</span>
          <div class="text-2xl font-extrabold text-emerald-600 dark:text-emerald-400 mt-1">{{ activeSbgCount }}</div>
        </div>
        <div class="p-3 bg-emerald-50 dark:bg-emerald-950/40 text-emerald-600 dark:text-emerald-400 rounded-xl">
          <AppIcon name="check-circle" class="w-6 h-6" />
        </div>
      </div>

      <div class="bg-white dark:bg-slate-800 p-5 rounded-2xl border border-slate-200/80 dark:border-slate-700/80 shadow-xs flex items-center justify-between">
        <div>
          <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">Sumber Leads Terbanyak</span>
          <div class="text-xl font-bold text-amber-600 dark:text-amber-400 mt-1 truncate max-w-[180px]">{{ topSource }}</div>
        </div>
        <div class="p-3 bg-amber-50 dark:bg-amber-950/40 text-amber-600 dark:text-amber-400 rounded-xl">
          <AppIcon name="trending-up" class="w-6 h-6" />
        </div>
      </div>
    </div>

    <!-- Search Box Card -->
    <div class="bg-white dark:bg-slate-800 p-5 rounded-2xl border border-slate-200/80 dark:border-slate-700/80 shadow-xs">
      <label class="block text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider mb-1.5">Cari Data Customer & CIF</label>
      <div class="relative max-w-md">
        <input
          v-model="searchQuery"
          @input="debounceFetch"
          type="text"
          placeholder="Ketik CIF, nama, nomor handphone, atau kota..."
          class="w-full pl-10 pr-4 py-2.5 bg-white dark:bg-slate-900 border border-slate-300 dark:border-slate-700 rounded-xl text-sm text-slate-800 dark:text-slate-100 focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 dark:focus:ring-amber-900/50 transition-all"
        />
        <AppIcon name="search" class="w-4 h-4 text-slate-400 absolute left-3.5 top-3.5" />
      </div>
    </div>

    <!-- Table Card -->
    <div class="bg-white dark:bg-slate-800 rounded-2xl border border-slate-200/80 dark:border-slate-700/80 shadow-xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-sm">
          <thead>
            <tr class="bg-slate-50/80 dark:bg-slate-700/50 border-b border-slate-200 dark:border-slate-700 text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider">
              <th class="px-5 py-3.5">No</th>
              <th class="px-5 py-3.5">No. CIF / Identitas</th>
              <th class="px-5 py-3.5">Nama Customer</th>
              <th class="px-5 py-3.5">No. Handphone</th>
              <th class="px-5 py-3.5">Alamat / Kota</th>
              <th class="px-5 py-3.5">Sumber Leads</th>
              <th class="px-5 py-3.5">Status Terakhir</th>
              <th class="px-5 py-3.5 text-center">Total Inquiry</th>
              <th class="px-5 py-3.5 text-center">Aksi / Hubungi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 dark:divide-slate-700/60 text-slate-700 dark:text-slate-300">
            <tr v-if="isLoading">
              <td colspan="9" class="text-center py-12">
                <AppIcon name="loader" class="w-7 h-7  text-amber-500   animate-spin mx-auto mb-2" />
                <span class="text-xs font-medium text-slate-400">Memuat data Customer & CIF dari database...</span>
              </td>
            </tr>
            <tr v-else-if="customers.length === 0">
              <td colspan="9" class="text-center py-12 text-slate-400 font-medium">
                Tidak ada data customer ditemukan.
              </td>
            </tr>
            <tr v-for="(c, idx) in customers" :key="c.id_customer || idx" class="hover:bg-amber-50/30 dark:hover:bg-amber-950/20 transition-all">
              <td class="px-5 py-3.5 text-slate-400 font-medium text-xs">{{ idx + 1 }}</td>
              <td class="px-5 py-3.5 font-mono text-xs font-semibold text-amber-700 dark:text-amber-400">
                <span class="px-2 py-1 bg-amber-50 dark:bg-amber-950/40 border border-amber-200 dark:border-amber-800 rounded-lg inline-block">
                  {{ c.cif || ('CIF-' + c.id_customer) }}
                </span>
              </td>
              <td class="px-5 py-3.5 font-semibold text-slate-800 dark:text-slate-100">
                <button @click="openHistoryModal(c)" class="hover:text-amber-600 dark:hover:text-amber-400 text-left underline underline-offset-2 decoration-amber-300 cursor-pointer transition-colors">
                  {{ c.nama }}
                </button>
              </td>
              <td class="px-5 py-3.5 text-slate-600 dark:text-slate-300 font-mono text-xs">{{ c.hp || '-' }}</td>
              <td class="px-5 py-3.5 text-xs text-slate-500 dark:text-slate-400 max-w-[180px] truncate">
                {{ formatAddress(c.alamat, c.kota) }}
              </td>
              <td class="px-5 py-3.5">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-slate-100 dark:bg-slate-700 text-slate-600 dark:text-slate-300">
                  {{ c.source_name || 'Tidak diketahui' }}
                </span>
              </td>
              <td class="px-5 py-3.5">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold uppercase tracking-wide" :class="getStatusBadgeClass(c.latest_status)">
                  {{ c.latest_status }}
                </span>
              </td>
              <td class="px-5 py-3.5 text-center font-bold text-slate-800 dark:text-white">
                <span class="px-2 py-0.5 rounded-full bg-slate-100 dark:bg-slate-700 text-xs font-semibold text-slate-700 dark:text-slate-300">
                  {{ c.total_inquiries || 0 }} Rencana
                </span>
              </td>
              <td class="px-5 py-3.5 text-center">
                <div class="inline-flex items-center gap-1.5">
                  <button
                    @click="openHistoryModal(c)"
                    title="Lihat Detail & Riwayat Inquiry"
                    class="px-2.5 py-1.5 bg-amber-50 hover:bg-amber-100 dark:bg-amber-950/40 dark:hover:bg-amber-900/60 text-amber-700 dark:text-amber-300 border border-amber-200 dark:border-amber-800 rounded-lg text-xs font-semibold transition-all flex items-center gap-1 cursor-pointer"
                  >
                    <AppIcon name="eye" class="w-3.5 h-3.5" />
                    Detail Inquiry
                  </button>

                  <button
                    @click="openWaModal(c)"
                    title="Kirim WhatsApp"
                    class="p-1.5 bg-emerald-50 hover:bg-emerald-100 dark:bg-emerald-950/40 dark:hover:bg-emerald-900/60 text-emerald-700 dark:text-emerald-300 border border-emerald-200 dark:border-emerald-800 rounded-lg transition-all cursor-pointer"
                  >
                    <AppIcon name="message" class="w-3.5 h-3.5" />
                  </button>

                  <a
                    :href="`tel:${cleanPhone(c.hp)}`"
                    title="Telepon Direct"
                    class="p-1.5 bg-slate-50 hover:bg-slate-100 dark:bg-slate-700 dark:hover:bg-slate-600 text-slate-700 dark:text-slate-200 border border-slate-200 dark:border-slate-600 rounded-lg transition-all"
                  >
                    <AppIcon name="phone" class="w-3.5 h-3.5" />
                  </a>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- WHATSAPP MESSAGE MODAL WITH PRESETS -->
    <ClientOnly>
      <Teleport to="body">
        <div v-if="activeWaCustomer" class="fixed inset-0 w-screen h-screen min-h-screen top-0 left-0 right-0 bottom-0 z-[9999] bg-slate-900/60 backdrop-blur-md flex items-center justify-center p-4 !mt-0 overflow-y-auto">
          <div class="bg-white dark:bg-slate-800 rounded-2xl max-w-md w-full p-6 shadow-2xl space-y-4 border border-slate-100 dark:border-slate-700 my-auto">
            <div class="flex items-center justify-between border-b border-slate-100 dark:border-slate-700 pb-3">
              <h3 class="font-bold text-slate-800 dark:text-white text-base">Kirim Pesan WhatsApp</h3>
              <button @click="activeWaCustomer = null" class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 p-1 cursor-pointer"><AppIcon name="close" class="w-4 h-4" /></button>
            </div>

            <div class="bg-emerald-50/60 dark:bg-emerald-950/40 p-3.5 rounded-xl border border-emerald-200 dark:border-emerald-800/60">
              <p class="text-sm font-bold text-slate-800 dark:text-slate-100">{{ activeWaCustomer.nama }}</p>
              <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">Tujuan WA: <strong class="font-mono text-emerald-700 dark:text-emerald-400">{{ cleanPhone(activeWaCustomer.hp) }}</strong></p>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider mb-1.5">Pilih Template Pesan</label>
              <div class="space-y-2">
                <button
                  v-for="(tpl, i) in waTemplates"
                  :key="i"
                  @click="waMessage = tpl.text.replace('{nama}', activeWaCustomer.nama)"
                  class="w-full text-left p-2.5 bg-slate-50 dark:bg-slate-900 hover:bg-emerald-50 dark:hover:bg-emerald-950/40 rounded-xl border border-slate-200 dark:border-slate-700 text-xs text-slate-700 dark:text-slate-300 transition-all cursor-pointer"
                >
                  <strong class="text-slate-800 dark:text-slate-100 block mb-0.5">{{ tpl.title }}</strong>
                  <span class="text-slate-500 dark:text-slate-400 line-clamp-1">{{ tpl.text.replace('{nama}', activeWaCustomer.nama) }}</span>
                </button>
              </div>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider mb-1">Isi Pesan WhatsApp</label>
              <textarea v-model="waMessage" rows="4" class="w-full px-3.5 py-2.5 bg-white dark:bg-slate-900 border border-slate-300 dark:border-slate-700 rounded-xl text-sm text-slate-800 dark:text-slate-100 focus:outline-none focus:border-emerald-500"></textarea>
            </div>

            <div class="flex justify-end gap-3 pt-2">
              <button @click="activeWaCustomer = null" class="px-4 py-2 border border-slate-300 dark:border-slate-700 text-slate-600 dark:text-slate-300 rounded-xl text-sm font-medium cursor-pointer">Batal</button>
              <a
                :href="`https://wa.me/${cleanPhone(activeWaCustomer.hp)}?text=${encodeURIComponent(waMessage)}`"
                target="_blank"
                @click="activeWaCustomer = null"
                class="px-5 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-sm font-semibold shadow-sm flex items-center gap-1.5 cursor-pointer"
              >
                Buka WhatsApp
              </a>
            </div>
          </div>
        </div>
      </Teleport>
    </ClientOnly>

    <!-- CUSTOMER PROFILE & INQUIRY HISTORY MODAL FROM DB -->
    <ClientOnly>
      <Teleport to="body">
        <div v-if="historyCustomer" class="fixed inset-0 w-screen h-screen min-h-screen top-0 left-0 right-0 bottom-0 z-[9999] bg-slate-900/60 backdrop-blur-md flex items-center justify-center p-4 !mt-0 overflow-y-auto">
          <div class="bg-white dark:bg-slate-800 rounded-2xl max-w-2xl w-full p-6 shadow-2xl space-y-5 border border-slate-100 dark:border-slate-700 max-h-[90vh] overflow-y-auto my-auto">
        
        <!-- Modal Header -->
        <div class="flex items-center justify-between border-b border-slate-100 dark:border-slate-700 pb-3">
          <div class="flex items-center gap-2">
            <div class="w-9 h-9 rounded-xl bg-amber-50 dark:bg-amber-950/40 text-amber-600 dark:text-amber-400 flex items-center justify-center font-bold text-sm">
              CIF
            </div>
            <div>
              <h3 class="font-bold text-slate-800 dark:text-white text-base">Profil & Riwayat Inquiry Customer</h3>
              <p class="text-xs text-slate-500 dark:text-slate-400 font-mono">{{ historyCustomer.cif || ('CIF-' + historyCustomer.id_customer) }}</p>
            </div>
          </div>
          <button @click="historyCustomer = null" class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 p-1 cursor-pointer"><AppIcon name="close" class="w-4 h-4" /></button>
        </div>

        <!-- Customer Profile Summary Box -->
        <div class="bg-slate-50 dark:bg-slate-900/60 p-4 rounded-xl border border-slate-200/80 dark:border-slate-700 grid grid-cols-1 md:grid-cols-2 gap-3 text-xs">
          <div>
            <span class="text-slate-400 font-medium">Nama Nasabah:</span>
            <p class="text-sm font-bold text-slate-800 dark:text-slate-100 mt-0.5">{{ historyCustomer.nama }}</p>
          </div>
          <div>
            <span class="text-slate-400 font-medium">Nomor Handphone:</span>
            <p class="text-sm font-mono font-semibold text-slate-700 dark:text-slate-200 mt-0.5">{{ historyCustomer.hp || '-' }}</p>
          </div>
          <div>
            <span class="text-slate-400 font-medium">Alamat Lengkap:</span>
            <p class="text-xs text-slate-700 dark:text-slate-300 mt-0.5">{{ formatAddress(historyCustomer.alamat, historyCustomer.kota) }}</p>
          </div>
          <div>
            <span class="text-slate-400 font-medium">Sumber Leads & Officer:</span>
            <p class="text-xs text-slate-700 dark:text-slate-300 mt-0.5 font-medium">
              {{ historyCustomer.source_name || 'Walk In' }} • {{ historyCustomer.id_login || 'System' }}
            </p>
          </div>
        </div>

        <!-- Inquiry History Section -->
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <h4 class="text-xs font-bold uppercase tracking-wider text-slate-500 dark:text-slate-400">
              Riwayat Inquiry & Rencana Aktivitas (Database MySQL)
            </h4>
            <span class="text-xs font-semibold text-amber-600 dark:text-amber-400">
              Total: {{ customerInquiries.length }} Record
            </span>
          </div>

          <!-- Loading state -->
          <div v-if="isLoadingHistory" class="py-8 text-center bg-slate-50 dark:bg-slate-900/40 rounded-xl border border-slate-200 dark:border-slate-700">
            <AppIcon name="loader" class="w-6 h-6  text-amber-500   animate-spin mx-auto mb-2" />
            <span class="text-xs text-slate-400 font-medium">Memuat riwayat inquiry dari database...</span>
          </div>

          <!-- Empty state -->
          <div v-else-if="customerInquiries.length === 0" class="py-8 text-center bg-slate-50 dark:bg-slate-900/40 rounded-xl border border-slate-200 dark:border-slate-700 text-slate-400 text-xs">
            Belum ada catatan inquiry / rencana aktivitas untuk nasabah ini di database.
          </div>

          <!-- Timeline History List -->
          <div v-else class="space-y-3 max-h-[320px] overflow-y-auto pr-1">
            <div
              v-for="(inq, idx) in customerInquiries"
              :key="inq.id_rencana_det || idx"
              class="p-4 bg-white dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-700 shadow-2xs space-y-2 hover:border-amber-300 transition-all"
            >
              <div class="flex items-center justify-between gap-2 border-b border-slate-100 dark:border-slate-800 pb-2">
                <div class="flex items-center gap-2">
                  <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-bold uppercase tracking-wide" :class="getStatusBadgeClass(inq.nama_hslaktiv)">
                    {{ inq.nama_hslaktiv }}
                  </span>
                  <span v-if="inq.no_sbg" class="text-xs font-mono font-bold text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-950/40 px-2 py-0.5 rounded-md border border-emerald-200 dark:border-emerald-800">
                    SBG: {{ inq.no_sbg }}
                  </span>
                </div>
                <span class="text-xs font-mono text-slate-400 dark:text-slate-500">
                  {{ formatDate(inq.tgl_rencana) }}
                </span>
              </div>

              <div class="text-xs text-slate-700 dark:text-slate-300 space-y-1">
                <p v-if="inq.ket_rencana">
                  <strong class="text-slate-800 dark:text-slate-200">Keterangan Rencana:</strong> {{ inq.ket_rencana }}
                </p>
                <p v-if="inq.ket_aktivitas">
                  <strong class="text-slate-800 dark:text-slate-200">Hasil Aktivitas:</strong> {{ inq.ket_aktivitas }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Modal Footer Actions -->
        <div class="flex items-center justify-between pt-3 border-t border-slate-100 dark:border-slate-700">
          <NuxtLink
            :to="`/inquiry?search=${encodeURIComponent(historyCustomer.nama)}`"
            class="px-4 py-2 bg-amber-500 hover:bg-amber-600 text-white rounded-xl text-xs font-semibold shadow-xs flex items-center gap-1.5 transition-all"
          >
            <AppIcon name="external-link" class="w-3.5 h-3.5" />
            Buka Halaman Inquiry
          </NuxtLink>
          <button @click="historyCustomer = null" class="px-5 py-2 border border-slate-300 dark:border-slate-700 text-slate-600 dark:text-slate-300 rounded-xl text-xs font-semibold hover:bg-slate-100 dark:hover:bg-slate-700 cursor-pointer">
            Tutup
          </button>
        </div>

      </div>
        </div>
      </Teleport>
    </ClientOnly>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'

interface CustomerItem {
  id_customer: number
  cif: string
  nama: string
  hp: string
  alamat: string
  kota: string
  id_login: string
  id_sumbercust: string
  source_name?: string
  latest_status_id: number
  latest_status: string
  latest_tgl: string
  total_inquiries: number
}

interface InquiryHistoryItem {
  id_rencana: number
  id_rencana_det: number
  tgl_rencana: string
  id_hslaktiv: number
  nama_hslaktiv: string
  ket_rencana: string
  ket_aktivitas: string
  id_sumbercust: string
  no_sbg: string
}

const customers = ref<CustomerItem[]>([])
const totalCustomersCount = ref(0)
const searchQuery = ref('')
const isLoading = ref(false)
const activeWaCustomer = ref<CustomerItem | null>(null)
const waMessage = ref('')
const historyCustomer = ref<CustomerItem | null>(null)
const customerInquiries = ref<InquiryHistoryItem[]>([])
const isLoadingHistory = ref(false)

const waTemplates = [
  { title: 'Follow Up Penawaran Gadai', text: 'Halo Bpk/Ibu {nama}, salam dari PT Gadai Mulia. Mengenai simulasi pencairan gadai sebelumnya, apakah ada yang bisa kami bantu kembali?' },
  { title: 'Informasi Cek Taksiran Emas', text: 'Halo Bpk/Ibu {nama}, kami menginfokan bahwa harga acuan taksiran emas hari ini sedang bagus. Apabila berminat mensimulasikan kembali, silakan hubungi kami.' },
  { title: 'Greetings Customer Baru', text: 'Halo Bpk/Ibu {nama}, terima kasih telah berkonsultasi dengan Gadai. Kami siap memberikan layanan pencairan cepat dan sewa modal terjangkau.' }
]

const activeSbgCount = computed(() => {
  return customers.value.filter(c => (c.latest_status || '').toLowerCase().includes('sbg') || c.latest_status_id === 7).length
})

const topSource = computed(() => {
  if (customers.value.length === 0) return 'Walk In'
  const counts: Record<string, number> = {}
  for (const c of customers.value) {
    const src = c.source_name || 'Walk In'
    counts[src] = (counts[src] || 0) + 1
  }
  let maxSrc = 'Walk In'
  let maxCnt = 0
  for (const [src, cnt] of Object.entries(counts)) {
    if (cnt > maxCnt) {
      maxCnt = cnt
      maxSrc = src
    }
  }
  return maxSrc
})

onMounted(() => {
  fetchCustomers()
})

const fetchCustomers = async () => {
  isLoading.value = true

  try {
    // 1. Fetch Master Reference Sources
    let sourceMap = new Map<string, string>()
    try {
      const appSource = await useApi('/sumbercust')
      const sourceList = Array.isArray(appSource?.data) ? appSource.data : (appSource?.data?.data || [])
      if (Array.isArray(sourceList)) {
        for (const item of sourceList) {
          if (item.kd_ref_cust) {
            sourceMap.set(String(item.kd_ref_cust), item.nm_ref_cust || '')
          }
        }
      }
    } catch (e) {
      console.warn('Failed to load sumbercust', e)
    }

    // 2. Fetch Customer List from DB
    const res = await useApi('/customer/list', {
      params: {
        search: searchQuery.value || '',
        start: 0,
        length: 300
      }
    })

    const rawData = res?.data
    const items: CustomerItem[] = Array.isArray(rawData) ? rawData : (rawData?.data || [])
    totalCustomersCount.value = rawData?.recordsTotal || items.length

    // 3. Map Source Name
    for (const c of items) {
      if (c.id_sumbercust && sourceMap.has(String(c.id_sumbercust))) {
        c.source_name = sourceMap.get(String(c.id_sumbercust))
      } else {
        c.source_name = 'Walk In'
      }
    }

    customers.value = items
  } catch (err) {
    console.error('Error fetching customers:', err)
  } finally {
    isLoading.value = false
  }
}

let debounceTimer: any = null
const debounceFetch = () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    fetchCustomers()
  }, 400)
}

const openWaModal = (c: CustomerItem) => {
  activeWaCustomer.value = c
  waMessage.value = waTemplates[0]?.text.replace('{nama}', c.nama) ?? ''
}

const openHistoryModal = async (c: CustomerItem) => {
  historyCustomer.value = c
  customerInquiries.value = []
  isLoadingHistory.value = true

  try {
    const res = await useApi(`/customer/${c.id_customer}/inquiries`)
    const historyData = Array.isArray(res?.data) ? res.data : (res?.data?.data || [])
    customerInquiries.value = historyData
  } catch (err) {
    console.error('Failed to load customer inquiry history:', err)
  } finally {
    isLoadingHistory.value = false
  }
}

const cleanPhone = (phoneStr: string) => {
  if (!phoneStr) return '6283138794243'
  let clean = phoneStr.replace(/\D/g, '')
  if (clean.startsWith('0')) {
    clean = '62' + clean.slice(1)
  }
  return clean || '6283138794243'
}

const formatAddress = (alamat?: string, kota?: string) => {
  const parts = [alamat, kota].filter(Boolean)
  return parts.length > 0 ? parts.join(', ') : 'Belum diisi'
}

const formatDate = (dateStr?: string) => {
  if (!dateStr) return '-'
  try {
    const d = new Date(dateStr)
    if (isNaN(d.getTime())) return dateStr
    return d.toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return dateStr
  }
}

const getStatusBadgeClass = (status?: string) => {
  const s = (status || '').toLowerCase()
  if (s.includes('sbg') || s.includes('deal')) {
    return 'bg-emerald-100 dark:bg-emerald-950/60 text-emerald-800 dark:text-emerald-300 border border-emerald-200 dark:border-emerald-800'
  }
  if (s.includes('hot')) {
    return 'bg-rose-100 dark:bg-rose-950/60 text-rose-800 dark:text-rose-300 border border-rose-200 dark:border-rose-800'
  }
  if (s.includes('prospect')) {
    return 'bg-blue-100 dark:bg-blue-950/60 text-blue-800 dark:text-blue-300 border border-blue-200 dark:border-blue-800'
  }
  if (s.includes('batal')) {
    return 'bg-slate-100 dark:bg-slate-700 text-slate-600 dark:text-slate-400 border border-slate-200 dark:border-slate-600'
  }
  return 'bg-amber-100 dark:bg-amber-950/60 text-amber-800 dark:text-amber-300 border border-amber-200 dark:border-amber-800'
}
</script>
