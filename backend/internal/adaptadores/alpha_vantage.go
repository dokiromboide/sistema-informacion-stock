package adaptadores

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dokiromboide/sistema-informacion-stock/backend/internal/modelos"
)

const (
	nombreAlphaVantage = "AlphaVantage"

	// Simbolos populares del mercado estadounidense usados cuando no se solicita uno especifico
	simbolosPopularesCSV = "AAPL,GOOGL,MSFT,AMZN,TSLA,META,NVDA,JPM,V,JNJ"
)

// AlphaVantageAdaptador implementa AdaptadorStock usando la API de Alpha Vantage.
// Documentacion: https://www.alphavantage.co/documentation/
type AlphaVantageAdaptador struct {
	apiKey     string
	urlBase    string
	clienteHTTP *http.Client
}

// NuevoAlphaVantage construye un adaptador de Alpha Vantage listo para usar
func NuevoAlphaVantage(apiKey, urlBase string) *AlphaVantageAdaptador {
	return &AlphaVantageAdaptador{
		apiKey:  apiKey,
		urlBase: urlBase,
		clienteHTTP: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (av *AlphaVantageAdaptador) NombreProveedor() string {
	return nombreAlphaVantage
}

// ObtenerCotizacion obtiene la cotizacion en tiempo real de un simbolo bursatil
func (av *AlphaVantageAdaptador) ObtenerCotizacion(simbolo string) (*modelos.Accion, error) {
	url := fmt.Sprintf(
		"%s?function=GLOBAL_QUOTE&symbol=%s&apikey=%s",
		av.urlBase, simbolo, av.apiKey,
	)

	respuesta, err := av.clienteHTTP.Get(url)
	if err != nil {
		return nil, &ErrorAPI{
			Proveedor: nombreAlphaVantage,
			Operacion: "ObtenerCotizacion",
			Mensaje:   "error al realizar la solicitud HTTP: " + err.Error(),
		}
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode != http.StatusOK {
		return nil, &ErrorAPI{
			Proveedor:  nombreAlphaVantage,
			Operacion:  "ObtenerCotizacion",
			Mensaje:    "respuesta inesperada del servidor",
			CodigoHTTP: respuesta.StatusCode,
		}
	}

	cuerpo, err := io.ReadAll(respuesta.Body)
	if err != nil {
		return nil, &ErrorAPI{
			Proveedor: nombreAlphaVantage,
			Operacion: "ObtenerCotizacion",
			Mensaje:   "error al leer el cuerpo de la respuesta: " + err.Error(),
		}
	}

	return av.parsearCotizacionGlobal(cuerpo)
}

// respuestaCotizacionGlobal mapea la respuesta JSON de Alpha Vantage para GLOBAL_QUOTE
type respuestaCotizacionGlobal struct {
	GlobalQuote struct {
		Simbolo          string `json:"01. symbol"`
		Apertura         string `json:"02. open"`
		Maximo           string `json:"03. high"`
		Minimo           string `json:"04. low"`
		Precio           string `json:"05. price"`
		Volumen          string `json:"06. volume"`
		FechaAnterior    string `json:"07. latest trading day"`
		CierreAnterior   string `json:"08. previous close"`
		Cambio           string `json:"09. change"`
		CambioPct        string `json:"10. change percent"`
	} `json:"Global Quote"`
}

func (av *AlphaVantageAdaptador) parsearCotizacionGlobal(datos []byte) (*modelos.Accion, error) {
	var respuesta respuestaCotizacionGlobal
	if err := json.Unmarshal(datos, &respuesta); err != nil {
		return nil, &ErrorAPI{
			Proveedor: nombreAlphaVantage,
			Operacion: "parsearCotizacionGlobal",
			Mensaje:   "error al decodificar JSON: " + err.Error(),
		}
	}

	q := respuesta.GlobalQuote
	if q.Simbolo == "" {
		return nil, &ErrorAPI{
			Proveedor: nombreAlphaVantage,
			Operacion: "parsearCotizacionGlobal",
			Mensaje:   "simbolo no encontrado o limite de API alcanzado",
		}
	}

	precio, _ := strconv.ParseFloat(q.Precio, 64)
	apertura, _ := strconv.ParseFloat(q.Apertura, 64)
	maximo, _ := strconv.ParseFloat(q.Maximo, 64)
	minimo, _ := strconv.ParseFloat(q.Minimo, 64)
	volumen, _ := strconv.ParseInt(q.Volumen, 10, 64)
	cambioMonto, _ := strconv.ParseFloat(q.Cambio, 64)

	// El porcentaje viene como "1.23%" — eliminar el simbolo antes de parsear
	pctStr := strings.TrimSuffix(q.CambioPct, "%")
	pctStr = strings.TrimSpace(pctStr)
	cambioPct, _ := strconv.ParseFloat(pctStr, 64)

	return &modelos.Accion{
		Simbolo:             q.Simbolo,
		Precio:              precio,
		Apertura:            apertura,
		Maximo:              maximo,
		Minimo:              minimo,
		Volumen:             volumen,
		CambioMonto:         cambioMonto,
		CambioPct:           cambioPct,
		UltimaActualizacion: time.Now(),
	}, nil
}

// ObtenerSimbolosPopulares retorna una lista predefinida de simbolos del mercado
func (av *AlphaVantageAdaptador) ObtenerSimbolosPopulares() ([]string, error) {
	return strings.Split(simbolosPopularesCSV, ","), nil
}

// BuscarSimbolo busca instrumentos usando el endpoint SYMBOL_SEARCH de Alpha Vantage
func (av *AlphaVantageAdaptador) BuscarSimbolo(consulta string) ([]modelos.Accion, error) {
	url := fmt.Sprintf(
		"%s?function=SYMBOL_SEARCH&keywords=%s&apikey=%s",
		av.urlBase, consulta, av.apiKey,
	)

	respuesta, err := av.clienteHTTP.Get(url)
	if err != nil {
		return nil, &ErrorAPI{
			Proveedor: nombreAlphaVantage,
			Operacion: "BuscarSimbolo",
			Mensaje:   "error al realizar la solicitud: " + err.Error(),
		}
	}
	defer respuesta.Body.Close()

	cuerpo, err := io.ReadAll(respuesta.Body)
	if err != nil {
		return nil, &ErrorAPI{
			Proveedor: nombreAlphaVantage,
			Operacion: "BuscarSimbolo",
			Mensaje:   "error al leer respuesta: " + err.Error(),
		}
	}

	return av.parsearResultadosBusqueda(cuerpo)
}

type respuestaBusqueda struct {
	BestMatches []struct {
		Simbolo string `json:"1. symbol"`
		Nombre  string `json:"2. name"`
	} `json:"bestMatches"`
}

func (av *AlphaVantageAdaptador) parsearResultadosBusqueda(datos []byte) ([]modelos.Accion, error) {
	var respuesta respuestaBusqueda
	if err := json.Unmarshal(datos, &respuesta); err != nil {
		return nil, &ErrorAPI{
			Proveedor: nombreAlphaVantage,
			Operacion: "parsearResultadosBusqueda",
			Mensaje:   "error al decodificar JSON: " + err.Error(),
		}
	}

	resultados := make([]modelos.Accion, 0, len(respuesta.BestMatches))
	for _, match := range respuesta.BestMatches {
		resultados = append(resultados, modelos.Accion{
			Simbolo: match.Simbolo,
			Nombre:  match.Nombre,
		})
	}

	return resultados, nil
}
