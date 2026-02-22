# Arem-Shop — Architecture du Projet

**API REST multi-tenant** de gestion de boutique (produits, transactions, reporting), écrite en **Go** avec le framework **Gin**.

---

## Arborescence

```
arem-shop/
├── cmd/
│   └── api/
│       └── main.go              ← Point d'entrée, DI & routing
├── config/
│   └── .env                     ← Variables d'environnement (gitignored)
├── internal/
│   ├── config/
│   │   └── config.go            ← Chargement de la configuration (.env → struct)
│   ├── database/
│   │   └── postgres.go          ← Connexion PostgreSQL via GORM
│   ├── models/
│   │   ├── shop.go              ← Entité Shop (tenant)
│   │   ├── user.go              ← Entité User + rôles
│   │   ├── product.go           ← Entité Product
│   │   └── transaction.go       ← Entité Transaction (Sale/Expense/Withdrawal)
│   ├── dto/
│   │   ├── auth_dto.go          ← DTOs d'authentification
│   │   ├── product_dto.go       ← DTOs produit (create/update/response)
│   │   ├── public_product_dto.go← DTOs vitrine publique
│   │   ├── report_dto.go        ← DTOs dashboard
│   │   ├── shop_dto.go          ← DTOs mise à jour boutique
│   │   └── transaction_dto.go   ← DTOs transaction
│   ├── repository/
│   │   ├── user_repository.go
│   │   ├── shop_repository.go
│   │   ├── product_repository.go
│   │   └── transaction_repository.go
│   ├── services/
│   │   ├── auth_service.go      ← Inscription & login
│   │   ├── product_service.go   ← CRUD produits
│   │   ├── public_service.go    ← Catalogue public
│   │   ├── shop_service.go      ← Mise à jour boutique (nom, WhatsApp)
│   │   ├── transaction_service.go ← Ventes, dépenses, retraits
│   │   └── report_service.go    ← Dashboard financier
│   ├── handlers/
│   │   ├── auth_handler.go
│   │   ├── product_handler.go
│   │   ├── public_handler.go
│   │   ├── shop_handler.go      ← Mise à jour boutique
│   │   ├── transaction_handler.go
│   │   └── report_handler.go
│   ├── middleware/
│   │   ├── auth_middleware.go           ← Validation JWT
│   │   ├── role_middleware.go           ← Contrôle RBAC
│   │   └── shop_isolation_middleware.go ← Isolation multi-tenant
│   └── utils/
│       ├── jwt.go               ← Génération & parsing JWT
│       ├── password.go          ← Hashing bcrypt
│       ├── response.go          ← Helpers réponse JSON
│       └── whatsapp.go          ← Utilitaire WhatsApp
├── migrations/
│   ├── 000001_init.up.sql       ← Migration initiale
│   └── 000001_init.down.sql     ← Rollback
├── scripts/
│   └── first_run.sh             ← Script de premier lancement
├── docs/
│   └── openapi.yaml             ← Spécification OpenAPI
├── go.mod
├── go.sum
├── Makefile
└── .gitignore
```

---

## Architecture en Couches

```mermaid
graph TB
    Client["🌐 Client HTTP"]
    
    subgraph "Gin Router"
        MW_Auth["🔒 AuthMiddleware<br/>JWT validation"]
        MW_Shop["🏪 ShopIsolationMiddleware<br/>Multi-tenant guard"]
        MW_Role["👤 RoleMiddleware<br/>RBAC check"]
    end
    
    subgraph "Handlers"
        H_Auth["AuthHandler"]
        H_Product["ProductHandler"]
        H_Transaction["TransactionHandler"]
        H_Report["ReportHandler"]
        H_Public["PublicHandler"]
    end
    
    subgraph "Services — Logique métier"
        S_Auth["AuthService"]
        S_Product["ProductService"]
        S_Transaction["TransactionService"]
        S_Report["ReportService"]
        S_Public["PublicService"]
    end
    
    subgraph "Repositories — Accès données"
        R_User["UserRepository"]
        R_Shop["ShopRepository"]
        R_Product["ProductRepository"]
        R_Transaction["TransactionRepository"]
    end
    
    DB[("🐘 PostgreSQL")]
    
    Client --> MW_Auth --> MW_Shop --> MW_Role --> H_Auth & H_Product & H_Transaction & H_Report
    Client -.->|"sans auth"| H_Public
    H_Auth --> S_Auth
    H_Product --> S_Product
    H_Transaction --> S_Transaction
    H_Report --> S_Report
    H_Public --> S_Public
    S_Auth --> R_User & R_Shop
    S_Product --> R_Product
    S_Transaction --> R_Product & R_Transaction
    S_Report --> R_Transaction & R_Product
    S_Public --> R_Shop & R_Product
    R_User & R_Shop & R_Product & R_Transaction --> DB
```

---

## Stack Technique

