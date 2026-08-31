<template>
  <div class="space-y-6 pb-12 font-sans">
    <!-- Header Section -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-slate-800 tracking-tight">Dashboard Ringkasan</h1>
        <p class="text-slate-500 text-sm">Selamat datang kembali, pantau performa Anda hari ini.</p>
      </div>
      <div class="flex items-center gap-3">
        <NuxtLink 
          to="/inquiry" 
          class="inline-flex items-center gap-2 bg-amber-500 hover:bg-amber-600 active:scale-95 text-white font-semibold px-5 py-2.5 rounded-xl text-sm transition-all shadow-md shadow-amber-200"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
          Leads Baru
        </NuxtLink>
      </div>
    </div>

    <!-- HERO & LEADERBOARD ROW -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Hero Card (2 cols) -->
      <div class="lg:col-span-2 bg-[#1e293b] text-white rounded-3xl p-6 md:p-8 shadow-xl relative overflow-hidden flex flex-col justify-between min-h-[220px]">
        <!-- Background Decorative Glow -->
        <div class="absolute -top-24 -right-24 w-72 h-72 bg-amber-500/10 rounded-full blur-3xl pointer-events-none"></div>
        <div class="absolute -bottom-24 -left-24 w-72 h-72 bg-emerald-500/10 rounded-full blur-3xl pointer-events-none"></div>

        <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-6 relative z-10">
          <div class="space-y-3 flex-1">
            <div class="flex items-center gap-3">
              <span class="text-xs font-bold uppercase tracking-wider text-amber-400">Estimasi Insentif Bulan Ini</span>
              
              <!-- Privacy Toggle Button -->
              <button 
                @click="isPrivate = !isPrivate"
                type="button"
                :class="isPrivate ? 'bg-rose-500/20 text-rose-300 border-rose-500/30' : 'bg-white/10 text-slate-300 border-white/20 hover:bg-white/20'"
                class="p-1.5 rounded-lg border text-xs flex items-center gap-1.5 transition-all cursor-pointer"
                :title="isPrivate ? 'Tampilkan Nilai' : 'Sembunyikan Nilai'"
              >
                <!-- Eye off icon when private -->
                <svg v-if="isPrivate" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858-5.908a10.04 10.04 0 013.682-.763c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m-2.115 2.115a3.5 3.5 0 11-4.95-4.95"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3l18 18"/></svg>
                <!-- Eye icon when visible -->
                <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/></svg>
              </button>
            </div>

            <!-- Incentive Display -->
            <div class="flex items-baseline gap-2">
              <span v-if="!isPrivate" class="text-2xl text-slate-400 font-light">Rp</span>
              <h2 class="text-4xl md:text-5xl font-black tracking-tight text-white select-none">
                <template v-if="isPrivate">******</template>
                <template v-else>{{ formatNum(estIncentive) }}</template>
              </h2>
            </div>

            <!-- Qualification Status Badge -->
            <div class="pt-1">
              <span v-if="isPrivate" class="text-sm font-bold text-slate-400">*****</span>
              <div 
                v-else 
                class="inline-flex items-center gap-2 px-3 py-1.5 rounded-xl text-xs font-bold"
                :class="isQualified ? 'bg-emerald-500/10 border border-emerald-500/20 text-emerald-400' : 'bg-amber-500/10 border border-amber-500/20 text-amber-400'"
              >
                <svg v-if="isQualified" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
                <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/></svg>
                <span>Status: {{ isQualified ? 'Qualified' : 'Dalam Proses' }} | Grade: {{ userGrade }} (Rate: {{ appliedRate }}%)</span>
              </div>
            </div>
          </div>

          <!-- Target Progress Card (Right Side inside Hero) -->
          <div class="w-full md:w-80 bg-slate-800/40 backdrop-blur-md rounded-2xl p-5 border border-white/10 space-y-3">
            <div class="flex justify-between items-center text-xs">
              <span class="text-slate-400 font-medium uppercase tracking-wider">Target Progress</span>
              <span class="font-bold text-white">{{ pointTotal }}/{{ targetGoalPoints }} Poin</span>
            </div>
            <div class="w-full bg-slate-700/50 rounded-full h-3 overflow-hidden p-0.5">
              <div 
                class="bg-amber-400 h-full rounded-full transition-all duration-1000 ease-out" 
                :style="{ width: `${Math.min(100, Math.max(0, targetProgressPct))}%` }"
              ></div>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-xs text-emerald-400 font-bold">{{ targetProgressPct }}%</span>
              <span class="text-[11px] text-slate-400">{{ targetDescription }}</span>
            </div>
            <div class="pt-1 border-t border-white/5 flex items-center gap-2 text-xs text-slate-300">
              <span class="animate-bounce text-sm">🚀</span>
              <span>{{ motivationText }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Leaderboard / Ranking Card (1 col) -->
      <div class="bg-white rounded-3xl p-6 border border-slate-200 shadow-sm flex flex-col items-center justify-center text-center relative group overflow-hidden">
        <div class="absolute top-0 right-0 p-4 opacity-5 group-hover:scale-110 transition-transform duration-500">
          <svg class="w-24 h-24 text-slate-900" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z"/></svg>
        </div>
        <div class="w-16 h-16 bg-amber-50 text-amber-600 rounded-2xl flex items-center justify-center mb-3 shadow-inner rotate-3 group-hover:rotate-0 transition-transform">
          <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z"/></svg>
        </div>
        <h2 class="text-amber-600 font-bold uppercase tracking-wide text-sm mb-1">{{ salesName }}</h2>
        <h3 class="text-slate-400 font-semibold uppercase tracking-widest text-xs mb-2">Posisi Leaderboard</h3>
        <div class="text-3xl md:text-4xl font-black text-slate-800 mb-2">
          #{{ leaderboard.rank }} <span class="text-sm font-normal text-slate-500">dari {{ leaderboard.total_participants }} Sales</span>
        </div>
        <p class="text-xs text-slate-500 max-w-[240px] leading-relaxed">
          <template v-if="leaderboard.rank === 1">
            <strong class="text-emerald-600 font-bold">Luar biasa!</strong> Anda memimpin peringkat #1! 🎉
          </template>
          <template v-else>
            Butuh <strong class="text-emerald-600 font-bold">Rp {{ formatNum(leaderboard.gap_to_next) }}</strong> omzet lagi untuk menyalip ke peringkat #{{ leaderboard.rank - 1 }}.
          </template>
        </p>
      </div>
    </div>

    <!-- METRICS ROW -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- Card 1: Emas -->
      <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-sm hover:shadow-md transition-all border-l-4 border-l-amber-500 flex flex-col justify-between">
        <div>
          <div class="flex justify-between items-start mb-3">
            <div class="p-2.5 bg-amber-50 text-amber-600 rounded-xl flex items-center justify-center shadow-inner">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
            </div>
            <span class="text-[10px] font-extrabold text-amber-700 bg-amber-100/80 px-2.5 py-1 rounded-full uppercase tracking-wider">EMAS</span>
          </div>
          <h4 class="text-slate-500 text-xs font-bold uppercase tracking-wider mb-1">Pencairan Gadai Emas</h4>
          <div class="text-2xl font-black text-slate-800 mb-1 tracking-tight">Rp {{ formatNum(stats.emas) }}</div>
          <div class="text-xs text-slate-500 font-semibold mb-3">
            Total Gramasi: <strong class="text-amber-600 font-bold">{{ gramasiEmas }}g</strong> / {{ targetEmasG }}g
          </div>
        </div>

        <div class="pt-3 border-t border-slate-100 space-y-2">
          <div class="flex items-center justify-between text-[11px]">
            <span class="font-extrabold text-amber-600 flex items-center gap-1">
              <span>🎁</span> Tier Bonus: 300g = 1 Jt
            </span>
            <span class="font-bold text-slate-600">{{ pctEmas }}%</span>
          </div>
          <div class="w-full bg-slate-100 rounded-full h-2 overflow-hidden">
            <div class="h-full rounded-full bg-gradient-to-r from-amber-400 to-amber-500 transition-all duration-1000" :style="{ width: `${pctEmas}%` }"></div>
          </div>
          <p class="text-[11px] text-slate-400 font-medium">Target berikutnya: 500g (Bonus 2 Jt)</p>
        </div>
      </div>

      <!-- Card 2: Non Emas / Elektronik -->
      <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-sm hover:shadow-md transition-all border-l-4 border-l-emerald-500 flex flex-col justify-between">
        <div>
          <div class="flex justify-between items-start mb-3">
            <div class="p-2.5 bg-sky-50 text-sky-600 rounded-xl flex items-center justify-center shadow-inner">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/></svg>
            </div>
            <span class="text-[10px] font-extrabold text-sky-700 bg-sky-100/80 px-2.5 py-1 rounded-full uppercase tracking-wider">NON EMAS</span>
          </div>
          <h4 class="text-slate-500 text-xs font-bold uppercase tracking-wider mb-1">Pencairan Gadai Non Emas</h4>
          <div class="text-2xl font-black text-slate-800 mb-1 tracking-tight">Rp {{ formatNum(stats.nonemas) }}</div>
          <div class="text-xs text-slate-500 font-semibold mb-3">
            Target UP: <strong class="text-sky-600 font-bold">50 Jt</strong> (Bonus 1 Jt)
          </div>
        </div>

        <div class="pt-3 border-t border-slate-100 space-y-2">
          <div class="flex items-center justify-between text-[11px]">
            <span class="font-extrabold text-sky-600 flex items-center gap-1">
              <span>🚀</span> Progress UP Non-Emas
            </span>
            <span class="font-bold text-slate-600">{{ pctNonEmas }}%</span>
          </div>
          <div class="w-full bg-slate-100 rounded-full h-2 overflow-hidden">
            <div class="h-full rounded-full bg-gradient-to-r from-sky-400 to-sky-500 transition-all duration-1000" :style="{ width: `${pctNonEmas}%` }"></div>
          </div>
          <p class="text-[11px] text-slate-400 font-medium">Realisasi omzet non-emas bulan ini</p>
        </div>
      </div>

      <!-- Card 3: Nasabah Baru (New CIF) -->
      <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-sm hover:shadow-md transition-all border-l-4 border-l-blue-500 flex flex-col justify-between">
        <div>
          <div class="flex justify-between items-start mb-4">
            <div class="p-2.5 bg-blue-50 text-blue-600 rounded-xl flex items-center justify-center shadow-inner">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"/></svg>
            </div>
            <span class="text-[10px] font-bold text-blue-600 bg-blue-50 px-2 py-0.5 rounded-full uppercase tracking-wider">Nasabah</span>
          </div>
          <h4 class="text-slate-500 text-xs font-bold uppercase tracking-wider mb-1">New CIF (Nasabah Baru)</h4>
          <div class="text-2xl font-black text-slate-800 mb-1 tracking-tight">{{ stats.newcif }}</div>
        </div>
        <div class="pt-3 border-t border-slate-100 text-[11px] font-semibold text-blue-600 flex items-center gap-1.5">
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"/></svg>
          CIF Bulan Lalu: {{ lastMonthCif }}
        </div>
      </div>

      <!-- Card 4: New Leads -->
      <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-sm hover:shadow-md transition-all border-l-4 border-l-rose-500 flex flex-col justify-between">
        <div>
          <div class="flex justify-between items-start mb-4">
            <div class="p-2.5 bg-rose-50 text-rose-600 rounded-xl flex items-center justify-center shadow-inner">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/></svg>
            </div>
            <span class="text-[10px] font-bold text-rose-600 bg-rose-50 px-2 py-0.5 rounded-full uppercase tracking-wider">Leads</span>
          </div>
          <h4 class="text-slate-500 text-xs font-bold uppercase tracking-wider mb-1">New Leads</h4>
          <div class="text-2xl font-black text-slate-800 mb-1 tracking-tight">{{ stats.leads_monthly || stats.leads }}</div>
        </div>
        <div class="pt-3 border-t border-slate-100 text-[11px] font-semibold text-rose-600 flex items-center gap-1.5">
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
          Leads Bulan Lalu: {{ lastMonthLeads }}
        </div>
      </div>
    </div>

    <!-- SALES FUNNEL SECTION (PRESERVED) -->
    <div class="bg-white rounded-3xl border border-slate-200 p-6 shadow-sm">
      <div class="flex items-center gap-2 mb-6">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#64748B" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"></polygon>
        </svg>
        <h3 class="font-bold text-slate-800 text-lg">Sales Funnel Progress</h3>
      </div>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Akumulasi Keseluruhan -->
        <div class="bg-slate-50 rounded-2xl p-5 border border-slate-100 space-y-4">
          <h4 class="text-xs font-bold uppercase tracking-wider text-slate-500 text-center">Akumulasi Keseluruhan</h4>
          <div class="space-y-2">
            <div class="flex items-center justify-between text-xs font-semibold text-slate-700">
              <span>Leads</span>
              <span>{{ stats.leads }}</span>
            </div>
            <div class="w-full bg-slate-200/60 rounded-lg h-9 p-1 flex items-center">
              <div class="bg-rose-600 h-full rounded-md flex items-center justify-center text-white text-xs font-bold transition-all duration-500" style="width: 100%">
                Leads - {{ stats.leads }}
              </div>
            </div>

            <div class="flex items-center justify-between text-xs font-semibold text-slate-700 pt-1">
              <span>Prospek</span>
              <span>{{ stats.prospect }}</span>
            </div>
            <div class="w-full bg-slate-200/60 rounded-lg h-9 p-1 flex items-center">
              <div class="bg-rose-800 h-full rounded-md flex items-center justify-center text-white text-xs font-bold transition-all duration-500" :style="{ width: getWidth(stats.prospect, stats.leads) }">
                Prospek - {{ stats.prospect }}
              </div>
            </div>

            <div class="flex items-center justify-between text-xs font-semibold text-slate-700 pt-1">
              <span>Hot Prospek</span>
              <span>{{ stats.hotprospect }}</span>
            </div>
            <div class="w-full bg-slate-200/60 rounded-lg h-9 p-1 flex items-center">
              <div class="bg-amber-500 h-full rounded-md flex items-center justify-center text-white text-xs font-bold transition-all duration-500" :style="{ width: getWidth(stats.hotprospect, stats.leads) }">
                Hot Prospek - {{ stats.hotprospect }}
              </div>
            </div>

            <div class="flex items-center justify-between text-xs font-semibold text-slate-700 pt-1">
              <span>DO (SBG)</span>
              <span>{{ stats.sbg }}</span>
            </div>
            <div class="w-full bg-slate-200/60 rounded-lg h-9 p-1 flex items-center">
              <div class="bg-emerald-500 h-full rounded-md flex items-center justify-center text-white text-xs font-bold transition-all duration-500" :style="{ width: getWidth(stats.sbg, stats.leads) }">
                DO - {{ stats.sbg }}
              </div>
            </div>
          </div>
        </div>

        <!-- Bulan Ini -->
        <div class="bg-slate-50 rounded-2xl p-5 border border-slate-100 space-y-4">
          <h4 class="text-xs font-bold uppercase tracking-wider text-slate-500 text-center">Bulan Ini</h4>
          <div class="space-y-2">
            <div class="flex items-center justify-between text-xs font-semibold text-slate-700">
              <span>Leads</span>
              <span>{{ stats.leads_monthly || stats.leads }}</span>
            </div>
            <div class="w-full bg-slate-200/60 rounded-lg h-9 p-1 flex items-center">
              <div class="bg-rose-600 h-full rounded-md flex items-center justify-center text-white text-xs font-bold transition-all duration-500" style="width: 100%">
                Leads - {{ stats.leads_monthly || stats.leads }}
              </div>
            </div>

            <div class="flex items-center justify-between text-xs font-semibold text-slate-700 pt-1">
              <span>Prospek</span>
              <span>{{ stats.prospect }}</span>
            </div>
            <div class="w-full bg-slate-200/60 rounded-lg h-9 p-1 flex items-center">
              <div class="bg-rose-800 h-full rounded-md flex items-center justify-center text-white text-xs font-bold transition-all duration-500" :style="{ width: getWidth(stats.prospect, stats.leads_monthly || stats.leads) }">
                Prospek - {{ stats.prospect }}
              </div>
            </div>

            <div class="flex items-center justify-between text-xs font-semibold text-slate-700 pt-1">
              <span>Hot Prospek</span>
              <span>{{ stats.hotprospect }}</span>
            </div>
            <div class="w-full bg-slate-200/60 rounded-lg h-9 p-1 flex items-center">
              <div class="bg-amber-500 h-full rounded-md flex items-center justify-center text-white text-xs font-bold transition-all duration-500" :style="{ width: getWidth(stats.hotprospect, stats.leads_monthly || stats.leads) }">
                Hot Prospek - {{ stats.hotprospect }}
              </div>
            </div>

            <div class="flex items-center justify-between text-xs font-semibold text-slate-700 pt-1">
              <span>DO (SBG)</span>
              <span>{{ stats.sbg }}</span>
            </div>
            <div class="w-full bg-slate-200/60 rounded-lg h-9 p-1 flex items-center">
              <div class="bg-emerald-500 h-full rounded-md flex items-center justify-center text-white text-xs font-bold transition-all duration-500" :style="{ width: getWidth(stats.sbg, stats.leads_monthly || stats.leads) }">
                DO - {{ stats.sbg }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- BOTTOM ROW: PACING CHART & LIVE FEED -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Pacing Chart (2 cols) -->
      <div class="lg:col-span-2 bg-white rounded-3xl border border-slate-200 p-6 shadow-sm flex flex-col relative">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-6">
          <div>
            <h3 class="font-bold text-slate-800 text-lg">Kinerja Harian</h3>
            <p class="text-xs text-slate-500">Statistik realisasi dibandingkan target kumulatif harian.</p>
          </div>
          <div class="flex items-center gap-4 text-xs">
            <div class="flex items-center gap-1.5">
              <span class="w-2.5 h-2.5 rounded-full bg-emerald-500"></span>
              <span class="font-bold text-slate-600">Pacing Selada (Submit)</span>
            </div>
            <div class="flex items-center gap-1.5">
              <span class="w-2.5 h-2.5 rounded-full bg-sky-500"></span>
              <span class="font-bold text-slate-600">Pacing Sopiga Pro (Pencairan)</span>
            </div>
          </div>
        </div>

        <!-- Pacing Chart Bars Container -->
        <div class="flex-1 min-h-[220px] flex items-end justify-between gap-1 pt-8 pb-2 relative">
          <div 
            v-for="d in daysInMonth" 
            :key="d"
            class="flex-1 flex flex-col items-center justify-end h-full group cursor-pointer relative"
            :class="d === currentDay ? 'ring-2 ring-amber-400 ring-offset-1 z-10 rounded-t-sm' : ''"
          >
            <!-- Tooltip -->
            <div 
              class="absolute bottom-full mb-2 bg-slate-800 text-white text-[10px] p-2.5 rounded-xl shadow-2xl opacity-0 group-hover:opacity-100 transition-opacity z-50 whitespace-nowrap pointer-events-none space-y-1"
              :class="d <= 4 ? 'left-0' : (d >= daysInMonth - 4 ? 'right-0' : 'left-1/2 -translate-x-1/2')"
            >
              <div class="font-bold border-b border-slate-700 pb-1 text-amber-400">
                Tgl {{ d }} {{ d === currentDay ? '(Hari Ini)' : '' }}
              </div>
              <div class="flex items-center gap-1.5 text-emerald-400 font-semibold">
                <span class="w-2 h-2 rounded-full bg-emerald-400"></span>
                <span>Selada: Rp {{ formatNum(getPacingVal(d, 'selada')) }}</span>
              </div>
              <div class="flex items-center gap-1.5 text-sky-400 font-semibold">
                <span class="w-2 h-2 rounded-full bg-sky-400"></span>
                <span>Sopiga Pro: Rp {{ formatNum(getPacingVal(d, 'sopiga')) }}</span>
              </div>
            </div>

            <!-- Dual Bar -->
            <div class="w-full flex items-end justify-center gap-[2px] h-full px-[1px]">
              <div 
                class="w-1/2 rounded-t-sm transition-all duration-500" 
                :class="getPacingVal(d, 'selada') > 0 ? 'bg-emerald-500 shadow-sm' : (d <= currentDay ? 'bg-emerald-200/60' : 'bg-slate-100')"
                :style="{ height: `${getPacingHeight(d, 'selada')}%`, opacity: d <= currentDay ? 1 : 0.4 }"
              ></div>
              <div 
                class="w-1/2 rounded-t-sm transition-all duration-500"
                :class="getPacingVal(d, 'sopiga') > 0 ? 'bg-sky-500 shadow-sm' : (d <= currentDay ? 'bg-sky-200/60' : 'bg-slate-100')"
                :style="{ height: `${getPacingHeight(d, 'sopiga')}%`, opacity: d <= currentDay ? 1 : 0.4 }"
              ></div>
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
          <h3 class="font-bold text-slate-800 text-lg flex items-center gap-2">
            Live Insentif
            <span class="flex h-2 w-2 relative">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
            </span>
          </h3>
        </div>

        <div class="flex-1 space-y-3 overflow-y-auto max-h-[350px] pr-1">
          <div 
            v-for="(feed, idx) in liveFeeds" 
            :key="idx"
            class="p-3 rounded-xl bg-slate-50 border border-slate-100 flex items-center justify-between group hover:bg-white hover:shadow-md transition-all"
          >
            <div class="flex items-center gap-3">
              <div 
                class="w-9 h-9 rounded-full flex items-center justify-center shadow-sm border text-xs font-bold"
                :class="feed.isBonus ? 'bg-amber-50 text-amber-600 border-amber-200' : 'bg-emerald-50 text-emerald-600 border-emerald-200'"
              >
                <span>{{ feed.iconText }}</span>
              </div>
              <div>
                <p class="text-xs font-bold text-slate-800">{{ feed.type }}</p>
                <p class="text-[10px] text-slate-500 uppercase tracking-tight">{{ feed.code }} • {{ feed.customer }}</p>
              </div>
            </div>
            <div class="text-right">
              <p class="text-xs font-black text-emerald-600">+Rp {{ formatNum(feed.amount) }}</p>
              <p class="text-[9px] text-amber-600 font-semibold">{{ feed.sub }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'

// Privacy toggle state (default censored)
const isPrivate = ref(true)

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
  leads_monthly: 0
})

// Extended state matching v_dashboard.php
const estIncentive = ref(1250000)
const isQualified = ref(true)
const userGrade = ref('Grade A')
const appliedRate = ref(100)
const pointTotal = ref(45)
const targetGoalPoints = ref(100)
const targetDescription = ref('Target Point Harian')
const motivationText = ref('Semangat! Kejar 10 Poin lagi untuk Bonus Emas! 🚀')

const salesName = ref('ADI TYA')
const leaderboard = ref({
  rank: 1,
  total_participants: 12,
  gap_to_next: 2500000
})

const gramasiEmas = ref(12.5)
const targetEmasG = ref(300)
const lastMonthCif = ref(3)
const lastMonthLeads = ref(8)

// Computed helpers
const targetProgressPct = computed(() => {
  if (!targetGoalPoints.value) return 0
  return Math.round((pointTotal.value / targetGoalPoints.value) * 100)
})

const pctEmas = computed(() => {
  if (!targetEmasG.value) return 0
  return Math.min(100, Math.round((gramasiEmas.value / targetEmasG.value) * 100))
})

const pctNonEmas = computed(() => {
  const targetUP = 50000000
  if (!stats.value.nonemas) return 25
  return Math.min(100, Math.round((stats.value.nonemas / targetUP) * 100))
})

// Date calculations for Pacing Chart
const now = new Date()
const currentDay = now.getDate()
const daysInMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0).getDate()

