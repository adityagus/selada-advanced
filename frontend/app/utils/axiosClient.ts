import axios from 'axios'
import { clearAuthCookies, isTokenExpired } from '~/utils/auth'

const getCookieClient = (name: string): string | null => {
  if (typeof document === 'undefined') return null
  const nameEQ = name + '='
  const ca = document.cookie.split(';')
  for (let i = 0; i < ca.length; i++) {
    let c: any = ca[i]
    while (c.charAt(0) === ' ') c = c.substring(1, c.length)
    if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length)
  }
  return null
}

const axiosClient = axios.create({
  baseURL: 'http://localhost:3000/api',
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  }
})

// Request Interceptor: Attach JWT Bearer Token
axiosClient.interceptors.request.use((config) => {
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

  const isLoginEndpoint = config.url?.includes('/login')

  if (!isLoginEndpoint && isTokenExpired(tokenVal, tokenExpiryVal)) {
    clearAuthCookies()
    if (typeof window !== 'undefined') {
      window.location.href = '/login'
    }
    return Promise.reject(new Error('Token kedaluwarsa atau tidak terdeteksi. Silakan login kembali.'))
  }

  if (tokenVal) {
    config.headers.Authorization = `Bearer ${tokenVal}`
  }

  return config
}, (error) => {
  return Promise.reject(error)
})

// Response Interceptor: Handle 401 / 403 Unauthorized Auto-Logout
axiosClient.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error.response?.status
    if (status === 401 || status === 403) {
      clearAuthCookies()
      if (typeof window !== 'undefined') {
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  }
)

export default axiosClient
