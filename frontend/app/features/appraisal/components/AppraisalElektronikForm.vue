<template>
  <div class="space-y-6">
    <!-- Detail Barang Form Card -->
    <div class="border border-slate-200 rounded-2xl p-4 md:p-6 bg-white space-y-4">
      <span class="font-bold text-slate-800 text-base md:text-lg block">Detail Barang</span>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div>
          <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Jenis Barang</label>
          <select
            v-model="elekJenisBarang"
            @change="onElekJenisChange"
            class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-800 focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all"
          >
            <option value="">Pilih jenis barang</option>
            <option
              v-for="cat in elekCategories"
              :key="cat.kd_jenis_barang"
              :value="cat.kd_jenis_barang"
            >
              {{ cat.nm_jenis_barang }}
            </option>
          </select>
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Nama Barang</label>
          <div class="flex gap-2">
            <input
              type="text"
              v-model="elekDetailName"
              readonly
              placeholder="Pilih nama barang"
              class="flex-1 px-3.5 py-2.5 bg-slate-50 border border-slate-300 rounded-xl text-sm text-slate-800 cursor-pointer focus:outline-none"
              @click="openElekSearchModal"
            />
            <button
              type="button"
              @click="openElekSearchModal"
              class="px-3.5 py-2.5 bg-slate-100 hover:bg-slate-200 text-slate-700 border border-slate-300 rounded-xl text-sm font-medium transition-all cursor-pointer"
            >
              <AppIcon name="search" class="w-4 h-4 text-slate-600" />
            </button>
          </div>
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Wilayah</label>
          <select
            v-model="elekWilayah"
            @change="loadElekNilaiTaksir"
            class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-800 focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all"
          >
            <option value="">- Pilih Wilayah -</option>
            <option
              v-for="w in elekWilayahList"
              :key="w.region_id"
              :value="w.region_id"
            >
              {{ w.region_name }}
            </option>
          </select>
        </div>
      </div>

      <!-- Dynamic Assessment / Scoring Matrix Area -->
      <div v-if="elekScoringSections.length > 0" class="pt-4 border-t border-dashed border-slate-200 space-y-6">
        <div
          v-for="(section, sIdx) in elekScoringSections"
          :key="section.section_id"
          class="space-y-3"
        >
          <span class="font-bold text-slate-800 text-xs md:text-sm uppercase tracking-wider block text-amber-700">{{ section.section_name }}</span>

          <!-- SECTION TYPE: REJECT -->
          <div v-if="section.section_type === 'reject'" class="border border-slate-200 rounded-2xl p-4 bg-slate-50/50 space-y-3">
            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
              <span class="text-sm font-semibold text-slate-700">Item Status</span>
              <select
                v-model="section.selectedOption"
                @change="calculateElekResult"
                class="px-3.5 py-2 bg-white border border-slate-300 rounded-xl text-sm font-medium text-slate-800 focus:outline-none"
              >
                <option :value="null">Pilih status</option>
                <option
                  v-for="(opt, oIdx) in section.options"
                  :key="opt.option_id"
                  :value="oIdx"
                >
                  {{ opt.option_name }}
                </option>
              </select>
            </div>
          </div>

          <!-- SECTION TYPE: POSITIF / NEGATIF -->
          <div v-else class="space-y-3">
            <!-- 1. DESKTOP VIEW: HTML TABLE (hidden on mobile) -->
            <div class="hidden md:block overflow-x-auto border border-slate-200 rounded-2xl bg-white shadow-2xs">
              <table class="w-full text-sm text-left border-collapse min-w-[650px]">
                <thead>
                  <tr class="bg-blue-50/70 border-b border-slate-200">
                    <th class="p-3.5 font-bold text-blue-900 uppercase border-r border-slate-200" style="width: 35%;">Item</th>
                    <th
                      v-for="opt in section.options"
                      :key="opt.option_id"
                      class="p-3.5 font-bold text-blue-900 text-center uppercase text-nowrap border-r border-slate-200"
                    >
                      {{ opt.option_name }}
                    </th>
                    <th class="p-3.5 font-bold text-blue-900 text-center uppercase text-nowrap w-24">Skor</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="(item, iIdx) in section.items"
                    :key="item.item_id"
                    class="border-b border-slate-100 hover:bg-blue-50/30 transition-all"
                  >
                    <td class="p-3.5 font-semibold text-slate-800 uppercase border-r border-slate-200 align-middle">
                      {{ item.item_name }}
                    </td>
                    <td
                      v-for="(opt, oIdx) in section.options"
                      :key="opt.option_id"
                      @click="selectElekItemScore(sIdx, iIdx, oIdx)"
                      class="p-3.5 text-center cursor-pointer border-r border-slate-200 align-middle hover:bg-blue-100/40"
                    >
                      <input
                        type="radio"
                        :name="'d_elek_' + sIdx + '_' + iIdx"
                        :checked="item.selected === oIdx"
                        class="w-4 h-4 text-amber-500 cursor-pointer pointer-events-none"
                      />
                    </td>
                    <td class="p-3.5 font-bold text-center text-slate-800 align-middle">
                      {{ getItemRowScore(section, item) }}%
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- 2. MOBILE VIEW: RESPONSIVE CARDS (visible only on mobile) -->
            <div class="block md:hidden space-y-3">
              <div
                v-for="(item, iIdx) in section.items"
                :key="item.item_id"
                class="p-3.5 bg-white border border-slate-200 rounded-2xl shadow-2xs space-y-2.5"
              >
                <div class="flex items-center justify-between pb-2 border-b border-slate-100">
                  <span class="font-bold text-slate-800 text-xs uppercase">{{ item.item_name }}</span>
                  <span class="font-extrabold text-amber-700 text-xs bg-amber-50 px-2 py-0.5 rounded border border-amber-200">
                    {{ getItemRowScore(section, item) }}%
                  </span>
                </div>

                <div class="grid grid-cols-2 gap-2">
                  <label
                    v-for="(opt, oIdx) in section.options"
                    :key="opt.option_id"
                    @click="selectElekItemScore(sIdx, iIdx, oIdx)"
                    class="flex items-center gap-2 p-2.5 rounded-xl border text-xs font-medium transition-all cursor-pointer"
                    :class="item.selected === oIdx ? 'bg-amber-50 border-amber-500 text-amber-900 font-bold shadow-2xs ring-1 ring-amber-400' : 'bg-slate-50 border-slate-200 text-slate-700 hover:bg-slate-100'"
                  >
                    <input
                      type="radio"
                      :name="'m_elek_' + sIdx + '_' + iIdx"
                      :checked="item.selected === oIdx"
                      class="w-4 h-4 text-amber-500 cursor-pointer pointer-events-none"
                    />
                    <span class="truncate">{{ opt.option_name }}</span>
                  </label>
                </div>
              </div>
            </div>

            <!-- Positif Section Total Summary -->
            <div v-if="section.section_type === 'positif' && (section.items || []).length > 1" class="flex items-center justify-between p-3.5 bg-slate-50 border border-slate-200 rounded-xl text-xs md:text-sm font-bold text-slate-800">
              <span>Total Skor {{ section.section_name }}</span>
              <span class="text-amber-700">{{ getSectionTotalPositifScore(section) }}%</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Result Summary Bar -->
      <div v-if="elekShowResult" class="pt-4 border-t border-slate-200 space-y-4">
        <div class="grid grid-cols-2 md:grid-cols-4 gap-3 border border-slate-200 rounded-2xl p-4 bg-slate-50/80 text-center">
          <div class="border-r border-slate-200/80 pr-2">
            <small class="text-slate-500 block text-xs font-semibold uppercase">HPS</small>
            <span class="font-extrabold text-lg md:text-xl text-slate-800">Rp {{ formatNum(elekResult.hps) }}</span>
          </div>
          <div class="border-r border-slate-200/80 pr-2">
            <small class="text-slate-500 block text-xs font-semibold uppercase">Skor</small>
            <span class="font-extrabold text-lg md:text-xl text-amber-600">{{ elekResult.score }}%</span>
          </div>
          <div class="border-r border-slate-200/80 pr-2">
            <small class="text-slate-500 block text-xs font-semibold uppercase">Grade</small>
            <span class="font-extrabold text-lg md:text-xl" :class="elekResult.isRejected ? 'text-rose-600' : 'text-emerald-600'">{{ elekResult.grade }}</span>
          </div>
          <div>
            <small class="text-slate-500 block text-xs font-semibold uppercase">Keterangan</small>
            <span class="font-bold text-xs md:text-sm text-slate-700 leading-tight block mt-0.5">{{ elekResult.note }}</span>
          </div>
        </div>

        <!-- Ringkasan Max Loan & UP -->
        <div class="space-y-2">
          <p class="font-bold text-slate-800 text-sm">Ringkasan Nilai Pinjaman</p>
          <div class="flex flex-col md:flex-row md:items-center justify-between gap-3 border border-slate-200 rounded-2xl p-4 bg-slate-50/80">
            <div class="flex items-center justify-between md:justify-start gap-2 text-sm">
              <span class="text-slate-600 font-medium">Max Uang Pinjaman :</span>
              <span class="text-amber-600 font-extrabold text-lg md:text-xl">Rp {{ formatNum(elekResult.maxLoan) }}</span>
            </div>
            <div class="flex items-center justify-between md:justify-start gap-2 text-sm">
              <span class="text-slate-500 font-medium">Total Pinjaman (UP):</span>
              <div class="flex items-center bg-white border border-slate-300 rounded-xl overflow-hidden shadow-2xs">
                <span class="px-3.5 py-1.5 font-extrabold text-emerald-600 text-base">Rp {{ formatNum(elekResult.requestedLoan) }}</span>
                <span class="bg-slate-100 px-2.5 py-1.5 text-xs text-slate-600 font-bold border-l border-slate-200">LTV 95%</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ELEKTRONIK CATALOG / SEARCH ITEM MODAL -->
    <ClientOnly>
      <Teleport to="body">
        <div v-if="showElekSearchModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs z-50 flex items-center justify-center p-4 !mt-0">
          <div class="bg-white rounded-3xl max-w-xl w-full p-6 shadow-2xl space-y-4 max-h-[85vh] flex flex-col">
            <div class="flex items-center justify-between border-b border-slate-100 pb-3">
              <h3 class="font-bold text-slate-800 text-base md:text-lg">Pilih Barang</h3>
              <button @click="showElekSearchModal = false" class="text-slate-400 hover:text-slate-600 p-1 text-lg">
                <AppIcon name="close" class="w-4 h-4" />
              </button>
            </div>

            <div class="flex gap-2">
              <input
                v-model="elekSearchQuery"
                type="text"
                @keyup.enter="searchElekItems"
                placeholder="Cari nama barang..."
                class="flex-1 px-3.5 py-2 bg-white border border-slate-300 rounded-xl text-sm focus:outline-none focus:border-amber-500"
                autofocus
              />
              <button @click="searchElekItems" class="px-4 py-2 bg-amber-500 hover:bg-amber-600 text-white font-medium rounded-xl text-sm cursor-pointer">
                Cari
              </button>
            </div>

            <div class="flex-1 overflow-y-auto space-y-2 pr-1">
              <div v-if="isSearchingElek" class="py-8 text-center text-slate-400 text-sm">
                Mencari data barang...
              </div>
              <div v-else-if="elekSearchResult.length === 0 && searchDisabled === true" class="py-8 text-center text-slate-400 text-sm">
                Tidak ada barang ditemukan. Silakan cari dengan kata kunci lain.
              </div>
              <div
                v-else
                v-for="item in elekSearchResult"
                :key="item.id_barang_detail"
                @click="selectElekCatalogItem(item)"
                class="p-3.5 rounded-2xl border border-slate-200 hover:border-amber-400 hover:bg-amber-50/50 transition-all cursor-pointer flex justify-between items-center"
              >
                <div>
                  <p class="text-sm font-semibold text-slate-800">{{ item.nm_barang || item.nama }}</p>
                  <p class="text-xs text-slate-500 font-medium">Kode: {{ item.kd_barang || item.code }}</p>
                </div>
                <span class="text-xs font-semibold text-amber-600 bg-amber-50 px-3 py-1 rounded-lg border border-amber-200/60">Pilih</span>
              </div>
            </div>
          </div>
        </div>
      </Teleport>
    </ClientOnly>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'

