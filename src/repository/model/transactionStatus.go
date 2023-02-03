package model

import (
	"pedroprado.transaction.api/src/core/domain/entity"
	"pedroprado.transaction.api/src/core/domain/values"
	"time"
)

type TransactionStatus struct {
	TransactionStatusID string `gorm:"primaryKey"`
	Status              string
	TransactionID       string `gorm:"uniqueIndex:idx_transaction_status_transaction_id"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func NewTransactionStatusFromDomain(transactionStatus entity.TransactionStatus) TransactionStatus {
	return TransactionStatus{
		TransactionStatusID: transactionStatus.TransactionStatusID,
		Status:              string(transactionStatus.Status),
		TransactionID:       transactionStatus.TransactionID,
		CreatedAt:           transactionStatus.CreatedAt,
		UpdatedAt:           transactionStatus.UpdatedAt,
	}
}

func (record TransactionStatus) ToDomain() *entity.TransactionStatus {
	return &entity.TransactionStatus{
		TransactionStatusID: record.TransactionStatusID,
		Status:              values.TransactionStatus(record.Status),
		TransactionID:       record.TransactionID,
		CreatedAt:           record.CreatedAt,
		UpdatedAt:           record.UpdatedAt,
	}
}
