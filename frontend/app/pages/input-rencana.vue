<template>
  <div class="space-y-6 pb-12 font-sans">
    <!-- Header with Back Button & Quick Actions -->
    <div
      class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-200/80 pb-4"
    >
      <div class="flex items-center gap-3">
        <NuxtLink
          to="/inquiry"
          class="p-2.5 bg-white hover:bg-slate-100 border border-slate-200 rounded-xl text-slate-600 transition-all flex items-center gap-1 text-sm font-medium shadow-2xs"
        >
          <AppIcon name="arrow-left" class="w-5 h-5" />
        </NuxtLink>
        <div>
          <h1 class="text-2xl font-bold text-slate-800 tracking-tight">
            Progress Sales
          </h1>
          <p class="text-xs text-slate-500">
            Kelola status inquiry dari Leads hingga Realisasi.
          </p>
        </div>
      </div>
      <div class="flex items-center gap-2">
        <span
          class="text-xs font-semibold px-3 py-1.5 bg-amber-50 text-amber-800 border border-amber-200 rounded-xl font-mono shadow-2xs"
        >
          ID Rencana: #{{ idRencana || "Baru" }}
        </span>
      </div>
    </div>

    <!-- Loading State -->
    <div
      v-if="isLoading"
      class="bg-white p-12 rounded-2xl border border-slate-200/80 text-center shadow-xs"
    >
      <AppIcon
        name="loader"
        class="w-8 h-8 text-amber-500 animate-spin mx-auto mb-3"
      />
      <p class="text-slate-500 text-sm font-medium">
        Memuat data inquiry & riwayat customer ID #{{ idRencana }}...
      </p>
    </div>

    <template v-else>
      <!-- STEPPER WIZARD TRACKER (Leads -> Prospect -> Hot Prospect -> SBG) -->
      <div class="rounded-2xlp-6 shadow-xs space-y-6">
        <div
          class="flex flex-col sm:flex-row sm:items-center justify-end gap-3"
        >
          <div class="flex items-center gap-2">
            <span
              class="text-xs font-semibold px-3 py-1 rounded-full border"
              :class="getStatusBadgeClass(activeTab)"
            >
              Tahap Aktif: {{ activeTab }}
            </span>
          </div>
        </div>

        <!-- Progress Bar & Stepper Buttons -->
        <div class="py-2 px-2 max-w-2xl mx-auto">
          <div class="flex items-center justify-between relative">
            <div
              class="absolute left-1 right-8 top-1/2 -translate-y-1/2 h-1 bg-slate-300 -z-0 -mt-[10px]"
            ></div>
            <div
              class="absolute left-1 top-1/2 -translate-y-1/2 h-1 bg-amber-500 transition-all duration-500 -z-0 -mt-[10px]"
              :style="{ width: getStepProgressWidth(activeTab) }"
            ></div>

            <!-- Step 1: Leads -->
            <button
              type="button"
              @click="selectTab('Leads')"
              class="relative z-10 flex flex-col items-center gap-1.5 cursor-pointer transition-all wizard"
            >
              <div
                class="w-10 h-10 rounded-full flex items-center justify-center text-sm font-bold transition-all"
                :class="
                  isStepActive('Leads', activeTab)
                    ? 'bg-amber-500 text-white shadow-md shadow-amber-200 ring-4 ring-amber-100'
                    : 'bg-slate-100 text-slate-400'
                "
              >
                1
              </div>
              <span
                class="text-xs font-semibold uppercase tracking-wider"
                :class="
                  activeTab === 'Leads'
                    ? 'text-amber-700 font-bold'
                    : 'text-slate-500'
                "
                >1. Leads</span
              >
            </button>

            <!-- Step 2: Prospect -->
            <button
              type="button"
              @click="selectTab('Prospect')"
              class="relative z-10 flex flex-col items-center gap-1.5 cursor-pointer transition-all wizard"
            >
              <div
                class="w-10 h-10 rounded-full flex items-center justify-center text-sm font-bold transition-all"
                :class="
                  isStepActive('Prospect', activeTab)
                    ? 'bg-blue-600 text-white shadow-md shadow-blue-200 ring-4 ring-blue-100'
                    : 'bg-slate-100 text-slate-400'
                "
              >
                2
              </div>
              <span
                class="text-xs font-semibold uppercase tracking-wider"
                :class="
                  activeTab === 'Prospect'
                    ? 'text-blue-700 font-bold'
                    : 'text-slate-400'
                "
                >2. Prospek</span
              >
            </button>

            <!-- Step 3: Hot Prospect -->
            <button
              type="button"
              @click="selectTab('Hot Prospect')"
              class="relative z-10 flex flex-col items-center gap-1.5 cursor-pointer transition-all wizard"
            >
              <div
                class="w-10 h-10 rounded-full flex items-center justify-center text-sm font-bold transition-all"
                :class="
                  isStepActive('Hot Prospect', activeTab)
                    ? 'bg-rose-600 text-white shadow-md shadow-rose-200 ring-4 ring-rose-100'
                    : 'bg-slate-100 text-slate-400'
                "
              >
                3
              </div>
              <span
                class="text-xs font-semibold uppercase tracking-wider"
                :class="
                  activeTab === 'Hot Prospect'
                    ? 'text-rose-700 font-bold'
                    : 'text-slate-400'
                "
                >3. Hot Prospek</span
              >
            </button>

            <!-- Step 4: SBG (Realisasi) -->
            <button
              type="button"
              @click="selectTab('SBG')"
              class="relative z-10 flex flex-col items-center gap-1.5 cursor-pointer transition-all wizard"
            >
              <div
                class="w-10 h-10 rounded-full flex items-center justify-center text-sm font-bold transition-all"
                :class="
                  isStepActive('SBG', activeTab)
                    ? 'bg-emerald-600 text-white shadow-md shadow-emerald-200 ring-4 ring-emerald-100'
                    : 'bg-slate-100 text-slate-400'
                "
              >
                4
              </div>
              <span
                class="text-xs font-semibold uppercase tracking-wider"
                :class="
                  activeTab === 'SBG'
                    ? 'text-emerald-700 font-bold'
                    : 'text-slate-400'
                "
                >4. SBG</span
              >
            </button>
          </div>
        </div>

        <!-- Warning Alert Banner -->
        <div
          v-if="stepWarning"
          class="p-3.5 bg-amber-50 border border-amber-300 rounded-xl text-xs font-medium text-amber-800 flex items-center gap-2.5"
        >
          <AppIcon
            name="warning"
            class="w-4 h-4 text-amber-600 flex-shrink-0"
          />
          <span>{{ stepWarning }}</span>
        </div>
      </div>

      <!-- RIWAYAT FOLLOW-UP & LOG AKTIVITAS -->
      <div
        v-if="historyList.length > 0"
        class="bg-white rounded-2xl border border-slate-200/80 p-6 shadow-xs space-y-4"
      >
        <div
          class="flex items-center justify-between border-b border-slate-100 pb-3"
        >
          <div>
            <h2 class="font-semibold text-slate-800 text-base">
              Riwayat Aktivitas
            </h2>
          </div>
          <span
            class="text-xs font-semibold px-3 py-1 bg-amber-50 text-amber-800 border border-amber-200 rounded-full"
          >
            {{ historyList.length }} Catatan Aktivitas
          </span>
        </div>

        <div class="relative border-l-2 border-slate-200 ml-3 space-y-4 py-1">
          <div v-for="(h, idx) in historyList" :key="idx" class="relative pl-6">
            <div
              class="absolute -left-[9px] top-1.5 w-4 h-4 rounded-full bg-white border-2 border-amber-500"
            ></div>
            <div
              class="bg-slate-50 border border-slate-200/80 rounded-xl p-3 space-y-1"
            >
              <div class="flex items-center justify-between">
                <span
                  class="text-xs font-semibold px-2 py-0.5 rounded-md border"
                  :class="getStatusBadgeClass(h.status)"
                >
                  {{ h.status }}
                </span>
                <span class="text-xs text-slate-400 font-mono">{{
                  h.tgl
                }}</span>
              </div>
              <p class="text-xs text-slate-700 leading-relaxed font-medium">
                {{ h.keterangan }}
              </p>
              <span class="text-[11px] text-slate-400 block"
                >Sumber: {{ h.source }}</span
              >
            </div>
          </div>
        </div>
      </div>

      <!-- ========================================================================= -->
      <!-- PANEL 1: LEADS (DATA KONTAK AWAL) -->
      <!-- ========================================================================= -->
      <div
        v-if="activeTab === 'Leads'"
        class="bg-white rounded-2xl border border-amber-200 p-6 shadow-xs space-y-5"
      >
        <div
          class="flex items-center justify-between border-b border-slate-100 pb-3"
        >
          <div class="flex items-center gap-2">
            <span
              class="w-7 h-7 rounded-full bg-amber-500 text-white font-bold text-xs flex items-center justify-center"
              >1</span
            >
            <h2 class="font-bold text-slate-800 text-base">
              PANEL STEP 1 — LEADS (Data Kontak Awal)
            </h2>
          </div>
          <span
            class="text-xs font-medium text-amber-800 bg-amber-50 px-3 py-1 rounded-lg border border-amber-200"
            >Status: Leads</span
          >
        </div>

        <div
          class="p-3.5 bg-amber-50/60 border border-amber-200 rounded-xl text-xs text-amber-900 flex items-center gap-2"
        >
          <AppIcon name="info" class="w-4 h-4 text-amber-600 flex-shrink-0" />
          <span
            >Isi data kontak awal customer. Nomor HP dari sistem dikunci
            (readonly) demi validasi.</span
          >
        </div>

        <div class="grid grid-cols-1 md:grid-cols-1 gap-4">
          <div>
            <label
              class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1"
              >Nama Lengkap Customer <span class="text-rose-500">*</span></label
            >
            <input
              v-model="customerInfo.nama"
              type="text"
              class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm font-semibold text-slate-800"
            />
          </div>

          <div>
            <label
              class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1"
              >No. Handphone / WhatsApp (Readonly)</label
            >
            <input
              v-model="customerInfo.hp"
              type="text"
              class="w-full px-3.5 py-2.5 bg-slate-100 border border-slate-300 rounded-xl text-sm font-medium text-slate-700 cursor-not-allowed"
              readonly
            />
          </div>
        </div>

        <!-- 3 Inputan Vertical: SOA, Aktifitas, Jenis Aplikasi (Baris ke Bawah) -->
        <div class="space-y-4 pt-1">
          <!-- Input 1: SOA -->
          <div>
            <label
              class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1"
            >
              SOA <span class="text-rose-500">*</span>
            </label>
            <select
              v-model="customerInfo.soa"
              @change="onCustomerSoaChange"
              class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm font-medium text-slate-800 focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all"
              required
            >
              <option value="">- Pilih SOA -</option>
              <option v-for="soaOpt in soaList" :key="soaOpt" :value="soaOpt">
                {{ soaOpt }}
              </option>
            </select>
          </div>

          <!-- Input 2: Aktifitas -->
          <div>
            <label
              class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1"
            >
              Aktifitas <span class="text-rose-500">*</span>
            </label>
            <select
              v-model="customerInfo.aktifitas"
              :disabled="!customerInfo.soa"
              class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm font-medium text-slate-800 focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all disabled:bg-slate-100 disabled:cursor-not-allowed"
              required
            >
              <option value="">- Pilih Aktifitas -</option>
              <option
                v-for="aktOpt in availableCustomerAktifitasOptions"
                :key="aktOpt"
                :value="aktOpt"
              >
                {{ aktOpt }}
              </option>
            </select>
          </div>

          <!-- Input 3: Jenis Aplikasi -->
          <div>
            <label
              class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1"
            >
              Jenis Aplikasi <span class="text-rose-500">*</span>
            </label>
            <select
              v-model="customerInfo.jenisAplikasi"
              :disabled="!customerInfo.soa"
              class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm font-medium text-slate-800 focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all disabled:bg-slate-100 disabled:cursor-not-allowed"
              required
            >
              <option value="">- Pilih Jenis Aplikasi -</option>
              <option
                v-for="jenisOpt in availableCustomerJenisAplikasiOptions"
                :key="jenisOpt"
                :value="jenisOpt"
              >
                {{ jenisOpt }}
              </option>
            </select>
          </div>
        </div>

        <!-- <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1">Tanggal Follow Up
              Perdana</label>
            <input v-model="form.tglFollowUp" type="date"
              class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-800" />
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1">Jam Follow
              Up</label>
            <input v-model="form.jamFollowUp" type="time"
              class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-800" />
          </div>
        </div> -->

        <div
          class="flex flex-wrap items-center justify-end gap-3 pt-3 border-t border-slate-100"
        >
          <div class="flex items-center gap-2">
            <button
              @click="submitWithStatus('Leads', false)"
              type="button"
              class="px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-xl text-xs font-semibold transition-all cursor-pointer"
              :disabled="isSubmitting"
            >
              Simpan Data Leads
            </button>
            <button
              @click="submitWithStatus('Leads', true)"
              type="button"
              class="px-5 py-2 bg-amber-500 hover:bg-amber-600 text-white rounded-xl text-xs font-bold shadow-xs flex items-center gap-1.5 cursor-pointer transition-all"
              :disabled="isSubmitting"
            >
              <AppIcon
                name="loader"
                v-if="isSubmitting"
                class="w-3.5 h-3.5 animate-spin"
              />
              <span v-else>Lanjut ke Prospek</span>
              <AppIcon
                name="arrow-right"
                v-if="!isSubmitting"
                class="w-4 h-4"
              />
            </button>
          </div>
        </div>
      </div>

      <!-- ========================================================================= -->
      <!-- PANEL 2: PROSPEK (DATA ALAMAT & SERVE CONSULTATION)                       -->
      <!-- ========================================================================= -->
      <div
        v-else-if="activeTab === 'Prospect'"
        class="bg-white rounded-2xl border border-blue-200 p-6 shadow-xs space-y-5"
      >
        <div
          class="flex items-center justify-between border-b border-slate-100 pb-3"
        >
          <div class="flex items-center gap-2">
            <span
              class="w-7 h-7 rounded-full bg-blue-600 text-white font-bold text-xs flex items-center justify-center"
              >2</span
            >
            <h2 class="font-bold text-slate-800 text-base">
              PANEL STEP 2 — PROSPEK (Alamat & Konsultasi Taksiran)
            </h2>
          </div>
          <span
            class="text-xs font-medium text-blue-800 bg-blue-50 px-3 py-1 rounded-lg border border-blue-200"
            >Status: Prospek</span
          >
        </div>

        <div
          class="p-3.5 bg-blue-50/60 border border-blue-200 rounded-xl text-xs text-blue-900 flex items-center gap-2"
        >
          <AppIcon name="pin" class="w-4 h-4 text-blue-600 flex-shrink-0" />
          <span
            >Lengkapi alamat domisili customer dan catatan hasil konsultasi awal
            taksiran barang.</span
          >
        </div>

        <!-- Alamat Form -->
        <div class="space-y-4">
          <div>
            <label
              class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1"
              >Jalan / Alamat Lengkap
              <span class="text-rose-500">*</span></label
            >
            <input
              v-model="customerInfo.alamat"
              type="text"
              placeholder="Jalan, Perumahan, No. Rumah"
              class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm"
            />
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-medium text-slate-500 mb-1"
                >RT</label
              >
              <input
                v-model="customerInfo.rt"
                type="text"
                placeholder="001"
                class="w-full px-3 py-2 bg-white border border-slate-300 rounded-xl text-sm"
              />
            </div>
            <div>
              <label class="block text-xs font-medium text-slate-500 mb-1"
                >RW</label
              >
              <input
                v-model="customerInfo.rw"
                type="text"
                placeholder="002"
                class="w-full px-3 py-2 bg-white border border-slate-300 rounded-xl text-sm"
              />
            </div>
          </div>

          <div class="grid grid-cols-2 md:grid-cols-3 gap-3">
            <!-- Search Select Kelurahan -->
            <div class="relative col-span-2 md:col-span-1">
              <label
                class="text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1 flex items-center justify-between"
              >
                <span>Kelurahan</span>
                <span class="text-amber-600 text-[10px] lowercase font-normal"
                  >Cari Master</span
                >
              </label>
              <div class="relative">
                <input
                  :value="customerInfo.kelurahan"
                  @input="onKelurahanInput"
                  @focus="onKelurahanFocus"
                  @blur="onKelurahanBlur"
                  type="text"
                  placeholder="Ketik kelurahan..."
                  class="w-full px-3 py-2 bg-white border border-slate-300 rounded-xl text-sm font-medium focus:ring-2 focus:ring-amber-200 focus:border-amber-500 outline-none"
                />
                <AppIcon
                  name="loader"
                  v-if="isSearchingKelurahan"
                  class="absolute right-2.5 top-2.5 w-3.5 h-3.5 text-amber-500 animate-spin"
                />
              </div>

              <!-- Search Results Dropdown -->
              <div
                v-if="
                  showKelurahanDropdown &&
                  kelurahanList &&
                  kelurahanList.length > 0
                "
                class="absolute left-0 right-0 top-full mt-1 bg-white border border-slate-200 rounded-xl shadow-lg z-50 max-h-48 overflow-y-auto divide-y divide-slate-100"
              >
                <div
                  v-for="(k, idx) in kelurahanList"
                  :key="idx"
                  @mousedown.prevent="onSelectKelurahan(k)"
                  class="p-2.5 hover:bg-amber-50 cursor-pointer transition-all text-xs"
                >
                  <div class="font-bold text-slate-800">
                    {{ k.nm_kelurahan }}
                  </div>
                  <div
                    class="text-slate-500 text-[11px] flex items-center gap-1"
                  >
                    <span>Kec. {{ k.nm_kecamatan }}</span>
                    <span>•</span>
                    <span class="text-amber-700 font-medium">{{
                      k.nm_kota
                    }}</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Auto-filled Kecamatan -->
            <div>
              <label
                class="block text-xs font-medium text-slate-500 mb-1 flex items-center justify-between"
              >
                <span>Kecamatan</span>
                <span
                  v-if="customerInfo.kecamatan"
                  class="text-[10px] text-emerald-600 font-semibold"
                  >Auto</span
                >
              </label>
              <input
                v-model="customerInfo.kecamatan"
                type="text"
                placeholder="Kecamatan"
                disabled
                class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-xl text-sm font-medium text-slate-800"
              />
            </div>

            <!-- Auto-filled Kota -->
            <div class="col-span-2 md:col-span-1">
              <label
                class="text-xs font-medium text-slate-500 mb-1 flex items-center justify-between"
              >
                <span>Kota / Kab</span>
                <span
                  v-if="customerInfo.kd_kota || customerInfo.kota"
                  class="text-[10px] text-emerald-600 font-semibold"
                  >Kode: {{ customerInfo.kd_kota || customerInfo.kota }}</span
                >
              </label>
              <input
                :value="customerInfo.nm_kota ? `${customerInfo.nm_kota} (${customerInfo.kd_kota || customerInfo.kota})` : (customerInfo.kd_kota || customerInfo.kota)"
                type="text"
                placeholder="Kota / Kab"
                disabled
                class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-xl text-sm font-medium text-slate-800"
              />
            </div>
          </div>
        </div>

        <div>
          <label
            class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1"
            >Keterangan Aktivitas Prospek</label
          >
          <textarea
            v-model="customerInfo.keterangan"
            rows="2"
            placeholder="Hasil konsultasi taksiran atau respon customer..."
            class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm"
          ></textarea>
        </div>

        <!-- Step Warning in Panel 2 -->
        <div
          v-if="stepWarning"
          class="p-3 bg-rose-50 border border-rose-200 text-rose-700 rounded-xl text-xs flex items-center gap-2"
        >
          <AppIcon name="info" class="w-4 h-4 flex-shrink-0 text-rose-500" />
          <span>{{ stepWarning }}</span>
        </div>

        <div
          class="flex flex-wrap items-center justify-between gap-3 pt-3 border-t border-slate-100"
        >
          <button
            @click="prevStep"
            type="button"
            class="px-4 py-2 border border-slate-300 text-slate-700 hover:bg-slate-50 rounded-xl text-xs font-semibold transition-all flex items-center gap-1.5 cursor-pointer"
          >
            <AppIcon name="arrow-left" class="w-4 h-4" />
            <span>Kembali</span>
          </button>

          <div class="flex flex-wrap items-center gap-2">
            <!-- Save progress and redirect to Inquiry page -->
            <button
              @click="submitWithStatus('Prospect', false, true)"
              type="button"
              class="px-4 py-2 bg-blue-50 hover:bg-blue-100 active:scale-95 text-blue-800 border border-blue-200 rounded-xl text-xs font-bold transition-all cursor-pointer flex items-center gap-1.5"
              :disabled="isSubmitting"
            >
              <AppIcon
                name="loader"
                v-if="isSubmitting"
                class="w-3.5 h-3.5 animate-spin"
              />
              <span>Simpan Progress</span>
            </button>

            <!-- Save Prospect and transition to Hot Prospect panel -->
            <button
              @click="submitWithStatus('Prospect', true)"
              type="button"
              class="px-5 py-2 bg-rose-600 hover:bg-rose-700 text-white rounded-xl text-xs font-bold flex items-center gap-1.5 shadow-sm cursor-pointer transition-all"
              :disabled="isSubmitting"
            >
              <AppIcon
                name="loader"
                v-if="isSubmitting"
                class="w-3.5 h-3.5 animate-spin"
              />
              <span v-else>Lanjut ke Hot Prospek</span>
              <AppIcon
                name="arrow-right"
                v-if="!isSubmitting"
                class="w-4 h-4"
              />
            </button>
          </div>
        </div>
      </div>

      <!-- ========================================================================= -->
      <!-- PANEL 3: HOT PROSPEK (Detail Transaksi, Produk & Barang Jaminan DB)      -->
      <!-- ========================================================================= -->
      <div
        v-else-if="activeTab === 'Hot Prospect'"
        class="bg-white rounded-2xl border border-rose-200 p-6 shadow-xs space-y-6"
      >
        <div
          class="flex items-center justify-between border-b border-slate-100 pb-3"
        >
          <div class="flex items-center gap-2">
            <span
              class="w-7 h-7 rounded-full bg-rose-600 text-white font-bold text-xs flex items-center justify-center"
              >3</span
            >
            <h2 class="font-bold text-slate-800 text-base">
              PANEL STEP 3 — HOT PROSPEK (Detail Transaksi & Barang Jaminan DB)
            </h2>
          </div>
          <span
            class="text-xs font-bold text-rose-800 bg-rose-50 px-3 py-1 rounded-lg border border-rose-200"
          >
            Status DB: {{ dbStatus }}
          </span>
        </div>

        <div
          class="p-3.5 bg-rose-50/60 border border-rose-200 rounded-xl text-xs text-rose-900 flex items-center gap-2"
        >
          <AppIcon name="flame" class="w-4 h-4 text-rose-600 flex-shrink-0" />
          <span
            >Langkah final sebelum realisasi SBG. Data produk & barang jaminan
            terhubung langsung ke Master Database Selada.</span
          >
        </div>

        <!-- 1. Pilihan Kartu Jenis Gadai (Gadai vs Cicilan) -->
        <div class="space-y-2">
          <label
            class="block text-xs font-semibold text-slate-700 uppercase tracking-wider"
          >
            Jenis Gadai <span class="text-rose-500">*</span>
          </label>
          <div class="grid grid-cols-2 gap-4 max-w-md">
            <div
              @click="selectJenisGadai('Gadai')"
              class="p-4 rounded-xl border-2 text-center cursor-pointer transition-all"
              :class="
                sbgData.jenisTransaksi === 'Gadai'
                  ? 'border-amber-500 bg-amber-50/60 shadow-xs'
                  : 'border-slate-200 hover:border-slate-300 bg-white'
              "
            >
              <div class="text-2xl mb-1">
                <AppIcon
                  name="diamond"
                  class="inline-block w-5 h-5 align-middle"
                />
              </div>
              <div class="font-bold text-sm text-slate-800">Gadai</div>
              <div class="text-xs text-slate-500">Gadai Barang</div>
            </div>

            <div
              @click="selectJenisGadai('Cicilan')"
              class="p-4 rounded-xl border-2 text-center cursor-pointer transition-all"
              :class="
                sbgData.jenisTransaksi === 'Cicilan'
                  ? 'border-amber-500 bg-amber-50/60 shadow-xs'
                  : 'border-slate-200 hover:border-slate-300 bg-white'
              "
            >
              <div class="text-2xl mb-1">
                <AppIcon
                  name="card"
                  class="inline-block w-5 h-5 align-middle"
                />
              </div>
              <div class="font-bold text-sm text-slate-800">Cicilan</div>
              <div class="text-xs text-slate-500">Cicil Barang</div>
            </div>
          </div>
        </div>

        <!-- 2. Nama Produk dari Database (product.product_offerings) -->
        <div
          class="grid grid-cols-1 md:grid-cols-2 gap-4 pt-2 border-t border-slate-100"
        >
          <div>
            <label
              class="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1"
            >
              Nama Produk <span class="text-rose-500">*</span>
              <span class="text-[10px] text-amber-600 lowercase font-normal"
                >(dari Master DB)</span
              >
            </label>
            <select
              v-model="selectedProdukCode"
              @change="onProdukChange"
              class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm font-medium"
            >
              <option value="">-- Pilih Produk --</option>
              <option
                v-for="p in productOfferings"
                :key="p.kd_produk"
                :value="p.kd_produk"
              >
                {{ p.nm_produk }}
              </option>
            </select>
          </div>
          <div>
            <label
              class="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1"
            >
              Jenis Barang Jaminan <span class="text-rose-500">*</span>
            </label>
            <div class="grid grid-cols-2 gap-2">
              <button
                type="button"
                @click="toggleJenisBarang('Emas')"
                class="py-2.5 px-4 rounded-xl text-xs font-bold border transition-all flex items-center justify-center gap-1.5 cursor-pointer"
                :class="
                  barangInfo.jenisBarang === 'Emas'
                    ? 'bg-amber-500 text-white border-amber-600 shadow-xs'
                    : 'bg-white text-slate-600 border-slate-200 hover:bg-slate-50'
                "
              >
                <span
                  ><AppIcon
                    name="ring"
                    class="inline-block w-5 h-5 align-middle"
                /></span>
                Emas
              </button>
              <button
                type="button"
                @click="toggleJenisBarang('Elektronik')"
                class="py-2.5 px-4 rounded-xl text-xs font-bold border transition-all flex items-center justify-center gap-1.5 cursor-pointer"
                :class="
                  barangInfo.jenisBarang === 'Elektronik'
                    ? 'bg-blue-600 text-white border-blue-700 shadow-xs'
                    : 'bg-white text-slate-600 border-slate-200 hover:bg-slate-50'
                "
              >
                <span
                  ><AppIcon
                    name="package"
                    class="inline-block w-5 h-5 align-middle"
                /></span>
                Non Emas (Elektronik)
              </button>
            </div>
          </div>
        </div>

        <!-- 3. CARI & PILIH BARANG JAMINAN DARI DATABASE -->
        <div class="space-y-4 pt-2 border-t border-slate-100">
          <div class="flex items-center justify-between">
            <label
              class="block text-xs font-semibold text-slate-700 uppercase tracking-wider"
            >
              Cari & Pilih Barang Jaminan di Database
              <span class="text-rose-500">*</span>
            </label>
            <span class="text-xs font-semibold text-slate-500">
              {{
                barangInfo.jenisBarang === "Emas"
                  ? "Master Emas DB"
                  : "Master Elektronik DB"
              }}
            </span>
          </div>

          <!-- IF EMAS: Multiselect dari /api/taksiran/barang-emas -->
          <div v-if="barangInfo.jenisBarang === 'Emas'" class="space-y-3">
            <label class="block text-xs font-medium text-slate-600">
              Pilih Barang Emas dari Database
              <span class="text-rose-500">*</span>
              <span class="text-[11px] text-amber-600 font-normal ml-1"
                >(Bisa pilih lebih dari 1 barang)</span
              >
            </label>

            <!-- Chips barang emas terpilih -->
            <div
              v-if="selectedEmasItems.length > 0"
              class="flex flex-wrap items-center gap-2 p-3 bg-amber-50/50 border border-amber-200 rounded-xl"
            >
              <span
                v-for="(it, idx) in selectedEmasItems"
                :key="it.id_barang_detail || idx"
                class="inline-flex items-center gap-2 px-3 py-1.5 bg-white border border-amber-300 text-amber-900 text-xs font-semibold rounded-lg"
              >
                {{ it.nama_barang }}
                <button
                  type="button"
                  @click="removeEmasItem(idx)"
                  class="hover:text-rose-600 text-slate-400 cursor-pointer"
                >
                  <AppIcon name="close" class="w-3.5 h-3.5" />
                </button>
              </span>
            </div>

            <!-- Daftar barang emas (checkbox) -->
            <div
              class="max-h-56 overflow-y-auto border border-slate-200 rounded-xl divide-y divide-slate-100"
            >
              <div
                v-if="emasDbList.length === 0"
                class="py-6 text-center text-slate-400 text-xs"
              >
                Data barang emas belum tersedia.
              </div>
              <label
                v-for="b in emasDbList"
                :key="b.id_barang_detail"
                class="flex items-center gap-3 px-3.5 py-2.5 cursor-pointer transition-all text-sm"
                :class="
                  isEmasSelected(b)
                    ? 'bg-amber-50 text-amber-900'
                    : 'hover:bg-slate-50 text-slate-700'
                "
              >
                <input
                  type="checkbox"
                  :checked="isEmasSelected(b)"
                  @change="toggleEmasItem(b)"
                  class="w-4 h-4 rounded text-amber-500 focus:ring-amber-400 cursor-pointer"
                />
                <span class="font-medium">{{ b.nama_barang }}</span>
                <span
                  v-if="b.status_antam"
                  class="ml-auto text-[11px] text-slate-400"
                  >{{ b.status_antam }}</span
                >
              </label>
            </div>
          </div>

          <!-- IF ELEKTRONIK: Modal / Search dari /api/taksiran/barang-umum (Hanya 1 barang) -->
          <div v-else class="space-y-3">
            <div class="flex items-center justify-between">
              <label class="block text-xs font-medium text-slate-600">
                Pilih Barang Elektronik dari Master DB
                <span class="text-rose-500">*</span>
                <span class="text-[11px] text-blue-600 font-normal ml-1"
                  >(Hanya 1 barang)</span
                >
              </label>
              <button
                type="button"
                @click="openElekDbSearchModal"
                class="px-3.5 py-1.5 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-bold transition-all shadow-xs cursor-pointer flex items-center gap-1.5"
              >
                <AppIcon name="search" class="w-3.5 h-3.5" />
                <span>{{
                  selectedElekItems.length > 0 ? "Ganti Barang" : "Cari Barang DB"
                }}</span>
              </button>
            </div>

            <!-- Chips list of selected electronic items -->
            <div
              v-if="selectedElekItems.length > 0"
              class="flex flex-wrap items-center gap-2 p-3 bg-slate-50 border border-slate-200 rounded-xl min-h-[52px]"
            >
              <span
                v-for="(it, idx) in selectedElekItems"
                :key="idx"
                class="inline-flex items-center gap-2 px-3 py-1.5 bg-white border border-blue-200 text-blue-900 text-xs font-semibold rounded-lg shadow-2xs"
              >
                <span
                  ><AppIcon
                    name="package"
                    class="inline-block w-5 h-5 align-middle"
                  />
                  {{ it.nm_barang || it.name }}</span
                >
                <span
                  v-if="it.kd_barang"
                  class="text-[10px] text-slate-400 font-mono"
                  >({{ it.kd_barang }})</span
                >
                <button
                  type="button"
                  @click="removeElekItem(idx)"
                  class="hover:text-rose-600 text-slate-400 font-bold ml-0.5 cursor-pointer"
                >
                  <AppIcon name="close" class="w-4 h-4" />
                </button>
              </span>
            </div>

            <div
              v-else
              @click="openElekDbSearchModal"
              class="p-4 border-2 border-dashed border-slate-200 hover:border-blue-300 rounded-xl text-center text-slate-400 text-xs cursor-pointer hover:bg-blue-50/40 transition-all flex items-center justify-center gap-2"
            >
              <AppIcon name="plus" class="w-4 h-4 text-blue-500" />
              <span
                >Belum ada barang elektronik dipilih. Klik untuk mencari &
                memilih 1 barang dari Master DB</span
              >
            </div>
          </div>
        </div>

        <!-- 4. FOTO BARANG (Ambil Foto / Upload) -->
        <div class="space-y-3 pt-2 border-t border-slate-100">
          <label
            class="block text-xs font-semibold text-slate-700 uppercase tracking-wider"
          >
            Foto Barang Jaminan <span class="text-rose-500">*</span>
          </label>
          <div class="flex flex-wrap items-center gap-4">
            <div
              class="w-24 h-24 rounded-2xl border-2 border-dashed border-slate-300 flex items-center justify-center overflow-hidden bg-slate-50"
            >
              <img
                v-if="barangFotoPreview"
                :src="barangFotoPreview"
                alt="Foto Barang"
                class="w-full h-full object-cover"
              />
              <div v-else class="text-center text-slate-400 p-2">
                <AppIcon name="camera" class="w-6 h-6 mx-auto mb-1" />
                <span class="text-[10px] block leading-tight"
                  >Belum Ada Foto</span
                >
              </div>
            </div>
              <div class="flex items-center gap-2">
                <label
                  class="inline-flex items-center gap-2 px-4 py-2 bg-slate-800 hover:bg-slate-900 text-white rounded-xl text-xs font-semibold cursor-pointer shadow-xs transition-all"
                >
                  <AppIcon name="camera" class="w-4 h-4" />
                  {{ barangFotoPreview ? 'Ganti Foto' : 'Ambil / Upload Foto Barang' }}
                  <input
                    type="file"
                    accept="image/*"
                    capture="environment"
                    @change="onFotoBarangChange"
                    class="hidden"
                  />
                </label>
                <button
                  v-if="barangFotoPreview"
                  type="button"
                  @click="removeBarangFoto"
                  class="px-3 py-2 border border-rose-200 text-rose-600 hover:bg-rose-50 rounded-xl text-xs font-semibold transition-all cursor-pointer inline-flex items-center gap-1.5"
                  title="Hapus Foto"
                >
                  <AppIcon name="trash" class="w-3.5 h-3.5" />
                  <span>Hapus</span>
                </button>
              </div>
              <span class="block text-[11px] text-slate-400 mt-1"
                >Format: JPG/PNG/WebP, otomatis dioptimasi</span
              >
            </div>
          </div>

        <!-- 5. KETERANGAN BARANG -->
        <div class="space-y-4 pt-2 border-t border-slate-100">
          <div>
            <label
              class="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1"
            >
              Keterangan Barang <span class="text-rose-500">*</span>
            </label>
            <input
              v-model="customerInfo.keteranganBarang"
              type="text"
              placeholder="Masukkan keterangan detail barang jaminan..."
              class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm"
            />
          </div>
        </div>

        <div
          class="flex flex-wrap items-center justify-between gap-3 pt-3 border-t border-slate-100"
        >
          <button
            @click="prevStep"
            type="button"
            class="px-4 py-2.5 border border-slate-300 text-slate-700 hover:bg-slate-50 rounded-xl text-xs font-semibold transition-all flex items-center gap-1.5 cursor-pointer"
          >
            <AppIcon name="arrow-left" class="w-4 h-4" />
            <span>Kembali</span>
          </button>

          <div class="flex items-center gap-2">
            <button
              @click="submitWithStatus('Prospect', false)"
              type="button"
              class="px-4 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-semibold transition-all cursor-pointer shadow-xs"
              :disabled="isSubmitting"
            >
              Simpan Data Prospek
            </button>

            <button
              @click="promptSubmitHotProspect"
              type="button"
              class="px-6 py-2.5 bg-rose-600 hover:bg-rose-700 text-white rounded-xl text-xs font-bold shadow-md shadow-rose-200 flex items-center gap-2 cursor-pointer transition-all"
              :disabled="isSubmitting"
            >
              <AppIcon
                name="loader"
                v-if="isSubmitting"
                class="w-3.5 h-3.5 animate-spin"
              />
              <span v-else>Submit Hot Prospect</span>
              <AppIcon name="check" v-if="!isSubmitting" class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>

      <!-- ========================================================================= -->
      <!-- PANEL 4: SBG / DEAL (FORM REALISASI AKAD & TRANSAKSI) -->
      <!-- ========================================================================= -->
      <div
        v-else-if="activeTab === 'SBG'"
        class="bg-white rounded-2xl border border-emerald-300 p-6 shadow-xs space-y-5"
      >
        <div
          class="flex items-center justify-between border-b border-slate-100 pb-3"
        >
          <div class="flex items-center gap-2">
            <span
              class="w-7 h-7 rounded-full bg-emerald-600 text-white font-bold text-xs flex items-center justify-center"
              >4</span
            >
            <h2 class="font-bold text-slate-800 text-base">
              PANEL STEP 4 — SBG (Realisasi Transaksi Deal SBG)
            </h2>
          </div>
          <span
            class="text-xs font-bold text-emerald-800 bg-emerald-100 px-3 py-1 rounded-lg border border-emerald-300"
            >Status: SBG Deal</span
          >
        </div>

        <div
          class="p-3.5 bg-emerald-50 border border-emerald-200 rounded-xl text-xs text-emerald-900 flex items-center gap-2"
        >
          <AppIcon
            name="check-circle"
            class="w-4 h-4 text-emerald-600 flex-shrink-0"
          />
          <span
            >Form pencairan Surat Bukti Gadai (SBG). Pastikan nomor SBG &
            tanggal jatuh tempo terverifikasi.</span
          >
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label
              class="block text-xs font-semibold text-emerald-900 uppercase tracking-wider mb-1"
              >Nomor SBG Transaksi <span class="text-rose-500">*</span></label
            >
            <input
              v-model="sbgData.noSbg"
              type="text"
              placeholder="SBG-202609-001"
              class="w-full px-3.5 py-2.5 bg-white border border-emerald-300 rounded-xl text-sm font-mono font-bold text-emerald-900"
              required
            />
          </div>

          <div>
            <label
              class="block text-xs font-semibold text-emerald-900 uppercase tracking-wider mb-1"
              >Skema Transaksi</label
            >
            <input
              v-model="sbgData.jenisTransaksi"
              type="text"
              class="w-full px-3.5 py-2.5 bg-slate-50 border border-emerald-300 rounded-xl text-sm font-medium text-emerald-900"
              readonly
            />
          </div>

          <div>
            <label
              class="block text-xs font-semibold text-emerald-900 uppercase tracking-wider mb-1"
              >Jenis Produk Gadai</label
            >
            <select
              v-model="sbgData.jenisBarang"
              class="w-full px-3.5 py-2.5 bg-white border border-emerald-300 rounded-xl text-sm font-medium text-emerald-900"
            >
              <option value="Emas">Emas / Perhiasan</option>
              <option value="Kendaraan">Kendaraan (Motor/Mobil)</option>
              <option value="Elektronik">Elektronik (HP/Laptop)</option>
            </select>
          </div>
        </div>

        <div>
          <label
            class="block text-xs font-semibold text-emerald-900 uppercase tracking-wider mb-1"
            >Deskripsi Barang Gadai SBG</label
          >
          <input
            v-model="sbgData.deskripsiBarang"
            type="text"
            placeholder="Deskripsi barang jaminan akad..."
            class="w-full px-3.5 py-2.5 bg-white border border-emerald-300 rounded-xl text-sm"
          />
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label
              class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1"
              >Tanggal Pencairan SBG</label
            >
            <input
              v-model="sbgData.tglPencairan"
              type="date"
              class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm"
            />
          </div>

          <div>
            <label
              class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1"
              >Tanggal Jatuh Tempo</label
            >
            <input
              v-model="sbgData.tglJatuhTempo"
              type="date"
              class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm"
            />
          </div>
        </div>

        <div>
          <label
            class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1"
            >Catatan Realisasi SBG / Akad</label
          >
          <textarea
            v-model="sbgData.catatanSbg"
            rows="2"
            placeholder="Catatan transaksi deal SBG..."
            class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm"
          ></textarea>
        </div>

        <div
          class="flex flex-wrap items-center justify-between gap-3 pt-3 border-t border-slate-100"
        >
          <button
            @click="prevStep"
            type="button"
            class="px-4 py-2 border border-slate-300 text-slate-700 hover:bg-slate-50 rounded-xl text-xs font-semibold transition-all flex items-center gap-1.5 cursor-pointer"
          >
            <AppIcon name="arrow-left" class="w-4 h-4" />
            <span>Kembali</span>
          </button>

          <button
            @click="submitWithStatus('SBG', true)"
            type="button"
            class="px-6 py-2.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-xs font-bold shadow-md shadow-emerald-100 flex items-center gap-1.5 cursor-pointer transition-all"
            :disabled="isSubmitting"
          >
            <AppIcon
              name="loader"
              v-if="isSubmitting"
              class="w-3.5 h-3.5 animate-spin"
            />
            <span v-else>Simpan</span>
          </button>
        </div>
      </div>

      <!-- MODAL SEARCH BARANG ELEKTRONIK DARI DATABASE (MULTISELECT) -->
      <ClientOnly>
        <Teleport to="body">
          <div
            v-if="showElekDbModal"
            class="fixed inset-0 w-screen h-screen min-h-screen top-0 left-0 right-0 bottom-0 z-[9999] bg-slate-900/60 backdrop-blur-md flex items-center justify-center p-4 !mt-0 overflow-y-auto"
          >
            <div
              class="bg-white rounded-2xl max-w-xl w-full p-6 shadow-2xl space-y-4 max-h-[85vh] flex flex-col my-auto"
            >
              <div
                class="flex items-center justify-between border-b border-slate-100 pb-3"
              >
                <div>
                  <h3 class="font-bold text-slate-800 text-base">
                    Cari & Pilih Barang Elektronik
                  </h3>
                  <p class="text-xs text-slate-400">
                    Pilih 1 barang dari Master Database
                  </p>
                </div>
                <button
                  @click="showElekDbModal = false"
                  class="text-slate-400 hover:text-slate-600 p-1 text-lg cursor-pointer"
                >
                  <AppIcon name="close" class="w-4 h-4" />
                </button>
              </div>

              <div class="flex gap-2">
                <input
                  v-model="elekDbSearch"
                  @keyup.enter="searchElekDb"
                  type="text"
                  placeholder="Ketik merk / tipe HP, Laptop, TV, dsb..."
                  class="flex-1 px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm font-medium outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-200"
                />
                <button
                  @click="searchElekDb"
                  class="px-5 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-bold cursor-pointer"
                >
                  Cari
                </button>
              </div>

              <div
                class="flex-1 overflow-y-auto space-y-2 pr-1 divide-y divide-slate-100"
              >
                <div
                  v-if="isSearchingElekDb"
                  class="py-8 text-center text-slate-400 text-xs"
                >
                  <AppIcon
                    name="loader"
                    class="w-5 h-5 text-blue-600 animate-spin mx-auto mb-2"
                  />
                  Mencari barang di database...
                </div>
                <div
                  v-else-if="elekDbList.length === 0"
                  class="py-8 text-center text-slate-400 text-xs"
                >
                  Ketik pencarian barang dan tekan tombol Cari.
                </div>
                <div
                  v-for="item in elekDbList"
                  :key="item.id_barang_detail || item.id_detail"
                  @click="toggleElekItem(item)"
                  class="p-3 rounded-xl cursor-pointer transition-all flex items-center justify-between gap-3 border"
                  :class="
                    isElekSelected(item)
                      ? 'bg-blue-50/80 border-blue-300 text-blue-900'
                      : 'hover:bg-slate-50 border-transparent text-slate-800'
                  "
                >
                  <div class="flex items-center gap-3">
                    <input
                      type="radio"
                      name="elek-db-item"
                      :checked="isElekSelected(item)"
                      @click.stop="toggleElekItem(item)"
                      class="w-4 h-4 text-blue-600 focus:ring-blue-500 cursor-pointer"
                    />
                    <div>
                      <span class="font-semibold text-sm block">{{
                        item.nm_barang
                      }}</span>
                      <span class="text-xs text-slate-400 font-mono">{{
                        item.kd_barang
                      }}</span>
                    </div>
                  </div>
                  <div class="text-right">
                    <span class="text-xs text-slate-400 block"
                      >Taksiran Standar</span
                    >
                    <span class="text-xs font-bold text-emerald-700"
                      >Rp
                      {{
                        Number(
                          item.harga_taksir || item.harga || 0,
                        ).toLocaleString("id-ID")
                      }}</span
                    >
                  </div>
                </div>
              </div>

              <div
                class="flex items-center justify-between pt-3 border-t border-slate-100"
              >
                <span class="text-xs text-slate-500 font-medium truncate">
                  <template v-if="selectedElekItems.length > 0">
                    Terpilih:
                    <strong class="text-blue-600">{{
                      selectedElekItems[0].nm_barang ||
                      selectedElekItems[0].name
                    }}</strong>
                  </template>
                  <template v-else>Belum ada barang terpilih</template>
                </span>
                <div class="flex items-center gap-2">
                  <button
                    @click="showElekDbModal = false"
                    class="px-4 py-2 border border-slate-300 text-slate-600 rounded-xl text-xs font-semibold hover:bg-slate-50 cursor-pointer"
                  >
                    Batal
                  </button>
                  <button
                    @click="showElekDbModal = false"
                    class="px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-bold cursor-pointer"
                  >
                    Selesai
                  </button>
                </div>
              </div>
            </div>
          </div>
        </Teleport>
      </ClientOnly>
      <!-- MODAL KONFIRMASI SUBMIT HOT PROSPECT (TIDAK ADA PENGEMBALIAN) -->
      <ClientOnly>
        <Teleport to="body">
          <div
            v-if="showHotProspectConfirmModal"
            class="fixed inset-0 w-screen h-screen min-h-screen top-0 left-0 right-0 bottom-0 z-[9999] bg-slate-900/60 backdrop-blur-md flex items-center justify-center p-4 !mt-0 overflow-y-auto"
          >
            <div
              class="bg-white rounded-2xl max-w-md w-full p-6 shadow-2xl space-y-4 my-auto"
            >
              <div
                class="flex items-center gap-3 text-rose-600 pb-2 border-b border-rose-100"
              >
                <div
                  class="w-10 h-10 rounded-full bg-rose-100 flex items-center justify-center flex-shrink-0"
                >
                  <AppIcon name="warning" class="w-6 h-6 text-rose-600" />
                </div>
                <div>
                  <h3 class="font-bold text-slate-900 text-base">
                    Konfirmasi Submit Hot Prospect
                  </h3>
                  <p class="text-xs text-rose-600 font-semibold">
                    Tindakan Ini Tidak Dapat Dibatalkan
                  </p>
                </div>
              </div>

              <div
                class="p-3.5 bg-rose-50 border border-rose-200 rounded-xl text-xs text-rose-900 leading-relaxed space-y-2"
              >
                <p class="font-bold flex items-center gap-1.5 text-rose-800">
                  <AppIcon name="warning" class="w-4 h-4" /> Peringatan
                  Penguncian Status:
                </p>
                <p>
                  Setelah di-submit sebagai <strong>Hot Prospect</strong>,
                  status data nasabah ini akan dikunci dan
                  <strong>TIDAK DAPAT DIKEMBALIKAN</strong> kembali ke status
                  Leads atau Prospek.
                </p>
              </div>

              <div
                class="flex items-center justify-end gap-3 pt-2 border-t border-slate-100"
              >
                <button
                  type="button"
                  @click="showHotProspectConfirmModal = false"
                  class="px-4 py-2.5 border border-slate-300 text-slate-700 rounded-xl text-xs font-semibold hover:bg-slate-50 transition-all cursor-pointer"
                >
                  Batal
                </button>
                <button
                  type="button"
                  @click="confirmSubmitHotProspect"
                  class="px-5 py-2.5 bg-rose-600 hover:bg-rose-700 text-white rounded-xl text-xs font-bold shadow-md shadow-rose-200 transition-all cursor-pointer flex items-center gap-1.5"
                >
                  <span>Ya, Submit Hot Prospect</span>
                  <AppIcon name="check" class="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        </Teleport>
      </ClientOnly>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onUnmounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useQueryClient } from "@tanstack/vue-query";
