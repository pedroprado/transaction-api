package generateCompensationTx

import (
	"pedroprado.transaction.api/src/core/_interfaces"
	"pedroprado.transaction.api/src/core/domain/entity"
	"pedroprado.transaction.api/src/core/domain/values"
)

type generateCompensationTxService struct{}

func NewGenerateCompensationTxService() _interfaces.GenerateCompensationTxService {
	return &generateCompensationTxService{}
}

func (ref *generateCompensationTxService) Generate(
	transactionID string,
	transactionStatusRepoTx _interfaces.TransactionStatusRepository,
	balanceProvisionRepoTx _interfaces.BalanceProvisionRepository,
) error {
	transactionStatus, err := transactionStatusRepoTx.FindByTransactionID(transactionID)
	if err != nil {
		return err
	}
	if transactionStatus == nil {
		return values.NewErrorValidation("transaction status not found")
	}
	if transactionStatus.Status != values.TransactionStatusOpen {
		return values.NewErrorValidation("cannot complete a transaction when not OPEN")
	}

	balanceProvisions, err := balanceProvisionRepoTx.FindByTransactionID(transactionID)
	if err != nil {
		return err
	}
	balanceProvisionAddOpen := balanceProvisions.FindProvision(values.ProvisionTypeAdd, values.ProvisionStatusOpen)
	if balanceProvisionAddOpen == nil {
		return values.NewErrorValidation("no provision found to complete")
	}

	balanceProvisionAddOpen.Status = values.ProvisionStatusClosed
	transactionStatus.Status = values.TransactionStatusFailed

	_, err = balanceProvisionRepoTx.Update(*balanceProvisionAddOpen)
	if err != nil {
		return err
	}

	_, err = transactionStatusRepoTx.Update(*transactionStatus)
	if err != nil {
		return err
	}

	voidBalanceProvision := entity.BalanceProvision{
		Value:                balanceProvisionAddOpen.Value,
		OriginAccountID:      balanceProvisionAddOpen.OriginAccountID,
		DestinationAccountID: balanceProvisionAddOpen.DestinationAccountID,
		Type:                 values.ProvisionTypeVoid,
		Status:               values.ProvisionStatusOpen,
		TransactionID:        balanceProvisionAddOpen.TransactionID,
	}

	_, err = balanceProvisionRepoTx.Create(voidBalanceProvision)
	if err != nil {
		return err
	}

	return nil
}
