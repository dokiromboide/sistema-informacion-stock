<template>
  <div class="min-h-screen bg-gray-50">
    <nav class="bg-white border-b border-gray-200 shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center gap-2">
            <span class="text-blue-700 text-2xl font-bold">Stock</span>
            <span class="text-gray-500 text-sm font-medium">Sistema de Informacion</span>
          </div>
          <div class="flex gap-6">
            <RouterLink
              to="/"
              class="text-gray-600 hover:text-blue-700 text-sm font-medium transition-colors"
              active-class="text-blue-700 font-semibold"
            >
              Acciones
            </RouterLink>
            <RouterLink
              to="/recomendaciones"
              class="text-gray-600 hover:text-blue-700 text-sm font-medium transition-colors"
              active-class="text-blue-700 font-semibold"
            >
              Recomendaciones
            </RouterLink>
          </div>
        </div>
      </div>
    </nav>

    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>
