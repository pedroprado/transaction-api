package responses

import (
	"pedroprado.transaction.api/src/core/domain/entity"
	"time"
)

type BalanceProvision struct {
	ProvisionID          string    `json:"provision_id"`
	Value                float64   `json:"value"`
	OriginAccountID      string    `json:"origin_account_id"`
	DestinationAccountID string    `json:"destination_account_id"`
	Type                 string    `json:"type"`
	Status               string    `json:"status"`
	TransactionID        string    `json:"transaction_id"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func BalanceProvisionFromDomain(balanceProvision entity.BalanceProvision) BalanceProvision {
	return BalanceProvision{
		ProvisionID:          balanceProvision.ProvisionID,
		Value:                balanceProvision.Value,
		OriginAccountID:      balanceProvision.OriginAccountID,
		DestinationAccountID: balanceProvision.DestinationAccountID,
		Type:                 string(balanceProvision.Type),
		Status:               string(balanceProvision.Status),
		TransactionID:        balanceProvision.TransactionID,
		CreatedAt:            balanceProvision.CreatedAt,
		UpdatedAt:            balanceProvision.UpdatedAt,
	}
}

func BalanceProvisionsFromDomain(balanceProvisions entity.BalanceProvisions) []BalanceProvision {
	responses := make([]BalanceProvision, len(balanceProvisions))
	for i := range balanceProvisions {
		responses[i] = BalanceProvisionFromDomain(balanceProvisions[i])
	}

	return responses
}