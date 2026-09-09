<template>
  <div class="space-y-6 pb-12 font-sans">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-semibold text-slate-800 tracking-tight">Data Customer & CIF</h1>
        <p class="text-slate-500 text-sm">Direktori pelanggan, kontak langsung WhatsApp, dan riwayat transaksi gadai.</p>
      </div>
    </div>

    <!-- Metric Summary Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-xs flex items-center justify-between">
        <div>
          <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">Total Customer</span>
          <div class="text-2xl font-bold text-slate-800 mt-1">{{ customers.length }}</div>
        </div>
        <div class="p-3 bg-blue-50 text-blue-600 rounded-xl">
          <svg data-v-433a9abd="" class="nav-icon" width="20" height="20"  fill="none"><path data-v-433a9abd="" d="M17 21V19C17 17.9391 16.5786 16.9217 15.8284 16.1716C15.0783 15.4214 14.0609 15 13 15H5C3.93913 15 2.92172 15.4214 2.17157 16.1716C1.42143 16.9217 1 17.9391 1 19V21" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path><path data-v-433a9abd="" d="M9 11C11.2091 11 13 9.20914 13 7C13 4.79086 11.2091 3 9 3C6.79086 3 5 4.79086 5 7C5 9.20914 6.79086 11 9 11Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path></svg>
        </div>
      </div>

      <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-xs flex items-center justify-between">
        <div>
          <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">Customer Active SBG</span>
          <div class="text-2xl font-bold text-emerald-600 mt-1">{{ activeSbgCount }}</div>
        </div>
        <div class="p-3 bg-emerald-50 text-emerald-600 rounded-xl">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
        </div>
      </div>

      <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-xs flex items-center justify-between">
        <div>
          <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">Sumber Terbanyak</span>
          <div class="text-xl font-bold text-amber-600 mt-1">Walk In</div>
        </div>
        <div class="p-3 bg-amber-50 text-amber-600 rounded-xl">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"/></svg>
        </div>
      </div>
    </div>

    <!-- Search Box Card -->
    <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-xs">
      <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Cari Customer</label>
      <div class="relative max-w-md">
        <input
          v-model="searchQuery"
          @input="debounceFetch"
          type="text"
          placeholder="Ketik nama atau nomor telepon..."
          class="w-full pl-10 pr-4 py-2.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-800 focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all"
        />
        <svg class="w-4 h-4 text-slate-400 absolute left-3.5 top-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
      </div>
    </div>

    <!-- Table Card -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-sm">
          <thead>
            <tr class="bg-slate-50/80 border-b border-slate-200 text-xs font-semibold text-slate-500 uppercase tracking-wider">
              <th class="px-5 py-3.5">No</th>
              <th class="px-5 py-3.5">Nama Customer</th>
              <th class="px-5 py-3.5">No. Handphone</th>
              <th class="px-5 py-3.5">Sumber Leads</th>
              <th class="px-5 py-3.5">Status Terakhir</th>
              <th class="px-5 py-3.5 text-center">Hubungi Langsung</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-slate-700">
            <tr v-if="isLoading">
              <td colspan="6" class="text-center py-8">
                <div class="w-6 h-6 border-2 border-amber-500 border-t-transparent rounded-full animate-spin mx-auto"></div>
              </td>
            </tr>
            <tr v-else-if="customers.length === 0">
              <td colspan="6" class="text-center py-8 text-slate-400 font-medium">
                Tidak ada data customer ditemukan.
              </td>
            </tr>
            <tr v-for="(c, idx) in customers" :key="idx" class="hover:bg-amber-50/30 transition-all">
              <td class="px-5 py-3.5 text-slate-500">{{ idx + 1 }}</td>
              <td class="px-5 py-3.5 font-semibold text-slate-800">
                <span @click="openHistoryModal(c)" class="hover:text-amber-600 cursor-pointer underline underline-offset-2 decoration-amber-300">
                  {{ c.nama }}
                </span>
              </td>
              <td class="px-5 py-3.5 text-slate-600 font-mono text-sm">{{ c.phone || '-' }}</td>
              <td class="px-5 py-3.5">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium text-slate-600">
                  {{ c.source || 'Tidak ditemukan' }}
                </span>
              </td>
              <td class="px-5 py-3.5">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold uppercase tracking-wide" :class="getStatusBadgeClass(c.status)">
                  {{ c.status }}
                </span>
              </td>
              <td class="px-5 py-3.5 text-center">
                <div class="inline-flex items-center gap-2">
                  <a :href="`tel:${cleanPhone(c.phone)}`" class="px-3 py-1.5 bg-amber-50 hover:bg-amber-100 text-amber-700 border border-amber-200 rounded-lg text-xs font-medium transition-all flex items-center gap-1">
                    <svg width="12" height="12" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"/></svg>
                    Telepon
                  </a>
                  <button @click="openWaModal(c)" class="px-3 py-1.5 bg-emerald-50 hover:bg-emerald-100 text-emerald-700 border border-emerald-200 rounded-lg text-xs font-medium transition-all flex items-center gap-1 cursor-pointer">
                    <svg width="12" height="12" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/></svg>
                    WhatsApp
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- WHATSAPP MESSAGE MODAL WITH PRESETS -->
    <div v-if="activeWaCustomer" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs z-50 flex items-center justify-center p-4">
      <div class="bg-white rounded-2xl max-w-md w-full p-6 shadow-2xl space-y-4">
        <div class="flex items-center justify-between border-b border-slate-100 pb-3">
          <h3 class="font-semibold text-slate-800 text-base">Kirim Pesan WhatsApp</h3>
          <button @click="activeWaCustomer = null" class="text-slate-400 hover:text-slate-600 p-1">✕</button>
        </div>

        <div class="bg-emerald-50/60 p-3 rounded-xl border border-emerald-200">
          <p class="text-sm font-semibold text-slate-800">{{ activeWaCustomer.nama }}</p>
          <p class="text-xs text-slate-500">Tujuan WA: <strong class="font-mono text-emerald-700">{{ cleanPhone(activeWaCustomer.phone) }}</strong></p>
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Pilih Template Pesan</label>
          <div class="space-y-2">
            <button
              v-for="(tpl, i) in waTemplates"
              :key="i"
              @click="waMessage = tpl.text.replace('{nama}', activeWaCustomer.nama)"
              class="w-full text-left p-2.5 bg-slate-50 hover:bg-emerald-50 rounded-xl border border-slate-200 text-xs text-slate-700 transition-all"
            >
              <strong class="text-slate-800 block mb-0.5">{{ tpl.title }}</strong>
              <span class="text-slate-500 line-clamp-1">{{ tpl.text.replace('{nama}', activeWaCustomer.nama) }}</span>
            </button>
          </div>
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">Isi Pesan WhatsApp</label>
          <textarea v-model="waMessage" rows="4" class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-800"></textarea>
        </div>

        <div class="flex justify-end gap-3 pt-2">
          <button @click="activeWaCustomer = null" class="px-4 py-2 border border-slate-300 text-slate-600 rounded-xl text-sm font-medium">Batal</button>
          <a
            :href="`https://wa.me/${cleanPhone(activeWaCustomer.phone)}?text=${encodeURIComponent(waMessage)}`"
            target="_blank"
            @click="activeWaCustomer = null"
            class="px-5 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-sm font-medium shadow-sm flex items-center gap-1.5"
          >
            Buka WhatsApp
          </a>
        </div>
      </div>
    </div>

  </div>

      <!-- CUSTOMER HISTORY MODAL -->
    <div v-if="historyCustomer" class="space-y-none fixed inset-0 bg-slate-900/40 backdrop-blur-xs space-y-0  z-50 flex items-center justify-center p-4">
      <div class="bg-white rounded-2xl max-w-lg w-full p-6 shadow-2xl space-y-4">
        <div class="flex items-center justify-between border-b border-slate-100 pb-3">
          <h3 class="font-semibold text-slate-800 text-base">Profil & Riwayat Customer</h3>
          <button @click="historyCustomer = null" class="text-slate-400 hover:text-slate-600 p-1">✕</button>
        </div>

        <div class="bg-slate-50 p-4 rounded-xl border border-slate-200 space-y-1">
          <p class="text-base font-semibold text-slate-800">{{ historyCustomer.nama }}</p>
          <p class="text-xs text-slate-500">Nomor Telepon: <span class="font-mono font-medium text-slate-700">{{ historyCustomer.phone }}</span></p>
          <p class="text-xs text-slate-500">Sumber Leads: <span class="font-medium text-slate-700">{{ historyCustomer.source || 'Tidak Ditemukan' }}</span></p>
          <p class="text-xs text-slate-500">Status Terakhir: <span class="font-semibold text-amber-600">{{ historyCustomer.status }}</span></p>
        </div>

        <div class="space-y-2">
          <h4 class="text-xs font-semibold uppercase tracking-wider text-slate-500">Catatan Ringkasan Aktivitas</h4>
          <p class="p-3 bg-white border border-slate-200 rounded-xl text-xs text-slate-700 leading-relaxed">
            Customer terdaftar melalui channel {{ historyCustomer.source || 'Tidak Ditemukan' }} dengan minat produk gadai. Status terakhir adalah {{ historyCustomer.status }}.
          </p>
        </div>

        <div class="flex justify-end pt-2">
          <button @click="historyCustomer = null" class="px-5 py-2 border border-slate-300 text-slate-600 rounded-xl text-sm font-medium">Tutup</button>
        </div>
      </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'

