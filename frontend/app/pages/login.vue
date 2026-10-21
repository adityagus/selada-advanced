<template>
  <div class="min-h-screen grid grid-cols-1 lg:grid-cols-12 bg-slate-900 font-sans text-slate-100">
    <!-- LEFT SIDE: BRAND & HERO BANNER (Desktop Only) -->
    <div class="lg:col-span-7 hidden lg:flex relative overflow-hidden flex-col justify-between p-12 bg-slate-950 border-r border-slate-800">
      <!-- Background Banner Image -->
      <img
        src="/assets/img/banner.jpg"
        alt="Gadai Mulia Banner"
        class="absolute inset-0 w-full h-full object-cover opacity-20"
      />
      <div class="absolute inset-0 bg-gradient-to-t from-slate-950 via-slate-950/10 to-transparent"></div>

      <!-- Top Brand Header -->
      <div class="relative z-10 flex items-center gap-3">
        <img src="/assets/img/logo-gadaimulia.png" alt="Logo Gadai Mulia" class="w-10 h-10 object-contain" />
        <div>
          <h2 class="font-bold text-white text-base leading-tight">PT Gadai Mulia Sejahtera</h2>
          <p class="text-xs text-slate-400">SELADA Advanced</p>
        </div>
      </div>

      <!-- Hero Statement -->


      <!-- Footer Copyright -->
      <div>
              <div class="relative z-10 max-w-lg mb-2">
        <h1 class="text-2xl lg:text-3xl font-bold text-white tracking-tight leading-snug mb-3">
          Sales Lead & Data Acquisition
        </h1>
        <p class="text-slate-400 text-xs leading-relaxed">
          Platform manajemen operasional cabang, perhitunganCeTar, dan pencapaian target SBG Gadai Mulia.
        </p>
      </div>
      <span class="relative z-10 text-xs text-slate-500">
        © 2026 PT Gadai Mulia Sejahtera. All rights reserved.
      </span>
      </div>
    </div>

    <!-- RIGHT SIDE: LOGIN FORM CARD -->
    <div class="lg:col-span-5 flex flex-col justify-between p-6 sm:p-12 bg-slate-50 border-l border-slate-800/60 min-h-screen lg:min-h-0">
      <div class="my-auto max-w-sm w-full mx-auto">
        <!-- Logo & Header -->
        <div class="mb-8">
          <div class="flex items-center gap-3 mb-6">
            <img src="/assets/img/logo-gadaimulia.png" alt="Logo Gadai Mulia" class="w-10 h-10 object-contain" />
            <div>
              <span class="text-xs font-semibold text-amber-500 block uppercase tracking-wider">Gadai Mulia</span>
              <span class="text-xl font-bold text-slate-800 tracking-tight">SELADA <span class="text-amber-500 font-normal">Advanced</span></span>
            </div>
          </div>
        </div>  

        <!-- Error Message Banner -->
        <div v-if="errorMessage" role="alert" class="mb-6 p-3.5 rounded-xl bg-red-500/10 border border-red-500/30 text-red-600 text-xs flex items-center gap-2.5">
          <AppIcon name="warning" class="w-4 h-4 flex-shrink-0" />
          <span>{{ errorMessage }}</span>
        </div>

        <!-- Form -->
        <form @submit.prevent="handleLogin" class="space-y-4">
          <div>
            <label class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1.5" for="username">
              Username
            </label>
            <input
              id="username"
              v-model="username"
              @blur="validateUsername"
              @input="validateUsername"
              type="text"
              placeholder="Username"
              class="w-full bg-white border text-slate-950 placeholder-slate-500 text-sm rounded-xl px-4 py-3 outline-none transition-colors"
              :class="[usernameError ? 'border-red-500' : 'border-slate-800 focus:border-amber-500']"
              :disabled="isLoading"
              autocomplete="username"
            />
            <span v-if="usernameError" class="text-xs text-red-600 mt-1 block">{{ usernameError }}</span>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1.5" for="password">
              Password
            </label>
            <div class="relative">
              <input
                id="password"
                v-model="password"
                @blur="validatePassword"
                @input="validatePassword"
                :type="showPassword ? 'text' : 'password'"
                placeholder="••••••••"
                class="w-full bg-white border text-slate-950 placeholder-slate-500 text-sm rounded-xl pl-4 pr-10 py-3 outline-none transition-colors"
                :class="[passwordError ? 'border-red-500' : 'border-slate-800 focus:border-amber-500']"
                :disabled="isLoading"
                autocomplete="current-password"
              />
              <button
                type="button"
                @click="togglePassword"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-500 hover:text-slate-800 transition-colors cursor-pointer"
                :disabled="isLoading"
                :aria-label="showPassword ? 'Sembunyikan password' : 'Tampilkan password'"
                :aria-pressed="showPassword"
              >
                <AppIcon name="eye" v-if="!showPassword" class="w-4 h-4" />
                <AppIcon name="eye-off" v-else class="w-4 h-4" />
              </button>
            </div>
            <span v-if="passwordError" class="text-xs text-red-600 mt-1 block">{{ passwordError }}</span>
          </div>

          <button
            type="submit"
            class="w-full mt-4 bg-amber-500 hover:bg-amber-600 text-slate-950 font-normal py-3 px-4 rounded-xl transition-colors cursor-pointer flex items-center justify-center gap-2 disabled:opacity-50"
            :disabled="isLoading || !!usernameError || !!passwordError"
          >
            <AppIcon name="loader" v-if="isLoading" class="w-4 h-4  text-slate-950   animate-spin" />
            <span>{{ isLoading ? 'Memproses...' : 'Masuk' }}</span>
          </button>
        </form>
      </div>

      <!-- Mobile Copyright Footer -->
      <div class="text-center text-xs text-slate-500 lg:hidden mt-8">
        © 2026 PT Gadai Mulia Sejahtera
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useCookie } from '#app'
import { useApi } from '~/composables/useApi'
import { parseJwt, SESSION_IDLE_TIMEOUT_MS, SESSION_MAX_AGE_SEC } from '~/utils/auth'
import { loginSchema } from '~/schemas/authSchema'

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

