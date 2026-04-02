import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAccionesAlmacen } from '@/almacen/accionesAlmacen'

// Simular el modulo de la API para no realizar llamadas HTTP reales
vi.mock('@/servicios/apiStock', () => ({
  servicioAcciones: {
    obtenerTodas: vi.fn(),
    obtenerDetalle: vi.fn(),
    sincronizar: vi.fn()
  }
}))

import { servicioAcciones } from '@/servicios/apiStock'

const accionesMock = [
  {
    id: 1, simbolo: 'AAPL', nombre: 'Apple Inc.', precio: 189.50,
    apertura: 188, maximo: 190, minimo: 187, volumen: 52000000,
    cambio_monto: 1.5, cambio_pct: 0.79, ultima_actualizacion: '2024-01-15T14:30:00Z'
  },
  {
    id: 2, simbolo: 'GOOGL', nombre: 'Alphabet Inc.', precio: 140.00,
    apertura: 138, maximo: 141, minimo: 137, volumen: 20000000,
    cambio_monto: 2.0, cambio_pct: 1.45, ultima_actualizacion: '2024-01-15T14:30:00Z'
  }
]

describe('accionesAlmacen', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('estado inicial tiene acciones vacias y sin error', () => {
    const almacen = useAccionesAlmacen()
    expect(almacen.acciones).toEqual([])
    expect(almacen.error).toBeNull()
    expect(almacen.cargando).toBe(false)
  })

  it('cargarAcciones popula el estado con los datos de la API', async () => {
    vi.mocked(servicioAcciones.obtenerTodas).mockResolvedValue({
      datos: accionesMock,
      pagina: 1
    })

    const almacen = useAccionesAlmacen()
    await almacen.cargarAcciones()

    expect(almacen.acciones).toHaveLength(2)
    expect(almacen.acciones[0].simbolo).toBe('AAPL')
    expect(almacen.error).toBeNull()
  })

  it('cargarAcciones guarda el error cuando la API falla', async () => {
    vi.mocked(servicioAcciones.obtenerTodas).mockRejectedValue(new Error('Sin conexion'))

    const almacen = useAccionesAlmacen()
    await almacen.cargarAcciones()

    expect(almacen.error).toBe('Sin conexion')
    expect(almacen.acciones).toEqual([])
  })

  it('getter hayAcciones retorna false cuando no hay datos', () => {
    const almacen = useAccionesAlmacen()
    expect(almacen.hayAcciones).toBe(false)
  })

  it('getter hayAcciones retorna true cuando hay acciones cargadas', async () => {
    vi.mocked(servicioAcciones.obtenerTodas).mockResolvedValue({
      datos: accionesMock,
      pagina: 1
    })

    const almacen = useAccionesAlmacen()
    await almacen.cargarAcciones()

    expect(almacen.hayAcciones).toBe(true)
  })

  it('limpiarError resetea el campo de error', async () => {
    vi.mocked(servicioAcciones.obtenerTodas).mockRejectedValue(new Error('Error'))

    const almacen = useAccionesAlmacen()
    await almacen.cargarAcciones()
    expect(almacen.error).not.toBeNull()

    almacen.limpiarError()
    expect(almacen.error).toBeNull()
  })

  it('sincronizarDatos muestra mensaje de exito', async () => {
    vi.mocked(servicioAcciones.sincronizar).mockResolvedValue({
      mensaje: 'Sincronizacion completada',
      sincronizados: 10
    })
    vi.mocked(servicioAcciones.obtenerTodas).mockResolvedValue({
      datos: accionesMock,
      pagina: 1
    })

    const almacen = useAccionesAlmacen()
    await almacen.sincronizarDatos()

    expect(almacen.mensajeSincronizacion).toContain('10 acciones actualizadas')
  })
})
