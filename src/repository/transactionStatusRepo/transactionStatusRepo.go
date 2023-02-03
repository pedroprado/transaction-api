package transactionStatusRepo

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"pedroprado.transaction.api/src/core/_interfaces"
	"pedroprado.transaction.api/src/core/domain/entity"
	"pedroprado.transaction.api/src/core/domain/values"
	"pedroprado.transaction.api/src/repository/errorHandler"
	"pedroprado.transaction.api/src/repository/model"
	"time"
)

type transactionStatusRepository struct {
	db           *gorm.DB
	errorHandler errorHandler.ErrorHandler
}

func NewTransactionStatusRepository(
	db *gorm.DB,
	errorHandler errorHandler.ErrorHandler,
) _interfaces.TransactionStatusRepository {
	return &transactionStatusRepository{
		db:           db,
		errorHandler: errorHandler,
	}
}

func (ref *transactionStatusRepository) FindByTransactionID(transactionID string) (*entity.TransactionStatus, error) {
	var record model.TransactionStatus

	result := ref.db.Find(&record, "transaction_id = ?", transactionID)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return record.ToDomain(), nil
}

func (ref *transactionStatusRepository) Create(transactionStatus entity.TransactionStatus) (*entity.TransactionStatus, error) {
	record := model.NewTransactionStatusFromDomain(transactionStatus)
	record.TransactionID = uuid.NewString()
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()

	if result := ref.db.Create(&record); result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	return record.ToDomain(), nil
}

func (ref *transactionStatusRepository) Update(transactionStatus entity.TransactionStatus) (*entity.TransactionStatus, error) {
	record := model.NewTransactionStatusFromDomain(transactionStatus)
	record.UpdatedAt = time.Now()

	result := ref.db.Model(&record).Updates(record)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, values.NewErrorNotFound("Transaction Status")
	}

	return record.ToDomain(), nil
}

func (ref *transactionStatusRepository) WithTransaction(tx *gorm.DB) _interfaces.TransactionStatusRepository {
	return NewTransactionStatusRepository(tx, ref.errorHandler)

}
