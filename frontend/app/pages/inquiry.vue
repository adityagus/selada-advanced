<template>
  <div class="space-y-6 pb-12 font-sans">
    <!-- Header Section -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-semibold text-slate-800 tracking-tight">Inquiry & Rencana Aktivitas</h1>
        <p class="text-slate-500 text-sm">Kelola pipeline sales, update status prospek, dan catat janji temu customer.</p>
      </div>
      <div class="flex items-center gap-3">
        <button @click="openModal = true" class="inline-flex items-center gap-2 bg-amber-500 hover:bg-amber-600 active:scale-95 text-white font-medium px-5 py-2.5 rounded-xl text-sm transition-all shadow-md shadow-amber-200 cursor-pointer">
          <AppIcon name="plus" class="w-4 h-4" />
          Add Leads Baru
        </button>
      </div>
    </div>

    <!-- Summary Metrics Bar -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
      <div @click="selectedStatus = 'Leads'; fetchInquiries()" class="bg-white p-4 rounded-xl border border-slate-200/80 shadow-xs hover:border-amber-400 transition-all cursor-pointer" :class="{ 'ring-2 ring-amber-400 bg-amber-50/20': selectedStatus === 'Leads' }">
        <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">Leads</span>
        <div class="text-2xl font-bold text-amber-600 mt-1">{{ summaryCounts.leads }}</div>
      </div>
      <div @click="selectedStatus = 'Prospect'; fetchInquiries()" class="bg-white p-4 rounded-xl border border-slate-200/80 shadow-xs hover:border-blue-400 transition-all cursor-pointer" :class="{ 'ring-2 ring-blue-400 bg-blue-50/20': selectedStatus === 'Prospect' }">
        <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">Prospect</span>
        <div class="text-2xl font-bold text-blue-600 mt-1">{{ summaryCounts.prospect }}</div>
      </div>
      <div @click="selectedStatus = 'Hot Prospect'; fetchInquiries()" class="bg-white p-4 rounded-xl border border-slate-200/80 shadow-xs hover:border-rose-400 transition-all cursor-pointer" :class="{ 'ring-2 ring-rose-400 bg-rose-50/20': selectedStatus === 'Hot Prospect' }">
        <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">Hot Prospect</span>
        <div class="text-2xl font-bold text-rose-600 mt-1">{{ summaryCounts.hotProspect }}</div>
      </div>
      <div @click="selectedStatus = 'all'; fetchInquiries()" class="bg-white p-4 rounded-xl border border-slate-200/80 shadow-xs hover:border-emerald-400 transition-all cursor-pointer" :class="{ 'ring-2 ring-emerald-400 bg-emerald-50/20': selectedStatus === 'all' }">
        <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">Semua Active</span>
        <div class="text-2xl font-bold text-emerald-600 mt-1">{{ summaryCounts.leads + summaryCounts.prospect + summaryCounts.hotProspect }}</div>
      </div>
    </div>

    <!-- Search & Filter Card -->
    <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-xs flex flex-col md:flex-row gap-4 justify-between items-end">
      <div class="flex-1 w-full md:w-auto">
        <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Cari Customer / Kontak</label>
        <div class="relative">
          <input
            v-model="searchQuery"
            @input="debounceFetch"
            type="text"
            placeholder="Ketik nama atau nomor handphone..."
            class="w-full pl-10 pr-4 py-2.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-800 focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all"
          />
          <AppIcon name="search" class="w-4 h-4 text-slate-400 absolute left-3.5 top-3.5" />
        </div>
      </div>

      <div class="w-full md:w-56">
        <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Filter Status</label>
        <select v-model="selectedStatus" @change="fetchInquiries" class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-800 focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all">
          <option value="all">Semua Status Aktif</option>
          <option value="Leads">Leads</option>
          <option value="Prospect">Prospect</option>
          <option value="Hot Prospect">Hot Prospect</option>
        </select>
      </div>

      <div class="w-full md:w-48">
        <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Sumber Leads</label>
        <select v-model="selectedSource" @change="fetchInquiries" class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all">
          <option value="all">Semua Sumber</option>
          <option v-for="s in sourcesList" :key="s" :value="s">{{ s }}</option>
        </select>
      </div>
    </div>

    <!-- Data Table Card -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-sm">
          <thead>
            <tr class="bg-slate-50/80 border-b border-slate-200 text-xs font-semibold text-slate-500 uppercase tracking-wider">
              <th class="px-5 py-3.5">Tanggal</th>
              <th class="px-5 py-3.5">Nama Customer</th>
              <th class="px-5 py-3.5 text-center">Sumber</th>
              <th class="px-5 py-3.5 text-center">Status</th>
              <th class="px-5 py-3.5">Aktivitas Terakhir</th>
              <th class="px-5 py-3.5">Keterangan Barang</th>
              <th class="px-5 py-3.5 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-slate-700">
            <tr v-if="isLoadingTable">
              <td colspan="7" class="text-center py-8">
                <AppIcon name="loader" class="w-6 h-6  text-amber-500   animate-spin mx-auto" />
              </td>
            </tr>
            <tr v-else-if="inquiries.length === 0">
              <td colspan="7" class="text-center py-8 text-slate-400 font-medium">
                Tidak ada data inquiry ditemukan.
              </td>
            </tr>
            <!-- <pre>{{ inquiries }}</pre> -->
            <tr v-for="item in inquiries" :key="item.id_rencana" v-show="!isLoadingTable" class="hover:bg-amber-50/30 transition-all">
              <td class="px-5 py-3.5 whitespace-nowrap text-slate-500">{{ formatDate(item.tgl_rencana) }}</td>
              <td class="px-5 py-3.5 font-semibold text-slate-800">
                {{ item.nama }}
                <!-- <span v-if="item.hp" class="block text-xs font-normal text-slate-400">{{ item.hp }}</span> -->
              </td>
              <td class="px-5 py-3.5">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-slate-100 text-slate-600">
                  {{ item.source }}
                </span>
              </td>
              <td class="px-5 py-3.5 text-center">
                <span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold uppercase tracking-wider" :class="getStatusBadgeClass(item.status)">
                  {{ item.status }}
                </span>
              </td>
              <td class="px-5 py-3.5">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold uppercase tracking-wide">
                  {{ item.aktivitas }}
                </span>
              </td>
              <td class="px-5 py-3.5">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold uppercase tracking-wide">
                  {{ item.keterangan_barang || '-' }}
                </span>
              </td>
              <td class="px-5 py-3.5 text-center">
                <div class="inline-flex items-center gap-1.5" >
                  <NuxtLink :to="`/input-sbg?id=${item.id_rencana || item.id_customer}`" class="p-1.5 text-emerald-600 hover:bg-emerald-50 rounded-lg border border-emerald-200 transition-all" title="Submit SBG" v-if="item.id_hslaktiv == 3">
                    <AppIcon name="check-circle" class="w-4 h-4" />
                  </NuxtLink>
                  <NuxtLink :to="`/input-rencana?id=${item.id_rencana || item.id_customer}`" class="p-1.5 text-amber-600 hover:bg-amber-50 rounded-lg border border-amber-200 transition-all" title="Buka Halaman Progress & Form SBG" v-if="item.id_hslaktiv != 3">
                    <AppIcon name="edit" class="w-4 h-4" />
                  </NuxtLink>
                  <button @click="showDetail(item)" class="p-1.5 text-slate-600 hover:bg-slate-100 rounded-lg border border-slate-200 transition-all" title="Detail">
                    <AppIcon name="eye" class="w-4 h-4" />
                  </button>
                  <button @click="handleDelete(item.id_rencana)" class="p-1.5 text-rose-600 hover:bg-rose-50 rounded-lg border border-rose-200 transition-all" title="Hapus">
                    <AppIcon name="trash" class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- MODAL: BUAT RENCANA LEADS BARU -->
    <ClientOnly>
      <Teleport to="body">
        <div v-if="openModal" class="fixed inset-0 w-screen h-screen min-h-screen top-0 left-0 right-0 bottom-0 z-[9999] bg-slate-900/60 backdrop-blur-md flex items-center justify-center p-4 !mt-0 overflow-y-auto">
          <div class="bg-white rounded-2xl max-w-lg w-full p-6 shadow-2xl space-y-4 max-h-[90vh] overflow-y-auto my-auto">
            <div class="flex items-center justify-between border-b border-slate-100 pb-3">
              <h3 class="font-semibold text-slate-800 text-base">Buat Rencana Leads Baru</h3>
              <button @click="closeFormModal" class="text-slate-400 hover:text-slate-600 p-1 cursor-pointer"><AppIcon name="close" class="w-4 h-4" /></button>
            </div>

            <form @submit.prevent="handleSubmitRencana" class="space-y-4">
              <div>
                <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Tipe Customer</label>
                <div class="grid grid-cols-2 gap-3">
                  <label class="flex items-center gap-2 p-3 border rounded-xl cursor-pointer transition-all" :class="form.cust === 'Baru' ? 'border-amber-500 bg-amber-50/50 text-amber-900 font-semibold' : 'border-slate-200 text-slate-600'">
                    <input type="radio" v-model="form.cust" value="Baru" class="text-amber-500 focus:ring-amber-400" />
                    <span class="text-xs">Customer Baru</span>
                  </label>
                  <label class="flex items-center gap-2 p-3 border rounded-xl cursor-pointer transition-all" :class="form.cust === 'Lama' ? 'border-amber-500 bg-amber-50/50 text-amber-900 font-semibold' : 'border-slate-200 text-slate-600'">
                    <input type="radio" v-model="form.cust" value="Lama" class="text-amber-500 focus:ring-amber-400" />
                    <span class="text-xs">Pelanggan Lama</span>
                  </label>
                </div>
              </div>

              <template v-if="form.cust === 'Baru'">
                <div>
                  <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">Nama Lengkap</label>
                  <input v-model="form.nm" type="text" placeholder="Masukkan nama lengkap customer" class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm" required />
                </div>
                <div>
                  <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">No. Handphone / WA</label>
                  <input v-model="form.hp" type="text" placeholder="08xxxxxxxxxx" class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm" required />
                </div>
              </template>

              <template v-else>
                <div>
                  <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">Pilih Customer Existing</label>
                  <select v-model="form.custlama" class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm" required>
                    <option value="">- Pilih Customer -</option>
                    <option v-for="c in existingCustomers" :key="c.id_customer" :value="c.id_customer">
                      {{ c.nama }} ({{ c.hp }})
                    </option>
                  </select>
                </div>
              </template>
              <!-- 3 Inputan Vertical: SOA, Aktifitas, Jenis Aplikasi -->
              <div class="space-y-4 pt-1">
                <!-- Input 1: SOA -->
                <div>
                  <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">
                    SOA <span class="text-rose-500">*</span>
                  </label>
                  <select v-model="form.soa" @change="onSoaChange" class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm font-medium text-slate-800 focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all" required>
                    <option value="">- Pilih SOA -</option>
                    <option v-for="soaOpt in soaList" :key="soaOpt" :value="soaOpt">{{ soaOpt }}</option>
                  </select>
                </div>

                <!-- Input 2: Aktifitas -->
                <div>
                  <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">
                    Aktifitas <span class="text-rose-500">*</span>
                  </label>
                  <select v-model="form.aktifitas" :disabled="!form.soa" class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm font-medium text-slate-800 focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all disabled:bg-slate-100 disabled:cursor-not-allowed" required>
                    <option value="">- Pilih Aktifitas -</option>
                    <option v-for="aktOpt in availableAktifitasOptions" :key="aktOpt" :value="aktOpt">{{ aktOpt }}</option>
                  </select>
                </div>

                <!-- Input 3: Jenis Aplikasi -->
                <div>
                  <label class="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">
                    Jenis Aplikasi <span class="text-rose-500">*</span>
                  </label>
                  <select v-model="form.jenisAplikasi" :disabled="!form.soa" class="w-full px-3.5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm font-medium text-slate-800 focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all disabled:bg-slate-100 disabled:cursor-not-allowed" required>
                    <option value="">- Pilih Jenis Aplikasi -</option>
                    <option v-for="jenisOpt in availableJenisAplikasiOptions" :key="jenisOpt" :value="jenisOpt">{{ jenisOpt }}</option>
                  </select>
                </div>
              </div>


              <div v-if="modalError" class="p-3 bg-rose-50 border border-rose-200 rounded-xl text-xs text-rose-600">
                {{ modalError }}
              </div>

              <div class="flex justify-end gap-3 pt-2">
                <button type="button" @click="closeFormModal" class="px-4 py-2 border border-slate-300 text-slate-600 rounded-xl text-sm font-medium cursor-pointer">Batal</button>
                <button type="submit" class="px-5 py-2 bg-amber-500 hover:bg-amber-600 text-white rounded-xl text-sm font-medium shadow-sm cursor-pointer" :disabled="isSubmitting">
                  <span v-if="isSubmitting">Menyimpan...</span>
                  <span v-else>Simpan Rencana</span>
                </button>
              </div>
            </form>
          </div>
        </div>
      </Teleport>
    </ClientOnly>

    <!-- MODAL: DETAIL INQUIRY -->
    <ClientOnly>
      <Teleport to="body">
        <div
          v-if="selectedItem"
          class="fixed inset-0 w-screen h-screen min-h-screen top-0 left-0 right-0 bottom-0 z-[9999] bg-slate-900/60 backdrop-blur-xs flex items-center justify-center p-4 !mt-0 overflow-y-auto"
          @click.self="closeDetailModal"
        >
          <div
            class="bg-white rounded-2xl max-w-2xl w-full shadow-2xl border border-slate-200 flex flex-col max-h-[90vh] overflow-hidden my-auto"
            @click.stop
          >
            <!-- MODAL HEADER -->
            <div class="px-6 py-4 border-b border-slate-100 flex items-center justify-between">
              <div>
                <h3 class="font-semibold text-slate-800 text-base">Detail Aktivitas</h3>
                <span class="text-xs text-slate-400">ID Rencana: #{{ selectedItem.id_rencana || selectedItem.id_customer }}</span>
              </div>
              <button
                type="button"
                @click="closeDetailModal"
                class="text-slate-400 hover:text-slate-600 p-1.5 rounded-lg hover:bg-slate-100 transition-all cursor-pointer"
                title="Tutup (ESC)"
              >
                <AppIcon name="close" class="w-4 h-4" />
              </button>
            </div>

            <!-- MODAL BODY -->
            <div class="p-6 overflow-y-auto space-y-5 text-xs sm:text-sm">
              <!-- INFORMASI UTAMA & STATUS -->
              <div class="grid grid-cols-2 sm:grid-cols-3 gap-4 border-b border-slate-100 pb-5">
                <div>
                  <span class="text-xs text-slate-400 block mb-0.5">Nama Customer</span>
                  <span class="font-semibold text-slate-800">{{ selectedItem.nama }}</span>
                </div>
                <div>
                  <span class="text-xs text-slate-400 block mb-0.5">No. Handphone</span>
                  <span class="font-medium text-slate-700">{{ selectedItem.hp || '-' }}</span>
                </div>
                <div>
                  <span class="text-xs text-slate-400 block mb-0.5">Status</span>
                  <span class="inline-flex items-center px-2.5 py-0.5 rounded-md text-xs font-semibold" :class="getStatusBadgeClass(selectedItem.status)">
                    {{ selectedItem.status }}
                  </span>
                </div>
                <div>
                  <span class="text-xs text-slate-400 block mb-0.5">Sumber Lead</span>
                  <span class="text-slate-700">{{ selectedItem.source || 'Walk In' }}</span>
                </div>
                <div>
                  <span class="text-xs text-slate-400 block mb-0.5">Tanggal Rencana</span>
                  <span class="text-slate-700">{{ formatDate(selectedItem.tgl_rencana) }}</span>
                </div>
                <div>
                  <span class="text-xs text-slate-400 block mb-0.5">Jam Rencana</span>
                  <span class="text-slate-700">{{ selectedItem.jam_rencana ? selectedItem.jam_rencana + ' WIB' : '-' }}</span>
                </div>
              </div>

              <!-- SOA, AKTIVITAS & DOMISILI -->
              <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 border-b border-slate-100 pb-5">
                <div>
                  <span class="text-xs text-slate-400 block mb-0.5">SOA</span>
                  <span class="font-medium text-slate-800">{{ selectedItem.soa || selectedItem.SOA || '-' }}</span>
                </div>
                <div>
                  <span class="text-xs text-slate-400 block mb-0.5">Aktifitas</span>
                  <span class="font-medium text-slate-800">{{ selectedItem.aktifitas || selectedItem.ket_aktivitas || '-' }}</span>
                </div>
                <div>
                  <span class="text-xs text-slate-400 block mb-0.5">Jenis Aplikasi</span>
                  <span class="font-medium text-slate-800">{{ selectedItem.jenis_aplikasi || selectedItem.jenisAplikasi || '-' }}</span>
                </div>
                <div class="sm:col-span-2">
                  <span class="text-xs text-slate-400 block mb-0.5">Alamat</span>
                  <span class="text-slate-700">{{ selectedItem.alamat || '-' }}</span>
                </div>
                <div>
                  <span class="text-xs text-slate-400 block mb-0.5">RT / RW</span>
                  <span class="text-slate-700">{{ selectedItem.rt || '-' }} / {{ selectedItem.rw || '-' }}</span>
                </div>
                <div>
                  <span class="text-xs text-slate-400 block mb-0.5">Kelurahan</span>
                  <span class="text-slate-700">{{ selectedItem.kelurahan || '-' }}</span>
                </div>
                <div>
                  <span class="text-xs text-slate-400 block mb-0.5">Kecamatan</span>
                  <span class="text-slate-700">{{ selectedItem.kecamatan || '-' }}</span>
                </div>
                <div>
                  <span class="text-xs text-slate-400 block mb-0.5">Kota</span>
                  <span class="text-slate-700">{{ selectedItem.kota || '-' }}</span>
                </div>
              </div>

              <!-- KETERANGAN BARANG & CATATAN -->
              <div class="space-y-4 border-b border-slate-100 pb-5">
                <div>
                  <span class="text-xs text-slate-400 block mb-1">Keterangan Barang Jaminan</span>
                  <div class="p-3 bg-slate-50 border border-slate-200 rounded-xl text-slate-800 text-xs">
                    {{ selectedItem.keterangan_barang || '-' }}
                  </div>
                </div>

                <div>
                  <span class="text-xs text-slate-400 block mb-1">Catatan Keterangan Progres</span>
                  <div class="p-3 bg-slate-50 border border-slate-200 rounded-xl text-slate-700 text-xs leading-relaxed whitespace-pre-wrap">
                    {{ selectedItem.keterangan || selectedItem.ket || '-' }}
                  </div>
                </div>
              </div>

              <!-- FOTO BARANG JAMINAN -->
              <div>
                <span class="text-xs text-slate-400 block mb-1.5">Foto Barang Jaminan</span>
                <div v-if="detailFotoUrl" class="inline-block">
                  <img
                    :src="detailFotoUrl"
                    alt="Foto Barang Jaminan"
                    @click="openLightbox"
                    class="max-h-56 rounded-xl border border-slate-200 object-contain cursor-pointer hover:opacity-90 transition-opacity bg-slate-50 p-1"
                    title="Klik untuk memperbesar foto"
                  />
                  <span class="block text-[11px] text-slate-400 mt-1 cursor-pointer" @click="openLightbox">
                    Klik foto untuk melihat ukuran penuh
                  </span>
                </div>
                <div v-else class="text-xs text-slate-400 italic">
                  Tidak ada foto terlampir.
                </div>
              </div>
            </div>

            <!-- MODAL FOOTER -->
            <div class="px-6 py-3.5 border-t border-slate-100 bg-slate-50/50 flex justify-end">
              <button
                type="button"
                @click="closeDetailModal"
                class="px-5 py-2 border border-slate-300 hover:bg-slate-100 text-slate-700 rounded-xl text-xs font-semibold transition-all cursor-pointer"
              >
                Tutup
              </button>
            </div>
          </div>
        </div>
      </Teleport>
    </ClientOnly>

    <!-- FULLSCREEN IMAGE PREVIEW -->
    <ClientOnly>
      <Teleport to="body">
        <div
          v-if="showLightbox && detailFotoUrl"
          class="fixed inset-0 w-screen h-screen min-h-screen top-0 left-0 right-0 bottom-0 z-[10000] bg-slate-950/85 backdrop-blur-xs flex flex-col justify-between p-4 sm:p-6 select-none"
          @click.self="closeLightbox"
        >
          <div class="flex items-center justify-between text-white max-w-4xl mx-auto w-full">
            <span class="text-sm font-medium text-slate-200 truncate">
              Foto Barang: {{ selectedItem?.nama }}
            </span>
            <button
              type="button"
              @click="closeLightbox"
              class="p-1.5 text-slate-300 hover:text-white rounded-lg hover:bg-white/10 transition-all cursor-pointer"
              title="Tutup (ESC)"
            >
              <AppIcon name="close" class="w-5 h-5" />
            </button>
          </div>

          <div class="flex-1 flex items-center justify-center p-2" @click="closeLightbox">
            <img
              :src="detailFotoUrl"
              alt="Foto Barang Jaminan"
              class="max-w-[90vw] max-h-[82vh] object-contain rounded-xl shadow-2xl"
              @click.stop
            />
          </div>

          <div class="text-center text-xs text-slate-400">
            Klik di luar gambar atau tekan ESC untuk menutup
          </div>
        </div>
      </Teleport>
    </ClientOnly>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useApi } from '~/composables/useApi'
