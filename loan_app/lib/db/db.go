package db

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Get() *gorm.DB {
	return db
}

func Init(url string, maxIdleConnections, maxOpenConnections int) {
	var err error
	dbConn, err := sql.Open("pgx", url)
	if err != nil {
		panic(fmt.Errorf("unable to configure new relic postgres %w", err))
	}

	db, err = gorm.Open(postgres.New(postgres.Config{Conn: dbConn}))
	if err != nil {
		panic(fmt.Errorf("unable to connect to db %w", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("unable to configure db: sql db get error: %w", err))
	}

	if err = sqlDB.Ping(); err != nil {
		panic(fmt.Errorf("unable to connect to db: sql db get error: %w", err))
	}

	sqlDB.SetMaxIdleConns(maxIdleConnections)
	sqlDB.SetMaxOpenConns(maxOpenConnections)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func Close() {
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("unable to close db: sql db get error: %w", err))
	}
	if err = sqlDB.Close(); err != nil {
		panic(fmt.Errorf("unable to close db: closing sql db error: %w", err))
	}
}
