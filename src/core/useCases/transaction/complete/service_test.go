package complete

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	mocks "pedroprado.transaction.api/src/core/_mocks"
	"testing"
)

func TestComplete_shouldCompleteSuccessfully(t *testing.T) {
	transactionID := uuid.NewString()
	transactionRepoTx := &mocks.TransactionRepository{}
	transactionStatusRepoTx := &mocks.TransactionStatusRepository{}
	balanceProvisionRepoTx := &mocks.BalanceProvisionRepository{}
	accountRepoTx := &mocks.AccountRepository{}
	completeTransactionTxService := &mocks.CompleteTransactionTxService{}
	generateCompensationTxService := &mocks.GenerateCompensationTxService{}

	completeTransactionTxService.On("Complete",
		transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx).Return(nil)

	err := complete(transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx,
		completeTransactionTxService, generateCompensationTxService)

	assert.Nil(t, err)
	generateCompensationTxService.AssertNumberOfCalls(t, "Generate", 0)
}

func TestComplete_shouldGenerateCompensationWhenFailedCompletingTransaction(t *testing.T) {
	transactionID := uuid.NewString()
	transactionRepoTx := &mocks.TransactionRepository{}
	transactionStatusRepoTx := &mocks.TransactionStatusRepository{}
	balanceProvisionRepoTx := &mocks.BalanceProvisionRepository{}
	accountRepoTx := &mocks.AccountRepository{}
	completeTransactionTxService := &mocks.CompleteTransactionTxService{}
	generateCompensationTxService := &mocks.GenerateCompensationTxService{}

	completeTransactionTxService.On("Complete",
		transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx).Return(errors.New("error completing transaction"))
	generateCompensationTxService.On("Generate", transactionID, transactionStatusRepoTx, balanceProvisionRepoTx).Return(nil)

	err := complete(transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx,
		completeTransactionTxService, generateCompensationTxService)

	assert.Nil(t, err)
	generateCompensationTxService.AssertNumberOfCalls(t, "Generate", 1)
}
