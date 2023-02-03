package _interfaces

import (
	"gorm.io/gorm"
	"pedroprado.transaction.api/src/core/domain/entity"
)

type BalanceProvisionService interface {
	FindByTransactionID(transactionID string) (entity.BalanceProvisions, error)
}

type BalanceProvisionRepository interface {
	Get(balanceProvisionID string) (*entity.BalanceProvision, error)
	FindByTransactionID(transactionID string) (entity.BalanceProvisions, error)
	Create(balanceProvision entity.BalanceProvision) (*entity.BalanceProvision, error)
	Update(balanceProvision entity.BalanceProvision) (*entity.BalanceProvision, error)
	WithTransaction(tx *gorm.DB) BalanceProvisionRepository
}
