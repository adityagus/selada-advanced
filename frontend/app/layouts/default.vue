<template>
  <MobileHeader :userName="userName || 'User'" @logout="handleLogout" />
  <div class="app-container">
    <!-- Sidebar (Desktop only) -->
    <aside class="sidebar" :class="{ expanded: isExpanded }">
      <!-- Brand -->
      <div class="sidebar-brand">
        <div class="brand-row">
          <div class="brand-logo">
            <NuxtImg src="/assets/img/logo-gadaimulia.png" alt="Logo Optimasi" style="width: 40px;" />
          </div>
          <span v-show="isExpanded" class="brand-name">SELADA</span>
        </div>
        <button @click="toggleSidebar" class="sidebar-toggle-btn" :class="{ expanded: isExpanded }">
          <svg :class="['toggle-icon', isExpanded ? 'rotated' : '']" width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M15 18L9 12L15 6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
        </button>
      </div>

      <!-- Navigation -->
      <nav class="sidebar-nav">
        <NuxtLink to="/dashboard" class="nav-link" active-class="nav-active" :class="{ 'icon-only': !isExpanded }">
          <svg class="nav-icon" width="22" height="22" viewBox="0 0 24 24" fill="none"><path d="M3 9L12 2L21 9V20C21 20.5304 20.7893 21.0391 20.4142 21.4142C20.0391 21.7893 19.5304 22 19 22H5C4.46957 22 3.96086 21.7893 3.58579 21.4142C3.21071 21.0391 3 20.5304 3 20V9Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
          <span v-show="isExpanded" class="nav-text">Dashboard</span>
        </NuxtLink>

        <NuxtLink to="/inquiry" class="nav-link" active-class="nav-active" :class="{ 'icon-only': !isExpanded }">
          <svg class="nav-icon" width="22" height="22" viewBox="0 0 24 24" fill="none"><rect x="3" y="3" width="7" height="7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><rect x="14" y="3" width="7" height="7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><rect x="14" y="14" width="7" height="7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><rect x="3" y="14" width="7" height="7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
          <span v-show="isExpanded" class="nav-text">Inquiry</span>
        </NuxtLink>

        <NuxtLink to="/data-customer" class="nav-link" active-class="nav-active" :class="{ 'icon-only': !isExpanded }">
          <svg class="nav-icon" width="22" height="22" viewBox="0 0 24 24" fill="none"><path d="M17 21V19C17 17.9391 16.5786 16.9217 15.8284 16.1716C15.0783 15.4214 14.0609 15 13 15H5C3.93913 15 2.92172 15.4214 2.17157 16.1716C1.42143 16.9217 1 17.9391 1 19V21" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M9 11C11.2091 11 13 9.20914 13 7C13 4.79086 11.2091 3 9 3C6.79086 3 5 4.79086 5 7C5 9.20914 6.79086 11 9 11Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
          <span v-show="isExpanded" class="nav-text">Data Customer</span>
        </NuxtLink>

        <NuxtLink to="/cek-taksiran" class="nav-link" active-class="nav-active" :class="{ 'icon-only': !isExpanded }">
          <svg class="nav-icon" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" width="22" height="22"><path stroke-linecap="round" stroke-linejoin="round" d="M12 3v17.25m0 0c-1.472 0-2.882.265-4.185.75M12 20.25c1.472 0 2.882.265 4.185.75M18.75 4.97A48.416 48.416 0 0 0 12 4.5c-2.291 0-4.545.16-6.75.47m13.5 0c1.01.143 2.01.317 3 .52m-3-.52 2.62 10.726c.122.499-.106 1.028-.589 1.202a5.988 5.988 0 0 1-2.031.352 5.988 5.988 0 0 1-2.031-.352c-.483-.174-.711-.703-.59-1.202L18.75 4.971Zm-16.5.52c.99-.203 1.99-.377 3-.52m0 0 2.62 10.726c.122.499-.106 1.028-.589 1.202a5.989 5.989 0 0 1-2.031.352 5.989 5.989 0 0 1-2.031-.352c-.483-.174-.711-.703-.59-1.202L5.25 4.971Z" /></svg>
          <span v-show="isExpanded" class="nav-text">CeTar</span>
        </NuxtLink>

        <NuxtLink to="/report" class="nav-link" active-class="nav-active" :class="{ 'icon-only': !isExpanded }">
          <svg class="nav-icon" width="22" height="22" viewBox="0 0 24 24" fill="none"><path d="M18 20V10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M12 20V4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M6 20V14" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
          <span v-show="isExpanded" class="nav-text">Report</span>
        </NuxtLink>
      </nav>

      <!-- User Footer -->
      <div class="sidebar-footer">
        <div class="user-info">
          <div class="user-avatar">
            <img src="https://api.dicebear.com/7.x/avataaars/svg?seed=user" alt="User" />
          </div>
          <div v-show="isExpanded" class="user-meta">
            <span class="user-name">{{ userName || 'User' }}</span>
            <span class="user-branch">Cabang {{ userBranch || '-' }}</span>
          </div>
        </div>
        <button @click="handleLogout" class="logout-btn" :title="isExpanded ? '' : 'Logout'">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M15.75 9V5.25A2.25 2.25 0 0 0 13.5 3h-6a2.25 2.25 0 0 0-2.25 2.25v13.5A2.25 2.25 0 0 0 7.5 21h6a2.25 2.25 0 0 0 2.25-2.25V15M12 9l-3 3m0 0 3 3m-3-3h12.75" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </button>
      </div>
    </aside>

    <!-- Main Content -->
    <div class="main-content" :style="{ marginLeft: mainMargin }">
      <slot />
    </div>

    <!-- Bottom Nav (Mobile only) -->
    <BottomNav />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useCookie, navigateTo } from '#app'
