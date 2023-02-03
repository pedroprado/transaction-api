package responses

import (
	"pedroprado.transaction.api/src/core/domain/entity"
	"time"
)

type Account struct {
	AccountID string    `json:"account_id"`
	Bank      string    `json:"bank"`
	Number    string    `json:"number"`
	Agency    string    `json:"agency"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func AccountFromDomain(account entity.Account) Account {
	return Account{
		AccountID: account.AccountID,
		Bank:      account.Bank,
		Number:    account.Number,
		Agency:    account.Agency,
		Balance:   account.Balance,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}
}
