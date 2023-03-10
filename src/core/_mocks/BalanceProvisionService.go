// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	entity "pedroprado.transaction.api/src/core/domain/entity"
)

// BalanceProvisionService is an autogenerated mock type for the BalanceProvisionService type
type BalanceProvisionService struct {
	mock.Mock
}

// FindByTransactionID provides a mock function with given fields: transactionID
func (_m *BalanceProvisionService) FindByTransactionID(transactionID string) (entity.BalanceProvisions, error) {
	ret := _m.Called(transactionID)

	var r0 entity.BalanceProvisions
	if rf, ok := ret.Get(0).(func(string) entity.BalanceProvisions); ok {
		r0 = rf(transactionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(entity.BalanceProvisions)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(transactionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewBalanceProvisionService interface {
	mock.TestingT
	Cleanup(func())
}

// NewBalanceProvisionService creates a new instance of BalanceProvisionService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBalanceProvisionService(t mockConstructorTestingTNewBalanceProvisionService) *BalanceProvisionService {
	mock := &BalanceProvisionService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
