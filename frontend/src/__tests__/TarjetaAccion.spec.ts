import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import TarjetaAccion from '@/componentes/TarjetaAccion.vue'
import type { Accion } from '@/tipos'

// Router minimo requerido por RouterLink dentro del componente
const enrutadorPrueba = createRouter({
  history: createMemoryHistory(),
  routes: [
    { path: '/', component: { template: '<div />' } },
    { path: '/acciones/:simbolo', component: { template: '<div />' } }
  ]
})

const accionMock: Accion = {
  id: 1,
  simbolo: 'AAPL',
  nombre: 'Apple Inc.',
  precio: 189.50,
  apertura: 188.00,
  maximo: 190.20,
  minimo: 187.50,
  volumen: 52000000,
  cambio_monto: 1.50,
  cambio_pct: 0.79,
  ultima_actualizacion: '2024-01-15T14:30:00Z'
}

describe('TarjetaAccion', () => {
  it('muestra el simbolo de la accion', async () => {
    const wrapper = mount(TarjetaAccion, {
      props: { accion: accionMock },
      global: { plugins: [enrutadorPrueba] }
    })
    expect(wrapper.text()).toContain('AAPL')
  })

  it('muestra el precio formateado correctamente', async () => {
    const wrapper = mount(TarjetaAccion, {
      props: { accion: accionMock },
      global: { plugins: [enrutadorPrueba] }
    })
    expect(wrapper.text()).toContain('189.50')
  })

  it('muestra el porcentaje de cambio positivo con el signo +', async () => {
    const wrapper = mount(TarjetaAccion, {
      props: { accion: accionMock },
      global: { plugins: [enrutadorPrueba] }
    })
    expect(wrapper.text()).toContain('+0.79%')
  })

  it('aplica clase de color verde para cambio positivo', async () => {
    const wrapper = mount(TarjetaAccion, {
      props: { accion: accionMock },
      global: { plugins: [enrutadorPrueba] }
    })
    const etiqueta = wrapper.find('.bg-green-100')
    expect(etiqueta.exists()).toBe(true)
  })

  it('aplica clase de color rojo para cambio negativo', async () => {
    const accionNegativa: Accion = { ...accionMock, cambio_monto: -2.0, cambio_pct: -1.05 }
    const wrapper = mount(TarjetaAccion, {
      props: { accion: accionNegativa },
      global: { plugins: [enrutadorPrueba] }
    })
    const etiqueta = wrapper.find('.bg-red-100')
    expect(etiqueta.exists()).toBe(true)
  })

  it('formatea el volumen en millones', async () => {
    const wrapper = mount(TarjetaAccion, {
      props: { accion: accionMock },
      global: { plugins: [enrutadorPrueba] }
    })
    expect(wrapper.text()).toContain('52.0M')
  })

  it('muestra "Sin nombre" cuando el nombre esta vacio', async () => {
    const accionSinNombre: Accion = { ...accionMock, nombre: '' }
    const wrapper = mount(TarjetaAccion, {
      props: { accion: accionSinNombre },
      global: { plugins: [enrutadorPrueba] }
    })
    expect(wrapper.text()).toContain('Sin nombre')
  })
})
