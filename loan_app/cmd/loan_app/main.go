package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/goProjects/loan_app/app"
	"github.com/goProjects/loan_app/lib/db"
	"github.com/goProjects/loan_app/lib/logger"
	"github.com/goProjects/loan_app/lib/server"
	"github.com/goProjects/loan_app/lib/utils"
	"github.com/julienschmidt/httprouter"
)

func main() {
	webApp := app.NewApp()

	setupLogger(webApp)
	setupDB(webApp)

	defer logger.Sync()
	defer db.Close()

	server := setupServer(webApp)
	if err := server.Run(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}

func setupLogger(webApp *app.App) {
	var logLevel int
	env := webApp.Config().ENV

	if env == utils.ProductionEnv {
		logLevel = logger.INFO
	} else {
		logLevel = logger.DEBUG // TODO: once the service stabilizes, make log level as INFO
	}
	logger.Init(logLevel, env)
}

func setupDB(webApp *app.App) {
	db.Init(webApp.Config().DATABASE_URL, webApp.Config().DB_MAX_IDLE_CONNECTIONS, webApp.Config().DB_MAX_OPEN_CONNECTIONS)
}

func setupServer(webApp *app.App) *server.GracefulShutdownServer {
	router := httprouter.New()
	app.InitRoutes(router)

	return &server.GracefulShutdownServer{
		Server: &http.Server{
			Addr: fmt.Sprintf("0.0.0.0:%d", webApp.Config().PORT),
			// set timeouts to avoid slowloris attacks
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
			IdleTimeout:  60 * time.Second,
			Handler:      router,
		},
	}
}
