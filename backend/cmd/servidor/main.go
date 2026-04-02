package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dokiromboide/sistema-informacion-stock/backend/internal/adaptadores"
	"github.com/dokiromboide/sistema-informacion-stock/backend/internal/handlers"
	"github.com/dokiromboide/sistema-informacion-stock/backend/internal/repositorio"
	"github.com/dokiromboide/sistema-informacion-stock/backend/internal/servicios"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno desde el archivo .env
	if err := godotenv.Load(); err != nil {
		log.Println("Archivo .env no encontrado, usando variables de entorno del sistema")
	}

	// Conectar a CockroachDB
	ctx := context.Background()
	urlDB := os.Getenv("DB_URL")
	if urlDB == "" {
		log.Fatal("DB_URL no esta definida en las variables de entorno")
	}

	bd, err := repositorio.NuevaBaseDeDatos(ctx, urlDB)
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}
	defer bd.Cerrar()

	// Inicializar esquema de la base de datos
	repoAcciones := repositorio.NuevoRepositorioAcciones(bd)
	if err := repoAcciones.CrearTablasSiNoExisten(ctx); err != nil {
		log.Fatalf("Error al inicializar el esquema: %v", err)
	}

	// Crear el adaptador de API segun la configuracion
	adaptador, err := adaptadores.NuevoAdaptador()
	if err != nil {
		log.Fatalf("No se pudo inicializar el adaptador de API: %v", err)
	}
	log.Printf("Proveedor de API: %s", adaptador.NombreProveedor())

	// Inicializar servicios
	servicioStock := servicios.NuevoServicioStock(adaptador, repoAcciones)
	servicioRecomendaciones := servicios.NuevoServicioRecomendaciones(repoAcciones)

	// Inicializar handlers y enrutador
	handlerAcciones := handlers.NuevoHandlerAcciones(servicioStock)
	handlerRecomendaciones := handlers.NuevoHandlerRecomendaciones(servicioRecomendaciones)
	enrutador := handlers.ConfigurarEnrutador(handlerAcciones, handlerRecomendaciones)

	// Configurar y arrancar el servidor HTTP
	puerto := os.Getenv("PUERTO")
	if puerto == "" {
		puerto = "8080"
	}

	servidor := &http.Server{
		Addr:         ":" + puerto,
		Handler:      enrutador,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Arrancar en goroutine para poder hacer shutdown graceful
	go func() {
		log.Printf("Servidor iniciado en http://localhost:%s", puerto)
		if err := servidor.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error al iniciar el servidor: %v", err)
		}
	}()

	// Esperar senal de terminacion del sistema operativo
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Cerrando el servidor...")
	ctxShutdown, cancelar := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelar()

	if err := servidor.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Error al cerrar el servidor: %v", err)
	}

	log.Println("Servidor cerrado correctamente")
}
