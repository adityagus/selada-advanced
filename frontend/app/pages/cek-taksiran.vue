<template>
  <div class="space-y-6 pb-12 font-sans">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl md:text-3xl font-extrabold text-slate-900 tracking-tight">Cek Taksiran</h1>
        <p class="text-slate-500 text-sm md:text-base">Simulasi presisi & riwayat taksiran Elektronik & Emas berdasarkan pembuat (sales).</p>
      </div>
      <div class="flex items-center gap-3">
        <!-- Main Mode Switcher: Simulasi vs History -->
        <div class="flex p-1.5 bg-slate-200/80 rounded-2xl gap-1">
          <button
            type="button"
            @click="mainMode = 'simulasi'"
            class="px-4 py-2 rounded-xl text-xs md:text-sm font-bold transition-all cursor-pointer flex items-center gap-2"
            :class="mainMode === 'simulasi' ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-600 hover:text-slate-900'"
          >
            <svg class="w-4 h-4 text-amber-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z"/></svg>
            Simulasi Calculator
          </button>
          <button
            type="button"
            @click="switchToHistoryMode"
            class="px-4 py-2 rounded-xl text-xs md:text-sm font-bold transition-all cursor-pointer flex items-center gap-2"
            :class="mainMode === 'history' ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-600 hover:text-slate-900'"
          >
            <svg class="w-4 h-4 text-emerald-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
            History Taksiran
          </button>
        </div>

        <button
          v-if="mainMode === 'simulasi'"
          type="button"
          @click="resetAll"
          class="px-3.5 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 text-xs md:text-sm font-semibold rounded-xl transition-all border border-slate-200"
        >
          Reset
        </button>
      </div>
    </div>

    <!-- MAIN MODE 1: SIMULASI CALCULATOR -->
    <div v-if="mainMode === 'simulasi'" class="bg-white rounded-3xl border border-slate-200/80 p-4 md:p-8 shadow-sm space-y-6">
      <!-- Tab Switcher -->
      <div class="flex p-1 bg-slate-100 rounded-2xl gap-1 max-w-md">
        <button
          type="button"
          @click="activeTab = 'elektronik'"
          class="flex-1 py-2.5 px-4 rounded-xl text-sm font-semibold transition-all cursor-pointer flex items-center justify-center gap-2"
          :class="activeTab === 'elektronik' ? 'bg-white text-amber-700 shadow-sm border border-slate-200/60 font-bold' : 'text-slate-600 hover:text-slate-900'"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/></svg>
          Elektronik
        </button>
        <button
          type="button"
          @click="activeTab = 'emas'"
          class="flex-1 py-2.5 px-4 rounded-xl text-sm font-semibold transition-all cursor-pointer flex items-center justify-center gap-2"
          :class="activeTab === 'emas' ? 'bg-white text-amber-700 shadow-sm border border-slate-200/60 font-bold' : 'text-slate-600 hover:text-slate-900'"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
          Emas
        </button>
      </div>

      <!-- =================================================== -->
      <!-- TAB: ELEKTRONIK                                    -->
      <!-- =================================================== -->
      <div v-if="activeTab === 'elektronik'" class="space-y-6">
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
                  <svg class="w-4 h-4 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
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
      </div>

      <!-- =================================================== -->
      <!-- TAB: EMAS                                          -->
      <!-- =================================================== -->
      <div v-else class="space-y-6">
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
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
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
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
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
            <input v-model="leadCustomerName" type="text" placeholder="Masukkan nama customer..." class="w-full px-3.5 py-2 bg-white border border-slate-300 rounded-xl text-sm focus:outline-none" />
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">No. Handphone / WA</label>
            <input v-model="leadCustomerPhone" type="text" placeholder="08xxxxxxxxxx" class="w-full px-3.5 py-2 bg-white border border-slate-300 rounded-xl text-sm focus:outline-none" />
          </div>
        </div>

        <div class="flex justify-end gap-3 pt-2">
          <button
            type="button"
            @click="handleSubmit"
            class="w-full sm:w-auto px-6 py-3 bg-amber-500 hover:bg-amber-600 text-white rounded-xl text-sm font-bold transition-all shadow-md shadow-amber-200 flex items-center justify-center gap-2 cursor-pointer"
            :disabled="isSubmitting"
          >
            <span v-if="isSubmitting" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
            <span v-else>Simpan Transaksi Taksiran</span>
          </button>
        </div>

        <div v-if="statusMessage" class="p-3.5 rounded-xl text-sm font-medium text-center" :class="statusClass">
          {{ statusMessage }}
        </div>
      </div>
    </div>

    <!-- MAIN MODE 2: HISTORY TAKSIRAN -->
    <div v-else-if="mainMode === 'history'" class="bg-white rounded-3xl border border-slate-200/80 p-4 md:p-8 shadow-sm space-y-6">
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
            <svg class="w-4 h-4 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/></svg>
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

    <!-- ELEKTRONIK CATALOG / SEARCH ITEM MODAL -->
    <div v-if="showElekSearchModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs z-50 flex items-center justify-center p-4">
      <div class="bg-white rounded-3xl max-w-xl w-full p-6 shadow-2xl space-y-4 max-h-[85vh] flex flex-col">
        <div class="flex items-center justify-between border-b border-slate-100 pb-3">
          <h3 class="font-bold text-slate-800 text-base md:text-lg">Pilih Barang</h3>
          <button @click="showElekSearchModal = false" class="text-slate-400 hover:text-slate-600 p-1 text-lg">✕</button>
        </div>

        <div class="flex gap-2">
          <input
            v-model="elekSearchQuery"
            type="text"
            @keyup.enter="searchElekItems"
            placeholder="Cari nama barang..."
            class="flex-1 px-3.5 py-2 bg-white border border-slate-300 rounded-xl text-sm focus:outline-none focus:border-amber-500"
          />
          <button @click="searchElekItems" class="px-4 py-2 bg-amber-500 hover:bg-amber-600 text-white font-medium rounded-xl text-sm cursor-pointer">
            Cari
          </button>
        </div>

        <div class="flex-1 overflow-y-auto space-y-2 pr-1">
          <div v-if="isSearchingElek" class="py-8 text-center text-slate-400 text-sm">
            Mencari data barang...
          </div>
          <div v-else-if="elekSearchResult.length === 0" class="py-8 text-center text-slate-400 text-sm">
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
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'

