package main

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"pedroprado.transaction.api/src/infra"
	rest "pedroprado.transaction.api/src/presentation"
	"pedroprado.transaction.api/src/presentation/transactionApi"
	"pedroprado.transaction.api/src/repository"
)

var (
	sslMode = false
)

const (
	apiPort      = "8098"
	relativePath = "/snapfi"
)

func main() {

	ginServer := rest.NewServerHttpGin()
	ginRouterGroup := ginServer.GetGinRouterGroup(relativePath)
	transactionApi.RegisterTransactionApi(ginRouterGroup)
	rest.RegisterInfraApi(ginRouterGroup, false)
	if err := ginServer.StartServer(apiPort); err != nil {
		logrus.Fatal(errors.WithStack(err))
	}
}

func getDb() *gorm.DB {
	db, err := infra.OpenConnection(sslMode)
	if err != nil {
		logrus.Fatal(errors.WithStack(err))
	}
	err = repository.Migrate(db)
	if err != nil {
		logrus.Fatal(errors.WithStack(err))
	}
	return db
}
