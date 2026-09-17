<template>
  <div class="space-y-6 pb-12 font-sans mx-auto max-w-5xl">
    <!-- Header with Back Button & Quick Actions -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-200/80 pb-4">
      <div class="flex items-center gap-3">
        <NuxtLink to="/inquiry"
          class="p-2.5 bg-white hover:bg-slate-100 border border-slate-200 rounded-xl text-slate-600 transition-all flex items-center gap-1 text-sm font-medium shadow-2xs">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
        </NuxtLink>
        <div>
          <h1 class="text-2xl font-bold text-slate-800 tracking-tight">Progress Sales & Pipeline Wizard</h1>
          <p class="text-xs text-slate-500">Kelola status inquiry dari Leads hingga Realisasi SBG Deal.</p>
        </div>
      </div>
      <div class="flex items-center gap-2">
        <span
          class="text-xs font-semibold px-3 py-1.5 bg-amber-50 text-amber-800 border border-amber-200 rounded-xl font-mono shadow-2xs">
          ID Rencana: #{{ idRencana || 'Baru' }}
        </span>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="bg-white p-12 rounded-2xl border border-slate-200/80 text-center shadow-xs">
      <div class="w-8 h-8 border-3 border-amber-500 border-t-transparent rounded-full animate-spin mx-auto mb-3"></div>
      <p class="text-slate-500 text-sm font-medium">Memuat data inquiry & riwayat customer ID #{{ idRencana }}...</p>
    </div>

    <template v-else>
      <!-- STEPPER WIZARD TRACKER (Leads -> Prospect -> Hot Prospect -> SBG) -->
      <div class="bg-white rounded-2xl border border-slate-200/80 p-6 shadow-xs space-y-6">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
          <div>
            <h2 class="font-bold text-slate-800 text-base">Wizard Status Pipeline</h2>
            <p class="text-xs text-slate-500">Progress status nasabah saat ini di database: <strong class="text-amber-700 font-bold">{{ dbStatus }}</strong></p>
          </div>
          <div class="flex items-center gap-2">
            <span class="text-xs font-semibold px-3 py-1 rounded-full border"
              :class="getStatusBadgeClass(activeTab)">
              Tahap Aktif: {{ activeTab }}
            </span>
          </div>
        </div>

        <!-- Progress Bar & Stepper Buttons -->
        <div class="py-2 px-2">
          <div class="flex items-center justify-between relative">
            <div class="absolute left-0 right-0 top-1/2 -translate-y-1/2 h-1 bg-slate-100 -z-0"></div>
            <div class="absolute left-0 top-1/2 -translate-y-1/2 h-1 bg-amber-500 transition-all duration-500 -z-0"
              :style="{ width: getStepProgressWidth(activeTab) }"></div>

            <!-- Step 1: Leads -->
            <button type="button" @click="selectTab('Leads')"
              class="relative z-10 flex flex-col items-center gap-1.5 cursor-pointer transition-all">
              <div class="w-10 h-10 rounded-full flex items-center justify-center text-sm font-bold transition-all"
                :class="isStepActive('Leads', activeTab) ? 'bg-amber-500 text-white shadow-md shadow-amber-200 ring-4 ring-amber-100' : 'bg-slate-100 text-slate-400'">
                1</div>
              <span class="text-xs font-semibold uppercase tracking-wider"
                :class="activeTab === 'Leads' ? 'text-amber-700 font-bold' : 'text-slate-400'">1. Leads</span>
            </button>

            <!-- Step 2: Prospect -->
            <button type="button" @click="selectTab('Prospect')"
              class="relative z-10 flex flex-col items-center gap-1.5 cursor-pointer transition-all">
              <div class="w-10 h-10 rounded-full flex items-center justify-center text-sm font-bold transition-all"
                :class="isStepActive('Prospect', activeTab) ? 'bg-amber-600 text-white shadow-md shadow-amber-200 ring-4 ring-amber-100' : 'bg-slate-100 text-slate-400'">
                2</div>
              <span class="text-xs font-semibold uppercase tracking-wider"
                :class="activeTab === 'Prospect' ? 'text-amber-700 font-bold' : 'text-slate-400'">2. Prospek</span>
            </button>

            <!-- Step 3: Hot Prospect -->
            <button type="button" @click="selectTab('Hot Prospect')"
              class="relative z-10 flex flex-col items-center gap-1.5 cursor-pointer transition-all">
              <div class="w-10 h-10 rounded-full flex items-center justify-center text-sm font-bold transition-all"
                :class="isStepActive('Hot Prospect', activeTab) ? 'bg-amber-700 text-white shadow-md shadow-amber-200 ring-4 ring-amber-100' : 'bg-slate-100 text-slate-400'">
                3</div>
              <span class="text-xs font-semibold uppercase tracking-wider"
                :class="activeTab === 'Hot Prospect' ? 'text-amber-800 font-bold' : 'text-slate-400'">3. Hot Prospek</span>
            </button>

            <!-- Step 4: SBG / Deal -->
            <button type="button" @click="selectTab('SBG')"
              class="relative z-10 flex flex-col items-center gap-1.5 cursor-pointer transition-all">
              <div class="w-10 h-10 rounded-full flex items-center justify-center text-xs font-bold transition-all"
                :class="isStepActive('SBG', activeTab) ? 'bg-yellow-600 text-white shadow-md shadow-amber-200 ring-4 ring-amber-100' : 'bg-slate-100 text-slate-400'">
                4</div>
              <span class="text-xs font-semibold uppercase tracking-wider"
                :class="activeTab === 'SBG' ? 'text-yellow-700 font-bold' : 'text-slate-400'">4. SBG (Deal)</span>
            </button>
          </div>
        </div>

        <!-- Warning Alert Banner -->
        <div v-if="stepWarning"
          class="p-3.5 bg-amber-50 border border-amber-300 rounded-xl text-xs font-medium text-amber-800 flex items-center gap-2.5">
          <svg class="w-4 h-4 text-amber-600 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
          <span>{{ stepWarning }}</span>
        </div>
      </div>

      <!-- RIWAYAT FOLLOW-UP & LOG AKTIVITAS -->
      <div v-if="historyList.length > 0"
        class="bg-white rounded-2xl border border-slate-200/80 p-6 shadow-xs space-y-4">
        <div class="flex items-center justify-between border-b border-slate-100 pb-3">
          <div>
            <h2 class="font-semibold text-slate-800 text-base">Riwayat Aktivitas & Timeline (v_his_inq)</h2>
            <p class="text-xs text-slate-500">Catatan perkembangan komunikasi sebelumnya untuk customer ini.</p>
          </div>
          <span class="text-xs font-semibold px-3 py-1 bg-amber-50 text-amber-800 border border-amber-200 rounded-full">
            {{ historyList.length }} Catatan Aktivitas
          </span>
        </div>

        <div class="relative border-l-2 border-slate-200 ml-3 space-y-4 py-1">
          <div v-for="(h, idx) in historyList" :key="idx" class="relative pl-6">
            <div class="absolute -left-[9px] top-1.5 w-4 h-4 rounded-full bg-white border-2 border-amber-500"></div>
            <div class="bg-slate-50 border border-slate-200/80 rounded-xl p-3 space-y-1">
              <div class="flex items-center justify-between">
                <span class="text-xs font-semibold px-2 py-0.5 rounded-md border"
                  :class="getStatusBadgeClass(h.status)">
                  {{ h.status }}
                </span>
                <span class="text-xs text-slate-400 font-mono">{{ h.tgl }}</span>
              </div>
              <p class="text-xs text-slate-700 leading-relaxed font-medium">{{ h.keterangan }}</p>
              <span class="text-[11px] text-slate-400 block">Sumber: {{ h.source }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- ========================================================================= -->
      <!-- PANEL 1: LEADS (DATA KONTAK AWAL) -->
      <!-- ========================================================================= -->
      <div v-if="activeTab === 'Leads'"
        class="bg-white rounded-2xl border border-amber-200 p-6 shadow-xs space-y-5">
        <div class="flex items-center justify-between border-b border-slate-100 pb-3">
          <div class="flex items-center gap-2">
            <span
              class="w-7 h-7 rounded-full bg-amber-500 text-white font-bold text-xs flex items-center justify-center">1</span>
            <h2 class="font-bold text-slate-800 text-base">PANEL STEP 1 — LEADS (Data Kontak Awal)</h2>
          </div>
          <span
            class="text-xs font-medium text-amber-800 bg-amber-50 px-3 py-1 rounded-lg border border-amber-200">Status:
            Leads</span>
        </div>

        <div
          class="p-3.5 bg-amber-50/60 border border-amber-200 rounded-xl text-xs text-amber-900 flex items-center gap-2">
          <svg class="w-4 h-4 text-amber-600 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span>Isi data kontak awal customer. Nomor HP dari sistem dikunci (readonly) demi validasi.</span>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1">Nama Lengkap
              Customer <span class="text-rose-500">*</span></label>
            <input v-model="customerInfo.nama" type="text"
              class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm font-semibold text-slate-800" />
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1">No. Handphone /
              WhatsApp (Readonly)</label>
            <input v-model="customerInfo.hp" type="text"
              class="w-full px-3.5 py-2.5 bg-slate-100 border border-slate-300 rounded-xl text-sm font-medium text-slate-700 cursor-not-allowed"
              readonly />
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1">Sumber Customer
              (Source) <span class="text-rose-500">*</span></label>
            <select v-model="customerInfo.source"
              class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-800">
              <option v-for="s in refSources" :key="s.kd_ref_cust" :value="s.nm_ref_cust">
                {{ s.nm_ref_cust }}
              </option>
              <option v-if="refSources.length === 0" value="Walk In">Walk In (Kunjungan Outlet)</option>
              <option v-if="refSources.length === 0" value="Media Sosial">Media Sosial (FB / IG / WA)</option>
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

        <div>
          <label class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1">Catatan Follow Up
            Leads</label>
          <textarea v-model="form.note" rows="2" placeholder="Catatan kontak awal..."
            class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm"></textarea>
        </div>

        <div class="flex flex-wrap items-center justify-between gap-3 pt-3 border-t border-slate-100">
          <NuxtLink to="/inquiry" class="px-4 py-2 border border-slate-300 text-slate-600 rounded-xl text-xs font-semibold hover:bg-slate-50 transition-all">
            Batal & Kembali
          </NuxtLink>

          <div class="flex items-center gap-2">
            <button @click="submitWithStatus('Leads', false)" type="button"
              class="px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-xl text-xs font-semibold transition-all cursor-pointer"
              :disabled="isSubmitting">
              Simpan Data Leads
            </button>
            <button @click="submitWithStatus('Prospect', true)" type="button"
              class="px-5 py-2 bg-amber-500 hover:bg-amber-600 text-white rounded-xl text-xs font-bold shadow-xs flex items-center gap-1.5 cursor-pointer transition-all"
              :disabled="isSubmitting">
              <span v-if="isSubmitting" class="w-3.5 h-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
              <span v-else>Simpan & Lanjut ke Prospek</span>
              <svg v-if="!isSubmitting" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- ========================================================================= -->
      <!-- PANEL 2: PROSPEK (DATA ALAMAT & SERVE CONSULTATION) -->
      <!-- ========================================================================= -->
      <div v-else-if="activeTab === 'Prospect'"
        class="bg-white rounded-2xl border border-blue-200 p-6 shadow-xs space-y-5">
        <div class="flex items-center justify-between border-b border-slate-100 pb-3">
          <div class="flex items-center gap-2">
            <span
              class="w-7 h-7 rounded-full bg-blue-600 text-white font-bold text-xs flex items-center justify-center">2</span>
            <h2 class="font-bold text-slate-800 text-base">PANEL STEP 2 — PROSPEK (Alamat & Konsultasi Taksiran)</h2>
          </div>
          <span class="text-xs font-medium text-blue-800 bg-blue-50 px-3 py-1 rounded-lg border border-blue-200">Status:
            Prospek</span>
        </div>

        <div
          class="p-3.5 bg-blue-50/60 border border-blue-200 rounded-xl text-xs text-blue-900 flex items-center gap-2">
          <svg class="w-4 h-4 text-blue-600 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
          </svg>
          <span>Lengkapi alamat domisili customer dan catatan hasil konsultasi awal taksiran barang.</span>
        </div>

        <!-- Alamat Form -->
        <div class="space-y-4">
          <div>
            <label class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1">Jalan / Alamat
              Lengkap <span class="text-rose-500">*</span></label>
            <input v-model="customerInfo.alamat" type="text" placeholder="Jalan, Perumahan, No. Rumah"
              class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm" />
          </div>

          <div class="grid grid-cols-2 md:grid-cols-5 gap-3">
            <div>
              <label class="block text-xs font-medium text-slate-500 mb-1">RT</label>
              <input v-model="customerInfo.rt" type="text" placeholder="001"
                class="w-full px-3 py-2 bg-white border border-slate-300 rounded-xl text-sm" />
            </div>
            <div>
              <label class="block text-xs font-medium text-slate-500 mb-1">RW</label>
              <input v-model="customerInfo.rw" type="text" placeholder="002"
                class="w-full px-3 py-2 bg-white border border-slate-300 rounded-xl text-sm" />
            </div>

            <!-- Search Select Kelurahan -->
            <div class="relative col-span-2 md:col-span-1">
              <label
                class="text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1 flex items-center justify-between">
                <span>Kelurahan</span>
                <span class="text-amber-600 text-[10px] lowercase font-normal">Cari Master</span>
              </label>
              <div class="relative">
                <input :value="customerInfo.kelurahan" @input="onKelurahanInput"
                  @focus="showKelurahanDropdown = (customerInfo.kelurahan || '').length >= 2"
                  @blur="setTimeout(() => { showKelurahanDropdown = false }, 200)" type="text"
                  placeholder="Ketik kelurahan..."
                  class="w-full px-3 py-2 bg-white border border-slate-300 rounded-xl text-sm font-medium focus:ring-2 focus:ring-amber-200 focus:border-amber-500 outline-none" />
                <span v-if="isSearchingKelurahan"
                  class="absolute right-2.5 top-2.5 w-3.5 h-3.5 border-2 border-amber-500 border-t-transparent rounded-full animate-spin"></span>
              </div>

              <!-- Search Results Dropdown -->
              <div v-if="showKelurahanDropdown && kelurahanList && kelurahanList.length > 0"
                class="absolute left-0 right-0 top-full mt-1 bg-white border border-slate-200 rounded-xl shadow-lg z-50 max-h-48 overflow-y-auto divide-y divide-slate-100">
                <div v-for="(k, idx) in kelurahanList" :key="idx" @mousedown.prevent="onSelectKelurahan(k)"
                  class="p-2.5 hover:bg-amber-50 cursor-pointer transition-all text-xs">
                  <div class="font-bold text-slate-800">{{ k.nm_kelurahan }}</div>
                  <div class="text-slate-500 text-[11px] flex items-center gap-1">
                    <span>Kec. {{ k.nm_kecamatan }}</span>
                    <span>•</span>
                    <span class="text-amber-700 font-medium">{{ k.nm_kota }}</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Auto-filled Kecamatan -->
            <div>
              <label class="block text-xs font-medium text-slate-500 mb-1 flex items-center justify-between">
                <span>Kecamatan</span>
                <span v-if="customerInfo.kecamatan" class="text-[10px] text-emerald-600 font-semibold">Auto</span>
              </label>
              <input v-model="customerInfo.kecamatan" type="text" placeholder="Kecamatan"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-xl text-sm font-medium text-slate-800" />
            </div>

            <!-- Auto-filled Kota -->
            <div class="col-span-2 md:col-span-1">
              <label class="text-xs font-medium text-slate-500 mb-1 flex items-center justify-between">
                <span>Kota / Kab</span>
                <span v-if="customerInfo.kota" class="text-[10px] text-emerald-600 font-semibold">Auto</span>
              </label>
              <input v-model="customerInfo.kota" type="text" placeholder="Kota / Kab"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-xl text-sm font-medium text-slate-800" />
            </div>
          </div>
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1">Keterangan Aktivitas
            Prospek</label>
          <textarea v-model="customerInfo.keterangan" rows="2"
            placeholder="Hasil konsultasi taksiran atau respon customer..."
            class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm"></textarea>
        </div>

        <div class="flex flex-wrap items-center justify-between gap-3 pt-3 border-t border-slate-100">
          <NuxtLink to="/inquiry" class="px-4 py-2 border border-slate-300 text-slate-600 rounded-xl text-xs font-semibold hover:bg-slate-50 transition-all">
            Batal & Kembali
          </NuxtLink>

          <div class="flex flex-wrap items-center gap-2">
            <!-- Save while maintaining Prospect status -->
            <button @click="submitWithStatus('Prospect', false)" type="button"
              class="px-4 py-2 bg-blue-50 hover:bg-blue-100 text-blue-800 border border-blue-200 rounded-xl text-xs font-bold transition-all cursor-pointer"
              :disabled="isSubmitting">
              Simpan Update Prospek
            </button>

            <!-- Save and upgrade status to Hot Prospect -->
            <button @click="submitWithStatus('Hot Prospect', true)" type="button"
              class="px-5 py-2 bg-rose-600 hover:bg-rose-700 text-white rounded-xl text-xs font-bold flex items-center gap-1.5 shadow-sm cursor-pointer transition-all"
              :disabled="isSubmitting">
              <span v-if="isSubmitting" class="w-3.5 h-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
              <span v-else>Simpan & Lanjut ke Hot Prospek</span>
              <svg v-if="!isSubmitting" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- ========================================================================= -->
      <!-- PANEL 3: HOT PROSPEK (Detail Transaksi, Produk & Barang Jaminan DB)      -->
      <!-- ========================================================================= -->
      <div v-else-if="activeTab === 'Hot Prospect'"
        class="bg-white rounded-2xl border border-rose-200 p-6 shadow-xs space-y-6">
        <div class="flex items-center justify-between border-b border-slate-100 pb-3">
          <div class="flex items-center gap-2">
            <span
              class="w-7 h-7 rounded-full bg-rose-600 text-white font-bold text-xs flex items-center justify-center">3</span>
            <h2 class="font-bold text-slate-800 text-base">PANEL STEP 3 — HOT PROSPEK (Detail Transaksi & Barang Jaminan DB)</h2>
          </div>
          <span class="text-xs font-bold text-rose-800 bg-rose-50 px-3 py-1 rounded-lg border border-rose-200">
            Status DB: {{ dbStatus }}
          </span>
        </div>

        <div
          class="p-3.5 bg-rose-50/60 border border-rose-200 rounded-xl text-xs text-rose-900 flex items-center gap-2">
          <svg class="w-4 h-4 text-rose-600 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M17.657 18.657A8 8 0 016.343 7.343S7 9 9 10c0-2 .5-5 2.986-7C14 5 16.09 5.777 17.656 7.343A7.975 7.975 0 0120 13a7.975 7.975 0 01-2.343 5.657z" />
          </svg>
          <span>Langkah final sebelum realisasi SBG. Data produk & barang jaminan terhubung langsung ke Master Database Selada.</span>
        </div>

        <!-- 1. Pilihan Kartu Jenis Gadai (Gadai vs Cicilan) -->
        <div class="space-y-2">
          <label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider">
            Jenis Gadai <span class="text-rose-500">*</span>
          </label>
          <div class="grid grid-cols-2 gap-4 max-w-md">
            <div @click="selectJenisGadai('Gadai')"
              class="p-4 rounded-xl border-2 text-center cursor-pointer transition-all"
              :class="sbgData.jenisTransaksi === 'Gadai' ? 'border-amber-500 bg-amber-50/60 shadow-xs' : 'border-slate-200 hover:border-slate-300 bg-white'">
              <div class="text-2xl mb-1">💎</div>
              <div class="font-bold text-sm text-slate-800">Gadai</div>
              <div class="text-xs text-slate-500">Gadai Barang</div>
            </div>

            <div @click="selectJenisGadai('Cicilan')"
              class="p-4 rounded-xl border-2 text-center cursor-pointer transition-all"
              :class="sbgData.jenisTransaksi === 'Cicilan' ? 'border-amber-500 bg-amber-50/60 shadow-xs' : 'border-slate-200 hover:border-slate-300 bg-white'">
              <div class="text-2xl mb-1">💳</div>
              <div class="font-bold text-sm text-slate-800">Cicilan</div>
              <div class="text-xs text-slate-500">Cicil Barang</div>
            </div>
          </div>
        </div>

        <!-- 2. Nama Produk dari Database (product.product_offerings) -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 pt-2 border-t border-slate-100">
          <div>
            <label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">
              Nama Produk <span class="text-rose-500">*</span> <span class="text-[10px] text-amber-600 lowercase font-normal">(dari Master DB)</span>
            </label>
            <select v-model="selectedProdukCode" @change="onProdukChange" class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm font-medium">
              <option value="">-- Pilih Produk --</option>
              <option v-for="p in productOfferings" :key="p.kd_produk" :value="p.kd_produk">
                {{ p.nm_produk }}
              </option>
            </select>
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">
              Jenis Barang Jaminan <span class="text-rose-500">*</span>
            </label>
            <div class="grid grid-cols-2 gap-2">
              <button type="button" @click="toggleJenisBarang('Emas')"
                class="py-2.5 px-4 rounded-xl text-xs font-bold border transition-all flex items-center justify-center gap-1.5 cursor-pointer"
                :class="barangInfo.jenisBarang === 'Emas' ? 'bg-amber-500 text-white border-amber-600 shadow-xs' : 'bg-white text-slate-600 border-slate-200 hover:bg-slate-50'">
                <span>💍</span> Emas
              </button>
              <button type="button" @click="toggleJenisBarang('Elektronik')"
                class="py-2.5 px-4 rounded-xl text-xs font-bold border transition-all flex items-center justify-center gap-1.5 cursor-pointer"
                :class="barangInfo.jenisBarang === 'Elektronik' ? 'bg-blue-600 text-white border-blue-700 shadow-xs' : 'bg-white text-slate-600 border-slate-200 hover:bg-slate-50'">
                <span>📦</span> Non Emas (Elektronik)
              </button>
            </div>
          </div>
        </div>

        <!-- 3. CARI & PILIH BARANG JAMINAN DARI DATABASE -->
        <div class="space-y-4 pt-2 border-t border-slate-100">
          <div class="flex items-center justify-between">
            <label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider">
              Cari & Pilih Barang Jaminan di Database <span class="text-rose-500">*</span>
            </label>
            <span class="text-xs font-semibold text-slate-500">
              {{ barangInfo.jenisBarang === 'Emas' ? 'Master Emas DB' : 'Master Elektronik DB' }}
            </span>
          </div>

          <!-- IF EMAS: Dropdown / Search dari /api/taksiran/barang-emas -->
          <div v-if="barangInfo.jenisBarang === 'Emas'" class="space-y-4">
            <div>
              <label class="block text-xs font-medium text-slate-600 mb-1">Pilih Barang Emas dari Database <span class="text-rose-500">*</span></label>
              <select v-model="selectedEmasDbItem" @change="onSelectEmasDb" class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm font-medium">
                <option value="">-- Pilih Barang Emas --</option>
                <option v-for="b in emasDbList" :key="b.id_barang_detail" :value="b">
                  {{ b.nama_barang }} ({{ b.status_antam }})
                </option>
              </select>
            </div>
          </div>

          <!-- IF ELEKTRONIK: Modal / Search dari /api/taksiran/barang-umum (Bisa Multiselect) -->
          <div v-else class="space-y-3">
            <div class="flex items-center justify-between">
              <label class="block text-xs font-medium text-slate-600">
                Pilih Barang Elektronik dari Master DB <span class="text-rose-500">*</span>
                <span class="text-[11px] text-blue-600 font-normal ml-1">(Bisa pilih lebih dari 1 barang)</span>
              </label>
              <button type="button" @click="openElekDbSearchModal"
                class="px-3.5 py-1.5 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-bold transition-all shadow-xs cursor-pointer flex items-center gap-1.5">
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
                <span>Cari & Tambah Barang DB</span>
              </button>
            </div>

            <!-- Chips list of selected electronic items -->
            <div v-if="selectedElekItems.length > 0" class="flex flex-wrap items-center gap-2 p-3 bg-slate-50 border border-slate-200 rounded-xl min-h-[52px]">
              <span v-for="(it, idx) in selectedElekItems" :key="idx"
                class="inline-flex items-center gap-2 px-3 py-1.5 bg-white border border-blue-200 text-blue-900 text-xs font-semibold rounded-lg shadow-2xs">
                <span>📦 {{ it.nm_barang || it.name }}</span>
                <span v-if="it.kd_barang" class="text-[10px] text-slate-400 font-mono">({{ it.kd_barang }})</span>
                <button type="button" @click="removeElekItem(idx)" class="hover:text-rose-600 text-slate-400 font-bold ml-0.5 cursor-pointer">✕</button>
              </span>
            </div>

            <div v-else @click="openElekDbSearchModal"
              class="p-4 border-2 border-dashed border-slate-200 hover:border-blue-300 rounded-xl text-center text-slate-400 text-xs cursor-pointer hover:bg-blue-50/40 transition-all flex items-center justify-center gap-2">
              <svg class="w-4 h-4 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
              <span>Belum ada barang elektronik dipilih. Klik untuk mencari & memilih barang dari Master DB (Bisa Multiselect)</span>
            </div>
          </div>
        </div>

        <!-- 4. FOTO BARANG (Ambil Foto / Upload) -->
        <div class="space-y-3 pt-2 border-t border-slate-100">
          <label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider">
            Foto Barang Jaminan <span class="text-rose-500">*</span>
          </label>
          <div class="flex flex-wrap items-center gap-4">
            <div class="w-24 h-24 rounded-2xl border-2 border-dashed border-slate-300 flex items-center justify-center overflow-hidden bg-slate-50">
              <img v-if="barangFotoPreview" :src="barangFotoPreview" alt="Foto Barang" class="w-full h-full object-cover" />
              <div v-else class="text-center text-slate-400 p-2">
                <svg class="w-6 h-6 mx-auto mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"/></svg>
                <span class="text-[10px] block leading-tight">Belum Ada Foto</span>
              </div>
            </div>
            <div>
              <label class="inline-flex items-center gap-2 px-4 py-2.5 bg-slate-800 hover:bg-slate-900 text-white rounded-xl text-xs font-semibold cursor-pointer shadow-xs transition-all">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"/></svg>
                Ambil / Upload Foto Barang
                <input type="file" accept="image/*" @change="onFotoBarangChange" class="hidden" />
              </label>
              <span class="block text-[11px] text-slate-400 mt-1">Format: JPG/PNG, ukuran maks 5MB</span>
            </div>
          </div>
        </div>

        <!-- 5. KETERANGAN BARANG -->
        <div class="space-y-4 pt-2 border-t border-slate-100">
          <div>
            <label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">
              Keterangan Barang <span class="text-rose-500">*</span>
            </label>
            <input v-model="customerInfo.keterangan" type="text" placeholder="Masukkan keterangan detail barang jaminan..."
              class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm" />
          </div>
        </div>

        <div class="flex flex-wrap items-center justify-between gap-3 pt-3 border-t border-slate-100">
          <NuxtLink to="/inquiry" class="px-4 py-2.5 border border-slate-300 text-slate-600 rounded-xl text-xs font-semibold hover:bg-slate-50 transition-all">
            Batal & Kembali
          </NuxtLink>

          <button @click="submitHotProspect" type="button"
            class="px-6 py-2.5 bg-rose-600 hover:bg-rose-700 text-white rounded-xl text-xs font-bold shadow-md shadow-rose-200 flex items-center gap-2 cursor-pointer transition-all"
            :disabled="isSubmitting">
            <span v-if="isSubmitting" class="w-3.5 h-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
            <span v-else>Submit Hot Prospect</span>
            <svg v-if="!isSubmitting" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </button>
        </div>
      </div>

      <!-- ========================================================================= -->
      <!-- PANEL 4: SBG / DEAL (FORM REALISASI AKAD & TRANSAKSI) -->
      <!-- ========================================================================= -->
      <div v-else-if="activeTab === 'SBG'"
        class="bg-white rounded-2xl border border-emerald-300 p-6 shadow-xs space-y-5">
        <div class="flex items-center justify-between border-b border-slate-100 pb-3">
          <div class="flex items-center gap-2">
            <span
              class="w-7 h-7 rounded-full bg-emerald-600 text-white font-bold text-xs flex items-center justify-center">4</span>
            <h2 class="font-bold text-slate-800 text-base">PANEL STEP 4 — SBG (Realisasi Transaksi Deal SBG)</h2>
          </div>
          <span
            class="text-xs font-bold text-emerald-800 bg-emerald-100 px-3 py-1 rounded-lg border border-emerald-300">Status:
            SBG Deal</span>
        </div>

        <div
          class="p-3.5 bg-emerald-50 border border-emerald-200 rounded-xl text-xs text-emerald-900 flex items-center gap-2">
          <svg class="w-4 h-4 text-emerald-600 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span>Form pencairan Surat Bukti Gadai (SBG). Pastikan nomor SBG & tanggal jatuh tempo terverifikasi.</span>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label class="block text-xs font-semibold text-emerald-900 uppercase tracking-wider mb-1">Nomor SBG
              Transaksi <span class="text-rose-500">*</span></label>
            <input v-model="sbgData.noSbg" type="text" placeholder="SBG-202609-001"
              class="w-full px-3.5 py-2.5 bg-white border border-emerald-300 rounded-xl text-sm font-mono font-bold text-emerald-900"
              required />
          </div>

          <div>
            <label class="block text-xs font-semibold text-emerald-900 uppercase tracking-wider mb-1">Skema
              Transaksi</label>
            <input v-model="sbgData.jenisTransaksi" type="text"
              class="w-full px-3.5 py-2.5 bg-slate-50 border border-emerald-300 rounded-xl text-sm font-medium text-emerald-900"
              readonly />
          </div>

          <div>
            <label class="block text-xs font-semibold text-emerald-900 uppercase tracking-wider mb-1">Jenis Produk
              Gadai</label>
            <select v-model="sbgData.jenisBarang"
              class="w-full px-3.5 py-2.5 bg-white border border-emerald-300 rounded-xl text-sm font-medium text-emerald-900">
              <option value="Emas">Emas / Perhiasan</option>
              <option value="Kendaraan">Kendaraan (Motor/Mobil)</option>
              <option value="Elektronik">Elektronik (HP/Laptop)</option>
            </select>
          </div>
        </div>

        <div>
          <label class="block text-xs font-semibold text-emerald-900 uppercase tracking-wider mb-1">Deskripsi Barang
            Gadai SBG</label>
          <input v-model="sbgData.deskripsiBarang" type="text" placeholder="Deskripsi barang jaminan akad..."
            class="w-full px-3.5 py-2.5 bg-white border border-emerald-300 rounded-xl text-sm" />
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1">Tanggal Pencairan
              SBG</label>
            <input v-model="sbgData.tglPencairan" type="date"
              class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm" />
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1">Tanggal Jatuh
              Tempo</label>
            <input v-model="sbgData.tglJatuhTempo" type="date"
              class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm" />
          </div>
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1">Catatan Realisasi SBG
            / Akad</label>
          <textarea v-model="sbgData.catatanSbg" rows="2" placeholder="Catatan transaksi deal SBG..."
            class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm"></textarea>
        </div>

        <div class="flex flex-wrap items-center justify-between gap-3 pt-3 border-t border-slate-100">
          <NuxtLink to="/inquiry" class="px-4 py-2 border border-slate-300 text-slate-600 rounded-xl text-xs font-semibold hover:bg-slate-50 transition-all">
            Batal & Kembali
          </NuxtLink>

          <button @click="submitWithStatus('SBG', true)" type="button"
            class="px-6 py-2.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-xs font-bold shadow-md shadow-emerald-100 flex items-center gap-1.5 cursor-pointer transition-all"
            :disabled="isSubmitting">
            <span v-if="isSubmitting" class="w-3.5 h-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
            <span v-else>Simpan & Submit Deal SBG</span>
          </button>
        </div>
      </div>

      <!-- MODAL SEARCH BARANG ELEKTRONIK DARI DATABASE (MULTISELECT) -->
      <div v-if="showElekDbModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs z-50 flex items-center justify-center p-4">
        <div class="bg-white rounded-2xl max-w-xl w-full p-6 shadow-2xl space-y-4 max-h-[85vh] flex flex-col">
          <div class="flex items-center justify-between border-b border-slate-100 pb-3">
            <div>
              <h3 class="font-bold text-slate-800 text-base">Cari & Pilih Barang Elektronik</h3>
              <p class="text-xs text-slate-400">Pilih satu atau beberapa barang dari Master Database</p>
            </div>
            <button @click="showElekDbModal = false" class="text-slate-400 hover:text-slate-600 p-1 text-lg">✕</button>
          </div>

          <div class="flex gap-2">
            <input
              v-model="elekDbSearch"
              @keyup.enter="searchElekDb"
              type="text"
              placeholder="Ketik merk / tipe HP, Laptop, TV, dsb..."
              class="flex-1 px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm font-medium outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-200"
            />
            <button @click="searchElekDb" class="px-5 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-bold cursor-pointer">
              Cari
            </button>
          </div>

          <div class="flex-1 overflow-y-auto space-y-2 pr-1 divide-y divide-slate-100">
            <div v-if="isSearchingElekDb" class="py-8 text-center text-slate-400 text-xs">
              <div class="w-5 h-5 border-2 border-blue-600 border-t-transparent rounded-full animate-spin mx-auto mb-2"></div>
              Mencari barang di database...
            </div>
            <div v-else-if="elekDbList.length === 0" class="py-8 text-center text-slate-400 text-xs">
              Ketik pencarian barang dan tekan tombol Cari.
            </div>
            <div
              v-for="item in elekDbList"
              :key="item.id_barang_detail || item.id_detail"
              @click="toggleElekItem(item)"
              class="p-3 rounded-xl cursor-pointer transition-all flex items-center justify-between gap-3 border"
              :class="isElekSelected(item) ? 'bg-blue-50/80 border-blue-300 text-blue-900' : 'hover:bg-slate-50 border-transparent text-slate-800'"
            >
              <div class="flex items-center gap-3">
                <input
                  type="checkbox"
                  :checked="isElekSelected(item)"
                  @click.stop="toggleElekItem(item)"
                  class="w-4 h-4 rounded text-blue-600 focus:ring-blue-500 cursor-pointer"
                />
                <div>
                  <span class="font-semibold text-sm block">{{ item.nm_barang }}</span>
                  <span class="text-xs text-slate-400 font-mono">{{ item.kd_barang }}</span>
                </div>
              </div>
              <div class="text-right">
                <span class="text-xs text-slate-400 block">Taksiran Standar</span>
                <span class="text-xs font-bold text-emerald-700">Rp {{ Number(item.harga_taksir || item.harga || 0).toLocaleString('id-ID') }}</span>
              </div>
            </div>
          </div>

          <div class="flex items-center justify-between pt-3 border-t border-slate-100">
            <span class="text-xs text-slate-500 font-medium">
              <strong class="text-blue-600">{{ selectedElekItems.length }}</strong> barang terpilih
            </span>
            <div class="flex items-center gap-2">
              <button @click="showElekDbModal = false" class="px-4 py-2 border border-slate-300 text-slate-600 rounded-xl text-xs font-semibold hover:bg-slate-50">
                Batal
              </button>
              <button @click="showElekDbModal = false" class="px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-bold cursor-pointer">
                Selesai ({{ selectedElekItems.length }})
              </button>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQueryClient } from '@tanstack/vue-query'
