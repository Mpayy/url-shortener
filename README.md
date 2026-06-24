![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![MySQL](https://img.shields.io/badge/MySQL-8.0+-4479A1?style=for-the-badge&logo=mysql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-7+-DC382D?style=for-the-badge&logo=redis&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

# 🔗 URL Shortener API

A production-grade REST API for URL shortening built with Go — featuring JWT authentication with Redis-backed session management, cache-aside pattern for sub-millisecond redirects, and async analytics tracking. Designed with **Clean Architecture** and **modular monolith** principles to demonstrate real-world backend engineering practices.

> This is not a tutorial project. Every architectural decision — from compile-time dependency injection to transaction propagation via context — was made deliberately to mirror how production Go services are built at scale.

---

## 📋 Table of Contents

- [Architecture Overview](#-architecture-overview)
- [Tech Stack](#-tech-stack)
- [Redirect Flow — Redis Cache-Aside Pattern](#-redirect-flow--redis-cache-aside-pattern)
- [Getting Started](#-getting-started)
- [API Endpoints](#-api-endpoints)
- [Request & Response Examples](#-request--response-examples)
- [Project Structure](#-project-structure)
- [Key Design Decisions](#-key-design-decisions)
- [Future Improvements](#-future-improvements)

---

## 🏗 Architecture Overview

The project follows **Clean Architecture** with clear layer separation and dependency inversion:

```
┌─────────────────────────────────────────────────────────────┐
│                      HTTP Layer (Gin)                       │
│            Controllers · Middleware · Routes                │
├─────────────────────────────────────────────────────────────┤
│                      Usecase Layer                          │
│          Business Logic · Orchestration · Caching           │
├─────────────────────────────────────────────────────────────┤
│                     Repository Layer                        │
│           Database Access · Query Building                  │
├─────────────────────────────────────────────────────────────┤
│                      Entity Layer                           │
│              Domain Models · Table Mapping                  │
└─────────────────────────────────────────────────────────────┘
```

Each layer communicates through **interfaces**, making the codebase testable and loosely coupled. Dependencies are resolved at compile time via **Google Wire** — wiring errors are caught during build, not at runtime.

The project is designed as a **modular monolith**: modules (auth, url) are self-contained and do not share database tables or foreign keys. This makes it straightforward to extract modules into separate services if needed in the future.

---

## 🛠 Tech Stack

| Technology | Purpose | Why This Choice |
|---|---|---|
| **Go 1.22+** | Language | Strong concurrency, fast compilation, ideal for API services |
| **Gin** | HTTP Framework | Lightweight, high-performance router with middleware support |
| **GORM** | ORM | Mature Go ORM with migration, hooks, and transaction support |
| **MySQL** | Database | Battle-tested relational DB, widely used in production |
| **Redis** | Cache + Session Store | In-memory store for fast URL lookups and JWT session management |
| **Google Wire** | Dependency Injection | Compile-time DI — catches wiring errors at build time, not runtime |
| **Logrus** | Logging | Structured JSON logging with configurable log levels |
| **Viper** | Configuration | Reads from `.env` files with type-safe access |
| **golang-migrate** | Migration | Versioned SQL migrations, clean up/down support |
| **JWT (golang-jwt)** | Authentication | Stateless token generation with Redis-backed session validation |
| **bcrypt** | Password Hashing | Industry-standard adaptive hashing with configurable cost |
| **Swagger** | API Documentation | Interactive API docs (planned) |

---

## 🔄 Redirect Flow — Redis Cache-Aside Pattern

The redirect endpoint (`GET /:short_code`) uses a **cache-aside pattern** with async hit counting to ensure redirects are as fast as possible:

```mermaid
flowchart TD
    A["Client: GET /:short_code"] --> B{"Redis Cache\nLookup"}
    B -- "Cache HIT ✅" --> C["Return Original URL"]
    B -- "Cache MISS ❌" --> D["Query MySQL"]
    D -- "Not Found" --> E["Return 404"]
    D -- "Found" --> F["Store in Redis\n(TTL: 24 hours)"]
    F --> C
    C --> G["🔀 HTTP 302 Redirect"]
    C --> H["Async: Increment\nHits Counter"]

    style B fill:#DC382D,color:#fff
    style D fill:#4479A1,color:#fff
    style G fill:#00ADD8,color:#fff
    style H fill:#FF9800,color:#fff
```

**Why this matters:**
- 🚀 **Cache hit** → Redis lookup only, no database query (~0.5ms)
- 💾 **Cache miss** → Query DB once, then cache for 24 hours
- 📊 **Hit counting** → Runs in a goroutine, never blocks the redirect response

---

## 🚀 Getting Started

### Prerequisites

- Go 1.22 or higher
- MySQL 8.0+
- Redis 7+
- [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

### Installation

```bash
# 1. Clone the repository
git clone https://github.com/your-username/url-shortener.git
cd url-shortener

# 2. Copy and configure environment variables
cp config.env.example config.env
```

Edit `config.env` with your credentials:

```env
DATABASE_HOST=localhost
DATABASE_PORT=3306
DATABASE_NAME=url_shortener_db
DATABASE_USERNAME=root
DATABASE_PASSWORD=your_password

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DB=0

APP_HOST=localhost
APP_PORT=8080

LOG_LEVEL=debug
JWT_SECRET_KEY=your-secret-key-minimum-32-characters!!
```

> ⚠️ `JWT_SECRET_KEY` must be at least **32 characters**. The application will panic on startup if this requirement is not met.

```bash
# 3. Create the database
mysql -u root -p -e "CREATE DATABASE url_shortener_db;"

# 4. Run database migrations
migrate -path db/migration -database "mysql://root:your_password@tcp(localhost:3306)/url_shortener_db" up

# 5. Install dependencies
go mod download

# 6. Run the server
go run cmd/web/main.go
```

The server will start at `http://localhost:8080`.

---

## 📡 API Endpoints

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/api/v1/auth/register` | ❌ | Register a new user |
| `POST` | `/api/v1/auth/login` | ❌ | Login and receive JWT token |
| `DELETE` | `/api/v1/auth/logout` | ✅ | Revoke current session |
| `POST` | `/api/v1/urls` | ✅ | Create a shortened URL |
| `GET` | `/api/v1/urls` | ✅ | List all URLs owned by user |
| `DELETE` | `/api/v1/urls/:short_code` | ✅ | Delete a shortened URL |
| `GET` | `/:short_code` | ❌ | Redirect to original URL |

> Protected endpoints require `Authorization: Bearer <token>` header.

---

## 📝 Request & Response Examples

### Register

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "rifai",
    "email": "rifai@example.com",
    "password": "securepassword123"
  }'
```

**201 Created**
```json
{
  "data": {
    "id": 1,
    "username": "rifai",
    "email": "rifai@example.com"
  }
}
```

**400 Bad Request** — Validation Error
```json
{
  "error": {
    "email": "must be a valid email",
    "password": "must be at least 8 characters long"
  }
}
```

**409 Conflict** — Duplicate
```json
{
  "error": "email already exists"
}
```

---

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "rifai@example.com",
    "password": "securepassword123"
  }'
```

**200 OK**
```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6..."
  }
}
```

**401 Unauthorized** — Wrong credentials
```json
{
  "error": "unauthorized"
}
```

---

### Logout

```bash
curl -X DELETE http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6..."
```

**200 OK**
```json
{
  "data": true
}
```

---

### Create Short URL

```bash
curl -X POST http://localhost:8080/api/v1/urls \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6..." \
  -H "Content-Type: application/json" \
  -d '{
    "original_url": "https://github.com/golang/go"
  }'
```

**201 Created**
```json
{
  "data": {
    "id": 1,
    "short_code": "aB3xYz",
    "original_url": "https://github.com/golang/go",
    "hits": 0
  }
}
```

---

### List My URLs

```bash
curl -X GET http://localhost:8080/api/v1/urls \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6..."
```

**200 OK**
```json
{
  "data": [
    {
      "id": 1,
      "short_code": "aB3xYz",
      "original_url": "https://github.com/golang/go",
      "hits": 42
    },
    {
      "id": 2,
      "short_code": "kL9mNp",
      "original_url": "https://go.dev/doc/effective_go",
      "hits": 7
    }
  ]
}
```

---

### Delete URL

```bash
curl -X DELETE http://localhost:8080/api/v1/urls/aB3xYz \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6..."
```

**200 OK**
```json
{
  "data": true
}
```

**404 Not Found**
```json
{
  "error": "data not found"
}
```

---

### Redirect

```bash
curl -L http://localhost:8080/aB3xYz
# → 302 Redirect to https://github.com/golang/go
```

---

## 📁 Project Structure

```
url-shortener/
├── cmd/
│   └── web/
│       ├── main.go              # Entry point with graceful shutdown
│       ├── application.go       # Application container struct
│       ├── injector.go          # Wire dependency injection definitions
│       └── wire_gen.go          # Wire generated code (auto)
│
├── db/
│   └── migration/               # Versioned SQL migration files
│       ├── ..._create_users_table.up.sql
│       ├── ..._create_users_table.down.sql
│       ├── ..._create_urls_table.up.sql
│       └── ..._create_urls_table.down.sql
│
├── internal/
│   ├── config/                   # Infrastructure setup
│   │   ├── app.go               # App container (Gin, DB, Logger, Config)
│   │   ├── gin.go               # Gin engine with logging middleware
│   │   ├── gorm.go              # MySQL connection + pool configuration
│   │   ├── logrus.go            # Structured JSON logger setup
│   │   ├── redis.go             # Redis client with cache operations
│   │   ├── validator.go         # Request validator instance
│   │   └── viper.go             # Environment config loader
│   │
│   ├── entity/                   # Database models (GORM)
│   │   ├── url.go               # URL entity
│   │   └── user.go              # User entity
│   │
│   ├── model/                    # DTOs (Request/Response)
│   │   ├── auth.go              # Auth context model
│   │   ├── url_model.go         # URL request & response
│   │   ├── user_model.go        # User request & response
│   │   ├── web.go               # Generic response wrappers
│   │   └── converter/           # Entity ↔ DTO mappers
│   │       ├── url_converter.go
│   │       └── user_converter.go
│   │
│   ├── repository/               # Data access layer
│   │   ├── url_repository.go    # URL repository interface
│   │   ├── url_repository_impl.go
│   │   ├── user_repository.go   # User repository interface
│   │   └── user_repository_impl.go
│   │
│   ├── usecase/                  # Business logic
│   │   ├── auth_usecase.go      # Auth usecase interface
│   │   ├── auth_usecase_impl.go # Register, Login, Logout logic
│   │   ├── url_usecase.go       # URL usecase interface
│   │   └── url_usecase_impl.go  # CRUD, Redirect, Cache logic
│   │
│   ├── delivery/http/            # HTTP transport layer
│   │   ├── auth_controller.go   # Auth controller interface
│   │   ├── auth_controller_impl.go
│   │   ├── url_controller.go    # URL controller interface
│   │   ├── url_controller_impl.go
│   │   ├── middleware/
│   │   │   └── auth_middleware.go # JWT + Redis session validation
│   │   └── route/
│   │       └── route.go         # Route definitions (public/protected)
│   │
│   ├── exception/                # Sentinel errors
│   │   └── error.go             # Domain errors + validation translator
│   │
│   └── util/                     # Shared utilities
│       ├── response.go          # HTTP response helpers (generic)
│       ├── shortcode_util.go    # Crypto-random shortcode generator
│       ├── token_util.go        # JWT creation & parsing
│       └── transaction.go       # DB transaction via context pattern
│
├── config.env.example            # Environment template
├── go.mod
└── go.sum
```

---

## 💡 Key Design Decisions

### Why modular monolith?

Tables are intentionally **not linked with foreign keys**. Each module (auth, url) owns its data independently. This mirrors a modular monolith approach where modules can be extracted into microservices without schema coupling.

### Why Redis for both caching and sessions?

- **Sessions**: Storing JWT tokens in Redis enables **token revocation** (logout). Pure JWT is stateless and can't be invalidated without a blocklist — Redis serves as that blocklist.
- **URL caching**: Redirect is the hottest path. Redis cache-aside with 24h TTL means most redirects never hit the database.

### Why compile-time DI (Wire) over runtime DI?

Runtime DI frameworks (like Uber's `fx`) discover wiring errors at application startup. Wire catches them at **compile time** — if your dependencies don't line up, `go generate` fails immediately. For a production service, this is a significant reliability win.

### Why async hit counting?

The `IncrementHits` call runs in a goroutine with its own context and timeout. A slow database write never delays the user's redirect. This is a conscious trade-off: hits may be slightly delayed, but redirect latency is guaranteed.

---

## 🔮 Future Improvements

- [ ] **Unit & Integration Tests** — Usecase-layer tests with mocked repositories, integration tests with testcontainers
- [ ] **Swagger/OpenAPI Documentation** — Interactive API docs via `swaggo/swag`
- [ ] **Rate Limiting** — Protect login and URL creation endpoints from abuse
- [ ] **Custom URL Validation** — Enforce HTTPS-only URLs to prevent redirecting to insecure sites
- [ ] **URL Expiration** — TTL-based URL expiry with configurable duration
- [ ] **Click Analytics** — Track redirect metadata (referrer, country, device)
- [ ] **Custom Short Codes** — Allow users to choose their own short code alias
- [ ] **Pagination** — Paginate the URL list endpoint for users with many URLs
- [ ] **Dockerize** — Docker Compose setup with MySQL, Redis, and the API
- [ ] **CI/CD Pipeline** — GitHub Actions for lint, test, and build

---

## 📄 License

This project is licensed under the MIT License.
