# Sistema de Informacion de Stock

Sistema web completo para consultar precios de acciones bursatiles en tiempo real, buscar empresas, ordenar resultados y recibir recomendaciones de inversion automaticas. Construido con tecnologias modernas: Go en el servidor, Vue 3 en el navegador y CockroachDB como base de datos.

---

## Que hace este sistema

- Conecta con la API de Alpha Vantage para obtener precios reales de acciones del mercado
- Guarda los datos en una base de datos local para acceso rapido
- Muestra una tabla con todas las acciones disponibles, con opcion de buscar y ordenar
- Permite ver el detalle de cada accion (precio actual, apertura, maximo, minimo, volumen)
- Genera recomendaciones automaticas de las mejores acciones para invertir ese dia
- Tiene una arquitectura que permite cambiar facilmente de proveedor de datos (Alpha Vantage, Polygon.io, Yahoo Finance, etc.)

---

## Requisitos previos

Antes de instalar el proyecto necesitas tener instalados los siguientes programas en tu computadora. A continuacion se explica como instalar cada uno.

### 1. Go (lenguaje del servidor)

Go es el lenguaje de programacion que usa el servidor de esta aplicacion.

**Como instalar Go:**

1. Ir a https://go.dev/dl/
2. Descargar el instalador para Windows (archivo `.msi`)
3. Ejecutar el instalador y seguir los pasos (aceptar todo por defecto)
4. Abrir una nueva ventana de terminal y verificar la instalacion ejecutando:

```
go version
```

Deberia aparecer algo como: `go version go1.22.0 windows/amd64`

**Version requerida:** Go 1.22 o superior

---

### 2. Node.js (necesario para el frontend)

Node.js es el entorno que permite ejecutar el proyecto de Vue 3 (la interfaz visual).

**Como instalar Node.js:**

1. Ir a https://nodejs.org/
2. Descargar la version **LTS** (la recomendada, aparece a la izquierda)
3. Ejecutar el instalador y seguir los pasos
4. Verificar la instalacion:

```
node --version
npm --version
```

Deberia aparecer algo como: `v20.11.0` y `10.2.4`

**Version requerida:** Node.js 20 o superior

---

### 3. Docker Desktop (para la base de datos)

Docker permite levantar la base de datos CockroachDB sin necesidad de instalarla manualmente.

**Como instalar Docker Desktop:**

1. Ir a https://www.docker.com/products/docker-desktop/
2. Descargar Docker Desktop para Windows
3. Ejecutar el instalador
4. Reiniciar el computador si lo solicita
5. Abrir Docker Desktop (aparece en el menu de inicio) y esperar a que diga "Running"
6. Verificar en terminal:

```
docker --version
```

Deberia aparecer algo como: `Docker version 25.0.0`

**Nota importante:** Docker Desktop debe estar abierto y corriendo antes de ejecutar el proyecto.

---

### 4. Git (para clonar el proyecto)

Git permite descargar el codigo del proyecto desde GitHub.

**Como instalar Git:**

1. Ir a https://git-scm.com/download/win
2. Descargar el instalador para Windows
3. Ejecutar el instalador (aceptar todas las opciones por defecto)
4. Verificar:

```
git --version
```

---

## Paso 1: Descargar el proyecto

Abrir una terminal (Git Bash o PowerShell) y ejecutar:

```bash
git clone https://github.com/dokiromboide/sistema-informacion-stock.git
cd sistema-informacion-stock
```

---

## Paso 2: Obtener la clave de la API de Alpha Vantage

El sistema necesita una clave gratuita de Alpha Vantage para obtener datos de acciones reales.

**Pasos para obtener la clave:**

1. Ir a https://www.alphavantage.co/support/#api-key
2. Completar el formulario:
   - **First Name:** tu nombre
   - **Last Name:** tu apellido
   - **Email:** tu correo electronico
   - En la pregunta sobre uso, seleccionar "Student / Individual"
