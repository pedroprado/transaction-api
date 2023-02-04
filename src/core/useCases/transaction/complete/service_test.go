package complete

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	mocks "pedroprado.transaction.api/src/core/_mocks"
	"pedroprado.transaction.api/src/core/domain/entity"
	"pedroprado.transaction.api/src/core/domain/values"
	"testing"
)

func TestTransactionService_CompleteAddTransaction(t *testing.T) {
	transactionID := uuid.NewString()
	transactionRepoTx := &mocks.TransactionRepository{}
	transactionStatusRepoTx := &mocks.TransactionStatusRepository{}
	balanceProvisionRepoTx := &mocks.BalanceProvisionRepository{}
	accountRepoTx := &mocks.AccountRepository{}

	transaction := &entity.Transaction{}
	transactionStatusOpen := &entity.TransactionStatus{
		TransactionStatusID: uuid.NewString(),
		Status:              values.TransactionStatusOpen}
	balanceProvisions := entity.BalanceProvisions{
		{
			ProvisionID:          uuid.NewString(),
			Type:                 values.ProvisionTypeAdd,
			Status:               values.ProvisionStatusOpen,
			DestinationAccountID: uuid.NewString(),
			Value:                50,
		},
	}
	intermediaryAccount := &entity.Account{AccountID: values.IntermediaryAccountID, Balance: 100}
	destinationAccount := &entity.Account{AccountID: uuid.NewString(), Balance: 50}

	intermediaryAccountUpdated := *intermediaryAccount
	intermediaryAccountUpdated.Balance = 50
	destinationAccountUpdated := *destinationAccount
	destinationAccountUpdated.Balance = 100
	balanceProvisionUpdated := balanceProvisions[0]
	balanceProvisionUpdated.Status = values.ProvisionStatusClosed
	transactionStatusUpdated := *transactionStatusOpen
	transactionStatusUpdated.Status = values.TransactionStatusBooked

	transactionRepoTx.On("Get", transactionID).Return(transaction, nil)
	transactionStatusRepoTx.On("FindByTransactionID", transactionID).Return(transactionStatusOpen, nil)
	balanceProvisionRepoTx.On("FindByTransactionID", transactionID).Return(balanceProvisions, nil)
	accountRepoTx.On("GetLock", values.IntermediaryAccountID).Return(intermediaryAccount, nil)
	accountRepoTx.On("GetLock", balanceProvisions[0].DestinationAccountID).Return(destinationAccount, nil)

	accountRepoTx.On("Update", intermediaryAccountUpdated).Return(&intermediaryAccountUpdated, nil)
	accountRepoTx.On("Update", destinationAccountUpdated).Return(&destinationAccountUpdated, nil)
	balanceProvisionRepoTx.On("Update", balanceProvisionUpdated).Return(&balanceProvisionUpdated, nil)
	transactionStatusRepoTx.On("Update", transactionStatusUpdated).Return(&transactionStatusUpdated, nil)

	err := completeAddTransaction(transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx)

	assert.Nil(t, err)
}

func TestTransactionService_CompleteAddTransaction_ShouldNotCompleteWhenTransactionNotFound(t *testing.T) {
	transactionID := uuid.NewString()
	transactionRepoTx := &mocks.TransactionRepository{}
	transactionStatusRepoTx := &mocks.TransactionStatusRepository{}
	balanceProvisionRepoTx := &mocks.BalanceProvisionRepository{}
	accountRepoTx := &mocks.AccountRepository{}

	var transaction *entity.Transaction

	transactionRepoTx.On("Get", transactionID).Return(transaction, nil)

	err := completeAddTransaction(transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx)

	assert.Equal(t, values.NewErrorValidation("transaction not found"), err)
}

func TestTransactionService_CompleteAddTransaction_ShouldNotCompleteWhenTransactionStatusNotFound(t *testing.T) {
	transactionID := uuid.NewString()
	transactionRepoTx := &mocks.TransactionRepository{}
	transactionStatusRepoTx := &mocks.TransactionStatusRepository{}
	balanceProvisionRepoTx := &mocks.BalanceProvisionRepository{}
	accountRepoTx := &mocks.AccountRepository{}

	transaction := &entity.Transaction{}
	var transactionStatusOpen *entity.TransactionStatus

	transactionRepoTx.On("Get", transactionID).Return(transaction, nil)
	transactionStatusRepoTx.On("FindByTransactionID", transactionID).Return(transactionStatusOpen, nil)

	err := completeAddTransaction(transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx)

	assert.Equal(t, values.NewErrorValidation("transaction status not found"), err)
}

func TestTransactionService_CompleteAddTransaction_ShouldNotCompleteWhenTransactionStatusNotOpen(t *testing.T) {
	transactionID := uuid.NewString()
	transactionRepoTx := &mocks.TransactionRepository{}
	transactionStatusRepoTx := &mocks.TransactionStatusRepository{}
	balanceProvisionRepoTx := &mocks.BalanceProvisionRepository{}
	accountRepoTx := &mocks.AccountRepository{}

	transaction := &entity.Transaction{}
	transactionStatusOpen := &entity.TransactionStatus{
		Status: values.TransactionStatusBooked,
	}

	transactionRepoTx.On("Get", transactionID).Return(transaction, nil)
	transactionStatusRepoTx.On("FindByTransactionID", transactionID).Return(transactionStatusOpen, nil)

	err := completeAddTransaction(transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx)

	assert.Equal(t, values.NewErrorValidation("cannot complete a transaction when not OPEN"), err)
}

func TestTransactionService_CompleteAddTransaction_ShouldNotCompleteWhenNotBalanceProvisionToComplete(t *testing.T) {
	transactionID := uuid.NewString()
	transactionRepoTx := &mocks.TransactionRepository{}
	transactionStatusRepoTx := &mocks.TransactionStatusRepository{}
	balanceProvisionRepoTx := &mocks.BalanceProvisionRepository{}
	accountRepoTx := &mocks.AccountRepository{}

	transaction := &entity.Transaction{}
	transactionStatusOpen := &entity.TransactionStatus{
		Status: values.TransactionStatusOpen,
	}
	balanceProvisions := entity.BalanceProvisions{}

	transactionRepoTx.On("Get", transactionID).Return(transaction, nil)
	transactionStatusRepoTx.On("FindByTransactionID", transactionID).Return(transactionStatusOpen, nil)
	balanceProvisionRepoTx.On("FindByTransactionID", transactionID).Return(balanceProvisions, nil)

	err := completeAddTransaction(transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx)

	assert.Equal(t, values.NewErrorValidation("no provision found to complete"), err)
}