export interface ElekCalculationOutput {
  detailName: string
  maxLoan: number
  requestedLoan: number
  hps: number
  score: number
  grade: string
  isRejected: boolean
  isReady: boolean
}

const emit = defineEmits<{
  (e: 'update:result', output: ElekCalculationOutput): void
  (e: 'error', message: string): void
}>()

const elekCategories = ref<any[]>([])
const elekJenisBarang = ref('')
const searchDisabled = ref(true)
const elekDetailId = ref<number | string>('')
const elekCollateralItemId = ref<number | string>('')
const elekDetailName = ref('')

const elekWilayahList = ref<any[]>([])
const elekWilayah = ref('')

const showElekSearchModal = ref(false)
const elekSearchQuery = ref('')
const isSearchingElek = ref(false)
const elekSearchResult = ref<any[]>([])

const elekScoringSections = ref<any[]>([])
const elekShowResult = ref(false)
const elekResult = ref({
  hps: 0,
  score: 100,
  isRejected: false,
  grade: 'A',
  note: 'Kondisi barang prima / 100%',
  maxLoan: 0,
  requestedLoan: 0
})

const onElekJenisChange = () => {
  elekSearchResult.value = []
  searchDisabled.value = true
  elekDetailId.value = ''
  elekCollateralItemId.value = ''
  elekDetailName.value = ''
  elekWilayah.value = ''
  elekWilayahList.value = []
  elekScoringSections.value = []
  elekShowResult.value = false
  elekResult.value = { hps: 0, score: 100, isRejected: false, grade: 'A', note: 'Kondisi barang prima / 100%', maxLoan: 0, requestedLoan: 0 }

  emitOutput()
}

