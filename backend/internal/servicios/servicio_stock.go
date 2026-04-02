// Package servicios contiene la logica de negocio del sistema
package servicios

import (
	"context"
	"fmt"
	"log"

	"github.com/dokiromboide/sistema-informacion-stock/backend/internal/adaptadores"
	"github.com/dokiromboide/sistema-informacion-stock/backend/internal/modelos"
	"github.com/dokiromboide/sistema-informacion-stock/backend/internal/repositorio"
)

// ServicioStock coordina la obtencion de datos desde la API externa y su persistencia
type ServicioStock struct {
	adaptador   adaptadores.AdaptadorStock
	repositorio *repositorio.RepositorioAcciones
}

// NuevoServicioStock construye el servicio con sus dependencias inyectadas
func NuevoServicioStock(adaptador adaptadores.AdaptadorStock, repo *repositorio.RepositorioAcciones) *ServicioStock {
	return &ServicioStock{
		adaptador:   adaptador,
		repositorio: repo,
	}
}

// SincronizarAcciones descarga los simbolos populares desde la API y los persiste en la base de datos
func (s *ServicioStock) SincronizarAcciones(ctx context.Context) (int, error) {
	simbolos, err := s.adaptador.ObtenerSimbolosPopulares()
	if err != nil {
		return 0, fmt.Errorf("error al obtener simbolos populares: %w", err)
	}

	sincronizados := 0
	for _, simbolo := range simbolos {
		accion, err := s.adaptador.ObtenerCotizacion(simbolo)
		if err != nil {
			log.Printf("Advertencia: no se pudo obtener cotizacion de %s: %v", simbolo, err)
			continue
		}

		if err := s.repositorio.GuardarOActualizar(ctx, accion); err != nil {
			log.Printf("Advertencia: no se pudo guardar %s: %v", simbolo, err)
			continue
		}

		sincronizados++
		log.Printf("Sincronizado: %s = $%.2f (%.2f%%)", accion.Simbolo, accion.Precio, accion.CambioPct)
	}

	return sincronizados, nil
}

// ObtenerAcciones retorna acciones con filtros aplicados
func (s *ServicioStock) ObtenerAcciones(ctx context.Context, filtro modelos.FiltroAcciones) ([]modelos.Accion, error) {
	return s.repositorio.ObtenerTodas(ctx, filtro)
}

// ObtenerDetalle retorna el detalle de una accion, actualizando desde la API si es posible
func (s *ServicioStock) ObtenerDetalle(ctx context.Context, simbolo string) (*modelos.Accion, error) {
	// Intentar actualizar desde la API antes de retornar
	accionAPI, err := s.adaptador.ObtenerCotizacion(simbolo)
	if err == nil {
		// Si la API responde, guardar los datos frescos y retornarlos
		if guardadoErr := s.repositorio.GuardarOActualizar(ctx, accionAPI); guardadoErr != nil {
			log.Printf("No se pudo persistir actualizacion de %s: %v", simbolo, guardadoErr)
		}
		return accionAPI, nil
	}

	// Si la API falla, intentar retornar desde la base de datos
	log.Printf("API no disponible para %s, usando datos en cache: %v", simbolo, err)
	return s.repositorio.ObtenerPorSimbolo(ctx, simbolo)
}

// BuscarAcciones busca simbolos por nombre o ticker
func (s *ServicioStock) BuscarAcciones(ctx context.Context, consulta string) ([]modelos.Accion, error) {
	return s.adaptador.BuscarSimbolo(consulta)
}
