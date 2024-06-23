package test_helper

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/goProjects/loan_app/lib/db"
)

var (
	_, b, _, _            = runtime.Caller(0)
	relPath               = filepath.Join(filepath.Dir(b), "../../migrations/db/structure.sql")
	defaultSqlFilePath, _ = filepath.Abs(relPath)
)

func SetupDatabase() {
	db.Init(os.Getenv("DATABASE_URL"), 5, 10)

	ClearDataFromPostgres()
}

func TeardownDatabase() {
	db.Close()
}

func loadSchemaFile(db *sql.DB, schemaFilePath string) {
	content, err := os.ReadFile(schemaFilePath)
	if err != nil {
		panic(fmt.Errorf("error loading schema file %s: %w", schemaFilePath, err))
	}

	_, err = db.Exec(string(content))
	if err != nil {
		panic(fmt.Errorf("error executing schema file %s: %w", schemaFilePath, err))
	}
}

func ClearDataFromPostgres() {
	gormDB := db.Get()

	sqlDB, err := gormDB.DB()
	if err != nil {
		panic(fmt.Errorf("error in fetching sql db: %w", err))
	}

	loadSchemaFile(sqlDB, defaultSqlFilePath)
}
