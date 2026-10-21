import { defineEventHandler, getRequestURL, proxyRequest } from 'h3'
import { useRuntimeConfig } from '#imports'

export default defineEventHandler((event) => {
  const config = useRuntimeConfig(event)
  const target = config.apiProxyTarget.replace(/\/+$/, '')
  const { pathname, search } = getRequestURL(event)

  return proxyRequest(event, `${target}${pathname}${search}`)
})
