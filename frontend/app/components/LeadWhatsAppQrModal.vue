<script setup lang="ts">
import { ref, computed } from 'vue'

interface LeadData {
  id_rencana?: number | string
  id_customer?: number | string
  nama: string
  hp: string
  soa?: string
  ket?: string
  keterangan_barang?: string
}

const props = defineProps<{
  modelValue: boolean
  lead: LeadData | null
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'close'): void
}>()

const isCopiedLink = ref(false)
const isCopiedMessage = ref(false)
const customMessage = ref('')
const isEditingMessage = ref(false)

// Format phone number to international WhatsApp format (628xxx)
const formattedPhone = computed(() => {
  if (!props.lead?.hp) return ''
  let phone = String(props.lead.hp).replace(/\D/g, '')
  if (phone.startsWith('0')) {
    phone = '62' + phone.substring(1)
  } else if (phone.startsWith('8')) {
    phone = '62' + phone
  } else if (!phone.startsWith('62') && phone.length > 5) {
    phone = '62' + phone
  }
  return phone
})

const defaultMessage = computed(() => {
  const leadId = props.lead?.id_rencana || props.lead?.id_customer || ''
  const idTag = leadId ? ` (No. Rencana #${leadId})` : ''
  const barangTag = props.lead?.keterangan_barang ? ` untuk estimasi taksiran barang ${props.lead.keterangan_barang}` : ''
  return `Halo Kak ${props.lead?.nama || 'Nasabah'}, terima kasih telah menghubungi Selada Gadai. Pengajuan Leads Anda telah kami terima${idTag}${barangTag}. Mohon konfirmasi kebutuhan taksiran / gadai Anda untuk kami proses lebih lanjut.`
})

const activeMessage = computed(() => {
  return customMessage.value || defaultMessage.value
})

const waUrl = computed(() => {
  if (!formattedPhone.value) return ''
  const encoded = encodeURIComponent(activeMessage.value)
  return `https://wa.me/${formattedPhone.value}?text=${encoded}`
})

const qrCodeUrl = computed(() => {
  if (!waUrl.value) return ''
  return `https://api.qrserver.com/v1/create-qr-code/?size=260x260&margin=8&data=${encodeURIComponent(waUrl.value)}`
})

const copyLink = async () => {
  if (!waUrl.value) return
  try {
    await navigator.clipboard.writeText(waUrl.value)
    isCopiedLink.value = true
    setTimeout(() => {
      isCopiedLink.value = false
    }, 2000)
  } catch (err) {
    console.error('Failed to copy WA url', err)
  }
}

const copyMessage = async () => {
  if (!activeMessage.value) return
  try {
    await navigator.clipboard.writeText(activeMessage.value)
    isCopiedMessage.value = true
    setTimeout(() => {
      isCopiedMessage.value = false
    }, 2000)
  } catch (err) {
    console.error('Failed to copy message text', err)
  }
}

const openWhatsApp = () => {
  if (waUrl.value) {
    window.open(waUrl.value, '_blank', 'noopener,noreferrer')
  }
}

const closeModal = () => {
  emit('update:modelValue', false)
  emit('close')
  isEditingMessage.value = false
  customMessage.value = ''
}
</script>

