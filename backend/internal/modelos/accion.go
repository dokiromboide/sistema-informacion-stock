package modelos

import "time"

// Accion representa una accion bursatil con sus datos de mercado
type Accion struct {
	ID            int64     `json:"id" db:"id"`
	Simbolo       string    `json:"simbolo" db:"simbolo"`
	Nombre        string    `json:"nombre" db:"nombre"`
	Precio        float64   `json:"precio" db:"precio"`
	Apertura      float64   `json:"apertura" db:"apertura"`
	Maximo        float64   `json:"maximo" db:"maximo"`
	Minimo        float64   `json:"minimo" db:"minimo"`
	Volumen       int64     `json:"volumen" db:"volumen"`
	CambioMonto   float64   `json:"cambio_monto" db:"cambio_monto"`
	CambioPct     float64   `json:"cambio_pct" db:"cambio_pct"`
	UltimaActualizacion time.Time `json:"ultima_actualizacion" db:"ultima_actualizacion"`
}

// PuntajeRecomendacion contiene la accion y su puntaje calculado por el algoritmo
type PuntajeRecomendacion struct {
	Accion  Accion  `json:"accion"`
	Puntaje float64 `json:"puntaje"`
	Razon   string  `json:"razon"`
}

// FiltroAcciones define los parametros de busqueda y ordenamiento
type FiltroAcciones struct {
	Busqueda  string `form:"q"`
	OrdenarPor string `form:"ordenar_por"`
	Direccion string `form:"direccion"`
	Pagina    int    `form:"pagina"`
	PorPagina int    `form:"por_pagina"`
}
