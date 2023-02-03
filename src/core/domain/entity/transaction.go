package entity

import (
	"pedroprado.transaction.api/src/core/domain/values"
	"time"
)

type Transaction struct {
	TransactionID        string
	Type                 values.TransactionType
	OriginAccountID      string
	DestinationAccountID string
	Value                float64
	CreatedAt            time.Time
	UpdatedAt            time.Time
	TransactionStatus    TransactionStatus
	BalanceProvisions    []BalanceProvision
}
