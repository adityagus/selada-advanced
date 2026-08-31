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
    let c = ca[i]
    while (c.charAt(0) === ' ') c = c.substring(1, c.length)
    if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length)
  }
  return null
}

// Client-side helper to delete/write cookies
const setCookieClient = (name: string, value: string, days?: number) => {
  if (typeof document === 'undefined') return
  let expires = ""
  if (days) {
    const date = new Date()
    date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000))
    expires = "; expires=" + date.toUTCString()
  }
  document.cookie = name + "=" + (value || "")  + expires + "; path=/"
}

export const useApi = async <T = any>(endpoint: string, options: FetchOptions = {}) => {
  let apiBase = 'http://localhost:3000/api'
  
  // Try getting runtime config safely
  try {
    const config = useRuntimeConfig()
    if (config?.public?.apiBase) {
      apiBase = config.public.apiBase
    }
  } catch (e) {
    // Context lost - use default
  }

  const url = `${apiBase}${endpoint.startsWith('/') ? endpoint : '/' + endpoint}`

  // Standardize request headers
  const headers: Record<string, string> = {
    ...options.headers,
  }

  // Retrieve token safely
  let tokenVal: string | null = null
  try {
    const token = useCookie('auth_token')
    tokenVal = token.value || null
  } catch (e) {
    tokenVal = getCookieClient('auth_token')
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
    // Auto-logout in case of 401 Unauthorized responses
    if (err.status === 401) {
      // Clear token safely
      try {
        const token = useCookie('auth_token')
        token.value = null
      } catch (e) {
        setCookieClient('auth_token', '', -1)
        setCookieClient('user_name', '', -1)
        setCookieClient('user_branch', '', -1)
        setCookieClient('fk_user', '', -1)
      }

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
      error: err.data?.error || err.message || 'Terjadi kesalahan sistem'
    }
  }
}
