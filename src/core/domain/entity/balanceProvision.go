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

type BalanceProvisions []BalanceProvision

func (provisions BalanceProvisions) FindProvision(provisionType values.ProvisionType, provisionStatus values.ProvisionStatus) *BalanceProvision {
	for i := range provisions {
		provision := provisions[i]
		if provision.Type == provisionType &&
			provision.Status == provisionStatus {
			return &provision
		}
	}

	return nil
}
