import axios from 'axios'
import type { Accion, FiltroAcciones, RespuestaAcciones, RespuestaRecomendaciones } from '@/tipos'

const cliente = axios.create({
  baseURL: '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Interceptor para manejo centralizado de errores
cliente.interceptors.response.use(
  (respuesta) => respuesta,
  (error) => {
    const mensaje = error.response?.data?.error ?? 'Error de conexion con el servidor'
    return Promise.reject(new Error(mensaje))
  }
)

export const servicioAcciones = {
  async obtenerTodas(filtro: FiltroAcciones = {}): Promise<RespuestaAcciones> {
    const { data } = await cliente.get<RespuestaAcciones>('/acciones', { params: filtro })
    return data
  },

  async obtenerDetalle(simbolo: string): Promise<Accion> {
    const { data } = await cliente.get<Accion>(`/acciones/${simbolo}`)
    return data
  },

  async buscar(consulta: string): Promise<Accion[]> {
    const { data } = await cliente.get<{ datos: Accion[] }>('/acciones/buscar', {
      params: { q: consulta }
    })
    return data.datos
  },

  async sincronizar(): Promise<{ mensaje: string; sincronizados: number }> {
    const { data } = await cliente.post('/sincronizar')
    return data
  },

  async obtenerRecomendaciones(): Promise<RespuestaRecomendaciones> {
    const { data } = await cliente.get<RespuestaRecomendaciones>('/recomendaciones')
    return data
  }
}
