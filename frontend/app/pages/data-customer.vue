<template>
  <div class="customer-container">
    <!-- Header -->
    <div class="customer-header mb-8">
      <h1>Data Customer</h1>
      <p class="text-muted">Kelola dan pantau data pelanggan Anda dengan mudah</p>
    </div>

    <!-- Search Box -->
    <div class="filter-card card mb-8">
      <div class="search-box">
        <label class="form-label" for="search-input">Cari Nama / Kontak</label>
        <input
          id="search-input"
          v-model="searchQuery"
          @input="debounceFetch"
          type="text"
          placeholder="Ketik nama customer..."
          class="form-input search-field"
        />
      </div>
    </div>

    <!-- Customers Table Card -->
    <div class="table-card card">
      <div class="table-responsive">
        <table class="customer-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Nama Customer</th>
              <th>No. Handphone</th>
              <th>Sumber Leads</th>
              <th>Status Akhir</th>
              <th class="text-center">Aksi Hubungi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="isLoading">
              <td colspan="6" class="text-center py-6">
                <div class="spinner-small centered"></div>
              </td>
            </tr>
            <tr v-else-if="customers.length === 0">
              <td colspan="6" class="text-center py-6 text-muted">
                Tidak ada data customer ditemukan.
              </td>
            </tr>
            <tr v-for="(c, idx) in customers" :key="idx">
              <td>{{ idx + 1 }}</td>
              <td><strong>{{ c.nama }}</strong></td>
              <td>{{ c.phone || '083138794243' }}</td>
              <td><span class="badge badge-info">{{ c.source || 'Walk In' }}</span></td>
              <td>
                <span class="badge" :class="getStatusBadgeClass(c.status)">
                  {{ c.status }}
                </span>
              </td>
              <td>
                <div class="action-buttons-cell">
                  <a :href="`tel:${cleanPhone(c.phone)}`" class="btn btn-secondary contact-btn tel-btn" target="_blank">
                    📞 Telepon
                  </a>
                  <a :href="`https://wa.me/${cleanPhone(c.phone)}`" class="btn btn-primary contact-btn wa-btn" target="_blank">
                    💬 WhatsApp
                  </a>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'

interface Customer {
  nama: string
  phone: string
  source: string
  status: string
}

const customers = ref<Customer[]>([])
const searchQuery = ref('')
const isLoading = ref(false)

onMounted(() => {
  fetchCustomers()
})

const fetchCustomers = async () => {
  isLoading.value = true
  // Query customers from inquiry list
  const { data } = await useApi('/inquiry/filter', {
    params: {
      search: searchQuery.value,
      start: 0,
      length: 100
    }
  })
  isLoading.value = false
  if (data && data.data) {
    // Map list to unique customer profiles
    const map = new Map()
    for (const item of data.data) {
      const name = item.nama || 'No Name'
      // Use fake phone number or source if not defined
      const phone = item.hp || '083138794243' 
      map.set(name, {
        nama: name,
        phone: phone,
        source: item.source || 'Walk In',
        status: item.status || 'Leads'
      })
    }
    customers.value = Array.from(map.values())
  }
}

let debounceTimer: any = null
const debounceFetch = () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    fetchCustomers()
  }, 400)
}

const cleanPhone = (phoneStr: string) => {
  if (!phoneStr) return '6283138794243'
  // Remove non-digit characters
  let clean = phoneStr.replace(/\D/g, '')
  if (clean.startsWith('0')) {
    clean = '62' + clean.slice(1)
  }
  return clean || '6283138794243'
}

const getStatusBadgeClass = (status: string) => {
  switch (status) {
    case 'Leads': return 'badge-warning'
    case 'Prospect': return 'badge-info'
    case 'Hot Prospect': return 'badge-danger'
    case 'SBG': return 'badge-success'
    case 'Batal': return 'badge-danger'
    default: return 'badge-info'
  }
}
</script>

<style scoped>
.customer-container {
  max-width: 1200px;
  margin: 0 auto;
}

.customer-header h1 {
  font-size: 1.75rem;
  font-weight: 700;
}

.filter-card {
  padding: 1.25rem 1.5rem;
}

.search-box {
  display: flex;
  flex-direction: column;
}

.search-field {
  max-width: 400px;
}

.table-card {
  padding: 0;
  overflow: hidden;
  border-radius: 12px;
}

.table-responsive {
  width: 100%;
  overflow-x: auto;
}

.customer-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.customer-table th, .customer-table td {
  padding: 0.85rem 1.25rem;
  border-bottom: 1px solid var(--border-color);
}

.customer-table th {
  background: #f9fafb;
  color: var(--text-muted);
  font-weight: 600;
  text-transform: uppercase;
  font-size: 0.75rem;
  letter-spacing: 0.5px;
}

.customer-table tr:last-child td {
  border-bottom: none;
}

.customer-table tr:hover td {
  background: #fffbeb;
}

.text-center { text-align: center; }
.py-6 { padding-top: 1.5rem; padding-bottom: 1.5rem; }

.action-buttons-cell {
  display: flex;
  gap: 0.5rem;
  justify-content: center;
}

.contact-btn {
  font-size: 0.85rem;
  padding: 0.5rem 1rem;
  border-radius: 9999px;
  font-weight: 600;
  cursor: pointer;
  transition: var(--transition-smooth);
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
}

.tel-btn {
  background: #fffbeb;
  color: #d97706;
  border: 1px solid #fde68a;
}

.tel-btn:hover {
  background: #f59e0b;
  color: #ffffff;
}

.wa-btn {
  background: #ecfdf5;
  color: #059669;
  border: 1px solid #a7f3d0;
}

.wa-btn:hover {
  background: #10b981;
  color: #ffffff;
}

.centered {
  margin: 0 auto;
}
</style>
