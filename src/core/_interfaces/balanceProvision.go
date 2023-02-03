package _interfaces

import (
	"gorm.io/gorm"
	"pedroprado.transaction.api/src/core/domain/entity"
)

type BalanceProvisionService interface {
	Get(balanceProvisionID string) (*entity.BalanceProvision, error)
	Create(balanceProvision entity.BalanceProvision) (*entity.BalanceProvision, error)
	Update(balanceProvision entity.BalanceProvision) (*entity.BalanceProvision, error)
}

type BalanceProvisionRepository interface {
	Get(balanceProvisionID string) (*entity.BalanceProvision, error)
	Create(balanceProvision entity.BalanceProvision) (*entity.BalanceProvision, error)
	Update(balanceProvision entity.BalanceProvision) (*entity.BalanceProvision, error)
	WithTransaction(tx *gorm.DB) BalanceProvisionRepository
}
