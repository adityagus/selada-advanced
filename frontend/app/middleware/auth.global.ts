import { isTokenExpired, clearAuthCookies, refreshSessionActivity } from '~/utils/auth'

export default defineNuxtRouteMiddleware((to) => {
  const publicPages = ['/login']
  const isPublicPage = publicPages.includes(to.path)

  const token = useCookie('auth_token').value
  const tokenExpiry = useCookie('token_expiry').value

  const expired = isTokenExpired(token, tokenExpiry)

  // 1. If token is missing or expired and user accesses protected route -> clear cookies & redirect to /login
  if ((!token || expired) && !isPublicPage) {
    clearAuthCookies()
    const redirectUrl = (token && expired) ? '/login?reason=timeout' : '/login'
    return navigateTo(redirectUrl)
  }

  // 2. If token is valid and user accesses /login -> redirect to /dashboard
  if (token && !expired && isPublicPage) {
    return navigateTo('/dashboard')
  }

  // 3. Active navigation on protected routes touches/refreshes the session
  if (token && !expired && !isPublicPage) {
    refreshSessionActivity()
  }
})

