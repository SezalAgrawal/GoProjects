package app

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/goProjects/loan_app/lib/utils"
)

type App struct {
	config config
}

type config struct {
	DATABASE_URL            string            `env:"DATABASE_URL,notEmpty"`
	DB_MAX_IDLE_CONNECTIONS int               `env:"DB_MAX_IDLE_CONNECTIONS,notEmpty"`
	DB_MAX_OPEN_CONNECTIONS int               `env:"DB_MAX_OPEN_CONNECTIONS,notEmpty"`
	ENV                     utils.Environment `env:"ENV,notEmpty"`
	PORT                    int               `env:"PORT,notEmpty"`
	SERVICE_AUTH_CONFIG     string            `env:"SERVICE_AUTH_CONFIG"`
}

func NewApp() *App {
	opts := env.Options{RequiredIfNoDef: true}

	cfg := &config{}
	if err := env.Parse(cfg, opts); err != nil {
		panic(fmt.Errorf("unable to parse env config %w", err))
	}
	return &App{config: *cfg}
}

func (a *App) Config() config {
	return a.config
}
