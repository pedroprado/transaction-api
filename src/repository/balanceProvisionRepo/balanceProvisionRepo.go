package balanceProvisionRepo

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

type balanceProvisionRepository struct {
	db           *gorm.DB
	errorHandler errorHandler.ErrorHandler
}

func NewBalanceProvisionRepository(
	db *gorm.DB,
	errorHandler errorHandler.ErrorHandler,
) _interfaces.BalanceProvisionRepository {
	return &balanceProvisionRepository{
		db:           db,
		errorHandler: errorHandler,
	}
}

func (ref *balanceProvisionRepository) Get(balanceProvisionID string) (*entity.BalanceProvision, error) {
	var record model.BalanceProvision

	result := ref.db.Find(&record, "provision_id = ?", balanceProvisionID)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return record.ToDomain(), nil
}

func (ref *balanceProvisionRepository) FindByTransactionID(transactionID string) (entity.BalanceProvisions, error) {
	var records []model.BalanceProvision

	result := ref.db.Find(&records, "transaction_id = ?", transactionID)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	balanceProvisions := make([]entity.BalanceProvision, len(records))
	for i := range records {
		record := records[i]
		balanceProvisions[i] = *record.ToDomain()
	}

	return balanceProvisions, nil
}

func (ref *balanceProvisionRepository) Create(balanceProvision entity.BalanceProvision) (*entity.BalanceProvision, error) {
	record := model.NewBalanceProvisionFromDomain(balanceProvision)
	record.ProvisionID = uuid.NewString()
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()

	if result := ref.db.Create(&record); result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	return record.ToDomain(), nil
}

func (ref *balanceProvisionRepository) Update(balanceProvision entity.BalanceProvision) (*entity.BalanceProvision, error) {
	record := model.NewBalanceProvisionFromDomain(balanceProvision)
	record.UpdatedAt = time.Now()

	result := ref.db.Model(&record).Updates(record)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, values.NewErrorNotFound("Balance Provision")
	}

	return record.ToDomain(), nil
}

func (ref *balanceProvisionRepository) WithTransaction(tx *gorm.DB) _interfaces.BalanceProvisionRepository {
	return NewBalanceProvisionRepository(tx, ref.errorHandler)

}
