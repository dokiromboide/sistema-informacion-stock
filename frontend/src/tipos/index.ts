// Tipos TypeScript que reflejan las estructuras del backend

export interface Accion {
  id: number
  simbolo: string
  nombre: string
  precio: number
  apertura: number
  maximo: number
  minimo: number
  volumen: number
  cambio_monto: number
  cambio_pct: number
  ultima_actualizacion: string
}

export interface PuntajeRecomendacion {
  accion: Accion
  puntaje: number
  razon: string
}

export interface RespuestaAcciones {
  datos: Accion[]
  pagina: number
}

export interface FiltroAcciones {
  q?: string
  ordenar_por?: string
  direccion?: 'asc' | 'desc'
  pagina?: number
  por_pagina?: number
}

export interface RespuestaRecomendaciones {
  datos: PuntajeRecomendacion[]
}
