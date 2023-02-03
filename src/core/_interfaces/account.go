package _interfaces

import (
	"gorm.io/gorm"
	"pedroprado.transaction.api/src/core/domain/entity"
)

type AccountService interface {
	Get(accountID string) (*entity.Account, error)
	Create(account entity.Account) (*entity.Account, error)
}

type AccountRepository interface {
	Get(accountID string) (*entity.Account, error)
	GetLock(accountID string) (*entity.Account, error)
	Create(account entity.Account) (*entity.Account, error)
	Update(account entity.Account) (*entity.Account, error)
	WithTransaction(tx *gorm.DB) AccountRepository
}