import { useApi } from '~/composables/useApi'

const token = useCookie('auth_token')
const userName = useCookie('user_name')
const userBranch = useCookie('user_branch')

const isExpanded = ref(false)

const mainMargin = computed(() => {
  if (typeof window !== 'undefined' && window.innerWidth < 768) return '0'
  return isExpanded.value ? '256px' : '80px'
})

const toggleSidebar = () => {
  isExpanded.value = !isExpanded.value
}

onMounted(() => {
  if (!token.value) {
    navigateTo('/login')
  }
})

const handleLogout = async () => {
  await useApi('/logout', { method: 'POST' })
  token.value = null
  userName.value = null
  userBranch.value = null
  navigateTo('/login')
}
</script>

<style scoped>
/* --- Sidebar --- */
.sidebar {
  display: none;
  position: fixed;
  left: 0;
  top: 0;
  height: 100vh;
  width: 80px;
  background: linear-gradient(180deg, var(--selada-400), var(--selada-500));
  flex-direction: column;
  padding: 1.25rem 1rem;
  z-index: 50;
  transition: width 0.3s ease;
  overflow: hidden;
}

@media (min-width: 768px) {
  .sidebar {
    display: flex;
  }
}

.sidebar.expanded {
  width: 256px;
}

/* --- Brand --- */
.sidebar-brand {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 2rem;
  padding: 0 0.25rem;
}

.brand-row {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.brand-logo {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.brand-name {
  color: #ffffff;
  font-weight: 700;
  font-size: 1.15rem;
  white-space: nowrap;
}

.sidebar-toggle-btn {
  width: 30px;
  height: 30px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.2);
  border: none;
  color: rgba(255, 255, 255, 0.9);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: var(--transition-smooth);
  flex-shrink: 0;
}

.sidebar-toggle-btn:hover {
  background: rgba(255, 255, 255, 0.8);
  color: var(--selada-500);
}

.toggle-icon {
  transition: transform 0.3s ease;
}

.toggle-icon.rotated {
  transform: rotate(180deg);
}

/* --- Navigation --- */
.sidebar-nav {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
  flex: 1;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  border-radius: 12px;
  color: rgba(255, 255, 255, 0.75);
  transition: var(--transition-smooth);
  white-space: nowrap;
}

.nav-link.icon-only {
  width: 48px;
  height: 48px;
  justify-content: center;
  padding: 0;
}

.nav-link:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #ffffff;
}

.nav-link.nav-active {
  background: #ffffff;
  color: var(--selada-400);
  box-shadow: 0 4px 6px -1px rgba(249, 115, 22, 0.1), 0 2px 4px -1px rgba(249, 115, 22, 0.06);
}

.nav-icon {
  flex-shrink: 0;
}

.nav-text {
  font-weight: 500;
  font-size: 0.95rem;
}

/* --- Sidebar Footer --- */
.sidebar-footer {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
  padding-top: 1rem;
  border-top: 1px solid rgba(255, 255, 255, 0.2);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  width: 100%;
  justify-content: center;
}

.user-avatar {
  width: 38px;
  height: 38px;
  border-radius: 50%;
  overflow: hidden;
  border: 2px solid #ffffff;
  flex-shrink: 0;
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.user-meta {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.user-name {
  color: #ffffff;
  font-size: 0.85rem;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-branch {
  color: rgba(255, 255, 255, 0.7);
  font-size: 0.7rem;
}

.logout-btn {
  width: 38px;
  height: 38px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: none;
  border: none;
  color: #ffffff;
  cursor: pointer;
  opacity: 0.8;
  transition: var(--transition-smooth);
}

.logout-btn:hover {
  opacity: 1;
}

/* --- Main Content Responsive --- */
@media (max-width: 767px) {
  .main-content {
    margin-left: 0 !important;
    padding-bottom: 6rem;
  }
}
</style>
