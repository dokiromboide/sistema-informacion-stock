# Arquitectura del Sistema de Informacion de Stock

Este documento describe como esta construido el sistema, que tecnologias se eligieron y por que, como fluyen los datos entre sus componentes, y los patrones de diseno implementados.

---

## Vision general

El sistema esta compuesto por tres capas independientes que se comunican entre si:

```
+------------------------------------------+
|          USUARIO EN EL NAVEGADOR         |
+------------------------------------------+
                     |
              HTTP / WebSocket
                     |
                     v
+------------------------------------------+
|    FRONTEND  (Vue 3 - Puerto 5173)       |
|                                          |
|  VistaAcciones   VistaDetalle   Recom.   |
|       |               |           |      |
|       +-------+-------+-----------+      |
|               |                          |
|           Pinia Store                    |
|               |                          |
|          Servicio API (Axios)            |
+------------------------------------------+
                     |
              HTTP / JSON (REST)
                     |
                     v
+------------------------------------------+
|    BACKEND   (Go / Gin - Puerto 8080)    |
|                                          |
|   Handlers  ->  Servicios  ->  Repos     |
|                    |              |      |
|               Adaptadores    CockroachDB |
|                    |        (Puerto 26257)|
|             API Alpha Vantage            |
|               (Internet)                 |
+------------------------------------------+
```

Cada capa tiene una responsabilidad clara:

- **Frontend:** Presentar la informacion al usuario e interactuar con sus acciones
- **Backend:** Procesar la logica de negocio, exponer la API REST y coordinar con la base de datos y las APIs externas
- **Base de datos:** Persistir los datos de acciones para evitar consultar la API externa en cada solicitud

---

## Tecnologias y por que se eligieron

### Backend: Go con Gin

**Go** es un lenguaje compilado con tipado estatico, orientado a la simplicidad y al rendimiento. Se eligio por:

- Alto rendimiento en servidores HTTP con bajo consumo de recursos
- Manejo de errores explicito que reduce comportamientos inesperados
- Soporte nativo de interfaces, lo que facilita el patron adaptador
- Compilacion rapida y binarios pequenos sin dependencias externas

**Gin** es el framework HTTP mas usado en el ecosistema Go. Provee enrutamiento eficiente, middleware de CORS, serializacion JSON automatica y recuperacion ante panics.

### Frontend: Vue 3 con TypeScript

**Vue 3** con la Composition API permite escribir componentes con logica clara, reutilizable y bien encapsulada. La API de reactividad de Vue actualiza la interfaz automaticamente cuando cambia el estado, sin manipulacion manual del DOM.

**TypeScript** agrega tipado estatico al JavaScript del frontend. Detecta errores en tiempo de desarrollo, antes de que lleguen al navegador del usuario. También mejora el autocompletado en el editor.

**Pinia** es la libreria oficial de gestion de estado para Vue 3. Centraliza el estado de la aplicacion (lista de acciones, filtros, estado de carga) en un unico lugar accesible desde cualquier componente.

**Vite** es el servidor de desarrollo y compilador del frontend. Inicia en menos de un segundo y aplica cambios en caliente sin recargar la pagina completa.

**Tailwind CSS** permite aplicar estilos directamente en el HTML mediante clases utilitarias. Elimina la necesidad de escribir CSS personalizado para la mayoria de los casos.

### Base de datos: CockroachDB con Docker

**CockroachDB** es una base de datos compatible con PostgreSQL al 100%. Se eligio porque:

- El mismo driver `pgx` de Go que se usa con PostgreSQL funciona sin cambios
- Se distribuye como imagen Docker oficial, sin instalacion manual
- Incluye una interfaz web de administracion en el puerto 8081
- Soporta SQL estandar, transacciones ACID y operaciones de `INSERT ... ON CONFLICT UPDATE`

**Docker** permite ejecutar CockroachDB como contenedor aislado. El archivo `docker-compose.yml` define toda la configuracion necesaria, de modo que un solo comando levanta la base de datos con el esquema correcto.

