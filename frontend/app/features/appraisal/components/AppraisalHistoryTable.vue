<template>
  <div class="bg-white rounded-3xl border border-slate-200/80 p-4 md:p-8 shadow-sm space-y-6">
    <!-- Filter Bar: Sales User Filter & Search -->
    <div class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-4 bg-slate-50 p-4 rounded-2xl border border-slate-200">
      <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 flex-1">
        <div class="w-full sm:w-64">
          <label class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-1.5">Pembuat (Sales User)</label>
          <select
            v-model="historySalesUser"
            @change="fetchTaksiranHistory"
            class="w-full px-3.5 py-2 bg-white border border-slate-300 rounded-xl text-sm font-semibold text-slate-800 focus:outline-none focus:border-amber-500"
          >
            <option value="all">Semua Sales / User</option>
            <option
              v-for="u in salesUsers"
              :key="u.id_login || u.username"
              :value="u.id_login || u.username"
            >
              {{ u.nama || u.username }} ({{ u.id_login || u.username }})
            </option>
          </select>
        </div>

        <div class="flex-1">
          <label class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-1.5">Cari Barang / User</label>
          <div class="flex gap-2">
            <input
              v-model="historySearch"
              @keyup.enter="fetchTaksiranHistory"
              type="text"
              placeholder="Ketik nama barang / user..."
              class="w-full px-3.5 py-2 bg-white border border-slate-300 rounded-xl text-sm text-slate-800 focus:outline-none focus:border-amber-500"
            />
            <button
              type="button"
              @click="fetchTaksiranHistory"
              class="px-4 py-2 bg-amber-500 hover:bg-amber-600 text-white rounded-xl text-sm font-semibold transition-all cursor-pointer"
            >
              Cari
            </button>
          </div>
        </div>
      </div>

      <div class="flex items-end">
        <button
          type="button"
          @click="fetchTaksiranHistory"
          class="px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-xl text-sm font-semibold transition-all border border-slate-200 flex items-center gap-2 cursor-pointer"
        >
          <AppIcon name="refresh" class="w-4 h-4 text-slate-500" />
          Refresh
        </button>
      </div>
    </div>

    <!-- DESKTOP TABLE VIEW -->
    <div class="hidden md:block overflow-x-auto border border-slate-200 rounded-2xl bg-white shadow-2xs">
      <table class="w-full text-sm text-left border-collapse min-w-[700px]">
        <thead>
          <tr class="bg-slate-100/80 border-b border-slate-200 text-slate-700 font-bold uppercase text-xs">
            <th class="p-3.5 text-center w-12">No</th>
            <th class="p-3.5">Tanggal & Waktu</th>
            <th class="p-3.5">Pembuat (Sales)</th>
            <th class="p-3.5">Jenis Barang</th>
            <th class="p-3.5">Nama Barang</th>
            <th class="p-3.5 text-right">Nilai Taksir</th>
            <th class="p-3.5 text-right">Nilai Pinjaman (UP)</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="isLoadingHistory">
            <td colspan="7" class="p-8 text-center text-slate-400 text-sm">Memuat riwayat taksiran...</td>
          </tr>
          <tr v-else-if="historyList.length === 0">
            <td colspan="7" class="p-8 text-center text-slate-400 text-sm">Belum ada riwayat simulasi taksiran untuk user yang dipilih.</td>
          </tr>
          <tr
            v-else
            v-for="(item, idx) in historyList"
            :key="item.id || idx"
            class="border-b border-slate-100 hover:bg-slate-50 transition-all"
          >
            <td class="p-3.5 text-center font-medium text-slate-500">{{ idx + 1 }}</td>
            <td class="p-3.5 text-slate-700 font-medium">{{ formatDate(item.created_at) }}</td>
            <td class="p-3.5">
              <span class="font-bold text-slate-800 bg-slate-100 px-2 py-1 rounded border border-slate-200 uppercase text-xs">{{ item.created_by }}</span>
            </td>
            <td class="p-3.5">
              <span
                class="font-bold px-2 py-1 rounded text-xs uppercase"
                :class="item.jenis_barang === 'emas' ? 'bg-amber-100 text-amber-800' : 'bg-blue-100 text-blue-800'"
              >
                {{ item.jenis_barang }}
              </span>
            </td>
            <td class="p-3.5 font-semibold text-slate-800">{{ item.nama_barang }}</td>
            <td class="p-3.5 text-right font-bold text-slate-800">Rp {{ formatNum(item.nilai_taksir) }}</td>
            <td class="p-3.5 text-right font-extrabold text-amber-600">Rp {{ formatNum(item.nilai_pinjaman) }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- MOBILE CARDS VIEW -->
    <div class="block md:hidden space-y-3">
      <div v-if="isLoadingHistory" class="p-8 text-center text-slate-400 text-sm">
        Memuat riwayat taksiran...
      </div>
      <div v-else-if="historyList.length === 0" class="p-8 text-center text-slate-400 text-sm">
        Belum ada riwayat simulasi taksiran.
      </div>
      <div
        v-else
        v-for="(item, idx) in historyList"
        :key="item.id || idx"
        class="p-4 bg-white border border-slate-200 rounded-2xl shadow-2xs space-y-2.5"
      >
        <div class="flex items-center justify-between border-b border-slate-100 pb-2">
          <span class="font-bold text-slate-800 bg-slate-100 px-2 py-0.5 rounded border border-slate-200 uppercase text-xs">{{ item.created_by }}</span>
          <span class="text-xs text-slate-400 font-medium">{{ formatDate(item.created_at) }}</span>
        </div>

        <div class="flex items-center justify-between">
          <span class="font-bold text-slate-900 text-sm">{{ item.nama_barang }}</span>
          <span
            class="font-bold px-2 py-0.5 rounded text-[10px] uppercase"
            :class="item.jenis_barang === 'emas' ? 'bg-amber-100 text-amber-800' : 'bg-blue-100 text-blue-800'"
          >
            {{ item.jenis_barang }}
          </span>
        </div>

        <div class="flex items-center justify-between pt-2 border-t border-slate-100 text-xs">
          <div>
            <span class="text-slate-400 block text-[10px] uppercase">Nilai Taksir</span>
            <span class="font-bold text-slate-700">Rp {{ formatNum(item.nilai_taksir) }}</span>
          </div>
          <div class="text-right">
            <span class="text-slate-400 block text-[10px] uppercase">Nilai Pinjaman</span>
            <span class="font-extrabold text-amber-600 text-sm">Rp {{ formatNum(item.nilai_pinjaman) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'

const salesUsers = ref<any[]>([])
const historySalesUser = ref('all')
const historySearch = ref('')
const historyList = ref<any[]>([])
const isLoadingHistory = ref(false)

const fetchSalesUsers = async () => {
  try {
    const { data } = await useApi('/sales/list')
    salesUsers.value = Array.isArray(data) ? data : []
  } catch (e) {
    console.error('Failed to load sales users:', e)
  }
}

const fetchTaksiranHistory = async () => {
  isLoadingHistory.value = true
  try {
    const { data } = await useApi('/taksiran/history', {
      params: {
        sales_user: historySalesUser.value,
        search: historySearch.value
      }
    })
    if (data && Array.isArray(data.data)) {
      historyList.value = data.data
    } else {
      historyList.value = []
    }
  } catch (e) {
    console.error('Failed to fetch taksiran history:', e)
    historyList.value = []
  } finally {
    isLoadingHistory.value = false
  }
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  try {
    const d = new Date(dateStr)
    return d.toLocaleString('id-ID', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return dateStr
  }
}

const formatNum = (val: number) => {
  if (!val && val !== 0) return '0'
  return new Intl.NumberFormat('id-ID').format(Math.round(val))
}

defineExpose({
  fetchTaksiranHistory
})

onMounted(async () => {
  await fetchSalesUsers()
  await fetchTaksiranHistory()
})
</script>
