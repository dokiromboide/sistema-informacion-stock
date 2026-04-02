<template>
  <div>
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900">Recomendaciones de Inversion</h1>
      <p class="text-sm text-gray-500 mt-1">
        Acciones seleccionadas por el algoritmo de analisis de momentum y volumen
      </p>
    </div>

    <!-- Cargando -->
    <div v-if="cargando" class="space-y-4">
      <div v-for="i in 3" :key="i" class="tarjeta animate-pulse flex items-center gap-4">
        <div class="h-12 w-12 bg-gray-200 rounded-full"></div>
        <div class="flex-1">
          <div class="h-4 bg-gray-200 rounded w-24 mb-2"></div>
          <div class="h-3 bg-gray-100 rounded w-64"></div>
        </div>
      </div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="tarjeta text-red-600 text-sm">{{ error }}</div>

    <!-- Lista de recomendaciones -->
    <div v-else-if="recomendaciones.length > 0" class="space-y-4">
      <div
        v-for="(item, indice) in recomendaciones"
        :key="item.accion.simbolo"
        class="tarjeta flex items-center gap-4"
      >
        <!-- Posicion -->
        <div
          class="w-10 h-10 rounded-full flex items-center justify-center font-bold text-sm flex-shrink-0"
          :class="colorPosicion(indice)"
        >
          {{ indice + 1 }}
        </div>

        <!-- Datos de la accion -->
        <div class="flex-1 min-w-0">
          <div class="flex items-center gap-2 mb-0.5">
            <RouterLink
              :to="{ name: 'detalle-accion', params: { simbolo: item.accion.simbolo } }"
              class="font-bold text-gray-900 hover:text-blue-700"
            >
              {{ item.accion.simbolo }}
            </RouterLink>
            <span
              class="text-xs font-semibold px-2 py-0.5 rounded-full"
              :class="item.accion.cambio_pct >= 0 ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'"
            >
              {{ item.accion.cambio_pct >= 0 ? '+' : '' }}{{ item.accion.cambio_pct.toFixed(2) }}%
            </span>
          </div>
          <p class="text-sm text-gray-500 truncate">{{ item.razon }}</p>
        </div>

        <!-- Precio y puntaje -->
        <div class="text-right flex-shrink-0">
          <p class="font-bold text-gray-900">${{ item.accion.precio.toFixed(2) }}</p>
          <p class="text-xs text-blue-600 font-medium">Puntaje: {{ item.puntaje.toFixed(1) }}</p>
        </div>
      </div>
    </div>

    <!-- Sin datos -->
    <div v-else class="text-center py-16 text-gray-400">
      <svg class="w-12 h-12 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
          d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
      </svg>
      <p class="text-lg font-medium text-gray-500">No hay recomendaciones disponibles</p>
      <p class="text-sm mt-1">Sincroniza las acciones primero desde la pagina principal</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import type { PuntajeRecomendacion } from '@/tipos'
import { servicioAcciones } from '@/servicios/apiStock'

const recomendaciones = ref<PuntajeRecomendacion[]>([])
const cargando = ref(false)
const error = ref<string | null>(null)

onMounted(async () => {
  cargando.value = true
  try {
    const respuesta = await servicioAcciones.obtenerRecomendaciones()
    recomendaciones.value = respuesta.datos ?? []
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Error al cargar recomendaciones'
  } finally {
    cargando.value = false
  }
})

function colorPosicion(indice: number): string {
  const colores = [
    'bg-yellow-100 text-yellow-700',
    'bg-gray-100 text-gray-600',
    'bg-orange-100 text-orange-700'
  ]
  return colores[indice] ?? 'bg-blue-50 text-blue-600'
}
</script>