// Simulated daily pacing chart data
const pacingSelada = ref<Record<number, number>>({ 1: 1500000, 3: 2000000, 5: 3500000, 8: 1200000, 12: 4000000, [currentDay]: 2500000 })
const pacingSopiga = ref<Record<number, number>>({ 2: 1000000, 4: 2500000, 6: 1800000, 10: 3000000, 14: 2200000, [currentDay]: 1800000 })

const getPacingVal = (day: number, type: 'selada' | 'sopiga') => {
  const data = type === 'selada' ? pacingSelada.value : pacingSopiga.value
  return data[day] || 0
}

const getPacingHeight = (day: number, type: 'selada' | 'sopiga') => {
  const val = getPacingVal(day, type)
  const maxVal = 5000000
  if (day > currentDay) return 8
  if (val === 0) return 12
  return Math.min(100, Math.max(15, (val / maxVal) * 85 + 15))
}

// Simulated live feed events
const liveFeeds = ref([
  { type: 'Gadai Emas', code: 'TRX-8412', customer: 'Budi Santoso', amount: 125000, sub: '2 Poin • Pinjaman Rp 5.000.000', iconText: '💎', isBonus: false },
  { type: 'Nasabah Baru', code: 'TRX-7741', customer: 'Siti Rahma', amount: 25000, sub: 'New CIF Bonus', iconText: '👤', isBonus: false },
  { type: 'Gadai Elektronik', code: 'TRX-6109', customer: 'Ahmad Fauzi', amount: 50000, sub: '1 Poin • Laptop Asus', iconText: '💻', isBonus: false },
  { type: 'Bonus Gramasi', code: 'BONUS-101', customer: 'Achievement 300g', amount: 1000000, sub: 'Target Emas Level 1', iconText: '🏆', isBonus: true }
])

// Formatting helpers
const formatNum = (val: number) => {
  if (!val && val !== 0) return '0'
  return new Intl.NumberFormat('id-ID').format(val)
}

const getWidth = (count: number, base: number) => {
  if (!base || count > base) return '100%'
  return `${Math.max(25, (count / base) * 100)}%`
}

onMounted(async () => {
  try {
    const { data } = await useApi('/inquiry/counts')
    if (data) {
      stats.value = { ...stats.value, ...data }
    }
  } catch (err) {
    console.error('Failed to load inquiry counts:', err)
  }
})
</script>
