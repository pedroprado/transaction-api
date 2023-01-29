package model

import "time"

type Transaction struct {
	TransactionID        string `gorm:"primaryKey"`
	Type                 string
	OriginAccountID      string
	DestinationAccountID string
	Value                float64
	CreatedAt            time.Time
	UpdatedAt            time.Time

	TransactionStatus TransactionStatus  `gorm:"foreignKey:TransactionID"`
	BalanceProvisions []BalanceProvision `gorm:"foreignKey:TransactionID"`
}
