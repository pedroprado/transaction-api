package accountRepo

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"pedroprado.transaction.api/src/core/_interfaces"
	"pedroprado.transaction.api/src/core/domain/entity"
	"pedroprado.transaction.api/src/core/domain/values"
	"pedroprado.transaction.api/src/repository/errorHandler"
	"pedroprado.transaction.api/src/repository/model"
	"time"
)

type accountRepository struct {
	db           *gorm.DB
	errorHandler errorHandler.ErrorHandler
}

func NewAccountRepository(db *gorm.DB, errorHandler errorHandler.ErrorHandler) _interfaces.AccountRepository {
	return &accountRepository{db: db, errorHandler: errorHandler}
}

func (ref *accountRepository) Get(accountID string) (*entity.Account, error) {
	var record model.Account

	result := ref.db.Find(&record, "account_id = ?", accountID)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return record.ToDomain(), nil
}

func (ref *accountRepository) Create(account entity.Account) (*entity.Account, error) {
	record := model.NewAccountRecordFromDomain(account)
	record.AccountID = uuid.NewString()
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()

	if result := ref.db.Create(&record); result.Error != nil {
		if ref.errorHandler.IsRecordDuplicated(result) {
			account, err := ref.Get(record.AccountID)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return nil, values.NewErrorDuplicated(*account)
		}

		return nil, errors.WithStack(result.Error)
	}

	return ref.Get(record.AccountID)
}

func (ref *accountRepository) GetLock(accountID string) (*entity.Account, error) {
	var record model.Account

	result := ref.db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&record, "account_id = ?", accountID)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return record.ToDomain(), nil
}

func (ref *accountRepository) Update(account entity.Account) (*entity.Account, error) {
	record := model.NewAccountRecordFromDomain(account)
	record.UpdatedAt = time.Now()

	result := ref.db.Model(&record).Updates(record)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, values.NewErrorNotFound("Account")
	}

	return ref.Get(account.AccountID)
}

func (ref *accountRepository) WithTransaction(tx *gorm.DB) _interfaces.AccountRepository {
	return NewAccountRepository(tx, ref.errorHandler)
}
