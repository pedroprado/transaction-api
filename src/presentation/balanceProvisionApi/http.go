package balanceProvisionApi

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"pedroprado.transaction.api/src/core/_interfaces"
	rest "pedroprado.transaction.api/src/presentation"
	"pedroprado.transaction.api/src/presentation/responses"
)

type balanceProvisionApi struct {
	balanceProvisionService _interfaces.BalanceProvisionService
}

func RegisterBalanceProvisionApi(ginRouterGroup *gin.RouterGroup, balanceProvisionService _interfaces.BalanceProvisionService) {
	balanceProvisionApi := balanceProvisionApi{
		balanceProvisionService: balanceProvisionService,
	}

	ginRouterGroup.GET("/balance_provisions", balanceProvisionApi.Search)

}

// Search BalanceProvisions for a Transaction godoc
// @Summary BalanceProvisions for a Transaction
// @Description BalanceProvisions a Transaction
// @Tags BalanceProvisions
// @Produce json
// @Param transaction_id query string true "Transaction ID"
// @Success 200 {object}
// @Failure 400 {object} rest.ErrorResponse
// @Failure 500 {object} rest.ErrorResponse
// @Router /balance_provisions [get]
func (ref *balanceProvisionApi) Search(c *gin.Context) {
	var request SearchBalanceProvisionsRequest
	if err := c.ShouldBindQuery(request); err != nil {
		rest.SendBadRequestError(c, err)
		return
	}

	balanceProvisions, err := ref.balanceProvisionService.FindByTransactionID(request.TransactionID)
	if err != nil {
		rest.NewErrorHandler(errors.WithStack(err)).Handle(c)
		return
	}

	c.JSON(200, responses.BalanceProvisionsFromDomain(balanceProvisions))
}