---

## Patron adaptador para proveedores de datos

### El problema

Existen multiples proveedores de datos bursatiles: Alpha Vantage, Polygon.io, Yahoo Finance, Finnhub, IEX Cloud, entre otros. Cada uno tiene una URL diferente, parametros de autenticacion distintos y un formato de respuesta propio.

Si el codigo del sistema dependiera directamente de Alpha Vantage, cambiar de proveedor implicaria modificar los servicios, los modelos de datos, los tests y posiblemente los handlers. Ademas, seria imposible testear la logica de negocio sin hacer llamadas reales a internet.

### La solucion: interfaz + adaptador + fabrica

El sistema define una **interfaz** en Go que actua como contrato comun para todos los proveedores:

```
AdaptadorStock (interfaz)
    |
    +-- AlphaVantageAdaptador    (implementacion actual)
    |
    +-- PolygonAdaptador         (implementacion futura)
    |
    +-- YahooFinanceAdaptador    (implementacion futura)
```

**Archivo:** `backend/internal/adaptadores/interfaz.go`

La interfaz define exactamente cuatro operaciones:

| Metodo | Parametros | Retorna | Descripcion |
|--------|-----------|---------|-------------|
| `ObtenerCotizacion` | `simbolo string` | `Accion, error` | Precio actual y datos del dia |
| `ObtenerSimbolosPopulares` | ninguno | `[]string, error` | Lista de simbolos por defecto |
| `BuscarSimbolo` | `consulta string` | `[]Accion, error` | Busqueda por nombre o ticker |
| `NombreProveedor` | ninguno | `string` | Nombre identificador del proveedor |

El resto del sistema (servicios, handlers, base de datos) solo conoce la interfaz. Nunca importa ni menciona `AlphaVantageAdaptador` directamente.

### Fabrica de adaptadores

La **fabrica** (`backend/internal/adaptadores/fabrica.go`) es el unico lugar donde se decide que proveedor usar. Lee la variable de entorno `PROVEEDOR_API` y construye el adaptador correspondiente:

```
Variable de entorno: PROVEEDOR_API=alpha_vantage
                            |
                            v
                    FabricaAdaptador()
                            |
                            v
              AlphaVantageAdaptador{apiKey, url}
                            |
                            v
              ServicioStock (recibe AdaptadorStock)
                            |
                      usa los 4 metodos
                    (sin conocer Alpha Vantage)
```

Cambiar de proveedor solo requiere:

1. Crear un archivo nuevo que implemente los 4 metodos de `AdaptadorStock`
2. Agregar un `case` en la fabrica
3. Cambiar la variable `PROVEEDOR_API` en el archivo `.env`

---

## Flujo de datos

### Flujo 1: Sincronizacion de datos

Cuando el usuario hace clic en "Sincronizar":

```
1.  Usuario hace clic en "Sincronizar" (navegador)
2.  Frontend envia: POST /api/sincronizar
3.  HandlerAcciones.Sincronizar() recibe la peticion HTTP
4.  Llama a ServicioStock.SincronizarAcciones()
5.  ServicioStock solicita al Adaptador la lista de simbolos populares
6.  Adaptador consulta Alpha Vantage API (internet)
7.  Para cada simbolo, Adaptador solicita la cotizacion actual a Alpha Vantage
8.  ServicioStock recibe cada cotizacion y llama a Repositorio.GuardarOActualizar()
9.  Repositorio ejecuta INSERT ... ON CONFLICT UPDATE en CockroachDB
10. El ciclo se repite para cada uno de los 10 simbolos configurados
11. HandlerAcciones responde: {"sincronizados": 10}
12. Frontend muestra notificacion de exito y recarga la lista
```

### Flujo 2: Listado de acciones con filtros

