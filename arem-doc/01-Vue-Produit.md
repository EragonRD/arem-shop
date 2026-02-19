# Vue produit

## Mission de l'API

Arem Shop expose une API REST pour gerer plusieurs boutiques d'electronique (tenants) avec isolation stricte entre shops.

## Acteurs

- `SuperAdmin`: pilote son shop, cree des utilisateurs, accede au dashboard.
- `Admin`: gere produits et transactions de son shop.
- `Guest`: consulte le catalogue public d'un shop via URL publique.

## Capacites metier

- Authentification par shop (`/auth/login` avec `shopID`).
- Gestion de catalogue (`/products` CRUD) scopee par shop.
- Transactions financieres (`Sale`, `Expense`, `Withdrawal`).
- Dashboard financier (`/reports/dashboard`) reserve `SuperAdmin`.
- Catalogue public (`/public/:shopID/products`) avec liens WhatsApp.

## Invariants produit

- Isolation tenant obligatoire sur toutes les routes privees.
- `purchasePrice` jamais expose en public, et masque pour `Admin`.
- Les ventes decremente le stock de facon atomique (transaction SQL + lock row).

## Suite

- [[02-Architecture-API]]
- [[03-Modele-Donnees]]
