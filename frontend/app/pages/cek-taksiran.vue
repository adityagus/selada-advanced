<template>
  <div class="space-y-6 pb-12 font-sans">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-xl md:text-2xl font-semibold text-slate-900 tracking-tight">Cek Taksiran</h1>
      </div>
      <div class="flex items-center gap-3">
        <!-- Main Mode Switcher: Simulasi vs History -->
        <div class="flex p-1.5 bg-slate-200/80 rounded-2xl gap-1">
          <button
            type="button"
            @click="mainMode = 'simulasi'"
            class="px-4 py-2 rounded-xl text-xs md:text-sm font-normal transition-all cursor-pointer flex items-center gap-2"
            :class="mainMode === 'simulasi' ? 'bg-white text-slate-900 shadow-sm font-semibold' : 'text-slate-600 hover:text-slate-900'"
          >
            <AppIcon name="calculator" class="w-4 h-4 text-amber-500" />
            Simulasi Calculator
          </button>
          <button
            type="button"
            @click="switchToHistoryMode"
            class="px-4 py-2 rounded-xl text-xs md:text-sm font-normal transition-all cursor-pointer flex items-center gap-2"
            :class="mainMode === 'history' ? 'bg-white text-slate-900 shadow-sm font-semibold' : 'text-slate-600 hover:text-slate-900'"
          >
            <AppIcon name="clock" class="w-4 h-4 text-emerald-500" />
            History Taksiran
          </button>
        </div>
      </div>
    </div>

    <!-- API Error Alert -->
    <p v-if="apiError" role="alert" class="rounded-xl border border-rose-200 bg-rose-50 p-4 text-sm text-rose-700">
      {{ apiError }}
    </p>

    <!-- MAIN MODE 1: SIMULASI CALCULATOR -->
    <div v-if="mainMode === 'simulasi'" class="bg-white rounded-3xl border border-slate-200/80 p-4 md:p-8 shadow-sm space-y-6">
      <!-- Tab Switcher -->
      <div class="flex p-1 bg-slate-100 rounded-2xl gap-1 max-w-md">
        <button
          type="button"
          @click="activeTab = 'elektronik'"
          class="flex-1 py-2.5 px-2 rounded-xl text-sm font-normal transition-all cursor-pointer flex items-center justify-center gap-2"
          :class="activeTab === 'elektronik' ? 'bg-white text-amber-700 shadow-sm border border-slate-200/60 font-semibold' : 'text-slate-600 hover:text-slate-900'"
        >
          <AppIcon name="monitor" class="w-4 h-4" />
          Elektronik
        </button>
        <button
          type="button"
          @click="activeTab = 'emas'"
          class="flex-1 py-2.5 px-2 rounded-xl text-sm font-normal transition-all cursor-pointer flex items-center justify-center gap-2"
          :class="activeTab === 'emas' ? 'bg-white text-amber-700 shadow-sm border border-slate-200/60 font-semibold' : 'text-slate-600 hover:text-slate-900'"
        >
          <AppIcon name="coin" class="w-4 h-4" />
          Emas
        </button>
      </div>

      <!-- Tab: Elektronik -->
      <div v-show="activeTab === 'elektronik'">
        <AppraisalElektronikForm
          ref="elekFormRef"
          @update:result="onElekResult"
          @error="handleApiError"
        />
      </div>

      <!-- Tab: Emas -->
      <div v-show="activeTab === 'emas'">
        <AppraisalEmasForm
          ref="emasFormRef"
          @update:savedList="onEmasSavedList"
          @update:totalTaksiran="onEmasTotalTaksiran"
          @error="handleApiError"
        />
      </div>

      <!-- Save As Lead Option & Submission -->
      <div class="pt-4 border-t border-slate-100 space-y-4">
        <div class="flex items-center justify-between">
          <label class="inline-flex items-center gap-2 cursor-pointer">
            <input type="checkbox" v-model="saveAsLead" class="w-4 h-4 rounded text-amber-500 focus:ring-amber-400" />
            <span class="text-xs md:text-sm font-medium text-slate-800">Simpan Hasil Taksiran Sebagai Inquiry Customer (Leads Baru)</span>
          </label>
        </div>

        <div v-if="saveAsLead" class="grid grid-cols-1 md:grid-cols-2 gap-4 bg-slate-50 p-4 rounded-2xl border border-slate-200">
          <div>
            <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Nama Customer</label>
            <input
              v-model="leadCustomerName"
              type="text"
              placeholder="Masukkan nama customer..."
              class="w-full px-3.5 py-2 bg-white border border-slate-300 rounded-xl text-sm focus:outline-none"
            />
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">No. Handphone / WA</label>
            <input
              v-model="leadCustomerPhone"
              type="text"
              placeholder="08xxxxxxxxxx"
              class="w-full px-3.5 py-2 bg-white border border-slate-300 rounded-xl text-sm focus:outline-none"
            />
          </div>
        </div>

        <div class="flex justify-end gap-3 pt-2">
          <button
            type="button"
            @click="handleSubmit"
            class="w-full sm:w-auto px-6 py-3 bg-amber-500 hover:bg-amber-600 text-white rounded-xl text-sm font-bold transition-all shadow-md shadow-amber-200 flex items-center justify-center gap-2 cursor-pointer"
            :disabled="isSubmitting"
          >
            <AppIcon name="loader" v-if="isSubmitting" class="w-4 h-4 animate-spin" />
            <span v-else>Simpan Transaksi Taksiran</span>
          </button>
        </div>

        <div v-if="statusMessage" class="p-3.5 rounded-xl text-sm font-medium text-center" :class="statusClass">
          {{ statusMessage }}
        </div>
      </div>
    </div>

    <!-- MAIN MODE 2: HISTORY TAKSIRAN -->
    <div v-else-if="mainMode === 'history'">
      <AppraisalHistoryTable ref="historyTableRef" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useApi } from '~/composables/useApi'
