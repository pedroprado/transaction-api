package transactionStatus

import (
	"pedroprado.transaction.api/src/core/_interfaces"
	"pedroprado.transaction.api/src/core/domain/entity"
)

type transactionStatusService struct {
	transactionStatusRepository _interfaces.TransactionStatusRepository
}

func NewTransactionStatusService(transactionStatusRepository _interfaces.TransactionStatusRepository) _interfaces.TransactionStatusService {
	return &transactionStatusService{
		transactionStatusRepository: transactionStatusRepository,
	}
}

func (ref *transactionStatusService) FindByTransactionID(transactionID string) (*entity.TransactionStatus, error) {
	return ref.transactionStatusRepository.FindByTransactionID(transactionID)
}
