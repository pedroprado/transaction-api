package transactionApi

import (
	"pedroprado.transaction.api/src/core/domain/entity"
	"pedroprado.transaction.api/src/core/domain/values"
)

type CreateTransactionRequest struct {
	Type                 string  `json:"transaction_type" binding:"required"`
	OriginAccountID      string  `json:"origin_account_id" binding:"required"`
	DestinationAccountID string  `json:"destination_account_id" binding:"required"`
	Value                float64 `json:"value" binding:"required"s`
}

func (request CreateTransactionRequest) ToDomain() entity.Transaction {
	return entity.Transaction{
		Type:                 values.TransactionType(request.Type),
		OriginAccountID:      request.OriginAccountID,
		DestinationAccountID: request.DestinationAccountID,
		Value:                request.Value,
	}
}

type CompleteTransactionRequest struct {
	TransactionID string `uri:"transaction_id" binding:"required"`
}

type CompensateTransactionRequest struct {
	TransactionID string `uri:"transaction_id" binding:"required"`
}