import { getRencanaDetails } from '~/utils/rencanaDetails'
import { SOA_LIST, getAktifitasOptions, getJenisAplikasiOptions } from '~/utils/soa'
import { getFotoBarang } from '~/utils/imageStorage'

const inquiries = ref<any[]>([])
const searchQuery = ref('')
const selectedStatus = ref('all')
const selectedSource = ref('all')
const isLoadingTable = ref(false)
const openModal = ref(false)
const isSubmitting = ref(false)
const modalError = ref('')
const existingCustomers = ref<any[]>([])
const selectedItem = ref<any>(null)

// Detail Modal & Photo State
const detailFotoUrl = ref<string | null>(null)
const showLightbox = ref(false)

const summaryCounts = reactive({
  leads: 0,
  prospect: 0,
  hotProspect: 0,
  sbg: 0
})

const form = reactive({
  cust: 'Baru',
  custlama: '',
  nm: '',
  hp: '',
  sumbercust: '',
  hslaktiv: 1,
  katcust: '1',
  ket: '',
  aktiv: 1,
  soa: '',
  aktifitas: '',
  jenisAplikasi: ''
})

const soaList = SOA_LIST

const availableAktifitasOptions = computed(() => getAktifitasOptions(form.soa))
const availableJenisAplikasiOptions = computed(() => getJenisAplikasiOptions(form.soa))

