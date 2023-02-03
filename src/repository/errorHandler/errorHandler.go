package errorHandler

import (
	"gorm.io/gorm"
)

type ErrorHandler interface {
	IsRecordDuplicated(db *gorm.DB) (result bool)
	ViolatesFKConstraint(db *gorm.DB) (result bool)
	GetDetail(db *gorm.DB) (detail string)
	IsRecordNotFound(db *gorm.DB) bool
}
