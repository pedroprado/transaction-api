package model

import "time"

type TransactionStatus struct {
	TransactionStatusID string `gorm:"primaryKey"`
	Status              string
	TransactionID       string `gorm:"uniqueIndex:idx_transaction_status_transaction_id"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
