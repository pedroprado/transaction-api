package create

import (
	"gorm.io/gorm"
	"pedroprado.transaction.api/src/core/_interfaces"
	"pedroprado.transaction.api/src/core/domain/entity"
	"pedroprado.transaction.api/src/core/domain/values"
)

type createTransactionService struct {
	transactionRepository       _interfaces.TransactionRepository
	transactionStatusRepository _interfaces.TransactionStatusRepository
	accountRepository           _interfaces.AccountRepository
	balanceProvisionRepository  _interfaces.BalanceProvisionRepository
	postgresClient              *gorm.DB
}

func NewCreateTransactionService(
	transactionRepository _interfaces.TransactionRepository,
	transactionStatusRepository _interfaces.TransactionStatusRepository,
	accountRepository _interfaces.AccountRepository,
	balanceProvisionRepository _interfaces.BalanceProvisionRepository,
	postgresClient *gorm.DB,
) _interfaces.CreateTransactionService {
	return &createTransactionService{
		transactionRepository:       transactionRepository,
		transactionStatusRepository: transactionStatusRepository,
		accountRepository:           accountRepository,
		balanceProvisionRepository:  balanceProvisionRepository,
		postgresClient:              postgresClient,
	}
}

func (ref *createTransactionService) Create(transaction entity.Transaction) (*entity.Transaction, error) {
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
		return nil, values.ErrorOriginAccountNotFound
	}
	intermediaryAccount, err := accountRepoTx.GetLock(values.IntermediaryAccountID)
	if err != nil {
		return nil, err
	}
	if intermediaryAccount == nil {
		return nil, values.ErrorIntermediaryAccountNotFound
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
