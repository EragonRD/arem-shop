package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

//-----OLD-----
/*
package dto
// CreateTransactionRequest payload de creation transaction.
// NB: Pas de ShopID ici, il sera injecté via le JWT pour la sécurité multi-tenant.
type CreateTransactionRequest struct {
	Type      string          `json:"type" binding:"required,oneof=Sale Expense Withdrawal"`
	ProductID *string         `json:"productID,omitempty"`
	Quantity  *int            `json:"quantity,omitempty"`
	Amount    decimal.Decimal `json:"amount"`
}

*/

// -----NEW-----
// CreateTransactionRequest est le payload attendu pour créer une transaction.
type CreateTransactionRequest struct {
	Type      string          `json:"type" binding:"required,oneof=Sale Expense Withdrawal"`
	ProductID *uuid.UUID      `json:"productID,omitempty"` // <-- L'astérisque permet la comparaison avec nil
	Quantity  int             `json:"quantity" binding:"gte=0"`
	Amount    decimal.Decimal `json:"amount" binding:"required"`
}

//-----OLD-----
/*
// TransactionResponse representation API d'une transaction.
// TransactionResponsense est l'objet renvoyé au client pour masquer les détails internes de la DB
type TransactionResponse struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	ProductID *string         `json:"productID,omitempty"`
	Quantity  int             `json:"quantity"`
	Amount    decimal.Decimal `json:"amount"`
	ShopID    string          `json:"shopID"`
	CreatedAt time.Time       `json:"createdAt"`
}
*/

// -----NEW-----
// TransactionResponse est l'objet renvoyé au client
type TransactionResponse struct {
	ID        uuid.UUID       `json:"id"`
	Type      string          `json:"type"`
	ProductID *uuid.UUID      `json:"productID,omitempty"` // <-- Idem ici
	Quantity  int             `json:"quantity"`
	Amount    decimal.Decimal `json:"amount"`
	ShopID    uuid.UUID       `json:"shopID"`
	CreatedAt time.Time       `json:"createdAt"`
}
