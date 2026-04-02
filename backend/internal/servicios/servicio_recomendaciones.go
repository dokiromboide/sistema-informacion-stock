package servicios

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/dokiromboide/sistema-informacion-stock/backend/internal/modelos"
	"github.com/dokiromboide/sistema-informacion-stock/backend/internal/repositorio"
)

// ServicioRecomendaciones analiza los datos bursatiles almacenados y genera
// una lista priorizada de acciones recomendadas para invertir.
//
// Algoritmo utilizado: puntuacion compuesta de momentum y volumen.
//   - Momentum (60%): mide que tan fuerte y positivo es el movimiento del precio
//   - Volumen relativo (25%): volumen actual vs promedio estimado del sector
//   - Estabilidad (15%): penaliza acciones con alta volatilidad intradiaria
//
// Nota: el algoritmo trabaja unicamente con datos de un dia de mercado.
// Para produccion se recomienda incorporar series de tiempo historicas.
type ServicioRecomendaciones struct {
	repositorio *repositorio.RepositorioAcciones
}

// NuevoServicioRecomendaciones construye el servicio de recomendaciones
func NuevoServicioRecomendaciones(repo *repositorio.RepositorioAcciones) *ServicioRecomendaciones {
	return &ServicioRecomendaciones{repositorio: repo}
}

// Recomendar retorna las mejores acciones ordenadas de mayor a menor puntaje
func (s *ServicioRecomendaciones) Recomendar(ctx context.Context) ([]modelos.PuntajeRecomendacion, error) {
	acciones, err := s.repositorio.ObtenerTodasParaRecomendacion(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al obtener acciones para recomendacion: %w", err)
	}

	if len(acciones) == 0 {
		return []modelos.PuntajeRecomendacion{}, nil
	}

	// Calcular el volumen promedio del universo para la normalizacion
	volumenPromedio := calcularVolumenPromedio(acciones)

	resultados := make([]modelos.PuntajeRecomendacion, 0, len(acciones))
	for _, accion := range acciones {
		puntaje, razon := calcularPuntaje(accion, volumenPromedio)
		resultados = append(resultados, modelos.PuntajeRecomendacion{
			Accion:  accion,
			Puntaje: puntaje,
			Razon:   razon,
		})
	}

	// Ordenar de mayor a menor puntaje
	sort.Slice(resultados, func(i, j int) bool {
		return resultados[i].Puntaje > resultados[j].Puntaje
	})

	// Retornar el top 5 o menos si hay menos acciones
	limite := 5
	if len(resultados) < limite {
		limite = len(resultados)
	}

	return resultados[:limite], nil
}

// calcularPuntaje asigna un puntaje de 0 a 100 a la accion
func calcularPuntaje(accion modelos.Accion, volumenPromedio float64) (float64, string) {
	// --- Componente 1: Momentum del precio (60%) ---
	// Normalizamos el cambio porcentual al rango [0, 100]
	// Consideramos que un cambio de +5% es excelente (100 puntos)
	puntajeMomentum := normalizarMomentum(accion.CambioPct)

	// --- Componente 2: Volumen relativo (25%) ---
	// Una accion con el doble del volumen promedio obtiene el maximo en esta metrica
	puntajeVolumen := normalizarVolumen(float64(accion.Volumen), volumenPromedio)

	// --- Componente 3: Estabilidad intradiaria (15%) ---
	// Acciones con menor spread maximo-minimo relativo son mas estables
	puntajeEstabilidad := calcularEstabilidad(accion.Maximo, accion.Minimo, accion.Precio)

	puntajeFinal := (puntajeMomentum * 0.60) + (puntajeVolumen * 0.25) + (puntajeEstabilidad * 0.15)

	razon := generarRazon(puntajeMomentum, puntajeVolumen, accion.CambioPct)

	return math.Round(puntajeFinal*100) / 100, razon
}

// normalizarMomentum convierte el cambio % al rango [0, 100]
func normalizarMomentum(cambioPct float64) float64 {
	const cambioMaximoEsperado = 5.0
	if cambioPct <= 0 {
		// Acciones en baja reciben puntaje proporcional a la caida
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

// normalizarVolumen compara el volumen de la accion contra el promedio del universo
func normalizarVolumen(volumen, promedio float64) float64 {
	if promedio == 0 {
		return 50
	}
	ratio := volumen / promedio
	// Ratio de 2x o mas = 100 puntos
	if ratio >= 2 {
		return 100
	}
	return ratio * 50
}

// calcularEstabilidad penaliza la volatilidad intradiaria alta
func calcularEstabilidad(maximo, minimo, precio float64) float64 {
	if precio == 0 {
		return 50
	}
	// Spread porcentual maximo-minimo
	spread := ((maximo - minimo) / precio) * 100
	// Spread de 0% = 100 puntos, spread de 10% o mas = 0 puntos
	if spread >= 10 {
		return 0
	}
	return 100 - (spread * 10)
}

// calcularVolumenPromedio retorna el promedio de volumen del universo de acciones
func calcularVolumenPromedio(acciones []modelos.Accion) float64 {
	if len(acciones) == 0 {
		return 0
	}
	total := 0.0
	for _, a := range acciones {
		total += float64(a.Volumen)
	}
	return total / float64(len(acciones))
}

// generarRazon produce un texto explicativo de la recomendacion
func generarRazon(puntajeMomentum, puntajeVolumen, cambioPct float64) string {
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
