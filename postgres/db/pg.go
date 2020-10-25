package db

import (
	"fmt"
)

// Ok check if database is running or not
func (db *database) Ok() (bool, error) {
	if err := db.session.Ping(); err != nil {
		return false, err
	}
	return true, nil
}

// Info provides information about the database such as
// the name of the database currently in use.
func (db *database) Info() string {
	return fmt.Sprintf("Database in use: %s url: %s", db.mainDB, db.url)
}

// Close closes the databases connections
func (db *database) Close() error {
	return db.session.Close()
}