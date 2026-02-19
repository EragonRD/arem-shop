# Runbook technique

## Demarrage local

1. Configurer `.env` (voir `config/.env.example`).
2. Creer la base `arem_shop`.
3. Appliquer `migrations/000001_init.up.sql`.
4. Lancer l'API: `go run ./cmd/api`.
5. Verifier: `GET /health`.

## Verification minimale post-deploiement

- Healthcheck HTTP OK.
- Database status `up` dans `/health`.
- Login valide avec un user de seed.
- Test d'une route privee et d'une route publique.

## Risques a surveiller

- Regressions multi-tenant (cross-shop).
- Concurrence stock sur `Sale`.
- Exposition accidentelle de `purchasePrice`.
- Mauvaise configuration `JWT_SECRET`.

## Evolutions conseillees

- Completer `docs/openapi.yaml` avec tous les endpoints.
- Ajouter des tests d'integration e2e multi-tenant.
- Ajouter tracing/metrics (latence, taux d'erreur, saturation DB).

## Retour au plan

- [[00-Plan-Explication]]
