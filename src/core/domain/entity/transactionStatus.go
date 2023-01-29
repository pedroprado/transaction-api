package entity

import "pedroprado.transaction.api/src/core/domain/values"

type TransactionStatus struct {
	TransactionStatusID string
	Status              values.TransactionStatus
	TransactionID       string
}
