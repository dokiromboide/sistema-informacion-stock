# Documentacion de la API REST

Base URL: `http://localhost:8080`

## Endpoints

### Salud del servidor

```
GET /salud
```

Respuesta:
```json
{ "estado": "activo" }
```

---

### Acciones

#### Listar acciones

```
GET /api/acciones
```

Parametros de query opcionales:

| Parametro   | Tipo   | Descripcion                                      | Default |
|-------------|--------|--------------------------------------------------|---------|
| q           | string | Texto de busqueda (nombre o simbolo)             | -       |
| ordenar_por | string | Campo: `simbolo`, `precio`, `cambio_pct`, `volumen` | simbolo |
| direccion   | string | `asc` o `desc`                                  | asc     |
| pagina      | int    | Numero de pagina                                 | 1       |
| por_pagina  | int    | Items por pagina (max 100)                       | 20      |

Respuesta:
```json
{
  "datos": [
    {
      "id": 1,
      "simbolo": "AAPL",
      "nombre": "Apple Inc.",
      "precio": 189.50,
      "apertura": 188.00,
      "maximo": 190.20,
      "minimo": 187.50,
      "volumen": 52000000,
      "cambio_monto": 1.50,
      "cambio_pct": 0.79,
      "ultima_actualizacion": "2024-01-15T14:30:00Z"
    }
  ],
  "pagina": 1
}
```

#### Detalle de una accion

```
GET /api/acciones/:simbolo
```

Ejemplo: `GET /api/acciones/AAPL`

Respuesta: objeto `Accion` (ver estructura arriba)

#### Buscar acciones

```
GET /api/acciones/buscar?q=apple
```

Busca en la API externa por nombre o ticker.

Respuesta:
```json
{
  "datos": [
    { "simbolo": "AAPL", "nombre": "Apple Inc." }
  ]
}
```

#### Sincronizar datos

```
POST /api/sincronizar
```

Descarga los simbolos populares desde la API externa y los persiste en la base de datos.

Respuesta:
```json
{
  "mensaje": "Sincronizacion completada",
  "sincronizados": 10
}
```

---

### Recomendaciones

```
GET /api/recomendaciones
```

Retorna las mejores acciones para invertir segun el algoritmo de analisis de momentum y volumen.

Respuesta:
```json
{
  "datos": [
    {
      "accion": { "simbolo": "NVDA", "precio": 875.00, ... },
      "puntaje": 92.5,
      "razon": "Momentum positivo fuerte con volumen elevado"
    }
  ]
}
```

## Codigos de error

| Codigo | Descripcion |
|--------|-------------|
| 400    | Parametros invalidos o faltantes |
| 404    | Recurso no encontrado |
| 500    | Error interno del servidor |