const openElekSearchModal = () => {
  showElekSearchModal.value = true
  if (elekSearchResult.value.length === 0) {
    searchElekItems()
  }
}

const searchElekItems = async () => {
  isSearchingElek.value = true
  try {
    const { data } = await useApi('/taksiran/barang-umum', {
      params: {
        jenisBarang: elekJenisBarang.value,
        namaBarang: elekSearchQuery.value
      }
    })
    elekSearchResult.value = Array.isArray(data) ? data : []
  } catch (e) {
    console.error('Failed to search electronic items:', e)
    elekSearchResult.value = []
  } finally {
    isSearchingElek.value = false
  }
}

const selectElekCatalogItem = (item: any) => {
  onElekJenisChange()
  elekDetailId.value = item.id_barang_detail || item.id || ''
  elekCollateralItemId.value = item.collateral_item_id || elekDetailId.value
  elekDetailName.value = item.nm_barang || item.name || ''
  showElekSearchModal.value = false

  loadElekWilayah(elekDetailId.value)
  loadElekScoring(elekCollateralItemId.value)
}

const loadElekWilayah = async (detailId: any) => {
  if (!detailId) return
  try {
    const { data } = await useApi('/taksiran/wilayah', {
      params: { wilayah: detailId, type: 'umum' }
    })
    elekWilayahList.value = Array.isArray(data) ? data : []
    if (elekWilayahList.value.length === 1 && !elekWilayah.value) {
      elekWilayah.value = elekWilayahList.value[0].region_id
      await loadElekNilaiTaksir()
    }
  } catch (e) {
    console.error('Failed to load regions:', e)
  }
}

