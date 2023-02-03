package account

import (
	"pedroprado.transaction.api/src/core/_interfaces"
	"pedroprado.transaction.api/src/core/domain/entity"
)

type accountService struct {
	accountRepository _interfaces.AccountRepository
}

func NewAccountService(accountRepository _interfaces.AccountRepository) _interfaces.AccountService {
	return &accountService{accountRepository: accountRepository}
}

func (ref *accountService) Get(accountID string) (*entity.Account, error) {
	return ref.accountRepository.Get(accountID)
}

func (ref *accountService) Create(account entity.Account) (*entity.Account, error) {
	return ref.accountRepository.Create(account)
}
