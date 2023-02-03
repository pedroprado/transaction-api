package transactionApi

import "github.com/gin-gonic/gin"

type transactionApi struct {
}

func RegisterTransactionApi(ginRouterGroup *gin.RouterGroup) {
	transactionApi := transactionApi{}

	ginRouterGroup.POST("/transactions", transactionApi.CreateTransaction)
	ginRouterGroup.POST("/transaction/:transaction_id/complete", transactionApi.CompleteTransaction)
	ginRouterGroup.POST("/transaction/:transaction_id/compensate", transactionApi.CompensateTransaction)
}

func (ref *transactionApi) CreateTransaction(c *gin.Context) {

	c.JSON(201, map[string]string{"result": "OK"})
}

func (ref *transactionApi) CompleteTransaction(c *gin.Context) {

	c.Status(200)
}

func (ref *transactionApi) CompensateTransaction(c *gin.Context) {

	c.Status(200)
}
