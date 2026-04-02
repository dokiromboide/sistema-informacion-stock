package tests

import (
	"context"
	"testing"

	"github.com/dokiromboide/sistema-informacion-stock/backend/internal/modelos"
	"github.com/dokiromboide/sistema-informacion-stock/backend/internal/servicios"
)

// repoAccionesSimulado implementa un repositorio en memoria para pruebas.
// Permite probar el algoritmo sin necesidad de una base de datos real.
type repoAccionesSimulado struct {
	acciones []modelos.Accion
}

func (r *repoAccionesSimulado) ObtenerTodasParaRecomendacion(_ context.Context) ([]modelos.Accion, error) {
	return r.acciones, nil
}

// adaptarRepoParaServicio convierte el simulado a la interfaz necesaria
func nuevoServicioConDatos(acciones []modelos.Accion) *servicios.ServicioRecomendacionesTesteable {
	return servicios.NuevoServicioRecomendacionesTesteable(acciones)
}

func TestRecomendaciones_AccionPositivaObtieneMayorPuntaje(t *testing.T) {
	acciones := []modelos.Accion{
		{
			Simbolo:     "POSITIVA",
			Precio:      100.0,
			Apertura:    95.0,
			Maximo:      102.0,
			Minimo:      94.0,
			Volumen:     5000000,
			CambioPct:   5.0,
			CambioMonto: 5.0,
		},
		{
			Simbolo:     "NEGATIVA",
			Precio:      100.0,
			Apertura:    105.0,
			Maximo:      106.0,
			Minimo:      99.0,
			Volumen:     1000000,
			CambioPct:   -3.0,
			CambioMonto: -3.0,
		},
	}

	servicio := servicios.NuevoServicioRecomendacionesTesteable(acciones)
	resultados, err := servicio.Recomendar(context.Background())

	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}
	if len(resultados) == 0 {
		t.Fatal("Se esperaba al menos una recomendacion")
	}
	// La accion positiva debe tener el mayor puntaje
	if resultados[0].Accion.Simbolo != "POSITIVA" {
		t.Errorf("Se esperaba POSITIVA como primera recomendacion, se obtuvo: %s", resultados[0].Accion.Simbolo)
	}
}

func TestRecomendaciones_SinAccionesRetornaVacio(t *testing.T) {
	servicio := servicios.NuevoServicioRecomendacionesTesteable([]modelos.Accion{})
	resultados, err := servicio.Recomendar(context.Background())

	if err != nil {
		t.Fatalf("Error inesperado con universo vacio: %v", err)
	}
	if len(resultados) != 0 {
		t.Errorf("Se esperaba slice vacio, se obtuvieron %d recomendaciones", len(resultados))
	}
}

func TestRecomendaciones_RetornaMaximoTopCinco(t *testing.T) {
	acciones := make([]modelos.Accion, 10)
	for i := range acciones {
		acciones[i] = modelos.Accion{
			Simbolo:   "ACC" + string(rune('A'+i)),
			Precio:    float64(100 + i),
			Maximo:    float64(105 + i),
			Minimo:    float64(95 + i),
			Volumen:   int64(1000000 * (i + 1)),
			CambioPct: float64(i),
		}
	}

	servicio := servicios.NuevoServicioRecomendacionesTesteable(acciones)
	resultados, err := servicio.Recomendar(context.Background())

	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}
	if len(resultados) > 5 {
		t.Errorf("Se esperaba maximo 5 recomendaciones, se obtuvieron: %d", len(resultados))
	}
}

func TestRecomendaciones_OrdenadosDeMayorAMenorPuntaje(t *testing.T) {
	acciones := []modelos.Accion{
		{Simbolo: "A1", Precio: 100, Maximo: 101, Minimo: 99, Volumen: 1000000, CambioPct: 1.0},
		{Simbolo: "A2", Precio: 100, Maximo: 103, Minimo: 97, Volumen: 5000000, CambioPct: 4.5},
		{Simbolo: "A3", Precio: 100, Maximo: 102, Minimo: 98, Volumen: 2000000, CambioPct: 2.0},
	}

	servicio := servicios.NuevoServicioRecomendacionesTesteable(acciones)
	resultados, err := servicio.Recomendar(context.Background())

	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}

	for i := 1; i < len(resultados); i++ {
		if resultados[i].Puntaje > resultados[i-1].Puntaje {
			t.Errorf("Los resultados no estan ordenados de mayor a menor puntaje en posicion %d", i)
		}
	}
}
