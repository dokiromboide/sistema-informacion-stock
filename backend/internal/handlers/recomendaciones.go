package handlers

import (
	"net/http"

	"github.com/dokiromboide/sistema-informacion-stock/backend/internal/servicios"
	"github.com/gin-gonic/gin"
)

// HandlerRecomendaciones expone el endpoint para obtener recomendaciones de inversion
type HandlerRecomendaciones struct {
	servicio *servicios.ServicioRecomendaciones
}

// NuevoHandlerRecomendaciones construye el handler con su dependencia inyectada
func NuevoHandlerRecomendaciones(servicio *servicios.ServicioRecomendaciones) *HandlerRecomendaciones {
	return &HandlerRecomendaciones{servicio: servicio}
}

// ObtenerRecomendaciones godoc
// GET /api/recomendaciones
// Retorna las mejores acciones para invertir segun el algoritmo de analisis
func (h *HandlerRecomendaciones) ObtenerRecomendaciones(c *gin.Context) {
	recomendaciones, err := h.servicio.Recomendar(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al calcular recomendaciones: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"datos": recomendaciones,
	})
}
