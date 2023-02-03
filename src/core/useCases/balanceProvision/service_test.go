package balanceProvision

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	mocks "pedroprado.transaction.api/src/core/_mocks"
	"pedroprado.transaction.api/src/core/domain/entity"
	"testing"
)

func TestBalanceProvisionService_FindByTransactionID(t *testing.T) {
	balanceProvisionRepo := &mocks.BalanceProvisionRepository{}
	service := NewBalanceProvisionService(balanceProvisionRepo)

	transactionID := uuid.NewString()
	balanceProvisions := entity.BalanceProvisions{{ProvisionID: uuid.NewString()}}

	balanceProvisionRepo.On("FindByTransactionID", transactionID).Return(balanceProvisions, nil)

	expected := balanceProvisions
	received, err := service.FindByTransactionID(transactionID)

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}
