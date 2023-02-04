package main

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"pedroprado.transaction.api/src/core/useCases/account"
	"pedroprado.transaction.api/src/core/useCases/balanceProvision"
	"pedroprado.transaction.api/src/core/useCases/transaction"
	"pedroprado.transaction.api/src/core/useCases/transaction/compensate"
	"pedroprado.transaction.api/src/core/useCases/transaction/complete"
	"pedroprado.transaction.api/src/core/useCases/transaction/create"
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

	_ "pedroprado.transaction.api/src/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	sslMode = false
)

const (
	apiPort      = "8098"
	relativePath = "/snapfi"
)

func main() {
	setLoggerFormatter()
	db := getDb()
	errHandler := errorHandler.NewPostgresErrorHandler()

	accountRepository := accountRepo.NewAccountRepository(db, errHandler)
	transactionRepository := transactionRepo.NewTransactionRepository(db, errHandler)
	transactionStatusRepository := transactionStatusRepo.NewTransactionStatusRepository(db, errHandler)
	balanceProvisionRepository := balanceProvisionRepo.NewBalanceProvisionRepository(db, errHandler)

	accountService := account.NewAccountService(accountRepository)
	createTransactionService := create.NewCreateTransactionService(transactionRepository,
		transactionStatusRepository, accountRepository, balanceProvisionRepository, db)
	completeTransactionService := complete.NewCompleteTransactionService(transactionRepository,
		transactionStatusRepository, accountRepository, balanceProvisionRepository, db)
	compensateTransactionService := compensate.NewCompensateTransactionService(transactionRepository,
		transactionStatusRepository, accountRepository, balanceProvisionRepository, db)
	transactionService := transaction.NewTransactionService(transactionRepository,
		createTransactionService, completeTransactionService, compensateTransactionService)
	balanceProvisionService := balanceProvision.NewBalanceProvisionService(balanceProvisionRepository)
	transactionStatusService := transactionStatus.NewTransactionStatusService(transactionStatusRepository)

	ginServer := rest.NewServerHttpGin()
	ginRouterGroup := ginServer.GetGinRouterGroup(relativePath)
	ginRouterGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	accountApi.RegisterAccountApi(ginRouterGroup, accountService)
	transactionApi.RegisterTransactionApi(ginRouterGroup, transactionService)
	balanceProvisionApi.RegisterBalanceProvisionApi(ginRouterGroup, balanceProvisionService)
	transactionStatusApi.RegisterTransactionStatusApi(ginRouterGroup, transactionStatusService)
	rest.RegisterInfraApi(ginRouterGroup, false)
	if err := ginServer.StartServer(apiPort); err != nil {
		logrus.Fatal(errors.WithStack(err))
	}
}

func setLoggerFormatter() {
	formatter := &logrus.JSONFormatter{}
	formatter.TimestampFormat = "2006-01-02T15:04:05.999999999Z"
	logrus.SetFormatter(formatter)
	logrus.SetOutput(&infra.LogWriter{})
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
