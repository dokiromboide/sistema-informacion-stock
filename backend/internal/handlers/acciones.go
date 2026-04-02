// Package handlers contiene los controladores HTTP de la API REST
package handlers

import (
	"net/http"
	"strconv"

	"github.com/dokiromboide/sistema-informacion-stock/backend/internal/modelos"
	"github.com/dokiromboide/sistema-informacion-stock/backend/internal/servicios"
	"github.com/gin-gonic/gin"
)

// HandlerAcciones expone los endpoints HTTP para la gestion de acciones bursatiles
type HandlerAcciones struct {
	servicio *servicios.ServicioStock
}

// NuevoHandlerAcciones construye el handler con su dependencia inyectada
func NuevoHandlerAcciones(servicio *servicios.ServicioStock) *HandlerAcciones {
	return &HandlerAcciones{servicio: servicio}
}

// ObtenerAcciones godoc
// GET /api/acciones
// Parametros de query:
//   - q          : texto de busqueda (nombre o simbolo)
//   - ordenar_por: campo de ordenamiento (simbolo, precio, cambio_pct, volumen)
//   - direccion  : asc | desc
//   - pagina     : numero de pagina (default 1)
//   - por_pagina : items por pagina (default 20, max 100)
func (h *HandlerAcciones) ObtenerAcciones(c *gin.Context) {
	filtro := modelos.FiltroAcciones{
		Busqueda:   c.Query("q"),
		OrdenarPor: c.Query("ordenar_por"),
		Direccion:  c.Query("direccion"),
	}

	pagina, _ := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	porPagina, _ := strconv.Atoi(c.DefaultQuery("por_pagina", "20"))
	filtro.Pagina = pagina
	filtro.PorPagina = porPagina

	acciones, err := h.servicio.ObtenerAcciones(c.Request.Context(), filtro)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al obtener las acciones",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"datos":  acciones,
		"pagina": pagina,
	})
}

// ObtenerDetalle godoc
// GET /api/acciones/:simbolo
func (h *HandlerAcciones) ObtenerDetalle(c *gin.Context) {
	simbolo := c.Param("simbolo")
	if simbolo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El simbolo es requerido"})
		return
	}

	accion, err := h.servicio.ObtenerDetalle(c.Request.Context(), simbolo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Accion no encontrada: " + simbolo,
		})
		return
	}

	c.JSON(http.StatusOK, accion)
}

// BuscarAcciones godoc
// GET /api/acciones/buscar?q=texto
func (h *HandlerAcciones) BuscarAcciones(c *gin.Context) {
	consulta := c.Query("q")
	if consulta == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El parametro q es requerido"})
		return
	}

	resultados, err := h.servicio.BuscarAcciones(c.Request.Context(), consulta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al buscar acciones",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"datos": resultados})
}

// SincronizarAcciones godoc
// POST /api/sincronizar
// Dispara la sincronizacion de datos desde la API externa hacia la base de datos
func (h *HandlerAcciones) SincronizarAcciones(c *gin.Context) {
	sincronizados, err := h.servicio.SincronizarAcciones(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error durante la sincronizacion: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje":       "Sincronizacion completada",
		"sincronizados": sincronizados,
	})
}