const onSoaChange = () => {
  form.aktifitas = ''
  form.jenisAplikasi = ''
}

const sourcesList = ref<string[]>(['Walk In', 'Media Sosial', 'Event', 'Referensi'])

onMounted(() => {
  fetchSummaryCounts()
  fetchInquiries()
  fetchExistingCustomers()
  fetchSources()
  if (typeof window !== 'undefined') {
    window.addEventListener('keydown', handleKeydown)
  }
})

onUnmounted(() => {
  if (typeof window !== 'undefined') {
    window.removeEventListener('keydown', handleKeydown)
  }
})

const fetchSummaryCounts = async () => {
  try {
    const { data } = await useApi('/inquiry/counts')
    if (data) {
      summaryCounts.leads = data.leads || 0
      summaryCounts.prospect = data.prospect || 0
      summaryCounts.hotProspect = data.hotprospect || 0
      summaryCounts.sbg = data.sbg || 0
    }
  } catch (e) {
    console.error('Failed to fetch inquiry counts:', e)
  }
}

const fetchSources = async () => {
  const { data } = await useApi('/sumbercust')
  if (data && Array.isArray(data) && data.length > 0) {
    sourcesList.value = data.map((s: any) => s.nm_ref_cust || s.kd_ref_cust)
    if (sourcesList.value.length > 0) {
      form.sumbercust = sourcesList.value[0]
    }
  }
}

