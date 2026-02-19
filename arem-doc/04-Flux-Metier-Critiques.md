# Flux metier critiques

## Login JWT multi-shop

![[schemas/sequence-login.png]]

## Vente avec decrement de stock atomique

![[schemas/sequence-transaction-sale.png]]

## Catalogue public et lien WhatsApp

![[schemas/sequence-public-catalog.png]]

## Lecture CTO

- Login: le `shopID` est obligatoire pour lever l'ambiguite email multi-tenant.
- Vente: `SELECT ... FOR UPDATE` empeche les conflits concurrents sur le stock.
- Public: aucune donnee sensible (`purchasePrice`) n'est exposee.

## Suite

- [[05-Contrats-API]]
