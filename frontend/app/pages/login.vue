<template>
  <div class="min-h-screen grid grid-cols-1 lg:grid-cols-12 bg-slate-950 font-sans text-slate-100 overflow-hidden relative selection:bg-amber-500 selection:text-slate-950">
    <!-- Ambient Gold Glow Accents -->
    <div class="bg-amber-500/10 blur-[130px] rounded-full w-[450px] h-[450px] absolute -top-32 -left-32 pointer-events-none z-0"></div>
    <div class="bg-amber-600/10 blur-[150px] rounded-full w-[550px] h-[550px] absolute -bottom-40 -right-40 pointer-events-none z-0"></div>

    <!-- LEFT VISUAL HERO SHOWCASE (Desktop lg: 7 Columns) -->
    <div class="lg:col-span-7 xl:col-span-7 hidden lg:flex relative overflow-hidden flex-col justify-between p-12 lg:p-16 z-10 border-r border-slate-800/60">
      <!-- Background Image with Luxury Overlay -->
      <img
        src="/assets/img/banner.jpg"
        alt="Selada Visual Banner"
        class="absolute inset-0 w-full h-full object-cover object-center opacity-25 mix-blend-luminosity scale-105"
      />
      <div class="absolute inset-0 bg-gradient-to-tr from-slate-950 via-slate-950/85 to-amber-950/30"></div>

      <!-- Top Header Branding -->
      <div class="relative z-10 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="w-12 h-12 rounded-2xl bg-amber-500/10 border border-amber-500/30 backdrop-blur-md flex items-center justify-center p-2 shadow-lg shadow-amber-500/10">
            <img src="/assets/img/logo-gadaimulia.png" alt="Logo Gadai Mulia" class="w-full h-full object-contain" />
          </div>
          <div>
            <span class="text-xs font-bold uppercase tracking-widest text-amber-400 block">PT Gadai Mulia Indonesia</span>
            <span class="text-sm font-semibold text-slate-300">SELADA Advanced System</span>
          </div>
        </div>
        
        <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 text-xs font-medium">
          <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
          System v2.0 Operational
        </span>
      </div>

      <!-- Hero Headline & Subtitle -->
      <div class="relative z-10 my-auto max-w-xl py-12">
        <div class="inline-flex items-center gap-2 px-3 py-1.5 rounded-xl bg-amber-400/10 border border-amber-400/20 text-amber-300 text-xs font-semibold mb-6">
          <svg class="w-4 h-4 text-amber-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-6z" />
          </svg>
          Sistem Rencana & Estimasi Taksiran
        </div>

        <h1 class="text-4xl lg:text-5xl font-extrabold text-white tracking-tight leading-tight mb-4">
          Kelola Target & Taksiran Gadai Secara <span class="text-transparent bg-clip-text bg-gradient-to-r from-amber-300 via-amber-400 to-amber-500">Presisi & Real-Time</span>
        </h1>

        <p class="text-slate-400 text-base leading-relaxed mb-10">
          Akses terpadu kalkulator CeTar, pemantauan leads & rencana SBG cabang, serta laporan performa sales dalam satu platform yang responsif.
        </p>

        <!-- Feature Highlight Cards -->
        <div class="grid grid-cols-3 gap-4">
          <div class="p-4 rounded-2xl bg-slate-900/60 border border-slate-800/80 backdrop-blur-md">
            <div class="w-9 h-9 rounded-xl bg-amber-500/10 border border-amber-500/20 flex items-center justify-center text-amber-400 mb-3">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
              </svg>
            </div>
            <h4 class="text-sm font-bold text-white mb-1">CeTar Motor</h4>
            <p class="text-xs text-slate-400 leading-snug">Perhitungan taksiran otomatis & presisi</p>
          </div>

          <div class="p-4 rounded-2xl bg-slate-900/60 border border-slate-800/80 backdrop-blur-md">
            <div class="w-9 h-9 rounded-xl bg-amber-500/10 border border-amber-500/20 flex items-center justify-center text-amber-400 mb-3">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <h4 class="text-sm font-bold text-white mb-1">Input & Rencana SBG</h4>
            <p class="text-xs text-slate-400 leading-snug">Manajemen leads hingga pencairan</p>
          </div>

          <div class="p-4 rounded-2xl bg-slate-900/60 border border-slate-800/80 backdrop-blur-md">
            <div class="w-9 h-9 rounded-xl bg-amber-500/10 border border-amber-500/20 flex items-center justify-center text-amber-400 mb-3">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <h4 class="text-sm font-bold text-white mb-1">Sales Performance</h4>
            <p class="text-xs text-slate-400 leading-snug">Monitoring target cabang real-time</p>
          </div>
        </div>
      </div>

      <!-- Footer Info -->
      <div class="relative z-10 flex items-center justify-between text-xs text-slate-500 border-t border-slate-800/60 pt-6">
        <span>© 2026 PT Gadai Mulia Indonesia. Hak cipta dilindungi.</span>
        <span class="flex items-center gap-1">
          <svg class="w-3.5 h-3.5 text-amber-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
          </svg>
          256-Bit Encrypted Portal
        </span>
      </div>
    </div>

    <!-- RIGHT LOGIN FORM CARD (Desktop lg: 5 Columns) -->
    <div class="lg:col-span-5 xl:col-span-5 flex flex-col justify-between p-6 sm:p-12 z-10 bg-slate-900/80 backdrop-blur-2xl border-l border-slate-800/80 shadow-2xl relative min-h-screen lg:min-h-0">
      <div class="my-auto max-w-md w-full mx-auto">

        <!-- Mobile Branding Logo (Visible on mobile/tablet) -->
        <div class="flex lg:hidden items-center justify-center gap-3 mb-8">
          <div class="w-12 h-12 rounded-2xl bg-amber-500/10 border border-amber-500/30 flex items-center justify-center p-2 shadow-lg shadow-amber-500/10">
            <img src="/assets/img/logo-gadaimulia.png" alt="Logo Gadai Mulia" class="w-full h-full object-contain" />
          </div>
          <div>
            <span class="text-xs font-bold uppercase tracking-widest text-amber-400 block">PT Gadai Mulia Indonesia</span>
            <span class="text-base font-extrabold text-white tracking-tight">SELADA Advanced</span>
          </div>
        </div>

        <!-- Form Title Header -->
        <div class="mb-8 text-center lg:text-left">
          <div class="hidden lg:flex items-center gap-3 mb-6">
            <div class="w-12 h-12 rounded-2xl bg-amber-500/10 border border-amber-500/30 flex items-center justify-center p-2.5 shadow-lg shadow-amber-500/10">
              <img src="/assets/img/logo-gadaimulia.png" alt="Logo Gadai Mulia" class="w-full h-full object-contain" />
            </div>
            <div>
              <h1 class="text-xl font-extrabold text-white tracking-tight">
                SELADA <span class="text-amber-400">Advanced</span>
              </h1>
              <span class="text-xs font-semibold text-slate-400">Portal Otentikasi Pengguna</span>
            </div>
          </div>

          <h2 class="text-2xl sm:text-3xl font-black text-white tracking-tight mb-2">Selamat Datang</h2>
          <p class="text-slate-400 text-sm">Masukan kredensial akun Anda untuk mengakses sistem.</p>
        </div>

        <!-- Error Message Banner -->
        <transition enter-active-class="transition duration-200 ease-out" enter-from-class="transform opacity-0 -translate-y-2" enter-to-class="transform opacity-100 translate-y-0" leave-active-class="transition duration-150 ease-in" leave-from-class="transform opacity-100 translate-y-0" leave-to-class="transform opacity-0 -translate-y-2">
          <div v-if="errorMessage" class="mb-6 p-4 rounded-2xl bg-red-500/10 border border-red-500/30 text-red-300 text-sm flex items-start gap-3 shadow-lg shadow-red-500/5">
            <svg class="w-5 h-5 text-red-400 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <div class="flex-1">
              <span class="font-semibold block text-red-200">Gagal Masuk</span>
              <span>{{ errorMessage }}</span>
            </div>
          </div>
        </transition>

        <!-- Login Form -->
        <form @submit.prevent="handleLogin" class="space-y-5">
          <!-- Username Input Field -->
          <div class="space-y-2">
            <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider" for="username">
              Username
            </label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
              </div>
              <input
                id="username"
                v-model="username"
                @blur="validateUsername"
                @input="validateUsername"
                type="text"
                placeholder="Masukkan username Anda"
                class="w-full bg-slate-950/70 border text-white placeholder-slate-500 text-sm rounded-xl pl-11 pr-4 py-3.5 outline-none transition-all duration-200"
                :class="[
                  usernameError
                    ? 'border-red-500/80 focus:ring-2 focus:ring-red-500/20'
                    : 'border-slate-800 focus:border-amber-400 focus:ring-2 focus:ring-amber-400/20'
                ]"
                :disabled="isLoading"
                autocomplete="username"
              />
            </div>
            <p v-if="usernameError" class="text-xs text-red-400 font-medium flex items-center gap-1.5 mt-1">
              <svg class="w-3.5 h-3.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              {{ usernameError }}
            </p>
          </div>

          <!-- Password Input Field -->
          <div class="space-y-2">
            <div class="flex items-center justify-between">
              <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider" for="password">
                Password
              </label>
            </div>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
              </div>
              <input
                id="password"
                v-model="password"
                @blur="validatePassword"
                @input="validatePassword"
                :type="showPassword ? 'text' : 'password'"
                placeholder="••••••••••••"
                class="w-full bg-slate-950/70 border text-white placeholder-slate-500 text-sm rounded-xl pl-11 pr-11 py-3.5 outline-none transition-all duration-200"
                :class="[
                  passwordError
                    ? 'border-red-500/80 focus:ring-2 focus:ring-red-500/20'
                    : 'border-slate-800 focus:border-amber-400 focus:ring-2 focus:ring-amber-400/20'
                ]"
                :disabled="isLoading"
                autocomplete="current-password"
              />
              <button
                type="button"
                @click="togglePassword"
                class="absolute inset-y-0 right-0 pr-3.5 flex items-center text-slate-400 hover:text-amber-400 transition-colors cursor-pointer"
                :disabled="isLoading"
                tabindex="-1"
                :title="showPassword ? 'Sembunyikan password' : 'Tampilkan password'"
              >
                <svg v-if="!showPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.477 0 8.268 2.943 9.542 7-1.274 4.057-5.065 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
                <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M3 3l18 18" />
                </svg>
              </button>
            </div>
            <p v-if="passwordError" class="text-xs text-red-400 font-medium flex items-center gap-1.5 mt-1">
              <svg class="w-3.5 h-3.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              {{ passwordError }}
            </p>
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            class="w-full mt-2 bg-gradient-to-r from-amber-500 via-amber-400 to-amber-500 hover:from-amber-400 hover:to-amber-600 text-slate-950 font-bold py-3.5 px-6 rounded-xl shadow-lg shadow-amber-500/25 hover:shadow-amber-500/40 hover:-translate-y-0.5 active:translate-y-0 transition-all duration-200 flex items-center justify-center gap-2 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
            :disabled="isLoading || !!usernameError || !!passwordError"
          >
            <template v-if="isLoading">
              <svg class="animate-spin w-5 h-5 text-slate-950" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span>Memproses Otentikasi...</span>
            </template>
            <template v-else>
              <span>Masuk ke Sistem</span>
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
              </svg>
            </template>
          </button>
        </form>

        <!-- Security Footer Badge -->
        <div class="mt-8 text-center border-t border-slate-800/80 pt-6">
          <p class="text-xs text-slate-500 flex items-center justify-center gap-1.5">
            <svg class="w-4 h-4 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
            </svg>
            Koneksi Terenkripsi & Terlindungi SSL
          </p>
        </div>

      </div>

      <!-- Mobile Footer -->
      <div class="mt-8 text-center text-xs text-slate-600 lg:hidden">
        © 2026 PT Gadai Mulia Indonesia. All rights reserved.
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useCookie, navigateTo } from '#app'
import { useApi } from '~/composables/useApi'
import { parseJwt } from '~/utils/auth'