```
1.  Usuario escribe en el campo de busqueda o cambia el orden
2.  Frontend envia: GET /api/acciones?q=apple&ordenar_por=precio&direccion=desc
3.  HandlerAcciones.ListarAcciones() extrae los parametros de la URL
4.  Llama a ServicioStock.ObtenerAcciones(filtro)
5.  ServicioStock llama a Repositorio.ObtenerTodas(filtro)
6.  Repositorio ejecuta SELECT ... WHERE simbolo ILIKE $1 ORDER BY precio DESC
7.  CockroachDB devuelve las filas que coinciden con el filtro
8.  HandlerAcciones serializa el resultado a JSON y responde
9.  Frontend actualiza la grilla de tarjetas con los nuevos datos
```

### Flujo 3: Detalle de una accion

```
1.  Usuario hace clic en una tarjeta de accion
2.  Frontend navega a la ruta /acciones/AAPL
3.  Frontend envia: GET /api/acciones/AAPL
4.  HandlerAcciones.ObtenerDetalle() recibe el simbolo de la URL
5.  ServicioStock intenta obtener un dato fresco desde el Adaptador (Alpha Vantage)
6a. Si Alpha Vantage responde correctamente:
    - Guarda el dato actualizado en CockroachDB
    - Retorna el dato fresco al frontend
6b. Si Alpha Vantage falla (limite de API diario alcanzado, sin conexion):
    - Retorna el dato almacenado en CockroachDB sin error al usuario
7.  Frontend muestra la pantalla de detalle con los datos recibidos
```

### Flujo 4: Algoritmo de recomendaciones

```
1.  Usuario navega a /recomendaciones
2.  Frontend envia: GET /api/recomendaciones
3.  HandlerRecomendaciones.ObtenerRecomendaciones() recibe la peticion
4.  Llama a ServicioRecomendaciones.Recomendar()
5.  ServicioRecomendaciones obtiene todas las acciones de CockroachDB
6.  Para cada accion calcula un puntaje compuesto de tres factores:
    a. Momentum del precio (60% del puntaje):
       puntaje_momentum = cambio_porcentual_del_dia (normalizado entre 0 y 1)
    b. Volumen relativo (25% del puntaje):
       puntaje_volumen = volumen_accion / promedio_volumen_universo
    c. Estabilidad del precio (15% del puntaje):
       puntaje_estabilidad = 1 - ((maximo - minimo) / precio_apertura)
    puntaje_total = 0.60 * momentum + 0.25 * volumen + 0.15 * estabilidad
7.  Ordena todas las acciones de mayor a menor puntaje
8.  Toma las primeras 5 y construye una descripcion en texto para cada una
9.  Responde con el ranking de 5 acciones con puntaje y razon
10. Frontend muestra el ranking con posiciones y colores
```

---

## Estructura de capas del backend

```
+----------------------------------------------+
|              CAPA HTTP (handlers/)            |
|  Recibe peticiones HTTP, valida parametros,   |
|  serializa/deserializa JSON, responde         |
|  No conoce la base de datos ni la API externa |
+----------------------------------------------+
                      |
                      v
+----------------------------------------------+
|           CAPA DE SERVICIOS (servicios/)      |
|  Logica de negocio: cuando sincronizar,       |
|  como calcular puntajes, que datos retornar   |
|  Coordina entre el adaptador y el repositorio |
|  No conoce HTTP ni SQL directamente           |
+----------------------------------------------+
            |                    |
            v                    v
+------------------+  +----------------------+
|   ADAPTADORES    |  |     REPOSITORIO       |
|  (adaptadores/)  |  |   (repositorio/)      |
|                  |  |                       |
|  Hablan con APIs |  |  Hablan con           |
|  externas usando |  |  CockroachDB usando   |
|  HTTP/JSON       |  |  SQL parametrizado    |
+------------------+  +----------------------+
        |                       |
        v                       v
  Alpha Vantage            CockroachDB
    (internet)            (contenedor Docker)
```

Esta separacion en capas tiene ventajas practicas:

