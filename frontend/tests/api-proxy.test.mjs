import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { createServer } from 'node:http'
import { test } from 'node:test'
import { runInNewContext } from 'node:vm'
import { createApp, defineEventHandler, getRequestURL, proxyRequest, toNodeListener } from 'h3'

const source = readFileSync(new URL('../server/api/[...path].ts', import.meta.url), 'utf8')
  .replace(/^import .*$/gm, '')
  .replace('export default ', 'const handler = ')

async function listen(server) {
  await new Promise(resolve => server.listen(0, '127.0.0.1', resolve))
  return `http://127.0.0.1:${server.address().port}`
}

async function close(server) {
  await new Promise((resolve, reject) => {
    server.close(error => error ? reject(error) : resolve())
    server.closeAllConnections()
  })
}

test('same-origin API forwards login JSON, query strings, auth headers and backend status', async (t) => {
  const received = []
  const backend = createServer(async (request, response) => {
    let body = ''
    for await (const chunk of request) body += chunk
    received.push({ url: request.url, method: request.method, body, authorization: request.headers.authorization })
    response.setHeader('Content-Type', 'application/json')
    if (request.url === '/api/login' && body === '{}') {
      response.writeHead(400).end(JSON.stringify({ error: 'Username dan password wajib diisi' }))
    } else if (request.url === '/api/login') {
      response.end(JSON.stringify({ token: 'test-token' }))
    } else {
      response.writeHead(401).end(JSON.stringify({ error: 'Unauthorized' }))
    }
  })
  const backendURL = await listen(backend)
  t.after(() => close(backend))

  const handler = runInNewContext(`${source}\nhandler`, {
    defineEventHandler, getRequestURL, proxyRequest,
    useRuntimeConfig: () => ({ apiProxyTarget: `${backendURL}/` }),
  })
  const frontend = createServer(toNodeListener(createApp().use(handler)))
  const frontendURL = await listen(frontend)
  t.after(() => close(frontend))

  const loginBody = JSON.stringify({ user: 'test-user', pass: 'test-password' })
  const login = await fetch(`${frontendURL}/api/login`, {
    method: 'POST', headers: { 'Content-Type': 'application/json' }, body: loginBody,
  })
  assert.equal(login.status, 200)
  assert.match(login.headers.get('content-type'), /application\/json/)
  assert.deepEqual(await login.json(), { token: 'test-token' })
  assert.deepEqual(received[0], { url: '/api/login', method: 'POST', body: loginBody, authorization: undefined })

  const invalid = await fetch(`${frontendURL}/api/login`, {
    method: 'POST', headers: { 'Content-Type': 'application/json' }, body: '{}',
  })
  assert.equal(invalid.status, 400)
  assert.equal((await invalid.json()).error, 'Username dan password wajib diisi')

  const protectedResponse = await fetch(`${frontendURL}/api/dashboard?branch=3&search=a%20b`, {
    headers: { Authorization: 'Bearer test-token' },
  })
  assert.equal(protectedResponse.status, 401)
  assert.equal((await protectedResponse.json()).error, 'Unauthorized')
  assert.equal(received[2].url, '/api/dashboard?branch=3&search=a%20b')
  assert.equal(received[2].authorization, 'Bearer test-token')
})
