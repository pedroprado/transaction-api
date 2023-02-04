package complete

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"pedroprado.transaction.api/src/core/_interfaces"
)

type completeTransactionService struct {
	transactionRepository         _interfaces.TransactionRepository
	transactionStatusRepository   _interfaces.TransactionStatusRepository
	accountRepository             _interfaces.AccountRepository
	balanceProvisionRepository    _interfaces.BalanceProvisionRepository
	postgresClient                *gorm.DB
	completeTransactionTxService  _interfaces.CompleteTransactionTxService
	generateCompensationTxService _interfaces.GenerateCompensationTxService
}

func NewCompleteTransactionService(
	transactionRepository _interfaces.TransactionRepository,
	transactionStatusRepository _interfaces.TransactionStatusRepository,
	accountRepository _interfaces.AccountRepository,
	balanceProvisionRepository _interfaces.BalanceProvisionRepository,
	postgresClient *gorm.DB,
	completeTransactionTxService _interfaces.CompleteTransactionTxService,
	generateCompensationTxService _interfaces.GenerateCompensationTxService,
) _interfaces.CompleteTransactionService {
	return &completeTransactionService{
		transactionRepository:         transactionRepository,
		transactionStatusRepository:   transactionStatusRepository,
		accountRepository:             accountRepository,
		balanceProvisionRepository:    balanceProvisionRepository,
		postgresClient:                postgresClient,
		completeTransactionTxService:  completeTransactionTxService,
		generateCompensationTxService: generateCompensationTxService,
	}
}

func (ref *completeTransactionService) Complete(transactionID string) error {
	return ref.postgresClient.Transaction(func(tx *gorm.DB) error {
		transactionRepoTx := ref.transactionRepository.WithTransaction(tx)
		transactionStatusRepoTx := ref.transactionStatusRepository.WithTransaction(tx)
		balanceProvisionRepoTx := ref.balanceProvisionRepository.WithTransaction(tx)
		accountRepoTx := ref.accountRepository.WithTransaction(tx)

		return complete(transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx,
			ref.completeTransactionTxService, ref.generateCompensationTxService)
	})
}

func complete(
	transactionID string,
	transactionRepoTx _interfaces.TransactionRepository,
	transactionStatusRepoTx _interfaces.TransactionStatusRepository,
	balanceProvisionRepoTx _interfaces.BalanceProvisionRepository,
	accountRepoTx _interfaces.AccountRepository,
	completeTransactionTxService _interfaces.CompleteTransactionTxService,
	generateCompensationTxService _interfaces.GenerateCompensationTxService,
) error {
	err := completeTransactionTxService.Complete(transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx)
	if err != nil {
		logrus.Errorf("cannot complete transaction %s, closing and generating compensation", transactionID)
		err = generateCompensationTxService.Generate(transactionID, transactionStatusRepoTx, balanceProvisionRepoTx)
		if err != nil {
			logrus.Errorf("cannot generate compensation for transaction %s", transactionID)

			return err
		}
	}

	return nil
}
