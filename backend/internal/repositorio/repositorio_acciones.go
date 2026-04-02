package repositorio

import (
	"context"
	"fmt"

	"github.com/dokiromboide/sistema-informacion-stock/backend/internal/modelos"
)

// RepositorioAcciones gestiona la persistencia de acciones bursatiles en CockroachDB
type RepositorioAcciones struct {
	db *BaseDeDatos
}

// NuevoRepositorioAcciones construye un repositorio de acciones
func NuevoRepositorioAcciones(db *BaseDeDatos) *RepositorioAcciones {
	return &RepositorioAcciones{db: db}
}

// CrearTablasSiNoExisten inicializa el esquema de base de datos
func (r *RepositorioAcciones) CrearTablasSiNoExisten(ctx context.Context) error {
	sql := `
	CREATE TABLE IF NOT EXISTS acciones (
		id                    BIGSERIAL PRIMARY KEY,
		simbolo               VARCHAR(20)  NOT NULL UNIQUE,
		nombre                VARCHAR(255) NOT NULL DEFAULT '',
		precio                DECIMAL(18,4) NOT NULL DEFAULT 0,
		apertura              DECIMAL(18,4) NOT NULL DEFAULT 0,
		maximo                DECIMAL(18,4) NOT NULL DEFAULT 0,
		minimo                DECIMAL(18,4) NOT NULL DEFAULT 0,
		volumen               BIGINT        NOT NULL DEFAULT 0,
		cambio_monto          DECIMAL(18,4) NOT NULL DEFAULT 0,
		cambio_pct            DECIMAL(10,4) NOT NULL DEFAULT 0,
		ultima_actualizacion  TIMESTAMPTZ   NOT NULL DEFAULT NOW()
	)`

	_, err := r.db.pool.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("error al crear tabla acciones: %w", err)
	}
	return nil
}

// GuardarOActualizar inserta una nueva accion o actualiza sus datos si el simbolo ya existe
func (r *RepositorioAcciones) GuardarOActualizar(ctx context.Context, accion *modelos.Accion) error {
	sql := `
	INSERT INTO acciones (
		simbolo, nombre, precio, apertura, maximo, minimo,
		volumen, cambio_monto, cambio_pct, ultima_actualizacion
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	ON CONFLICT (simbolo) DO UPDATE SET
		nombre               = EXCLUDED.nombre,
		precio               = EXCLUDED.precio,
		apertura             = EXCLUDED.apertura,
		maximo               = EXCLUDED.maximo,
		minimo               = EXCLUDED.minimo,
		volumen              = EXCLUDED.volumen,
		cambio_monto         = EXCLUDED.cambio_monto,
		cambio_pct           = EXCLUDED.cambio_pct,
		ultima_actualizacion = EXCLUDED.ultima_actualizacion
	RETURNING id`

	err := r.db.pool.QueryRow(ctx, sql,
		accion.Simbolo, accion.Nombre, accion.Precio, accion.Apertura,
		accion.Maximo, accion.Minimo, accion.Volumen, accion.CambioMonto,
		accion.CambioPct, accion.UltimaActualizacion,
	).Scan(&accion.ID)

	if err != nil {
		return fmt.Errorf("error al guardar accion %s: %w", accion.Simbolo, err)
	}
	return nil
}

