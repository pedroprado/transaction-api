package responses

import (
	"pedroprado.transaction.api/src/core/domain/entity"
	"time"
)

type Transaction struct {
	TransactionID        string    `json:"transaction_id"`
	Type                 string    `json:"type"`
	OriginAccountID      string    `json:"origin_account_id"`
	DestinationAccountID string    `json:"destination_account_id"`
	Value                float64   `json:"value"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func TransactionFromDomain(transaction entity.Transaction) Transaction {
	return Transaction{
		TransactionID:        transaction.TransactionID,
		Type:                 string(transaction.Type),
		OriginAccountID:      transaction.OriginAccountID,
		DestinationAccountID: transaction.DestinationAccountID,
		Value:                transaction.Value,
		CreatedAt:            transaction.CreatedAt,
		UpdatedAt:            transaction.UpdatedAt,
	}
}
