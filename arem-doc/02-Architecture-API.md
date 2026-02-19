# Architecture API

## Vue globale

![[schemas/architecture-globale.png]]

## Pipeline de securite prive

![[schemas/pipeline-securite.png]]

## Empilement technique

- Transport HTTP: Gin (`cmd/api/main.go`)
- Couches applicatives: `handlers -> services -> repository`
- Persistance: PostgreSQL via GORM
- Auth: JWT HS256
- Passwords: bcrypt

## Groupes de routes

- Public: `/health`, `/auth/login`, `/public/:shopID/products`
- Prive (`Admin`, `SuperAdmin`): `/products`, `/transactions`
- Prive (`SuperAdmin`): `/auth/register`, `/reports/dashboard`

## Points de conception

- Les handlers restent minces: validation/erreur HTTP.
- Les services portent les regles metier.
- Les repositories centralisent le filtrage SQL par `shop_id`.
- Le middleware injecte le `shop_id` source JWT pour eviter le cross-tenant.

## Suite

- [[03-Modele-Donnees]]
- [[04-Flux-Metier-Critiques]]