interface Customer {
  nama: string
  phone: string
  source: string
  status: string
}

const customers = ref<Customer[]>([])
const searchQuery = ref('')
const isLoading = ref(false)
const activeWaCustomer = ref<Customer | null>(null)
const waMessage = ref<[] | any>('')
const historyCustomer = ref<Customer | null>(null)

const waTemplates = [
  { title: 'Follow Up Penawaran Gadai', text: 'Halo Bpk/Ibu {nama}, salam dari PT Gadai Mulia. Mengenai simulasi pencairan gadai sebelumnya, apakah ada yang bisa kami bantu kembali?' },
  { title: 'Informasi Cek Taksiran Emas', text: 'Halo Bpk/Ibu {nama}, kami menginfokan bahwa harga acuan taksiran emas hari ini sedang bagus. Apabila berminat mensimulasikan kembali, silakan hubungi kami.' },
  { title: 'Greetings Customer Baru', text: 'Halo Bpk/Ibu {nama}, terima kasih telah berkonsultasi dengan Gadai. Kami siap memberikan layanan pencairan cepat dan sewa modal terjangkau.' }
]

const activeSbgCount = computed(() => {
  return customers.value.filter(c => (c.status || '').toLowerCase().includes('sbg')).length
})

onMounted(() => {
  fetchCustomers()
})

