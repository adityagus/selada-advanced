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
  token: string | null | undefined,
  tokenExpiryCookie?: string | number | null
): boolean => {
  if (!token || typeof token !== 'string' || !token.trim()) {
    return true
  }

  // 1. Check token_expiry cookie if provided
  if (tokenExpiryCookie) {
    const expiryNum = Number(tokenExpiryCookie)
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
 * Clears all auth-related cookies safely.
 */
export const clearAuthCookies = () => {
  try {
    const token = useCookie('auth_token')
    const tokenExpiry = useCookie('token_expiry')
    const userName = useCookie('user_name')
    const userBranch = useCookie('user_branch')
    const fkUser = useCookie('fk_user')

    token.value = null
    tokenExpiry.value = null
    userName.value = null
    userBranch.value = null
    fkUser.value = null
  } catch (e) {
    // Fallback if Nuxt context is unavailable (client side)
    if (typeof document !== 'undefined') {
      const cookies = ['auth_token', 'token_expiry', 'user_name', 'user_branch', 'fk_user']
      cookies.forEach((c) => {
        document.cookie = `${c}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`
      })
    }
  }
}