import {
  useKelurahanQuery,
  useSumberCustQuery,
} from "~/composables/useMasterQueries";
import { getRencanaDetails } from "~/utils/rencanaDetails";
import {
  SOA_LIST,
  getAktifitasOptions,
  getJenisAplikasiOptions,
} from "~/utils/soa";
import {
  compressImage,
  getFotoBarang,
  setFotoBarang,
  removeFotoBarang,
} from "~/utils/imageStorage";

const route = useRoute();
const router = useRouter();

const currentIdRencana = ref<string | number | null>(
  (route.query.id as string) ||
    (route.query.id_rencana as string) ||
    (route.query.id_customer as string) ||
    null,
);
const idRencana = computed(() =>
  currentIdRencana.value ? String(currentIdRencana.value) : null,
);

watch(
  () => [route.query.id, route.query.id_rencana, route.query.id_customer],
  ([newId, newIdRencana, newIdCustomer]) => {
    const validId =
      (newId as string) ||
      (newIdRencana as string) ||
      (newIdCustomer as string);
    if (validId) {
      currentIdRencana.value = validId;
    }
  },
);
const { mutate: inputRencanaMutate, isPending: isMutating } =
  useInputRencanaMutation();
const isSubmitting = computed(() => isMutating.value);
const stepWarning = ref("");

