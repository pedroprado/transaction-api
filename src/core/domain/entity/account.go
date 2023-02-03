package entity

import (
	"pedroprado.transaction.api/src/core/domain/values"
	"time"
)

var (
	errorAccountDoesNotHaveFunds = values.NewErrorValidation("account does not have funds")
)

type Account struct {
	AccountID string
	Bank      string
	Number    string
	Agency    string
	Balance   float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (account Account) AddFunds(funds float64) *Account {
	account.Balance = account.Balance + funds
	return &account
}

func (account Account) RemoveFunds(funds float64) (*Account, error) {
	if (account.Balance - funds) < 0 {
		return nil, errorAccountDoesNotHaveFunds
	}
	account.Balance = account.Balance - funds
	return &account, nil
}
