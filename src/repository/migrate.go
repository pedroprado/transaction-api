package repository

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"pedroprado.transaction.api/src/repository/model"
	"time"
)

var (
	intermediaryAccount = model.Account{
		AccountID: "12345",
		Bank:      "10",
		Number:    "10",
		Agency:    "9",
		Balance:   9999999,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
)

func Migrate(db *gorm.DB) error {
	err := prevMigrations(db)
	if err != nil {
		return errors.WithStack(err)
	}

	err = fixedMigrations(db)
	if err != nil {
		return errors.WithStack(err)
	}

	err = temporaryMigrations(db)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func prevMigrations(db *gorm.DB) error {
	return nil
}

func fixedMigrations(db *gorm.DB) error {
	autoMigrateResult := db.AutoMigrate(
		&model.Transaction{},
		&model.TransactionStatus{},
		&model.Account{},
		&model.BalanceProvision{},
	)
	if autoMigrateResult != nil {
		return errors.WithStack(autoMigrateResult)
	}

	result := db.Create(&intermediaryAccount)
	if result.Error != nil {
		logrus.Warnf(result.Error.Error())
	}

	return nil
}

func temporaryMigrations(db *gorm.DB) error {

	return nil
}
