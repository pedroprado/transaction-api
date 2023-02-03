package transactionStatusApi

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"pedroprado.transaction.api/src/core/_interfaces"
	rest "pedroprado.transaction.api/src/presentation"
	"pedroprado.transaction.api/src/presentation/responses"
)

type transactionStatusApi struct {
	transactionStatusService _interfaces.TransactionStatusService
}

func RegisterTransactionStatusApi(ginRouterGroup *gin.RouterGroup, transactionStatusService _interfaces.TransactionStatusService) {
	transactionStatusApi := transactionStatusApi{
		transactionStatusService: transactionStatusService,
	}

	ginRouterGroup.GET("/transaction_status", transactionStatusApi.Search)
}

// Search BalanceProvisions for a Transaction godoc
// @Summary BalanceProvisions for a Transaction
// @Description BalanceProvisions a Transaction
// @Tags BalanceProvisions
// @Produce json
// @Param transaction_id query string true "Transaction ID"
// @Success 200 {object} responses.TransactionStatus
// @Failure 400 {object} rest.ErrorResponse
// @Failure 500 {object} rest.ErrorResponse
// @Router /transaction_status [get]
func (ref *transactionStatusApi) Search(c *gin.Context) {
	var request SearchTransactionStatusRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		rest.SendBadRequestError(c, err)
		return
	}

	transactionStatus, err := ref.transactionStatusService.FindByTransactionID(request.TransactionID)
	if err != nil {
		rest.NewErrorHandler(errors.WithStack(err)).Handle(c)
		return
	}
	if transactionStatus == nil {
		c.Status(204)
		return
	}

	c.JSON(200, responses.TransactionStatusFromDomain(*transactionStatus))
}
