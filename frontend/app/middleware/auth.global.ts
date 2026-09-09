import { isTokenExpired, clearAuthCookies } from '~/utils/auth'

export default defineNuxtRouteMiddleware((to) => {
  const publicPages = ['/login']
  const isPublicPage = publicPages.includes(to.path)

  const token = useCookie('auth_token').value
  const tokenExpiry = useCookie('token_expiry').value

  const expired = isTokenExpired(token, tokenExpiry)

  // 1. If token is missing or expired and user accesses protected route -> clear cookies & redirect to /login
  if ((!token || expired) && !isPublicPage) {
    clearAuthCookies()
    return navigateTo('/login')
  }

  // 2. If token is valid and user accesses /login -> redirect to /dashboard
  if (token && !expired && isPublicPage) {
    return navigateTo('/dashboard')
  }
})

