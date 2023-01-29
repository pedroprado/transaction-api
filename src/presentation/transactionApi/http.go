package transactionApi

import "github.com/gin-gonic/gin"

type transactionApi struct {
}

func RegisterTransactionApi(ginRouterGroup *gin.RouterGroup) {
	transactionApi := transactionApi{}

	ginRouterGroup.POST("/transaction", transactionApi.Create)
}

func (ref *transactionApi) Create(c *gin.Context) {

	c.JSON(201, map[string]string{"result": "OK"})
}
