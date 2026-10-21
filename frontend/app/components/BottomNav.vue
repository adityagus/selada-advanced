<template>
  <nav class="bottom-nav">
    <div class="bottom-nav-inner">
      <NuxtLink to="/dashboard" class="bottom-nav-item" :class="{ active: currentRoute === '/dashboard' }">
        <AppIcon name="dashboard" width="20" height="20" />
        <span>Dashboard</span>
      </NuxtLink>

      <NuxtLink to="/inquiry" class="bottom-nav-item" :class="{ active: currentRoute === '/inquiry' }">
        <AppIcon name="inquiry" width="20" height="20" />
        <span>Inquiry</span>
      </NuxtLink>

      <NuxtLink to="/data-customer" class="bottom-nav-item" :class="{ active: currentRoute === '/data-customer' }">
        <AppIcon name="users" width="20" height="20" />
        <span>Customer</span>
      </NuxtLink>

      <NuxtLink to="/cek-taksiran" class="bottom-nav-item" :class="{ active: currentRoute === '/cek-taksiran' }">
        <AppIcon name="calculator" width="20" height="20" />
        <span>CeTar</span>
      </NuxtLink>

      <NuxtLink to="/report" class="bottom-nav-item" :class="{ active: currentRoute === '/report' }">
        <AppIcon name="chart" width="20" height="20" />
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

