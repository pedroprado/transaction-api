package responses

import (
	"pedroprado.transaction.api/src/core/domain/entity"
	"time"
)

type TransactionStatus struct {
	TransactionStatusID string    `json:"transaction_status_id"`
	Status              string    `json:"status"`
	TransactionID       string    `json:"transaction_id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func TransactionStatusFromDomain(transactionStatus entity.TransactionStatus) TransactionStatus {
	return TransactionStatus{
		TransactionStatusID: transactionStatus.TransactionStatusID,
		Status:              string(transactionStatus.Status),
		TransactionID:       transactionStatus.TransactionID,
		CreatedAt:           transactionStatus.CreatedAt,
		UpdatedAt:           transactionStatus.UpdatedAt,
	}
}
