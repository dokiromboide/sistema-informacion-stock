package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// ConfigurarEnrutador crea y configura el enrutador principal de la API REST.
// Registra todos los endpoints y aplica middlewares de CORS y logging.
func ConfigurarEnrutador(
	handlerAcciones *HandlerAcciones,
	handlerRecomendaciones *HandlerRecomendaciones,
) *gin.Engine {
	enrutador := gin.New()
	enrutador.Use(gin.Logger())
	enrutador.Use(gin.Recovery())

	// CORS: permite peticiones desde el frontend en desarrollo y produccion
	configCORS := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	enrutador.Use(cors.New(configCORS))

	// Endpoint de salud para verificar que el servidor esta activo
	enrutador.GET("/salud", func(c *gin.Context) {
		c.JSON(200, gin.H{"estado": "activo"})
	})

	// Agrupacion de endpoints bajo el prefijo /api
	api := enrutador.Group("/api")
	{
		// Acciones
		api.GET("/acciones", handlerAcciones.ObtenerAcciones)
		api.GET("/acciones/buscar", handlerAcciones.BuscarAcciones)
		api.GET("/acciones/:simbolo", handlerAcciones.ObtenerDetalle)
		api.POST("/sincronizar", handlerAcciones.SincronizarAcciones)

		// Recomendaciones
		api.GET("/recomendaciones", handlerRecomendaciones.ObtenerRecomendaciones)
	}

	return enrutador
}