import { useKelurahanQuery, useSumberCustQuery } from '~/composables/queries/useMasterQueries'

const route = useRoute()
const router = useRouter()

const idRencana = computed(() => (route.query.id as string) || null)
const { mutate: inputRencanaMutate, isPending: isMutating } = useInputRencanaMutation()
const isSubmitting = computed(() => isMutating.value)
const stepWarning = ref('')

const pipelineSteps = ['Leads', 'Prospect', 'Hot Prospect', 'SBG']
const refSources = ref<any[]>([])
const historyList = ref<any[]>([])

// Wizard Tab & Status Control
const activeTab = ref<string>('Leads')
const dbStatus = ref<string>('Leads')

const customerInfo = reactive({
  id_customer: 0,
  custType: 'Pelanggan Lama',
  nama: 'Customer',
  source: 'Walk In',
  katcust: '1',
  statusRo: 'Baru',
  alamat: '',
  rt: '',
  hp: 0,
  rw: '',
  kelurahan: '',
  kecamatan: '',
  kota: '',
  keterangan: ''
})

const barangInfo = reactive({
  jenisBarang: 'Emas',
  subJenis: 'Logam Mulia',
  tipeBarang: 'Logam Mulia Antam',
  warnaSpesifikasi: '',
  beratGram: 5,
  beratNetto: 5,
  kadarEmas: '24K',
  kondisiSurat: 'Lengkap Surat / Nota Toko',
  kondisiFisikEmas: 'Mulus Utuh & Segel Toko',
  kondisiElektronik: 'Segel Mulus 99% Normal',
  kelengkapanElektronik: 'Fullset (Box + Charger Original)',
  kelengkapan: 'Kotak & Sertifikat Toko Lengkap',
  kondisi: 'Mulus 99% Utuh & Berfungsi Normal',
  nilaiTaksiran: 6000000,
  maxUp: 5400000
})