- Los tests del algoritmo de recomendaciones no necesitan internet ni base de datos: usan un repositorio en memoria
- Los tests del adaptador no necesitan base de datos: usan un servidor HTTP de prueba que simula Alpha Vantage
- Cambiar la base de datos requiere solo reescribir el repositorio, sin tocar los servicios ni los handlers

---

## Estructura de capas del frontend

```
+----------------------------------------------+
|               VISTAS (vistas/)               |
|  Paginas completas de la aplicacion          |
|  VistaAcciones, VistaDetalle, Recomendaciones|
|  Leen del almacen, componen componentes      |
+----------------------------------------------+
                      |
                      v
+----------------------------------------------+
|            COMPONENTES (componentes/)         |
|  Piezas visuales reutilizables               |
|  TarjetaAccion, BarraBusqueda, SelectorOrden |
|  Reciben datos via props, emiten eventos     |
+----------------------------------------------+
                      |
                      v
+----------------------------------------------+
|           ALMACEN PINIA (almacen/)            |
|  Estado global de la aplicacion              |
|  acciones, accionActual, cargando, error     |
|  Llama al servicio de API, guarda resultados |
+----------------------------------------------+
                      |
                      v
+----------------------------------------------+
|          SERVICIO DE API (servicios/)         |
|  Encapsula todas las llamadas HTTP           |
|  Usa Axios configurado con la URL del backend |
|  Un unico lugar donde se define el protocolo |
+----------------------------------------------+
                      |
                   HTTP/JSON
                      |
                      v
             Backend (Puerto 8080)
```

---

## Modelo de datos

La entidad principal del sistema es `Accion`, que representa una cotizacion bursatil en un momento dado.

**Estructura en Go** (`backend/internal/modelos/accion.go`):

```
Accion {
    Simbolo          string    // Ticker: "AAPL", "GOOGL", etc.
    Nombre           string    // Nombre completo: "Apple Inc."
    Precio           float64   // Precio actual de la accion
    PrecioApertura   float64   // Precio al abrir el mercado ese dia
    PrecioMaximo     float64   // Precio mas alto del dia
    PrecioMinimo     float64   // Precio mas bajo del dia
    Volumen          int64     // Cantidad de acciones operadas en el dia
    Cambio           float64   // Variacion en dolares respecto al dia anterior
    CambioPorcentual float64   // Variacion en porcentaje respecto al dia anterior
    UltimaActual.   time.Time  // Timestamp de la ultima sincronizacion
}
```

**Tabla en CockroachDB:**

```sql
CREATE TABLE IF NOT EXISTS acciones (
    simbolo            TEXT PRIMARY KEY,
    nombre             TEXT,
    precio             DECIMAL,
    precio_apertura    DECIMAL,
    precio_maximo      DECIMAL,
    precio_minimo      DECIMAL,
    volumen            BIGINT,
    cambio             DECIMAL,
    cambio_porcentual  DECIMAL,
    ultima_actualizacion TIMESTAMPTZ DEFAULT NOW()
);
```

---

## Seguridad implementada

### Prevencion de inyeccion SQL

El repositorio nunca concatena strings para construir consultas SQL. Todas las consultas usan parametros posicionales (`$1`, `$2`, ...) que el driver `pgx` maneja de forma segura:

```
// Correcto: parametro posicional
SELECT * FROM acciones WHERE simbolo = $1

// Incorrecto (vulnerable a inyeccion SQL): NO se usa en este sistema
SELECT * FROM acciones WHERE simbolo = '" + input + "'"
```

Para el ordenamiento dinamico, donde SQL no permite parametrizar nombres de columnas, el repositorio usa una lista blanca de columnas validas en la funcion `columnaValida()`. Cualquier valor que no este en la lista es rechazado y reemplazado por el valor por defecto.

### CORS configurado

El servidor solo acepta peticiones de los origenes configurados (`localhost:5173` y `localhost:3000`). Peticiones de otros origenes son rechazadas por el navegador antes de llegar al servidor.

### Variables sensibles fuera del codigo fuente

