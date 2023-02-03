package entity

import (
	"pedroprado.transaction.api/src/core/domain/values"
	"time"
)

type TransactionStatus struct {
	TransactionStatusID string
	Status              values.TransactionStatus
	TransactionID       string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
