import { createRouter, createWebHistory } from 'vue-router'
import VistaAcciones from '@/vistas/VistaAcciones.vue'

const enrutador = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'acciones',
      component: VistaAcciones,
      meta: { titulo: 'Acciones del Mercado' }
    },
    {
      path: '/acciones/:simbolo',
      name: 'detalle-accion',
      component: () => import('@/vistas/VistaDetalleAccion.vue'),
      meta: { titulo: 'Detalle de Accion' }
    },
    {
      path: '/recomendaciones',
      name: 'recomendaciones',
      component: () => import('@/vistas/VistaRecomendaciones.vue'),
      meta: { titulo: 'Recomendaciones de Inversion' }
    }
  ]
})

// Actualizar el titulo de la pagina en cada navegacion
enrutador.beforeEach((to) => {
  const titulo = to.meta.titulo as string
  document.title = titulo ? `${titulo} | Sistema de Stock` : 'Sistema de Stock'
})

export default enrutador
