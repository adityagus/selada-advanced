<template>
  <div class="space-y-6 pb-12 font-sans">
    <!-- Header Section -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-slate-900 tracking-tight">
          Dashboard
        </h1>
        <p v-if="isReview && selectedSalesUser" class="text-amber-600 font-bold text-xs uppercase tracking-wider mt-1 flex items-center gap-1.5">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
          </svg>
          Meninjau Dashboard Sales: <span class="underline">{{ activeSalesName }}</span>
        </p>
        <!-- <p v-else class="text-slate-500 text-sm">
          Selamat datang kembali, pantau performa Anda hari ini.
        </p> -->
      </div>

      <div class="flex items-center gap-3">
        <NuxtLink to="/inquiry"
          class="inline-flex items-center gap-2 bg-amber-500 hover:bg-amber-600 active:scale-95 text-white font-medium px-5 py-2.5 rounded-xl text-sm transition-all shadow-md shadow-amber-200">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Leads Baru
        </NuxtLink>
      </div>
    </div>

    <!-- Review Sales Controls Card (matching screenshot style) -->
    <div v-if="isReview && salesList.length > 0" class="bg-white border border-slate-200 rounded-2xl p-4 flex flex-wrap items-center justify-between gap-4 shadow-sm">
      <div class="flex items-center gap-2.5 font-bold text-xs text-amber-600 uppercase tracking-wider">
        <span class="w-2.5 h-2.5 rounded-full bg-amber-500 animate-pulse"></span>
        <svg class="w-4 h-4 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
        </svg>
        <span>PILIH SALES PERSON:</span>
      </div>
      <div class="flex items-center gap-3 flex-wrap">
        <select v-model="selectedSalesUser"
          class="px-4 py-2 bg-slate-50 border border-slate-300 rounded-xl text-xs sm:text-sm font-semibold text-slate-800 focus:outline-none focus:border-amber-500 focus:ring-2 focus:ring-amber-200 transition-all min-w-[240px]">
          <option v-for="s in salesList" :key="s.username" :value="s.username">
            {{ s.nama }} - {{ s.nm_jabatan || "SALES" }}
          </option>
        </select>
        <button @click="handleReviewSales" :disabled="isFetching"
          class="inline-flex items-center gap-2 px-5 py-2 bg-[#d97706] hover:bg-[#b45309] active:scale-95 text-white font-bold text-xs sm:text-sm rounded-xl shadow-md shadow-amber-200 transition-all cursor-pointer disabled:opacity-50">
          <svg class="w-4 h-4" :class="{ 'animate-spin': isFetching }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
          </svg>
          <span>Ganti Sales</span>
        </button>
        <button v-if="selectedSalesUser" @click="resetReview" class="text-xs text-slate-500 hover:text-slate-800 font-semibold underline cursor-pointer ml-1">
          Reset
        </button>
      </div>
    </div>

    <!-- HERO & LEADERBOARD ROW -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Hero Card (2 cols) -->
      <div
        class="lg:col-span-2 bg-[#1e293b] text-white rounded-3xl p-6 md:p-8 shadow-xl relative overflow-hidden flex flex-col min-h-[220px] justify-center">
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
                <svg v-if="isPrivate" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M3 3l18 18"/>
                </svg>
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
              <h2 class="text-4xl md:text-5xl font-bold tracking-tight text-white select-none">
                <template v-if="isPrivate">******</template>
                <template v-else>{{ formatNum(estIncentive) }}</template>
              </h2>
            </div>

            <!-- Qualification Status Badge -->
            <div class="pt-1">
              <span v-if="isPrivate" class="text-sm font-normal text-slate-400">*****</span>
              <div v-else class="inline-flex items-center gap-2 px-3 py-1.5 rounded-xl text-xs font-medium" :class="isQualified
                ? 'bg-emerald-500/10 border border-emerald-500/20 text-emerald-400'
                : 'bg-slate-800/80 border border-slate-700 text-slate-300'
                ">
                <svg v-if="isQualified" class="w-4 h-4 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <svg v-else class="w-4 h-4 text-amber-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
                <span>Status: {{ isQualified ? "Qualified" : "Potensi Terkunci" }} |
                  Grade: {{ userGrade }} (Rate: {{ isQualified ? appliedRate : (appliedRate || 25) }}%)</span>
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
              <span class="text-[11px] text-slate-400">{{ targetDescription }}</span>
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
        class="bg-white rounded-3xl p-6 border border-slate-200 shadow-sm flex flex-col items-center justify-between text-center relative group overflow-hidden">
        <div class="absolute top-0 right-0 p-4 opacity-5 group-hover:scale-110 transition-transform duration-500">
          <svg class="w-24 h-24 text-slate-900" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z" />
          </svg>
        </div>

        <div>
          <div
            class="w-14 h-14 rounded-2xl bg-amber-50 text-amber-500 flex items-center justify-center mx-auto mb-3 shadow-inner border border-amber-100">
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
          <p class="text-xs text-slate-500 max-w-[240px] leading-relaxed mx-auto">
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

        <button @click="showLeaderboardModal = true"
          class="mt-4 px-4 py-2 text-xs font-semibold bg-amber-500 hover:bg-amber-600 active:scale-95 text-white rounded-xl shadow-sm transition-all flex items-center justify-center gap-2 w-full cursor-pointer">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
          </svg>
          <span>Lihat Leaderboard Lengkap</span>
        </button>
      </div>
    </div>

    <!-- METRICS ROW (Matching Shopeers Stat Cards style) -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
      <!-- Card 1: Gadai Emas -->
      <div class="bg-white p-5 rounded-3xl border border-slate-200/80 shadow-xs hover:shadow-md transition-all duration-300 flex flex-col justify-between group">
        <div>
          <div class="flex items-center justify-between mb-4">
            <h4 class="font-bold text-slate-800 text-sm md:text-base tracking-tight">
              Emas
            </h4>
            <div class="w-10 h-10 rounded-2xl bg-amber-50 text-amber-600 border border-amber-200/60 flex items-center justify-center flex-shrink-0 shadow-2xs group-hover:scale-105 transition-transform">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
          </div>

          <div class="flex items-baseline justify-between gap-2 mb-2 flex-wrap">
            <span class="text-xl lg:text-2xl font-semibold text-slate-900 tracking-tight">
              Rp {{ formatNum(stats.emas) }}
            </span>
            <span class="inline-flex items-center gap-1 bg-emerald-50 text-emerald-600 border border-emerald-200/60 font-normal text-xs px-2.5 py-1 rounded-full">
              <svg class="w-3 h-3 text-emerald-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 10l7-7m0 0l7 7m-7-7v18"/>
              </svg>
              {{ gramasiEmas }}g
            </span>
          </div>
        </div>

        <div class="pt-3 border-t border-slate-100 space-y-1.5">
          <div class="flex items-center justify-between text-xs">
            <span class="font-medium text-slate-500">Target {{ emasCard.target_emas_g }}g</span>
            <span class="font-bold text-amber-600">{{ emasCard.pct_emas }}%</span>
          </div>
          <div class="w-full bg-slate-100 rounded-full h-2 overflow-hidden">
            <div class="h-full rounded-full transition-all duration-1000" :class="emasCard.is_emas_bonus_done ? 'bg-emerald-500' : 'bg-amber-500'" :style="{ width: `${emasCard.pct_emas}%` }"></div>
          </div>
          <p class="text-[11px] text-slate-400 font-medium truncate">{{ emasCard.bonus_emas_title }}</p>
        </div>
      </div>

      <!-- Card 2: Gadai Non Emas / Elektronik -->
      <div class="bg-white p-5 rounded-3xl border border-slate-200/80 shadow-xs hover:shadow-md transition-all duration-300 flex flex-col justify-between group">
        <div>
          <div class="flex items-center justify-between mb-4">
            <h4 class="font-bold text-slate-800 text-sm md:text-base tracking-tight">
              Elektronik
            </h4>
            <div class="w-10 h-10 rounded-2xl bg-amber-50 text-amber-600 border border-amber-200/60 flex items-center justify-center flex-shrink-0 shadow-2xs group-hover:scale-105 transition-transform">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
              </svg>
            </div>
          </div>

          <div class="flex items-baseline justify-between gap-2 mb-2 flex-wrap">
            <span class="text-xl lg:text-2xl font-bold text-slate-900 tracking-tight">
              Rp {{ formatNum(stats.nonemas) }}
            </span>
            <span class="inline-flex items-center gap-1 bg-amber-50 text-amber-700 border border-amber-200/60 font-medium text-xs px-2.5 py-1 rounded-full">
              Target 50M
            </span>
          </div>
        </div>

        <div class="pt-3 border-t border-slate-100 space-y-1.5">
          <div class="flex items-center justify-between text-xs">
            <span class="font-medium text-slate-500">Pencapaian UP</span>
            <span class="font-bold text-amber-600">{{ nonEmasCard.pct_nonemas }}%</span>
          </div>
          <div class="w-full bg-slate-100 rounded-full h-2 overflow-hidden">
            <div class="h-full rounded-full transition-all duration-1000" :class="nonEmasCard.is_nonemas_bonus_done ? 'bg-emerald-500' : 'bg-amber-500'" :style="{ width: `${nonEmasCard.pct_nonemas}%` }"></div>
          </div>
          <p class="text-[11px] text-slate-400 font-medium truncate">{{ nonEmasCard.bonus_title }}</p>
        </div>
      </div>

      <!-- Card 3: New CIF (Nasabah Baru) - Gold & Black Style -->
      <div class="bg-white p-5 rounded-3xl border border-slate-200/80 shadow-xs hover:shadow-md transition-all duration-300 flex flex-col justify-between group">
        <div>
          <div class="flex items-center justify-between mb-4">
            <h4 class="font-bold text-slate-800 text-sm md:text-base tracking-tight">
              New CIF (Nasabah Baru)
            </h4>
            <div class="w-10 h-10 rounded-2xl bg-slate-900 text-amber-400 border border-slate-800 flex items-center justify-center flex-shrink-0 shadow-2xs group-hover:scale-105 transition-transform">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
              </svg>
            </div>
          </div>

          <div class="flex items-baseline justify-between gap-2 mb-2">
            <span class="text-2xl lg:text-3xl font-extrabold text-slate-900 tracking-tight">
              {{ formatNum(stats.newcif) }}
            </span>
            <span class="inline-flex items-center gap-1 bg-amber-50 text-amber-700 border border-amber-200/60 font-bold text-xs px-2.5 py-1 rounded-full">
              <svg class="w-3 h-3 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 10l7-7m0 0l7 7m-7-7v18"/>
              </svg>
              Active
            </span>
          </div>
        </div>

        <div class="pt-3 border-t border-slate-100 flex items-center justify-between text-xs text-slate-400 font-medium">
          <span>vs. <span class="font-semibold text-slate-500"> {{ formatNum(lastMonthCif) }} Bulan Lalu</span></span>
        </div>
      </div>

      <!-- Card 4: New Leads -->
      <div class="bg-white p-5 rounded-3xl border border-slate-200/80 shadow-xs hover:shadow-md transition-all duration-300 flex flex-col justify-between group">
        <div>
          <div class="flex items-center justify-between mb-4">
            <h4 class="font-bold text-slate-800 text-sm md:text-base tracking-tight">
              New Leads
            </h4>
            <div class="w-10 h-10 rounded-2xl bg-amber-50 text-amber-600 border border-amber-200/60 flex items-center justify-center flex-shrink-0 shadow-2xs group-hover:scale-105 transition-transform">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
          </div>

          <div class="flex items-baseline justify-between gap-2 mb-2">
            <span class="text-2xl lg:text-3xl font-extrabold text-slate-900 tracking-tight">
              {{ formatNum(monthlyFunnel.leads) }}
            </span>
            <span class="inline-flex items-center gap-1 bg-amber-50 text-amber-700 border border-amber-200/60 font-bold text-xs px-2.5 py-1 rounded-full">
              <svg class="w-3 h-3 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 10l7-7m0 0l7 7m-7-7v18"/>
              </svg>
              Leads
            </span>
          </div>
        </div>

        <div class="pt-3 border-t border-slate-100 flex items-center justify-between text-xs text-slate-400 font-medium">
          <span>vs. <span class="font-semibold text-slate-500"> {{ formatNum(lastMonthLeads) }} Bulan Lalu</span></span>
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
          class="flex items-center gap-1.5 text-xs font-semibold px-3 py-1.5 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-xl transition-all disabled:opacity-50 cursor-pointer">
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
            <div class="sales-percent flex flex-col justify-around">
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
            <p class="text-xs text-slate-500">
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
                <span>Sopiga Pro: Rp {{ formatNum(getPacingVal(d, "sopiga")) }}</span>
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
            class="p-3 rounded-2xl bg-slate-50 border border-slate-100 flex items-center justify-between group hover:bg-white hover:shadow-md transition-all">
            <div class="flex items-center gap-3">
              <div class="w-9 h-9 rounded-xl flex items-center justify-center shadow-xs border transition-all"
                :class="getFeedBadgeClass(feed)">
                <!-- Gold / Emas SVG Icon -->
                <svg v-if="getFeedCategory(feed) === 'emas'" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <!-- Non-Emas / Laptop SVG Icon -->
                <svg v-else-if="getFeedCategory(feed) === 'nonemas'" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                </svg>
                <!-- CIF / User SVG Icon -->
                <svg v-else-if="getFeedCategory(feed) === 'cif'" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
                </svg>
                <!-- Bonus / Star SVG Icon -->
                <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
                </svg>
              </div>
              <div>
                <p class="text-xs font-semibold text-slate-800">
                  {{ feed.type }}
                </p>
                <p class="text-[10px] text-slate-500 uppercase tracking-tight font-medium">
                  {{ feed.code }} • {{ feed.customer }}
                </p>
              </div>
            </div>
            <div class="text-right">
              <p class="text-xs font-bold text-emerald-600">
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

    <!-- FULL LEADERBOARD MODAL -->
    <div v-if="showLeaderboardModal" class="fixed inset-0 z-50 bg-slate-900/60 backdrop-blur-sm flex items-center justify-center p-4">
      <div class="bg-white rounded-3xl max-w-2xl w-full p-6 shadow-2xl border border-slate-100 space-y-5 animate-in fade-in zoom-in duration-200">
        <div class="flex items-center justify-between border-b border-slate-100 pb-4">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-2xl bg-amber-50 text-amber-600 flex items-center justify-center font-bold text-lg shadow-inner border border-amber-200">
              <svg class="w-5 h-5 text-amber-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z"/>
              </svg>
            </div>
            <div>
              <h3 class="font-bold text-slate-800 text-lg">Top Sales Leaderboard</h3>
              <p class="text-xs text-slate-500">Peringkat performa pencairan omzet sales bulan ini</p>
            </div>
          </div>
          <button @click="showLeaderboardModal = false" class="p-2 text-slate-400 hover:text-slate-600 rounded-xl hover:bg-slate-100 transition-all cursor-pointer">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <div class="overflow-x-auto max-h-[400px]">
          <table class="w-full text-left text-xs">
            <thead class="bg-slate-50 text-slate-500 font-semibold uppercase tracking-wider sticky top-0">
              <tr>
                <th class="py-3 px-4 rounded-l-xl">Rank</th>
                <th class="py-3 px-4">Nama Sales</th>
                <th class="py-3 px-4 text-right">Gadai Emas</th>
                <th class="py-3 px-4 text-right">Non Emas</th>
                <th class="py-3 px-4 text-right">Total Omzet</th>
                <th class="py-3 px-4 text-right rounded-r-xl">Poin</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 font-medium">
              <tr v-for="p in leaderboard.top_performers" :key="p.rank"
                :class="p.name === salesName || (selectedSalesUser && p.name.includes(selectedSalesUser)) ? 'bg-amber-50/70 text-amber-900 font-bold' : 'hover:bg-slate-50'">
                <td class="py-3 px-4">
                  <span v-if="p.rank === 1" class="inline-flex items-center gap-1 px-2.5 py-0.5 bg-amber-100 text-amber-800 border border-amber-200 rounded-lg font-bold text-xs">
                    <svg class="w-3.5 h-3.5 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z"/>
                    </svg>
                    #1
                  </span>
                  <span v-else-if="p.rank === 2" class="inline-flex items-center gap-1 px-2.5 py-0.5 bg-slate-200 text-slate-800 border border-slate-300 rounded-lg font-bold text-xs">
                    #2
                  </span>
                  <span v-else-if="p.rank === 3" class="inline-flex items-center gap-1 px-2.5 py-0.5 bg-orange-100 text-orange-800 border border-orange-200 rounded-lg font-bold text-xs">
                    #3
                  </span>
                  <span v-else class="text-slate-600 font-semibold px-2">#{{ p.rank }}</span>
                </td>
                <td class="py-3 px-4 text-slate-800 font-semibold">{{ p.name }}</td>
                <td class="py-3 px-4 text-right text-slate-600">Rp {{ formatNum(p.total_emas) }}</td>
                <td class="py-3 px-4 text-right text-slate-600">Rp {{ formatNum(p.total_non_emas) }}</td>
                <td class="py-3 px-4 text-right font-bold text-slate-900">Rp {{ formatNum(p.total_omzet) }}</td>
                <td class="py-3 px-4 text-right text-amber-600 font-bold">{{ p.points }} Poin</td>
              </tr>
              <tr v-if="!leaderboard.top_performers || leaderboard.top_performers.length === 0">
                <td colspan="6" class="py-8 text-center text-slate-400 font-normal">Belum ada data leaderboard bulan ini</td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="pt-3 border-t border-slate-100 flex justify-end">
          <button @click="showLeaderboardModal = false" class="px-4 py-2 bg-slate-800 hover:bg-slate-900 text-white rounded-xl text-xs font-semibold transition-all cursor-pointer">
            Tutup
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from "vue";
import { useDashboardMetricsQuery } from "~/composables/useDashboardQueries";
import type { LeaderboardInfo } from "~/types/dashboard";