const pipelineSteps = ["Leads", "Prospect", "Hot Prospect", "SBG"];
const refSources = ref<any[]>([]);
const historyList = ref<any[]>([]);

// Wizard Tab & Status Control
const activeTab = ref<string>("Leads");
const dbStatus = ref<string>("Leads");

const customerInfo = reactive({
  id_customer: 0,
  custType: "Pelanggan Lama",
  nama: "Customer",
  source: "Walk In",
  katcust: "1",
  statusRo: "Baru",
  alamat: "",
  rt: "",
  hp: 0,
  rw: "",
  kelurahan: "",
  kecamatan: "",
  kota: "",
  kd_kota: "",
  nm_kota: "",
  keterangan: "",
  keteranganBarang: "",
  soa: "",
  aktifitas: "",
  jenisAplikasi: "",
});

const soaList = SOA_LIST;

const availableCustomerAktifitasOptions = computed(() =>
  getAktifitasOptions(customerInfo.soa),
);
const availableCustomerJenisAplikasiOptions = computed(() =>
  getJenisAplikasiOptions(customerInfo.soa),
);

const onCustomerSoaChange = () => {
  customerInfo.aktifitas = "";
  customerInfo.jenisAplikasi = "";
};

const barangInfo = reactive({
  jenisBarang: "Emas",
  subJenis: "Logam Mulia",
  tipeBarang: "Logam Mulia Antam",
  warnaSpesifikasi: "",
  beratGram: 5,
  beratNetto: 5,
  kadarEmas: "24K",
  kondisiSurat: "Lengkap Surat / Nota Toko",
  kondisiFisikEmas: "Mulus Utuh & Segel Toko",
  kondisiElektronik: "Segel Mulus 99% Normal",
  kelengkapanElektronik: "Fullset (Box + Charger Original)",
  kelengkapan: "Kotak & Sertifikat Toko Lengkap",
  kondisi: "Mulus 99% Utuh & Berfungsi Normal",
  nilaiTaksiran: 6000000,
  maxUp: 5400000,
});

