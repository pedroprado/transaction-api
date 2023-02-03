package accountApi

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"pedroprado.transaction.api/src/core/_interfaces"
	rest "pedroprado.transaction.api/src/presentation"
	"pedroprado.transaction.api/src/presentation/responses"
)

type accountApi struct {
	accountService _interfaces.AccountService
}

func RegisterAccountApi(ginRouterGroup *gin.RouterGroup, accountService _interfaces.AccountService) {
	accountApi := accountApi{accountService: accountService}

	ginRouterGroup.GET("/account/:account_id", accountApi.Get)
	ginRouterGroup.POST("/account", accountApi.Create)
}

// Get Account godoc
// @Summary Get an Account by id
// @Description Get an Account by id
// @Tags Account
// @Produce json
// @Param account_id path string true "Account ID"
// @Success 200 {object} responses.Account
// @Success 204
// @Failure 400 {object} rest.ErrorResponse
// @Failure 500 {object} rest.ErrorResponse
// @Router /account/{account_id} [get]
func (ref *accountApi) Get(c *gin.Context) {
	var request GetAccountRequest
	if err := c.ShouldBindUri(&request); err != nil {
		rest.SendBadRequestError(c, err)
		return
	}

	account, err := ref.accountService.Get(request.AccountID)
	if err != nil {
		rest.NewErrorHandler(errors.WithStack(err)).Handle(c)
		return
	}
	if account == nil {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, responses.AccountFromDomain(*account))
}

// Create Account godoc
// @Summary Create an Account
// @Description Create an Account
// @Tags Account
// @Produce json
// @Success 201 {object} responses.Account
// @Failure 400 {object} rest.ErrorResponse
// @Failure 500 {object} rest.ErrorResponse
// @Router /account/{account_id} [get]
func (ref *accountApi) Create(c *gin.Context) {
	var request CreateAccountRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		rest.SendBadRequestError(c, err)
		return
	}

	created, err := ref.accountService.Create(request.ToDomain())
	if err != nil {
		rest.NewErrorHandler(errors.WithStack(err)).Handle(c)
		return
	}

	c.JSON(201, responses.AccountFromDomain(*created))
}
