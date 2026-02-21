# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Arem Shop is a multi-tenant REST API + Next.js frontend for managing electronics shops. Each shop is isolated at the database level via `shop_id` foreign keys and enforced through middleware. Two roles exist: **SuperAdmin** (full access, sees purchase prices) and **Admin** (limited, no purchase price visibility).

## Common Commands

### Backend (Go)
```bash
make run           # Run API server locally (go run ./cmd/api)
make build         # Compile (go build ./...)
make test          # Run all Go tests (go test ./...)
make fmt           # Format Go code (gofmt -w ./cmd ./internal)
go test ./internal/services/...   # Run tests for a single package
```

### Frontend (Next.js)
```bash
cd frontend
npm run dev        # Dev server on port 3000
npm run build      # Production build
npm test           # Run all tests (vitest run)
npm run test:watch # Watch mode (vitest)
```

### Docker
```bash
make docker-up         # Start all services (postgres, api, frontend)
make docker-up-detach  # Start detached
make docker-down       # Stop services
make docker-clean      # Stop and remove volumes
make docker-test       # Run API smoke tests in Docker
```

### First-time Setup
Copy `.env.example` to `.env`, then run `bash scripts/first_run.sh` to initialize DB, run migrations, and seed a demo shop with a SuperAdmin user.

## Architecture

### Backend Layers (Go + Gin + GORM + PostgreSQL)

```
cmd/api/main.go  →  Entry point, DI wiring, route registration
internal/
  config/        →  Env-based configuration (.env via godotenv)
  database/      →  GORM PostgreSQL connection
  models/        →  Domain entities (User, Product, Transaction, Shop)
  dto/           →  Request/response data transfer objects
  repository/    →  Data access (all queries filtered by shop_id)
  services/      →  Business logic (auth, product, transaction, report, public)
  handlers/      →  HTTP handlers (parse request → call service → JSON response)
  middleware/    →  AuthMiddleware (JWT) → ShopIsolationMiddleware → RoleMiddleware
  utils/         →  JWT helpers, bcrypt, response formatters, WhatsApp link builder
migrations/      →  SQL schema (000001_init) and seed data (000002_seed_demo)
```

Dependency flow: `handlers → services → repositories → GORM/DB`

### Middleware Chain (private routes)
1. **AuthMiddleware** — validates JWT, extracts claims (userID, shopID, role) into Gin context
2. **ShopIsolationMiddleware** — blocks cross-shop access (403 if shopID mismatch)
3. **RoleMiddleware** — checks user role against allowed roles

### Frontend (Next.js 14 App Router + TypeScript + Tailwind)

```
frontend/src/
  app/              →  File-based routing (App Router)
    (private)/      →  Auth-guarded layout for dashboard, products, transactions
    public/[shopID] →  Public shop catalog (no auth required)
  components/       →  Reusable UI (ProductCard, ProductForm, TransactionForm, etc.)
  lib/
    services/api.ts →  HTTP client wrapping fetch, auto-injects JWT from localStorage
    services/adapters.ts → Maps API responses to frontend types
    auth/           →  Auth helpers, guards, localStorage token management
    i18n/           →  I18nProvider context (French + English)
```

Auth tokens stored in localStorage as `arem_token` and `arem_user`. On 401, session is cleared and user redirected to `/login`.

`NEXT_PUBLIC_DATA_MODE` env var switches between `api` (real backend) and `mock` (test data).

### API Routes

| Route | Auth | Description |
|-------|------|-------------|
| `GET /health` | Public | Health check |
| `POST /auth/login` | Public | Login (returns JWT) |
| `POST /auth/register` | SuperAdmin | Create user |
| `GET/POST /products` | Admin+ | List/create products |
| `GET/PUT/DELETE /products/:id` | Admin+ | Single product CRUD |
| `POST /transactions` | Admin+ | Create sale/expense/withdrawal |
| `PUT /shop` | SuperAdmin | Update shop name and WhatsApp |
| `GET /reports/dashboard` | SuperAdmin | Financial dashboard |
| `GET /public/:shopID/products` | Public | Public catalog |

### Multi-tenant Isolation
- Every repository query includes a `WHERE shop_id = ?` clause
- JWT contains `shopID` claim; middleware enforces it on all private routes
- Stock updates use `SELECT ... FOR UPDATE` row-level locking inside GORM transactions

### Database (PostgreSQL 15)
4 tables: `shops`, `users`, `products`, `transactions`. All linked by `shop_id` FK. Users have shop-scoped unique email constraint. Products have composite index on `(shop_id, category)`.

### Environment
Key env vars (see `.env.example`): `APP_PORT` (default 8080), `DB_*` for Postgres, `JWT_SECRET`, `JWT_TTL_HOURS`, `BCRYPT_COST`, `LOW_STOCK_THRESHOLD`, `NEXT_PUBLIC_API_BASE_URL`, `NEXT_PUBLIC_DATA_MODE`.

## Error Response Format

```json
{"error": "message", "code": "ERROR_CODE", "info": ["optional details"]}
```

Common codes: `INVALID_CREDENTIALS`, `CROSS_SHOP_FORBIDDEN`, `INSUFFICIENT_STOCK`, `PRODUCT_NOT_FOUND`.