const sbgData = reactive({
  noSbg: "",
  jenisTransaksi: "Gadai",
  jenisBarang: "Emas",
  deskripsiBarang: "",
  nominalUp: 5000000,
  sewaModalPercent: 1.2,
  sewaModal: 60000,
  biayaAdmin: 15000,
  tglPencairan: new Date().toISOString().split("T")[0],
  tglJatuhTempo: new Date(Date.now() + 15 * 86400000)
    .toISOString()
    .split("T")[0],
  catatanSbg: "",
});

const calculateSbgFees = () => {
  if (barangInfo.jenisBarang === "Emas") {
    const karatVal =
      barangInfo.kadarEmas === "24K"
        ? 24
        : barangInfo.kadarEmas === "22K"
          ? 22
          : barangInfo.kadarEmas === "18K"
            ? 18
            : 16;
    const baseRate = 1200000;
    barangInfo.nilaiTaksiran = Math.round(
      (barangInfo.beratGram || 0) * baseRate * (karatVal / 24),
    );
    barangInfo.maxUp = Math.round(barangInfo.nilaiTaksiran * 0.9);
    sbgData.sewaModalPercent = 1.2;
  } else {
    barangInfo.maxUp = Math.round((barangInfo.nilaiTaksiran || 0) * 0.85);
    sbgData.sewaModalPercent = 1.5;
  }

  const rate = (sbgData.sewaModalPercent || 1.2) / 100;
  sbgData.sewaModal = Math.round((sbgData.nominalUp || 0) * rate);
  sbgData.biayaAdmin =
    (sbgData.nominalUp || 0) > 10000000
      ? 35000
      : (sbgData.nominalUp || 0) > 5000000
        ? 25000
        : 15000;
};

