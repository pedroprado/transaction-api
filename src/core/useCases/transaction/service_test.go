package transaction

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	mocks "pedroprado.transaction.api/src/core/_mocks"
	"pedroprado.transaction.api/src/core/domain/entity"
	"testing"
)

func TestTransactionService_Get(t *testing.T) {
	transactionRepository := &mocks.TransactionRepository{}
	createTransactionService := &mocks.CreateTransactionService{}
	completeTransactionService := &mocks.CompleteTransactionService{}
	compensateTransactionService := &mocks.CompensateTransactionService{}
	service := NewTransactionService(transactionRepository, createTransactionService, completeTransactionService, compensateTransactionService)

	transactionID := uuid.NewString()
	transaction := &entity.Transaction{TransactionID: transactionID}

	transactionRepository.On("Get", transactionID).Return(transaction, nil)

	expected := transaction
	received, err := service.Get(transactionID)

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestTransactionService_Create(t *testing.T) {
	transactionRepository := &mocks.TransactionRepository{}
	createTransactionService := &mocks.CreateTransactionService{}
	completeTransactionService := &mocks.CompleteTransactionService{}
	compensateTransactionService := &mocks.CompensateTransactionService{}
	service := NewTransactionService(transactionRepository, createTransactionService, completeTransactionService, compensateTransactionService)

	transaction := entity.Transaction{TransactionID: uuid.NewString()}
	created := &entity.Transaction{TransactionID: uuid.NewString()}

	createTransactionService.On("Create", transaction).Return(created, nil)

	expected := created
	received, err := service.Create(transaction)

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestTransactionService_Complete(t *testing.T) {
	transactionRepository := &mocks.TransactionRepository{}
	createTransactionService := &mocks.CreateTransactionService{}
	completeTransactionService := &mocks.CompleteTransactionService{}
	compensateTransactionService := &mocks.CompensateTransactionService{}
	service := NewTransactionService(transactionRepository, createTransactionService, completeTransactionService, compensateTransactionService)

	transactionID := uuid.NewString()

	completeTransactionService.On("Complete", transactionID).Return(nil)

	err := service.Complete(transactionID)

	assert.Nil(t, err)
}

func TestTransactionService_Compensate(t *testing.T) {
	transactionRepository := &mocks.TransactionRepository{}
	createTransactionService := &mocks.CreateTransactionService{}
	completeTransactionService := &mocks.CompleteTransactionService{}
	compensateTransactionService := &mocks.CompensateTransactionService{}
	service := NewTransactionService(transactionRepository, createTransactionService, completeTransactionService, compensateTransactionService)

	transactionID := uuid.NewString()

	compensateTransactionService.On("Compensate", transactionID).Return(nil)

	err := service.Compensate(transactionID)

	assert.Nil(t, err)
}
