package completeTx

import (
	"github.com/sirupsen/logrus"
	"pedroprado.transaction.api/src/core/_interfaces"
	"pedroprado.transaction.api/src/core/domain/values"
)

type completeTransactionTxService struct {
}

func NewCompleteTransactionTxService() _interfaces.CompleteTransactionTxService {
	return &completeTransactionTxService{}
}

func (ref *completeTransactionTxService) Complete(
	transactionID string,
	transactionRepoTx _interfaces.TransactionRepository,
	transactionStatusRepoTx _interfaces.TransactionStatusRepository,
	balanceProvisionRepoTx _interfaces.BalanceProvisionRepository,
	accountRepoTx _interfaces.AccountRepository,
) error {
	transaction, err := transactionRepoTx.GetLock(transactionID)
	if err != nil {
		return err
	}
	if transaction == nil {
		return values.NewErrorValidation("transaction not found")
	}

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
		logrus.Errorf("no provision found to complete for transaction %s", transactionID)

		return values.NewErrorValidation("no provision found to complete")
	}

	intermediaryAccount, err := accountRepoTx.GetLock(values.IntermediaryAccountID)
	if err != nil {
		return err
	}
	if intermediaryAccount == nil {
		return values.ErrorIntermediaryAccountNotFound
	}

	destinationAccount, err := accountRepoTx.GetLock(balanceProvisionAddOpen.DestinationAccountID)
	if err != nil {
		return err
	}
	if destinationAccount == nil {
		return values.ErrorDestinationAccountNotFound
	}

	intermediaryAccount, err = intermediaryAccount.RemoveFunds(balanceProvisionAddOpen.Value)
	if err != nil {
		return err
	}
	destinationAccount = destinationAccount.AddFunds(balanceProvisionAddOpen.Value)

	balanceProvisionAddOpen.Status = values.ProvisionStatusClosed
	transactionStatus.Status = values.TransactionStatusBooked

	_, err = accountRepoTx.Update(*intermediaryAccount)
	if err != nil {
		return err
	}

	_, err = accountRepoTx.Update(*destinationAccount)
	if err != nil {
		return err
	}

	_, err = balanceProvisionRepoTx.Update(*balanceProvisionAddOpen)
	if err != nil {
		return err
	}

	_, err = transactionStatusRepoTx.Update(*transactionStatus)
	if err != nil {
		return err
	}

	return nil
}