// ObtenerTodas retorna todas las acciones almacenadas con soporte de busqueda y ordenamiento
func (r *RepositorioAcciones) ObtenerTodas(ctx context.Context, filtro modelos.FiltroAcciones) ([]modelos.Accion, error) {
	columnaOrden := columnaValida(filtro.OrdenarPor)
	direccion := direccionValida(filtro.Direccion)
	pagina, porPagina := paginacionValida(filtro.Pagina, filtro.PorPagina)

	sql := fmt.Sprintf(`
	SELECT id, simbolo, nombre, precio, apertura, maximo, minimo,
	       volumen, cambio_monto, cambio_pct, ultima_actualizacion
	FROM acciones
	WHERE ($1 = '' OR simbolo ILIKE '%%' || $1 || '%%' OR nombre ILIKE '%%' || $1 || '%%')
	ORDER BY %s %s
	LIMIT $2 OFFSET $3`, columnaOrden, direccion)

	desplazamiento := (pagina - 1) * porPagina
	filas, err := r.db.pool.Query(ctx, sql, filtro.Busqueda, porPagina, desplazamiento)
	if err != nil {
		return nil, fmt.Errorf("error al consultar acciones: %w", err)
	}
	defer filas.Close()

	acciones := make([]modelos.Accion, 0)
	for filas.Next() {
		var a modelos.Accion
		if err := filas.Scan(
			&a.ID, &a.Simbolo, &a.Nombre, &a.Precio, &a.Apertura,
			&a.Maximo, &a.Minimo, &a.Volumen, &a.CambioMonto,
			&a.CambioPct, &a.UltimaActualizacion,
		); err != nil {
			return nil, fmt.Errorf("error al leer fila: %w", err)
		}
		acciones = append(acciones, a)
	}

	return acciones, filas.Err()
}

// ObtenerPorSimbolo retorna una accion especifica por su simbolo bursatil
func (r *RepositorioAcciones) ObtenerPorSimbolo(ctx context.Context, simbolo string) (*modelos.Accion, error) {
	sql := `
	SELECT id, simbolo, nombre, precio, apertura, maximo, minimo,
	       volumen, cambio_monto, cambio_pct, ultima_actualizacion
	FROM acciones
	WHERE simbolo = $1`

	var a modelos.Accion
	err := r.db.pool.QueryRow(ctx, sql, simbolo).Scan(
		&a.ID, &a.Simbolo, &a.Nombre, &a.Precio, &a.Apertura,
		&a.Maximo, &a.Minimo, &a.Volumen, &a.CambioMonto,
		&a.CambioPct, &a.UltimaActualizacion,
	)
	if err != nil {
		return nil, fmt.Errorf("accion %s no encontrada: %w", simbolo, err)
	}

	return &a, nil
}

// ObtenerTodasParaRecomendacion retorna todas las acciones sin paginacion para el algoritmo
func (r *RepositorioAcciones) ObtenerTodasParaRecomendacion(ctx context.Context) ([]modelos.Accion, error) {
	sql := `
	SELECT id, simbolo, nombre, precio, apertura, maximo, minimo,
	       volumen, cambio_monto, cambio_pct, ultima_actualizacion
	FROM acciones
	ORDER BY simbolo`

	filas, err := r.db.pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("error al consultar acciones para recomendacion: %w", err)
	}
	defer filas.Close()

	acciones := make([]modelos.Accion, 0)
	for filas.Next() {
		var a modelos.Accion
		if err := filas.Scan(
			&a.ID, &a.Simbolo, &a.Nombre, &a.Precio, &a.Apertura,
			&a.Maximo, &a.Minimo, &a.Volumen, &a.CambioMonto,
			&a.CambioPct, &a.UltimaActualizacion,
		); err != nil {
			return nil, fmt.Errorf("error al leer fila: %w", err)
		}
		acciones = append(acciones, a)
	}

	return acciones, filas.Err()
}

// columnaValida previene inyeccion SQL en la clausula ORDER BY
func columnaValida(columna string) string {
	columnasPermitidas := map[string]bool{
		"simbolo":     true,
		"nombre":      true,
		"precio":      true,
		"cambio_pct":  true,
		"volumen":     true,
	}
	if columnasPermitidas[columna] {
		return columna
	}
	return "simbolo"
}

// direccionValida previene inyeccion SQL en la direccion de ordenamiento
func direccionValida(direccion string) string {
	if direccion == "desc" {
		return "DESC"
	}
	return "ASC"
}

func paginacionValida(pagina, porPagina int) (int, int) {
	if pagina < 1 {
		pagina = 1
	}
	if porPagina < 1 || porPagina > 100 {
		porPagina = 20
	}
	return pagina, porPagina
}