| Composant       | Technologie                        |
|-----------------|------------------------------------|
| Langage         | **Go 1.18**                        |
| Framework HTTP  | **Gin v1.8**                       |
| ORM             | **GORM v1.23** + driver PostgreSQL |
| Base de données | **PostgreSQL**                     |
| Auth            | **JWT** (golang-jwt/jwt/v4)        |
| Hashing         | **bcrypt** (golang.org/x/crypto)   |
| Décimaux        | **shopspring/decimal**             |
| UUID            | **google/uuid**                    |
| Config          | **godotenv**                       |

---

## Modèle de Données (ERD)

```mermaid
erDiagram
    SHOP {
        uuid id PK
        string name
        bool active
        string whatsAppNumber
        timestamp createdAt
    }
    
    USER {
        uuid id PK
        string name
        string email
        string password
        enum role "SuperAdmin | Admin"
        uuid shopID FK
        timestamp createdAt
    }
    
    PRODUCT {
        uuid id PK
        string name
        string description
        string category
        decimal purchasePrice
        decimal sellingPrice
        int stock
        string imageURL
        uuid shopID FK
        timestamp createdAt
    }
    
    TRANSACTION {
        uuid id PK
        enum type "Sale | Expense | Withdrawal"
        uuid productID FK "nullable"
        int quantity
        decimal amount
        uuid shopID FK
        timestamp createdAt
    }
    
    SHOP ||--o{ USER : "emploie"
    SHOP ||--o{ PRODUCT : "possède"
    SHOP ||--o{ TRANSACTION : "enregistre"
    PRODUCT ||--o{ TRANSACTION : "concerne"
```

---

## Endpoints API

### Publics (sans authentification)
| Méthode | Route                        | Description            |
|---------|------------------------------|------------------------|
| `GET`   | `/health`                    | Health check           |
| `POST`  | `/auth/login`                | Connexion (retourne JWT) |
| `GET`   | `/public/:shopID/products`   | Catalogue public       |
| `GET`   | `/uploads/*`                 | Fichiers images statiques|

### Protégés — SuperAdmin + Admin
| Méthode  | Route              | Description            |
|----------|--------------------|------------------------|
| `GET`    | `/products`        | Lister les produits    |
| `POST`   | `/products`        | Créer un produit       |
| `PUT`    | `/products/:id`    | Modifier un produit    |
| `DELETE` | `/products/:id`    | Supprimer un produit   |
| `POST`   | `/transactions`    | Créer une transaction  |
| `POST`   | `/upload`          | Upload d'une image     |

### Protégés — SuperAdmin uniquement
| Méthode | Route                 | Description                |
|---------|-----------------------|----------------------------|
| `POST`  | `/auth/register`      | Enregistrer un utilisateur |
| `PUT`   | `/shop`               | Modifier nom et WhatsApp   |
| `GET`   | `/reports/dashboard`  | Dashboard financier        |

---

## Sécurité & Multi-Tenancy

```mermaid
sequenceDiagram
    participant C as Client
    participant Auth as AuthMiddleware
    participant Shop as ShopIsolationMiddleware
    participant Role as RoleMiddleware
    participant H as Handler
    
    C->>Auth: Request + Bearer JWT
    Auth->>Auth: Valider & décoder JWT
    Auth->>Shop: Claims injectés dans le contexte
    Shop->>Shop: Vérifier shopID du JWT<br/>vs shopID de la route
    Shop->>Role: shopID validé dans le contexte
    Role->>Role: Vérifier le rôle de l'utilisateur<br/>vs rôles autorisés
    Role->>H: ✅ Accès autorisé
    
    Note over Auth: ❌ 401 si token absent/invalide
    Note over Shop: ❌ 403 si cross-shop détecté
    Note over Role: ❌ 403 si rôle insuffisant
```

### Principes clés
- **JWT** : Chaque token contient `userID`, `shopID`, `role`
- **Isolation multi-tenant** : Un utilisateur ne peut accéder qu'aux données de son shop
- **RBAC** : Deux rôles — `SuperAdmin` (gestion complète) et `Admin` (opérations courantes)
- **Protection des données sensibles** : Le `purchasePrice` n'est visible que par les `SuperAdmin`

---

## Injection de Dépendances

Le point d'entrée `main.go` assemble manuellement les couches (**pas de framework DI**) :

```
Config → Database → Repositories → Services → Handlers → Router
```

Cela garantit un **couplage faible** et une **testabilité maximale** — les services utilisent des **interfaces** pour leurs dépendances repository.

---

## Transactions Base de Données

Les ventes (`TransactionSale`) utilisent une **transaction SQL** avec un verrou pessimiste (`SELECT ... FOR UPDATE`) pour :
1. Verrouiller le produit concerné
2. Vérifier le stock disponible
3. Décrémenter le stock
4. Créer la transaction

Cela empêche les **race conditions** sur le stock.