// Database Master Offerings & Items
const productOfferings = ref<any[]>([])
const selectedProdukCode = ref('')
const emasDbList = ref<any[]>([])
const selectedEmasDbItem = ref<any>(null)
const elekDbList = ref<any[]>([])
const showElekDbModal = ref(false)
const elekDbSearch = ref('')
const isSearchingElekDb = ref(false)
const barangFotoPreview = ref('')
const barangFotoBase64 = ref('')

const selectJenisGadai = (type: string) => {
  sbgData.jenisTransaksi = type
  fetchProductOfferings()
}

const toggleJenisBarang = (type: string) => {
  barangInfo.jenisBarang = type
  sbgData.jenisBarang = type
  fetchProductOfferings()
  calculateSbgFees()
}

const fetchProductOfferings = async () => {
  try {
    const isEmas = barangInfo.jenisBarang === 'Emas'
    const { data } = await useApi('/taksiran/products', {
      params: { type: isEmas ? 'emas' : 'elektronik' }
    })
    const list = Array.isArray(data) ? data : []
    productOfferings.value = list
    if (list.length > 0 && !selectedProdukCode.value) {
      selectedProdukCode.value = list[0].kd_produk
    }
  } catch (e) {
    console.error('Failed to fetch product offerings:', e)
  }
}