// Database Master Offerings & Items
const productOfferings = ref<any[]>([]);
const selectedProdukCode = ref("");
const emasDbList = ref<any[]>([]);
const selectedEmasItems = ref<any[]>([]);
const elekDbList = ref<any[]>([]);
const showElekDbModal = ref(false);
const elekDbSearch = ref("");
const isSearchingElekDb = ref(false);
const barangFotoPreview = ref("");
const barangFotoBase64 = ref("");

const selectJenisGadai = (type: string) => {
  sbgData.jenisTransaksi = type;
  fetchProductOfferings();
};

const toggleJenisBarang = (type: string) => {
  barangInfo.jenisBarang = type;
  sbgData.jenisBarang = type;
  fetchProductOfferings();
  calculateSbgFees();
};

const fetchProductOfferings = async () => {
  try {
    const isEmas = barangInfo.jenisBarang === "Emas";
    const { data } = await useApi("/taksiran/products", {
      params: { type: isEmas ? "emas" : "elektronik" },
    });
    const list = Array.isArray(data) ? data : [];
    productOfferings.value = list;
    if (list.length > 0 && !selectedProdukCode.value) {
      selectedProdukCode.value = list[0].kd_produk;
    }
  } catch (e) {
    console.error("Failed to fetch product offerings:", e);
  }
};

