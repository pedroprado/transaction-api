package compensate

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	mocks "pedroprado.transaction.api/src/core/_mocks"
	"pedroprado.transaction.api/src/core/domain/entity"
	"pedroprado.transaction.api/src/core/domain/values"
	"testing"
)

func TestCompensate_ShouldCompensateTransactionSuccessfully(t *testing.T) {
	transactionID := uuid.NewString()
	transactionRepoTx := &mocks.TransactionRepository{}
	transactionStatusRepoTx := &mocks.TransactionStatusRepository{}
	balanceProvisionRepoTx := &mocks.BalanceProvisionRepository{}
	accountRepoTx := &mocks.AccountRepository{}

	transaction := &entity.Transaction{}
	transactionStatusFailed := &entity.TransactionStatus{
		TransactionStatusID: uuid.NewString(),
		Status:              values.TransactionStatusFailed,
	}
	balanceProvisions := entity.BalanceProvisions{
		{
			ProvisionID:     uuid.NewString(),
			Type:            values.ProvisionTypeVoid,
			Status:          values.ProvisionStatusOpen,
			OriginAccountID: uuid.NewString(),
			Value:           50,
		},
	}
	intermediaryAccount := &entity.Account{AccountID: values.IntermediaryAccountID, Balance: 100}
	originAccount := &entity.Account{AccountID: uuid.NewString(), Balance: 50}

	intermediaryAccountUpdated := *intermediaryAccount
	intermediaryAccountUpdated.Balance = 50
	originAccountUpdated := *originAccount
	originAccountUpdated.Balance = 100
	balanceProvisionUpdated := balanceProvisions[0]
	balanceProvisionUpdated.Status = values.ProvisionStatusClosed

	transactionRepoTx.On("GetLock", transactionID).Return(transaction, nil)
	transactionStatusRepoTx.On("FindByTransactionID", transactionID).Return(transactionStatusFailed, nil)
	balanceProvisionRepoTx.On("FindByTransactionID", transactionID).Return(balanceProvisions, nil)
	accountRepoTx.On("GetLock", values.IntermediaryAccountID).Return(intermediaryAccount, nil)
	accountRepoTx.On("GetLock", balanceProvisions[0].OriginAccountID).Return(originAccount, nil)

	accountRepoTx.On("Update", intermediaryAccountUpdated).Return(&intermediaryAccountUpdated, nil)
	accountRepoTx.On("Update", originAccountUpdated).Return(&originAccountUpdated, nil)
	balanceProvisionRepoTx.On("Update", balanceProvisionUpdated).Return(&balanceProvisionUpdated, nil)

	err := compensate(transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx)

	assert.Nil(t, err)
}

func TestCompensate_ShouldNotCompensateWhenTransactionNotFound(t *testing.T) {
	transactionID := uuid.NewString()
	transactionRepoTx := &mocks.TransactionRepository{}
	transactionStatusRepoTx := &mocks.TransactionStatusRepository{}
	balanceProvisionRepoTx := &mocks.BalanceProvisionRepository{}
	accountRepoTx := &mocks.AccountRepository{}

	var transaction *entity.Transaction

	transactionRepoTx.On("GetLock", transactionID).Return(transaction, nil)

	err := compensate(transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx)

	assert.Equal(t, values.NewErrorValidation("transaction not found"), err)
	transactionStatusRepoTx.AssertNumberOfCalls(t, "FindByTransactionID", 0)
	balanceProvisionRepoTx.AssertNumberOfCalls(t, "FindByTransactionID", 0)
	accountRepoTx.AssertNumberOfCalls(t, "GetLock", 0)
	accountRepoTx.AssertNumberOfCalls(t, "Update", 0)
	balanceProvisionRepoTx.AssertNumberOfCalls(t, "Update", 0)
}

func TestCompensate_ShouldNotCompensateWhenTransactionStatusNotFound(t *testing.T) {
	transactionID := uuid.NewString()
	transactionRepoTx := &mocks.TransactionRepository{}
	transactionStatusRepoTx := &mocks.TransactionStatusRepository{}
	balanceProvisionRepoTx := &mocks.BalanceProvisionRepository{}
	accountRepoTx := &mocks.AccountRepository{}

	transaction := &entity.Transaction{}
	var transactionStatus *entity.TransactionStatus

	transactionRepoTx.On("GetLock", transactionID).Return(transaction, nil)
	transactionStatusRepoTx.On("FindByTransactionID", transactionID).Return(transactionStatus, nil)

	err := compensate(transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx)

	assert.Equal(t, values.NewErrorValidation("transaction status not found"), err)
	balanceProvisionRepoTx.AssertNumberOfCalls(t, "FindByTransactionID", 0)
	accountRepoTx.AssertNumberOfCalls(t, "GetLock", 0)
	accountRepoTx.AssertNumberOfCalls(t, "Update", 0)
	balanceProvisionRepoTx.AssertNumberOfCalls(t, "Update", 0)
}

