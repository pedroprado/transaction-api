package _interfaces

import (
	"gorm.io/gorm"
	"pedroprado.transaction.api/src/core/domain/entity"
)

type TransactionService interface {
	Get(transactionID string) (*entity.Transaction, error)
	CreateTransactionService
	CompleteTransactionService
	CompensateTransactionService
}

type CreateTransactionService interface {
	Create(transaction entity.Transaction) (*entity.Transaction, error)
}

type CompleteTransactionService interface {
	Complete(transactionID string) error
}

type CompensateTransactionService interface {
	Compensate(transactionID string) error
}

type CompleteTransactionTxService interface {
	Complete(
		transactionID string,
		transactionRepoTx TransactionRepository,
		transactionStatusRepoTx TransactionStatusRepository,
		balanceProvisionRepoTx BalanceProvisionRepository,
		accountRepoTx AccountRepository,
	) error
}

type GenerateCompensationTxService interface {
	Generate(
		transactionID string,
		transactionStatusRepoTx TransactionStatusRepository,
		balanceProvisionRepoTx BalanceProvisionRepository,
	) error
}

type TransactionRepository interface {
	Get(transactionID string) (*entity.Transaction, error)
	GetLock(transactionID string) (*entity.Transaction, error)
	Create(transaction entity.Transaction) (*entity.Transaction, error)
	WithTransaction(tx *gorm.DB) TransactionRepository
}
