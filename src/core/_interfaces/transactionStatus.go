package _interfaces

import (
	"gorm.io/gorm"
	"pedroprado.transaction.api/src/core/domain/entity"
)

type TransactionStatusService interface {
	FindByTransactionID(transactionID string) (*entity.TransactionStatus, error)
}

type TransactionStatusRepository interface {
	FindByTransactionID(transactionID string) (*entity.TransactionStatus, error)
	Create(transactionStatus entity.TransactionStatus) (*entity.TransactionStatus, error)
	Update(transactionStatus entity.TransactionStatus) (*entity.TransactionStatus, error)
	WithTransaction(tx *gorm.DB) TransactionStatusRepository
}
