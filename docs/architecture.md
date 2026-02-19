# Architecture Mermaid (HD)

Ce document donne une vue complete du fonctionnement de l'API multi-tenant.

## 1) Vue globale de l'architecture

```mermaid
flowchart LR
    %% Clients
    subgraph C[Clients]
        SA[SuperAdmin Client]
        AD[Admin Client]
        GU[Guest Client]
    end

    %% API
    subgraph A[Go API - Gin]
        R[Router]

        subgraph RG[Route Groups]
            H[/GET /health/]
            L[/POST /auth/login/]
            R1[/POST /auth/register/]
            P1[/GET /products/]
            P2[/POST /products/]
            P3[/PUT /products/:id/]
            P4[/DELETE /products/:id/]
            T1[/POST /transactions/]
            D1[/GET /reports/dashboard/]
            U1[/GET /public/:shopID/products/]
        end

        subgraph MW[Middlewares]
            M1[AuthMiddleware]
            M2[ShopIsolationMiddleware]
            M3[RoleMiddleware]
        end

        subgraph HN[Handlers]
            AH[AuthHandler]
            PH[ProductHandler]
            TH[TransactionHandler]
            RH[ReportHandler]
            UH[PublicHandler]
        end

        subgraph SV[Services]
            AS[AuthService]
            PS[ProductService]
            TS[TransactionService]
            RS[ReportService]
            US[PublicService]
        end

        subgraph RP[Repositories]
            UR[UserRepository]
            SR[ShopRepository]
            PR[ProductRepository]
            TR[TransactionRepository]
        end

        subgraph UT[Utils]
            JWT[JWT Utils]
            BCR[Bcrypt Utils]
            WA[WhatsApp Utils]
            RESP[Response Utils]
        end
    end

    %% DB
    subgraph DB[PostgreSQL]
        S[(shops)]
        U[(users)]
        P[(products)]
        T[(transactions)]
    end

    SA --> R
    AD --> R
    GU --> R

    R --> H
    R --> L
    R --> R1
    R --> P1
    R --> P2
    R --> P3
    R --> P4
    R --> T1
    R --> D1
    R --> U1

    %% Private route chain
    R1 --> M1 --> M2 --> M3 --> AH
    P1 --> M1 --> M2 --> M3 --> PH
    P2 --> M1 --> M2 --> M3 --> PH
    P3 --> M1 --> M2 --> M3 --> PH
    P4 --> M1 --> M2 --> M3 --> PH
    T1 --> M1 --> M2 --> M3 --> TH
    D1 --> M1 --> M2 --> M3 --> RH

    %% Public routes
    H --> RESP
    L --> AH
    U1 --> UH

    AH --> AS
    PH --> PS
    TH --> TS
    RH --> RS
    UH --> US

    AS --> UR
    AS --> SR
    AS --> JWT
    AS --> BCR

    PS --> PR
    TS --> PR
    TS --> TR
    RS --> PR
    RS --> TR
    US --> SR
    US --> PR
    US --> WA

    UR --> U
    SR --> S
    PR --> P
    TR --> T

    P --> S
    U --> S
    T --> S
    T --> P
```

## 2) Pipeline de securite des routes privees

```mermaid
flowchart TB
    REQ[Incoming private request]
    AUTH[AuthMiddleware\n- Read Authorization\n- Validate JWT signature/exp\n- Inject claims]
    ISO[ShopIsolationMiddleware\n- Read claims.shopID\n- Validate UUID\n- Reject cross-shop path/query\n- Inject shop_id context]
    ROLE[RoleMiddleware\n- Compare claims.role with allowed roles]
    HANDLER[Handler]
    SERVICE[Service]
    REPO[Repository]
    DB[(PostgreSQL)]
    RESP[HTTP response]

    REQ --> AUTH
    AUTH -->|401| RESP
    AUTH --> ISO
    ISO -->|401/403| RESP
    ISO --> ROLE
    ROLE -->|403| RESP
    ROLE --> HANDLER --> SERVICE --> REPO --> DB --> RESP
```

## 3) Sequence detaillee: POST /transactions (type=Sale)