func TestCompensate_ShouldNotCompensateWhenTransactionNotFailed(t *testing.T) {
	transactionID := uuid.NewString()
	transactionRepoTx := &mocks.TransactionRepository{}
	transactionStatusRepoTx := &mocks.TransactionStatusRepository{}
	balanceProvisionRepoTx := &mocks.BalanceProvisionRepository{}
	accountRepoTx := &mocks.AccountRepository{}

	transaction := &entity.Transaction{}
	transactionStatus := &entity.TransactionStatus{Status: values.TransactionStatusOpen}

	transactionRepoTx.On("GetLock", transactionID).Return(transaction, nil)
	transactionStatusRepoTx.On("FindByTransactionID", transactionID).Return(transactionStatus, nil)

	err := compensate(transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx)

	assert.Equal(t, values.NewErrorValidation("cannot compensate a transaction when it has not FAILED"), err)
	balanceProvisionRepoTx.AssertNumberOfCalls(t, "FindByTransactionID", 0)
	accountRepoTx.AssertNumberOfCalls(t, "GetLock", 0)
	accountRepoTx.AssertNumberOfCalls(t, "Update", 0)
	balanceProvisionRepoTx.AssertNumberOfCalls(t, "Update", 0)
}

func TestCompensate_ShouldNotCompensateWhenNoProvisionVoidOpen(t *testing.T) {
	transactionID := uuid.NewString()
	transactionRepoTx := &mocks.TransactionRepository{}
	transactionStatusRepoTx := &mocks.TransactionStatusRepository{}
	balanceProvisionRepoTx := &mocks.BalanceProvisionRepository{}
	accountRepoTx := &mocks.AccountRepository{}

	transaction := &entity.Transaction{}
	transactionStatus := &entity.TransactionStatus{Status: values.TransactionStatusFailed}
	balanceProvisions := entity.BalanceProvisions{}

	transactionRepoTx.On("GetLock", transactionID).Return(transaction, nil)
	transactionStatusRepoTx.On("FindByTransactionID", transactionID).Return(transactionStatus, nil)
	balanceProvisionRepoTx.On("FindByTransactionID", transactionID).Return(balanceProvisions, nil)

	err := compensate(transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx)

	assert.Equal(t, values.NewErrorValidation("no provision found to complete"), err)
	accountRepoTx.AssertNumberOfCalls(t, "GetLock", 0)
	accountRepoTx.AssertNumberOfCalls(t, "Update", 0)
	balanceProvisionRepoTx.AssertNumberOfCalls(t, "Update", 0)
}

func TestCompensate_ShouldNotCompensateWhenIntermediaryAccountNotFound(t *testing.T) {
	transactionID := uuid.NewString()
	transactionRepoTx := &mocks.TransactionRepository{}
	transactionStatusRepoTx := &mocks.TransactionStatusRepository{}
	balanceProvisionRepoTx := &mocks.BalanceProvisionRepository{}
	accountRepoTx := &mocks.AccountRepository{}

	transaction := &entity.Transaction{}
	transactionStatus := &entity.TransactionStatus{Status: values.TransactionStatusFailed}
	balanceProvisions := entity.BalanceProvisions{{Type: values.ProvisionTypeVoid, Status: values.ProvisionStatusOpen}}

	transactionRepoTx.On("GetLock", transactionID).Return(transaction, nil)
	transactionStatusRepoTx.On("FindByTransactionID", transactionID).Return(transactionStatus, nil)
	balanceProvisionRepoTx.On("FindByTransactionID", transactionID).Return(balanceProvisions, nil)
	accountRepoTx.On("GetLock", values.IntermediaryAccountID).Return(nil, nil)

	err := compensate(transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx)

	assert.Equal(t, values.ErrorIntermediaryAccountNotFound, err)
	accountRepoTx.AssertNumberOfCalls(t, "GetLock", 1)
	accountRepoTx.AssertNumberOfCalls(t, "Update", 0)
	balanceProvisionRepoTx.AssertNumberOfCalls(t, "Update", 0)
}

func TestCompensate_ShouldNotCompensateWhenOriginAccountNotFound(t *testing.T) {
	transactionID := uuid.NewString()
	transactionRepoTx := &mocks.TransactionRepository{}
	transactionStatusRepoTx := &mocks.TransactionStatusRepository{}
	balanceProvisionRepoTx := &mocks.BalanceProvisionRepository{}
	accountRepoTx := &mocks.AccountRepository{}

	transaction := &entity.Transaction{}
	transactionStatus := &entity.TransactionStatus{Status: values.TransactionStatusFailed}
	balanceProvisions := entity.BalanceProvisions{{
		OriginAccountID: uuid.NewString(),
		Type:            values.ProvisionTypeVoid,
		Status:          values.ProvisionStatusOpen,
	}}

	transactionRepoTx.On("GetLock", transactionID).Return(transaction, nil)
	transactionStatusRepoTx.On("FindByTransactionID", transactionID).Return(transactionStatus, nil)
	balanceProvisionRepoTx.On("FindByTransactionID", transactionID).Return(balanceProvisions, nil)
	accountRepoTx.On("GetLock", values.IntermediaryAccountID).Return(&entity.Account{}, nil)
	accountRepoTx.On("GetLock", balanceProvisions[0].OriginAccountID).Return(nil, nil)

	err := compensate(transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx)

	assert.Equal(t, values.ErrorOriginAccountNotFound, err)
	accountRepoTx.AssertNumberOfCalls(t, "GetLock", 2)
	accountRepoTx.AssertNumberOfCalls(t, "Update", 0)
	balanceProvisionRepoTx.AssertNumberOfCalls(t, "Update", 0)
}