const fetchInquiries = async () => {
  isLoadingTable.value = true
  const statusParam = selectedStatus.value === 'all' ? '' : selectedStatus.value

  try {
    const appSource = await useApi('/sumbercust')
    const sourcesData = Array.isArray(appSource?.data) ? appSource.data : []

    const { data, error, unauthorized } = await useApi('/inquiry/filter', {
      params: {
        status: statusParam,
        status_filter: statusParam,
        search: searchQuery.value || '',
        start: 0,
        length: 300
      }
    })

    // Session expired: useApi already redirects to /login, so keep the current list as-is
    if (unauthorized) return
    if (error) {
      console.error('Failed to fetch inquiries:', error)
      return
    }

    const rawList = Array.isArray(data) ? data : (data?.data || [])

    if (rawList && rawList.length >= 0) {
      let list = rawList.map((i: any) => {
        const sourceObj = sourcesData.find((value: any) => value.kd_ref_cust === i.id_sumbercust)
        const source = sourceObj?.nm_ref_cust || i.id_sumbercust || 'Walk In'
        const statusName = i.nama_hslaktiv || (Number(i.id_hslaktiv) === 1 ? 'Leads' : Number(i.id_hslaktiv) === 2 ? 'Prospect' : Number(i.id_hslaktiv) === 3 ? 'Hot Prospect' : Number(i.id_hslaktiv) === 7 ? 'SBG' : 'Leads')
        return {
          ...i,
          id_customer: i.id_customer || i.id_rencana,
          nama: i.nama || i.nm_customer || i.name || 'Customer',
          hp: i.hp || i.phone || i.no_hp || '-',
          status: statusName,
          source: source,
          foto: i.url_img || i.foto || i.gambar || '',
          aktivitas: i.ket_aktivitas || '-',
          ...getRencanaDetails(i),
          keterangan: getRencanaDetails(i).keterangan_aktivitas
        }
      })

      // Exclude SBG / Deal status from Inquiry page list unless SBG selected
      if (selectedStatus.value !== 'SBG') {
        list = list.filter((i: any) => {
          const st = (i.status || i.nama_hslaktiv || '').toLowerCase()
          const idAktiv = Number(i.id_hslaktiv)
          return idAktiv !== 7 && !st.includes('sbg') && !st.includes('deal') && !st.includes('do')
        })
      }

      if (selectedSource.value !== 'all') {
        list = list.filter((i: any) => (i.source || 'Walk In').toLowerCase().includes(selectedSource.value.toLowerCase()))
      }

      inquiries.value = list

      // Recalculate summary metrics if summaryCounts are 0
      if (summaryCounts.leads === 0 && summaryCounts.prospect === 0 && summaryCounts.hotProspect === 0) {
        summaryCounts.leads = rawList.filter((i: any) => Number(i.id_hslaktiv) === 1 || (i.nama_hslaktiv || '').toLowerCase().includes('lead')).length
        summaryCounts.prospect = rawList.filter((i: any) => Number(i.id_hslaktiv) === 2 || (i.nama_hslaktiv || '').toLowerCase() === 'prospect').length
        summaryCounts.hotProspect = rawList.filter((i: any) => Number(i.id_hslaktiv) === 3 || (i.nama_hslaktiv || '').toLowerCase().includes('hot')).length
        summaryCounts.sbg = rawList.filter((i: any) => Number(i.id_hslaktiv) === 7 || (i.nama_hslaktiv || '').toLowerCase().includes('sbg')).length
      }
    } else {
      inquiries.value = []
    }
  } catch (err) {
    console.error('Failed to fetch inquiries:', err)
  } finally {
    isLoadingTable.value = false
  }
}

