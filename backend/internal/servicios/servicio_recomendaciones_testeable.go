package servicios

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/dokiromboide/sistema-informacion-stock/backend/internal/modelos"
)

// ServicioRecomendacionesTesteable permite probar el algoritmo sin base de datos.
// Recibe directamente el slice de acciones a analizar.
type ServicioRecomendacionesTesteable struct {
	acciones []modelos.Accion
}

// NuevoServicioRecomendacionesTesteable construye el servicio testeable con datos en memoria
func NuevoServicioRecomendacionesTesteable(acciones []modelos.Accion) *ServicioRecomendacionesTesteable {
	return &ServicioRecomendacionesTesteable{acciones: acciones}
}

// Recomendar aplica el mismo algoritmo que ServicioRecomendaciones sobre los datos en memoria
func (s *ServicioRecomendacionesTesteable) Recomendar(_ context.Context) ([]modelos.PuntajeRecomendacion, error) {
	if len(s.acciones) == 0 {
		return []modelos.PuntajeRecomendacion{}, nil
	}

	volumenPromedio := calcularVolumenPromedio(s.acciones)

	resultados := make([]modelos.PuntajeRecomendacion, 0, len(s.acciones))
	for _, accion := range s.acciones {
		puntaje, razon := calcularPuntaje(accion, volumenPromedio)
		resultados = append(resultados, modelos.PuntajeRecomendacion{
			Accion:  accion,
			Puntaje: puntaje,
			Razon:   razon,
		})
	}

	sort.Slice(resultados, func(i, j int) bool {
		return resultados[i].Puntaje > resultados[j].Puntaje
	})

	limite := 5
	if len(resultados) < limite {
		limite = len(resultados)
	}

	return resultados[:limite], nil
}

// Las funciones siguientes son duplicadas intencionalmente para que el paquete de tests
// no dependa de funciones no exportadas. En un proyecto mas grande se extraerian
// a un subpaquete compartido.

func calcularPuntajeTesteable(accion modelos.Accion, volumenPromedio float64) (float64, string) {
	puntajeMomentum := normalizarMomentumTesteable(accion.CambioPct)
	puntajeVolumen := normalizarVolumenTesteable(float64(accion.Volumen), volumenPromedio)
	puntajeEstabilidad := calcularEstabilidadTesteable(accion.Maximo, accion.Minimo, accion.Precio)
	puntajeFinal := (puntajeMomentum * 0.60) + (puntajeVolumen * 0.25) + (puntajeEstabilidad * 0.15)
	razon := generarRazonTesteable(puntajeMomentum, puntajeVolumen, accion.CambioPct)
	return math.Round(puntajeFinal*100) / 100, razon
}

func normalizarMomentumTesteable(cambioPct float64) float64 {
	const cambioMaximoEsperado = 5.0
	if cambioPct <= 0 {
		normalizado := (cambioPct / cambioMaximoEsperado) * 100
		if normalizado < -100 {
			return 0
		}
		return 50 + (normalizado / 2)
	}
	normalizado := (cambioPct / cambioMaximoEsperado) * 100
	if normalizado > 100 {
		return 100
	}
	return 50 + (normalizado / 2)
}

func normalizarVolumenTesteable(volumen, promedio float64) float64 {
	if promedio == 0 {
		return 50
	}
	ratio := volumen / promedio
	if ratio >= 2 {
		return 100
	}
	return ratio * 50
}

func calcularEstabilidadTesteable(maximo, minimo, precio float64) float64 {
	if precio == 0 {
		return 50
	}
	spread := ((maximo - minimo) / precio) * 100
	if spread >= 10 {
		return 0
	}
	return 100 - (spread * 10)
}

func generarRazonTesteable(puntajeMomentum, puntajeVolumen, cambioPct float64) string {
	switch {
	case puntajeMomentum >= 80 && puntajeVolumen >= 70:
		return fmt.Sprintf("Momentum positivo fuerte (+%.2f%%) con volumen muy elevado", cambioPct)
	case puntajeMomentum >= 80:
		return fmt.Sprintf("Momentum positivo fuerte (+%.2f%%)", cambioPct)
	case puntajeVolumen >= 80:
		return fmt.Sprintf("Volumen inusualmente alto con variacion de %.2f%%", cambioPct)
	case puntajeMomentum >= 60:
		return fmt.Sprintf("Tendencia al alza moderada (+%.2f%%)", cambioPct)
	default:
		return fmt.Sprintf("Movimiento estable con variacion de %.2f%%", cambioPct)
	}
}
