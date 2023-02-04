package transaction

import (
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
	transactionRepository        _interfaces.TransactionRepository
	createTransactionService     _interfaces.CreateTransactionService
	completeTransactionService   _interfaces.CompleteTransactionService
	compensateTransactionService _interfaces.CompensateTransactionService
}

func NewTransactionService(
	transactionRepository _interfaces.TransactionRepository,
	createTransactionService _interfaces.CreateTransactionService,
	completeTransactionService _interfaces.CompleteTransactionService,
	compensateTransactionService _interfaces.CompensateTransactionService,
) _interfaces.TransactionService {
	return &transactionService{
		transactionRepository:        transactionRepository,
		createTransactionService:     createTransactionService,
		completeTransactionService:   completeTransactionService,
		compensateTransactionService: compensateTransactionService,
	}
}

func (ref *transactionService) Get(transactionID string) (*entity.Transaction, error) {
	return ref.transactionRepository.Get(transactionID)
}

func (ref *transactionService) Create(transaction entity.Transaction) (*entity.Transaction, error) {
	return ref.createTransactionService.Create(transaction)
}

func (ref *transactionService) Complete(transactionID string) error {
	return ref.completeTransactionService.Complete(transactionID)
}

func (ref *transactionService) Compensate(transactionID string) error {
	return ref.compensateTransactionService.Compensate(transactionID)
}
