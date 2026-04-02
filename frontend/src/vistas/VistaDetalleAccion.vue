<template>
  <div>
    <!-- Navegacion de retorno -->
    <RouterLink
      to="/"
      class="inline-flex items-center gap-1 text-sm text-blue-600 hover:text-blue-800 mb-6"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
      </svg>
      Volver a acciones
    </RouterLink>

    <!-- Estado de carga -->
    <div v-if="almacen.cargando" class="tarjeta animate-pulse">
      <div class="h-8 bg-gray-200 rounded w-24 mb-2"></div>
      <div class="h-4 bg-gray-100 rounded w-48 mb-6"></div>
      <div class="h-12 bg-gray-200 rounded w-36"></div>
    </div>

    <!-- Error -->
    <div v-else-if="almacen.error" class="tarjeta text-red-600">
      {{ almacen.error }}
    </div>

    <!-- Detalle de la accion -->
    <div v-else-if="almacen.accionActual">
      <div class="tarjeta mb-6">
        <div class="flex items-start justify-between mb-4">
          <div>
            <h1 class="text-3xl font-bold text-gray-900">{{ almacen.accionActual.simbolo }}</h1>
            <p class="text-gray-500 mt-1">{{ almacen.accionActual.nombre || 'Sin nombre registrado' }}</p>
          </div>
          <span
            class="text-sm font-semibold px-3 py-1.5 rounded-full"
            :class="almacen.accionActual.cambio_pct >= 0 ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'"
          >
            {{ almacen.accionActual.cambio_pct >= 0 ? '+' : '' }}{{ almacen.accionActual.cambio_pct.toFixed(2) }}%
          </span>
        </div>

        <p class="text-4xl font-bold text-gray-900 mb-1">
          ${{ almacen.accionActual.precio.toFixed(2) }}
        </p>
        <p
          class="text-base font-medium"
          :class="almacen.accionActual.cambio_monto >= 0 ? 'etiqueta-positivo' : 'etiqueta-negativo'"
        >
          {{ almacen.accionActual.cambio_monto >= 0 ? '+' : '' }}{{ almacen.accionActual.cambio_monto.toFixed(2) }} hoy
        </p>
      </div>

      <!-- Estadisticas del dia -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div class="tarjeta text-center">
          <p class="text-xs text-gray-500 uppercase tracking-wide mb-1">Apertura</p>
          <p class="text-lg font-bold text-gray-900">${{ almacen.accionActual.apertura.toFixed(2) }}</p>
        </div>
        <div class="tarjeta text-center">
          <p class="text-xs text-gray-500 uppercase tracking-wide mb-1">Maximo</p>
          <p class="text-lg font-bold text-green-600">${{ almacen.accionActual.maximo.toFixed(2) }}</p>
        </div>
        <div class="tarjeta text-center">
          <p class="text-xs text-gray-500 uppercase tracking-wide mb-1">Minimo</p>
          <p class="text-lg font-bold text-red-600">${{ almacen.accionActual.minimo.toFixed(2) }}</p>
        </div>
        <div class="tarjeta text-center">
          <p class="text-xs text-gray-500 uppercase tracking-wide mb-1">Volumen</p>
          <p class="text-lg font-bold text-gray-900">{{ formatearVolumen(almacen.accionActual.volumen) }}</p>
        </div>
      </div>

      <p class="text-xs text-gray-400 mt-4 text-right">
        Ultima actualizacion: {{ formatearFecha(almacen.accionActual.ultima_actualizacion) }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { useAccionesAlmacen } from '@/almacen/accionesAlmacen'

const ruta = useRoute()
const almacen = useAccionesAlmacen()

onMounted(() => {
  const simbolo = ruta.params.simbolo as string
  almacen.cargarDetalle(simbolo)
})

function formatearVolumen(volumen: number): string {
  if (volumen >= 1_000_000) return (volumen / 1_000_000).toFixed(2) + 'M'
  if (volumen >= 1_000) return (volumen / 1_000).toFixed(1) + 'K'
  return volumen.toLocaleString('es')
}

function formatearFecha(fecha: string): string {
  return new Date(fecha).toLocaleString('es-CO', {
    dateStyle: 'medium',
    timeStyle: 'short'
  })
}
</script>