const onProdukChange = () => {
  // Triggers reactivity
}

const fetchEmasDb = async () => {
  try {
    const { data } = await useApi('/taksiran/barang-emas')
    emasDbList.value = Array.isArray(data) ? data : []
    if (emasDbList.value.length > 0 && !selectedEmasDbItem.value) {
      selectedEmasDbItem.value = emasDbList.value[0]
      barangInfo.tipeBarang = emasDbList.value[0].nama_barang
    }
  } catch (e) {
    console.error('Failed to fetch emas items from DB:', e)
  }
}

const onSelectEmasDb = () => {
  if (selectedEmasDbItem.value) {
    barangInfo.tipeBarang = selectedEmasDbItem.value.nama_barang || ''
    calculateSbgFees()
  }
}

const openElekDbSearchModal = () => {
  showElekDbModal.value = true
  if (elekDbList.value.length === 0) {
    searchElekDb()
  }
}

const searchElekDb = async () => {
  isSearchingElekDb.value = true
  try {
    const { data } = await useApi('/taksiran/barang-umum', {
      params: { namaBarang: elekDbSearch.value }
    })
    elekDbList.value = Array.isArray(data) ? data : []
  } catch (e) {
    console.error('Failed to search electronic items in DB:', e)
    elekDbList.value = []
  } finally {
    isSearchingElekDb.value = false
  }
}

