# Arem Shop API

API REST multi-tenant pour la gestion de boutiques d'electronique.

## 1) But du projet

Cette API permet de gerer plusieurs shops totalement isoles entre eux:

- Gestion des utilisateurs internes (`SuperAdmin`, `Admin`) par shop
- Gestion des produits par shop
- Gestion des transactions financieres avec impact stock
- Exposition publique du catalogue d'un shop via URL publique
- Dashboard financier pour `SuperAdmin`

Le point critique est l'isolation multi-tenant stricte:

- Cote prive: toutes les requetes utilisent `shopID` depuis le JWT
- Cote public: toutes les requetes utilisent `:shopID` dans l'URL

## 2) Schema comprehensif Mermaid (HD)

Le schema complet est disponible ici:

- `docs/architecture.md`

Ce document contient:

- Vue globale architecture (clients, middleware, handlers, services, repositories, DB)
- Pipeline de securite des routes privees
- Sequence detaillee de `POST /transactions` (verrouillage stock atomique)
- Sequence detaillee de `GET /public/:shopID/products`
- ERD complet des tables
- Matrice d'acces par role

## 3) Stack technique

- Go 1.18
- Gin (`github.com/gin-gonic/gin`)
- PostgreSQL
- GORM (`gorm.io/gorm`)
- JWT (`github.com/golang-jwt/jwt/v4`)
- Bcrypt (`golang.org/x/crypto/bcrypt`)
- Validation via tags Gin binding + validator

## 4) Structure du projet

```text
cmd/api/main.go
config/.env.example
internal/config
internal/database
internal/models
internal/dto
internal/repository
internal/services
internal/handlers
internal/middleware
internal/utils
migrations
docs
```

## 5) Prerequis

- Go >= 1.18
- PostgreSQL >= 13
- Un utilisateur PostgreSQL avec droits de creation DB/schema

## 6) Setup local

### 6.0 Docker (recommandé — zéro installation)

```bash
cp .env.example .env        # Créer la configuration locale
docker-compose up --build    # Build + launch (API + PostgreSQL + Frontend)
```

API is available on `http://localhost:8080` and frontend on `http://localhost:3000`.

```bash
curl http://localhost:8080/health   # Vérifier que tout fonctionne
```

Pour arrêter :

```bash
docker-compose down       # Arrêter les conteneurs
docker-compose down -v    # Arrêter + supprimer la base de données
```

### 6.1 Script one-shot (setup local sans Docker)

Le script `scripts/first_run.sh` automatise:

- creation de `.env` si absent (avec generation auto d'un `JWT_SECRET` fort)
- verification connexion PostgreSQL
- creation de la DB si absente
- application de la migration initiale si schema absent
- creation du premier shop + premier SuperAdmin si aucun SuperAdmin n'existe
- lancement de l'API

Commande:

```bash
bash scripts/first_run.sh
```

Mode preparation uniquement (sans demarrer l'API):

```bash
bash scripts/first_run.sh --no-run
```

Mode non-interactif (seed explicite):

```bash
FIRST_SHOP_NAME="Shop Demo" \
FIRST_SHOP_WHATSAPP="+212600000000" \
FIRST_SUPERADMIN_NAME="Owner Demo" \
FIRST_SUPERADMIN_EMAIL="owner@shopdemo.com" \
FIRST_SUPERADMIN_PASSWORD="ChangeMe123!" \
bash scripts/first_run.sh
```

### 6.1 Copier la configuration

```bash
cp config/.env.example .env
```

### 6.2 Creer la base PostgreSQL

```bash
createdb arem_shop
```

ou via `psql`:

```sql
CREATE DATABASE arem_shop;
```

### 6.3 Appliquer les migrations

```bash
psql -d arem_shop -f migrations/000001_init.up.sql
```

### 6.4 Lancer l'API

```bash
go run ./cmd/api
```

L'API ecoute par defaut sur `http://localhost:8080`.

### 6.5 Verifier le healthcheck

```bash
curl -s http://localhost:8080/health
```

## 7) Variables d'environnement

Reference: `config/.env.example`.

- `APP_NAME` nom app
- `APP_ENV` environment (`development`, `production`, ...)
- `APP_PORT` port HTTP
- `CORS_ALLOWED_ORIGINS` liste des origins frontend autorisees (CSV)
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`, `DB_TIMEZONE`
- `JWT_SECRET` secret de signature JWT (obligatoire)
- `JWT_TTL_HOURS` duree de vie token
- `BCRYPT_COST` cout bcrypt
- `LOW_STOCK_THRESHOLD` seuil dashboard low stock
- `FRONTEND_PORT` port service frontend Docker
- `NEXT_PUBLIC_API_BASE_URL` base URL API utilisee par le frontend
- `NEXT_PUBLIC_DATA_MODE` mode frontend (`mock` ou `api`)

## 8) Bootstrap initial (premier shop + premier SuperAdmin)

Important: `POST /auth/register` est protege et requiert deja un token `SuperAdmin`.

Sur une base vide, il faut donc creer un premier shop et un premier user `SuperAdmin` en base.

### 8.1 Creer un hash bcrypt

Exemple commande locale:

```bash
cat >/tmp/hashpass.go <<'GO'
package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    h, _ := bcrypt.GenerateFromPassword([]byte("ChangeMe123!"), 12)
    fmt.Println(string(h))
}
GO
go run /tmp/hashpass.go
```

### 8.2 Inserer shop + superadmin

```sql
INSERT INTO shops (id, name, active, whatsapp_number)
VALUES ('11111111-1111-1111-1111-111111111111', 'Shop Demo', true, '+212600000000');

INSERT INTO users (name, email, password, role, shop_id)
VALUES (
  'Owner Demo',
  'owner@shopdemo.com',
  '<COLLER_LE_HASH_BCRYPT>',
  'SuperAdmin',
  '11111111-1111-1111-1111-111111111111'
);
```

Ensuite, connectez-vous avec `/auth/login`.

## 9) Authentification JWT

Header attendu:

```http
Authorization: Bearer <token>
```

Claims JWT emises:

```json
{
  "userID": "uuid",
  "email": "user@shop.com",
  "role": "SuperAdmin|Admin",
  "shopID": "uuid-du-shop",
  "exp": 1712345678
}
```

## 10) Regles de securite et multi-tenant

### 10.1 Middleware prive

Routes privees passent par:

1. `AuthMiddleware`
2. `ShopIsolationMiddleware`
3. `RoleMiddleware`

### 10.2 Isolation

- Le `shopID` effectif d'une requete privee vient du JWT
- Toute tentative cross-shop via `:shopID` ou `?shopID` est rejetee (`403`)
- Les repositories filtrent toutes les operations critiques par `shop_id`

### 10.3 Confidentialite des prix d'achat

- `purchasePrice` visible uniquement pour `SuperAdmin`
- Routes publiques: `purchasePrice` jamais retourne

### 10.4 Stock

- `POST /transactions` de type `Sale` decremente le stock dans une transaction SQL atomique
- Verrouillage row-level (`SELECT ... FOR UPDATE`)
- Stock negatif impossible

## 11) Endpoints disponibles

| Methode | Endpoint | Acces | Description |
|---|---|---|---|
| GET | `/health` | Public | Etat API + DB |
| POST | `/auth/login` | Public | Login (email + password + shopID) |
| POST | `/auth/register` | SuperAdmin | Creer un user dans le meme shop |
| GET | `/products` | Admin, SuperAdmin | Lister produits du shop JWT |
| POST | `/products` | Admin, SuperAdmin | Creer produit dans le shop JWT |
| PUT | `/products/:id` | Admin, SuperAdmin | Modifier produit du shop JWT |
| DELETE | `/products/:id` | Admin, SuperAdmin | Supprimer produit du shop JWT |
| POST | `/transactions` | Admin, SuperAdmin | Creer transaction (Sale/Expense/Withdrawal) |
| PUT | `/shop` | SuperAdmin | Modifier nom boutique et numero WhatsApp |
| GET | `/reports/dashboard` | SuperAdmin | Dashboard financier shop JWT |
| GET | `/public/:shopID/products` | Public | Catalogue public du shop |
| POST | `/upload` | Admin, SuperAdmin | Uploader une image (multipart) et recuperer l'url absolue |
| GET | `/uploads/*` | Public | Servir les fichiers statiques des images uploadées |

## 12) Usage API detaille

### 12.1 Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "owner@shopdemo.com",
    "password": "ChangeMe123!",
    "shopID": "11111111-1111-1111-1111-111111111111"
  }'
```

### 12.2 Register (SuperAdmin seulement)

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Authorization: Bearer <TOKEN_SUPERADMIN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Admin Shop",
    "email": "admin@shopdemo.com",
    "password": "StrongPass123",
    "role": "Admin",
    "shopID": "11111111-1111-1111-1111-111111111111"
  }'
```

### 12.3 Produits prives

### GET /products

```bash
curl -H "Authorization: Bearer <TOKEN>" \
  http://localhost:8080/products
```

Exemple reponse:

```json
{
  "success": true,
  "data": [
    {
      "id": "6ac4f1ba-28df-4d7e-b0d7-19348d1e95a6",
      "name": "iPhone 14",
      "description": "128GB",
      "category": "Smartphones",
      "sellingPrice": 5999.99,
      "stock": 5,
      "imageURL": "https://example.com/iphone14.jpg",
      "shopID": "11111111-1111-1111-1111-111111111111",
      "createdAt": "2026-02-19T10:00:00Z"
    }
  ]
}
```

Regles role-based:

- `SuperAdmin` voit `purchasePrice`
- `Admin` ne voit jamais `purchasePrice`

### POST /products

- Cas `SuperAdmin`: `purchasePrice` requis et strictement `> 0`
- Cas `Admin`: `purchasePrice` interdit dans le payload; si omis, la valeur interne est forcee a `0`
- Validation commune: `sellingPrice > 0` et `stock >= 0`

Exemple:

```bash
curl -X POST http://localhost:8080/products \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 14",
    "description": "128GB",
    "category": "Smartphones",
    "purchasePrice": 4500.00,
    "sellingPrice": 5999.99,
    "stock": 5,
    "imageURL": "https://example.com/iphone14.jpg"
  }'
```

Exemple erreur (`Admin` qui envoie `purchasePrice`):

```json
{
  "success": false,
  "error": "purchasePrice is not allowed for Admin"
}
```

### PUT /products/:id

```bash
curl -X PUT http://localhost:8080/products/<PRODUCT_ID> \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 14",
    "description": "128GB - Updated",
    "category": "Smartphones",
    "purchasePrice": 4400.00,
    "sellingPrice": 5899.99,
    "stock": 7,
    "imageURL": "https://example.com/iphone14-new.jpg"
  }'
```

Regles:

- `Admin` ne peut pas fournir `purchasePrice`
- `SuperAdmin` peut mettre a jour `purchasePrice`, avec contrainte `> 0` si fourni
- Toujours isole par `shopID` issu du JWT (jamais depuis le body)

### DELETE /products/:id

```bash
curl -X DELETE http://localhost:8080/products/<PRODUCT_ID> \
  -H "Authorization: Bearer <TOKEN>"
```

Exemple reponse:

```json
{
  "success": true,
  "data": {
    "message": "product deleted"
  }
}
```

### 12.4 Transactions

### POST /transactions (Sale)

```bash
curl -X POST http://localhost:8080/transactions \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "Sale",
    "productID": "<PRODUCT_ID>",
    "quantity": 2,
    "amount": 11999.98
  }'
```

### POST /transactions (Expense)

```bash
curl -X POST http://localhost:8080/transactions \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "Expense",
    "amount": 300.00
  }'
```

Regles:

- `Sale`: `productID` obligatoire, `quantity > 0`
- `Expense`/`Withdrawal`: `productID` interdit, `quantity = 0`
- `amount > 0`

### 12.5 Dashboard (SuperAdmin)

```bash
curl -H "Authorization: Bearer <TOKEN_SUPERADMIN>" \
  http://localhost:8080/reports/dashboard
```

Exemple reponse:

```json
{
  "totalSales": 50000,
  "totalExpenses": 30000,
  "netProfit": 20000,
  "lowStockProducts": 3,
  "shopID": "11111111-1111-1111-1111-111111111111"
}
```

### 12.6 Catalogue public

```bash
curl http://localhost:8080/public/11111111-1111-1111-1111-111111111111/products
```

Exemple reponse:

```json
[
  {
    "id": "...",
    "name": "iPhone 14",
    "description": "128GB",
    "category": "Smartphones",
    "sellingPrice": 5999.99,
    "stock": 5,
    "imageURL": "https://example.com/iphone14.jpg",
    "whatsappLink": "https://wa.me/212600000000?text=Bonjour%20je%20veux%20plus%20d%27information%20sur%20iPhone%2014"
  }
]
```

## 13) Format des erreurs

Format standard:

```json
{
  "error": "message",
  "code": "ERROR_CODE",
  "info": ["details optionnels"]
}
```

Exemples de codes frequents:

- `INVALID_CREDENTIALS`
- `CROSS_SHOP_FORBIDDEN`
- `FORBIDDEN`
- `PRODUCT_NOT_FOUND`
- `INSUFFICIENT_STOCK`
- `SHOP_INACTIVE`

## 14) Choses a eviter

- Ne jamais appeler `POST /auth/register` sans token `SuperAdmin`
- Ne jamais oublier `shopID` dans `/auth/login`
- Ne jamais faire de requete metier sans filtre `shop_id` en repository
- Ne jamais exposer `purchasePrice` sur une route publique ou pour role `Admin`
- Ne jamais decremeter le stock hors transaction SQL pour les ventes
- Ne jamais accepter `Sale` sans `productID` et `quantity > 0`
- Ne jamais accepter stock ou montants negatifs
- Ne jamais utiliser un `JWT_SECRET` faible en production

## 15) Tests et validation

### 15.1 Tests code (rapides)

```bash
# Depuis la racine du projet
make test
make build
```

Equivalent sans Makefile:

```bash
go test ./...
go build ./...
```

Tests unitaires deja inclus:

- Isolation middleware (`shopID` match/mismatch)
- Generation lien WhatsApp
- Regles produits (`purchasePrice` role-based)
- Regles transactions (stock insuffisant, sale, expense)

### 15.2 Smoke test API (end-to-end avec curl)

Prerequis:

- API lancee (`go run ./cmd/api`)
- Base migree (`migrations/000001_init.up.sql`)
- Shop + SuperAdmin deja crees (section 8)
- `jq` installe pour parser les reponses JSON

Commandes:

```bash
set -euo pipefail

export BASE_URL="http://localhost:8080"
export SHOP_ID="11111111-1111-1111-1111-111111111111"
export SUPERADMIN_EMAIL="owner@shopdemo.com"
export SUPERADMIN_PASSWORD="ChangeMe123!"

echo "1) Healthcheck"
curl -s "$BASE_URL/health" | jq

echo "2) Login SuperAdmin"
TOKEN=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$SUPERADMIN_EMAIL\",
    \"password\": \"$SUPERADMIN_PASSWORD\",
    \"shopID\": \"$SHOP_ID\"
  }" | jq -r '.token')

test -n "$TOKEN" && test "$TOKEN" != "null"
echo "TOKEN obtenu"

echo "3) Create product (SuperAdmin)"
PRODUCT_ID=$(curl -s -X POST "$BASE_URL/products" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Produit Test",
    "description": "Produit pour smoke test",
    "category": "Tests",
    "purchasePrice": 100.00,
    "sellingPrice": 150.00,
    "stock": 10,
    "imageURL": "https://example.com/test.jpg"
  }' | jq -r '.id')

test -n "$PRODUCT_ID" && test "$PRODUCT_ID" != "null"
echo "PRODUCT_ID=$PRODUCT_ID"

echo "4) List private products"
curl -s "$BASE_URL/products" \
  -H "Authorization: Bearer $TOKEN" | jq

echo "5) Public products by shop"
curl -s "$BASE_URL/public/$SHOP_ID/products" | jq

echo "6) Create Sale transaction (stock decrement)"
curl -s -X POST "$BASE_URL/transactions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"type\": \"Sale\",
    \"productID\": \"$PRODUCT_ID\",
    \"quantity\": 2,
    \"amount\": 300.00
  }" | jq

echo "7) Dashboard SuperAdmin"
curl -s "$BASE_URL/reports/dashboard" \
  -H "Authorization: Bearer $TOKEN" | jq

echo "8) Delete product"
curl -s -o /dev/null -w "%{http_code}\n" -X DELETE "$BASE_URL/products/$PRODUCT_ID" \
  -H "Authorization: Bearer $TOKEN"
```

### 15.3 Verification de securite (cross-shop)

Le but est de verifier qu'un SuperAdmin ne peut pas creer un user dans un autre shop:

```bash
curl -i -X POST "$BASE_URL/auth/register" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Cross Shop Test",
    "email": "cross-shop@test.com",
    "password": "StrongPass123",
    "role": "Admin",
    "shopID": "22222222-2222-2222-2222-222222222222"
  }'
```

Resultat attendu: `403 Forbidden` avec code `CROSS_SHOP_FORBIDDEN`.

## 16) Limitations actuelles

Le coeur metier principal est present, mais certains modules restent a etendre:

- CRUD complet users
- CRUD complet transactions (aujourd'hui `POST` uniquement)
- OpenAPI complet (le fichier actuel est minimal)
- Hardening infra: rate limiting, CORS strict, logs structures, CI/CD

## 17) License

Usage interne projet.

## 18) Frontend (Person 4)

The project now includes a Next.js frontend in `frontend/` with a Rackoon-like visual direction and Arem temporary branding.

### 18.1 Pages

- `/login`
- `/dashboard` (SuperAdmin financial view)
- `/products`
- `/products/new`
- `/products/:id/edit`
- `/transactions/new`
- `/profile` (SuperAdmin — modification du nom de la boutique et du numero WhatsApp)
- `/public/:shopID`

### 18.2 Local frontend run (without Docker)

```bash
cd frontend
npm install
npm run dev
```

Frontend is available at `http://localhost:3000`.

### 18.3 Docker run (API + DB + Frontend)

```bash
cp .env.example .env
docker compose up --build
```

- API: `http://localhost:8080`
- Frontend: `http://localhost:3000`

### 18.4 CORS and API connection

- Backend reads `CORS_ALLOWED_ORIGINS` (default: `http://localhost:3000`).
- Frontend reads:
  - `NEXT_PUBLIC_API_BASE_URL`
  - `NEXT_PUBLIC_DATA_MODE=mock|api` (default recommended: `api`)

### 18.5 Auth and token flow

- Login stores JWT in `localStorage` (`arem_token`) and user info in `arem_user`.
- Private pages redirect to `/login` if token/session is missing.
- In `api` mode, frontend adds `Authorization: Bearer <token>` on protected calls.
- On `401`, frontend clears local session and forces re-login.

### 18.6 Product visibility rules in UI

- `Admin`: purchase price is hidden and cannot be submitted in forms.
- `SuperAdmin`: purchase price is visible and editable.

### 18.7 Public catalog behavior

- Public storefront uses `/public/:shopID/products`.
- No authentication required.
- Purchase price is never shown.

## 19) Comptes de test et liens

### 19.1 URLs de l'application

| Service | URL |
|---------|-----|
| Frontend | `http://localhost:3000` |
| API | `http://localhost:8080` |
| Health check | `http://localhost:8080/health` |
| Vitrine publique Shop 1 | `http://localhost:3000/public/11111111-1111-1111-1111-111111111111` |
| Vitrine publique Shop 2 | `http://localhost:3000/public/22222222-2222-2222-2222-222222222222` |

### 19.2 Comptes de connexion

| Shop | Email | Mot de passe | Shop ID |
|------|-------|--------------|---------|
| Shop 1 (Shop Demo) | `owner@shopdemo.com` | `ChangeMe123!` | `11111111-1111-1111-1111-111111111111` |
| Shop 2 | `owner2@shopdemo.com` | `Password456!` | `22222222-2222-2222-2222-222222222222` |

### 19.3 Tester les features (apres connexion)

| Feature | Ou tester | Description |
|---------|-----------|-------------|
| Dashboard financier | `http://localhost:3000/dashboard` | Affiche Total Sales, Expenses, Net Profit, Low Stock |
| Catalogue produits | `http://localhost:3000/products` | Liste les produits de la boutique |
| Creer un produit | `http://localhost:3000/products/new` | Formulaire avec nom, prix, stock, image |
| Nouvelle transaction | `http://localhost:3000/transactions/new` | Vente (decremente stock) ou depense |
| Profil boutique | `http://localhost:3000/profile` | Modifier le nom de la boutique et le numero WhatsApp |
| Vitrine publique | `http://localhost:3000/public/11111111-1111-1111-1111-111111111111` | Catalogue public avec bouton WhatsApp |
