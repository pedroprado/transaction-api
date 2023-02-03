package model

import (
	"pedroprado.transaction.api/src/core/domain/entity"
	"pedroprado.transaction.api/src/core/domain/values"
	"time"
)

type BalanceProvision struct {
	ProvisionID          string `gorm:"primaryKey"`
	Value                float64
	OriginAccountID      string
	DestinationAccountID string
	Type                 string
	Status               string
	TransactionID        string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func NewBalanceProvisionFromDomain(balanceProvision entity.BalanceProvision) BalanceProvision {
	return BalanceProvision{
		ProvisionID:          balanceProvision.ProvisionID,
		Value:                balanceProvision.Value,
		OriginAccountID:      balanceProvision.OriginAccountID,
		DestinationAccountID: balanceProvision.DestinationAccountID,
		Type:                 string(balanceProvision.Type),
		Status:               string(balanceProvision.Status),
		TransactionID:        balanceProvision.TransactionID,
		CreatedAt:            balanceProvision.CreatedAt,
		UpdatedAt:            balanceProvision.UpdatedAt,
	}
}

func (record BalanceProvision) ToDomain() *entity.BalanceProvision {
	return &entity.BalanceProvision{
		ProvisionID:          record.ProvisionID,
		Value:                record.Value,
		OriginAccountID:      record.OriginAccountID,
		DestinationAccountID: record.DestinationAccountID,
		Type:                 values.ProvisionType(record.Type),
		Status:               values.ProvisionStatus(record.Status),
		TransactionID:        record.TransactionID,
		CreatedAt:            record.CreatedAt,
		UpdatedAt:            record.UpdatedAt,
	}
}