const selectedElekItems = ref<any[]>([])

const isElekSelected = (item: any) => {
  const id = item.id_barang_detail || item.id_detail || item.kd_barang
  return selectedElekItems.value.some(i => (i.id_barang_detail || i.id_detail || i.kd_barang) === id)
}

const toggleElekItem = (item: any) => {
  const id = item.id_barang_detail || item.id_detail || item.kd_barang
  const idx = selectedElekItems.value.findIndex(i => (i.id_barang_detail || i.id_detail || i.kd_barang) === id)
  if (idx >= 0) {
    selectedElekItems.value.splice(idx, 1)
  } else {
    selectedElekItems.value.push(item)
  }
  barangInfo.tipeBarang = selectedElekItems.value.map(i => i.nm_barang || i.name).join(', ')
}

const removeElekItem = (index: number) => {
  selectedElekItems.value.splice(index, 1)
  barangInfo.tipeBarang = selectedElekItems.value.map(i => i.nm_barang || i.name).join(', ')
}

const submitHotProspect = () => {
  submitWithStatus('Hot Prospect', true)
}

const onFotoBarangChange = (e: Event) => {
  const target = e.target as HTMLInputElement
  if (target.files && target.files[0]) {
    const file = target.files[0]
    const reader = new FileReader()
    reader.onload = (event) => {
      barangFotoPreview.value = (event.target?.result as string) || ''
      barangFotoBase64.value = barangFotoPreview.value
    }
    reader.readAsDataURL(file)
  }
}