3. Hacer clic en **"GET FREE API KEY"**
4. Anotar la clave que aparece en pantalla, tiene un formato similar a: `ABCDEF1234567890`

**Limitaciones de la clave gratuita:**
- 25 solicitudes por dia
- Suficiente para pruebas y uso personal

---

## Paso 3: Configurar las variables de entorno

Las variables de entorno son configuraciones que el sistema necesita para funcionar (clave de API, conexion a la base de datos, etc.).

1. Abrir la carpeta `backend` dentro del proyecto
2. Encontrar el archivo llamado `.env.ejemplo`
3. Hacer una copia de ese archivo y nombrarla `.env` (sin la palabra "ejemplo")

**En terminal:**

```bash
cd backend
cp .env.ejemplo .env
```

4. Abrir el archivo `.env` con cualquier editor de texto (Bloc de notas, VS Code, etc.)
5. Reemplazar `tu_api_key_aqui` con la clave que obtuviste en el Paso 2:

```
PUERTO=8080
PROVEEDOR_API=alpha_vantage
ALPHA_VANTAGE_API_KEY=ABCDEF1234567890
ALPHA_VANTAGE_URL=https://www.alphavantage.co/query
DB_URL=postgresql://root@localhost:26257/stock_db?sslmode=disable
```

Guardar el archivo.

---

## Paso 4: Levantar la base de datos con Docker

La base de datos CockroachDB se levanta automaticamente con Docker Compose.

**Verificar que Docker Desktop este abierto** (debe estar corriendo en segundo plano).

Desde la raiz del proyecto (donde esta el archivo `docker-compose.yml`), ejecutar:

```bash
export PATH="$PATH:/c/Program Files/Docker/Docker/resources/bin"
docker-compose up -d
```

La primera vez esto puede tardar varios minutos porque descarga la imagen de CockroachDB.

**Verificar que funciono:**

```bash
docker ps
```

Debe aparecer un contenedor llamado `stock_cockroachdb` con estado `Up`.

Para ver la interfaz web de administracion de CockroachDB, abrir en el navegador:
`http://localhost:8081`

---

## Paso 5: Iniciar el servidor (backend)

El servidor es el programa en Go que expone la API REST y se comunica con la base de datos.

Abrir una terminal, ir a la carpeta `backend` y ejecutar:

```bash
cd backend
go mod download
go run cmd/servidor/main.go
```

La primera vez `go mod download` descargara las dependencias necesarias (puede tardar unos minutos con conexion a internet).

**Que deberia verse en la terminal:**

```
Proveedor de API: AlphaVantage
Servidor iniciado en http://localhost:8080
```

El servidor queda corriendo en esta terminal. No cerrar esta ventana.

---

## Paso 6: Cargar los datos de acciones

Con el servidor corriendo, ejecutar una sincronizacion inicial para cargar las acciones en la base de datos. Abrir otra terminal y ejecutar:

```bash
curl -X POST http://localhost:8080/api/sincronizar
```

O simplemente, una vez que el frontend este corriendo, hacer clic en el boton **"Sincronizar"** que aparece en la pantalla principal.

---

## Paso 7: Iniciar el frontend (interfaz de usuario)

El frontend es la interfaz visual que se ve en el navegador.

Abrir otra terminal, ir a la carpeta `frontend` y ejecutar:

```bash
cd frontend
npm install
npm run dev
```

`npm install` descarga las dependencias del proyecto (solo la primera vez).

**Que deberia verse:**

```
  VITE v5.3.2  ready in 800 ms

  Local:   http://localhost:5173/
```

---

## Paso 8: Acceder a la aplicacion

Abrir el navegador y entrar a:

```
http://localhost:5173
```

### Que se ve en la aplicacion

**Pantalla principal - Lista de acciones:**
- Una barra de busqueda para buscar por simbolo o nombre de empresa
- Un selector para ordenar por precio, variacion porcentual o volumen
- Tarjetas con las acciones disponibles mostrando: simbolo, nombre, precio actual y variacion del dia en color verde (subio) o rojo (bajo)
- Un boton "Sincronizar" para actualizar los datos desde Alpha Vantage

