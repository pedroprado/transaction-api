package model

import (
	"pedroprado.transaction.api/src/core/domain/entity"
	"time"
)

type Account struct {
	AccountID string `gorm:"primaryKey"`
	Bank      string
	Number    string
	Agency    string
	Balance   float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAccountRecordFromDomain(account entity.Account) Account {
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

func (record Account) ToDomain() *entity.Account {
	return &entity.Account{
		AccountID: record.AccountID,
		Bank:      record.Bank,
		Number:    record.Number,
		Agency:    record.Agency,
		Balance:   record.Balance,
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
	}
}
