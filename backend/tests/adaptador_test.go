package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dokiromboide/sistema-informacion-stock/backend/internal/adaptadores"
)

// servidorMock crea un servidor HTTP de prueba que responde con el cuerpo dado
func servidorMock(t *testing.T, cuerpo string, codigo int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(codigo)
		_, _ = w.Write([]byte(cuerpo))
	}))
}

func TestAlphaVantage_ObtenerCotizacion_Exitosa(t *testing.T) {
	respuestaJSON := `{
		"Global Quote": {
			"01. symbol": "AAPL",
			"02. open": "188.00",
			"03. high": "190.50",
			"04. low": "187.20",
			"05. price": "189.50",
			"06. volume": "52000000",
			"07. latest trading day": "2024-01-15",
			"08. previous close": "188.00",
			"09. change": "1.50",
			"10. change percent": "0.7979%"
		}
	}`

	servidor := servidorMock(t, respuestaJSON, http.StatusOK)
	defer servidor.Close()

	adaptador := adaptadores.NuevoAlphaVantage("clave-de-prueba", servidor.URL)
	accion, err := adaptador.ObtenerCotizacion("AAPL")

	if err != nil {
		t.Fatalf("Se esperaba nil como error, se obtuvo: %v", err)
	}
	if accion.Simbolo != "AAPL" {
		t.Errorf("Simbolo esperado: AAPL, obtenido: %s", accion.Simbolo)
	}
	if accion.Precio != 189.50 {
		t.Errorf("Precio esperado: 189.50, obtenido: %f", accion.Precio)
	}
	if accion.CambioPct != 0.7979 {
		t.Errorf("CambioPct esperado: 0.7979, obtenido: %f", accion.CambioPct)
	}
	if accion.Volumen != 52000000 {
		t.Errorf("Volumen esperado: 52000000, obtenido: %d", accion.Volumen)
	}
}

func TestAlphaVantage_ObtenerCotizacion_SimboloVacio(t *testing.T) {
	// Respuesta de limite de API alcanzado
	respuestaJSON := `{"Global Quote": {}}`

	servidor := servidorMock(t, respuestaJSON, http.StatusOK)
	defer servidor.Close()

	adaptador := adaptadores.NuevoAlphaVantage("clave-de-prueba", servidor.URL)
	_, err := adaptador.ObtenerCotizacion("INVALIDO")

	if err == nil {
		t.Fatal("Se esperaba un error para simbolo no encontrado, pero no se obtuvo ninguno")
	}
}

func TestAlphaVantage_ObtenerCotizacion_ErrorHTTP(t *testing.T) {
	servidor := servidorMock(t, "", http.StatusInternalServerError)
	defer servidor.Close()

	adaptador := adaptadores.NuevoAlphaVantage("clave-de-prueba", servidor.URL)
	_, err := adaptador.ObtenerCotizacion("AAPL")

	if err == nil {
		t.Fatal("Se esperaba un error para respuesta HTTP 500")
	}
}

func TestAlphaVantage_ObtenerSimbolosPopulares(t *testing.T) {
	adaptador := adaptadores.NuevoAlphaVantage("clave-de-prueba", "http://no-se-usa")
	simbolos, err := adaptador.ObtenerSimbolosPopulares()

	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}
	if len(simbolos) == 0 {
		t.Error("Se esperaba al menos un simbolo popular")
	}
}

func TestAlphaVantage_NombreProveedor(t *testing.T) {
	adaptador := adaptadores.NuevoAlphaVantage("clave", "http://test")
	if adaptador.NombreProveedor() == "" {
		t.Error("NombreProveedor no debe retornar cadena vacia")
	}
}
