<template>
  <div class="space-y-6 pb-12 font-sans">
    <!-- Header Section -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-semibold text-slate-800 tracking-tight">
          Dashboard Ringkasan
        </h1>
        <p class="text-slate-500 text-sm">
          Selamat datang kembali, pantau performa Anda hari ini.
        </p>
      </div>
      <div class="flex items-center gap-3">
        <div v-if="isReview && salesList.length > 0" class="flex items-center gap-2">
          <label class="text-xs font-semibold text-slate-500 uppercase tracking-wider hidden sm:inline">Review
            Sales:</label>
          <select v-model="selectedSalesUser" @change="fetchDashboardMetrics"
            class="px-3.5 py-2 bg-white border border-slate-300 rounded-xl text-xs sm:text-sm text-slate-800 focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all font-medium">
            <option v-for="s in salesList" :key="s.username" :value="s.username">
              {{ s.nama }} ({{ s.nm_jabatan || "Sales" }})
            </option>
          </select>
        </div>
        <NuxtLink to="/inquiry"
          class="inline-flex items-center gap-2 bg-amber-500 hover:bg-amber-600 active:scale-95 text-white font-medium px-5 py-2.5 rounded-xl text-sm transition-all shadow-md shadow-amber-200">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Leads Baru
        </NuxtLink>
      </div>
    </div>

    <!-- HERO & LEADERBOARD ROW -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Hero Card (2 cols) -->
      <div
        class="lg:col-span-2 bg-[#1e293b] text-white rounded-3xl p-6 md:p-8 shadow-xl relative overflow-hidden flex flex-col justify-between min-h-[220px]">
        <!-- Background Decorative Glow -->
        <div class="absolute -top-24 -right-24 w-72 h-72 bg-amber-500/10 rounded-full blur-3xl pointer-events-none">
        </div>
        <div class="absolute -bottom-24 -left-24 w-72 h-72 bg-emerald-500/10 rounded-full blur-3xl pointer-events-none">
        </div>

        <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-6 relative z-10">
          <div class="space-y-3 flex-1">
            <div class="flex items-center gap-3">
              <span class="text-lg font-semibold uppercase tracking-wider text-amber-400">
                {{
                  isQualified
                    ? "Estimasi Insentif Bulan Ini"
                    : "Potensi Insentif Terkumpul"
                }}
              </span>

              <!-- Privacy Toggle Button -->
              <button @click="isPrivate = !isPrivate" type="button" :class="isPrivate
                ? 'bg-rose-500/20 text-rose-300 border-rose-500/30'
                : 'bg-white/10 text-slate-300 border-white/20 hover:bg-white/20'
                " class="p-1.5 rounded-lg border text-xs flex items-center gap-1.5 transition-all cursor-pointer"
                :title="isPrivate ? 'Tampilkan Nilai' : 'Sembunyikan Nilai'">
                <!-- Eye off icon when private -->
                <svg v-if="isPrivate" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
                  fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
                  data-lucide="eye-off" aria-hidden="true" id="eyeIcon" class="lucide lucide-eye-off w-5 h-5">
                  <path
                    d="M10.733 5.076a10.744 10.744 0 0 1 11.205 6.575 1 1 0 0 1 0 .696 10.747 10.747 0 0 1-1.444 2.49">
                  </path>
                  <path d="M14.084 14.158a3 3 0 0 1-4.242-4.242"></path>
                  <path
                    d="M17.479 17.499a10.75 10.75 0 0 1-15.417-5.151 1 1 0 0 1 0-.696 10.75 10.75 0 0 1 4.446-5.143">
                  </path>
                  <path d="m2 2 20 20"></path>
                </svg>
                <!-- Eye icon when visible -->
                <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
              </button>
            </div>

            <!-- Incentive Display -->
            <div class="flex items-baseline gap-2">
              <span v-if="!isPrivate" class="text-2xl text-slate-400 font-light">Rp</span>
              <h2 class="text-4xl md:text-5xl font-normal tracking-tight text-white select-none">
                <template v-if="isPrivate">******</template>
                <template v-else>{{ formatNum(estIncentive) }}</template>
              </h2>
            </div>

            <!-- Qualification Status Badge -->
            <div class="pt-1">
              <span v-if="isPrivate" class="text-sm font-normal text-slate-400">*****</span>
              <div v-else class="inline-flex items-center gap-2 px-3 py-1.5 rounded-xl text-xs font-normal" :class="isQualified
                ? 'bg-emerald-500/10 border border-emerald-500/20 text-emerald-400'
                : 'bg-amber-500/10 border border-amber-500/20 text-amber-400'
                ">
                <svg v-if="isQualified" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
                <span>Status: {{ isQualified ? "Qualified" : "Dalam Proses" }} |
                  Grade: {{ userGrade }} (Rate: {{ appliedRate }}%)</span>
              </div>
            </div>
          </div>

          <!-- Target Progress Card (Right Side inside Hero) -->
          <div class="w-full md:w-80 bg-slate-800/40 backdrop-blur-md rounded-2xl p-5 border border-white/10 space-y-3">
            <div class="flex justify-between items-center text-xs">
              <span class="text-slate-400 font-medium uppercase tracking-wider">Target Progress</span>
              <span class="font-medium text-white">{{ pointTotal }}/{{ targetGoalPoints }} Poin</span>
            </div>
            <div class="w-full bg-slate-700/50 rounded-full h-3 overflow-hidden p-0.5">
              <div class="bg-amber-400 h-full rounded-full transition-all duration-1000 ease-out" :style="{
                width: `${Math.min(100, Math.max(0, targetProgressPct))}%`,
              }"></div>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-xs text-emerald-400 font-semibold">{{ targetProgressPct }}%</span>
              <span class="text-[11px] text-slate-400">{{
                targetDescription
              }}</span>
            </div>
            <div class="pt-1 border-t border-white/5 flex items-center gap-2 text-xs text-slate-300">
              <svg class="w-4 h-4 text-amber-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
              <span>{{ motivationText }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Leaderboard / Ranking Card (1 col) -->
      <div
        class="bg-white rounded-3xl p-6 border border-slate-200 shadow-sm flex flex-col items-center justify-center text-center relative group overflow-hidden">
        <div class="absolute top-0 right-0 p-4 opacity-5 group-hover:scale-110 transition-transform duration-500">
          <svg class="w-24 h-24 text-slate-900" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z" />
          </svg>
        </div>

        <div
          class="w-14 h-14 rounded-2xl bg-amber-50 text-amber-500 flex items-center justify-center mb-3 shadow-inner border border-amber-100">
          <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z" />
          </svg>
        </div>
        <h3 class="text-slate-400 text-xs font-semibold uppercase tracking-wider mb-1">
          Posisi Leaderboard
        </h3>
        <div class="text-3xl font-bold text-slate-800 tracking-tight mb-2">
          #{{ leaderboard.rank }}
          <span class="text-sm font-normal text-slate-500">dari {{ leaderboard.total_participants }} Sales</span>
        </div>
        <p class="text-xs text-slate-500 max-w-[240px] leading-relaxed">
          <template v-if="leaderboard.rank === 1">
            <strong class="text-emerald-600 font-semibold">Luar biasa!</strong>
            Anda memimpin peringkat #1!
          </template>
          <template v-else>
            Butuh
            <strong class="text-emerald-600 font-semibold">Rp {{ formatNum(leaderboard.gap_to_next) }}</strong>
            omzet lagi untuk menyalip ke peringkat #{{ leaderboard.rank - 1 }}.
          </template>
        </p>
      </div>
    </div>

    <!-- METRICS ROW -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- Card 1: Emas -->
      <div
        class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-sm hover:shadow-md transition-all border-l-4 border-l-amber-500 flex flex-col justify-between">
        <div>
          <div class="flex justify-between items-start mb-3">
            <div class="p-2.5 bg-amber-50 text-amber-600 rounded-xl flex items-center justify-center shadow-inner">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <span
              class="text-[10px] font-semibold text-amber-700 bg-amber-100/80 px-2.5 py-1 rounded-full uppercase tracking-wider">EMAS</span>
          </div>
          <h4 class="text-slate-500 text-xs font-semibold uppercase tracking-wider mb-1">
            Pencairan Gadai Emas
          </h4>
          <div class="text-2xl font-bold text-slate-800 mb-1 tracking-tight">
            Rp {{ formatNum(stats.emas) }}
          </div>
          <div class="text-xs text-slate-500 font-medium mb-3">
            Total Gramasi:
            <strong class="text-amber-600 font-semibold">{{ gramasiEmas }}g</strong>
            / {{ emasCard.target_emas_g }}g
          </div>
        </div>

        <div class="pt-3 border-t border-slate-100 space-y-2">
          <div class="flex items-center justify-between text-[11px]">
            <span class="font-medium flex items-center gap-1" :class="emasCard.is_emas_bonus_done
              ? 'text-emerald-600'
              : 'text-amber-600'
              ">
              {{ emasCard.bonus_emas_title }}
            </span>
            <span class="font-semibold text-slate-600">{{ emasCard.pct_emas }}%</span>
          </div>
          <div class="w-full bg-slate-100 rounded-full h-2 overflow-hidden">
            <div class="h-full rounded-full transition-all duration-1000" :class="emasCard.is_emas_bonus_done
              ? 'bg-emerald-500'
              : 'bg-gradient-to-r from-amber-400 to-amber-500'
              " :style="{ width: `${emasCard.pct_emas}%` }"></div>
          </div>
          <p class="text-[11px] text-slate-400 font-medium" v-html="emasCard.sub_emas_text"></p>
        </div>
      </div>

      <!-- Card 2: Non Emas / Elektronik -->
      <div
        class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-sm hover:shadow-md transition-all border-l-4 border-l-emerald-500 flex flex-col justify-between">
        <div>
          <div class="flex justify-between items-start mb-3">
            <div class="p-2.5 bg-sky-50 text-sky-600 rounded-xl flex items-center justify-center shadow-inner">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
              </svg>
            </div>
            <span
              class="text-[10px] font-semibold text-sky-700 bg-sky-100/80 px-2.5 py-1 rounded-full uppercase tracking-wider">NON
              EMAS</span>
          </div>
          <h4 class="text-slate-500 text-xs font-semibold uppercase tracking-wider mb-1">
            Pencairan Gadai Non Emas
          </h4>
          <div class="text-2xl font-bold text-slate-800 mb-1 tracking-tight">
            Rp {{ formatNum(stats.nonemas) }}
          </div>
          <div class="text-xs text-slate-500 font-medium mb-3">
            Target UP:
            <strong class="text-sky-600 font-semibold">50 Jt</strong> (Bonus 500
            Rb)
          </div>
        </div>

        <div class="pt-3 border-t border-slate-100 space-y-2">
          <div class="flex items-center justify-between text-[11px]">
            <span class="font-medium flex items-center gap-1" :class="nonEmasCard.is_nonemas_bonus_done
              ? 'text-emerald-600'
              : 'text-sky-600'
              ">
              {{ nonEmasCard.bonus_title }}
            </span>
            <span class="font-semibold text-slate-600">{{ nonEmasCard.pct_nonemas }}%</span>
          </div>
          <div class="w-full bg-slate-100 rounded-full h-2 overflow-hidden">
            <div class="h-full rounded-full transition-all duration-1000" :class="nonEmasCard.is_nonemas_bonus_done
              ? 'bg-emerald-500'
              : 'bg-gradient-to-r from-sky-400 to-sky-500'
              " :style="{ width: `${nonEmasCard.pct_nonemas}%` }"></div>
          </div>
          <p class="text-[11px] text-slate-400 font-medium" v-html="nonEmasCard.sub_text"></p>
        </div>
      </div>

      <!-- Card 3: Nasabah Baru (New CIF) -->
      <div
        class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-sm hover:shadow-md transition-all border-l-4 border-l-blue-500 flex flex-col justify-between">
        <div>
          <div class="flex justify-between items-start mb-4">
            <div class="p-2.5 bg-blue-50 text-blue-600 rounded-xl flex items-center justify-center shadow-inner">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
              </svg>
            </div>
            <span
              class="text-[10px] font-medium text-blue-600 bg-blue-50 px-2 py-0.5 rounded-full uppercase tracking-wider">Nasabah</span>
          </div>
          <h4 class="text-slate-500 text-xs font-semibold uppercase tracking-wider mb-1">
            New CIF (Nasabah Baru)
          </h4>
          <div class="text-2xl font-bold text-slate-800 mb-1 tracking-tight">
            {{ stats.newcif }}
          </div>
        </div>
        <div class="pt-3 border-t border-slate-100 text-[11px] font-medium text-blue-600 flex items-center gap-1.5">
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
          </svg>
          CIF Bulan Lalu: {{ lastMonthCif }}
        </div>
      </div>

      <!-- Card 4: New Leads -->
      <div
        class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-sm hover:shadow-md transition-all border-l-4 border-l-rose-500 flex flex-col justify-between">
        <div>
          <div class="flex justify-between items-start mb-4">
            <div class="p-2.5 bg-rose-50 text-rose-600 rounded-xl flex items-center justify-center shadow-inner">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <span
              class="text-[10px] font-medium text-rose-600 bg-rose-50 px-2 py-0.5 rounded-full uppercase tracking-wider">Leads</span>
          </div>
          <h4 class="text-slate-500 text-xs font-semibold uppercase tracking-wider mb-1">
            New Leads
          </h4>
          <div class="text-2xl font-bold text-slate-800 mb-1 tracking-tight">
            {{ monthlyFunnel.leads }}
          </div>
        </div>
        <div class="pt-3 border-t border-slate-100 text-[11px] font-medium text-rose-600 flex items-center gap-1.5">
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          Leads Bulan Lalu: {{ lastMonthLeads }}
        </div>
      </div>
    </div>

    <!-- SALES FUNNEL SECTION -->
    <div class="bg-white rounded-3xl border border-slate-200 p-6 shadow-sm space-y-6">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2.5">
          <div class="w-7 h-7 rounded-lg bg-slate-100 flex items-center justify-center text-slate-500">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
          </div>
          <h3 class="font-semibold text-slate-800 text-lg">Sales Funnel</h3>
        </div>
        <button @click="refetch()" :disabled="isFetching"
          class="flex items-center gap-1.5 text-xs font-semibold px-3 py-1.5 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-xl transition-all disabled:opacity-50">
          <svg class="w-3.5 h-3.5" :class="{ 'animate-spin': isFetching }" fill="none" stroke="currentColor"
            viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          <span>{{ isFetching ? "Refetching..." : "Refresh Data" }}</span>
        </button>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Column 1: Akumulasi Keseluruhan -->
        <div class="bg-slate-50/80 rounded-3xl p-6 border border-slate-100/90 space-y-5">
          <h4 class="text-sm font-bold text-slate-800 text-center">
            Akumulasi Keseluruhan
          </h4>

          <div class="flex">
            <div class="sales-funnel" style="width: 100%; max-width: none">
              <!-- Stage 1: Leads -->
              <div class="relative flex items-center justify-between gap-3">
                <div class="flex-1 flex justify-center">
                  <div
                    class="bg-[#E50000] text-white text-xs font-bold py-2.5 px-4 rounded-t-md rounded-b-3xl text-center shadow-2xs transition-all duration-300"
                    :style="{
                      width: getFunnelWidth(1, stats.leads, stats.leads),
                    }">
                    Leads - {{ stats.leads }}
                  </div>
                </div>
              </div>

              <!-- Stage 2: Prospek -->
              <div class="relative flex items-center justify-between gap-3">
                <div class="flex-1 flex justify-center">
                  <div
                    class="bg-[#8B0000] text-white text-xs font-bold py-2.5 px-4 rounded-b-3xl text-center shadow-2xs transition-all duration-300"
                    :style="{
                      width: getFunnelWidth(2, stats.prospect, stats.leads),
                    }">
                    Prospek - {{ stats.prospect }}
                  </div>
                </div>
              </div>

              <!-- Stage 3: Hot Prospek -->
              <div class="relative flex items-center justify-between gap-3">
                <div class="flex-1 flex justify-center">
                  <div
                    class="bg-[#F59E0B] text-white text-xs font-bold py-2.5 px-4 rounded-b-3xl text-center shadow-2xs transition-all duration-300"
                    :style="{
                      width: getFunnelWidth(
                        3,
                        stats.hotprospect,
                        stats.prospect || stats.leads,
                      ),
                    }">
                    Hot Prospek - {{ stats.hotprospect }}
                  </div>
                </div>
              </div>

              <!-- Stage 4: DO -->
              <div class="relative flex items-center justify-between gap-3">
                <div class="flex-1 flex justify-center">
                  <div
                    class="bg-[#00C853] text-white text-xs font-bold py-2.5 px-4 rounded-b-3xl text-center shadow-2xs transition-all duration-300"
                    :style="{
                      width: getFunnelWidth(
                        4,
                        stats.sbg,
                        stats.hotprospect || stats.leads,
                      ),
                    }">
                    DO - {{ stats.sbg }}
                  </div>
                </div>
              </div>
            </div>
            <div class="sales-percent flex flex-col justify-around">
              <div
                class="w-16 text-right font-semibold text-xs text-red-500 flex items-center justify-end gap-0.5 whitespace-nowrap">
                {{ calcDrop(stats.prospect, stats.leads) }}
              </div>
              <div
                class="w-16 text-right font-semibold text-xs text-amber-500 flex items-center justify-end gap-0.5 whitespace-nowrap">
                {{ calcDrop(stats.hotprospect, stats.prospect) }}
              </div>
              <div
                class="w-16 text-right font-semibold text-xs text-green-500 flex items-center justify-end gap-0.5 whitespace-nowrap">
                {{ calcDrop(stats.sbg, stats.hotprospect) }}
              </div>
            </div>
          </div>
        </div>

        <!-- Column 2: Bulan ini -->
        <div class="bg-slate-50/80 rounded-3xl p-6 border border-slate-100/90 space-y-5">
          <h4 class="text-sm font-bold text-slate-800 text-center">
            Bulan ini
          </h4>

          <div class="flex">
            <div class="sales-funnel" style="width: 100%; max-width: none">
              <!-- Stage 1: Leads -->
              <div class="relative flex items-center justify-between gap-3">
                <div class="flex-1 flex justify-center">
                  <div
                    class="bg-[#E50000] text-white text-xs font-bold py-2.5 px-4 rounded-t-md rounded-b-3xl text-center shadow-2xs transition-all duration-300"
                    :style="{
                      width: getFunnelWidth(
                        1,
                        monthlyFunnel.leads,
                        monthlyFunnel.leads,
                      ),
                    }">
                    Leads - {{ monthlyFunnel.leads }}
                  </div>
                </div>
              </div>

              <!-- Stage 2: Prospek -->
              <div class="relative flex items-center justify-between gap-3">
                <div class="flex-1 flex justify-center">
                  <div
                    class="bg-[#8B0000] text-white text-xs font-bold py-2.5 px-4 rounded-b-3xl text-center shadow-2xs transition-all duration-300"
                    :style="{
                      width: getFunnelWidth(
                        2,
                        monthlyFunnel.prospect,
                        monthlyFunnel.leads,
                      ),
                    }">
                    Prospek - {{ monthlyFunnel.prospect }}
                  </div>
                </div>
              </div>

              <!-- Stage 3: Hot Prospek -->
              <div class="relative flex items-center justify-between gap-3">
                <div class="flex-1 flex justify-center">
                  <div
                    class="bg-[#F59E0B] text-white text-xs font-bold py-2.5 px-4 rounded-b-3xl text-center shadow-2xs transition-all duration-300"
                    :style="{
                      width: getFunnelWidth(
                        3,
                        monthlyFunnel.hotprospect,
                        monthlyFunnel.prospect || monthlyFunnel.leads,
                      ),
                    }">
                    Hot Prospek - {{ monthlyFunnel.hotprospect }}
                  </div>
                </div>
              </div>

              <!-- Stage 4: DO -->
              <div class="relative flex items-center justify-between gap-3">
                <div class="flex-1 flex justify-center">
                  <div
                    class="bg-[#00C853] text-white text-xs font-bold py-2.5 px-4 rounded-b-3xl text-center shadow-2xs transition-all duration-300"
                    :style="{
                      width: getFunnelWidth(
                        4,
                        monthlyFunnel.sbg,
                        monthlyFunnel.hotprospect || monthlyFunnel.leads,
                      ),
                    }">
                    DO - {{ monthlyFunnel.sbg }}
                  </div>
                </div>
              </div>
            </div>
            <div class="sales-percent flex flex-col justify-around mt-0!">
              <div
                class="w-16 text-right font-semibold text-xs text-red-500 flex items-center justify-end gap-0.5 whitespace-nowrap">
                {{ calcDrop(monthlyFunnel.prospect, monthlyFunnel.leads) }}
              </div>
              <div
                class="w-16 text-right font-semibold text-xs text-amber-500 flex items-center justify-end gap-0.5 whitespace-nowrap">
                {{ calcDrop(monthlyFunnel.hotprospect, monthlyFunnel.prospect) }}
              </div>
              <div
                class="w-16 text-right font-semibold text-xs text-green-500 flex items-center justify-end gap-0.5 whitespace-nowrap">
                {{ calcDrop(monthlyFunnel.sbg, monthlyFunnel.hotprospect) }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- PACING CHART & LIVE FEED ROW -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Daily Pacing Chart (2 cols) -->
      <div
        class="lg:col-span-2 bg-white rounded-3xl border border-slate-200 p-6 shadow-sm flex flex-col justify-between">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 mb-4">
          <div>
            <h3 class="font-semibold text-slate-800 text-lg">
              Grafik Pacing Harian (Booking Selada vs Sopiga Pro)
            </h3>
            <p class="text-xs tdext-slate-500">
              Monitoring realisasi omzet harian bulan ini
            </p>
          </div>
          <div class="flex items-center gap-4 text-xs font-medium">
            <div class="flex items-center gap-1.5">
              <span class="w-3 h-3 rounded-full bg-emerald-500"></span>
              <span class="text-slate-600">Selada</span>
            </div>
            <div class="flex items-center gap-1.5">
              <span class="w-3 h-3 rounded-full bg-sky-500"></span>
              <span class="text-slate-600">Sopiga Pro</span>
            </div>
          </div>
        </div>

        <!-- Pacing Chart Bars -->
        <div class="h-44 flex items-end justify-between gap-1 pt-6 pb-2 px-1 relative">
          <div v-for="d in daysInMonth" :key="d"
            class="flex-1 h-full flex flex-col justify-end items-center relative group" :class="d === currentDay
              ? 'ring-2 ring-amber-400 ring-offset-1 z-10 rounded-t-sm'
              : ''
              ">
            <!-- Tooltip -->
            <div
              class="absolute bottom-full mb-2 bg-slate-800 text-white text-[10px] p-2.5 rounded-xl shadow-2xl opacity-0 group-hover:opacity-100 transition-opacity z-50 whitespace-nowrap pointer-events-none space-y-1"
              :class="d <= 4
                ? 'left-0'
                : d >= daysInMonth - 4
                  ? 'right-0'
                  : 'left-1/2 -translate-x-1/2'
                ">
              <div class="font-semibold border-b border-slate-700 pb-1 text-amber-400">
                Tgl {{ d }} {{ d === currentDay ? "(Hari Ini)" : "" }}
              </div>
              <div class="flex items-center gap-1.5 text-emerald-400 font-medium">
                <span class="w-2 h-2 rounded-full bg-emerald-400"></span>
                <span>Selada: Rp {{ formatNum(getPacingVal(d, "selada")) }}</span>
              </div>
              <div class="flex items-center gap-1.5 text-sky-400 font-medium">
                <span class="w-2 h-2 rounded-full bg-sky-400"></span>
                <span>Sopiga Pro: Rp
                  {{ formatNum(getPacingVal(d, "sopiga")) }}</span>
              </div>
            </div>

            <!-- Dual Bar -->
            <div class="w-full flex items-end justify-center gap-[2px] h-full px-[1px]">
              <div class="w-1/2 rounded-t-sm transition-all duration-500" :class="getPacingVal(d, 'selada') > 0
                ? 'bg-emerald-500 shadow-sm'
                : d <= currentDay
                  ? 'bg-emerald-200/60'
                  : 'bg-slate-100'
                " :style="{
                  height: `${getPacingHeight(d, 'selada')}%`,
                  opacity: d <= currentDay ? 1 : 0.4,
                }"></div>
              <div class="w-1/2 rounded-t-sm transition-all duration-500" :class="getPacingVal(d, 'sopiga') > 0
                ? 'bg-sky-500 shadow-sm'
                : d <= currentDay
                  ? 'bg-sky-200/60'
                  : 'bg-slate-100'
                " :style="{
                  height: `${getPacingHeight(d, 'sopiga')}%`,
                  opacity: d <= currentDay ? 1 : 0.4,
                }"></div>
            </div>
          </div>
        </div>
        <div class="flex justify-between items-center text-[10px] text-slate-400 pt-2 border-t border-slate-100">
          <span>Tgl 1</span>
          <span>Tgl 15</span>
          <span>Tgl {{ daysInMonth }}</span>
        </div>
      </div>

      <!-- Live Feed (1 col) -->
      <div class="bg-white rounded-3xl border border-slate-200 p-6 shadow-sm flex flex-col">
        <div class="flex items-center justify-between mb-4">
          <h3 class="font-semibold text-slate-800 text-lg flex items-center gap-2">
            Live Insentif
            <span class="flex h-2 w-2 relative">
              <span
                class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
            </span>
          </h3>
        </div>

        <div class="flex-1 space-y-3 overflow-y-auto max-h-[350px] pr-1">
          <div v-for="(feed, idx) in liveFeeds" :key="idx"
            class="p-3 rounded-xl bg-slate-50 border border-slate-100 flex items-center justify-between group hover:bg-white hover:shadow-md transition-all">
            <div class="flex items-center gap-3">
              <div
                class="w-9 h-9 rounded-full flex items-center justify-center shadow-sm border text-[10px] font-semibold uppercase tracking-wider"
                :class="feed.isBonus
                  ? 'bg-amber-50 text-amber-600 border-amber-200'
                  : 'bg-emerald-50 text-emerald-600 border-emerald-200'
                  ">
                <span>{{ feed.iconText }}</span>
              </div>
              <div>
                <p class="text-xs font-semibold text-slate-800">
                  {{ feed.type }}
                </p>
                <p class="text-[10px] text-slate-500 uppercase tracking-tight">
                  {{ feed.code }} • {{ feed.customer }}
                </p>
              </div>
            </div>
            <div class="text-right">
              <p class="text-xs font-semibold text-emerald-600">
                +Rp {{ formatNum(feed.amount) }}
              </p>
              <p class="text-[9px] text-amber-600 font-medium">
                {{ feed.sub }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from "vue";
import { useApi } from "~/composables/useApi";

// Privacy toggle state (default censored)
const isPrivate = ref(true);

// Stats from backend API
const stats = ref({
  all: 0,
  leads: 0,
  prospect: 0,
  hotprospect: 0,
  sbg: 0,
  batal: 0,
  newcif: 0,
  emas: 0,
  nonemas: 0,
  leads_monthly: 0,
});

// Review mode & sales list
const isReview = ref(false);
const salesList = ref<any[]>([]);
const selectedSalesUser = ref("");

// Monthly funnel counts
const monthlyFunnel = reactive({
  leads: 0,
  prospect: 0,
  hotprospect: 0,
  sbg: 0,
});

// Extended state matching v_dashboard.php
const estIncentive = ref(0);
const isQualified = ref(false);
const userGrade = ref("JUNIOR");
const appliedRate = ref(25);
const pointTotal = ref(0);
const targetGoalPoints = ref(25);
const dailyTargetProgress = ref(0);
const targetDescription = ref("Mengejar Qualified 1 (Target: 25 Poin)");
const motivationText = ref("Semangat! Kejar Poin untuk Bonus Insentif!");

const salesName = ref("Sales");
const leaderboard = ref({
  rank: 1,
  total_participants: 1,
  gap_to_next: 0,
});

const gramasiEmas = ref(0);
const lastMonthCif = ref(0);
const lastMonthLeads = ref(0);

const emasCard = ref({
  target_emas_g: 300,
  pct_emas: 0,
  bonus_emas_title: "Tier Bonus: 300g = 1 Jt",
  sub_emas_text: "Target berikutnya: 300g (Bonus 1 Jt)",
  is_emas_bonus_done: false,
});

const nonEmasCard = ref({
  pct_nonemas: 0,
  bonus_title: "Bonus UP 50 Jt = 500 Rb",
  sub_text: "Target UP: 50 Jt",
  is_nonemas_bonus_done: false,
});

// Computed helpers
const targetProgressPct = computed(() => {
  if (
    dailyTargetProgress.value !== undefined &&
    dailyTargetProgress.value !== null
  ) {
    return Math.min(100, Math.max(0, Math.round(dailyTargetProgress.value)));
  }
  if (!targetGoalPoints.value) return 0;
  return Math.min(
    100,
    Math.max(0, Math.round((pointTotal.value / targetGoalPoints.value) * 100)),
  );
});

// Date calculations for Pacing Chart
const now = new Date();
const currentDay = now.getDate();
const daysInMonth = new Date(
  now.getFullYear(),
  now.getMonth() + 1,
  0,
).getDate();

// Daily pacing chart data
const pacingSelada = ref<Record<number, number>>({});
const pacingSopiga = ref<Record<number, number>>({});

const getPacingVal = (day: number, type: "selada" | "sopiga") => {
  const data = type === "selada" ? pacingSelada.value : pacingSopiga.value;
  return data[day] || 0;
};

const getPacingHeight = (day: number, type: "selada" | "sopiga") => {
  const val = getPacingVal(day, type);
  const maxVal = 5000000;
  if (day > currentDay) return 8;
  if (val === 0) return 12;
  return Math.min(100, Math.max(15, (val / maxVal) * 85 + 15));
};

// Live feed events
const liveFeeds = ref([
  {
    type: "Gadai Emas",
    code: "TRX-8412",
    customer: "Budi Santoso",
    amount: 125000,
    sub: "2 Poin • Pinjaman Rp 5.000.000",
    iconText: "AU",
    isBonus: false,
  },
  {
    type: "Nasabah Baru",
    code: "TRX-7741",
    customer: "Siti Rahma",
    amount: 25000,
    sub: "New CIF Bonus",
    iconText: "CIF",
    isBonus: false,
  },
  {
    type: "Gadai Elektronik",
    code: "TRX-6109",
    customer: "Ahmad Fauzi",
    amount: 50000,
    sub: "1 Poin • Laptop Asus",
    iconText: "EL",
    isBonus: false,
  },
  {
    type: "Bonus Gramasi",
    code: "BONUS-101",
    customer: "Achievement 300g",
    amount: 1000000,
    sub: "Target Emas Level 1",
    iconText: "BON",
    isBonus: true,
  },
]);

// Formatting helpers
const formatNum = (val: number) => {
  if (!val && val !== 0) return "0";
  return new Intl.NumberFormat("id-ID").format(val);
};

const calcDrop  = (curr: number, prev: number) => {
  if (!prev || prev <= 0) return "0 %";
  if (curr >= prev) return "0 %";
  const drop = Math.round(100 - (((prev - curr) / prev) * 100));
  return `${drop} %`;
};

const calcConv = (curr: number, prev: number) => {
  if (!prev || prev <= 0) return 0;
  return Math.min(100, Math.round((curr / prev) * 100));
};

const getFunnelWidth = (stage: number, count?: number, base?: number) => {
  // Guaranteed visual tapering bands per stage (1: 100%, 2: 75-85%, 3: 55-65%, 4: 35-45%)
  const stageBases: Record<
    number,
    { min: number; default: number; max: number }
  > = {
    1: { min: 100, default: 100, max: 100 },
    2: { min: 72, default: 80, max: 85 },
    3: { min: 52, default: 60, max: 65 },
    4: { min: 32, default: 40, max: 45 },
  };

  const config = stageBases[stage] || { min: 40, default: 50, max: 60 };
  if (!base || base <= 0 || count === undefined || count === null) {
    return `${config.default}%`;
  }

  const ratio = Math.min(1, Math.max(0, count / base));
  const calculatedWidth = config.min + ratio * (config.max - config.min);
  return `${Math.round(calculatedWidth)}%`;
};

const getWidth = (count: number, base: number) => {
  if (!base || count > base) return "100%";
  return `${Math.max(25, (count / base) * 100)}%`;
};

const rounded = (before: number, after: number) => {
  if (!after) return "0";
  if (after > before) return "rounded-b-md";
  if (before === after) return "";
  // return `${Math.max(25, (before / after) * 100)}%`
};

// TanStack Query integration
const {
  data: metricsData,
  isPending,
  isFetching,
  refetch,
} = useDashboardMetricsQuery(selectedSalesUser);

// Watch and synchronize TanStack Query data into reactive refs
watch(
  metricsData,
  (data) => {
    if (!data) return;
    if (data.counts) {
      stats.value = { ...stats.value, ...data.counts };
    }
    if (data.newcif !== undefined) stats.value.newcif = data.newcif;
    if (data.last_month_cif !== undefined)
      lastMonthCif.value = data.last_month_cif;
    if (data.leads_monthly !== undefined)
      stats.value.leads_monthly = data.leads_monthly;
    if (data.last_month_leads !== undefined)
      lastMonthLeads.value = data.last_month_leads;
    if (data.emas !== undefined) stats.value.emas = data.emas;
    if (data.nonemas !== undefined) stats.value.nonemas = data.nonemas;
    if (data.gramasi_emas !== undefined) gramasiEmas.value = data.gramasi_emas;

    if (data.est_incentive !== undefined)
      estIncentive.value = data.est_incentive;
    if (data.is_qualified !== undefined) isQualified.value = data.is_qualified;
    if (data.user_grade !== undefined) userGrade.value = data.user_grade;
    if (data.applied_rate !== undefined) appliedRate.value = data.applied_rate;
    if (data.poin_total !== undefined) pointTotal.value = data.poin_total;
    if (data.target_goal_points !== undefined)
      targetGoalPoints.value = data.target_goal_points;
    if (data.daily_target_progress !== undefined)
      dailyTargetProgress.value = data.daily_target_progress;
    if (data.target_description !== undefined)
      targetDescription.value = data.target_description;
    if (data.motivation_text !== undefined)
      motivationText.value = data.motivation_text;

    if (data.is_review !== undefined) isReview.value = data.is_review;
    if (data.sales_list && Array.isArray(data.sales_list))
      salesList.value = data.sales_list;
    if (data.selected_sales_user && !selectedSalesUser.value)
      selectedSalesUser.value = data.selected_sales_user;

    if (data.emas_card)
      emasCard.value = { ...emasCard.value, ...data.emas_card };
    if (data.nonemas_card)
      nonEmasCard.value = { ...nonEmasCard.value, ...data.nonemas_card };

    if (data.leaderboard)
      leaderboard.value = { ...leaderboard.value, ...data.leaderboard };
    if (data.live_feeds && Array.isArray(data.live_feeds))
      liveFeeds.value = data.live_feeds;
    if (data.username) salesName.value = data.username;

    if (data.pacing_daily) {
      if (data.pacing_daily.selada) {
        const sMap: Record<number, number> = {};
        for (const [k, v] of Object.entries(data.pacing_daily.selada)) {
          sMap[Number(k)] = Number(v);
        }
        pacingSelada.value = sMap;
      }
      if (data.pacing_daily.sopiga) {
        const spMap: Record<number, number> = {};
        for (const [k, v] of Object.entries(data.pacing_daily.sopiga)) {
          spMap[Number(k)] = Number(v);
        }
        pacingSopiga.value = spMap;
      }
    }

    // Set monthly funnel values directly from TanStack Query payload
    monthlyFunnel.leads =
      data.leads_monthly ??
      data.counts?.leads_monthly ??
      data.counts?.leads ??
      stats.value.leads ??
      0;
    monthlyFunnel.prospect =
      data.prospect_monthly ?? data.counts?.prospect_monthly ?? 0;
    monthlyFunnel.hotprospect =
      data.hotprospect_monthly ?? data.counts?.hotprospect_monthly ?? 0;
    monthlyFunnel.sbg = data.sbg_monthly ?? data.counts?.sbg_monthly ?? 0;
  },
  { immediate: true },
);
</script>

<style>
* {
  font-family: "Poppins", sans-serif;
}
</style>
