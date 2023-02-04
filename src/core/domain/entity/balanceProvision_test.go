package entity

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"pedroprado.transaction.api/src/core/domain/values"
	"testing"
)

func TestFindBalanceProvision(t *testing.T) {
	balanceProvisions := BalanceProvisions{
		{
			ProvisionID: uuid.NewString(),
			Type:        values.ProvisionTypeAdd,
			Status:      values.ProvisionStatusOpen,
		},
		{
			ProvisionID: uuid.NewString(),
			Type:        values.ProvisionTypeAdd,
			Status:      values.ProvisionStatusClosed,
		},
		{
			ProvisionID: uuid.NewString(),
			Type:        values.ProvisionTypeVoid,
			Status:      values.ProvisionStatusOpen,
		},
		{
			ProvisionID: uuid.NewString(),
			Type:        values.ProvisionTypeVoid,
			Status:      values.ProvisionStatusClosed,
		},
	}

	cases := map[string]struct {
		provisionType   values.ProvisionType
		provisionStatus values.ProvisionStatus
		expected        *BalanceProvision
	}{
		"should-find-add-open": {
			provisionType:   values.ProvisionTypeAdd,
			provisionStatus: values.ProvisionStatusOpen,
			expected:        &balanceProvisions[0],
		},
		"should-find-add-closed": {
			provisionType:   values.ProvisionTypeAdd,
			provisionStatus: values.ProvisionStatusClosed,
			expected:        &balanceProvisions[1],
		},
		"should-find-void-open": {
			provisionType:   values.ProvisionTypeVoid,
			provisionStatus: values.ProvisionStatusOpen,
			expected:        &balanceProvisions[2],
		},
		"should-find-void-closed": {
			provisionType:   values.ProvisionTypeVoid,
			provisionStatus: values.ProvisionStatusClosed,
			expected:        &balanceProvisions[3],
		},
		"should-not-find-any": {
			expected: nil,
		},
	}

	for title, tc := range cases {
		t.Run(title, func(t *testing.T) {
			received := balanceProvisions.FindProvision(tc.provisionType, tc.provisionStatus)

			assert.Equal(t, tc.expected, received)
		})
	}
}