onMounted(() => {
  fetchProductOfferings()
  fetchEmasDb()
})

const form = reactive({
  newStatus: 'Leads',
  tglFollowUp: new Date().toISOString().split('T')[0],
  jamFollowUp: '10:00',
  note: ''
})

const sbgData = reactive({
  noSbg: '',
  jenisTransaksi: 'Gadai',
  jenisBarang: 'Emas',
  deskripsiBarang: '',
  nominalUp: 5000000,
  sewaModalPercent: 1.2,
  sewaModal: 60000,
  biayaAdmin: 15000,
  tglPencairan: new Date().toISOString().split('T')[0],
  tglJatuhTempo: new Date(Date.now() + 15 * 86400000).toISOString().split('T')[0],
  catatanSbg: ''
})

// Tab Switcher for manual inspection
const selectTab = (targetTab: string) => {
  stepWarning.value = ''
  if (targetTab === 'Prospect') {
    activeTab.value = 'Hot Prospect'
  } else {
    activeTab.value = targetTab
  }
}

// Kelurahan Search Select State & Query
const kelurahanSearchInput = ref('')
const showKelurahanDropdown = ref(false)
const { data: kelurahanList, isFetching: isSearchingKelurahan } = useKelurahanQuery(kelurahanSearchInput)

const onSelectKelurahan = (item: any) => {
  customerInfo.kelurahan = item.nm_kelurahan
  customerInfo.kecamatan = item.nm_kecamatan || ''
  customerInfo.kota = item.nm_kota || ''
  kelurahanSearchInput.value = item.nm_kelurahan
  showKelurahanDropdown.value = false
}

