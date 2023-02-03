package balanceProvision

import (
	"pedroprado.transaction.api/src/core/_interfaces"
	"pedroprado.transaction.api/src/core/domain/entity"
)

type balanceProvisionService struct {
	balanceProvisionRepository _interfaces.BalanceProvisionRepository
}

func NewBalanceProvisionService(balanceProvisionRepository _interfaces.BalanceProvisionRepository) _interfaces.BalanceProvisionService {
	return &balanceProvisionService{
		balanceProvisionRepository: balanceProvisionRepository,
	}
}

func (ref *balanceProvisionService) FindByTransactionID(transactionID string) (entity.BalanceProvisions, error) {
	return ref.balanceProvisionRepository.FindByTransactionID(transactionID)
}