<template>
  <ClientOnly>
    <Teleport to="body">
      <div
        v-if="modelValue && lead"
        class="fixed inset-0 w-screen h-screen min-h-screen top-0 left-0 right-0 bottom-0 z-[99999] bg-slate-900/60 backdrop-blur-sm flex items-center justify-center p-4 !mt-0 overflow-y-auto"
        @click.self="closeModal"
      >
        <div class="bg-white rounded-3xl max-w-lg w-full p-6 shadow-2xl border border-slate-100 space-y-5 animate-in fade-in zoom-in-95 duration-200 my-auto">
          <!-- Header -->
          <div class="flex items-start justify-between border-b border-slate-100 pb-4">
            <div class="flex items-center gap-3">
              <div class="w-11 h-11 rounded-2xl bg-emerald-500 text-white flex items-center justify-center shadow-lg shadow-emerald-500/20">
                <svg class="w-6 h-6 fill-current" viewBox="0 0 24 24">
                  <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 0 1-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 0 1-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 0 1 2.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0 0 12.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 0 0 5.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 0 0-3.48-8.413Z" />
                </svg>
              </div>
              <div>
                <h3 class="font-bold text-slate-800 text-lg leading-tight">WhatsApp QR Leads</h3>
                <p class="text-xs text-slate-500 mt-0.5">Scan QR untuk terhubung langsung ke WhatsApp Leads</p>
              </div>
            </div>
            <button
              @click="closeModal"
              type="button"
              class="w-8 h-8 rounded-full bg-slate-100 hover:bg-slate-200 text-slate-400 hover:text-slate-600 flex items-center justify-center transition-all cursor-pointer"
            >
              <AppIcon name="close" class="w-4 h-4" />
            </button>
          </div>

          <!-- Customer Info Chip -->
          <div class="bg-gradient-to-r from-emerald-50 to-teal-50 border border-emerald-100 rounded-2xl p-4 flex items-center justify-between">
            <div>
              <div class="text-[11px] font-bold uppercase tracking-wider text-emerald-800">Customer Leads</div>
              <div class="text-base font-bold text-slate-800 mt-0.5">{{ lead.nama }}</div>
              <div class="text-xs text-slate-600 font-mono mt-0.5 flex items-center gap-1.5">
                <AppIcon name="phone" class="w-3.5 h-3.5 text-emerald-600" />
                <span>{{ lead.hp }}</span>
                <span class="text-emerald-700 font-semibold" v-if="lead.soa">• SOA: {{ lead.soa }}</span>
              </div>
            </div>
            <div class="text-right">
              <span class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-bold bg-white text-emerald-800 border border-emerald-200 shadow-2xs">
                #{{ lead.id_rencana || lead.id_customer || 'Baru' }}
              </span>
            </div>
          </div>

          <!-- QR Code Display Area -->
          <div class="flex flex-col items-center justify-center py-2">
            <div class="p-3 bg-white rounded-2xl border-2 border-emerald-500/30 shadow-md relative group">
              <img
                v-if="qrCodeUrl"
                :src="qrCodeUrl"
                :alt="`QR WhatsApp ${lead.nama}`"
                class="w-56 h-56 rounded-xl object-contain"
                loading="eager"
              />
              <div
                v-else
                class="w-56 h-56 rounded-xl bg-slate-100 flex items-center justify-center text-slate-400 text-xs font-medium"
              >
                Nomor HP tidak valid
              </div>

              <!-- Center Logo Accent -->
              <div class="absolute inset-0 flex items-center justify-center pointer-events-none">
                <div class="w-10 h-10 rounded-full bg-white shadow-md border border-emerald-200 flex items-center justify-center text-emerald-600">
                  <svg class="w-5 h-5 fill-current" viewBox="0 0 24 24">
                    <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 0 1-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 0 1-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 0 1 2.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0 0 12.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 0 0 5.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 0 0-3.48-8.413Z" />
                  </svg>
                </div>
              </div>
            </div>

            <p class="text-[11px] text-slate-500 mt-2 text-center">
              Scan menggunakan kamera HP nasabah atau tim marketing untuk membuka obrolan langsung.
            </p>
          </div>

          <!-- Message Preview & Edit -->
          <div class="bg-slate-50 border border-slate-200 rounded-xl p-3 text-xs space-y-1.5">
            <div class="flex items-center justify-between text-slate-600 font-semibold">
              <span class="flex items-center gap-1">
                <AppIcon name="message" class="w-3.5 h-3.5 text-emerald-600" />
                Template Pesan WhatsApp
              </span>
              <button
                type="button"
                @click="isEditingMessage = !isEditingMessage"
                class="text-blue-600 hover:text-blue-700 text-[11px] font-medium cursor-pointer"
              >
                {{ isEditingMessage ? 'Selesai' : 'Edit Pesan' }}
              </button>
            </div>

            <textarea
              v-if="isEditingMessage"
              v-model="customMessage"
              rows="3"
              class="w-full p-2 bg-white border border-slate-300 rounded-lg text-xs outline-none focus:border-emerald-500"
              :placeholder="defaultMessage"
            ></textarea>
            <p v-else class="text-slate-700 text-xs italic line-clamp-3 leading-relaxed">
              "{{ activeMessage }}"
            </p>
          </div>

          <!-- Webhook Information Alert -->
          <div class="p-3 bg-blue-50 border border-blue-200 rounded-xl text-xs text-blue-900 flex items-start gap-2.5">
            <AppIcon name="sparkles" class="w-4 h-4 text-blue-600 flex-shrink-0 mt-0.5" />
            <div class="leading-relaxed">
              <strong class="font-bold">Otomatisasi Webhook Aktif:</strong>
              Saat customer merespon pesan ini (mengandung kode ID rencana), webhook WhatsApp akan secara otomatis meng-upgrade status Leads ini menjadi <strong>Prospect</strong>.
            </div>
          </div>

          <!-- Actions -->
          <div class="flex flex-col sm:flex-row items-center gap-2.5 pt-1">
            <button
              @click="openWhatsApp"
              type="button"
              class="w-full sm:flex-1 py-3 px-4 bg-emerald-600 hover:bg-emerald-700 active:scale-95 text-white font-bold rounded-xl text-sm flex items-center justify-center gap-2 transition-all shadow-md shadow-emerald-600/20 cursor-pointer"
            >
              <svg class="w-4 h-4 fill-current" viewBox="0 0 24 24">
                <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 0 1-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 0 1-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 0 1 2.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0 0 12.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 0 0 5.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 0 0-3.48-8.413Z" />
              </svg>
              <span>Buka WhatsApp Langsung</span>
            </button>

            <button
              @click="copyLink"
              type="button"
              class="w-full sm:w-auto py-3 px-4 border border-slate-300 hover:bg-slate-50 text-slate-700 font-semibold rounded-xl text-xs flex items-center justify-center gap-1.5 transition-all cursor-pointer"
            >
              <AppIcon :name="isCopiedLink ? 'check' : 'external-link'" class="w-4 h-4" :class="{ 'text-emerald-600': isCopiedLink }" />
              <span>{{ isCopiedLink ? 'Tersalin!' : 'Salin Link' }}</span>
            </button>

            <button
              @click="copyMessage"
              type="button"
              class="w-full sm:w-auto py-3 px-4 border border-slate-300 hover:bg-slate-50 text-slate-700 font-semibold rounded-xl text-xs flex items-center justify-center gap-1.5 transition-all cursor-pointer"
            >
              <AppIcon :name="isCopiedMessage ? 'check' : 'document'" class="w-4 h-4" :class="{ 'text-emerald-600': isCopiedMessage }" />
              <span>{{ isCopiedMessage ? 'Tersalin!' : 'Salin Pesan' }}</span>
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </ClientOnly>
</template>
