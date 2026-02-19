# Docker & Déploiement

## Infrastructure Docker

```mermaid
graph TB
    subgraph "Docker Network: arem-shop"
        subgraph "Container: arem-shop-api"
            API["🚀 API Go<br/>Binary compilé<br/>Port 8080"]
        end
        
        subgraph "Container: arem-shop-db"
            PG["🐘 PostgreSQL 15<br/>Port 5432"]
            VOL["📦 Volume: pgdata"]
        end
    end
    
    CLIENT["🌐 Client HTTP"] -->|":8080"| API
    API -->|"DB_HOST=postgres"| PG
    PG --- VOL
    
    ENV[".env"] -.->|env_file| API
    ENV -.->|env_file| PG
    MIG["migrations/<br/>000001_init.up.sql"] -.->|"initdb.d (auto)"| PG
```

## Build multi-stage Dockerfile

```mermaid
graph LR
    subgraph "Stage 1 — Builder"
        SRC["📂 Code source<br/>go.mod + go.sum<br/>+ internal/ cmd/"]
        BUILD["🔨 go build<br/>CGO_ENABLED=0<br/>binaire statique"]
        SRC --> BUILD
    end
    
    subgraph "Stage 2 — Runner"
        BIN["📦 Binaire seul<br/>+ migrations"]
        IMG["🐳 Alpine 3.18<br/>~20 Mo"]
        BIN --> IMG
    end
    
    BUILD -->|"COPY --from=builder"| BIN
```

## Cycle de démarrage

```mermaid
sequenceDiagram
    participant U as Utilisateur
    participant DC as docker-compose
    participant PG as PostgreSQL
    participant API as API Go

    U->>DC: docker-compose up --build
    DC->>PG: Démarrer conteneur postgres
    PG->>PG: Créer DB arem_shop
    PG->>PG: Exécuter initdb.d/01-init.sql
    PG->>DC: Healthcheck OK ✅
    DC->>API: Démarrer conteneur api
    API->>PG: Connexion GORM
    API->>API: Écoute :8080
    U->>API: curl /health
    API-->>U: {"app":"arem-shop","database":"up"}
```

## Scripts de lancement

| Script | Usage | Description |
|--------|-------|-------------|
| `scripts/docker-start.sh` | `bash scripts/docker-start.sh` | Vérifie Docker, crée `.env` si absent, build & lance tout |
| `scripts/docker-stop.sh` | `bash scripts/docker-stop.sh` | Arrête les conteneurs (données DB conservées) |
| `scripts/docker-stop.sh -v` | `bash scripts/docker-stop.sh -v` | Arrête + supprime le volume DB |
| `scripts/docker-logs.sh` | `bash scripts/docker-logs.sh` | Suit les logs de tous les services |
| `scripts/docker-logs.sh api` | `bash scripts/docker-logs.sh api` | Suit les logs de l'API uniquement |

Ou via le **Makefile** :

| Commande | Equivalent |
|----------|------------|
| `make docker-up` | `bash scripts/docker-start.sh` |
| `make docker-up-detach` | `bash scripts/docker-start.sh -d` |
| `make docker-down` | `bash scripts/docker-stop.sh` |
| `make docker-clean` | `bash scripts/docker-stop.sh --volumes` |
| `make docker-logs` | `bash scripts/docker-logs.sh` |

## Commandes Docker directes

| Commande | Description |
|----------|-------------|
| `cp .env.example .env` | Créer la config locale |
| `docker-compose up --build` | Build + démarrer tout |
| `docker-compose up -d` | Démarrer en arrière-plan |
| `docker-compose logs -f api` | Suivre les logs API |
| `docker-compose logs -f postgres` | Suivre les logs DB |
| `docker-compose down` | Arrêter les conteneurs |
| `docker-compose down -v` | Arrêter + supprimer le volume DB |
| `docker-compose build --no-cache` | Rebuild sans cache |

## Troubleshooting

### `docker-compose` v1 vs `docker compose` v2

```mermaid
graph LR
    subgraph "❌ v1 — Python (cassé)"
        V1["docker-compose 1.29.2<br/>Python"]
        V1 --> ERR["URLSchemeUnknown<br/>http+docker"]
        ERR --> FAIL["💥 Crash"]
    end
    
    subgraph "✅ v2 — Go plugin (fonctionne)"
        V2["docker compose v5.0.0<br/>Go natif"]
        V2 --> OK["Build + Start<br/>correctement"]
        OK --> RUN["🚀 API + DB"]
    end
```

**Symptôme** : `URLSchemeUnknown: Not supported URL scheme http+docker`

**Cause** : Incompatibilité entre `docker-compose` v1 (Python) et les versions récentes de `requests`/`urllib3`.

**Solution** : Nos scripts détectent automatiquement `docker compose` v2 en priorité. Si l'erreur persiste :

```bash
# Vérifier que v2 est disponible
docker compose version

# Utiliser directement v2
docker compose up --build
```

### Conflit de port PostgreSQL (5432)

```mermaid
graph LR
    subgraph "Machine hôte"
        LOCAL_PG["🐘 PostgreSQL local<br/>:5432"]
        DOCKER_PG["🐳 Docker PostgreSQL<br/>:5433 → :5432 interne"]
        API_LOCAL["💻 go run ./cmd/api<br/>DB_HOST=localhost:5432"]
        API_DOCKER["🐳 Docker API<br/>DB_HOST=postgres:5432"]
    end
    
    API_LOCAL -.->|"dev local"| LOCAL_PG
    API_DOCKER -->|"réseau Docker"| DOCKER_PG
```

**Symptôme** : `address already in use` sur le port 5432

**Cause** : Un PostgreSQL local tourne déjà sur le port 5432.

**Solution** : Le `docker-compose.yml` expose la DB Docker sur le port **5433** par défaut. L'API Docker communique en interne sur le port 5432 via le réseau Docker — aucun conflit.

## Fichiers impliqués

- `Dockerfile` → build multi-stage Go
- `docker-compose.yml` → orchestration api + postgres
- `.env.example` → template de configuration (versionné)
- `.env` → configuration effective (gitignored, **jamais pushé**)
- `migrations/000001_init.up.sql` → schéma appliqué automatiquement au premier démarrage

## Suite

- [[00-Plan-Explication]]
