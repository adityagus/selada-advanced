<template>
  <div class="space-y-6 pb-12 font-sans mx-auto">
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
          <h1 class="text-2xl font-bold text-slate-800 tracking-tight">Progress Sales
          </h1>
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
        <div class="flex items-center justify-between">
          <div>
            <h2 class="font-semibold text-slate-800 text-base">Wizard Status Pipeline</h2>
            <p class="text-xs text-slate-500">Klik langkah untuk berpindah form input spesifik tiap tahap status.</p>
          </div>
          <span class="text-xs font-semibold px-3 py-1 rounded-full border"
            :class="getStatusBadgeClass(form.newStatus)">
            Tahap Aktif: {{ form.newStatus }}
          </span>
        </div>

        <!-- Progress Bar & Stepper Buttons -->
        <div class="py-2 px-2">
          <div class="flex items-center justify-between relative">
            <div class="absolute left-0 right-0 top-1/2 -translate-y-1/2 h-1 bg-slate-100 -z-0"></div>
            <div class="absolute left-0 top-1/2 -translate-y-1/2 h-1 bg-amber-500 transition-all duration-500 -z-0"
              :style="{ width: getStepProgressWidth(form.newStatus) }"></div>

            <!-- Step 1: Leads -->
            <button type="button" @click="selectStep('Leads')"
              class="relative z-10 flex flex-col items-center gap-1.5 cursor-pointer transition-all">
              <div class="w-10 h-10 rounded-full flex items-center justify-center text-sm font-bold transition-all"
                :class="isStepActive('Leads', form.newStatus) ? 'bg-amber-500 text-white shadow-md shadow-amber-200 ring-4 ring-amber-100' : 'bg-slate-100 text-slate-400'">
                1</div>
              <span class="text-xs font-semibold uppercase tracking-wider"
                :class="form.newStatus === 'Leads' ? 'text-amber-600 font-bold' : 'text-slate-400'">1. Leads</span>
            </button>

            <!-- Step 2: Prospect -->
            <button type="button" @click="selectStep('Prospect')"
              class="relative z-10 flex flex-col items-center gap-1.5 cursor-pointer transition-all">
              <div class="w-10 h-10 rounded-full flex items-center justify-center text-sm font-bold transition-all"
                :class="isStepActive('Prospect', form.newStatus) ? 'bg-blue-600 text-white shadow-md shadow-blue-200 ring-4 ring-blue-100' : 'bg-slate-100 text-slate-400'">
                2</div>
              <span class="text-xs font-semibold uppercase tracking-wider"
                :class="form.newStatus === 'Prospect' ? 'text-blue-600 font-bold' : 'text-slate-400'">2. Prospek</span>
            </button>

            <!-- Step 3: Hot Prospect -->
            <button type="button" @click="selectStep('Hot Prospect')"
              class="relative z-10 flex flex-col items-center gap-1.5 cursor-pointer transition-all">
              <div class="w-10 h-10 rounded-full flex items-center justify-center text-sm font-bold transition-all"
                :class="isStepActive('Hot Prospect', form.newStatus) ? 'bg-rose-600 text-white shadow-md shadow-rose-200 ring-4 ring-rose-100' : 'bg-slate-100 text-slate-400'">
                3</div>
              <span class="text-xs font-semibold uppercase tracking-wider"
                :class="form.newStatus === 'Hot Prospect' ? 'text-rose-600 font-bold' : 'text-slate-400'">3. Hot
                Prospek</span>
            </button>

            <!-- Step 4: SBG / Deal -->
            <button type="button" @click="selectStep('SBG')"
              class="relative z-10 flex flex-col items-center gap-1.5 cursor-pointer transition-all">
              <div class="w-10 h-10 rounded-full flex items-center justify-center text-xs font-bold transition-all"
                :class="isStepActive('SBG', form.newStatus) ? 'bg-emerald-600 text-white shadow-md shadow-emerald-200 ring-4 ring-emerald-100' : 'bg-slate-100 text-slate-400'">
                4</div>
              <span class="text-xs font-semibold uppercase tracking-wider"
                :class="form.newStatus === 'SBG' ? 'text-emerald-600 font-bold' : 'text-slate-400'">4. SBG (Deal)</span>
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

      <!-- RIWAYAT FOLLOW-UP & LOG AKTIVITAS (SLIDE-IN HISTORY) -->
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
      <!-- PANEL 1: LEADS (DILENGKAPI DATA KONTAK AWAL) -->
      <!-- ========================================================================= -->
      <div v-if="form.newStatus === 'Leads'"
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

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
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
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1">Catatan Follow Up
            Leads</label>
          <textarea v-model="form.note" rows="2" placeholder="Catatan kontak awal..."
            class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm"></textarea>
        </div>

        <div class="flex justify-end gap-3 pt-3 border-t border-slate-100">
          <button @click="selectStep('Prospect')" type="button"
            class="px-5 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-semibold flex items-center gap-1.5 shadow-sm">
            <span>Lanjut ke Step 2 (Prospek)</span>
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
            </svg>
          </button>
        </div>
      </div>

      <!-- ========================================================================= -->
      <!-- PANEL 2: PROSPEK (DATA ALAMAT & SERVE CONSULTATION) -->
      <!-- ========================================================================= -->
      <div v-else-if="form.newStatus === 'Prospect'"
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

        <div class="flex justify-between items-center pt-3 border-t border-slate-100">
          <button @click="selectStep('Leads')" type="button"
            class="px-4 py-2 border flex items-center gap-2 border-slate-300 text-slate-600 rounded-xl text-xs font-medium">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
              stroke="currentColor" class="size-3">
              <path stroke-linecap="round" stroke-linejoin="round" d="M10.5 19.5 3 12m0 0 7.5-7.5M3 12h18" />
            </svg>
            Kembali ke Leads
          </button>

          <button @click="selectStep('Hot Prospect')" type="button"
            class="px-5 py-2.5 bg-rose-600 hover:bg-rose-700 text-white rounded-xl text-xs font-semibold flex items-center gap-1.5 shadow-sm">
            <span>Lanjut ke Step 3 (Hot Prospek)</span>
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
            </svg>
          </button>
        </div>
      </div>

      <!-- ========================================================================= -->
      <!-- PANEL 3: HOT PROSPEK (KARTU GADAI/CICILAN & SPESIFIKASI BARANG) -->
      <!-- ========================================================================= -->
      <div v-else-if="form.newStatus === 'Hot Prospect'"
        class="bg-white rounded-2xl border border-rose-200 p-6 shadow-xs space-y-5">
        <div class="flex items-center justify-between border-b border-slate-100 pb-3">
          <div class="flex items-center gap-2">
            <span
              class="w-7 h-7 rounded-full bg-rose-600 text-white font-bold text-xs flex items-center justify-center">3</span>
            <h2 class="font-bold text-slate-800 text-base">PANEL STEP 3 — HOT PROSPEK (Kartu Transaksi & Inventori)</h2>
          </div>
          <span class="text-xs font-medium text-rose-800 bg-rose-50 px-3 py-1 rounded-lg border border-rose-200">Status:
            Hot Prospek</span>
        </div>

        <div
          class="p-3.5 bg-rose-50/60 border border-rose-200 rounded-xl text-xs text-rose-900 flex items-center gap-2">
          <svg class="w-4 h-4 text-rose-600 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M17.657 18.657A8 8 0 016.343 7.343S7 9 9 10c0-2 .5-5 2.986-7C14 5 16.09 5.777 17.656 7.343A7.975 7.975 0 0120 13a7.975 7.975 0 01-2.343 5.657z" />
          </svg>
          <span>Pilih jenis skema transaksi (Gadai vs Cicilan) dan spesifikasi detail barang jaminan customer.</span>
        </div>

        <!-- Kartu Skema Transaksi (Gadai vs Cicil - Parity v_updateinq) -->
        <div>
          <label class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-2">Pilih Skema Transaksi
            <span class="text-rose-500">*</span></label>
          <div class="grid grid-cols-2 gap-4 max-w-sm">
            <div @click="sbgData.jenisTransaksi = 'Gadai'"
              class="p-4 rounded-xl border-2 text-center cursor-pointer transition-all"
              :class="sbgData.jenisTransaksi === 'Gadai' ? 'border-amber-500 bg-amber-50/60 shadow-xs' : 'border-slate-200 hover:border-slate-300 bg-white'">
              <div class="text-2xl mb-1">💎</div>
              <div class="font-bold text-sm text-slate-800">Gadai</div>
              <div class="text-xs text-slate-500">Pencairan Reguler</div>
            </div>

            <div @click="sbgData.jenisTransaksi = 'Cicilan'"
              class="p-4 rounded-xl border-2 text-center cursor-pointer transition-all"
              :class="sbgData.jenisTransaksi === 'Cicilan' ? 'border-amber-500 bg-amber-50/60 shadow-xs' : 'border-slate-200 hover:border-slate-300 bg-white'">
              <div class="text-2xl mb-1">💳</div>
              <div class="font-bold text-sm text-slate-800">Cicilan</div>
              <div class="text-xs text-slate-500">Angsuran Berkala</div>
            </div>
          </div>
        </div>

        <!-- Spesifikasi Barang Gadai -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 pt-2">
          <div>
            <label class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1">Jenis Barang
              Jaminan</label>
            <select v-model="barangInfo.jenisBarang"
              class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm">
              <option value="Emas">Emas / Perhiasan / Logam Mulia</option>
              <option value="Kendaraan">Kendaraan (Motor / Mobil)</option>
              <option value="Elektronik">Elektronik (HP / Laptop / TV)</option>
              <option value="Lainnya">Lain-lain</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1">Tipe / Merk
              Barang</label>
            <input v-model="barangInfo.tipeBarang" type="text" placeholder="Contoh: Honda Beat 2022 / iPhone 13 128GB"
              class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm" />
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1">Warna / Kadar
              Emas</label>
            <input v-model="barangInfo.warnaSpesifikasi" type="text" placeholder="Contoh: Hitam / Karat 24K (5 gram)"
              class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm" />
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1">Nominal UP Pinjaman
              (Rp) <span class="text-rose-500">*</span></label>
            <input v-model.number="sbgData.nominalUp" @input="calculateSbgFees" type="number"
              class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm font-bold text-amber-700" />
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">Sewa Modal (15
              Hari)</label>
            <input v-model.number="sbgData.sewaModal" type="number"
              class="w-full px-3.5 py-2.5 bg-slate-100 border border-slate-200 rounded-xl text-sm font-semibold"
              readonly />
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">Biaya
              Administrasi</label>
            <input v-model.number="sbgData.biayaAdmin" type="number"
              class="w-full px-3.5 py-2.5 bg-slate-100 border border-slate-200 rounded-xl text-sm font-semibold"
              readonly />
          </div>
        </div>

        <div class="flex justify-between items-center pt-3 border-t border-slate-100">
          <button @click="selectStep('Prospect')" type="button"
            class="px-4 py-2 border border-slate-300 text-slate-600 rounded-xl text-xs font-medium">
            &larr; Kembali ke Prospek
          </button>

          <button @click="selectStep('SBG')" type="button"
            class="px-5 py-2.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-xs font-semibold flex items-center gap-1.5 shadow-sm">
            <span>Lanjut ke Step 4 (SBG / Deal)</span>
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
            </svg>
          </button>
        </div>
      </div>

      <!-- ========================================================================= -->
      <!-- PANEL 4: SBG / DEAL (FORM REALISASI AKAD & TRANSAKSI) -->
      <!-- ========================================================================= -->
      <div v-else-if="form.newStatus === 'SBG'"
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

        <div class="flex justify-start pt-3 border-t border-slate-100">
          <button @click="selectStep('Hot Prospect')" type="button"
            class="px-4 py-2 border border-slate-300 text-slate-600 rounded-xl text-xs font-medium">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
              stroke="currentColor" class="size-11">
              <path stroke-linecap="round" stroke-linejoin="round" d="M10.5 19.5 3 12m0 0 7.5-7.5M3 12h18" />
            </svg>
            Kembali ke Hot Prospek
          </button>
        </div>
      </div>

      <!-- ACTION BUTTONS BAR -->
      <div
        class="flex flex-col sm:flex-row justify-between items-center bg-white p-6 rounded-2xl border border-slate-200/80 shadow-xs gap-4">
        <NuxtLink to="/inquiry"
          class="px-5 py-2.5 border border-slate-300 text-slate-600 hover:bg-slate-50 rounded-xl text-sm font-medium transition-all w-full sm:w-auto text-center">
          Batal & Kembali
        </NuxtLink>

        <button @click="submitAllProgress" type="button"
          class="px-7 py-2.5 bg-amber-500 hover:bg-amber-600 text-white font-bold rounded-xl text-sm shadow-md shadow-amber-200 flex items-center justify-center gap-2 cursor-pointer transition-all w-full sm:w-auto"
          :disabled="isSubmitting">
          <span v-if="isSubmitting"
            class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
          <span v-else>Simpan Transaksi Status {{ form.newStatus }}</span>
        </button>
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
  tipeBarang: '',
  warnaSpesifikasi: '',
  kelengkapan: [] as string[]
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
  sewaModal: 60000,
  biayaAdmin: 15000,
  tglPencairan: new Date().toISOString().split('T')[0],
  tglJatuhTempo: new Date(Date.now() + 15 * 86400000).toISOString().split('T')[0],
  catatanSbg: ''
})