// Privacy toggle state (default censored)
const isPrivate = ref(true);

// Modal state
const showLeaderboardModal = ref(false);

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

// Extended state matching legacy Selada Incentive Calculator
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
const leaderboard = ref<LeaderboardInfo>({
  rank: 1,
  total_participants: 1,
  gap_to_next: 0,
  top_performers: [],
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
const liveFeeds = ref<any[]>([
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

const calcDrop = (curr: number, prev: number) => {
  if (!prev || prev <= 0) return "0 %";
  if (curr >= prev) return "0 %";
  const drop = Math.round(100 - (((prev - curr) / prev) * 100));
  return `${drop} %`;
};

const getFunnelWidth = (stage: number, count?: number, base?: number) => {
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

// TanStack Query integration
const {
  data: metricsData,
  isFetching,
  refetch,
} = useDashboardMetricsQuery(selectedSalesUser);

const handleReviewSales = () => {
  refetch();
};

const resetReview = () => {
  selectedSalesUser.value = "";
  refetch();
};

const activeSalesName = computed(() => {
  if (!selectedSalesUser.value) return salesName.value;
  const targetUser = selectedSalesUser.value.toLowerCase().trim();
  const s = salesList.value.find(
    (item) => item.username && item.username.toLowerCase().trim() === targetUser
  );
  if (!s) return selectedSalesUser.value.toUpperCase();
  return s.nm_jabatan ? `${s.nama} (${s.nm_jabatan})` : s.nama;
});

const getFeedCategory = (feed: any) => {
  if (feed.isBonus || feed.is_bonus || (feed.type && feed.type.toLowerCase().includes('bonus'))) {
    return 'bonus';
  }
  const typeStr = (feed.type || '').toLowerCase();
  const iconStr = (feed.iconText || '').toLowerCase();
  if (typeStr.includes('cif') || typeStr.includes('nasabah') || iconStr === 'cif') {
    return 'cif';
  }
  if (typeStr.includes('elektronik') || typeStr.includes('laptop') || typeStr.includes('hp') || iconStr === 'el' || typeStr.includes('non emas')) {
    return 'nonemas';
  }
  return 'emas';
};

const getFeedBadgeClass = (feed: any) => {
  const cat = getFeedCategory(feed);
  if (cat === 'bonus') return 'bg-amber-50 text-amber-600 border-amber-200/80';
  if (cat === 'cif') return 'bg-blue-50 text-blue-600 border-blue-200/80';
  if (cat === 'nonemas') return 'bg-sky-50 text-sky-600 border-sky-200/80';
  return 'bg-emerald-50 text-emerald-600 border-emerald-200/80';
};

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

<style scoped>
* {
  font-family: "Poppins", sans-serif;
}

.sales-funnel {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.sales-percent {
  display: flex;
  flex-direction: column;
  justify-content: space-around;
  padding-left: 0.5rem;
}
</style>