const onProdukChange = () => {
  // Triggers reactivity
};

const fetchEmasDb = async () => {
  try {
    const { data } = await useApi("/taksiran/barang-emas");
    emasDbList.value = Array.isArray(data) ? data : [];
  } catch (e) {
    console.error("Failed to fetch emas items from DB:", e);
  }
};

const getEmasId = (item: any) =>
  item?.id_barang_detail ?? item?.id_detail ?? item?.nama_barang;

const syncEmasTipeBarang = () => {
  barangInfo.tipeBarang = selectedEmasItems.value
    .map((i) => i.nama_barang)
    .join(", ");
  calculateSbgFees();
};

// Emas: boleh pilih lebih dari 1 barang
const isEmasSelected = (item: any) =>
  selectedEmasItems.value.some((i) => getEmasId(i) === getEmasId(item));

const toggleEmasItem = (item: any) => {
  const idx = selectedEmasItems.value.findIndex(
    (i) => getEmasId(i) === getEmasId(item),
  );
  if (idx >= 0) {
    selectedEmasItems.value.splice(idx, 1);
  } else {
    selectedEmasItems.value.push(item);
  }
  syncEmasTipeBarang();
};

const removeEmasItem = (index: number) => {
  selectedEmasItems.value.splice(index, 1);
  syncEmasTipeBarang();
};

