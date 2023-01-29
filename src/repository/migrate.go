package repository

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"pedroprado.transaction.api/src/repository/model"
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
	)
	if autoMigrateResult != nil {
		return errors.WithStack(autoMigrateResult)
	}

	migrateResult := db.Model(&model.Transaction{})
	if migrateResult.Error != nil {
		return errors.WithStack(migrateResult.Error)
	}

	return nil
}

func temporaryMigrations(db *gorm.DB) error {

	return nil
}
