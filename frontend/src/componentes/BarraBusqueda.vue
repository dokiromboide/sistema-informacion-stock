<template>
  <div class="relative">
    <input
      v-model="terminoBusqueda"
      type="text"
      :placeholder="placeholder"
      class="w-full pl-10 pr-4 py-2.5 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
      @input="alEscribir"
    />
    <svg
      class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400"
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
        d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
    </svg>
    <button
      v-if="terminoBusqueda"
      class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
      @click="limpiar"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
      </svg>
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const props = withDefaults(defineProps<{
  placeholder?: string
  debounceMs?: number
}>(), {
  placeholder: 'Buscar por simbolo o nombre...',
  debounceMs: 400
})

const emit = defineEmits<{
  buscar: [valor: string]
}>()

const terminoBusqueda = ref('')
let temporizador: ReturnType<typeof setTimeout> | null = null

function alEscribir() {
  if (temporizador) clearTimeout(temporizador)
  temporizador = setTimeout(() => {
    emit('buscar', terminoBusqueda.value)
  }, props.debounceMs)
}

function limpiar() {
  terminoBusqueda.value = ''
  emit('buscar', '')
}
</script>
