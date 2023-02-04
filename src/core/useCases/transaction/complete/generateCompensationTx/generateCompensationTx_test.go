package generateCompensationTx

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	mocks "pedroprado.transaction.api/src/core/_mocks"
	"pedroprado.transaction.api/src/core/domain/entity"
	"pedroprado.transaction.api/src/core/domain/values"
	"testing"
)

func TestGenerateCompensation_shouldGenerateSuccessfully(t *testing.T) {
	transactionID := uuid.NewString()
	transactionStatusRepoTx := &mocks.TransactionStatusRepository{}
	balanceProvisionRepoTx := &mocks.BalanceProvisionRepository{}
	service := NewGenerateCompensationTxService()

	transactionStatus := &entity.TransactionStatus{
		TransactionStatusID: uuid.NewString(),
		Status:              values.TransactionStatusOpen,
	}

	balanceProvisions := entity.BalanceProvisions{
		{
			ProvisionID:          uuid.NewString(),
			Value:                100,
			OriginAccountID:      uuid.NewString(),
			DestinationAccountID: uuid.NewString(),
			TransactionID:        uuid.NewString(),
			Type:                 values.ProvisionTypeAdd,
			Status:               values.ProvisionStatusOpen,
		},
	}

	balanceProvisionAddOpen := balanceProvisions[0]
	balanceProvisionAddOpen.Status = values.ProvisionStatusClosed
	transactionStatusUpdated := *transactionStatus
	transactionStatusUpdated.Status = values.TransactionStatusFailed

	voidBalanceProvision := entity.BalanceProvision{
		Value:                balanceProvisionAddOpen.Value,
		OriginAccountID:      balanceProvisionAddOpen.OriginAccountID,
		DestinationAccountID: balanceProvisionAddOpen.DestinationAccountID,
		Type:                 values.ProvisionTypeVoid,
		Status:               values.ProvisionStatusOpen,
		TransactionID:        balanceProvisionAddOpen.TransactionID,
	}

	transactionStatusRepoTx.On("FindByTransactionID", transactionID).Return(transactionStatus, nil)
	balanceProvisionRepoTx.On("FindByTransactionID", transactionID).Return(balanceProvisions, nil)
	balanceProvisionRepoTx.On("Update", balanceProvisionAddOpen).Return(&balanceProvisionAddOpen, nil)
	transactionStatusRepoTx.On("Update", transactionStatusUpdated).Return(&transactionStatusUpdated, nil)
	balanceProvisionRepoTx.On("Create", voidBalanceProvision).Return(&voidBalanceProvision, nil)

	err := service.Generate(transactionID, transactionStatusRepoTx, balanceProvisionRepoTx)

	assert.Nil(t, err)
}

func TestGenerateCompensation_shouldNotGenerateWhenTransactionStatusNotFound(t *testing.T) {
	transactionID := uuid.NewString()
	transactionStatusRepoTx := &mocks.TransactionStatusRepository{}
	balanceProvisionRepoTx := &mocks.BalanceProvisionRepository{}
	service := NewGenerateCompensationTxService()

	var transactionStatus *entity.TransactionStatus

	transactionStatusRepoTx.On("FindByTransactionID", transactionID).Return(transactionStatus, nil)

	err := service.Generate(transactionID, transactionStatusRepoTx, balanceProvisionRepoTx)

	assert.Equal(t, values.NewErrorValidation("transaction status not found"), err)
	balanceProvisionRepoTx.AssertNumberOfCalls(t, "FindByTransactionID", 0)
	balanceProvisionRepoTx.AssertNumberOfCalls(t, "Update", 0)
	transactionStatusRepoTx.AssertNumberOfCalls(t, "Update", 0)
	balanceProvisionRepoTx.AssertNumberOfCalls(t, "Create", 0)
}

func TestGenerateCompensation_shouldNotGenerateWhenTransactionStatusNotOpen(t *testing.T) {
	transactionID := uuid.NewString()
	transactionStatusRepoTx := &mocks.TransactionStatusRepository{}
	balanceProvisionRepoTx := &mocks.BalanceProvisionRepository{}
	service := NewGenerateCompensationTxService()

	transactionStatus := &entity.TransactionStatus{
		Status: values.TransactionStatusBooked,
	}

	transactionStatusRepoTx.On("FindByTransactionID", transactionID).Return(transactionStatus, nil)

	err := service.Generate(transactionID, transactionStatusRepoTx, balanceProvisionRepoTx)

	assert.Equal(t, values.NewErrorValidation("cannot complete a transaction when not OPEN"), err)
	balanceProvisionRepoTx.AssertNumberOfCalls(t, "FindByTransactionID", 0)
	balanceProvisionRepoTx.AssertNumberOfCalls(t, "Update", 0)
	transactionStatusRepoTx.AssertNumberOfCalls(t, "Update", 0)
	balanceProvisionRepoTx.AssertNumberOfCalls(t, "Create", 0)
}

func TestGenerateCompensation_shouldNotGenerateWhenNoBalanceProvisionOpen(t *testing.T) {
	transactionID := uuid.NewString()
	transactionStatusRepoTx := &mocks.TransactionStatusRepository{}
	balanceProvisionRepoTx := &mocks.BalanceProvisionRepository{}
	service := NewGenerateCompensationTxService()

	transactionStatus := &entity.TransactionStatus{
		Status: values.TransactionStatusOpen,
	}
	balanceProvisions := entity.BalanceProvisions{
		{
			ProvisionID:          uuid.NewString(),
			Value:                100,
			OriginAccountID:      uuid.NewString(),
			DestinationAccountID: uuid.NewString(),
			TransactionID:        uuid.NewString(),
			Type:                 values.ProvisionTypeAdd,
			Status:               values.ProvisionStatusClosed,
		},
	}

	transactionStatusRepoTx.On("FindByTransactionID", transactionID).Return(transactionStatus, nil)
	balanceProvisionRepoTx.On("FindByTransactionID", transactionID).Return(balanceProvisions, nil)

	err := service.Generate(transactionID, transactionStatusRepoTx, balanceProvisionRepoTx)

	assert.Equal(t, values.NewErrorValidation("no provision found to complete"), err)
	balanceProvisionRepoTx.AssertNumberOfCalls(t, "Update", 0)
	transactionStatusRepoTx.AssertNumberOfCalls(t, "Update", 0)
	balanceProvisionRepoTx.AssertNumberOfCalls(t, "Create", 0)
}
