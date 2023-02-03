package _interfaces

import (
	"gorm.io/gorm"
	"pedroprado.transaction.api/src/core/domain/entity"
)

type TransactionService interface {
	Get(transactionID string) (*entity.Transaction, error)
	CreateTransaction(transaction entity.Transaction) (*entity.Transaction, error)
	CompleteTransaction(transactionID string) (*entity.Transaction, error)
	CompensateTransaction(transactionID string) (*entity.Transaction, error)
}

type TransactionRepository interface {
	Get(transactionID string) (*entity.Transaction, error)
	Create(transaction entity.Transaction) (*entity.Transaction, error)
	WithTransaction(tx *gorm.DB) TransactionRepository
}