// if (typeof window !== 'undefined' && window.location?.search?.includes('reason=timeout')) {
//   errorMessage.value = 'Sesi Anda telah berakhir karena tidak ada aktivitas selama 30 menit. Silakan masuk kembali.'
// }

const token = useCookie('auth_token', { maxAge: SESSION_MAX_AGE_SEC })
const tokenExpiry = useCookie('token_expiry', { maxAge: SESSION_MAX_AGE_SEC })
const name = useCookie('user_name', { maxAge: SESSION_MAX_AGE_SEC })
const branch = useCookie('user_branch', { maxAge: SESSION_MAX_AGE_SEC })
const branchName = useCookie('user_branch_name', { maxAge: SESSION_MAX_AGE_SEC })
const fkUser = useCookie('fk_user', { maxAge: SESSION_MAX_AGE_SEC })

const persistSession = (session: Record<string, string>) => {
  const maxAge = SESSION_MAX_AGE_SEC
  const secure = typeof window !== 'undefined' && window.location.protocol === 'https:' ? '; Secure' : ''

  Object.entries(session).forEach(([key, value]) => {
    document.cookie = `${key}=${encodeURIComponent(value)}; Path=/; Max-Age=${maxAge}; SameSite=Lax${secure}`
  })
}

const validateUsername = () => {
  const val = username.value.trim()
  if (!val) {
    usernameError.value = 'Username tidak boleh kosong'
    return false
  }
  if (val.length < 3) {
    usernameError.value = 'Username minimal 3 karakter'
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
  const result = loginSchema.shape.password.safeParse(val)
  passwordError.value = result.success ? '' : result.error.issues[0]?.message || 'Password tidak valid'
  return result.success
}

const togglePassword = () => {
  showPassword.value = !showPassword.value
}

const handleLogin = async () => {
  if (isLoading.value) return
  const isUsernameValid = validateUsername()
  const isPasswordValid = validatePassword()

  if (!isUsernameValid || !isPasswordValid) {
    errorMessage.value = 'Mohon periksa kembali username dan password Anda.'
    return
  }

  isLoading.value = true
  errorMessage.value = ''

  try {
    const { data, error } = await useApi('/login', {
      method: 'POST',
      body: {
        user: username.value.trim(),
        pass: password.value
      }
    })

    if (error) {
      errorMessage.value = typeof error === 'string' ? error : 'Username atau password Anda tidak sesuai.'
      return
    }

    if (data && typeof data.token === 'string' && data.token) {
      const payload = parseJwt(data.token)
      let jwtExpMs = 0
      if (payload) {
        if (typeof payload.exp === 'number') {
          jwtExpMs = payload.exp * 1000
        } else if (payload.expiry) {
          const expMs = new Date(payload.expiry).getTime()
          if (!isNaN(expMs)) {
            jwtExpMs = expMs
          }
        }
      }

      const idleExpiry = Date.now() + SESSION_IDLE_TIMEOUT_MS
      const expiryMs = (jwtExpMs && jwtExpMs > Date.now() && jwtExpMs < idleExpiry)
        ? jwtExpMs
        : idleExpiry

      token.value = data.token
      tokenExpiry.value = String(expiryMs)

      name.value = data.nama || ''
      branch.value = data.cabang || ''
      branchName.value = data.branch_name || ''
      fkUser.value = data.fk_user || ''

      if (typeof localStorage !== 'undefined') {
        try {
          localStorage.setItem('selada_last_activity', String(Date.now()))
        } catch {}
      }

      // Write directly so the authentication middleware can read the session during
      // the dashboard page load, then use a full browser navigation for reliability.
      if (typeof window !== 'undefined') {
        persistSession({
          auth_token: data.token,
          token_expiry: String(expiryMs),
          user_name: String(data.nama || ''),
          user_branch: String(data.cabang || ''),
          user_branch_name: String(data.branch_name || ''),
          fk_user: String(data.fk_user || '')
        })
        window.location.replace('/dashboard')
      }
    } else {
      errorMessage.value = 'Respons login dari server tidak valid. Silakan coba lagi.'
    }
  } catch {
    errorMessage.value = 'Login gagal diproses. Silakan coba lagi.'
  } finally {
    isLoading.value = false
  }
}
</script>
