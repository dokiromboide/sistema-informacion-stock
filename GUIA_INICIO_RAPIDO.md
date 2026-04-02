# Guia de Inicio Rapido

Esta guia asume que ya tienes instalado: Go 1.22+, Node.js 20+, Docker Desktop y Git.
Si no los tienes instalados, seguir primero las instrucciones del README.md.

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

Abrir `backend/.env` y reemplazar `tu_api_key_aqui` con tu clave de Alpha Vantage.
Obtener clave gratuita en: https://www.alphavantage.co/support/#api-key

---

## Paso 3 - Levantar la base de datos

```bash
export PATH="$PATH:/c/Program Files/Docker/Docker/resources/bin"
docker-compose up -d
```

Esperar a que el contenedor diga "Up" al ejecutar `docker ps`.

---

## Paso 4 - Iniciar el servidor

Abrir una terminal nueva y ejecutar:

```bash
cd backend
go mod download
go run cmd/servidor/main.go
```

Debe aparecer: `Servidor iniciado en http://localhost:8080`

Dejar esta terminal abierta.

---

## Paso 5 - Iniciar el frontend

Abrir otra terminal nueva y ejecutar:

```bash
cd frontend
npm install
npm run dev
```

Debe aparecer: `Local: http://localhost:5173/`

Dejar esta terminal abierta.

---

## Paso 6 - Abrir la aplicacion

Abrir el navegador y entrar a:

```
http://localhost:5173
```

Hacer clic en el boton **"Sincronizar"** para cargar los datos de acciones desde Alpha Vantage.

---

## Comandos de referencia rapida

| Accion | Comando |
|--------|---------|
| Iniciar base de datos | `docker-compose up -d` |
| Detener base de datos | `docker-compose down` |
| Iniciar servidor | `cd backend && go run cmd/servidor/main.go` |
| Iniciar frontend | `cd frontend && npm run dev` |
| Ejecutar pruebas Go | `cd backend && go test ./...` |
| Ejecutar pruebas Vue | `cd frontend && npm run test:unit` |
| Verificar servidor activo | `curl http://localhost:8080/salud` |
| Sincronizar datos manualmente | `curl -X POST http://localhost:8080/api/sincronizar` |

---

## Puertos utilizados

| Puerto | Servicio |
|--------|----------|
| 8080 | Servidor Go (API REST) |
| 5173 | Frontend Vue (interfaz de usuario) |
| 26257 | CockroachDB (base de datos) |
| 8081 | Panel de administracion de CockroachDB |
