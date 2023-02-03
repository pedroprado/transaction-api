package main

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"pedroprado.transaction.api/src/core/useCases/account"
	"pedroprado.transaction.api/src/core/useCases/balanceProvision"
	"pedroprado.transaction.api/src/core/useCases/transaction"
	"pedroprado.transaction.api/src/core/useCases/transactionStatus"
	"pedroprado.transaction.api/src/infra"
	rest "pedroprado.transaction.api/src/presentation"
	"pedroprado.transaction.api/src/presentation/accountApi"
	"pedroprado.transaction.api/src/presentation/balanceProvisionApi"
	"pedroprado.transaction.api/src/presentation/transactionApi"
	"pedroprado.transaction.api/src/presentation/transactionStatusApi"
	"pedroprado.transaction.api/src/repository"
	"pedroprado.transaction.api/src/repository/accountRepo"
	"pedroprado.transaction.api/src/repository/balanceProvisionRepo"
	"pedroprado.transaction.api/src/repository/errorHandler"
	"pedroprado.transaction.api/src/repository/transactionRepo"
	"pedroprado.transaction.api/src/repository/transactionStatusRepo"
)

var (
	sslMode = false
)

const (
	apiPort      = "8098"
	relativePath = "/snapfi"
)

func main() {
	db := getDb()
	errHandler := errorHandler.NewPostgresErrorHandler()

	accountRepository := accountRepo.NewAccountRepository(db, errHandler)
	transactionRepository := transactionRepo.NewTransactionRepository(db, errHandler)
	transactionStatusRepository := transactionStatusRepo.NewTransactionStatusRepository(db, errHandler)
	balanceProvisionRepository := balanceProvisionRepo.NewBalanceProvisionRepository(db, errHandler)

	accountService := account.NewAccountService(accountRepository)
	transactionService := transaction.NewTransactionService(transactionRepository,
		transactionStatusRepository, accountRepository, balanceProvisionRepository, db)
	balanceProvisionService := balanceProvision.NewBalanceProvisionService(balanceProvisionRepository)
	transactionStatusService := transactionStatus.NewTransactionStatusService(transactionStatusRepository)

	ginServer := rest.NewServerHttpGin()
	ginRouterGroup := ginServer.GetGinRouterGroup(relativePath)
	accountApi.RegisterAccountApi(ginRouterGroup, accountService)
	transactionApi.RegisterTransactionApi(ginRouterGroup, transactionService)
	balanceProvisionApi.RegisterBalanceProvisionApi(ginRouterGroup, balanceProvisionService)
	transactionStatusApi.RegisterTransactionStatusApi(ginRouterGroup, transactionStatusService)
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
