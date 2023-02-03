package entity

import (
	"pedroprado.transaction.api/src/core/domain/values"
	"time"
)

type BalanceProvision struct {
	ProvisionID          string
	Value                float64
	OriginAccountID      string
	DestinationAccountID string
	Type                 values.ProvisionType
	Status               values.ProvisionStatus
	TransactionID        string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
