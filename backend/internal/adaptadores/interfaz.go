// Package adaptadores implementa el patron adaptador para conectarse a diferentes
// proveedores de datos bursatiles (Alpha Vantage, Polygon.io, Yahoo Finance, etc.)
// sin modificar la logica de negocio del sistema.
package adaptadores

import "github.com/dokiromboide/sistema-informacion-stock/backend/internal/modelos"

// AdaptadorStock define el contrato que deben cumplir todos los proveedores de datos
// bursatiles. Al depender de esta interfaz, el resto del sistema permanece
// desacoplado del proveedor concreto utilizado.
type AdaptadorStock interface {
	// ObtenerCotizacion retorna la cotizacion actual de un simbolo bursatil
	ObtenerCotizacion(simbolo string) (*modelos.Accion, error)

	// ObtenerSimbolosPopulares retorna una lista de simbolos populares del mercado
	ObtenerSimbolosPopulares() ([]string, error)

	// BuscarSimbolo busca simbolos por nombre o ticker
	BuscarSimbolo(consulta string) ([]modelos.Accion, error)

	// NombreProveedor retorna el nombre del proveedor para logging y diagnostico
	NombreProveedor() string
}

// ErrorAPI representa un error al comunicarse con la API externa
type ErrorAPI struct {
	Proveedor  string
	Operacion  string
	Mensaje    string
	CodigoHTTP int
}

func (e *ErrorAPI) Error() string {
	return "[" + e.Proveedor + "] " + e.Operacion + ": " + e.Mensaje
}
