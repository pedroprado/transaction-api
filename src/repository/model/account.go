package model

import "time"

type Account struct {
	AccountID string `gorm:"primaryKey"`
	Bank      string
	Number    string
	Agency    string
	Balance   float64
	CreatedAt time.Time
	UpdatedAt time.Time
}
