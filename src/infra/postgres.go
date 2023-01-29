package infra

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

var (
	host     = os.Getenv("DATABASE_HOST")
	dbPort   = os.Getenv("DATABASE_PORT")
	user     = os.Getenv("DATABASE_USERNAME")
	dbName   = os.Getenv("DATABASE_NAME")
	password = os.Getenv("DATABASE_PASSWORD")
)

func OpenConnection(sslMode bool) (*gorm.DB, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", host, dbPort, user, dbName, password)

	if !sslMode {
		connString = connString + " sslmode=disable"
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: connString,
	}), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, errors.WithMessagef(err,
			"Failed to connect database, host: %s - port: %s - user: %s - dbname: %s",
			host, dbPort, user, dbName)
	}

	sqlDB, _ := db.DB()
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(50)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