const activeTab = ref<'elektronik' | 'emas'>('elektronik')

// ==========================================
// ELEKTRONIK MODULE
// ==========================================
const elekCategories = ref<any[]>([])
const elekJenisBarang = ref('')
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
  score: 0,
  isRejected: false,
  grade: '-',
  note: '-',
  maxLoan: 0,
  requestedLoan: 0
})

const onElekJenisChange = () => {
  elekDetailId.value = ''
  elekCollateralItemId.value = ''
  elekDetailName.value = ''
  elekWilayah.value = ''
  elekWilayahList.value = []
  elekScoringSections.value = []
  elekShowResult.value = false
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
  } catch (e) {
    console.error('Failed to load regions:', e)
  }
}

const loadElekScoring = async (collateralId: any) => {
  if (!collateralId) return
  try {
    const { data } = await useApi('/taksiran/electronic-scoring', {
      params: { collateral_id: collateralId }
    })

    if (data && Array.isArray(data.sections)) {
      elekScoringSections.value = data.sections.map((sec: any) => ({
        ...sec,
        selectedOption: null,
        items: (sec.items || []).map((itm: any) => ({
          ...itm,
          selected: null
        }))
      }))
    } else {
      elekScoringSections.value = []
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
  return Math.round(sum / section.items.length)
}

const calculateElekResult = async () => {
  let isRejected = false

  for (const sec of elekScoringSections.value) {
    if (sec.section_type === 'reject') {
      if (sec.selectedOption !== null && sec.selectedOption !== undefined) {
        const opt = sec.options[sec.selectedOption]
        if (opt && opt.is_reject) {
          isRejected = true;
          break;
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

  const finalScore = Math.max(0, totalPositif - totalNegatif)

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
  elekResult.value.note = isRejected ? 'Barang tidak memenuhi standar kelayakan gadai (Reject).' : `Lulus pengujian dengan Grade ${grade}`

  if (elekWilayah.value && elekDetailId.value) {
    await loadElekNilaiTaksir()
  }

  elekShowResult.value = true
}

const loadElekNilaiTaksir = async () => {
  if (!elekDetailId.value || !elekWilayah.value) return
  try {
    const { data } = await useApi('/taksiran/nilai-taksir', {
      params: {
        detail_id: elekDetailId.value,
        region_id: elekWilayah.value,
        type: 'umum'
      }
    })

    const hps = data && data.hps ? parseFloat(data.hps) : 3000000
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
  }
}

// ==========================================
// EMAS MODULE
// ==========================================
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

const emasSavedList = ref<any[]>([])

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

    const stle = (data && data.stle) ? parseFloat(data.stle) : 1200000
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

  // Clear form
  emasDetailId.value = ''
  emasCollateralName.value = ''
  emasKarat.value = 24
  emasTotal.value = 1
  emasGrossWeight.value = 0
  emasNetWeight.value = 0
  emasNotes.value = ''
  isEmasLM.value = false
  currentEmasRowPrice.value = 0
}

const removeEmasSavedItem = (idx: number) => {
  emasSavedList.value.splice(idx, 1)
}

const totalEmasTaksiran = computed(() => {
  return emasSavedList.value.reduce((acc, item) => acc + (item.price || 0), 0)
})

// ==========================================
// SHARED SUBMIT & GENERAL ACTIONS
// ==========================================
const saveAsLead = ref(false)
const leadCustomerName = ref('')
const leadCustomerPhone = ref('')
const isSubmitting = ref(false)
const statusMessage = ref('')
const statusClass = ref('')

const resetAll = () => {
  onElekJenisChange()
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
  saveAsLead.value = false
  leadCustomerName.value = ''
  leadCustomerPhone.value = ''
  statusMessage.value = ''
}

const handleSubmit = async () => {
  isSubmitting.value = true
  statusMessage.value = ''

  const isElek = activeTab.value === 'elektronik'
  const nilaiTaksir = isElek ? elekResult.value.maxLoan : totalEmasTaksiran.value
  const nilaiPinjaman = isElek ? elekResult.value.requestedLoan : Math.round(totalEmasTaksiran.value * 0.95)

  try {
    const { error } = await useApi('/taksiran/simpan', {
      method: 'POST',
      body: {
        jenis_barang: activeTab.value,
        nama_barang: isElek ? elekDetailName.value : 'Paket Barang Emas',
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
          ket: `Simulasi Taksiran ${isElek ? elekDetailName.value : 'Emas'} - Total Taksiran Rp ${formatNum(nilaiTaksir)}`
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
    }
  } catch (e: any) {
    isSubmitting.value = false
    statusMessage.value = 'Terjadi kesalahan sistem'
    statusClass.value = 'bg-rose-50 border border-rose-200 text-rose-700 font-bold'
  }
}

const mainMode = ref<'simulasi' | 'history'>('simulasi')
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

const switchToHistoryMode = () => {
  mainMode.value = 'history'
  if (salesUsers.value.length === 0) {
    fetchSalesUsers()
  }
  fetchTaksiranHistory()
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

onMounted(async () => {
  try {
    const { data: catData } = await useApi('/taksiran/jenis-barang')
    elekCategories.value = Array.isArray(catData) ? catData : []
  } catch (e) {
    console.error('Failed to load electronic categories:', e)
  }

  loadEmasBarang()
  fetchSalesUsers()
})

const formatNum = (val: number) => {
  if (!val && val !== 0) return '0'
  return new Intl.NumberFormat('id-ID').format(Math.round(val))
}
</script>
