package compensate

import (
	"gorm.io/gorm"
	"pedroprado.transaction.api/src/core/_interfaces"
)

type compensateTransactionService struct {
	transactionRepository       _interfaces.TransactionRepository
	transactionStatusRepository _interfaces.TransactionStatusRepository
	accountRepository           _interfaces.AccountRepository
	balanceProvisionRepository  _interfaces.BalanceProvisionRepository
	postgresClient              *gorm.DB
}

func NewCompensateTransactionService(
	transactionRepository _interfaces.TransactionRepository,
	transactionStatusRepository _interfaces.TransactionStatusRepository,
	accountRepository _interfaces.AccountRepository,
	balanceProvisionRepository _interfaces.BalanceProvisionRepository,
	postgresClient *gorm.DB,
) _interfaces.CompensateTransactionService {
	return &compensateTransactionService{
		transactionRepository:       transactionRepository,
		transactionStatusRepository: transactionStatusRepository,
		accountRepository:           accountRepository,
		balanceProvisionRepository:  balanceProvisionRepository,
		postgresClient:              postgresClient,
	}
}

func (ref *compensateTransactionService) Compensate(transactionID string) error {
	//TODO implement me
	panic("implement me")
}
