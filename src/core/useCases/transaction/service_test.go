package transaction

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	mocks "pedroprado.transaction.api/src/core/_mocks"
	"pedroprado.transaction.api/src/core/domain/entity"
	"pedroprado.transaction.api/src/core/domain/values"
	"testing"
)

func TestCreateTransactionTx_shouldCreateTransactionSuccefully(t *testing.T) {
	transaction := entity.Transaction{
		Type:                 values.TransactionTypePixOut,
		OriginAccountID:      uuid.NewString(),
		DestinationAccountID: uuid.NewString(),
		Value:                100,
	}
	accountRepoTx := &mocks.AccountRepository{}
	transactionRepoTx := &mocks.TransactionRepository{}
	transactionStatusRepoTx := &mocks.TransactionStatusRepository{}
	balanceProvisionRepoTx := &mocks.BalanceProvisionRepository{}

	originAccount := &entity.Account{Balance: 200}
	intermediaryAccount := &entity.Account{Balance: 300}
	originAccountRemovedFunds := &entity.Account{Balance: 100}
	intermediaryAccountAddedFunds := &entity.Account{Balance: 400}
	createdTransaction := &entity.Transaction{TransactionID: uuid.NewString()}
	balanceProvision := entity.BalanceProvision{
		Value:                transaction.Value,
		OriginAccountID:      transaction.OriginAccountID,
		DestinationAccountID: transaction.DestinationAccountID,
		Type:                 values.ProvisionTypeAdd,
		Status:               values.ProvisionStatusOpen,
		TransactionID:        createdTransaction.TransactionID,
	}

	accountRepoTx.On("GetLock", transaction.OriginAccountID).Return(originAccount, nil)
	accountRepoTx.On("GetLock", intermediaryAccountID).Return(intermediaryAccount, nil)
	transactionRepoTx.On("Create", transaction).Return(createdTransaction, nil)
	transactionStatusRepoTx.On("Create", entity.TransactionStatus{
		Status:        values.TransactionStatusOpen,
		TransactionID: createdTransaction.TransactionID,
	}).Return(&entity.TransactionStatus{}, nil)
	accountRepoTx.On("Update", *originAccountRemovedFunds).Return(originAccountRemovedFunds, nil)
	accountRepoTx.On("Update", *intermediaryAccountAddedFunds).Return(intermediaryAccountAddedFunds, nil)
	balanceProvisionRepoTx.On("Create", balanceProvision).Return(&balanceProvision, nil)

	expected := createdTransaction
	received, err := createTransactionTx(transaction, accountRepoTx, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx)

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestCreateTransactionTx_shouldNotCreateWhenOriginNotFound(t *testing.T) {
	transaction := entity.Transaction{
		Type:                 values.TransactionTypePixOut,
		OriginAccountID:      uuid.NewString(),
		DestinationAccountID: uuid.NewString(),
		Value:                100,
	}
	accountRepoTx := &mocks.AccountRepository{}
	transactionRepoTx := &mocks.TransactionRepository{}
	transactionStatusRepoTx := &mocks.TransactionStatusRepository{}
	balanceProvisionRepoTx := &mocks.BalanceProvisionRepository{}

	var originAccount *entity.Account

	accountRepoTx.On("GetLock", transaction.OriginAccountID).Return(originAccount, nil)

	received, err := createTransactionTx(transaction, accountRepoTx, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx)

	assert.Nil(t, received)
	assert.Equal(t, errorOriginAccountNotFound, err)
	accountRepoTx.AssertNumberOfCalls(t, "GetLock", 1)
	transactionRepoTx.AssertNumberOfCalls(t, "Create", 0)
	transactionStatusRepoTx.AssertNumberOfCalls(t, "Create", 0)
	accountRepoTx.AssertNumberOfCalls(t, "Update", 0)
	balanceProvisionRepoTx.AssertNumberOfCalls(t, "Create", 0)
}

func TestCreateTransactionTx_shouldNotCreateWhenIntermediaryAccountNotFound(t *testing.T) {
	transaction := entity.Transaction{
		Type:                 values.TransactionTypePixOut,
		OriginAccountID:      uuid.NewString(),
		DestinationAccountID: uuid.NewString(),
		Value:                100,
	}
	accountRepoTx := &mocks.AccountRepository{}
	transactionRepoTx := &mocks.TransactionRepository{}
	transactionStatusRepoTx := &mocks.TransactionStatusRepository{}
	balanceProvisionRepoTx := &mocks.BalanceProvisionRepository{}

	originAccount := &entity.Account{Balance: 200}
	var intermediaryAccount *entity.Account

	accountRepoTx.On("GetLock", transaction.OriginAccountID).Return(originAccount, nil)
	accountRepoTx.On("GetLock", intermediaryAccountID).Return(intermediaryAccount, nil)

	received, err := createTransactionTx(transaction, accountRepoTx, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx)

	assert.Nil(t, received)
	assert.Equal(t, errorIntermediaryAccountNotFound, err)
	transactionRepoTx.AssertNumberOfCalls(t, "Create", 0)
	transactionStatusRepoTx.AssertNumberOfCalls(t, "Create", 0)
	accountRepoTx.AssertNumberOfCalls(t, "Update", 0)
	balanceProvisionRepoTx.AssertNumberOfCalls(t, "Create", 0)
}

func TestCreateTransactionTx_shouldNotCreateWhenOriginHasNotFunds(t *testing.T) {
	transaction := entity.Transaction{
		Type:                 values.TransactionTypePixOut,
		OriginAccountID:      uuid.NewString(),
		DestinationAccountID: uuid.NewString(),
		Value:                100,
	}
	accountRepoTx := &mocks.AccountRepository{}
	transactionRepoTx := &mocks.TransactionRepository{}
	transactionStatusRepoTx := &mocks.TransactionStatusRepository{}
	balanceProvisionRepoTx := &mocks.BalanceProvisionRepository{}

	originAccount := &entity.Account{Balance: 0}
	intermediaryAccount := &entity.Account{Balance: 200}

	accountRepoTx.On("GetLock", transaction.OriginAccountID).Return(originAccount, nil)
	accountRepoTx.On("GetLock", intermediaryAccountID).Return(intermediaryAccount, nil)

	received, err := createTransactionTx(transaction, accountRepoTx, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx)

	assert.Nil(t, received)
	assert.NotNil(t, err)
	transactionRepoTx.AssertNumberOfCalls(t, "Create", 0)
	transactionStatusRepoTx.AssertNumberOfCalls(t, "Create", 0)
	accountRepoTx.AssertNumberOfCalls(t, "Update", 0)
	balanceProvisionRepoTx.AssertNumberOfCalls(t, "Create", 0)
}
