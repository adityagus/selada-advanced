// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: false },

  devServer: {
    host: '0.0.0.0',
    port: 3001
  },

  // Set to SPA mode (Single Page Application) for client-only rendering
  ssr: false,

  // Global CSS rules
  css: [
    '~/assets/css/main.css'
  ],

  // Public runtime configuration variables (API endpoints)
  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE
    }
  },

  // HTML Head Configuration
  app: {
    head: {
      title: 'Selada V2 - Sistem Rencana & Taksiran Gadai',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'Aplikasi Manajemen Target, Katalog Motor, dan Estimasi Taksiran Gadai Selada V2.' }
      ],
      link: []
    }
  },

  vite: {
    build: {
      cssMinify: true,
      minify: 'esbuild',
      chunkSizeWarningLimit: 1000
    }
  },

  nitro: {
    compressPublicAssets: true
  },

  modules: ['@nuxt/image', '@nuxtjs/tailwindcss']
})