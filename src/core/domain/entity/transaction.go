package entity

import "pedroprado.transaction.api/src/core/domain/values"

type Transaction struct {
	TransactionID        string
	Type                 values.TransactionType
	OriginAccountID      string
	DestinationAccountID string
	Value                float64
	TransactionStatus    TransactionStatus
	BalanceProvisions    []BalanceProvision
}
