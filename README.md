# Reto Técnico — Interseguro · Backend Developer 2026

Solución al Coding Challenge: dos APIs REST (Go + Node.js) que se comunican por HTTP,
un frontend en Svelte que las consume, todo contenerizado con Docker y seguro con JWT.

```
┌────────────┐  JWT   ┌──────────────┐  HTTP + JWT   ┌───────────────┐
│  Frontend  │ ─────▶ │  API en Go   │ ────────────▶ │  API en Node  │
│ (SvelteKit)│        │   (Fiber)    │   Q, R        │   (Express)   │
│   :5173    │ ◀───── │  QR :3000    │ ◀──────────── │  Stats :4000  │
└────────────┘        └──────────────┘  estadísticas └───────────────┘
```

- **API en Go (`go-api`, :3000)** — recibe la matriz, calcula su **factorización QR**,
  emite el JWT y envía `Q` y `R` a la API en Node.js. Es el punto de entrada del sistema.
- **API en Node.js (`node-api`, :4000)** — recibe `Q` y `R` y calcula las estadísticas.

## Flujo funcional

1. El usuario inicia sesión en el frontend → la **API en Go** emite un **JWT**.
2. El usuario envía una matriz → la **API en Go** calcula su **factorización QR** (`A = Q·R`).
3. La **API en Go** envía `Q` y `R` (con el JWT) a la **API en Node.js**.
4. La **API en Node.js** calcula las estadísticas y las devuelve.
5. La **API en Go** responde al frontend con `{ original, qr, statistics }`.

## Factorización QR (API en Go)

Algoritmo **Gram-Schmidt modificado** (numéricamente más estable que el clásico), sin
dependencias externas. `A = Q · R` con `Q` de columnas ortonormales y `R` triangular superior.

## Estadísticas (API en Node.js)

Sobre las matrices `Q` y `R`: **valor máximo**, **valor mínimo**, **promedio**,
**suma total** y si **alguna es diagonal**.

## Seguridad (JWT)

- `POST /api/v1/auth/login` (Go) emite un JWT firmado con HS256.
- Las rutas de negocio de **ambas** APIs exigen `Authorization: Bearer <token>`.
- La API en Go **reenvía el mismo JWT** a la API en Node.js: la comunicación entre
  servicios también está autenticada. Ambas comparten `JWT_SECRET`.
- Credenciales de demo: **admin / admin123** (configurables por variables de entorno).

## Arquitectura — Clean Architecture

**Go** (`go-api/`): `controllers → usecases → repositories → models`. La configuración se
carga y valida en `config/` y el ensamblado (inyección de dependencias + rutas + ciclo de
vida) vive en `server/`. La dependencia externa (API en Node) se define como interfaz
(`usecases.StatsRepository`) e implementa en `repositories/`. Los controladores dependen de
abstracciones (`controllers.Processor`, `Authenticator`) — inversión de dependencias (SOLID).

**Node** (`node-api/src/`): `controllers → usecases → adapters → repository`, con
`domain/` (entities + interfaces). El caso de uso depende de la interfaz del adaptador.

**Frontend** (`frontend/src/lib/`): misma separación — `domain / adapters / usecases /
repository / enviroment` — la UI solo invoca casos de uso.

## Cómo ejecutar

### Configuración (un solo `.env`)

Toda la configuración vive en **un único `.env` en la raíz**, junto al `docker-compose.yml`.
Compose lo lee automáticamente y lo distribuye a los servicios por interpolación (`${VAR}`).

```bash
cp .env.example .env    # ajusta secretos/puertos si hace falta
```

### Opción A — Docker Compose (todo junto)

```bash
docker compose up --build
```

- Frontend:   http://localhost:5173
- API Go:     http://localhost:3000
- API Node:   http://localhost:4000

### Opción B — Local (desarrollo)

`JWT_SECRET` es obligatorio (la API en Go aborta si falta). El resto usa valores por
defecto que apuntan a `localhost`.

```bash
cd node-api && npm install && npm run dev              # :4000
cd go-api   && JWT_SECRET=dev go run main.go           # :3000
cd frontend && npm install && npm run dev              # :5173
```

## Pruebas

**API en Go** — unitarias (QR, orquestación con mocks, autenticación) + integración (rutas
Fiber end-to-end con un stub de Node):

```bash
cd go-api && go test ./...
```

**API en Node.js** — unitarias (cálculo de estadísticas) + integración (endpoint con JWT,
usando supertest):

```bash
cd node-api && npm test
```

## Prueba manual rápida (cURL)

```bash
# 1. Login
TOKEN=$(curl -s -X POST http://localhost:3000/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin123"}' | jq -r .data.token)

# 2. Procesar matriz
curl -s -X POST http://localhost:3000/api/v1/matrix/process \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"matrix":[[12,-51,4],[6,167,-68],[-4,24,-41]]}' | jq
```

## Endpoints

| Servicio | Método | Ruta                     | Auth | Descripción                                 |
|----------|--------|--------------------------|------|---------------------------------------------|
| Go       | GET    | `/`                      | —    | Healthcheck                                 |
| Go       | POST   | `/api/v1/auth/login`     | —    | Emite JWT                                   |
| Go       | POST   | `/api/v1/matrix/process` | JWT  | Factoriza (QR) y delega stats de Q,R a Node |
| Node     | GET    | `/`                      | —    | Healthcheck                                 |
| Node     | POST   | `/api/v1/statistics`     | JWT  | Estadísticas de un set de matrices          |

## Despliegue en la nube

Cada servicio produce una imagen Docker autónoma, lista para desplegar en cualquier
proveedor (AWS ECS/App Runner, Google Cloud Run, Azure Container Apps, etc.). El
`docker-compose.yml` define la red y las variables; en producción se sustituyen los
secretos (`JWT_SECRET`, credenciales) por gestores de secretos del proveedor.

## Checklist del reto

- [x] API en Go con Fiber (factorización QR de la matriz)
- [x] API en Node.js con Express (estadísticas de Q y R)
- [x] Comunicación entre APIs por HTTP + JWT (Go → Node)
- [x] Docker en cada servicio + `docker-compose`
- [x] Frontend (SvelteKit) que consume las APIs *(opcional)*
- [x] Seguridad con JWT *(opcional)*
- [x] Pruebas unitarias e integración en ambas APIs *(opcional)*
- [x] Clean architecture y documentación
