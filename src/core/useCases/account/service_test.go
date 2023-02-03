package account

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	mocks "pedroprado.transaction.api/src/core/_mocks"
	"pedroprado.transaction.api/src/core/domain/entity"
	"testing"
)

func TestGet(t *testing.T) {
	accountRepository := &mocks.AccountRepository{}
	service := NewAccountService(accountRepository)

	accountID := uuid.NewString()
	account := &entity.Account{AccountID: uuid.NewString()}

	accountRepository.On("Get", accountID).Return(account, nil)

	expected := account
	received, err := service.Get(accountID)

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestCreate(t *testing.T) {
	accountRepository := &mocks.AccountRepository{}
	service := NewAccountService(accountRepository)

	account := entity.Account{AccountID: uuid.NewString()}
	created := &entity.Account{AccountID: uuid.NewString()}

	accountRepository.On("Create", account).Return(created, nil)

	expected := created
	received, err := service.Create(account)

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}