const fetchExistingCustomers = async () => {
  const { data } = await useApi('/customer/select')
  if (data && Array.isArray(data) && data.length > 0) {
    existingCustomers.value = data.map((c: any) => ({
      id_customer: c.id,
      nama: c.text,
      hp: ''
    }))
  } else {
    const { data: inqData } = await useApi('/inquiry/filter', {
      params: { start: 0, length: 100 }
    })
    if (inqData && inqData.data) {
      const map = new Map()
      for (const item of inqData.data) {
        map.set(item.nama, { id_customer: item.id_customer || item.id_rencana, nama: item.nama, hp: item.id_sumbercust || 'Walk In' })
      }
      existingCustomers.value = Array.from(map.values())
    }
  }
}

let debounceTimer: any = null
const debounceFetch = () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    fetchInquiries()
  }, 400)
}

const handleSubmitRencana = async () => {
  isSubmitting.value = true
  modalError.value = ''

  if (!form.soa || !form.aktifitas || !form.jenisAplikasi) {
    modalError.value = 'Silakan lengkapi pilihan SOA, Aktifitas, dan Jenis Aplikasi terlebih dahulu.'
    isSubmitting.value = false
    return
  }

  if (form.cust === 'Lama' && form.custlama) {
    const custId = Number(form.custlama)
    const existingCount = inquiries.value.filter(
      (i: any) => Number(i.id_customer) === custId
    ).length
    if (existingCount >= 2) {
      modalError.value = 'Batas maksimal rencana / leads per customer hanya boleh 2. Customer ini sudah memiliki 2 rencana.'
      isSubmitting.value = false
      return
    }
  }

  const payload: any = {
    cust: form.cust,
    sumbercust: form.soa,
    hslaktiv: 1, // Default status Leads
    katcust: '1',
    ket: form.ket,
    aktiv: form.aktiv,
    soa: form.soa,
    aktifitas: form.aktifitas,
    jenis_aplikasi: form.jenisAplikasi
  }

  console.log('payload', payload);

  if (form.cust === 'Lama') {
    payload.custlama = Number(form.custlama)
  } else {
    payload.nm = form.nm
    payload.hp = form.hp
  }

  const { data, error } = await useApi('/rencana/input', {
    method: 'POST',
    body: payload
  })

  isSubmitting.value = false

  if (error) {
    modalError.value = typeof error === 'string' ? error : (error.message || 'Gagal menyimpan rencana')
    return
  }

  fetchSummaryCounts()
  fetchInquiries()
  closeFormModal()
}

