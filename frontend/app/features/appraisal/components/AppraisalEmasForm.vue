<template>
  <div class="space-y-6">
    <!-- Form Input Barang Emas Card -->
    <div class="border border-slate-200 rounded-2xl p-4 md:p-6 bg-white space-y-4">
      <p class="font-bold text-slate-800 text-base md:text-lg">Tambah Barang Emas</p>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div>
          <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Nama Barang</label>
          <select
            v-model="emasDetailId"
            @change="onEmasDetailChange"
            class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-800 focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all"
          >
            <option value="">Pilih nama barang</option>
            <option
              v-for="b in emasOptionBarang"
              :key="b.id_barang_detail"
              :value="b.id_barang_detail"
            >
              {{ b.nama_barang || b.kode_barang }}
            </option>
          </select>
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Jenis Barang</label>
          <input
            type="text"
            v-model="emasCollateralName"
            readonly
            placeholder="Jenis Emas"
            class="w-full px-3.5 py-2.5 bg-slate-100 border border-slate-200 rounded-xl text-sm text-slate-700 cursor-not-allowed font-medium"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Karat (Min: {{ emasMinKarat }}K)</label>
          <input
            type="number"
            v-model.number="emasKarat"
            @input="calcEmasRowPrice"
            :readonly="isEmasLM"
            min="1"
            max="24"
            placeholder="0"
            class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-800 focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all"
            :class="isEmasLM ? 'bg-slate-100 cursor-not-allowed' : ''"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Jumlah</label>
          <input
            type="number"
            v-model.number="emasTotal"
            @input="calcEmasRowPrice"
            :readonly="isEmasLM"
            min="1"
            placeholder="0"
            class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-800 focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all"
            :class="isEmasLM ? 'bg-slate-100 cursor-not-allowed' : ''"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Berat Kotor (gr)</label>
          <input
            type="number"
            step="0.01"
            v-model.number="emasGrossWeight"
            placeholder="0"
            class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-800 focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Berat Bersih (gr)</label>
          <input
            type="number"
            step="0.01"
            v-model.number="emasNetWeight"
            @input="calcEmasRowPrice"
            placeholder="0"
            class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-800 focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all"
          />
        </div>

        <div class="md:col-span-2">
          <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Keterangan Barang</label>
          <input
            type="text"
            v-model="emasNotes"
            placeholder="Masukkan keterangan"
            class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-800 focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all"
          />
        </div>

        <div class="flex items-end">
          <button
            type="button"
            @click="saveEmasRow"
            class="w-full py-2.5 px-4 bg-amber-500 hover:bg-amber-600 text-white rounded-xl text-sm font-semibold transition-all shadow-sm flex items-center justify-center gap-2 cursor-pointer"
          >
            <AppIcon name="plus" class="w-4 h-4" />
            Tambah
          </button>
        </div>
      </div>

      <!-- Real-Time Row Price Preview -->
      <div class="flex items-center justify-between p-3.5 bg-slate-50 border border-slate-200 rounded-xl text-sm">
        <span class="text-slate-600 font-medium">Nilai Taksir :</span>
        <span class="text-lg md:text-xl font-extrabold text-amber-600">Rp {{ formatNum(currentEmasRowPrice) }}</span>
      </div>
    </div>

    <!-- Saved Emas Items List -->
    <div v-if="emasSavedList.length > 0" class="space-y-3 pt-2">
      <p class="font-bold text-slate-800 text-sm">Daftar Barang ({{ emasSavedList.length }})</p>

      <div class="space-y-2.5">
        <div
          v-for="(item, idx) in emasSavedList"
          :key="idx"
          class="p-4 bg-white border border-slate-200 rounded-2xl shadow-2xs flex flex-col md:flex-row md:items-center justify-between gap-3"
        >
          <div class="space-y-1">
            <div class="flex items-center gap-2">
              <strong class="text-slate-800 text-sm">{{ item.name }}</strong>
              <span v-if="item.notes" class="text-xs text-slate-500"> - ({{ item.notes }})</span>
            </div>
            <p class="text-xs text-slate-500 font-medium">
              {{ item.karat }} Karat • Jumlah: {{ item.total }} pcs • Bruto: {{ item.gross_weight }}g • Netto: {{ item.net_weight }}g
            </p>
          </div>

          <div class="flex items-center justify-between md:justify-end gap-4 border-t md:border-t-0 pt-2 md:pt-0 border-slate-100">
            <div class="text-left md:text-right">
              <span class="text-[10px] text-slate-400 block uppercase">Taksiran</span>
              <span class="text-base md:text-lg font-extrabold text-amber-600">Rp {{ formatNum(item.price) }}</span>
            </div>
            <button
              type="button"
              @click="removeEmasSavedItem(idx)"
              class="p-2 text-rose-500 hover:bg-rose-50 rounded-lg transition-all cursor-pointer"
              title="Hapus barang"
            >
              <AppIcon name="trash" class="w-5 h-5" />
            </button>
          </div>
        </div>
      </div>

      <!-- Total Ringkasan Emas -->
      <div class="border border-slate-200 p-4 rounded-2xl bg-slate-50/80 flex items-center justify-between">
        <span class="text-sm font-bold text-slate-800">Total Taksiran :</span>
        <span class="text-lg md:text-xl font-extrabold text-amber-600">Rp {{ formatNum(totalEmasTaksiran) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'

export interface EmasRowItem {
  id_detail: string | number
  name: string
  collateral_name: string
  karat: number
  total: number
  gross_weight: number
  net_weight: number
  notes: string
  price: number
}

const emit = defineEmits<{
  (e: 'update:savedList', list: EmasRowItem[]): void
  (e: 'update:totalTaksiran', total: number): void
  (e: 'error', message: string): void
}>()

const emasOptionBarang = ref<any[]>([])
const emasDetailId = ref<number | string>('')
const emasCollateralName = ref('')
const emasKarat = ref<number>(24)
const emasMinKarat = ref<number>(8)
const emasTotal = ref<number>(1)
const emasGrossWeight = ref<number>(0)
const emasNetWeight = ref<number>(0)
const emasNotes = ref('')
const isEmasLM = ref(false)
const currentEmasRowPrice = ref<number>(0)

const emasSavedList = ref<EmasRowItem[]>([])

const totalEmasTaksiran = computed(() => {
  return emasSavedList.value.reduce((acc, item) => acc + (item.price || 0), 0)
})

const loadEmasBarang = async () => {
  try {
    const { data } = await useApi('/taksiran/barang-emas')
    emasOptionBarang.value = Array.isArray(data) ? data : []
  } catch (e) {
    console.error('Failed to load gold items:', e)
  }
}

const onEmasDetailChange = () => {
  const item = emasOptionBarang.value.find(i => String(i.id_barang_detail) === String(emasDetailId.value))
  if (!item) {
    emasCollateralName.value = ''
    isEmasLM.value = false
    return
  }

  emasCollateralName.value = item.nama_barang || item.kode_barang || 'Emas'
  const lowerName = (emasCollateralName.value).toLowerCase()

  if (lowerName.includes('antam') || lowerName.includes('logam mulia') || lowerName.includes('lm')) {
    isEmasLM.value = true
    emasKarat.value = 24
    emasTotal.value = 1
  } else {
    isEmasLM.value = false
    if (!emasKarat.value) emasKarat.value = 24
    if (!emasTotal.value) emasTotal.value = 1
  }
  calcEmasRowPrice()
}

const calcEmasRowPrice = async () => {
  if (!emasDetailId.value || !emasKarat.value || !emasNetWeight.value || emasNetWeight.value <= 0) {
    currentEmasRowPrice.value = 0
    return
  }

  try {
    const { data } = await useApi('/taksiran/gold-stle', {
      params: {
        id_barang_detail: emasDetailId.value,
        karat: emasKarat.value
      }
    })

    if (data && data.min_karat) {
      emasMinKarat.value = parseFloat(data.min_karat) || 8
    }

    if (emasKarat.value < emasMinKarat.value || emasTotal.value <= 0 || emasNetWeight.value <= 0) {
      currentEmasRowPrice.value = 0
      return
    }

    const stle = Number(data?.stle) || 0
    if (stle <= 0) {
      currentEmasRowPrice.value = 0
      emit('error', 'Harga taksiran emas belum tersedia untuk barang dan karat ini.')
      return
    }

    currentEmasRowPrice.value = (emasKarat.value / 24) * emasNetWeight.value * stle * emasTotal.value
  } catch (e) {
    console.error('Failed to calculate gold row price:', e)
    currentEmasRowPrice.value = 0
  }
}

const saveEmasRow = () => {
  if (!emasDetailId.value) {
    alert('Pilih nama barang emas terlebih dahulu.')
    return
  }
  if (emasKarat.value < emasMinKarat.value) {
    alert(`Karat minimal adalah ${emasMinKarat.value}!`)
    return
  }
  if (currentEmasRowPrice.value <= 0) {
    alert('Pastikan Karat, Jumlah, dan Berat Bersih terisi dengan benar.')
    return
  }

  const itemObj = emasOptionBarang.value.find(i => String(i.id_barang_detail) === String(emasDetailId.value))
  emasSavedList.value.push({
    id_detail: emasDetailId.value,
    name: itemObj ? itemObj.nama_barang : 'Barang Emas',
    collateral_name: emasCollateralName.value,
    karat: emasKarat.value,
    total: emasTotal.value,
    gross_weight: emasGrossWeight.value,
    net_weight: emasNetWeight.value,
    notes: emasNotes.value,
    price: currentEmasRowPrice.value
  })

  // Clear row form
  emasDetailId.value = ''
  emasCollateralName.value = ''
  emasKarat.value = 24
  emasTotal.value = 1
  emasGrossWeight.value = 0
  emasNetWeight.value = 0
  emasNotes.value = ''
  isEmasLM.value = false
  currentEmasRowPrice.value = 0

  emitUpdates()
}

const removeEmasSavedItem = (idx: number) => {
  emasSavedList.value.splice(idx, 1)
  emitUpdates()
}

const emitUpdates = () => {
  emit('update:savedList', [...emasSavedList.value])
  emit('update:totalTaksiran', totalEmasTaksiran.value)
}

const reset = () => {
  emasSavedList.value = []
  emasDetailId.value = ''
  emasCollateralName.value = ''
  emasKarat.value = 24
  emasTotal.value = 1
  emasGrossWeight.value = 0
  emasNetWeight.value = 0
  emasNotes.value = ''
  isEmasLM.value = false
  currentEmasRowPrice.value = 0
  emitUpdates()
}

defineExpose({
  reset
})

const formatNum = (val: number) => {
  if (!val && val !== 0) return '0'
  return new Intl.NumberFormat('id-ID').format(Math.round(val))
}

onMounted(async () => {
  await loadEmasBarang()
})
</script>
