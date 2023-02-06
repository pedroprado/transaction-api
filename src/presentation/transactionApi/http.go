package transactionApi

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"pedroprado.transaction.api/src/core/_interfaces"
	rest "pedroprado.transaction.api/src/presentation"
	"pedroprado.transaction.api/src/presentation/responses"
)

type transactionApi struct {
	transactionService _interfaces.TransactionService
}

func RegisterTransactionApi(ginRouterGroup *gin.RouterGroup, transactionService _interfaces.TransactionService) {
	transactionApi := transactionApi{
		transactionService: transactionService,
	}

	ginRouterGroup.POST("/transactions", transactionApi.CreateTransaction)
	ginRouterGroup.POST("/transaction/:transaction_id/complete", transactionApi.CompleteTransaction)
	ginRouterGroup.POST("/transaction/:transaction_id/compensate", transactionApi.CompensateTransaction)
}

// CreateTransaction Create a Transaction godoc
// @Summary Create a Transaction
// @Description Create a Transaction
// @Tags Transaction
// @Produce json
// @Param transaction body CreateTransactionRequest true "Body"
// @Success 201 {object} responses.Transaction
// @Failure 400 {object} rest.ErrorResponse
// @Failure 500 {object} rest.ErrorResponse
// @Router /transactions [post]
func (ref *transactionApi) CreateTransaction(c *gin.Context) {
	var request CreateTransactionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		rest.SendBadRequestError(c, err)
		return
	}

	created, err := ref.transactionService.Create(request.ToDomain())
	if err != nil {
		rest.NewErrorHandler(errors.WithStack(err)).Handle(c)
		return
	}

	c.JSON(201, responses.TransactionFromDomain(*created))
}

// CompleteTransaction Complete a Transaction godoc
// @Summary Complete a Transaction
// @Description Complete a Transaction
// @Tags Transaction
// @Produce json
// @Param transaction_id path string true "Account ID"
// @Success 202
// @Failure 400 {object} rest.ErrorResponse
// @Failure 500 {object} rest.ErrorResponse
// @Router /transaction/{transaction_id}/complete [post]
func (ref *transactionApi) CompleteTransaction(c *gin.Context) {
	var request CompleteTransactionRequest
	if err := c.ShouldBindUri(&request); err != nil {
		rest.SendBadRequestError(c, err)
		return
	}

	err := ref.transactionService.Complete(request.TransactionID)
	if err != nil {
		rest.NewErrorHandler(errors.WithStack(err)).Handle(c)
		return
	}

	c.Status(202)
}

// CompensateTransaction Compensate a Transaction godoc
// @Summary Compensate a Transaction
// @Description Compensate a Transaction
// @Tags Transaction
// @Produce json
// @Param transaction_id path string true "Account ID"
// @Success 202
// @Failure 400 {object} rest.ErrorResponse
// @Failure 500 {object} rest.ErrorResponse
// @Router /transaction/{transaction_id}/compensate [post]
func (ref *transactionApi) CompensateTransaction(c *gin.Context) {
	var request CompensateTransactionRequest
	if err := c.ShouldBindUri(&request); err != nil {
		rest.SendBadRequestError(c, err)
		return
	}

	err := ref.transactionService.Compensate(request.TransactionID)
	if err != nil {
		rest.NewErrorHandler(errors.WithStack(err)).Handle(c)
		return
	}

	c.Status(202)
}
