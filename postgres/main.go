package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/goProjects/test/db"
	"go.caringcompany.co/simpliclaim/log"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var (
	database db.Provider
)

func mustInitDB() {
	url := os.Getenv("DATABASE_URL")
	mainDB := os.Getenv("DATABASE_NAME")
	url = "postgresql://sezalagrawal:@localhost:5432"
	mainDB = "mydatabasename"
	database = db.NewProvider(url, mainDB)
	if ok, err := database.Ok(); !ok {
		panic(fmt.Sprintf("can not establish connection to database: %s reason: %s", mainDB, err))
	}
	log.Info(database.Info())
}

func closeDBConnection() {
	if err := database.Close(); err != nil {
		log.Errorf("error while closing database connection: %s", err)
	}
}

func main() {
	mustInitDB()
	defer closeDBConnection()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-signals:
			log.Info("signal interrupt. exiting...")
			return
		}
	}
}
