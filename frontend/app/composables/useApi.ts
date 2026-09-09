import { clearAuthCookies, isTokenExpired } from '~/utils/auth'

interface FetchOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'
  body?: any
  params?: Record<string, any>
  headers?: Record<string, string>
}

// Client-side helper to read cookies when Nuxt context is unavailable
const getCookieClient = (name: string): string | null => {
  if (typeof document === 'undefined') return null
  const nameEQ = name + '='
  const ca = document.cookie.split(';')
  for (let i = 0; i < ca.length; i++) {
    let c : any  = ca[i]
    while (c.charAt(0) === ' ') c = c.substring(1, c.length)
    if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length)
  }
  return null
}

export const useApi = async <T = any>(endpoint: string, options: FetchOptions = {}) => {
  let apiBase = ''
  
  // Try getting runtime config safely
  try {
    const config = useRuntimeConfig()
    if (config?.public?.apiBase) {
      apiBase = config.public.apiBase
    }
  } catch (e) {
    // Context lost - use default
  }

  if (!apiBase || apiBase === 'http://localhost:3000/api') {
    if (typeof window !== 'undefined' && window.location && window.location.hostname) {
      const protocol = window.location.protocol || 'http:'
      const hostname = window.location.hostname
      apiBase = `${protocol}//${hostname}:3000/api`
    } else {
      apiBase = 'http://localhost:3000/api'
    }
  }

  const cleanEndpoint = endpoint.startsWith('/') ? endpoint : '/' + endpoint
  const url = `${apiBase}${cleanEndpoint}`
  const isLoginEndpoint = cleanEndpoint === '/login'

  // Retrieve token safely
  let tokenVal: string | null = null
  let tokenExpiryVal: string | null = null
  try {
    const token = useCookie('auth_token')
    const tokenExpiry = useCookie('token_expiry')
    tokenVal = token.value || null
    tokenExpiryVal = tokenExpiry.value || null
  } catch (e) {
    tokenVal = getCookieClient('auth_token')
    tokenExpiryVal = getCookieClient('token_expiry')
  }

  // Pre-flight check: if not a login endpoint and token is missing or expired, auto-redirect to login
  if (!isLoginEndpoint && isTokenExpired(tokenVal, tokenExpiryVal)) {
    clearAuthCookies()
    try {
      navigateTo('/login')
    } catch (e) {
      if (typeof window !== 'undefined') {
        window.location.href = '/login'
      }
    }
    return {
      data: null,
      error: 'Token tidak ada atau sudah kedaluwarsa. Silakan login kembali.'
    }
  }

  // Standardize request headers
  const headers: Record<string, string> = {
    ...options.headers,
  }

  if (tokenVal) {
    headers['Authorization'] = `Bearer ${tokenVal}`
  }

  try {
    const data = await $fetch<T>(url, {
      ...options,
      headers,
    })
    return { data, error: null }
  } catch (err: any) {
    const status = err.status || err.statusCode || err.response?.status
    // Auto-logout in case of 401 Unauthorized or 403 Forbidden responses
    if (status === 401 || status === 403) {
      clearAuthCookies()

      // Redirect safely
      try {
        navigateTo('/login')
      } catch (e) {
        if (typeof window !== 'undefined') {
          window.location.href = '/login'
        }
      }
    }
    
    return {
      data: null,
      error: err.data?.error || err.data?.message || err.message || 'Terjadi kesalahan sistem'
    }
  }
}

