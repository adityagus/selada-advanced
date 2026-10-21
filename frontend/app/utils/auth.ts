/**
 * Idle timeout duration: 30 minutes in milliseconds
 */
export const SESSION_IDLE_TIMEOUT_MS = 30 * 60 * 1000 // 1,800,000 ms

/**
 * Cookie Max-Age: 30 minutes in seconds
 */
export const SESSION_MAX_AGE_SEC = 30 * 60 // 1,800 s

/**
 * Safely parses JWT base64url payload on both Node (SSR) and Browser environments.
 */
export const parseJwt = (token: string): Record<string, any> | null => {
  try {
    const parts = token.split('.')
    let base64Url = ''
    if (parts.length === 3 && parts[1]) {
      base64Url = parts[1]
    } else if (parts.length === 2 && parts[0]) {
      base64Url = parts[0]
    } else {
      return null
    }

    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')

    let jsonPayload = ''
    const globalBuffer = typeof globalThis !== 'undefined' ? (globalThis as any).Buffer : null
    if (globalBuffer) {
      jsonPayload = globalBuffer.from(base64, 'base64').toString('utf-8')
    } else if (typeof atob !== 'undefined') {
      jsonPayload = decodeURIComponent(
        atob(base64)
          .split('')
          .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
          .join('')
      )
    } else {
      return null
    }

    return JSON.parse(jsonPayload)
  } catch (e) {
    return null
  }
}

/**
 * Checks if a token is missing, expired by cookie date, or expired by JWT exp/expiry claim.
 */
export const isTokenExpired = (
  token?: string | null,
  tokenExpiryCookie?: string | number | null
): boolean => {
  if (!token || typeof token !== 'string' || !token.trim()) {
    return true
  }

  let expiryToCheck = tokenExpiryCookie
  if (!expiryToCheck && typeof document !== 'undefined' && typeof document.cookie === 'string') {
    const match = document.cookie.match(/(?:^|;\s*)token_expiry=([^;]*)/)
    if (match && match[1]) {
      expiryToCheck = decodeURIComponent(match[1])
    }
  }

  // 1. Check token_expiry cookie if provided or found
  if (expiryToCheck) {
    const expiryNum = Number(expiryToCheck)
    if (!isNaN(expiryNum) && expiryNum > 0) {
      // Convert to milliseconds if given in seconds
      const expiryMs = expiryNum < 1e11 ? expiryNum * 1000 : expiryNum
      if (Date.now() >= expiryMs) {
        return true
      }
    }
  }

  // 2. Check exp / expiry claim in JWT payload
  const payload = parseJwt(token)
  if (payload) {
    if (typeof payload.exp === 'number') {
      const expMs = payload.exp * 1000
      if (Date.now() >= expMs) {
        return true
      }
    }
    if (payload.expiry) {
      const expTime = new Date(payload.expiry).getTime()
      if (!isNaN(expTime) && Date.now() >= expTime) {
        return true
      }
    }
  }

  return false
}

/**
 * Refreshes session activity timestamp and extends session cookies by 30 minutes.
 * Triggered by user interaction (mouse, touch, keypress, scroll) or API calls.
 */
export const refreshSessionActivity = () => {
  if (typeof document === 'undefined' || typeof document.cookie !== 'string') return

  const tokenMatch = document.cookie.match(/(?:^|;\s*)auth_token=([^;]*)/)
  if (!tokenMatch || !tokenMatch[1]) return

  const tokenVal = decodeURIComponent(tokenMatch[1])
  const expiryMatch = document.cookie.match(/(?:^|;\s*)token_expiry=([^;]*)/)
  const tokenExpiryVal = expiryMatch ? decodeURIComponent(expiryMatch[1]) : null

  if (isTokenExpired(tokenVal, tokenExpiryVal)) {
    return
  }

  const now = Date.now()
  const newExpiry = now + SESSION_IDLE_TIMEOUT_MS
  const secure = typeof window !== 'undefined' && window.location?.protocol === 'https:' ? '; Secure' : ''
  const maxAge = SESSION_MAX_AGE_SEC

  // Refresh token_expiry cookie with sliding 30-minute timestamp
  document.cookie = `token_expiry=${newExpiry}; Path=/; Max-Age=${maxAge}; SameSite=Lax${secure}`

  // Extend active session cookies to maintain 30-minute browser expiration
  const sessionCookieNames = ['auth_token', 'user_name', 'user_branch', 'user_branch_name', 'fk_user']
  for (const name of sessionCookieNames) {
    const match = document.cookie.match(new RegExp(`(?:^|;\\s*)${name}=([^;]*)`))
    if (match && match[1] !== undefined) {
      document.cookie = `${name}=${match[1]}; Path=/; Max-Age=${maxAge}; SameSite=Lax${secure}`
    }
  }

  if (typeof localStorage !== 'undefined') {
    try {
      localStorage.setItem('selada_last_activity', String(now))
    } catch (e) {
      // Ignore localStorage errors
    }
  }

  try {
    const expiryCookie = useCookie('token_expiry')
    if (expiryCookie) {
      expiryCookie.value = String(newExpiry)
    }
  } catch (e) {
    // Ignore outside Nuxt context
  }
}

/**
 * Clears all auth-related cookies safely.
 */
export const clearAuthCookies = () => {
  try {
    const token = useCookie('auth_token')
    const tokenExpiry = useCookie('token_expiry')
    const userName = useCookie('user_name')
    const userBranch = useCookie('user_branch')
    const userBranchName = useCookie('user_branch_name')
    const fkUser = useCookie('fk_user')

    token.value = null
    tokenExpiry.value = null
    userName.value = null
    userBranch.value = null
    userBranchName.value = null
    fkUser.value = null
  } catch (e) {
    // Fallback if Nuxt context is unavailable (client side)
  }

  if (typeof document !== 'undefined') {
    const cookies = ['auth_token', 'token_expiry', 'user_name', 'user_branch', 'user_branch_name', 'fk_user']
    cookies.forEach((c) => {
      document.cookie = `${c}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; Max-Age=0; path=/;`
    })
  }

  if (typeof localStorage !== 'undefined') {
    try {
      localStorage.removeItem('selada_last_activity')
    } catch (e) {
      // Ignore
    }
  }
}