const fetchCustomers = async () => {
  isLoading.value = true

  const  appSource  = await useApi('/sumbercust')
  console.log('appSource', appSource)
  
  const { data } = await useApi('/inquiry/filter', {
    params: {
      search: searchQuery.value || '',
      start: 0,
      length: 300
    }
  })
  isLoading.value = false

  const rawList = Array.isArray(data) ? data : (data?.data || [])

  if (rawList) {
    const map = new Map()
    for (const item of rawList) {
      const name = item.nama || item.nm_customer || item.name || 'No Name'
      const phone = item.hp || item.phone || item.no_hp || '-'
      const status = item.nama_hslaktiv || item.status || 'Leads'
      const source = appSource.data.find((value : any) => { return value.kd_ref_cust === item.id_sumbercust})?.nm_ref_cust
      map.set(name, {
        nama: name,
        phone: phone,
        source: source,
        status: status
      })
    }
    customers.value = Array.from(map.values())
  }
}

let debounceTimer: any = null
const debounceFetch = () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    fetchCustomers()
  }, 400)
}

const openWaModal = (c: Customer) => {
  activeWaCustomer.value = c
  waMessage.value = waTemplates[0]?.text.replace('{nama}', c.nama) ?? ''
}

const openHistoryModal = (c: Customer) => {
  historyCustomer.value = c
}

const cleanPhone = (phoneStr: string) => {
  if (!phoneStr) return '6283138794243'
  let clean = phoneStr.replace(/\D/g, '')
  if (clean.startsWith('0')) {
    clean = '62' + clean.slice(1)
  }
  return clean || '6283138794243'
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
</script>
