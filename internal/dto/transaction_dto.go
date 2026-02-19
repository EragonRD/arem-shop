package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

// CreateTransactionRequest payload de creation transaction.
type CreateTransactionRequest struct {
	Type      string          `json:"type" binding:"required,oneof=Sale Expense Withdrawal"`
	ProductID *string         `json:"productID,omitempty"`
	Quantity  *int            `json:"quantity,omitempty"`
	Amount    decimal.Decimal `json:"amount"`
}

// TransactionResponse representation API d'une transaction.
type TransactionResponse struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	ProductID *string         `json:"productID,omitempty"`
	Quantity  int             `json:"quantity"`
	Amount    decimal.Decimal `json:"amount"`
	ShopID    string          `json:"shopID"`
	CreatedAt time.Time       `json:"createdAt"`
}
