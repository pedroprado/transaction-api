package transactionRepo

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"pedroprado.transaction.api/src/core/_interfaces"
	"pedroprado.transaction.api/src/core/domain/entity"
	"pedroprado.transaction.api/src/core/domain/values"
	"pedroprado.transaction.api/src/repository/errorHandler"
	"pedroprado.transaction.api/src/repository/model"
	"time"
)

type transactionRepository struct {
	db           *gorm.DB
	errorHandler errorHandler.ErrorHandler
}

func NewTransactionRepository(
	db *gorm.DB,
	errorHandler errorHandler.ErrorHandler,
) _interfaces.TransactionRepository {
	return &transactionRepository{
		db:           db,
		errorHandler: errorHandler,
	}
}

func (ref *transactionRepository) Get(transactionID string) (*entity.Transaction, error) {
	var record model.Transaction

	result := ref.db.Find(&record, "transaction_id = ?", transactionID)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return record.ToDomain(), nil
}

func (ref *transactionRepository) Create(transaction entity.Transaction) (*entity.Transaction, error) {
	record := model.NewTransactionFromDomain(transaction)
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()

	if result := ref.db.Create(&record); result.Error != nil {
		if ref.errorHandler.IsRecordDuplicated(result) {
			transaction, err := ref.Get(record.TransactionID)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return nil, values.NewErrorDuplicated(*transaction)
		}

		if ref.errorHandler.ViolatesFKConstraint(result) {
			return nil, values.NewErrorValidation(fmt.Sprintf("transaction %v does not exist", record.TransactionID))
		}

		return nil, errors.WithStack(result.Error)
	}

	return record.ToDomain(), nil
}

func (ref *transactionRepository) WithTransaction(tx *gorm.DB) _interfaces.TransactionRepository {
	return NewTransactionRepository(tx, ref.errorHandler)
}
