<template>
  <div>
    <!-- Encabezado de la pagina -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Acciones del Mercado</h1>
        <p class="text-sm text-gray-500 mt-1">
          {{ almacen.totalAcciones }} acciones disponibles
        </p>
      </div>
      <button
        class="btn-primario flex items-center gap-2"
        :disabled="almacen.sincronizando"
        @click="almacen.sincronizarDatos()"
      >
        <svg
          class="w-4 h-4"
          :class="{ 'animate-spin': almacen.sincronizando }"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
        {{ almacen.sincronizando ? 'Sincronizando...' : 'Sincronizar' }}
      </button>
    </div>

    <!-- Mensaje de exito de sincronizacion -->
    <div
      v-if="almacen.mensajeSincronizacion"
      class="mb-4 px-4 py-3 bg-green-50 border border-green-200 text-green-700 rounded-lg text-sm"
    >
      {{ almacen.mensajeSincronizacion }}
    </div>

    <!-- Controles de busqueda y ordenamiento -->
    <div class="flex flex-col sm:flex-row gap-3 mb-6">
      <BarraBusqueda
        class="flex-1"
        @buscar="(q) => almacen.actualizarFiltro({ q })"
      />
      <SelectorOrden
        @cambiar="(columna, dir) => almacen.actualizarFiltro({ ordenar_por: columna, direccion: dir })"
      />
    </div>

    <!-- Estado de error -->
    <div
      v-if="almacen.error"
      class="mb-4 px-4 py-3 bg-red-50 border border-red-200 text-red-700 rounded-lg text-sm flex justify-between items-center"
    >
      <span>{{ almacen.error }}</span>
      <button class="text-red-500 hover:text-red-700" @click="almacen.limpiarError()">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <!-- Estado de carga -->
    <div v-if="almacen.cargando" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      <div v-for="i in 8" :key="i" class="tarjeta animate-pulse">
        <div class="h-4 bg-gray-200 rounded w-16 mb-2"></div>
        <div class="h-3 bg-gray-100 rounded w-32 mb-4"></div>
        <div class="h-6 bg-gray-200 rounded w-24"></div>
      </div>
    </div>

    <!-- Lista de acciones -->
    <div
      v-else-if="almacen.hayAcciones"
      class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4"
    >
      <TarjetaAccion
        v-for="accion in almacen.acciones"
        :key="accion.simbolo"
        :accion="accion"
      />
    </div>

    <!-- Estado vacio -->
    <div v-else class="text-center py-16 text-gray-400">
      <svg class="w-12 h-12 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
          d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
      </svg>
      <p class="text-lg font-medium text-gray-500">No hay acciones disponibles</p>
      <p class="text-sm mt-1">Haz clic en "Sincronizar" para cargar datos del mercado</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAccionesAlmacen } from '@/almacen/accionesAlmacen'
import TarjetaAccion from '@/componentes/TarjetaAccion.vue'
import BarraBusqueda from '@/componentes/BarraBusqueda.vue'
import SelectorOrden from '@/componentes/SelectorOrden.vue'

const almacen = useAccionesAlmacen()

onMounted(() => {
  almacen.cargarAcciones()
})
</script>
