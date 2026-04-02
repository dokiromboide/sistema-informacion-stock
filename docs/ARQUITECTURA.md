# Arquitectura del Sistema

Este documento describe como esta construido el sistema, como fluyen los datos entre sus partes y por que se tomaron ciertas decisiones de diseno.

---

## Vision general

El sistema esta dividido en tres capas independientes que se comunican entre si:

```
[ Usuario en el navegador ]
           |
           | HTTP (peticiones y respuestas JSON)
           v
[ Frontend - Vue 3 en puerto 5173 ]
           |
           | HTTP (llamadas a la API REST)
           v
[ Backend - Go en puerto 8080 ]
           |
    +------+------+
    |             |
    v             v
[ CockroachDB ] [ API Alpha Vantage ]
  puerto 26257    (internet)
```

Cada capa tiene una responsabilidad clara y no mezcla funciones con las demas.

---

## Tecnologias y por que se eligieron

### Backend: Go con Gin

**Go** es un lenguaje compilado, tipado estaticamente y con excelente soporte para concurrencia. Se eligio porque:
- Alto rendimiento para servicios HTTP
- Compilacion rapida y binarios pequenos
- Manejo de errores explicito que facilita el razonamiento sobre el codigo
- Excelente soporte para el patron adaptador a traves de interfaces

**Gin** es el framework HTTP mas popular del ecosistema Go. Provee enrutamiento rapido, middleware de CORS, logging y recovery de panics.

### Frontend: Vue 3 con TypeScript

**Vue 3** con la Composition API permite escribir componentes con logica clara y reutilizable. TypeScript agrega tipado estatico que detecta errores en tiempo de desarrollo antes de que lleguen al usuario.

**Pinia** es la libreria oficial de gestion de estado para Vue 3. Reemplaza a Vuex con una API mas simple basada en la Composition API.

**Tailwind CSS** permite construir interfaces directamente en el HTML sin escribir CSS personalizado, lo que acelera el desarrollo y mantiene consistencia visual.

### Base de datos: CockroachDB

CockroachDB es una base de datos compatible con PostgreSQL que soporta SQL estandar. Se eligio por:
- Compatibilidad con el driver pgx de Go (mismo que PostgreSQL)
- Distribucion con Docker sin configuracion adicional
- Interfaz web de administracion incluida en el puerto 8081

---

## Patron adaptador para APIs de stock

### El problema que resuelve

Existen multiples proveedores de datos bursatiles: Alpha Vantage, Polygon.io, Yahoo Finance, Finnhub, entre otros. Cada uno tiene una URL distinta, parametros distintos y formato de respuesta distinto. Si el codigo del sistema dependiera directamente de Alpha Vantage, cambiar de proveedor requeriria modificar decenas de archivos.

### La solucion: interfaz + adaptador

El sistema define una **interfaz** (contrato) que todos los proveedores deben cumplir:

```
AdaptadorStock (interfaz - el contrato universal)
    |
    +-- AlphaVantageAdaptador (implementacion concreta)
    |
    +-- PolygonAdaptador (futuro)
    |
    +-- YahooFinanceAdaptador (futuro)
```

La interfaz define exactamente cuatro operaciones:

| Operacion | Descripcion |
|-----------|-------------|
| `ObtenerCotizacion(simbolo)` | Precio actual de una accion |
| `ObtenerSimbolosPopulares()` | Lista de acciones populares |
| `BuscarSimbolo(consulta)` | Busqueda por nombre o ticker |
| `NombreProveedor()` | Identificador del proveedor |

El resto del sistema (servicios, handlers, base de datos) solo conoce la interfaz, nunca la implementacion concreta. Esto significa que cambiar de Alpha Vantage a Polygon.io requiere unicamente:

1. Crear un archivo nuevo que implemente los cuatro metodos de la interfaz
2. Agregar un caso en la fabrica de adaptadores
3. Cambiar una variable de entorno

### Fabrica de adaptadores

La fabrica (`fabrica.go`) es el unico lugar donde se decide que proveedor usar. Lee la variable de entorno `PROVEEDOR_API` y construye el adaptador correspondiente:

```
Variable PROVEEDOR_API=alpha_vantage
         |
         v
   FabricaAdaptador
         |
         v
   AlphaVantageAdaptador (implementa AdaptadorStock)
         |
         v
   ServicioStock (usa AdaptadorStock, no conoce Alpha Vantage)
```

---

## Flujo de datos completo

### Flujo 1: Sincronizacion de datos

Cuando el usuario hace clic en "Sincronizar":

```
1. Usuario hace clic en boton "Sincronizar" en el navegador
2. Frontend llama POST /api/sincronizar
3. HandlerAcciones recibe la peticion HTTP
4. Llama a ServicioStock.SincronizarAcciones()
5. ServicioStock pide la lista de simbolos populares al Adaptador
6. Adaptador consulta Alpha Vantage: GET alphavantage.co/query?function=SYMBOL_SEARCH
7. Para cada simbolo, Adaptador consulta la cotizacion actual: GET alphavantage.co/query?function=GLOBAL_QUOTE
8. ServicioStock recibe cada cotizacion y llama a RepositorioAcciones.GuardarOActualizar()
9. RepositorioAcciones ejecuta INSERT ... ON CONFLICT UPDATE en CockroachDB
10. Se repite para cada simbolo (10 simbolos populares por defecto)
11. HandlerAcciones responde con JSON: {"sincronizados": 10}
12. Frontend muestra mensaje de exito
```

### Flujo 2: Consulta de acciones con filtros

