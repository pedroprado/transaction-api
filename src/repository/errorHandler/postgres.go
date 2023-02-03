package errorHandler

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type postgresErrorHandler struct {
}

func NewPostgresErrorHandler() ErrorHandler {
	return &postgresErrorHandler{}
}

func (ref *postgresErrorHandler) IsRecordDuplicated(db *gorm.DB) bool {
	if db.Error != nil && db.Error.(*pgconn.PgError).Code == "23505" {
		return true
	}
	return false
}

func (ref *postgresErrorHandler) ViolatesFKConstraint(db *gorm.DB) (result bool) {
	if db.Error != nil && db.Error.(*pgconn.PgError).Code == "23503" {
		return true
	}
	return false
}

func (ref *postgresErrorHandler) GetDetail(db *gorm.DB) (detail string) {
	if db.Error == nil {
		return ""
	}
	return db.Error.(*pgconn.PgError).Detail
}

func (ref *postgresErrorHandler) IsRecordNotFound(db *gorm.DB) bool {
	return db.Error != nil && errors.Is(db.Error, gorm.ErrRecordNotFound)
}
