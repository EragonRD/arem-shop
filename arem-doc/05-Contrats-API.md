# Contrats API

## Matrice d'acces

| Methode | Endpoint | Acces |
|---|---|---|
| GET | `/health` | Public |
| POST | `/auth/login` | Public |
| POST | `/auth/register` | SuperAdmin |
| GET | `/products` | Admin, SuperAdmin |
| POST | `/products` | Admin, SuperAdmin |
| PUT | `/products/:id` | Admin, SuperAdmin |
| DELETE | `/products/:id` | Admin, SuperAdmin |
| POST | `/transactions` | Admin, SuperAdmin |
| GET | `/reports/dashboard` | SuperAdmin |
| GET | `/public/:shopID/products` | Public |

## Enveloppe d'erreur

Format standard:

```json
{
  "error": "message",
  "code": "ERROR_CODE",
  "info": ["details optionnels"]
}
```

## Regles de contrat importantes

- Header prive: `Authorization: Bearer <token>`
- Claims JWT attendues: `userID`, `email`, `role`, `shopID`, `exp`
- `shopID` effectif des routes privees vient du JWT (pas du client)
- Validation stricte des payloads via tags `binding`

## Lien avec l'implementation

- Router: `cmd/api/main.go`
- DTO: `internal/dto/*.go`
- Erreurs HTTP: `internal/utils/response.go`

## Suite

- [[06-Runbook-Technique]]
