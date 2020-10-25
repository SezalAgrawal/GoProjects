package db

import (
	"database/sql"
	"go.caringcompany.co/simpliclaim/utils/stringops"
)

const (
	// DefaultMainDB is the name of default database
	DefaultMainDB = "mydatabasename"
)

// DB is database type that encapsulates *sql.DB and implements PG interface.
type database struct {
	session *sql.DB
	url     string // #Warning do not print to console or log this
	mainDB  string
}

// Provider holds all database methods
type Provider interface {
	Close() error
	Ok() (bool, error)
	Info() string
}

// NewProvider panics if database client fail to connect with given database server.
func NewProvider(dbURL, mainDB string) Provider {
	var err error
	if stringops.IsBlank(mainDB) {
		mainDB = DefaultMainDB
	}
	db := &database{
		url:    dbURL,
		mainDB: mainDB,
	}
	if db.session, err = sql.Open("pgx", "postgresql://sezalagrawal:@localhost:5432/mydatabasename?sslmode=disable"); err != nil {
		panic(err.Error())
	}
	return db
}