import type { ElekCalculationOutput } from '~/features/appraisal/components/AppraisalElektronikForm.vue'
import type { EmasRowItem } from '~/features/appraisal/components/AppraisalEmasForm.vue'

const apiError = ref('')
const mainMode = ref<'simulasi' | 'history'>('simulasi')
const activeTab = ref<'elektronik' | 'emas'>('elektronik')

const elekFormRef = ref()
const emasFormRef = ref()
const historyTableRef = ref()

// Shared Lead & Submission State
const saveAsLead = ref(false)
const leadCustomerName = ref('')
const leadCustomerPhone = ref('')
const isSubmitting = ref(false)
const statusMessage = ref('')
const statusClass = ref('')

// State from child components
const elekResultData = ref<ElekCalculationOutput>({
  detailName: '',
  maxLoan: 0,
  requestedLoan: 0,
  hps: 0,
  score: 100,
  grade: 'A',
  isRejected: false,
  isReady: false
})
const emasSavedList = ref<EmasRowItem[]>([])
const totalEmasTaksiran = ref(0)

const onElekResult = (output: ElekCalculationOutput) => {
  elekResultData.value = output
  apiError.value = ''
}

const onEmasSavedList = (list: EmasRowItem[]) => {
  emasSavedList.value = list
}

const onEmasTotalTaksiran = (total: number) => {
  totalEmasTaksiran.value = total
  apiError.value = ''
}

const handleApiError = (msg: string) => {
  apiError.value = msg
}

const switchToHistoryMode = () => {
  mainMode.value = 'history'
  historyTableRef.value?.fetchTaksiranHistory?.()
}

const formatNum = (val: number) => {
  if (!val && val !== 0) return '0'
  return new Intl.NumberFormat('id-ID').format(Math.round(val))
}

const handleSubmit = async () => {
  const isElek = activeTab.value === 'elektronik'
  if (isElek ? elekResultData.value.requestedLoan <= 0 : totalEmasTaksiran.value <= 0) {
    apiError.value = 'Lengkapi simulasi hingga nilai pinjaman tersedia sebelum menyimpan.'
    return
  }

  isSubmitting.value = true
  statusMessage.value = ''

  const nilaiTaksir = isElek ? elekResultData.value.maxLoan : totalEmasTaksiran.value
  const nilaiPinjaman = isElek ? elekResultData.value.requestedLoan : Math.round(totalEmasTaksiran.value * 0.95)

  try {
    const { error } = await useApi('/taksiran/simpan', {
      method: 'POST',
      body: {
        jenis_barang: activeTab.value,
        nama_barang: isElek ? elekResultData.value.detailName : 'Paket Barang Emas',
        nilai_taksir: nilaiTaksir,
        nilai_pinjaman: nilaiPinjaman,
        items_emas: isElek ? [] : emasSavedList.value
      }
    })

    if (saveAsLead.value && leadCustomerName.value) {
      await useApi('/rencana/input', {
        method: 'POST',
        body: {
          cust: 'Baru',
          nm: leadCustomerName.value,
          hp: leadCustomerPhone.value || '081234567890',
          sumbercust: 'Walk In',
          hslaktiv: 1,
          katcust: '1',
          ket: `Simulasi Taksiran ${isElek ? elekResultData.value.detailName : 'Emas'} - Total Taksiran Rp ${formatNum(nilaiTaksir)}`
        }
      })
    }

    isSubmitting.value = false

    if (error) {
      statusMessage.value = 'Gagal menyimpan transaksi taksiran: ' + error
      statusClass.value = 'bg-rose-50 border border-rose-200 text-rose-700 font-bold'
    } else {
      statusMessage.value = 'Simulasi taksiran gadai berhasil disimpan!'
      statusClass.value = 'bg-emerald-50 border border-emerald-200 text-emerald-700 font-bold'
      // Reset child forms
      elekFormRef.value?.reset?.()
      emasFormRef.value?.reset?.()
      leadCustomerName.value = ''
      leadCustomerPhone.value = ''
      saveAsLead.value = false
    }
  } catch (e: any) {
    isSubmitting.value = false
    statusMessage.value = 'Terjadi kesalahan sistem'
    statusClass.value = 'bg-rose-50 border border-rose-200 text-rose-700 font-bold'
  }
}
</script>
