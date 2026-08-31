// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },

  devServer: {
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
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:3000/api'
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
      link: [
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: 'anonymous' },
        { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=Poppins:wght@300;400;500;600;700;800;900&display=swap' }
      ]
    }
  },

  modules: ['@nuxt/image', '@nuxtjs/tailwindcss']
})