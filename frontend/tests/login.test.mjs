import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { stripTypeScriptTypes } from 'node:module'
import { test } from 'node:test'
import { runInNewContext } from 'node:vm'
import { ref, watch, nextTick } from 'vue'
import { z } from 'zod'

const page = readFileSync(new URL('../app/pages/login.vue', import.meta.url), 'utf8')
const setup = page.match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1]
  .replace(/^import .*$/gm, '')
const auth = stripTypeScriptTypes(readFileSync(new URL('../app/utils/auth.ts', import.meta.url), 'utf8'))
  .replace(/export /g, '')
const schema = stripTypeScriptTypes(readFileSync(new URL('../app/schemas/authSchema.ts', import.meta.url), 'utf8'))
  .replace(/^import .*$/gm, '').replace(/export /g, '')

async function login(response) {
  const cookies = new Map()
  const stops = []
  let destination
  const context = {
    ref, nextTick, z, Buffer, definePageMeta() {},
    document: {
      set cookie(value) {
        const pair = value.split(';')[0]
        const separator = pair.indexOf('=')
        cookies.set(pair.slice(0, separator), decodeURIComponent(pair.slice(separator + 1)))
      },
    },
    window: {
      location: {
        protocol: 'http:',
        replace(path) {
          assert.ok(cookies.get('auth_token'), 'session must be persisted before navigation')
          assert.ok(Number(cookies.get('token_expiry')) > Date.now())
          destination = path
        },
      },
    },
    useCookie(name) {
      const value = ref(null)
      stops.push(watch(value, (v) => cookies.set(name, v)))
      return value
    },
    async useApi() {
      if (response instanceof Error) throw response
      return response
    },
    async navigateTo(path) {
      assert.ok(cookies.get('auth_token'), 'session cookie must be flushed before navigation')
      assert.ok(Number(cookies.get('token_expiry')) > Date.now())
      destination = path
    },
  }
  const state = runInNewContext(`${auth}\n${schema}\n${stripTypeScriptTypes(setup)}\n({ username, password, handleLogin, errorMessage, isLoading })`, context)
  state.username.value = 'test-user'
  state.password.value = 'testpass'
  try {
    await state.handleLogin()
    return { destination, error: state.errorMessage.value, loading: state.isLoading.value }
  } finally {
    stops.forEach(stop => stop())
  }
}

test('successful login flushes session cookies and navigates to dashboard', async () => {
  const payload = Buffer.from(JSON.stringify({ expiry: new Date(Date.now() + 60000).toISOString() })).toString('base64url')
  const result = await login({ data: { token: `${payload}.signature`, nama: 'Test' }, error: null })
  assert.equal(result.destination, '/dashboard')
  assert.equal(result.error, '')
  assert.equal(result.loading, false)
})

test('HTTP 200 without a token displays an error', async () => {
  const result = await login({ data: { message: 'success' }, error: null })
  assert.equal(result.destination, undefined)
  assert.match(result.error, /Respons login/)
  assert.equal(result.loading, false)
})

test('HTTP 200 containing the frontend HTML is not treated as a successful login', async () => {
  const result = await login({ data: '<!DOCTYPE html><html></html>', error: null })
  assert.equal(result.destination, undefined)
  assert.match(result.error, /Respons login/)
  assert.equal(result.loading, false)
})

test('unexpected request failure restores the submit button', async () => {
  const result = await login(new Error('Network failure'))
  assert.equal(result.destination, undefined)
  assert.match(result.error, /Login gagal/)
  assert.equal(result.loading, false)
})