**Pantalla de detalle de una accion:**
- Al hacer clic en cualquier tarjeta se ve el detalle completo:
- Precio actual, variacion del dia
- Precio de apertura, maximo y minimo del dia
- Volumen de operaciones

**Pantalla de recomendaciones:**
- Un ranking con las mejores 5 acciones del dia
- Cada recomendacion incluye el puntaje calculado y la razon de la recomendacion

---

## Estructura del proyecto

```
sistema-informacion-stock/
│
├── backend/                    Servidor en Go
│   ├── cmd/
│   │   └── servidor/
│   │       └── main.go         Punto de entrada del servidor
│   ├── internal/
│   │   ├── adaptadores/        Conexion con APIs externas de stock
│   │   │   ├── interfaz.go     Contrato que deben cumplir todos los proveedores
│   │   │   ├── alpha_vantage.go  Implementacion para Alpha Vantage
│   │   │   └── fabrica.go      Selecciona el proveedor segun configuracion
│   │   ├── modelos/            Estructuras de datos
│   │   ├── repositorio/        Acceso a la base de datos CockroachDB
│   │   ├── servicios/          Logica de negocio (sincronizacion, recomendaciones)
│   │   └── handlers/           Controladores de la API REST
│   ├── docs/
│   │   └── api.md              Documentacion de los endpoints
│   ├── tests/                  Pruebas unitarias del backend
│   ├── .env.ejemplo            Plantilla de configuracion
│   └── go.mod                  Dependencias de Go
│
├── frontend/                   Interfaz visual en Vue 3
│   ├── src/
│   │   ├── almacen/            Estado global de la aplicacion (Pinia)
│   │   ├── componentes/        Componentes visuales reutilizables
│   │   ├── servicios/          Llamadas a la API del backend
│   │   ├── tipos/              Tipos TypeScript
│   │   └── vistas/             Paginas de la aplicacion
│   ├── index.html
│   └── package.json
│
├── docker-compose.yml          Configuracion de CockroachDB con Docker
└── README.md                   Este archivo
```

---

## Como cambiar de proveedor de API

El sistema esta disenado para que sea facil cambiar de proveedor de datos bursatiles sin modificar el resto del codigo. Esto se logra mediante el **patron adaptador**.

### Que es el patron adaptador (explicacion simple)

Imagine que tiene un cargador de celular universal: sin importar que marca de celular tenga, el cargador universal se adapta a todos. El patron adaptador funciona igual: el sistema no conoce los detalles de cada API de stock, solo conoce una interfaz comun (el "conector universal"), y cada proveedor implementa esa interfaz a su manera.

### Como agregar un nuevo proveedor

**Ejemplo: agregar Polygon.io**

1. Crear un archivo nuevo: `backend/internal/adaptadores/polygon.go`
2. Implementar las funciones requeridas:
   - `ObtenerCotizacion(simbolo string)` - obtener el precio de una accion
   - `ObtenerSimbolosPopulares()` - obtener lista de simbolos populares
   - `BuscarSimbolo(consulta string)` - buscar acciones por nombre
   - `NombreProveedor()` - retornar el nombre del proveedor

3. Registrar el nuevo proveedor en `backend/internal/adaptadores/fabrica.go`

4. Cambiar en el archivo `.env`:

```
PROVEEDOR_API=polygon
POLYGON_API_KEY=tu_clave_de_polygon
```

El resto del sistema (base de datos, API REST, frontend) no necesita modificacion.

---

## Endpoints de la API

El servidor expone los siguientes endpoints que tambien pueden usarse directamente:

| Metodo | URL | Descripcion |
|--------|-----|-------------|
| GET | `http://localhost:8080/salud` | Verificar que el servidor esta activo |
| GET | `http://localhost:8080/api/acciones` | Listar todas las acciones |
| GET | `http://localhost:8080/api/acciones/AAPL` | Ver detalle de Apple |
| GET | `http://localhost:8080/api/acciones/buscar?q=apple` | Buscar acciones |
| POST | `http://localhost:8080/api/sincronizar` | Actualizar datos desde la API |
| GET | `http://localhost:8080/api/recomendaciones` | Ver recomendaciones del dia |