// Do not wrap login page with main default sidebar layout
definePageMeta({
  layout: false
})

const username = ref('')
const password = ref('')
const usernameError = ref('')
const passwordError = ref('')
const errorMessage = ref('')
const isLoading = ref(false)
const showPassword = ref(false)

// Declare session cookies at setup hook level for proper Nuxt reactivity
const token = useCookie('auth_token', { maxAge: 60 * 60 * 24 }) // 24 hours
const tokenExpiry = useCookie('token_expiry', { maxAge: 60 * 60 * 24 })
const name = useCookie('user_name')
const branch = useCookie('user_branch')
const fkUser = useCookie('fk_user')

const validateUsername = () => {
  const val = username.value.trim()
  if (!val) {
    usernameError.value = 'Username tidak boleh kosong'
    return false
  }
  if (val.length < 3) {
    usernameError.value = 'Username minimal harus 3 karakter'
    return false
  }
  if (/\s/.test(val)) {
    usernameError.value = 'Username tidak boleh mengandung spasi'
    return false
  }
  usernameError.value = ''
  return true
}

const validatePassword = () => {
  const val = password.value
  if (!val) {
    passwordError.value = 'Password tidak boleh kosong'
    return false
  }
  if (val.length < 6) {
    passwordError.value = 'Password minimal harus 6 karakter'
    return false
  }
  passwordError.value = ''
  return true
}

