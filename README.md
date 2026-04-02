# Sistema de Informacion de Stock

Sistema web completo para consultar precios de acciones bursatiles en tiempo real. Permite buscar empresas, filtrar resultados, ver detalles de cada accion y recibir recomendaciones de inversion automaticas basadas en un algoritmo de momentum y volumen.

Construido con Go en el servidor, Vue 3 en el navegador y CockroachDB como base de datos, corriendo en contenedor Docker.

---

## Indice

1. [Que hace este sistema](#que-hace-este-sistema)
2. [Requisitos previos](#requisitos-previos)
3. [Obtener la clave de API de Alpha Vantage](#obtener-la-clave-de-api-de-alpha-vantage)
4. [Instalacion paso a paso](#instalacion-paso-a-paso)
5. [Acceder a la aplicacion](#acceder-a-la-aplicacion)
6. [Estructura del proyecto](#estructura-del-proyecto)
7. [Como cambiar de proveedor de API](#como-cambiar-de-proveedor-de-api)
8. [Endpoints de la API](#endpoints-de-la-api)
9. [Ejecutar las pruebas](#ejecutar-las-pruebas)
10. [Solucion de problemas](#solucion-de-problemas)
11. [Tecnologias utilizadas](#tecnologias-utilizadas)

---

## Que hace este sistema

- Conecta con la API de Alpha Vantage para obtener precios reales del mercado bursatil
- Guarda los datos en una base de datos local para acceso rapido sin depender de la red
- Muestra un listado de acciones con busqueda por nombre o simbolo y opciones de ordenamiento
- Permite ver el detalle de cada accion: precio actual, apertura, maximo, minimo y volumen
- Genera un ranking de las 5 mejores acciones del dia con puntaje y justificacion
- Arquitectura con patron adaptador que permite cambiar de proveedor de datos sin modificar el resto del sistema

---

## Requisitos previos

Antes de instalar el proyecto, es necesario tener los siguientes programas en la computadora. Esta seccion explica como instalar cada uno desde cero.

---

### 1. Go (lenguaje del servidor)

Go es el lenguaje de programacion que ejecuta el servidor de esta aplicacion.

**Pasos de instalacion:**

1. Abrir el navegador e ir a: https://go.dev/dl/
2. Descargar el instalador para Windows (archivo con extension `.msi`)
3. Ejecutar el instalador descargado y seguir todos los pasos (aceptar las opciones por defecto)
4. Una vez terminada la instalacion, abrir una nueva ventana de terminal y escribir:

```
go version
```

El resultado debe ser similar a: `go version go1.22.0 windows/amd64`

Si aparece ese mensaje, Go esta correctamente instalado.

**Version minima requerida:** Go 1.22

---

### 2. Node.js (necesario para el frontend)

Node.js es el entorno que permite ejecutar el proyecto Vue 3, es decir, la interfaz visual que el usuario ve en el navegador.

**Pasos de instalacion:**

1. Ir a: https://nodejs.org/
2. Descargar la version **LTS** (aparece a la izquierda con la etiqueta "Recommended For Most Users")
3. Ejecutar el instalador descargado y seguir los pasos
4. Abrir una nueva ventana de terminal y verificar la instalacion:

```
node --version
npm --version
```

Deben aparecer versiones similares a: `v20.11.0` y `10.2.4`

**Version minima requerida:** Node.js 20

---

### 3. Docker Desktop (para la base de datos)

Docker permite ejecutar la base de datos CockroachDB dentro de un contenedor, sin necesidad de instalar la base de datos manualmente en el sistema operativo.

**Pasos de instalacion:**

1. Ir a: https://www.docker.com/products/docker-desktop/
2. Descargar Docker Desktop para Windows
3. Ejecutar el instalador
4. Reiniciar el computador si el instalador lo solicita
5. Abrir Docker Desktop desde el menu de inicio y esperar hasta que el icono en la barra de tareas muestre el estado "Running" (en funcionamiento)
6. Verificar en terminal:

```
docker --version
```

Debe aparecer algo similar a: `Docker version 25.0.0`

**Nota importante:** Docker Desktop debe estar abierto y en estado "Running" cada vez que se use el proyecto. Si Docker no esta corriendo, la base de datos no podra iniciarse.

---

### 4. Git (para descargar el proyecto)

Git es la herramienta que permite descargar el codigo fuente del proyecto desde GitHub.

**Pasos de instalacion:**

1. Ir a: https://git-scm.com/download/win
2. Descargar el instalador para Windows
3. Ejecutar el instalador (aceptar todas las opciones por defecto)
4. Verificar en terminal:

```
git --version
```

---

## Obtener la clave de API de Alpha Vantage

El sistema necesita una clave de acceso gratuita de Alpha Vantage para consultar los precios de acciones en tiempo real. El registro es gratuito y no requiere tarjeta de credito.

**Pasos para obtener la clave:**

1. Abrir el navegador e ir a: https://www.alphavantage.co/support/#api-key
2. Completar el formulario con los siguientes datos:
   - **First Name:** nombre
   - **Last Name:** apellido
   - **Email:** correo electronico (se usa para recibir la clave)
   - En la pregunta sobre el tipo de uso, seleccionar "Student / Individual"
3. Hacer clic en el boton **"GET FREE API KEY"**
4. La clave aparece inmediatamente en pantalla con un formato similar a: `ABCDEF1234567890`
5. Copiar y guardar esa clave, se necesitara en el siguiente paso

**Limitaciones del plan gratuito:**

- 25 solicitudes por dia al servidor de Alpha Vantage
- Suficiente para pruebas, desarrollo y uso personal
- Si se necesitan mas solicitudes, Alpha Vantage ofrece planes de pago

---

## Instalacion paso a paso

Con todos los requisitos instalados y la clave de API en mano, seguir estos pasos en orden.

---

### Paso 1: Descargar el proyecto

Abrir una terminal (Git Bash o PowerShell) y ejecutar:

```bash
git clone https://github.com/dokiromboide/sistema-informacion-stock.git
cd sistema-informacion-stock
```

Esto descarga todo el codigo del proyecto en la carpeta actual.

---

### Paso 2: Configurar las variables de entorno

Las variables de entorno son configuraciones del sistema: la clave de API, la direccion de la base de datos, el puerto del servidor, etc.

1. Copiar el archivo de ejemplo `.env.ejemplo` y crear el archivo `.env`:

```bash
cp backend/.env.ejemplo backend/.env
```

2. Abrir el archivo recien creado `backend/.env` con cualquier editor de texto (Bloc de notas, Notepad++, VS Code):

```
PUERTO=8080
PROVEEDOR_API=alpha_vantage
ALPHA_VANTAGE_API_KEY=tu_api_key_aqui
ALPHA_VANTAGE_URL=https://www.alphavantage.co/query
DB_URL=postgresql://root@localhost:26257/stock_db?sslmode=disable
```

3. Reemplazar `tu_api_key_aqui` con la clave obtenida en la seccion anterior, por ejemplo:

```
ALPHA_VANTAGE_API_KEY=ABCDEF1234567890
```

4. Guardar el archivo.

Los demas valores pueden dejarse exactamente como estan.

---

### Paso 3: Levantar la base de datos con Docker

Verificar que Docker Desktop este abierto y en estado "Running" antes de continuar.

Desde la carpeta raiz del proyecto (donde esta el archivo `docker-compose.yml`), ejecutar:

```bash
export PATH="$PATH:/c/Program Files/Docker/Docker/resources/bin"
docker-compose up -d
```

La primera vez, Docker descarga automaticamente la imagen de CockroachDB desde internet. Esto puede tardar varios minutos segun la velocidad de la conexion.

**Verificar que la base de datos inicio correctamente:**

```bash
docker ps
```

Debe aparecer en la lista un contenedor llamado `stock_cockroachdb` con estado `Up`.

El panel de administracion web de CockroachDB queda disponible en: `http://localhost:8081`

---

### Paso 4: Iniciar el servidor (backend)

El servidor es el programa en Go que procesa las consultas, se comunica con la base de datos y expone la API REST.

Abrir una terminal, ubicarse en la carpeta del proyecto y ejecutar:

```bash
cd backend
go mod download
go run cmd/servidor/main.go
```

`go mod download` descarga las dependencias de Go la primera vez. Con una conexion a internet normal, tarda menos de un minuto.

**Que debe aparecer en la terminal:**

```
Proveedor de API: AlphaVantage
Servidor iniciado en http://localhost:8080
```

Dejar esta terminal abierta. El servidor queda corriendo en segundo plano mientras la terminal este activa.

---

### Paso 5: Cargar los datos iniciales

Con el servidor corriendo, es necesario hacer una sincronizacion inicial para que la base de datos tenga acciones que mostrar.

Desde otra terminal, ejecutar:

```bash
curl -X POST http://localhost:8080/api/sincronizar
```

O bien, una vez que el frontend este corriendo (siguiente paso), hacer clic en el boton **"Sincronizar"** que aparece en la pantalla principal.

La sincronizacion descarga los datos de 10 acciones populares (AAPL, GOOGL, MSFT, AMZN, TSLA, META, NVDA, JPM, V, JNJ).

---

### Paso 6: Iniciar el frontend (interfaz de usuario)

El frontend es la aplicacion web que el usuario ve en el navegador.

Abrir una segunda terminal, ubicarse en la carpeta del proyecto y ejecutar:

```bash
cd frontend
npm install
npm run dev
```

`npm install` descarga las dependencias del frontend la primera vez (puede tardar un minuto).

**Que debe aparecer:**

```
  VITE v5.3.2  ready in 800 ms

  Local:   http://localhost:5173/
```

Dejar esta terminal abierta.

---

## Acceder a la aplicacion

Con el servidor y el frontend corriendo, abrir el navegador e ingresar a:

```
http://localhost:5173
```

### Pantallas disponibles

**Acciones (pantalla principal):**

- Barra de busqueda para filtrar por simbolo o nombre de empresa
- Selector de orden por: simbolo, precio, variacion porcentual o volumen
- Tarjetas con cada accion mostrando: simbolo, nombre de la empresa, precio actual y variacion del dia en verde (subio) o rojo (bajo)
- Boton "Sincronizar" para actualizar los datos desde Alpha Vantage

**Detalle de una accion:**

- Al hacer clic en una tarjeta se accede a la pagina de detalle
- Muestra: precio actual, variacion del dia, precio de apertura, maximo, minimo y volumen de operaciones

**Recomendaciones:**

- Ranking con las 5 mejores acciones del dia segun el algoritmo de momentum y volumen
- Cada recomendacion incluye el puntaje calculado y la justificacion en texto

---

## Estructura del proyecto

```
sistema-informacion-stock/
│
├── backend/                         Servidor en Go (API REST)
│   ├── cmd/
│   │   └── servidor/
│   │       └── main.go              Punto de entrada: inicia DB, servicios y servidor HTTP
│   ├── internal/
│   │   ├── adaptadores/             Conexion con APIs externas de stock
│   │   │   ├── interfaz.go          Contrato que todo proveedor de datos debe cumplir
│   │   │   ├── alpha_vantage.go     Implementacion concreta para Alpha Vantage
│   │   │   └── fabrica.go           Selecciona el proveedor segun la variable de entorno
│   │   ├── modelos/
│   │   │   └── accion.go            Estructuras de datos: Accion, Filtro, Recomendacion
│   │   ├── repositorio/             Acceso a CockroachDB
│   │   │   ├── base_de_datos.go     Gestion del pool de conexiones
│   │   │   └── repositorio_acciones.go  Operaciones CRUD sobre las acciones
│   │   ├── servicios/               Logica de negocio
│   │   │   ├── servicio_stock.go    Sincronizacion de datos desde la API externa
│   │   │   └── servicio_recomendaciones.go  Algoritmo de recomendaciones
│   │   └── handlers/                Controladores HTTP
│   │       ├── enrutador.go         Definicion de rutas y configuracion de CORS
│   │       ├── acciones.go          Endpoints de acciones
│   │       └── recomendaciones.go   Endpoint de recomendaciones
│   ├── docs/
│   │   └── api.md                   Documentacion de los endpoints REST
│   ├── tests/                       Pruebas unitarias del backend
│   ├── .env.ejemplo                 Plantilla del archivo de configuracion
│   └── go.mod                       Dependencias del proyecto Go
│
├── frontend/                        Interfaz visual en Vue 3 + TypeScript
│   ├── src/
│   │   ├── almacen/
│   │   │   └── accionesAlmacen.ts   Estado global de la aplicacion (Pinia)
│   │   ├── componentes/             Componentes visuales reutilizables
│   │   │   ├── BarraBusqueda.vue    Campo de busqueda
│   │   │   ├── SelectorOrden.vue    Desplegable de ordenamiento
│   │   │   └── TarjetaAccion.vue    Tarjeta de resumen de una accion
│   │   ├── servicios/
│   │   │   └── apiStock.ts          Llamadas HTTP al backend con Axios
│   │   ├── tipos/
│   │   │   └── index.ts             Interfaces TypeScript de los modelos
│   │   └── vistas/                  Paginas de la aplicacion
│   │       ├── VistaAcciones.vue    Listado principal con busqueda y orden
│   │       ├── VistaDetalleAccion.vue   Detalle de una accion
│   │       └── VistaRecomendaciones.vue  Ranking de recomendaciones
│   └── package.json                 Dependencias del proyecto frontend
│
├── docs/
│   └── ARQUITECTURA.md              Documentacion de la arquitectura del sistema
├── docker-compose.yml               Configuracion de CockroachDB con Docker
├── GUIA_INICIO_RAPIDO.md            Guia resumida para tener el sistema corriendo
└── README.md                        Este archivo
```

---

## Como cambiar de proveedor de API

El sistema esta disenado para que cambiar de proveedor de datos bursatiles sea sencillo y no requiera modificar la logica de negocio ni la base de datos. Esto se logra mediante el **patron adaptador**.

### Que es el patron adaptador

El patron adaptador define una interfaz comun (un contrato) que todos los proveedores de datos deben cumplir. El resto del sistema solo conoce esa interfaz, nunca los detalles de implementacion de cada proveedor.

Es similar a un tomacorriente universal: sin importar que aparato se conecte, el tomacorriente siempre expone la misma interfaz estandar. Cada aparato se adapta al tomacorriente, no al reves.

### Interfaz que debe implementar cada proveedor

Todo nuevo proveedor de datos debe implementar las siguientes cuatro operaciones definidas en `backend/internal/adaptadores/interfaz.go`:

| Metodo | Descripcion |
|--------|-------------|
| `ObtenerCotizacion(simbolo string)` | Retorna el precio actual y datos del dia de una accion |
| `ObtenerSimbolosPopulares()` | Retorna la lista de simbolos que se sincronizan por defecto |
| `BuscarSimbolo(consulta string)` | Busca acciones por nombre o ticker |
| `NombreProveedor() string` | Retorna el nombre del proveedor (usado en logs) |

### Pasos para agregar un nuevo proveedor (ejemplo: Polygon.io)

1. Crear el archivo `backend/internal/adaptadores/polygon.go`
2. Implementar los cuatro metodos de la interfaz `AdaptadorStock`
3. Agregar el nuevo proveedor en `backend/internal/adaptadores/fabrica.go`:

```go
case "polygon":
    return NuevoPolygonAdaptador(os.Getenv("POLYGON_API_KEY"), os.Getenv("POLYGON_URL")), nil
```

4. En el archivo `backend/.env`, cambiar las variables:

```
PROVEEDOR_API=polygon
POLYGON_API_KEY=tu_clave_de_polygon
POLYGON_URL=https://api.polygon.io
```

5. Reiniciar el servidor.

El frontend, los handlers, los servicios y la base de datos no requieren ninguna modificacion.

---

## Endpoints de la API

El servidor expone los siguientes endpoints que pueden consultarse directamente desde el navegador o herramientas como curl o Postman:

| Metodo | URL | Descripcion |
|--------|-----|-------------|
| GET | `http://localhost:8080/salud` | Verificar que el servidor esta activo |
| GET | `http://localhost:8080/api/acciones` | Listar todas las acciones con filtros opcionales |
| GET | `http://localhost:8080/api/acciones/AAPL` | Ver detalle de una accion por simbolo |
| GET | `http://localhost:8080/api/acciones/buscar?q=apple` | Buscar acciones por nombre o simbolo |
| POST | `http://localhost:8080/api/sincronizar` | Sincronizar datos desde la API externa |
| GET | `http://localhost:8080/api/recomendaciones` | Obtener las 5 mejores recomendaciones del dia |

**Parametros de filtro para `/api/acciones`:**

| Parametro | Valores posibles | Descripcion |
|-----------|-----------------|-------------|
| `q` | cualquier texto | Filtra por nombre o simbolo |
| `ordenar_por` | `simbolo`, `precio`, `cambio_porcentual`, `volumen` | Campo de ordenamiento |
| `direccion` | `asc`, `desc` | Orden ascendente o descendente |
| `limite` | numero entero | Cantidad de resultados por pagina |
| `pagina` | numero entero | Pagina de resultados |

---

## Ejecutar las pruebas

**Pruebas del servidor (Go):**

```bash
cd backend
go test ./...
```

**Pruebas del frontend (Vitest):**

```bash
cd frontend
npm run test:unit
```

---

## Solucion de problemas

### El servidor no inicia: "DB_URL no esta definida"

**Causa:** El archivo `backend/.env` no existe.

**Solucion:**

```bash
cp backend/.env.ejemplo backend/.env
```

Abrir `backend/.env` y completar los valores, especialmente `ALPHA_VANTAGE_API_KEY`.

---

### El servidor no inicia: "error al conectar con CockroachDB"

**Causa:** La base de datos no esta corriendo.

**Solucion:**

1. Verificar que Docker Desktop este abierto y en estado "Running"
2. Ejecutar:

```bash
export PATH="$PATH:/c/Program Files/Docker/Docker/resources/bin"
docker-compose up -d
```

3. Esperar 10 segundos y volver a intentar iniciar el servidor

---

### Docker dice "command not found"

**Causa:** Docker Desktop no agrego los ejecutables al PATH del sistema.

**Solucion:** Ejecutar esto en la terminal antes de usar cualquier comando de Docker:

```bash
export PATH="$PATH:/c/Program Files/Docker/Docker/resources/bin"
```

Para que sea permanente, agregar esa linea al archivo `~/.bashrc` o `~/.zshrc`.

---

### El frontend muestra "Error de conexion con el servidor"

**Causa:** El servidor Go no esta corriendo.

**Solucion:**

```bash
cd backend
go run cmd/servidor/main.go
```

Verificar que aparezca el mensaje `Servidor iniciado en http://localhost:8080` y recargar el navegador.

---

### La sincronizacion no trae datos

**Causa posible 1:** La clave de API no es valida o tiene el valor por defecto.

**Solucion:** Abrir `backend/.env` y verificar que `ALPHA_VANTAGE_API_KEY` tenga una clave real (no el texto `tu_api_key_aqui`).

**Causa posible 2:** Se alcanzaron las 25 solicitudes diarias del plan gratuito de Alpha Vantage.

**Solucion:** Esperar hasta el dia siguiente o registrar una nueva clave con otro correo electronico.

---

### "go: command not found"

**Causa:** Go no esta instalado o la terminal necesita reiniciarse despues de la instalacion.

**Solucion:**

1. Cerrar la terminal completamente y abrir una nueva
2. Ejecutar `go version` para verificar
3. Si sigue sin funcionar, reinstalar Go desde https://go.dev/dl/

---

### "npm: command not found"

**Causa:** Node.js no esta instalado o la terminal necesita reiniciarse.

**Solucion:**

1. Cerrar la terminal completamente y abrir una nueva
2. Ejecutar `node --version` para verificar
3. Si sigue sin funcionar, reinstalar Node.js desde https://nodejs.org/

---

## Tecnologias utilizadas

| Tecnologia | Version | Uso en el proyecto |
|------------|---------|-------------------|
| Go | 1.22 | Servidor y API REST |
| Gin | 1.10 | Framework HTTP para Go |
| Vue 3 | 3.4 | Interfaz de usuario reactiva |
| TypeScript | 5.4 | Tipado estatico del frontend |
| Pinia | 2.1 | Gestion de estado global del frontend |
| Vue Router | 4.3 | Navegacion entre pantallas |
| Axios | 1.7 | Cliente HTTP del frontend |
| Tailwind CSS | 3.4 | Estilos visuales utilitarios |
| Vite | 5.3 | Servidor de desarrollo y compilacion del frontend |
| CockroachDB | 23.2 | Base de datos SQL distribuida |
| Docker | - | Contenedor para la base de datos |
| Alpha Vantage | - | Proveedor de datos bursatiles en tiempo real |

---

## Licencia

MIT - Libre para uso personal y comercial.
