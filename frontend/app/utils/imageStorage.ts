/**
 * Utility for compressing, caching, and retrieving collateral/activity photos.
 * Keeps image sizes compact (typically ~100-250KB) so they fit safely within
 * browser client storage (localStorage) without quota overflow.
 */

export function compressImage(
  file: File,
  maxWidth = 1280,
  maxHeight = 1280,
  quality = 0.82
): Promise<string> {
  return new Promise((resolve, reject) => {
    if (!file || !file.type.startsWith('image/')) {
      return reject(new Error('File yang dipilih bukan gambar'))
    }

    const reader = new FileReader()
    reader.onerror = () => reject(new Error('Gagal membaca file gambar'))
    reader.onload = (e) => {
      const dataUrl = e.target?.result as string
      if (!dataUrl) {
        return reject(new Error('Data gambar kosong'))
      }

      // If in SSR environment where Image is not available, return as-is
      if (typeof window === 'undefined' || typeof Image === 'undefined') {
        return resolve(dataUrl)
      }

      const img = new Image()
      img.onerror = () => reject(new Error('Gagal memuat gambar untuk kompresi'))
      img.onload = () => {
        let width = img.width
        let height = img.height

        if (width > maxWidth || height > maxHeight) {
          const ratio = Math.min(maxWidth / width, maxHeight / height)
          width = Math.round(width * ratio)
          height = Math.round(height * ratio)
        }

        const canvas = document.createElement('canvas')
        canvas.width = width
        canvas.height = height

        const ctx = canvas.getContext('2d')
        if (!ctx) {
          return resolve(dataUrl)
        }

        // Fill white background for transparent PNGs converted to JPEG
        ctx.fillStyle = '#FFFFFF'
        ctx.fillRect(0, 0, width, height)
        ctx.drawImage(img, 0, 0, width, height)

        try {
          const compressed = canvas.toDataURL('image/jpeg', quality)
          resolve(compressed)
        } catch (err) {
          resolve(dataUrl)
        }
      }
      img.src = dataUrl
    }

    reader.readAsDataURL(file)
  })
}

export function getFotoBarang(id: string | number | undefined | null): string | null {
  if (!id || typeof window === 'undefined' || typeof localStorage === 'undefined') {
    return null
  }
  try {
    return (
      localStorage.getItem(`foto_barang_${id}`) ||
      localStorage.getItem(`foto_barang_${String(id)}`) ||
      null
    )
  } catch {
    return null
  }
}

export function setFotoBarang(
  id: string | number | undefined | null,
  dataUrl: string
): boolean {
  if (!id || !dataUrl || typeof window === 'undefined' || typeof localStorage === 'undefined') {
    return false
  }
  const key = `foto_barang_${id}`
  try {
    localStorage.setItem(key, dataUrl)
    return true
  } catch (err: any) {
    // Quota exceeded: prune older photos except the current one and retry
    try {
      const keysToRemove: string[] = []
      for (let i = 0; i < localStorage.length; i++) {
        const k = localStorage.key(i)
        if (k && k.startsWith('foto_barang_') && k !== key) {
          keysToRemove.push(k)
        }
      }
      // Remove older items if storage is full (keep up to 5 newest)
      if (keysToRemove.length > 5) {
        keysToRemove.slice(0, keysToRemove.length - 5).forEach((k) => localStorage.removeItem(k))
      }
      localStorage.setItem(key, dataUrl)
      return true
    } catch {
      console.warn('Gagal menyimpan foto ke localStorage (quota exceeded)')
      return false
    }
  }
}

export function removeFotoBarang(id: string | number | undefined | null): void {
  if (!id || typeof window === 'undefined' || typeof localStorage === 'undefined') {
    return
  }
  try {
    localStorage.removeItem(`foto_barang_${id}`)
  } catch {
    // Ignore error
  }
}
