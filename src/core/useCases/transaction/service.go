package transaction

import (
	"gorm.io/gorm"
	"pedroprado.transaction.api/src/core/_interfaces"
	"pedroprado.transaction.api/src/core/domain/entity"
	"pedroprado.transaction.api/src/core/domain/values"
)

var (
	errorDestinationAccountNotFound  = values.NewErrorValidation("destination account not found")
	errorOriginAccountNotFound       = values.NewErrorValidation("origin account not found")
	errorIntermediaryAccountNotFound = values.NewErrorValidation("intermediary account not found")
)

const (
	intermediaryAccountID = "12345"
)

type transactionService struct {
	transactionRepository       _interfaces.TransactionRepository
	transactionStatusRepository _interfaces.TransactionStatusRepository
	accountRepository           _interfaces.AccountRepository
	balanceProvisionRepository  _interfaces.BalanceProvisionRepository
	postgresClient              *gorm.DB
}

func NewTransactionService(
	transactionRepository _interfaces.TransactionRepository,
	transactionStatusRepository _interfaces.TransactionStatusRepository,
	accountRepository _interfaces.AccountRepository,
	balanceProvisionRepository _interfaces.BalanceProvisionRepository,
	postgresClient *gorm.DB,
) _interfaces.TransactionService {
	return &transactionService{
		transactionRepository:       transactionRepository,
		transactionStatusRepository: transactionStatusRepository,
		accountRepository:           accountRepository,
		balanceProvisionRepository:  balanceProvisionRepository,
		postgresClient:              postgresClient,
	}
}

func (ref *transactionService) Get(transactionID string) (*entity.Transaction, error) {
	return ref.transactionRepository.Get(transactionID)
}

func (ref *transactionService) CreateTransaction(transaction entity.Transaction) (*entity.Transaction, error) {
	if transaction.Type != values.TransactionTypePixOut {
		return nil, values.NewErrorValidation("only PIX OUT available")
	}

	destinationAccount, err := ref.accountRepository.Get(transaction.DestinationAccountID)
	if err != nil {
		return nil, err
	}
	if destinationAccount == nil {
		return nil, values.NewErrorValidation("destination account not found")
	}

	var createdTransaction *entity.Transaction

	transactionErr := ref.postgresClient.Transaction(func(tx *gorm.DB) error {
		accountRepoTx := ref.accountRepository.WithTransaction(tx)
		transactionRepoTx := ref.transactionRepository.WithTransaction(tx)
		transactionStatusRepoTx := ref.transactionStatusRepository.WithTransaction(tx)
		balanceProvisionRepoTx := ref.balanceProvisionRepository.WithTransaction(tx)

		createdTransaction, err = createTransactionTx(
			transaction, accountRepoTx, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx)
		if err != nil {
			return err
		}

		return nil
	})

	return createdTransaction, transactionErr
}

func (ref *transactionService) CompleteTransaction(transactionID string) error {

	return ref.postgresClient.Transaction(func(tx *gorm.DB) error {
		transactionRepoTx := ref.transactionRepository.WithTransaction(tx)
		transactionStatusRepoTx := ref.transactionStatusRepository.WithTransaction(tx)
		balanceProvisionRepoTx := ref.balanceProvisionRepository.WithTransaction(tx)
		accountRepoTx := ref.accountRepository.WithTransaction(tx)

		return completeAddTransaction(transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx)
	})
}

func (ref *transactionService) CompensateTransaction(transactionID string) error {

	return nil
}

func createTransactionTx(
	transaction entity.Transaction,
	accountRepoTx _interfaces.AccountRepository,
	transactionRepoTx _interfaces.TransactionRepository,
	transactionStatusRepoTx _interfaces.TransactionStatusRepository,
	balanceProvisionRepoTx _interfaces.BalanceProvisionRepository,
) (*entity.Transaction, error) {
	originAccount, err := accountRepoTx.GetLock(transaction.OriginAccountID)
	if err != nil {
		return nil, err
	}
	if originAccount == nil {
		return nil, errorOriginAccountNotFound
	}
	intermediaryAccount, err := accountRepoTx.GetLock(intermediaryAccountID)
	if err != nil {
		return nil, err
	}
	if intermediaryAccount == nil {
		return nil, errorIntermediaryAccountNotFound
	}

	originAccount, err = originAccount.RemoveFunds(transaction.Value)
	if err != nil {
		return nil, err
	}
	intermediaryAccount = intermediaryAccount.AddFunds(transaction.Value)

	createdTransaction, err := transactionRepoTx.Create(transaction)
	if err != nil {
		return nil, err
	}
	transactionStatus := entity.TransactionStatus{
		Status:        values.TransactionStatusOpen,
		TransactionID: createdTransaction.TransactionID,
	}
	_, err = transactionStatusRepoTx.Create(transactionStatus)
	if err != nil {
		return nil, err
	}

	_, err = accountRepoTx.Update(*originAccount)
	if err != nil {
		return nil, err
	}

	_, err = accountRepoTx.Update(*intermediaryAccount)
	if err != nil {
		return nil, err
	}

	balanceProvision := entity.BalanceProvision{
		Value:                transaction.Value,
		OriginAccountID:      transaction.OriginAccountID,
		DestinationAccountID: transaction.DestinationAccountID,
		Type:                 values.ProvisionTypeAdd,
		Status:               values.ProvisionStatusOpen,
		TransactionID:        createdTransaction.TransactionID,
	}
	_, err = balanceProvisionRepoTx.Create(balanceProvision)
	if err != nil {
		return nil, err
	}
	return createdTransaction, nil
}

func completeAddTransaction(
	transactionID string,
	transactionRepoTx _interfaces.TransactionRepository,
	transactionStatusRepoTx _interfaces.TransactionStatusRepository,
	balanceProvisionRepoTx _interfaces.BalanceProvisionRepository,
	accountRepoTx _interfaces.AccountRepository,
) error {
	// TODO: dar get e lock aqui também, para evitar multiplas açoes de fechamento de transacao concorrentes
	transaction, err := transactionRepoTx.Get(transactionID)
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

	intermediaryAccount, err := accountRepoTx.GetLock(intermediaryAccountID)
	if err != nil {
		return err
	}
	if intermediaryAccount == nil {
		return errorIntermediaryAccountNotFound
	}

	destinationAccount, err := accountRepoTx.GetLock(balanceProvisionToComplete.DestinationAccountID)
	if err != nil {
		return err
	}
	if destinationAccount == nil {
		return errorDestinationAccountNotFound
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