const initScoringTo100 = () => {
  for (const sec of elekScoringSections.value) {
    if (sec.section_type === 'reject') {
      const nonRejectIdx = (sec.options || []).findIndex((opt: any) => !opt.is_reject)
      sec.selectedOption = nonRejectIdx !== -1 ? nonRejectIdx : 0
    } else if (sec.section_type === 'positif') {
      for (const item of sec.items || []) {
        let bestIdx = 0
        let maxScore = -Infinity
        ;(sec.options || []).forEach((opt: any, oIdx: number) => {
          const optIdStr = String(opt.option_id)
          let score = 0
          if (Array.isArray(item.scores)) {
            const found = item.scores.find((entry: any) => String(entry.option_id) === optIdStr)
            score = Number(found?.score_percent) || 0
          } else if (item.scores && typeof item.scores === 'object') {
            score = parseFloat(item.scores[optIdStr]) || 0
          } else {
            score = parseFloat(opt.score) || 0
          }
          if (score > maxScore) {
            maxScore = score
            bestIdx = oIdx
          }
        })
        item.selected = bestIdx
      }
    } else if (sec.section_type === 'negatif') {
      for (const item of sec.items || []) {
        let bestIdx = 0
        let minScore = Infinity
        ;(sec.options || []).forEach((opt: any, oIdx: number) => {
          const optIdStr = String(opt.option_id)
          let score = 0
          if (Array.isArray(item.scores)) {
            const found = item.scores.find((entry: any) => String(entry.option_id) === optIdStr)
            score = Number(found?.score_percent) || 0
          } else if (item.scores && typeof item.scores === 'object') {
            score = parseFloat(item.scores[optIdStr]) || 0
          } else {
            score = parseFloat(opt.score) || 0
          }
          if (score < minScore) {
            minScore = score
            bestIdx = oIdx
          }
        })
        item.selected = bestIdx
      }
    }
  }
}

