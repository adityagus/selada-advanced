<template>
  <nav class="bottom-nav">
    <div class="bottom-nav-inner">
      <NuxtLink to="/dashboard" class="bottom-nav-item" :class="{ active: currentRoute === '/dashboard' }">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M3 9L12 2L21 9V20C21 20.5304 20.7893 21.0391 20.4142 21.4142C20.0391 21.7893 19.5304 22 19 22H5C4.46957 22 3.96086 21.7893 3.58579 21.4142C3.21071 21.0391 3 20.5304 3 20V9Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
        <span>Dashboard</span>
      </NuxtLink>

      <NuxtLink to="/inquiry" class="bottom-nav-item" :class="{ active: currentRoute === '/inquiry' }">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <rect x="3" y="3" width="7" height="7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
          <rect x="14" y="3" width="7" height="7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
          <rect x="14" y="14" width="7" height="7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
          <rect x="3" y="14" width="7" height="7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
        <span>Inquiry</span>
      </NuxtLink>

      <NuxtLink to="/data-customer" class="bottom-nav-item" :class="{ active: currentRoute === '/data-customer' }">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M17 21V19C17 17.9391 16.5786 16.9217 15.8284 16.1716C15.0783 15.4214 14.0609 15 13 15H5C3.93913 15 2.92172 15.4214 2.17157 16.1716C1.42143 16.9217 1 17.9391 1 19V21" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
          <path d="M9 11C11.2091 11 13 9.20914 13 7C13 4.79086 11.2091 3 9 3C6.79086 3 5 4.79086 5 7C5 9.20914 6.79086 11 9 11Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
        <span>Customer</span>
      </NuxtLink>

      <NuxtLink to="/cek-taksiran" class="bottom-nav-item" :class="{ active: currentRoute === '/cek-taksiran' }">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor" width="20" height="20">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 3v17.25m0 0c-1.472 0-2.882.265-4.185.75M12 20.25c1.472 0 2.882.265 4.185.75M18.75 4.97A48.416 48.416 0 0 0 12 4.5c-2.291 0-4.545.16-6.75.47m13.5 0c1.01.143 2.01.317 3 .52m-3-.52 2.62 10.726c.122.499-.106 1.028-.589 1.202a5.988 5.988 0 0 1-2.031.352 5.988 5.988 0 0 1-2.031-.352c-.483-.174-.711-.703-.59-1.202L18.75 4.971Zm-16.5.52c.99-.203 1.99-.377 3-.52m0 0 2.62 10.726c.122.499-.106 1.028-.589 1.202a5.989 5.989 0 0 1-2.031.352 5.989 5.989 0 0 1-2.031-.352c-.483-.174-.711-.703-.59-1.202L5.25 4.971Z" />
        </svg>
        <span>CeTar</span>
      </NuxtLink>

      <NuxtLink to="/report" class="bottom-nav-item" :class="{ active: currentRoute === '/report' }">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M18 20V10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
          <path d="M12 20V4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
          <path d="M6 20V14" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
        <span>Report</span>
      </NuxtLink>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from '#app'

const route = useRoute()
const currentRoute = computed(() => route.path)
</script>

<style scoped>
.bottom-nav {
  display: none;
  position: fixed;
  bottom: 0.75rem;
  left: 0.75rem;
  right: 0.75rem;
  z-index: 50;
}

@media (max-width: 767px) {
  .bottom-nav {
    display: block;
  }
}

.bottom-nav-inner {
  display: flex;
  justify-content: space-around;
  align-items: center;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(226, 232, 240, 0.9);
  border-radius: 1.25rem;
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.08), 0 8px 10px -6px rgba(0, 0, 0, 0.03);
  padding: 0.35rem 0.5rem;
  height: 4.25rem;
  transition: all 0.3s ease;
}

:deep(.dark) .bottom-nav-inner {
  background: rgba(15, 23, 42, 0.95);
  border-color: rgba(30, 41, 59, 0.9);
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.5);
}

.bottom-nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  padding: 0.35rem;
  color: #64748b;
  transition: all 0.2s ease;
  border-radius: 0.85rem;
  gap: 0.2rem;
}

:deep(.dark) .bottom-nav-item {
  color: #94a3b8;
}

.bottom-nav-item span {
  font-size: 0.65rem;
  font-weight: 600;
}

.bottom-nav-item.active {
  background: #fef3c7;
  color: #d97706;
  font-weight: 700;
}

:deep(.dark) .bottom-nav-item.active {
  background: rgba(217, 119, 6, 0.25);
  color: #fbbf24;
}

.bottom-nav-item:not(.active):hover {
  color: #0f172a;
}

:deep(.dark) .bottom-nav-item:not(.active):hover {
  color: #f8fafc;
}
</style>

