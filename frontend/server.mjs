// Production'da Railway/Nixpacks tarafından çalıştırılır.
// Vite ile derlenmiş `dist/` klasörünü servis eder ve SPA fallback uygular
// (her bilinmeyen path için index.html döner ki Vue Router client-side çalışsın).
import http from 'node:http'
import path from 'node:path'
import fs from 'node:fs/promises'
import { fileURLToPath } from 'node:url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const distDir = path.join(__dirname, 'dist')
const port = parseInt(process.env.PORT || '4173', 10)

const mime = {
  '.html': 'text/html; charset=utf-8',
  '.js':   'application/javascript; charset=utf-8',
  '.mjs':  'application/javascript; charset=utf-8',
  '.css':  'text/css; charset=utf-8',
  '.json': 'application/json; charset=utf-8',
  '.svg':  'image/svg+xml',
  '.png':  'image/png',
  '.jpg':  'image/jpeg',
  '.jpeg': 'image/jpeg',
  '.gif':  'image/gif',
  '.ico':  'image/x-icon',
  '.webp': 'image/webp',
  '.woff':  'font/woff',
  '.woff2': 'font/woff2',
  '.ttf':   'font/ttf',
  '.otf':   'font/otf',
  '.map':   'application/json; charset=utf-8',
  '.txt':   'text/plain; charset=utf-8',
}

async function tryFile(p) {
  try {
    const stat = await fs.stat(p)
    if (stat.isFile()) return p
  } catch { /* */ }
  return null
}

async function send(res, file, status = 200) {
  const ext = path.extname(file).toLowerCase()
  const data = await fs.readFile(file)
  const headers = {
    'Content-Type': mime[ext] || 'application/octet-stream',
    'Content-Length': data.length,
  }
  if (ext === '.html') {
    headers['Cache-Control'] = 'no-cache, no-store, must-revalidate'
  } else if (file.includes('/assets/')) {
    headers['Cache-Control'] = 'public, max-age=31536000, immutable'
  } else {
    headers['Cache-Control'] = 'public, max-age=600'
  }
  res.writeHead(status, headers)
  res.end(data)
}

const server = http.createServer(async (req, res) => {
  try {
    const url = new URL(req.url, `http://localhost:${port}`)
    let pathname = decodeURIComponent(url.pathname)
    if (pathname.includes('..')) {
      res.writeHead(400); return res.end('Bad Request')
    }
    if (pathname === '/' || pathname === '') pathname = '/index.html'

    const filePath = path.join(distDir, pathname)
    const file = await tryFile(filePath)
    if (file) return await send(res, file)

    // SPA fallback
    const index = path.join(distDir, 'index.html')
    return await send(res, index)
  } catch (err) {
    console.error('serve error:', err)
    res.writeHead(500)
    res.end('Internal Server Error')
  }
})

server.listen(port, () => {
  console.log(`kare-rehber frontend listening on :${port} (dist=${distDir})`)
})
