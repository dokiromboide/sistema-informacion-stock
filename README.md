# Sistema de Informacion de Stock

Sistema completo para consultar, analizar y recibir recomendaciones de acciones bursatiles. Construido con Go en el backend, Vue 3 en el frontend y CockroachDB como base de datos.

## Tecnologias utilizadas

| Capa | Tecnologia |
|------|-----------|
| Backend | Go 1.22 |
| Frontend | Vue 3, TypeScript, Pinia, Tailwind CSS |
| Base de datos | CockroachDB |
| API de stock | Alpha Vantage (patron adaptador, intercambiable) |
| Contenedores | Docker / Docker Compose |

## Arquitectura

El sistema implementa un **patron adaptador** para la integracion con APIs de stock externas. Esto permite cambiar de proveedor (Alpha Vantage, Yahoo Finance, Polygon.io, etc.) sin modificar la logica de negocio.

```
backend/
  cmd/servidor/         # Punto de entrada de la aplicacion
  internal/
    adaptadores/        # Patron adaptador para APIs de stock
    modelos/            # Estructuras de datos
    repositorio/        # Capa de acceso a CockroachDB
    servicios/          # Logica de negocio
    handlers/           # Controladores HTTP

frontend/
  src/
    componentes/        # Componentes Vue reutilizables
    vistas/             # Paginas de la aplicacion
    almacen/            # Estado global con Pinia
    servicios/          # Llamadas a la API del backend
    tipos/              # Tipos TypeScript
```

## Funcionalidades

1. **Conexion a API de stock** - Sincroniza datos desde Alpha Vantage y los persiste en CockroachDB
2. **API REST** - Endpoints para consultar, buscar, filtrar y ordenar acciones
3. **Interfaz de usuario** - Vista de lista, detalle, busqueda y ordenamiento
4. **Recomendaciones** - Algoritmo que analiza datos bursatiles y recomienda las mejores acciones del dia
5. **Pruebas unitarias** - Cobertura en backend (Go) y frontend (Vitest)

## Configuracion rapida

### Requisitos previos

- Go 1.22+
- Node.js 20+
- Docker y Docker Compose
- Cuenta en [Alpha Vantage](https://www.alphavantage.co/) (gratuita)

### 1. Clonar el repositorio

```bash
git clone https://github.com/dokiromboide/sistema-informacion-stock.git
cd sistema-informacion-stock
```

### 2. Configurar variables de entorno

```bash
cp backend/.env.ejemplo backend/.env
```

Editar `backend/.env` con tus valores:

```env
ALPHA_VANTAGE_API_KEY=tu_api_key_aqui
DB_URL=postgresql://root@localhost:26257/stock_db?sslmode=disable
PUERTO=8080
```

### 3. Levantar CockroachDB

```bash
docker-compose up -d
```

### 4. Ejecutar el backend

```bash
cd backend
go mod download
go run cmd/servidor/main.go
```

### 5. Ejecutar el frontend

```bash
cd frontend
npm install
npm run dev
```

La aplicacion estara disponible en `http://localhost:5173`

## API REST

| Metodo | Endpoint | Descripcion |
|--------|----------|-------------|
| GET | `/api/acciones` | Lista todas las acciones |
| GET | `/api/acciones/:simbolo` | Detalle de una accion |
| GET | `/api/acciones/buscar?q=texto` | Buscar acciones |
| GET | `/api/recomendaciones` | Mejores acciones para invertir hoy |
| POST | `/api/sincronizar` | Sincronizar datos desde la API externa |

## Pruebas

```bash
# Backend
cd backend
go test ./...

# Frontend
cd frontend
npm run test:unit
```

## Cambiar proveedor de API

El sistema implementa la interfaz `AdaptadorStock`. Para cambiar de proveedor:

1. Crear un nuevo archivo en `backend/internal/adaptadores/`
2. Implementar la interfaz `AdaptadorStock`
3. Cambiar la variable `PROVEEDOR_API` en `.env`

## Licencia

MIT
