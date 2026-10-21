import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { stripTypeScriptTypes } from 'node:module'
import { test } from 'node:test'
import { runInNewContext } from 'node:vm'

const authSource = stripTypeScriptTypes(
  readFileSync(new URL('../app/utils/auth.ts', import.meta.url), 'utf8')
).replace(/export /g, '')

function createContext(initialCookies = {}, initialStorage = {}) {
  const cookieStore = new Map(Object.entries(initialCookies))
  const storageStore = new Map(Object.entries(initialStorage))

  const documentObj = {
    get cookie() {
      return Array.from(cookieStore.entries())
        .map(([k, v]) => `${k}=${encodeURIComponent(v)}`)
        .join('; ')
    },
    set cookie(val) {
      const parts = val.split(';')
      const [nameVal] = parts
      const eqIdx = nameVal.indexOf('=')
      const name = nameVal.slice(0, eqIdx).trim()
      const value = decodeURIComponent(nameVal.slice(eqIdx + 1).trim())

      // Check max-age or expires for removal
      const isExpired = parts.some((p) => {
        const lower = p.trim().toLowerCase()
        return lower.startsWith('max-age=0') || lower.includes('1970')
      })

      if (isExpired || !value) {
        cookieStore.delete(name)
      } else {
        cookieStore.set(name, value)
      }
    }
  }

  const localStorageObj = {
    getItem(k) {
      return storageStore.has(k) ? storageStore.get(k) : null
    },
    setItem(k, v) {
      storageStore.set(k, String(v))
    },
    removeItem(k) {
      storageStore.delete(k)
    }
  }

  const context = {
    Buffer,
    document: documentObj,
    localStorage: localStorageObj,
    window: {
      location: { protocol: 'http:' }
    },
    useCookie(name) {
      return {
        get value() {
          return cookieStore.get(name) || null
        },
        set value(v) {
          if (v === null || v === undefined) {
            cookieStore.delete(name)
          } else {
            cookieStore.set(name, String(v))
          }
        }
      }
    }
  }

  const sandbox = runInNewContext(
    `${authSource}\n({ SESSION_IDLE_TIMEOUT_MS, SESSION_MAX_AGE_SEC, isTokenExpired, refreshSessionActivity, clearAuthCookies })`,
    context
  )

  return { sandbox, cookieStore, storageStore }
}

test('session timeout constants are set to 30 minutes', () => {
  const { sandbox } = createContext()
  assert.equal(sandbox.SESSION_IDLE_TIMEOUT_MS, 30 * 60 * 1000) // 1,800,000 ms
  assert.equal(sandbox.SESSION_MAX_AGE_SEC, 30 * 60) // 1,800 seconds
})

test('isTokenExpired correctly marks token expired after 30 minutes without activity', () => {
  const { sandbox } = createContext()
  const fakeToken = 'header.payload.signature'

  // 31 minutes ago (expired)
  const expiredTimestamp = Date.now() - 31 * 60 * 1000
  assert.equal(sandbox.isTokenExpired(fakeToken, expiredTimestamp), true)

  // 10 minutes from now (active)
  const activeTimestamp = Date.now() + 10 * 60 * 1000
  assert.equal(sandbox.isTokenExpired(fakeToken, activeTimestamp), false)
})

test('refreshSessionActivity extends token_expiry by 30 minutes and updates localStorage', () => {
  const initialExpiry = Date.now() + 5 * 60 * 1000 // 5 minutes left
  const { sandbox, cookieStore, storageStore } = createContext({
    auth_token: 'header.payload.signature',
    token_expiry: String(initialExpiry),
    user_name: 'Budi'
  })

  sandbox.refreshSessionActivity()

  const updatedExpiry = Number(cookieStore.get('token_expiry'))
  const expectedMinExpiry = Date.now() + 29 * 60 * 1000 // close to 30 mins
  assert.ok(updatedExpiry >= expectedMinExpiry, 'token_expiry should be extended by ~30 minutes')

  assert.ok(storageStore.has('selada_last_activity'), 'localStorage should record last activity')
})

test('clearAuthCookies wipes all cookies and clears activity storage', () => {
  const { sandbox, cookieStore, storageStore } = createContext(
    {
      auth_token: 'valid-token',
      token_expiry: String(Date.now() + 1800000),
      user_name: 'Budi',
      user_branch: '001',
      user_branch_name: 'Cabang Pusat',
      fk_user: '123'
    },
    { selada_last_activity: String(Date.now()) }
  )

  assert.equal(cookieStore.has('auth_token'), true)
  assert.equal(storageStore.has('selada_last_activity'), true)

  sandbox.clearAuthCookies()

  assert.equal(cookieStore.has('auth_token'), false)
  assert.equal(cookieStore.has('token_expiry'), false)
  assert.equal(cookieStore.has('user_name'), false)
  assert.equal(storageStore.has('selada_last_activity'), false)
})
