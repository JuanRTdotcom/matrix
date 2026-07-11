# 🧮 Reto Técnico — Interseguro · Backend Developer 2026

> Dos APIs (**Go** + **Node.js**) que se hablan por HTTP, un **frontend Svelte** que las usa,
> todo con **Docker** y protegido con **JWT**.

---

## 🌐 Demo en vivo (Render)

> [!IMPORTANT]
> Están en el plan **gratis** de Render, así que se **duermen** cuando nadie las usa.
> **Antes de probar**, abre los 3 enlaces para que "despierten" (tarda ~30–50s la primera vez).

| Servicio | Enlace | Ábrelo primero |
|----------|--------|:--------------:|
| 🟦 **Frontend** (Svelte) | https://front-svelte-2vw2.onrender.com/ | 👉 |
| 🐹 **API Go** (QR) | https://go-3sl1.onrender.com/ | 👉 |
| 🟩 **API Node** (Stats) | https://node-7jwe.onrender.com/ | 👉 |

**🔑 Login demo:** usuario `admin` · contraseña `admin123`

---

## 🗺️ Cómo funciona (de un vistazo)

```
   👤 Usuario
      │  1. login  → recibe 🔑 JWT
      ▼
┌─────────────┐   2. envía matriz A      ┌──────────────┐   3. manda Q y R (+🔑)   ┌───────────────┐
│  🟦 Frontend │ ───────────────────────▶ │  🐹 API Go    │ ───────────────────────▶ │  🟩 API Node   │
│  (SvelteKit)│                          │  QR: A = Q·R │                          │  Estadísticas │
│    :5173    │ ◀─────────────────────── │    :3000     │ ◀─────────────────────── │     :4000     │
└─────────────┘   5. { original, qr,     └──────────────┘   4. devuelve stats      └───────────────┘
                       statistics }
```

1. 👤 El usuario **inicia sesión** → la API Go emite un **🔑 JWT**.
2. 👤 Envía una **matriz** → la API Go calcula su **factorización QR** (`A = Q·R`).
3. 🐹 La API Go manda `Q` y `R` (con el 🔑) a la API Node.
4. 🟩 La API Node calcula las **estadísticas** y las devuelve.
5. 🐹 La API Go responde al frontend con `{ original, qr, statistics }`.

---

## 🧩 Qué hace cada pieza

| | Servicio | Puerto | Trabajo |
|---|----------|:------:|---------|
| 🐹 | **API Go** (Fiber) | `3000` | Emite el JWT · calcula la **factorización QR** · orquesta |
| 🟩 | **API Node** (Express) | `4000` | Calcula **estadísticas** de `Q` y `R` |
| 🟦 | **Frontend** (SvelteKit) | `5173` | UI: login + enviar matriz + ver resultados |

**📐 Factorización QR** → algoritmo *Gram-Schmidt modificado* (más estable), sin librerías externas.
`Q` con columnas ortonormales, `R` triangular superior.

**📊 Estadísticas** → máximo, mínimo, promedio, suma total y si alguna matriz es diagonal.

---

## 🔒 Seguridad (JWT)

- 🔑 `POST /api/v1/auth/login` (Go) emite un JWT firmado con **HS256**.
- 🚧 Las rutas de negocio de **ambas** APIs exigen `Authorization: Bearer <token>`.
- 🔁 La API Go **reenvía el mismo JWT** a la API Node → la comunicación entre servicios
  también va autenticada. Comparten `JWT_SECRET`.
- 👤 Credenciales demo: **admin / admin123** (configurables por variables de entorno).

---

## ▶️ Cómo ejecutarlo en local

### ⚙️ Configuración — un solo `.env`

Toda la config vive en **un `.env` en la raíz**, junto al `docker-compose.yml`.
Compose lo lee solo y lo reparte a los servicios (`${VAR}`).

```bash
cp .env.example .env    # ajusta secretos/puertos si hace falta
```

### 🐳 Levantar todo con Docker

```bash
docker compose up --build
```

| | URL |
|---|-----|
| 🟦 Frontend | http://localhost:5173 |
| 🐹 API Go | http://localhost:3000 |
| 🟩 API Node | http://localhost:4000 |

---

## 🧪 Pruebas

```bash
cd go-api   && go test ./...     # 🐹 unitarias (QR, orquestación, auth) + integración
cd node-api && npm test          # 🟩 unitarias (stats) + integración (JWT, supertest)
```

### ⚡ Prueba rápida con cURL

```bash
# 1. Login → obtener token
TOKEN=$(curl -s -X POST http://localhost:3000/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin123"}' | jq -r .data.token)

# 2. Procesar matriz
curl -s -X POST http://localhost:3000/api/v1/matrix/process \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"matrix":[[12,-51,4],[6,167,-68],[-4,24,-41]]}' | jq
```

---

## 🛣️ Endpoints

| | Método | Ruta | 🔒 Auth | Qué hace |
|---|--------|------|:------:|----------|
| 🐹 Go | `GET`  | `/`                      | — | Healthcheck |
| 🐹 Go | `POST` | `/api/v1/auth/login`     | — | Emite JWT |
| 🐹 Go | `POST` | `/api/v1/matrix/process` | ✅ | Factoriza (QR) y delega stats a Node |
| 🟩 Node | `GET`  | `/`                    | — | Healthcheck |
| 🟩 Node | `POST` | `/api/v1/statistics`   | ✅ | Estadísticas de un set de matrices |

---

## 🏗️ Arquitectura — Clean Architecture

Cada proyecto separa responsabilidades en capas (la UI/controllers nunca tocan detalles de infraestructura):

- **🐹 Go** — `controllers → usecases → repositories → models`. Config validada en `config/`,
  ensamblado (DI + rutas + ciclo de vida) en `server/`. La dependencia externa (API Node) es
  una interfaz (`usecases.StatsRepository`) → inversión de dependencias (SOLID).
- **🟩 Node** — `controllers → usecases → adapters → repository`, con `domain/` (entities + interfaces).
- **🟦 Frontend** — `domain / adapters / usecases / repository / enviroment`; la UI solo invoca casos de uso.

---

## ☁️ Despliegue en Render

Cada servicio genera su **imagen Docker autónoma** → desplegable en cualquier proveedor.
Aquí corre en **Render** (ver enlaces arriba ⬆️).

**Cómo se desplegó:**

1. 🌿 Cada proyecto se separó en su **propia rama** de Git.
2. 📦 Cada rama se desplegó como un **proyecto/servicio independiente** en Render (uno por uno).
3. 🔐 A cada proyecto se le cargaron sus **variables de entorno** — las mismas del `docker-compose.yml`.

> En Render **no hay red interna de Docker**, así que las URLs entre servicios son públicas:
> - `NODE_API_URL` → apunta a la URL pública de la API Node.
> - `PUBLIC_GO_API_URL` → apunta a la URL pública de la API Go.

**Variables por servicio:**

| Servicio | Variables de entorno |
|----------|----------------------|
| 🟩 Node | `PORT` · `JWT_SECRET` · `CORS_ORIGINS` |
| 🐹 Go | `SERVER_PORT` · `JWT_SECRET` · `TOKEN_EXPIRY_MINUTES` · `AUTH_USER` · `AUTH_PASS` · `NODE_API_URL` · `CORS_ORIGINS` |
| 🟦 Frontend | `PORT` · `PUBLIC_GO_API_URL` |

> `JWT_SECRET` debe ser **el mismo** en Go y Node (validan el mismo token).