const getNextStepName = computed(() => {
  const currentIdx = pipelineSteps.indexOf(form.newStatus)
  if (currentIdx >= 0 && currentIdx < pipelineSteps.length - 1) {
    return pipelineSteps[currentIdx + 1]
  }
  return null
})

const selectStep = (targetStatus: string) => {
  stepWarning.value = ''
  if (targetStatus === 'Batal') {
    form.newStatus = 'Batal'
    return
  }
  form.newStatus = targetStatus
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

    const queryStep = route.query.step as string
    if (queryStep && pipelineSteps.includes(queryStep)) {
      form.newStatus = queryStep
    } else {
      form.newStatus = found.nama_hslaktiv || found.status || 'Leads'
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

const cleanPhone = (phoneStr: string) => {
  if (!phoneStr) return ''
  let clean = phoneStr.replace(/\D/g, '')
  if (clean.startsWith('0')) {
    clean = '62' + clean.slice(1)
  }
  return clean
}

const getStatusBadgeClass = (status: string) => {
  switch (status) {
    case 'Leads': return 'bg-amber-100 text-amber-800 border-amber-200'
    case 'Prospect': return 'bg-blue-100 text-blue-800 border-blue-200'
    case 'Hot Prospect': return 'bg-rose-100 text-rose-800 border-rose-200'
    case 'SBG': return 'bg-emerald-100 text-emerald-800 border-emerald-200'
    case 'Batal': return 'bg-slate-200 text-slate-600 border-slate-300'
    default: return 'bg-slate-100 text-slate-600 border-slate-200'
  }
}

const isStepAllowed = (targetStatus: string) => true

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
  const rate = sbgData.jenisBarang === 'Emas' ? 0.015 : 0.012
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

const submitAllProgress = () => {
  stepWarning.value = ''
  const hslAktivVal = statusMap[form.newStatus] || 1
  let noteText = form.note ? `[${form.newStatus}] ${form.note}` : `Update progress status: ${form.newStatus}`

  if (form.newStatus === 'SBG' || form.newStatus === 'Hot Prospect') {
    noteText += ` | [SBG DEAL - ${sbgData.noSbg}] ${sbgData.deskripsiBarang} | UP Rp ${sbgData.nominalUp.toLocaleString('id-ID')} | Tgl ${sbgData.tglPencairan}`
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
      router.push('/inquiry')
    },
    onError: (err: any) => {
      stepWarning.value = err.response?.data?.error || err.data?.error || err.message || 'Gagal menyimpan transaksi'
    }
  })
}
</script>