```
1. Usuario escribe en la barra de busqueda o cambia el orden
2. Frontend llama GET /api/acciones?q=apple&ordenar_por=precio&direccion=desc
3. HandlerAcciones extrae los parametros de la URL
4. Llama a ServicioStock.ObtenerAcciones(filtro)
5. ServicioStock llama a RepositorioAcciones.ObtenerTodas(filtro)
6. RepositorioAcciones ejecuta SELECT con WHERE ILIKE y ORDER BY en CockroachDB
7. CockroachDB devuelve las filas que coinciden
8. HandlerAcciones responde con el JSON de acciones
9. Frontend actualiza la grilla de tarjetas
```

### Flujo 3: Ver detalle de una accion

```
1. Usuario hace clic en una tarjeta de accion
2. Frontend navega a /acciones/AAPL
3. Frontend llama GET /api/acciones/AAPL
4. HandlerAcciones llama a ServicioStock.ObtenerDetalle("AAPL")
5. ServicioStock intenta obtener el dato fresco desde Alpha Vantage (Adaptador)
6. Si Alpha Vantage responde: guarda el dato actualizado en CockroachDB y lo retorna
7. Si Alpha Vantage falla (limite de API, sin internet): retorna el dato guardado en CockroachDB
8. Frontend muestra los datos del detalle
```

### Flujo 4: Recomendaciones

```
1. Usuario navega a /recomendaciones
2. Frontend llama GET /api/recomendaciones
3. HandlerRecomendaciones llama a ServicioRecomendaciones.Recomendar()
4. ServicioRecomendaciones obtiene todas las acciones de CockroachDB
5. Para cada accion calcula un puntaje con tres componentes:
   a. Momentum (60%): que tan positivo es el cambio porcentual del dia
   b. Volumen relativo (25%): volumen vs promedio del universo
   c. Estabilidad (15%): que tan bajo es el spread maximo-minimo
6. Ordena las acciones de mayor a menor puntaje
7. Retorna las primeras 5 con el puntaje y una razon en texto
8. Frontend muestra el ranking con colores y posiciones
```

---

## Estructura de capas del backend

```
Capa HTTP (handlers/)
    Recibe peticiones, valida parametros, retorna JSON
    No conoce la base de datos ni la API externa
         |
         v
Capa de Servicios (servicios/)
    Logica de negocio: cuando actualizar, como calcular puntajes
    Coordina entre el adaptador y el repositorio
    No conoce HTTP ni SQL
         |
    +----+----+
    |         |
    v         v
Adaptadores  Repositorio
(adaptadores/) (repositorio/)
Hablan con    Hablan con
la API        CockroachDB
externa       con SQL
```

Esta separacion en capas tiene una ventaja practica: se puede probar cada capa de forma independiente. Los tests del servicio de recomendaciones, por ejemplo, no necesitan ni internet ni base de datos porque usan datos en memoria.

---

## Estructura de capas del frontend

```
Vistas (vistas/)
    Paginas completas: VistaAcciones, VistaDetalle, VistaRecomendaciones
    Componen componentes y leen del almacen
         |
         v
Componentes (componentes/)
    Piezas visuales reutilizables: TarjetaAccion, BarraBusqueda, SelectorOrden
    Reciben datos via props y emiten eventos
         |
         v
Almacen Pinia (almacen/)
    Estado global de la aplicacion
    Llama al servicio de API y guarda los resultados
         |
         v
Servicio de API (servicios/)
    Encapsula todas las llamadas HTTP al backend con Axios
    Un solo lugar donde se define como hablar con el servidor
```

---

## Seguridad implementada

### Prevencion de inyeccion SQL

El repositorio no concatena strings para construir consultas SQL. Usa parametros posicionales (`$1`, `$2`, ...) que el driver de base de datos maneja de forma segura.

Para el ordenamiento, donde no se pueden usar parametros (SQL no permite parametrizar nombres de columnas), se usa una lista blanca de columnas permitidas en la funcion `columnaValida()`. Cualquier valor que no este en la lista es reemplazado por el valor por defecto.

### CORS configurado

El servidor solo acepta peticiones de los origenes configurados (`localhost:5173` y `localhost:3000`). Peticiones de otros origenes son rechazadas por el navegador antes de llegar al servidor.

### Variables sensibles fuera del codigo

La clave de API y la URL de conexion a la base de datos se leen desde variables de entorno, nunca estan escritas directamente en el codigo fuente. El archivo `.env` esta en el `.gitignore` para que nunca se suba accidentalmente a GitHub.

---

## Pruebas unitarias

### Backend (Go)

Las pruebas estan en `backend/tests/` y cubren:

- **Adaptador Alpha Vantage:** Se crea un servidor HTTP de prueba que simula las respuestas de Alpha Vantage sin hacer llamadas reales a internet. Se prueban escenarios de exito, simbolo invalido y error HTTP 500.

- **Algoritmo de recomendaciones:** Se usa un repositorio en memoria (`ServicioRecomendacionesTesteable`) que permite probar el algoritmo sin base de datos. Se verifica que el orden sea correcto, que el limite de 5 se respete y que el universo vacio retorne una lista vacia.

### Frontend (Vitest)

Las pruebas estan en `frontend/src/__tests__/` y cubren:

- **TarjetaAccion:** Se monta el componente con datos simulados y se verifica que renderice el simbolo, precio, porcentaje, colores correctos y formato del volumen.

- **accionesAlmacen:** Se simula el modulo de API con `vi.mock` para no hacer llamadas HTTP reales. Se prueban el estado inicial, la carga exitosa, el manejo de errores y la sincronizacion.
