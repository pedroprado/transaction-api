package compensate

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"pedroprado.transaction.api/src/core/_interfaces"
	"pedroprado.transaction.api/src/core/domain/values"
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
	return ref.postgresClient.Transaction(func(tx *gorm.DB) error {
		transactionRepoTx := ref.transactionRepository.WithTransaction(tx)
		transactionStatusRepoTx := ref.transactionStatusRepository.WithTransaction(tx)
		balanceProvisionRepoTx := ref.balanceProvisionRepository.WithTransaction(tx)
		accountRepoTx := ref.accountRepository.WithTransaction(tx)

		return compensate(transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx)
	})
}

func compensate(
	transactionID string,
	transactionRepoTx _interfaces.TransactionRepository,
	transactionStatusRepoTx _interfaces.TransactionStatusRepository,
	balanceProvisionRepoTx _interfaces.BalanceProvisionRepository,
	accountRepoTx _interfaces.AccountRepository,
) error {
	transaction, err := transactionRepoTx.GetLock(transactionID)
	if err != nil {
		return err
	}
	if transaction == nil {
		return values.NewErrorValidation("transaction not found")
	}

	transactionStatus, err := transactionStatusRepoTx.FindByTransactionID(transactionID)
	if err != nil {
		return err
	}
	if transactionStatus == nil {
		return values.NewErrorValidation("transaction status not found")
	}

	if transactionStatus.Status != values.TransactionStatusFailed {
		return values.NewErrorValidation("cannot compensate a transaction when it has not FAILED")
	}

	balanceProvisions, err := balanceProvisionRepoTx.FindByTransactionID(transactionID)
	if err != nil {
		return err
	}

	balanceProvisionVoidOpen := balanceProvisions.FindProvision(values.ProvisionTypeVoid, values.ProvisionStatusOpen)
	if balanceProvisionVoidOpen == nil {
		logrus.Errorf("no provision found to complete for transaction %s", transactionID)

		return values.NewErrorValidation("no provision found to complete")
	}

	intermediaryAccount, err := accountRepoTx.GetLock(values.IntermediaryAccountID)
	if err != nil {
		return err
	}
	if intermediaryAccount == nil {
		return values.ErrorIntermediaryAccountNotFound
	}

	originAccount, err := accountRepoTx.GetLock(balanceProvisionVoidOpen.OriginAccountID)
	if err != nil {
		return err
	}
	if originAccount == nil {
		return values.ErrorOriginAccountNotFound
	}

	intermediaryAccount, err = intermediaryAccount.RemoveFunds(balanceProvisionVoidOpen.Value)
	if err != nil {
		return err
	}
	originAccount = originAccount.AddFunds(balanceProvisionVoidOpen.Value)

	balanceProvisionVoidOpen.Status = values.ProvisionStatusClosed

	_, err = accountRepoTx.Update(*intermediaryAccount)
	if err != nil {
		return err
	}

	_, err = accountRepoTx.Update(*originAccount)
	if err != nil {
		return err
	}

	_, err = balanceProvisionRepoTx.Update(*balanceProvisionVoidOpen)
	if err != nil {
		return err
	}

	return nil
}
