import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Accion, FiltroAcciones } from '@/tipos'
import { servicioAcciones } from '@/servicios/apiStock'

export const useAccionesAlmacen = defineStore('acciones', () => {
  // Estado
  const acciones = ref<Accion[]>([])
  const accionActual = ref<Accion | null>(null)
  const cargando = ref(false)
  const error = ref<string | null>(null)
  const paginaActual = ref(1)
  const sincronizando = ref(false)
  const mensajeSincronizacion = ref<string | null>(null)

  const filtro = ref<FiltroAcciones>({
    q: '',
    ordenar_por: 'simbolo',
    direccion: 'asc',
    pagina: 1,
    por_pagina: 20
  })

  // Getters
  const hayAcciones = computed(() => acciones.value.length > 0)
  const totalAcciones = computed(() => acciones.value.length)

  // Acciones del almacen
  async function cargarAcciones() {
    cargando.value = true
    error.value = null
    try {
      const respuesta = await servicioAcciones.obtenerTodas(filtro.value)
      acciones.value = respuesta.datos ?? []
      paginaActual.value = respuesta.pagina
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Error desconocido'
      acciones.value = []
    } finally {
      cargando.value = false
    }
  }

  async function cargarDetalle(simbolo: string) {
    cargando.value = true
    error.value = null
    accionActual.value = null
    try {
      accionActual.value = await servicioAcciones.obtenerDetalle(simbolo)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Error al cargar detalle'
    } finally {
      cargando.value = false
    }
  }

  async function sincronizarDatos() {
    sincronizando.value = true
    mensajeSincronizacion.value = null
    error.value = null
    try {
      const resultado = await servicioAcciones.sincronizar()
      mensajeSincronizacion.value = `${resultado.mensaje}: ${resultado.sincronizados} acciones actualizadas`
      await cargarAcciones()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Error al sincronizar'
    } finally {
      sincronizando.value = false
    }
  }

  function actualizarFiltro(nuevoFiltro: Partial<FiltroAcciones>) {
    filtro.value = { ...filtro.value, ...nuevoFiltro, pagina: 1 }
    cargarAcciones()
  }

  function limpiarError() {
    error.value = null
  }

  return {
    acciones,
    accionActual,
    cargando,
    error,
    paginaActual,
    sincronizando,
    mensajeSincronizacion,
    filtro,
    hayAcciones,
    totalAcciones,
    cargarAcciones,
    cargarDetalle,
    sincronizarDatos,
    actualizarFiltro,
    limpiarError
  }
})
