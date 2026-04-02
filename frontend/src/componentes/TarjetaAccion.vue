<template>
  <RouterLink
    :to="{ name: 'detalle-accion', params: { simbolo: accion.simbolo } }"
    class="tarjeta hover:shadow-md transition-shadow duration-200 block"
  >
    <div class="flex items-start justify-between mb-3">
      <div>
        <span class="text-lg font-bold text-gray-900">{{ accion.simbolo }}</span>
        <p class="text-sm text-gray-500 truncate max-w-[180px]">{{ accion.nombre || 'Sin nombre' }}</p>
      </div>
      <span
        class="text-xs font-semibold px-2 py-1 rounded-full"
        :class="accion.cambio_pct >= 0 ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'"
      >
        {{ accion.cambio_pct >= 0 ? '+' : '' }}{{ accion.cambio_pct.toFixed(2) }}%
      </span>
    </div>

    <div class="flex items-end justify-between">
      <div>
        <p class="text-2xl font-bold text-gray-900">${{ accion.precio.toFixed(2) }}</p>
        <p
          class="text-sm font-medium"
          :class="accion.cambio_monto >= 0 ? 'etiqueta-positivo' : 'etiqueta-negativo'"
        >
          {{ accion.cambio_monto >= 0 ? '+' : '' }}{{ accion.cambio_monto.toFixed(2) }}
        </p>
      </div>
      <div class="text-right text-xs text-gray-400">
        <p>Vol: {{ formatearVolumen(accion.volumen) }}</p>
      </div>
    </div>
  </RouterLink>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
import type { Accion } from '@/tipos'

defineProps<{ accion: Accion }>()

function formatearVolumen(volumen: number): string {
  if (volumen >= 1_000_000) return (volumen / 1_000_000).toFixed(1) + 'M'
  if (volumen >= 1_000) return (volumen / 1_000).toFixed(1) + 'K'
  return volumen.toString()
}
</script>