const onKelurahanInput = (e: Event) => {
  const val = (e.target as HTMLInputElement).value
  customerInfo.kelurahan = val
  kelurahanSearchInput.value = val
  showKelurahanDropdown.value = val.length >= 2
}

// TanStack Queries
const { data: sourcesQueryData } = useSumberCustQuery()
const { data: inquiryQueryData, isPending: isPendingInquiry } = useInquiryDetailQuery(idRencana)

const isLoading = computed(() => {
  if (!idRencana.value) return false
  return isPendingInquiry.value
})

watch([sourcesQueryData, inquiryQueryData], ([sourcesRes, rawList]) => {
  const sources = Array.isArray(sourcesRes) ? sourcesRes : (sourcesRes?.data || [])
  refSources.value = sources

  if (!rawList || rawList.length === 0) return
  const id = idRencana.value
  const found = rawList.find((i: any) =>
    String(i.id_rencana) === String(id) ||
    String(i.id_customer) === String(id) ||
    String(i.id) === String(id)
  ) || rawList[0]

  if (found) {
    customerInfo.id_customer = found.id_customer || found.id_rencana || found.id || 0
    customerInfo.nama = found.nama || found.nm_customer || found.name || 'Customer'
    customerInfo.hp = found.hp || found.phone || found.no_hp || ''

    const matchedSource = sources.find((s: any) => s.kd_ref_cust === found.id_sumbercust)
    customerInfo.source = matchedSource?.nm_ref_cust || found.nm_ref_cust || found.source || 'Walk In'
    customerInfo.keterangan = found.ket_rencana || found.ket_aktivitas || found.keterangan || ''
    customerInfo.katcust = found.katcust || '1'
    customerInfo.custType = found.cust === 'Baru' ? 'Nasabah Baru' : 'Pelanggan Lama'

    customerInfo.alamat = found.alamat || ''
    customerInfo.rt = found.rt || ''
    customerInfo.rw = found.rw || ''
    customerInfo.kelurahan = found.nm_kelurahan || found.kelurahan || ''
    customerInfo.kecamatan = found.nm_kecamatan || found.kecamatan || ''
    customerInfo.kota = found.nm_kota || found.kota || ''

    const statusFromDb = found.nama_hslaktiv || found.status || 'Leads'
    dbStatus.value = statusFromDb
    form.newStatus = statusFromDb

    const queryStep = route.query.step as string
    if (queryStep && queryStep.toLowerCase().includes('prospect')) {
      activeTab.value = 'Hot Prospect'
    } else if (queryStep && pipelineSteps.includes(queryStep)) {
      activeTab.value = queryStep
    } else {
      const st = statusFromDb.toLowerCase()
      if (st.includes('hot') || st.includes('prospect') || st.includes('prospek')) {
        activeTab.value = 'Hot Prospect'
      } else if (st.includes('sbg') || st.includes('deal')) {
        activeTab.value = 'SBG'
      } else {
        activeTab.value = 'Leads'
      }
    }

    // Build history entries matching customer
    const matches = rawList.filter((i: any) =>
      (i.id_customer && found.id_customer && String(i.id_customer) === String(found.id_customer)) ||
      (i.nama && found.nama && i.nama.toLowerCase() === found.nama.toLowerCase())
    )

    historyList.value = matches.map((m: any) => ({
      tgl: m.tgl_aktivitas || m.tgl_rencana || m.tgl || new Date().toISOString().split('T')[0],
      status: m.nama_hslaktiv || m.status || 'Leads',
      keterangan: m.ket_rencana || m.ket_aktivitas || m.keterangan || '-',
      source: sources.find((s: any) => s.kd_ref_cust === m.id_sumbercust)?.nm_ref_cust || m.nm_ref_cust || 'Walk In'
    }))

    // Init SBG defaults
    sbgData.noSbg = `SBG-${new Date().getFullYear()}${(new Date().getMonth() + 1).toString().padStart(2, '0')}-${Math.floor(100 + Math.random() * 900)}`
    sbgData.deskripsiBarang = customerInfo.keterangan || 'Barang Gadai Customer'
    sbgData.catatanSbg = `Deal SBG untuk ${customerInfo.nama}`
    calculateSbgFees()
  }
}, { immediate: true })

