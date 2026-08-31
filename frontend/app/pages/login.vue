<template>
  <div class="login-wrapper">
    <div class="login-decor-1"></div>
    <div class="login-decor-2"></div>

    <div class="login-card">
      <div class="login-header">
        <div class="login-logo">
          <svg width="40" height="40" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="#f59e0b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M2 17L12 22L22 17" stroke="#f59e0b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M2 12L12 17L22 12" stroke="#f59e0b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
        <h1>SELADA <span>V2</span></h1>
        <p class="text-muted">Masuk untuk mengakses sistem taksiran & rencana</p>
      </div>

      <!-- Error message banner -->
      <div v-if="errorMessage" class="error-banner">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><path d="M12 9V13M12 17H12.01M21 12C21 16.9706 16.9706 21 12 21C7.02944 21 3 16.9706 3 12C3 7.02944 7.02944 3 12 3C16.9706 3 21 7.02944 21 12Z" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
        {{ errorMessage }}
      </div>

      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label class="form-label" for="username">Username</label>
          <input
            id="username"
            v-model="username"
            @blur="validateUsername"
            @input="validateUsername"
            type="text"
            placeholder="Ketik username Anda"
            class="form-input"
            :class="{ 'has-error': usernameError }"
            :disabled="isLoading"
          />
          <span v-if="usernameError" class="error-text">{{ usernameError }}</span>
        </div>

        <div class="form-group">
          <label class="form-label" for="password">Password</label>
          <div class="password-container">
            <input
              id="password"
              v-model="password"
              @blur="validatePassword"
              @input="validatePassword"
              :type="showPassword ? 'text' : 'password'"
              placeholder="••••••••••••"
              class="form-input password-field"
              :class="{ 'has-error': passwordError }"
              :disabled="isLoading"
            />
            <button type="button" @click="togglePassword" class="password-toggle" :disabled="isLoading" tabindex="-1">
              <svg v-if="!showPassword" width="18" height="18" viewBox="0 0 24 24" fill="none"><path d="M2.458 12C3.732 7.943 7.523 5 12 5c4.477 0 8.268 2.943 9.542 7-1.274 4.057-5.065 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" stroke="currentColor" stroke-width="2"/><path d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" stroke="currentColor" stroke-width="2"/></svg>
              <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none"><path d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M3 3l18 18" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
            </button>
          </div>
          <span v-if="passwordError" class="error-text">{{ passwordError }}</span>
        </div>

        <button type="submit" class="btn btn-primary login-btn" :disabled="isLoading || !!usernameError || !!passwordError">
          <span v-if="isLoading" class="spinner-small"></span>
          <span v-else>Masuk Sekarang</span>
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useCookie, navigateTo } from '#app'
import { useApi } from '~/composables/useApi'

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
    errorMessage.value = 'Silakan perbaiki kesalahan input sebelum masuk.'
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
    errorMessage.value = error
    return
  }

  if (data && data.token) {
    token.value = data.token
    name.value = data.nama
    branch.value = data.cabang
    fkUser.value = data.fk_user
    await navigateTo('/dashboard')
  }
}
</script>

<style scoped>
.login-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #fef3c7 0%, #ffffff 50%, #fef3c7 100%);
  position: relative;
  overflow: hidden;
  padding: 1.5rem;
}

.login-decor-1 {
  position: absolute;
  top: -120px;
  right: -120px;
  width: 350px;
  height: 350px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(245, 158, 11, 0.2) 0%, transparent 70%);
  filter: blur(40px);
}

.login-decor-2 {
  position: absolute;
  bottom: -120px;
  left: -120px;
  width: 350px;
  height: 350px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(217, 119, 6, 0.15) 0%, transparent 70%);
  filter: blur(40px);
}

.login-card {
  width: 100%;
  max-width: 420px;
  background: #ffffff;
  border-radius: 20px;
  padding: 2.5rem;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.08);
  z-index: 10;
  border: 1px solid rgba(245, 158, 11, 0.1);
}

.login-header {
  text-align: center;
  margin-bottom: 2rem;
}

.login-logo {
  display: flex;
  justify-content: center;
  margin-bottom: 0.75rem;
}

.login-header h1 {
  font-size: 2rem;
  font-weight: 800;
  color: var(--selada-500);
}

.login-header h1 span {
  color: var(--selada-400);
  font-weight: 600;
}

.error-banner {
  background: #fef2f2;
  border: 1px solid #fecaca;
  color: #dc2626;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  margin-bottom: 1.5rem;
  font-size: 0.875rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.login-form {
  display: flex;
  flex-direction: column;
}

.password-container {
  position: relative;
  width: 100%;
}

.password-field {
  padding-right: 3rem !important;
}

.password-toggle {
  position: absolute;
  right: 0.75rem;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: var(--text-light);
  cursor: pointer;
  padding: 0.25rem;
  display: flex;
  align-items: center;
  transition: var(--transition-smooth);
}

.password-toggle:hover {
  color: var(--selada-400);
}

.has-error {
  border-color: var(--color-danger) !important;
}

.has-error:focus {
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.15) !important;
}

.error-text {
  display: block;
  font-size: 0.75rem;
  color: #dc2626;
  margin-top: 0.35rem;
}

.login-btn {
  width: 100%;
  margin-top: 0.5rem;
  padding: 0.875rem;
  font-size: 1rem;
}
</style>