const loadElekScoring = async (collateralId: any) => {
  if (!collateralId) return
  try {
    const { data } = await useApi('/taksiran/electronic-scoring', {
      params: { collateral_id: collateralId }
    })

    if (data && Array.isArray(data.sections) && data.sections.length > 0) {
      elekScoringSections.value = data.sections.map((sec: any) => ({
        ...sec,
        selectedOption: null,
        items: (sec.items || []).map((itm: any) => ({
          ...itm,
          selected: null
        }))
      }))
      initScoringTo100()
      await calculateElekResult()
    } else {
      elekScoringSections.value = []
      await calculateElekResult()
    }
  } catch (e) {
    console.error('Failed to load electronic scoring:', e)
  }
}

const selectElekItemScore = (sIdx: number, iIdx: number, oIdx: number) => {
  elekScoringSections.value[sIdx].items[iIdx].selected = oIdx
  calculateElekResult()
}

const getItemRowScore = (section: any, item: any) => {
  if (item.selected === null || item.selected === undefined) return 0
  const opt = section.options[item.selected]
  if (!opt) return 0

  const optIdStr = String(opt.option_id)
  if (Array.isArray(item.scores)) {
    const score = item.scores.find((entry: any) => String(entry.option_id) === optIdStr)
    return Number(score?.score_percent) || 0
  }
  if (item.scores && typeof item.scores === 'object') {
    if (item.scores[optIdStr] !== undefined) {
      return parseFloat(item.scores[optIdStr]) || 0
    }
  }
  return parseFloat(opt.score) || 0
}

const getSectionTotalPositifScore = (section: any) => {
  if (!section.items || section.items.length === 0) return 0
  let sum = 0
  for (const item of section.items) {
    sum += getItemRowScore(section, item)
  }
  return sum
}

