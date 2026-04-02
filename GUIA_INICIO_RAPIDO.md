# Guia de Inicio Rapido

Esta guia asume que ya tenes instalado: Go 1.22+, Node.js 20+, Docker Desktop y Git.
Si aun no los tenes, seguir primero las instrucciones completas del README.md.

---

## Paso 1 - Clonar el proyecto

```bash
git clone https://github.com/dokiromboide/sistema-informacion-stock.git
cd sistema-informacion-stock
```

---

## Paso 2 - Crear el archivo de configuracion

```bash
cp backend/.env.ejemplo backend/.env
```

Abrir `backend/.env` con cualquier editor de texto y reemplazar `tu_api_key_aqui` con tu clave de Alpha Vantage.

Clave gratuita en: https://www.alphavantage.co/support/#api-key

El archivo debe quedar asi (con tu clave real):

```
PUERTO=8080
PROVEEDOR_API=alpha_vantage
ALPHA_VANTAGE_API_KEY=ABCDEF1234567890
ALPHA_VANTAGE_URL=https://www.alphavantage.co/query
DB_URL=postgresql://root@localhost:26257/stock_db?sslmode=disable
```

---

## Paso 3 - Levantar la base de datos

Verificar que Docker Desktop este abierto y en estado "Running", luego ejecutar:

```bash
export PATH="$PATH:/c/Program Files/Docker/Docker/resources/bin"
docker-compose up -d
```

Verificar que el contenedor este corriendo:

```bash
docker ps
```

Debe aparecer `stock_cockroachdb` con estado `Up`.

---

## Paso 4 - Iniciar el servidor

Abrir una terminal y ejecutar:

```bash
cd backend
go mod download
go run cmd/servidor/main.go
```

Debe aparecer: `Servidor iniciado en http://localhost:8080`

Dejar esta terminal abierta.

---

## Paso 5 - Iniciar el frontend

Abrir otra terminal y ejecutar:

```bash
cd frontend
npm install
npm run dev
```

Debe aparecer: `Local: http://localhost:5173/`

Dejar esta terminal abierta.

---

## Paso 6 - Abrir la aplicacion y cargar datos

1. Abrir el navegador en: `http://localhost:5173`
2. Hacer clic en el boton **"Sincronizar"** para cargar las acciones desde Alpha Vantage
3. Esperar la confirmacion de sincronizacion exitosa

Listo. El sistema esta funcionando.

---

## Referencia de comandos

| Accion | Comando |
|--------|---------|
| Levantar base de datos | `docker-compose up -d` |
| Detener base de datos | `docker-compose down` |
| Iniciar servidor | `cd backend && go run cmd/servidor/main.go` |
| Iniciar frontend | `cd frontend && npm run dev` |
| Sincronizar datos (desde terminal) | `curl -X POST http://localhost:8080/api/sincronizar` |
| Verificar que el servidor esta activo | `curl http://localhost:8080/salud` |
| Ejecutar pruebas Go | `cd backend && go test ./...` |
| Ejecutar pruebas Vue | `cd frontend && npm run test:unit` |

---

## Puertos utilizados

| Puerto | Servicio |
|--------|----------|
| 8080 | Servidor Go (API REST) |
| 5173 | Frontend Vue (interfaz de usuario) |
| 26257 | CockroachDB (base de datos) |
| 8081 | Panel de administracion de CockroachDB |

---

## Si algo no funciona

Consultar la seccion **"Solucion de problemas"** en el `README.md` para errores comunes y sus soluciones.