```mermaid
sequenceDiagram
    autonumber
    participant C as Client (Admin/SuperAdmin)
    participant G as Gin Router
    participant AM as AuthMiddleware
    participant SIM as ShopIsolationMiddleware
    participant RM as RoleMiddleware
    participant TH as TransactionHandler
    participant TS as TransactionService
    participant PR as ProductRepository
    participant TR as TransactionRepository
    participant PG as PostgreSQL

    C->>G: POST /transactions + Bearer JWT
    G->>AM: Validate JWT
    AM-->>G: claims{userID, role, shopID}
    G->>SIM: Validate tenant isolation
    SIM-->>G: shop_id in context
    G->>RM: Check role in [Admin, SuperAdmin]
    RM-->>G: OK
    G->>TH: Create()
    TH->>TS: Create(shopID_from_context, payload)

    TS->>PG: BEGIN
    TS->>PR: FindByIDAndShopIDForUpdate(productID, shopID)
    PR->>PG: SELECT ... FOR UPDATE
    PG-->>PR: product row
    PR-->>TS: product(stock)

    alt stock < quantity
        TS->>PG: ROLLBACK
        TS-->>TH: ErrInsufficientStock
        TH-->>C: 400 INSUFFICIENT_STOCK
    else stock >= quantity
        TS->>PR: SaveWithTx(product.stock - quantity)
        PR->>PG: UPDATE products SET stock = stock - qty
        TS->>TR: Create(transaction)
        TR->>PG: INSERT INTO transactions
        TS->>PG: COMMIT
        TS-->>TH: TransactionResponse
        TH-->>C: 201 Created
    end
```

## 4) Sequence detaillee: GET /public/:shopID/products

```mermaid
sequenceDiagram
    autonumber
    participant G as Guest Client
    participant R as Gin Router
    participant UH as PublicHandler
    participant US as PublicService
    participant SR as ShopRepository
    participant PR as ProductRepository
    participant WA as WhatsApp Utils
    participant PG as PostgreSQL

    G->>R: GET /public/:shopID/products
    R->>UH: ListPublicProducts()
    UH->>US: ListProductsByShopID(shopID_from_path)
    US->>SR: FindByID(shopID)
    SR->>PG: SELECT * FROM shops WHERE id = ?
    PG-->>SR: shop row
    SR-->>US: shop(active, whatsapp_number)

    alt shop inactive/not found
        US-->>UH: ErrShopInactive/ErrShopNotFound
        UH-->>G: 403/404
    else shop active
        US->>PR: ListByShopID(shopID)
        PR->>PG: SELECT * FROM products WHERE shop_id = ?
        PG-->>PR: product rows
        PR-->>US: products
        loop each product
            US->>WA: GenerateWhatsAppLink(shop.whatsapp, product.name)
            WA-->>US: whatsappLink
        end
        US-->>UH: []PublicProductResponse (no purchasePrice)
        UH-->>G: 200 OK
    end
```

## 5) Modele de donnees (ERD)

```mermaid
erDiagram
    SHOPS ||--o{ USERS : owns
    SHOPS ||--o{ PRODUCTS : owns
    SHOPS ||--o{ TRANSACTIONS : owns
    PRODUCTS ||--o{ TRANSACTIONS : referenced_by

    SHOPS {
        uuid id PK
        string name
        bool active
        string whatsapp_number
        timestamptz created_at
    }

    USERS {
        uuid id PK
        string name
        string email
        string password_hash
        enum role "SuperAdmin|Admin"
        uuid shop_id FK
        timestamptz created_at
    }

    PRODUCTS {
        uuid id PK
        string name
        text description
        string category
        decimal purchase_price
        decimal selling_price
        int stock
        string image_url
        uuid shop_id FK
        timestamptz created_at
    }

    TRANSACTIONS {
        uuid id PK
        enum type "Sale|Expense|Withdrawal"
        uuid product_id FK "nullable"
        int quantity
        decimal amount
        uuid shop_id FK
        timestamptz created_at
    }
```

## 6) Matrice d'acces

```mermaid
flowchart LR
    SA[SuperAdmin]
    AD[Admin]
    GU[Guest]

    EP1[/auth/register/]
    EP2[/auth/login/]
    EP3[/products CRUD/]
    EP4[/transactions POST/]
    EP5[/reports/dashboard/]
    EP6[/public/:shopID/products/]

    SA --> EP1
    SA --> EP2
    SA --> EP3
    SA --> EP4
    SA --> EP5

    AD --> EP2
    AD --> EP3
    AD --> EP4

    GU --> EP6
```

## Notes importantes

- L'isolation multi-tenant est enforcee au niveau middleware + repository (filtre `shop_id`).
- Les routes publiques sont isolees par `:shopID` et ne renvoient jamais `purchasePrice`.
- Les ventes decremente le stock de facon atomique dans une transaction SQL.