const getStatusBadgeClass = (status: string) => {
  if (!status) return 'bg-amber-100 text-amber-800 border-amber-200'
  const st = status.toLowerCase()
  if (st.includes('lead')) return 'bg-amber-100 text-amber-800 border-amber-200'
  if (st.includes('hot')) return 'bg-rose-100 text-rose-800 border-rose-200'
  if (st.includes('prospect') || st.includes('prospek')) return 'bg-blue-100 text-blue-800 border-blue-200'
  if (st.includes('sbg') || st.includes('deal') || st.includes('do')) return 'bg-emerald-100 text-emerald-800 border-emerald-200'
  if (st.includes('batal')) return 'bg-slate-200 text-slate-600 border-slate-300'
  return 'bg-slate-100 text-slate-600 border-slate-200'
}

const getStepProgressWidth = (status: string) => {
  switch (status) {
    case 'Leads': return '0%'
    case 'Prospect': return '33%'
    case 'Hot Prospect': return '66%'
    case 'SBG': return '100%'
    default: return '0%'
  }
}

const isStepActive = (stepName: string, currentStatus: string) => {
  const currentIdx = pipelineSteps.indexOf(currentStatus)
  const stepIdx = pipelineSteps.indexOf(stepName)
  return stepIdx <= currentIdx
}

const calculateSbgFees = () => {
  if (barangInfo.jenisBarang === 'Emas') {
    const karatVal = barangInfo.kadarEmas === '24K' ? 24 : (barangInfo.kadarEmas === '22K' ? 22 : (barangInfo.kadarEmas === '18K' ? 18 : 16))
    const baseRate = 1200000
    barangInfo.nilaiTaksiran = Math.round((barangInfo.beratGram || 0) * baseRate * (karatVal / 24))
    barangInfo.maxUp = Math.round(barangInfo.nilaiTaksiran * 0.90)
    sbgData.sewaModalPercent = 1.2
  } else {
    barangInfo.maxUp = Math.round((barangInfo.nilaiTaksiran || 0) * 0.85)
    sbgData.sewaModalPercent = 1.5
  }

  const rate = (sbgData.sewaModalPercent || 1.2) / 100
  sbgData.sewaModal = Math.round((sbgData.nominalUp || 0) * rate)
  sbgData.biayaAdmin = (sbgData.nominalUp || 0) > 10000000 ? 35000 : ((sbgData.nominalUp || 0) > 5000000 ? 25000 : 15000)
}

const statusMap: Record<string, number> = {
  'Leads': 1,
  'Prospect': 2,
  'Hot Prospect': 3,
  'Batal': 6,
  'SBG': 7
}

const queryClient = useQueryClient()

const submitWithStatus = (targetStatus: string, navigateNext: boolean = false) => {
  stepWarning.value = ''
  const saveStatus = targetStatus || activeTab.value || dbStatus.value || 'Leads'
  form.newStatus = saveStatus
  const hslAktivVal = statusMap[saveStatus] || 1

  let noteText = form.note ? `[${saveStatus}] ${form.note}` : `Update progress status: ${saveStatus}`
  if (saveStatus === 'Hot Prospect') {
    const prodObj = productOfferings.value.find(p => p.kd_produk === selectedProdukCode.value)
    const prodName = prodObj?.nm_produk || ''
    const itemNames = barangInfo.jenisBarang === 'Emas'
      ? (selectedEmasDbItem.value?.nama_barang || barangInfo.tipeBarang || '-')
      : (selectedElekItems.value.length > 0 ? selectedElekItems.value.map(i => i.nm_barang || i.name).join(', ') : (barangInfo.tipeBarang || '-'))
    noteText += ` | [${sbgData.jenisTransaksi || 'Gadai'}${prodName ? ' - ' + prodName : ''}] ${barangInfo.jenisBarang}: ${itemNames}`
    if (customerInfo.keterangan) {
      noteText += ` (${customerInfo.keterangan})`
    }
  } else if (saveStatus === 'SBG') {
    noteText += ` | [${barangInfo.jenisBarang} - ${barangInfo.tipeBarang}] Nilai Taksiran Rp ${(barangInfo.nilaiTaksiran || 0).toLocaleString('id-ID')} | UP Rp ${(sbgData.nominalUp || 0).toLocaleString('id-ID')} | Sewa Modal Rp ${(sbgData.sewaModal || 0).toLocaleString('id-ID')}`
  }

  const payload: any = {
    id_rencana: idRencana.value ? Number(idRencana.value) : undefined,
    cust: idRencana.value || customerInfo.id_customer ? 'Lama' : 'Baru',
    custlama: customerInfo.id_customer || (idRencana.value ? Number(idRencana.value) : 0),
    nm: customerInfo.nama,
    hp: customerInfo.hp,
    sumbercust: customerInfo.source || 'Walk In',
    hslaktiv: hslAktivVal,
    katcust: customerInfo.katcust,
    ket: noteText,
    tgl_rencana: form.tglFollowUp,
    jam_rencana: form.jamFollowUp,
    alamat: customerInfo.alamat,
    rt: customerInfo.rt,
    rw: customerInfo.rw,
    kelurahan: customerInfo.kelurahan,
    kecamatan: customerInfo.kecamatan,
    kota: customerInfo.kota,
    aktiv: 1
  }

  inputRencanaMutate(payload, {
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inquiries'] })
      queryClient.invalidateQueries({ queryKey: ['dashboard-metrics'] })
      dbStatus.value = saveStatus
      
      if (navigateNext) {
        if (activeTab.value === 'Leads') {
          activeTab.value = 'Prospect'
        } else if (activeTab.value === 'Prospect') {
          activeTab.value = 'Hot Prospect'
        } else if (activeTab.value === 'Hot Prospect') {
          router.push('/inquiry')
        } else {
          router.push('/inquiry')
        }
      } else {
        activeTab.value = saveStatus
      }
    },
    onError: (err: any) => {
      stepWarning.value = err.response?.data?.error || err.data?.error || err.message || 'Gagal menyimpan transaksi'
    }
  })
}
</script>