const openElekDbSearchModal = () => {
  showElekDbModal.value = true;
  if (elekDbList.value.length === 0) {
    searchElekDb();
  }
};

const searchElekDb = async () => {
  isSearchingElekDb.value = true;
  try {
    const { data } = await useApi("/taksiran/barang-umum", {
      params: { namaBarang: elekDbSearch.value },
    });
    elekDbList.value = Array.isArray(data) ? data : [];
  } catch (e) {
    console.error("Failed to search electronic items in DB:", e);
    elekDbList.value = [];
  } finally {
    isSearchingElekDb.value = false;
  }
};

const selectedElekItems = ref<any[]>([]);

const isElekSelected = (item: any) => {
  const id = item.id_barang_detail || item.id_detail || item.kd_barang;
  return selectedElekItems.value.some(
    (i) => (i.id_barang_detail || i.id_detail || i.kd_barang) === id,
  );
};

// Non Emas (Elektronik): hanya boleh 1 barang
const toggleElekItem = (item: any) => {
  if (isElekSelected(item)) {
    selectedElekItems.value = [];
  } else {
    selectedElekItems.value = [item];
    showElekDbModal.value = false;
  }
  barangInfo.tipeBarang = selectedElekItems.value[0]
    ? selectedElekItems.value[0].nm_barang || selectedElekItems.value[0].name
    : "";
};

const removeElekItem = (index: number) => {
  selectedElekItems.value.splice(index, 1);
  barangInfo.tipeBarang = selectedElekItems.value
    .map((i) => i.nm_barang || i.name)
    .join(", ");
};

const showHotProspectConfirmModal = ref(false);

const promptSubmitHotProspect = () => {
  stepWarning.value = "";
  if (!customerInfo.alamat && !customerInfo.kelurahan) {
    stepWarning.value =
      "Mohon lengkapi Alamat atau Kelurahan customer terlebih dahulu.";
    return;
  }
  if (!customerInfo.keteranganBarang.trim()) {
    stepWarning.value = "Mohon isi Keterangan Barang terlebih dahulu.";
    return;
  }
  showHotProspectConfirmModal.value = true;
};

const confirmSubmitHotProspect = () => {
  showHotProspectConfirmModal.value = false;
  submitWithStatus("Hot Prospect", true);
};

const onFotoBarangChange = async (e: Event) => {
  const target = e.target as HTMLInputElement;
  if (target.files && target.files[0]) {
    try {
      const file = target.files[0];
      const compressed = await compressImage(file);
      barangFotoPreview.value = compressed;
      barangFotoBase64.value = compressed;
      const targetId = idRencana.value || customerInfo.id_customer;
      if (targetId) {
        setFotoBarang(targetId, compressed);
      }
    } catch (err) {
      console.error("Gagal memproses foto barang:", err);
    }
  }
};

const removeBarangFoto = () => {
  barangFotoPreview.value = "";
  barangFotoBase64.value = "";
  const targetId = idRencana.value || customerInfo.id_customer;
  if (targetId) {
    removeFotoBarang(targetId);
  }
};

onMounted(() => {
  fetchProductOfferings();
  fetchEmasDb();
});

const form = reactive({
  newStatus: "Leads",
  tglFollowUp: new Date().toISOString().split("T")[0],
  jamFollowUp: "10:00",
  note: "",
});

// Tab Switcher for manual inspection
const selectTab = (targetTab: string) => {
  stepWarning.value = "";
  activeTab.value = targetTab;
  router.replace({ query: { ...route.query, step: targetTab } });
};

// Navigate backwards one step in wizard (mundur step)
const prevStep = () => {
  stepWarning.value = "";
  if (activeTab.value === "SBG") {
    selectTab("Hot Prospect");
  } else if (activeTab.value === "Hot Prospect") {
    selectTab("Prospect");
  } else if (activeTab.value === "Prospect") {
    selectTab("Leads");
  }
};

const hasInitializedTab = ref(false);

// Kelurahan Search Select State & Query
const kelurahanSearchInput = ref("");
const showKelurahanDropdown = ref(false);
const { data: kelurahanList, isFetching: isSearchingKelurahan } =
  useKelurahanQuery(kelurahanSearchInput);
let kelurahanBlurTimer: ReturnType<typeof setTimeout> | null = null;

const onKelurahanFocus = () => {
  if (kelurahanBlurTimer) {
    clearTimeout(kelurahanBlurTimer);
    kelurahanBlurTimer = null;
  }
  showKelurahanDropdown.value = (customerInfo.kelurahan || "").length >= 2;
};

const onKelurahanBlur = () => {
  if (kelurahanBlurTimer) clearTimeout(kelurahanBlurTimer);
  kelurahanBlurTimer = setTimeout(() => {
    showKelurahanDropdown.value = false;
  }, 200);
};

const onSelectKelurahan = (item: any) => {
  if (kelurahanBlurTimer) {
    clearTimeout(kelurahanBlurTimer);
    kelurahanBlurTimer = null;
  }
  customerInfo.kelurahan = item.nm_kelurahan || "";
  customerInfo.kecamatan = item.nm_kecamatan || "";
  customerInfo.kd_kota = item.kd_kota ? String(item.kd_kota) : "";
  customerInfo.nm_kota = item.nm_kota || "";
  // Masuk ke field kota hanya kodenya saja
  customerInfo.kota = customerInfo.kd_kota;
  kelurahanSearchInput.value = item.nm_kelurahan || "";
  showKelurahanDropdown.value = false;
};

const onKelurahanInput = (e: Event) => {
  if (kelurahanBlurTimer) {
    clearTimeout(kelurahanBlurTimer);
    kelurahanBlurTimer = null;
  }
  const val = (e.target as HTMLInputElement).value;
  customerInfo.kelurahan = val;
  kelurahanSearchInput.value = val;
  showKelurahanDropdown.value = val.length >= 2;
};

onUnmounted(() => {
  if (kelurahanBlurTimer) {
    clearTimeout(kelurahanBlurTimer);
  }
});

// TanStack Queries
const { data: sourcesQueryData } = useSumberCustQuery();
const { data: inquiryQueryData, isPending: isPendingInquiry } =
  useInquiryDetailQuery(idRencana);

const isLoading = computed(() => {
  if (!idRencana.value) return false;
  return isPendingInquiry.value;
});

watch(
  [sourcesQueryData, inquiryQueryData],
  ([sourcesRes, rawList]) => {
    const sources = Array.isArray(sourcesRes)
      ? sourcesRes
      : sourcesRes?.data || [];
    refSources.value = sources;

    if (!rawList || rawList.length === 0) return;
    const id = idRencana.value;
    const found =
      rawList.find(
        (i: any) =>
          String(i.id_rencana) === String(id) ||
          String(i.id_customer) === String(id) ||
          String(i.id) === String(id),
      ) || rawList[0];

    if (found) {
      customerInfo.id_customer =
        found.id_customer || found.id_rencana || found.id || 0;
      customerInfo.nama =
        found.nama || found.nm_customer || found.name || "Customer";
      customerInfo.hp = found.hp || found.phone || found.no_hp || "";

      const matchedSource = sources.find(
        (s: any) => s.kd_ref_cust === found.id_sumbercust,
      );
      customerInfo.source =
        matchedSource?.nm_ref_cust ||
        found.nm_ref_cust ||
        found.source ||
        "Walk In";
      const details = getRencanaDetails(found);
      customerInfo.keterangan = details.keterangan_aktivitas;
      customerInfo.keteranganBarang = details.keterangan_barang;
      customerInfo.soa = details.soa;
      customerInfo.aktifitas = details.aktifitas;
      customerInfo.jenisAplikasi = details.jenis_aplikasi;
      customerInfo.katcust = found.katcust || "1";
      customerInfo.custType =
        found.cust === "Baru" ? "Nasabah Baru" : "Pelanggan Lama";

      customerInfo.alamat = found.alamat || "";
      customerInfo.rt = found.rt || "";
      customerInfo.rw = found.rw || "";
      customerInfo.kelurahan = found.nm_kelurahan || found.kelurahan || "";
      customerInfo.kecamatan = found.nm_kecamatan || found.kecamatan || "";
      customerInfo.kota = found.kota || "";
      customerInfo.kd_kota = found.kota || "";
      customerInfo.nm_kota = found.nm_kota || "";

      // Restore foto barang jaminan from database (url_img), inquiry object, or client storage
      const existingFoto =
        found.url_img ||
        found.foto ||
        found.gambar ||
        getFotoBarang(found.id_rencana) ||
        getFotoBarang(found.id_customer);
      if (existingFoto) {
        barangFotoPreview.value = existingFoto;
        barangFotoBase64.value = existingFoto;
      }

      const statusFromDb = found.nama_hslaktiv || found.status || "Leads";
      dbStatus.value = statusFromDb;
      form.newStatus = statusFromDb;

      if (!hasInitializedTab.value) {
        const queryStep = route.query.step as string;
        if (queryStep && pipelineSteps.includes(queryStep)) {
          activeTab.value = queryStep;
        } else {
          const st = statusFromDb.toLowerCase();
          if (st.includes("lead")) {
            activeTab.value = "Leads";
          } else if (st.includes("hot")) {
            activeTab.value = "Hot Prospect";
          } else if (st.includes("prospect") || st.includes("prospek")) {
            activeTab.value = "Prospect";
          } else {
            // Status di luar Leads, Prospect, dan Hot Prospect -> langsung SBG
            activeTab.value = "SBG";
          }
        }
        hasInitializedTab.value = true;
      }

      // Build history entries matching customer
      const matches = rawList.filter(
        (i: any) =>
          (i.id_customer &&
            found.id_customer &&
            String(i.id_customer) === String(found.id_customer)) ||
          (i.nama &&
            found.nama &&
            i.nama.toLowerCase() === found.nama.toLowerCase()),
      );

      historyList.value = matches.map((m: any) => ({
        tgl:
          m.tgl_aktivitas ||
          m.tgl_rencana ||
          m.tgl ||
          new Date().toISOString().split("T")[0],
        status: m.nama_hslaktiv || m.status || "Leads",
        keterangan: m.ket_aktivitas || m.ket_rencana || m.keterangan || "-",
        source:
          sources.find((s: any) => s.kd_ref_cust === m.id_sumbercust)
            ?.nm_ref_cust ||
          m.nm_ref_cust ||
          "Walk In",
      }));

      // Init SBG defaults
      sbgData.noSbg = `SBG-${new Date().getFullYear()}${(new Date().getMonth() + 1).toString().padStart(2, "0")}-${Math.floor(100 + Math.random() * 900)}`;
      sbgData.deskripsiBarang =
        customerInfo.keteranganBarang || "Barang Gadai Customer";
      sbgData.catatanSbg = `Deal SBG untuk ${customerInfo.nama}`;
      calculateSbgFees();
    }
  },
  { immediate: true },
);

