package complete

import (
	"gorm.io/gorm"
	"pedroprado.transaction.api/src/core/_interfaces"
	"pedroprado.transaction.api/src/core/domain/values"
)

type completeTransactionService struct {
	transactionRepository       _interfaces.TransactionRepository
	transactionStatusRepository _interfaces.TransactionStatusRepository
	accountRepository           _interfaces.AccountRepository
	balanceProvisionRepository  _interfaces.BalanceProvisionRepository
	postgresClient              *gorm.DB
}

func NewCompleteTransactionService(
	transactionRepository _interfaces.TransactionRepository,
	transactionStatusRepository _interfaces.TransactionStatusRepository,
	accountRepository _interfaces.AccountRepository,
	balanceProvisionRepository _interfaces.BalanceProvisionRepository,
	postgresClient *gorm.DB,
) _interfaces.CompleteTransactionService {
	return &completeTransactionService{
		transactionRepository:       transactionRepository,
		transactionStatusRepository: transactionStatusRepository,
		accountRepository:           accountRepository,
		balanceProvisionRepository:  balanceProvisionRepository,
		postgresClient:              postgresClient,
	}
}

func (ref *completeTransactionService) Complete(transactionID string) error {
	return ref.postgresClient.Transaction(func(tx *gorm.DB) error {
		transactionRepoTx := ref.transactionRepository.WithTransaction(tx)
		transactionStatusRepoTx := ref.transactionStatusRepository.WithTransaction(tx)
		balanceProvisionRepoTx := ref.balanceProvisionRepository.WithTransaction(tx)
		accountRepoTx := ref.accountRepository.WithTransaction(tx)

		return completeAddTransaction(transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx)
	})
}

func completeAddTransaction(
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

	if transactionStatus.Status != values.TransactionStatusOpen {
		return values.NewErrorValidation("cannot complete a transaction when not OPEN")
	}

	balanceProvisions, err := balanceProvisionRepoTx.FindByTransactionID(transactionID)
	if err != nil {
		return err
	}

	balanceProvisionToComplete := balanceProvisions.FindProvisionToComplete()
	if balanceProvisionToComplete == nil {
		return values.NewErrorValidation("no provision found to complete")
	}

	intermediaryAccount, err := accountRepoTx.GetLock(values.IntermediaryAccountID)
	if err != nil {
		return err
	}
	if intermediaryAccount == nil {
		return values.ErrorIntermediaryAccountNotFound
	}

	destinationAccount, err := accountRepoTx.GetLock(balanceProvisionToComplete.DestinationAccountID)
	if err != nil {
		return err
	}
	if destinationAccount == nil {
		return values.ErrorDestinationAccountNotFound
	}

	intermediaryAccount, err = intermediaryAccount.RemoveFunds(balanceProvisionToComplete.Value)
	if err != nil {
		return err
	}
	destinationAccount = destinationAccount.AddFunds(balanceProvisionToComplete.Value)

	balanceProvisionToComplete.Status = values.ProvisionStatusClosed
	transactionStatus.Status = values.TransactionStatusBooked

	_, err = accountRepoTx.Update(*intermediaryAccount)
	if err != nil {
		return err
	}

	_, err = accountRepoTx.Update(*destinationAccount)
	if err != nil {
		return err
	}

	_, err = balanceProvisionRepoTx.Update(*balanceProvisionToComplete)
	if err != nil {
		return err
	}

	_, err = transactionStatusRepoTx.Update(*transactionStatus)
	if err != nil {
		return err
	}

	return nil
}