---

## Solucion de problemas comunes

### El servidor no inicia y dice "DB_URL no esta definida"

**Causa:** El archivo `.env` no fue creado o no esta en la carpeta correcta.

**Solucion:**
1. Verificar que existe el archivo `backend/.env` (no `backend/.env.ejemplo`)
2. Si no existe, crearlo copiando el archivo ejemplo:
   ```bash
   cp backend/.env.ejemplo backend/.env
   ```
3. Abrir `backend/.env` y completar los valores

---

### El servidor no inicia y dice "error al conectar con CockroachDB"

**Causa:** La base de datos no esta corriendo.

**Solucion:**
1. Verificar que Docker Desktop este abierto y corriendo (icono en la barra de tareas)
2. Ejecutar:
   ```bash
   export PATH="$PATH:/c/Program Files/Docker/Docker/resources/bin"
   docker-compose up -d
   ```
3. Esperar 10 segundos y volver a intentar iniciar el servidor

---

### Docker dice "command not found"

**Causa:** Docker Desktop no agrego los comandos al PATH del sistema.

**Solucion:** Ejecutar esto en la terminal antes de usar Docker:

```bash
export PATH="$PATH:/c/Program Files/Docker/Docker/resources/bin"
```

---

### El frontend dice "Error de conexion con el servidor"

**Causa:** El servidor (backend) no esta corriendo.

**Solucion:**
1. Abrir una terminal y ejecutar:
   ```bash
   cd backend
   go run cmd/servidor/main.go
   ```
2. Verificar que diga "Servidor iniciado en http://localhost:8080"
3. Recargar el navegador

---

### La sincronizacion no trae datos o dice "simbolo no encontrado"

**Causa posible 1:** La clave de Alpha Vantage no es valida o no fue configurada.

**Solucion:**
1. Abrir `backend/.env`
2. Verificar que `ALPHA_VANTAGE_API_KEY` tenga una clave real (no el texto "tu_api_key_aqui")
3. Reiniciar el servidor

**Causa posible 2:** Se alcanzaron las 25 solicitudes diarias del plan gratuito.

**Solucion:** Esperar hasta el dia siguiente o usar una clave diferente.

---

### go run dice "go: command not found"

**Causa:** Go no fue instalado correctamente o la terminal necesita reiniciarse.

**Solucion:**
1. Cerrar la terminal completamente y abrirla de nuevo
2. Ejecutar `go version` para verificar
3. Si sigue sin funcionar, reinstalar Go desde https://go.dev/dl/

---

### npm dice "npm: command not found"

**Causa:** Node.js no fue instalado o la terminal necesita reiniciarse.

**Solucion:**
1. Cerrar la terminal y abrirla de nuevo
2. Ejecutar `node --version` para verificar
3. Si sigue sin funcionar, reinstalar Node.js desde https://nodejs.org/

---

## Tecnologias utilizadas

| Tecnologia | Uso | Version |
|------------|-----|---------|
| Go | Servidor y API REST | 1.22 |
| Gin | Framework HTTP de Go | 1.10 |
| Vue 3 | Interfaz de usuario | 3.4 |
| TypeScript | Tipado estatico del frontend | 5.4 |
| Pinia | Gestion de estado del frontend | 2.1 |
| Tailwind CSS | Estilos visuales | 3.4 |
| CockroachDB | Base de datos | 23.2 |
| Docker | Contenedor de la base de datos | - |
| Alpha Vantage | API de datos bursatiles | - |

---

## Ejecutar las pruebas

**Pruebas del servidor (Go):**

```bash
cd backend
go test ./...
```

**Pruebas de la interfaz (Vue):**

```bash
cd frontend
npm run test:unit
```

---

## Licencia

MIT - Libre para uso personal y comercial.
