# Modele de donnees

## ERD

![[schemas/erd.png]]

## Entites principales

- `shops`: tenant principal (etat actif + numero WhatsApp).
- `users`: utilisateurs internes (`SuperAdmin`, `Admin`) rattaches a un shop.
- `products`: catalogue prive d'un shop.
- `transactions`: flux financiers d'un shop.

## Contraintes importantes

- `users`: unicite de l'email **dans** un shop (`uq_users_shop_email`).
- `products`: prix >= 0, stock >= 0.
- `transactions`: contraintes de coherence selon le type:
  - `Sale`: `product_id` non null, `quantity > 0`
  - `Expense` / `Withdrawal`: `product_id` null, `quantity = 0`
- Cle etrangere composite `(product_id, shop_id)` pour empecher une reference cross-shop.

## Invariants data

- Un `user` n'appartient qu'a un seul `shop`.
- Un `product` n'est visible/modifiable que dans son shop.
- Une `transaction` ne peut pas pointer un produit hors tenant.

## Suite

- [[04-Flux-Metier-Critiques]]
- [[05-Contrats-API]]
