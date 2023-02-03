package entity

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"pedroprado.transaction.api/src/core/domain/values"
	"testing"
)

func TestFindBalanceProvisionOpen(t *testing.T) {
	balanceProvisions := BalanceProvisions{
		{
			ProvisionID: uuid.NewString(),
			Type:        values.ProvisionTypeAdd,
			Status:      values.ProvisionStatusOpen,
		},
		{
			ProvisionID: uuid.NewString(),
			Status:      values.ProvisionStatusClosed,
			Type:        values.ProvisionTypeVoid,
		},
	}

	expected := &balanceProvisions[0]

	received := balanceProvisions.FindProvisionToComplete()

	assert.Equal(t, expected, received)
}
