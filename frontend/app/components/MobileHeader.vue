<template>
  <header class="mobile-header">
    <div class="mobile-header-inner">
      <div class="mobile-header-left">
        <span class="mobile-welcome">Selamat datang</span>
        <span class="mobile-username">{{ userName }}</span>
      </div>
      <button @click="$emit('logout')" class="mobile-logout-btn">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M9 21H5C4.46957 21 3.96086 20.7893 3.58579 20.4142C3.21071 20.0391 3 19.5304 3 19V5C3 4.46957 3.21071 3.96086 3.58579 3.58579C3.96086 3.21071 4.46957 3 5 3H9" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
          <path d="M16 17L21 12L16 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
          <path d="M21 12H9" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
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
  background: linear-gradient(135deg, var(--selada-400), var(--selada-500));
  padding: 1rem 1.5rem;
  box-shadow: 0 2px 8px rgba(245, 158, 11, 0.3);
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
  color: rgba(255, 255, 255, 0.8);
}

.mobile-username {
  font-weight: 700;
  color: #ffffff;
  font-size: 1rem;
}

.mobile-logout-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 0.75rem;
  background: rgba(255, 255, 255, 0.2);
  border: none;
  border-radius: 8px;
  color: #ffffff;
  font-size: 0.85rem;
  font-weight: 500;
  cursor: pointer;
  transition: var(--transition-smooth);
  font-family: var(--font-body);
}

.mobile-logout-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}
</style>
