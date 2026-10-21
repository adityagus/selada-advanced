<template>
  <header class="mobile-header">
    <div class="mobile-header-inner">
      <div class="mobile-header-left">
        <span class="mobile-welcome">Selamat datang</span>
        <span class="mobile-username">{{ userName }}</span>
      </div>
      <button @click="$emit('logout')" class="mobile-logout-btn">
        <AppIcon name="logout" width="20" height="20" />
        <span>Logout</span>
      </button>
    </div>
  </header>
</template>

<script setup lang="ts">
defineProps<{
  userName: string
}>()

defineEmits(['logout'])
</script>

<style scoped>
.mobile-header {
  display: none;
  position: sticky;
  top: 0;
  z-index: 30;
  background: linear-gradient(180deg, #0b0f17 0%, #111827 100%);
  border-bottom: 1px solid rgba(245, 158, 11, 0.25);
  padding: 1rem 1.5rem;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.35);
}

@media (max-width: 767px) {
  .mobile-header {
    display: block;
  }
}

.mobile-header-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.mobile-header-left {
  display: flex;
  flex-direction: column;
}

.mobile-welcome {
  font-size: 0.75rem;
  color: rgba(229, 231, 235, 0.7);
}

.mobile-username {
  font-weight: 700;
  color: #fbbf24;
  font-size: 1rem;
}

.mobile-logout-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 0.75rem;
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 8px;
  color: #f87171;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  transition: var(--transition-smooth);
  font-family: var(--font-body);
}

.mobile-logout-btn:hover {
  background: rgba(239, 68, 68, 0.3);
  color: #ffffff;
}
</style>
