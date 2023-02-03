package model

import (
	"pedroprado.transaction.api/src/core/domain/entity"
	"pedroprado.transaction.api/src/core/domain/values"
	"time"
)

type Transaction struct {
	TransactionID        string `gorm:"primaryKey"`
	Type                 string
	OriginAccountID      string
	DestinationAccountID string
	Value                float64
	CreatedAt            time.Time
	UpdatedAt            time.Time

	TransactionStatus TransactionStatus  `gorm:"foreignKey:TransactionID;references:TransactionID"`
	BalanceProvisions []BalanceProvision `gorm:"foreignKey:TransactionID;references:TransactionID"`
}

func NewTransactionFromDomain(transaction entity.Transaction) Transaction {
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

func (record Transaction) ToDomain() *entity.Transaction {
	return &entity.Transaction{
		TransactionID:        record.TransactionID,
		Type:                 values.TransactionType(record.Type),
		OriginAccountID:      record.OriginAccountID,
		DestinationAccountID: record.DestinationAccountID,
		Value:                record.Value,
		CreatedAt:            record.CreatedAt,
		UpdatedAt:            record.UpdatedAt,
	}
}