const handleDelete = async (id: number) => {
  if (confirm('Apakah Anda yakin ingin menghapus rencana aktivitas ini?')) {
    const { error } = await useApi(`/rencana/${id}`, { method: 'DELETE' })
    if (!error) {
      fetchSummaryCounts()
      fetchInquiries()
    }
  }
}

const showDetail = (item: any) => {
  selectedItem.value = item

  // Load photo from database (tb_img url_img), item object, or client storage
  const photo =
    item.url_img ||
    item.foto ||
    item.gambar ||
    getFotoBarang(item.id_rencana) ||
    getFotoBarang(item.id_customer) ||
    null

  detailFotoUrl.value = photo
}

const closeDetailModal = () => {
  selectedItem.value = null
  showLightbox.value = false
}

const openLightbox = () => {
  showLightbox.value = true
}

const closeLightbox = () => {
  showLightbox.value = false
}

const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape') {
    if (showLightbox.value) {
      closeLightbox()
    } else if (selectedItem.value) {
      closeDetailModal()
    }
  }
}

const closeFormModal = () => {
  openModal.value = false
  modalError.value = ''
  form.cust = 'Baru'
  form.custlama = ''
  form.nm = ''
  form.hp = ''
  form.ket = ''
  form.soa = ''
  form.aktifitas = ''
  form.jenisAplikasi = ''
}

const getStatusBadgeClass = (status: string) => {
  switch (status) {
    case 'Leads': return 'bg-amber-100 text-amber-800 border border-amber-200'
    case 'Prospect': return 'bg-blue-100 text-blue-800 border border-blue-200'
    case 'Hot Prospect': return 'bg-rose-100 text-rose-800 border border-rose-200'
    case 'SBG': return 'bg-emerald-100 text-emerald-800 border border-emerald-200'
    case 'Batal': return 'bg-slate-100 text-slate-600 border border-slate-200'
    default: return 'bg-slate-100 text-slate-600 border border-slate-200'
  }
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

const formatFullDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return dateStr
  return d.toLocaleDateString('id-ID', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  })
}
</script>
