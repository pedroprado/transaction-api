package accountApi

import (
	"pedroprado.transaction.api/src/core/domain/entity"
)

type GetAccountRequest struct {
	AccountID string `uri:"account_id" binding:"required"`
}

type CreateAccountRequest struct {
	Bank    string  `json:"bank"`
	Number  string  `json:"number"`
	Agency  string  `json:"agency"`
	Balance float64 `json:"balance"`
}

func (request CreateAccountRequest) ToDomain() entity.Account {
	return entity.Account{
		Bank:    request.Bank,
		Number:  request.Number,
		Agency:  request.Agency,
		Balance: request.Balance,
	}
}

type PatchAccountRequest struct {
	AccountID string  `json:"account_id" binding:"required"`
	Balance   float64 `json:"balance" binding:"required"`
}