const getStatusBadgeClass = (status: string) => {
  if (!status) return "bg-amber-100 text-amber-800 border-amber-200";
  const st = status.toLowerCase();
  if (st.includes("lead"))
    return "bg-amber-100 text-amber-800 border-amber-200";
  if (st.includes("hot")) return "bg-rose-100 text-rose-800 border-rose-200";
  if (st.includes("prospect") || st.includes("prospek"))
    return "bg-blue-100 text-blue-800 border-blue-200";
  if (st.includes("sbg") || st.includes("deal") || st.includes("do"))
    return "bg-emerald-100 text-emerald-800 border-emerald-200";
  if (st.includes("batal"))
    return "bg-slate-200 text-slate-600 border-slate-300";
  return "bg-slate-100 text-slate-600 border-slate-200";
};

const getStepProgressWidth = (status: string) => {
  switch (status) {
    case "Leads":
      return "0%";
    case "Prospect":
      return "33%";
    case "Hot Prospect":
      return "66%";
    case "SBG":
      return "100%";
    default:
      return "0%";
  }
};

const isStepActive = (stepName: string, currentStatus: string) => {
  const currentIdx = pipelineSteps.indexOf(currentStatus);
  const stepIdx = pipelineSteps.indexOf(stepName);
  return stepIdx <= currentIdx;
};

const statusMap: Record<string, number> = {
  Leads: 1,
  Prospect: 2,
  "Hot Prospect": 3,
  Batal: 6,
  SBG: 7,
};

const queryClient = useQueryClient();

const submitWithStatus = (
  targetStatus: string,
  navigateNext: boolean = false,
  redirectToInquiry: boolean = false,
) => {
  stepWarning.value = "";
  let saveStatus = targetStatus || activeTab.value || dbStatus.value || "Leads";

  if (
    activeTab.value === "Leads" &&
    (!customerInfo.soa ||
      !customerInfo.aktifitas ||
      !customerInfo.jenisAplikasi)
  ) {
    stepWarning.value =
      "Silakan pilih SOA, Aktifitas, dan Jenis Aplikasi terlebih dahulu.";
    return;
  }

  // Guard against status downgrade if data is already Hot Prospect or SBG:
  // Instead of aborting, maintain the higher status so address/notes can still be updated
  const currentDbLevel = statusMap[dbStatus.value] || 1;
  const targetLevel = statusMap[saveStatus] || 1;
  if (currentDbLevel >= 3 && targetLevel < currentDbLevel) {
    saveStatus = dbStatus.value;
  }

  form.newStatus = saveStatus;
  const hslAktivVal = statusMap[saveStatus] || 1;

  let noteText = form.note
    ? `[${saveStatus}] ${form.note}`
    : `Update progress status: ${saveStatus}`;
  if (saveStatus === "Hot Prospect") {
    const prodObj = productOfferings.value.find(
      (p) => p.kd_produk === selectedProdukCode.value,
    );
    const prodName = prodObj?.nm_produk || "";
    const itemNames =
      barangInfo.jenisBarang === "Emas"
        ? selectedEmasItems.value.length > 0
          ? selectedEmasItems.value.map((i) => i.nama_barang).join(", ")
          : barangInfo.tipeBarang || "-"
        : selectedElekItems.value[0]?.nm_barang ||
          selectedElekItems.value[0]?.name ||
          barangInfo.tipeBarang ||
          "-";
    noteText += ` | [${sbgData.jenisTransaksi || "Gadai"}${prodName ? " - " + prodName : ""}] ${barangInfo.jenisBarang}: ${itemNames}`;
  } else if (saveStatus === "SBG") {
    noteText += ` | [${barangInfo.jenisBarang} - ${barangInfo.tipeBarang}] Nilai Taksiran Rp ${(barangInfo.nilaiTaksiran || 0).toLocaleString("id-ID")} | UP Rp ${(sbgData.nominalUp || 0).toLocaleString("id-ID")} | Sewa Modal Rp ${(sbgData.sewaModal || 0).toLocaleString("id-ID")}`;
  }

  const payload: any = {
    id_rencana: idRencana.value ? Number(idRencana.value) : undefined,
    cust: idRencana.value || customerInfo.id_customer ? "Lama" : "Baru",
    custlama:
      customerInfo.id_customer ||
      (idRencana.value ? Number(idRencana.value) : 0),
    nm: customerInfo.nama,
    hp: customerInfo.hp,
    sumbercust: customerInfo.soa || "Marketing",
    hslaktiv: hslAktivVal,
    katcust: customerInfo.katcust,
    ket: customerInfo.keterangan.trim() || noteText,
    // Always send exactly what the user typed so the backend never copies old metadata into it
    keterangan_barang: (customerInfo.keteranganBarang || "").trim(),
    tgl_rencana: form.tglFollowUp,
    jam_rencana: form.jamFollowUp,
    alamat: customerInfo.alamat,
    rt: customerInfo.rt,
    rw: customerInfo.rw,
    kelurahan: customerInfo.kelurahan,
    kecamatan: customerInfo.kecamatan,
    kota: customerInfo.kd_kota || customerInfo.kota,
    kd_kota: customerInfo.kd_kota || customerInfo.kota,
    aktiv: 1,
    soa: customerInfo.soa,
    aktifitas: customerInfo.aktifitas,
    jenis_aplikasi: customerInfo.jenisAplikasi,
  };

  inputRencanaMutate(payload, {
    onSuccess: (res: any) => {
      queryClient.invalidateQueries({ queryKey: ["inquiries"] });
      queryClient.invalidateQueries({ queryKey: ["dashboard-metrics"] });
      dbStatus.value = saveStatus;

      let nextTab = activeTab.value;
      if (navigateNext) {
        if (activeTab.value === "Leads") {
          nextTab = "Prospect";
        } else if (activeTab.value === "Prospect") {
          nextTab = "Hot Prospect";
        } else if (activeTab.value === "Hot Prospect") {
          router.push("/inquiry");
          return;
        } else {
          router.push("/inquiry");
          return;
        }
      } else {
        nextTab = saveStatus;
      }

      activeTab.value = nextTab;

      if (res?.id_rencana) {
        currentIdRencana.value = res.id_rencana;
        if (res.id_customer) {
          customerInfo.id_customer = res.id_customer;
        }
        if (barangFotoBase64.value) {
          setFotoBarang(res.id_rencana, barangFotoBase64.value);
          if (res.id_customer) {
            setFotoBarang(res.id_customer, barangFotoBase64.value);
          }
        }
      } else if (barangFotoBase64.value && customerInfo.id_customer) {
        setFotoBarang(customerInfo.id_customer, barangFotoBase64.value);
      }

      if (redirectToInquiry) {
        router.push("/inquiry");
        return;
      }

      if (res?.id_rencana) {
        router.replace({
          query: { ...route.query, id: String(res.id_rencana), step: nextTab },
        });
      } else {
        router.replace({ query: { ...route.query, step: nextTab } });
      }
    },
    onError: (err: any) => {
      stepWarning.value =
        err.response?.data?.error ||
        err.data?.error ||
        err.message ||
        "Gagal menyimpan transaksi";
    },
  });
};
</script>
