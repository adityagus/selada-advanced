// // app/router.options.ts
// import type { RouterConfig } from '@nuxt/schema'

// export default <RouterConfig> {
//   // Mengoverride halaman default Nuxt
//   routes: (_routes) => [
//     {
//       name: 'home',
//       path: '/',
//       component: () => import('~/pages/dashboard.vue') // Mengarah ke komponen bebas
//     },
//     {
//       name: 'tentang',
//       path: '/tentang-kami',
//       component: () => import('~/pages/cek-taksiran.vue')
//     },
//     {
//       name: 'produk-detail',
//       path: '/produk/:id',
//       component: () => import('~/pages/cek-taksiran.vue')
//     }
//   ]
// }
