package model

import "time"

type BalanceProvision struct {
	ProvisionID          string
	Value                float64
	OriginAccountID      string
	DestinationAccountID string
	Type                 string
	Status               string
	TransactionID        string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
