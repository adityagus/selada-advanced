import {
  clearAuthCookies,
  isTokenExpired,
  refreshSessionActivity,
  SESSION_IDLE_TIMEOUT_MS
} from '~/utils/auth'

/**
 * Client plugin that enforces a 30-minute idle/inactivity timeout.
 * - Any user interaction ("sentuhan atau kotak katik" like mouse, keyboard, touch, scroll)
 *   resets the 30-minute inactivity countdown (sliding expiration).
 * - If no interaction occurs for 30 minutes, cookies are cleared and the user is redirected to /login?reason=timeout.
 * - Synchronizes across multiple browser tabs via storage events.
 */
export default defineNuxtPlugin(() => {
  if (typeof window === 'undefined') return

  const token = useCookie('auth_token')
  const tokenExpiry = useCookie('token_expiry')

  let lastActivityTime = Date.now()
  let lastPersistTime = 0
  const PERSIST_THROTTLE_MS = 15_000 // Throttle cookie writes to at most once per 15s

  // Initialize last activity from storage if available
  if (typeof localStorage !== 'undefined') {
    try {
      const stored = localStorage.getItem('selada_last_activity')
      if (stored) {
        const parsed = Number(stored)
        if (!isNaN(parsed) && parsed > 0) {
          lastActivityTime = parsed
        }
      }
    } catch {}
  }

  const check = () => {
    if (window.location.pathname === '/login') return

    // If no token exists, immediately redirect to login
    if (!token.value) {
      clearAuthCookies()
      window.location.replace('/login')
      return
    }

    // Check localStorage in case another tab updated activity
    let latestActivity = lastActivityTime
    if (typeof localStorage !== 'undefined') {
      try {
        const stored = localStorage.getItem('selada_last_activity')
        if (stored) {
          const parsed = Number(stored)
          if (!isNaN(parsed) && parsed > latestActivity) {
            latestActivity = parsed
            lastActivityTime = parsed
          }
        }
      } catch {}
    }

    const isIdle = Date.now() - latestActivity >= SESSION_IDLE_TIMEOUT_MS
    const isExpired = isTokenExpired(token.value, tokenExpiry.value)

    if (isIdle || isExpired) {
      clearAuthCookies()
      if (typeof localStorage !== 'undefined') {
        try {
          localStorage.removeItem('selada_last_activity')
          localStorage.setItem('selada_logout', String(Date.now()))
        } catch {}
      }
      window.location.replace('/login?reason=timeout')
    }
  }

  // Activity handler triggered on user interaction ("sentuhan atau kotak-katik")
  const onUserActivity = () => {
    if (window.location.pathname === '/login') return
    if (!token.value) return

    const now = Date.now()
    lastActivityTime = now

    // Persist new 30-minute expiry to cookies/localStorage with throttling
    if (now - lastPersistTime >= PERSIST_THROTTLE_MS) {
      lastPersistTime = now
      refreshSessionActivity()
    }
  }

  // Interactive user events (desktop + mobile touch/scroll)
  const activityEvents = [
    'mousedown',
    'keydown',
    'touchstart',
    'touchmove',
    'scroll',
    'wheel',
    'mousemove',
    'click'
  ]

  activityEvents.forEach((eventName) => {
    window.addEventListener(eventName, onUserActivity, { passive: true, capture: true })
  })

  // Synchronize across multiple browser tabs
  window.addEventListener('storage', (event) => {
    if (event.key === 'selada_last_activity') {
      const parsed = Number(event.newValue)
      if (!isNaN(parsed) && parsed > 0) {
        lastActivityTime = parsed
      }
    } else if (event.key === 'selada_logout') {
      if (window.location.pathname !== '/login') {
        clearAuthCookies()
        window.location.replace('/login?reason=timeout')
      }
    }
  })

  // Check every 10 seconds for idle timeout expiration
  setInterval(check, 10_000)

  // Immediate check when tab becomes visible or focused
  document.addEventListener('visibilitychange', () => {
    if (document.visibilityState === 'visible') {
      check()
    }
  })
  window.addEventListener('focus', check)
})
