package adaptadores

import (
	"fmt"
	"os"
)

// NuevoAdaptador actua como fabrica y retorna el adaptador correcto segun
// la variable de entorno PROVEEDOR_API. Por defecto usa Alpha Vantage.
//
// Para agregar un nuevo proveedor:
//  1. Crear un archivo nuevo en este paquete implementando AdaptadorStock
//  2. Agregar un caso en este switch con el identificador del proveedor
//  3. Establecer PROVEEDOR_API en el .env con el nuevo identificador
func NuevoAdaptador() (AdaptadorStock, error) {
	proveedor := os.Getenv("PROVEEDOR_API")
	if proveedor == "" {
		proveedor = "alpha_vantage"
	}

	switch proveedor {
	case "alpha_vantage":
		apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("ALPHA_VANTAGE_API_KEY no esta definida en las variables de entorno")
		}
		urlBase := os.Getenv("ALPHA_VANTAGE_URL")
		if urlBase == "" {
			urlBase = "https://www.alphavantage.co/query"
		}
		return NuevoAlphaVantage(apiKey, urlBase), nil

	default:
		return nil, fmt.Errorf("proveedor de API desconocido: %q. Proveedores soportados: alpha_vantage", proveedor)
	}
}
