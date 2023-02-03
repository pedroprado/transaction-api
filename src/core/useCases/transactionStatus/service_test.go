package transactionStatus

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	mocks "pedroprado.transaction.api/src/core/_mocks"
	"pedroprado.transaction.api/src/core/domain/entity"
	"testing"
)

func TestTransactionStatusService_FindByTransactionID(t *testing.T) {
	transactionStatusRepo := &mocks.TransactionStatusRepository{}
	service := NewTransactionStatusService(transactionStatusRepo)

	transactionID := uuid.NewString()
	transactionStatus := &entity.TransactionStatus{TransactionStatusID: uuid.NewString()}

	transactionStatusRepo.On("FindByTransactionID", transactionID).Return(transactionStatus, nil)

	expected := transactionStatus
	received, err := service.FindByTransactionID(transactionID)

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}
