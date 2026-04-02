<template>
  <div class="flex items-center gap-2">
    <select
      v-model="columnaSeleccionada"
      class="text-sm border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
      @change="emitirCambio"
    >
      <option v-for="opcion in opciones" :key="opcion.valor" :value="opcion.valor">
        {{ opcion.etiqueta }}
      </option>
    </select>

    <button
      class="p-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
      :title="direccion === 'asc' ? 'Ascendente' : 'Descendente'"
      @click="alternarDireccion"
    >
      <svg class="w-4 h-4 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path
          v-if="direccion === 'asc'"
          stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12"
        />
        <path
          v-else
          stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M3 4h13M3 8h9m-9 4h9m5-4v12m0 0l-4-4m4 4l4-4"
        />
      </svg>
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits<{
  cambiar: [columna: string, direccion: 'asc' | 'desc']
}>()

const opciones = [
  { valor: 'simbolo', etiqueta: 'Simbolo' },
  { valor: 'precio', etiqueta: 'Precio' },
  { valor: 'cambio_pct', etiqueta: 'Variacion %' },
  { valor: 'volumen', etiqueta: 'Volumen' }
]

const columnaSeleccionada = ref('simbolo')
const direccion = ref<'asc' | 'desc'>('asc')

function alternarDireccion() {
  direccion.value = direccion.value === 'asc' ? 'desc' : 'asc'
  emitirCambio()
}

function emitirCambio() {
  emit('cambiar', columnaSeleccionada.value, direccion.value)
}
</script>