const togglePassword = () => {
  showPassword.value = !showPassword.value
}

const handleLogin = async () => {
  const isUsernameValid = validateUsername()
  const isPasswordValid = validatePassword()

  if (!isUsernameValid || !isPasswordValid) {
    errorMessage.value = 'Mohon periksa kembali username dan password Anda.'
    return
  }

  isLoading.value = true
  errorMessage.value = ''

  const { data, error } = await useApi('/login', {
    method: 'POST',
    body: {
      user: username.value.trim(),
      pass: password.value
    }
  })

  isLoading.value = false

  if (error) {
    errorMessage.value = typeof error === 'string' ? error : 'Username atau password Anda tidak sesuai.'
    return
  }

  if (data && data.token) {
    token.value = data.token
    
    // Store expiration time if exp or expiry claim exists in token
    const payload = parseJwt(data.token)
    let expiryMs = 0
    if (payload) {
      if (typeof payload.exp === 'number') {
        expiryMs = payload.exp * 1000
      } else if (payload.expiry) {
        const expMs = new Date(payload.expiry).getTime()
        if (!isNaN(expMs)) {
          expiryMs = expMs
        }
      }
    }
    if (!expiryMs || isNaN(expiryMs)) {
      expiryMs = Date.now() + 24 * 60 * 60 * 1000
    }
    tokenExpiry.value = String(expiryMs)

    name.value = data.nama || ''
    branch.value = data.cabang || ''
    fkUser.value = data.fk_user || ''

    if (typeof window !== 'undefined') {
      window.location.href = '/dashboard'
    } else {
      await navigateTo('/dashboard')
    }
  }
}
</script>
