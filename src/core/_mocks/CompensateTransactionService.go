// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// CompensateTransactionService is an autogenerated mock type for the CompensateTransactionService type
type CompensateTransactionService struct {
	mock.Mock
}

// Compensate provides a mock function with given fields: transactionID
func (_m *CompensateTransactionService) Compensate(transactionID string) error {
	ret := _m.Called(transactionID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(transactionID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewCompensateTransactionService interface {
	mock.TestingT
	Cleanup(func())
}

// NewCompensateTransactionService creates a new instance of CompensateTransactionService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCompensateTransactionService(t mockConstructorTestingTNewCompensateTransactionService) *CompensateTransactionService {
	mock := &CompensateTransactionService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