La clave de API y la cadena de conexion a la base de datos se leen desde variables de entorno en tiempo de ejecucion. Nunca estan escritas en el codigo fuente. El archivo `.env` esta incluido en `.gitignore` para evitar que se suba accidentalmente a GitHub.

---

## Pruebas unitarias

### Backend (Go testing)

Las pruebas estan en `backend/tests/`:

**Adaptador Alpha Vantage** (`adaptador_test.go`):

Se crea un servidor HTTP de prueba con `httptest.NewServer` que simula las respuestas de Alpha Vantage sin hacer llamadas reales a internet. Esto permite ejecutar los tests sin conexion y sin consumir el limite diario de solicitudes. Se cubren tres escenarios:

- Respuesta exitosa con datos validos
- Simbolo no encontrado (respuesta vacia de la API)
- Error del servidor HTTP 500

**Algoritmo de recomendaciones** (`recomendaciones_test.go`):

Se usa `ServicioRecomendacionesTesteable`, una variante del servicio que acepta un repositorio en memoria en lugar del repositorio real. Esto permite testear el algoritmo sin base de datos. Se verifican:

- El orden correcto del ranking (la accion con mayor puntaje va primera)
- Que el resultado se limite a 5 acciones como maximo
- Que con un universo de acciones vacio se retorne una lista vacia sin error

### Frontend (Vitest)

Las pruebas estan en `frontend/src/__tests__/`:

**TarjetaAccion** (`TarjetaAccion.spec.ts`):

Se monta el componente con datos simulados usando `@vue/test-utils` y se verifica:

- Que el simbolo y el nombre de la empresa se rendericen correctamente
- Que el precio muestre el formato correcto
- Que el color sea verde para variacion positiva y rojo para negativa
- Que el volumen se formatee con separadores de miles

**accionesAlmacen** (`accionesAlmacen.spec.ts`):

Se simula el modulo de API con `vi.mock` para no hacer llamadas HTTP reales. Se prueban:

- Estado inicial del almacen (listas vacias, sin errores, sin carga)
- Carga exitosa de acciones (el almacen se actualiza correctamente)
- Manejo de errores de red (el error queda almacenado en el estado)
- Proceso de sincronizacion (el flag `sincronizando` cambia correctamente)

---

## Diagrama de componentes del frontend

```
App.vue
└── RouterView
    ├── VistaAcciones.vue
    │   ├── BarraBusqueda.vue
    │   ├── SelectorOrden.vue
    │   └── TarjetaAccion.vue (x N acciones)
    ├── VistaDetalleAccion.vue
    └── VistaRecomendaciones.vue
```

Todos los componentes leen del `accionesAlmacen` (Pinia). Cuando el almacen cambia, Vue actualiza automaticamente solo los componentes afectados sin recargar la pagina.

---

## Decisiones de diseno

### Por que CockroachDB en lugar de SQLite o PostgreSQL directo

- **SQLite** no tiene servidor propio, lo que dificulta la migracion a produccion sin cambiar el codigo
- **PostgreSQL** requiere instalacion manual y configuracion adicional para desarrollo local
- **CockroachDB** con Docker ofrece lo mejor de ambos: compatibilidad total con PostgreSQL, sin instalacion manual, y levanta con un solo comando

### Por que Pinia en lugar de Vuex

Vuex era la libreria de estado estandar en Vue 2 y principios de Vue 3. Pinia la reemplaza oficialmente porque:

- API mas simple y menos codigo repetitivo (boilerplate)
- Soporte nativo de TypeScript sin configuracion adicional
- Devtools de Vue la soportan completamente
- Modular por defecto: cada almacen es independiente y se puede importar donde se necesite

### Por que el patron adaptador y no llamar Alpha Vantage directamente

- Permite cambiar de proveedor sin tocar la logica de negocio
- Permite testear los servicios sin hacer llamadas reales a internet
- Evita acoplar el nucleo del sistema a un proveedor externo especifico
- Facilita agregar multiples proveedores en paralelo o como fallback
