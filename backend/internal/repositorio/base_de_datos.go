// Package repositorio gestiona la persistencia de datos en CockroachDB
package repositorio

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// BaseDeDatos encapsula el pool de conexiones a CockroachDB
type BaseDeDatos struct {
	pool *pgxpool.Pool
}

// NuevaBaseDeDatos crea un pool de conexiones y verifica que la base de datos
// sea accesible antes de retornar
func NuevaBaseDeDatos(ctx context.Context, urlConexion string) (*BaseDeDatos, error) {
	pool, err := pgxpool.New(ctx, urlConexion)
	if err != nil {
		return nil, fmt.Errorf("error al crear el pool de conexiones: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("error al conectar con CockroachDB: %w", err)
	}

	return &BaseDeDatos{pool: pool}, nil
}

// Cerrar cierra todas las conexiones del pool
func (db *BaseDeDatos) Cerrar() {
	db.pool.Close()
}

// Pool expone el pool de conexiones para uso en los repositorios especificos
func (db *BaseDeDatos) Pool() *pgxpool.Pool {
	return db.pool
}
