// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	_interfaces "pedroprado.transaction.api/src/core/_interfaces"
)

// CompleteTransactionTxService is an autogenerated mock type for the CompleteTransactionTxService type
type CompleteTransactionTxService struct {
	mock.Mock
}

// Complete provides a mock function with given fields: transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx
func (_m *CompleteTransactionTxService) Complete(transactionID string, transactionRepoTx _interfaces.TransactionRepository, transactionStatusRepoTx _interfaces.TransactionStatusRepository, balanceProvisionRepoTx _interfaces.BalanceProvisionRepository, accountRepoTx _interfaces.AccountRepository) error {
	ret := _m.Called(transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, _interfaces.TransactionRepository, _interfaces.TransactionStatusRepository, _interfaces.BalanceProvisionRepository, _interfaces.AccountRepository) error); ok {
		r0 = rf(transactionID, transactionRepoTx, transactionStatusRepoTx, balanceProvisionRepoTx, accountRepoTx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewCompleteTransactionTxService interface {
	mock.TestingT
	Cleanup(func())
}

// NewCompleteTransactionTxService creates a new instance of CompleteTransactionTxService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCompleteTransactionTxService(t mockConstructorTestingTNewCompleteTransactionTxService) *CompleteTransactionTxService {
	mock := &CompleteTransactionTxService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