const calculateElekResult = async () => {
  let isRejected = false

  for (const sec of elekScoringSections.value) {
    if (sec.section_type === 'reject') {
      if (sec.selectedOption !== null && sec.selectedOption !== undefined) {
        const opt = sec.options[sec.selectedOption]
        if (opt && opt.is_reject) {
          isRejected = true
          break
        }
      }
    }
  }

  let totalPositif = 0
  let totalNegatif = 0

  for (const sec of elekScoringSections.value) {
    if (sec.section_type === 'positif') {
      totalPositif += getSectionTotalPositifScore(sec)
    } else if (sec.section_type === 'negatif') {
      if (sec.items) {
        for (const item of sec.items) {
          totalNegatif += getItemRowScore(sec, item)
        }
      }
    }
  }

  const hasPositif = elekScoringSections.value.some((s: any) => s.section_type === 'positif')
  let finalScore = 100
  if (hasPositif) {
    finalScore = Math.max(0, totalPositif - totalNegatif)
  } else {
    finalScore = Math.max(0, 100 - totalNegatif)
  }

  let grade = 'A'
  if (isRejected) {
    grade = 'REJECT'
  } else if (finalScore < 60) {
    grade = 'REJECT'
  } else if (finalScore < 70) {
    grade = 'D'
  } else if (finalScore < 80) {
    grade = 'C'
  } else if (finalScore < 90) {
    grade = 'B'
  } else {
    grade = 'A'
  }

  elekResult.value.score = finalScore
  elekResult.value.isRejected = isRejected || grade === 'REJECT'
  elekResult.value.grade = grade
  elekResult.value.note = elekResult.value.isRejected ? 'Barang tidak memenuhi standar kelayakan gadai (Reject).' : `Lulus pengujian dengan Grade ${grade}`

  if (elekWilayah.value && elekDetailId.value) {
    await loadElekNilaiTaksir()
  }

  elekShowResult.value = true
  emitOutput()
}

const loadElekNilaiTaksir = async () => {
  elekResult.value.hps = 0
  elekResult.value.maxLoan = 0
  elekResult.value.requestedLoan = 0
  if (!elekDetailId.value || !elekWilayah.value) return
  try {
    const { data } = await useApi('/taksiran/nilai-taksir', {
      params: {
        id_barang_detail: elekDetailId.value,
        wilayah: elekWilayah.value,
        type: 'umum'
      }
    })

    const hps = Number(data?.price) || 0
    if (hps <= 0) {
      emit('error', 'Harga taksiran barang untuk wilayah ini belum tersedia.')
      return
    }
    elekResult.value.hps = hps

    if (elekResult.value.isRejected) {
      elekResult.value.maxLoan = 0
      elekResult.value.requestedLoan = 0
    } else {
      let multiplier = 0.90
      if (elekResult.value.grade === 'B') multiplier = 0.85
      if (elekResult.value.grade === 'C') multiplier = 0.80
      if (elekResult.value.grade === 'D') multiplier = 0.75

      const maxLoan = Math.round(hps * (elekResult.value.score / 100) * multiplier)
      elekResult.value.maxLoan = maxLoan
      elekResult.value.requestedLoan = Math.round(maxLoan * 0.95)
    }
  } catch (e) {
    console.error('Failed to load HPS price:', e)
  } finally {
    emitOutput()
  }
}

const emitOutput = () => {
  emit('update:result', {
    detailName: elekDetailName.value,
    maxLoan: elekResult.value.maxLoan,
    requestedLoan: elekResult.value.requestedLoan,
    hps: elekResult.value.hps,
    score: elekResult.value.score,
    grade: elekResult.value.grade,
    isRejected: elekResult.value.isRejected,
    isReady: !elekResult.value.isRejected && elekResult.value.requestedLoan > 0
  })
}

const reset = () => {
  onElekJenisChange()
}

defineExpose({
  reset
})

const formatNum = (val: number) => {
  if (!val && val !== 0) return '0'
  return new Intl.NumberFormat('id-ID').format(Math.round(val))
}

onMounted(async () => {
  try {
    const { data: catData } = await useApi('/taksiran/jenis-barang')
    elekCategories.value = Array.isArray(catData) ? catData : []
  } catch (e) {
    console.error('Failed to load electronic categories:', e)
  }
})
</script>
