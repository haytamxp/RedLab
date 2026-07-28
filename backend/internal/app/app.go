package app

import (
	"fmt"

	"github.com/haytamxp/redlab/backend/internal/config"
	"github.com/haytamxp/redlab/backend/internal/logger"
	"github.com/haytamxp/redlab/backend/internal/router"
)

type App struct {
	Config *config.Config
	Router *router.Router
}

func New() (*App, error) {

	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	log, err := logger.New(cfg.Logger.Development)
	if err != nil {
		return nil, err
	}

	logger.Init(log)

	r := router.New()

	app := &App{
		Config: cfg,
		Router: r,
	}

	return app, nil
}

func (a *App) Start() error {

	logger.Log.Info("Starting RedLab Backend")

	address := fmt.Sprintf("%s:%s", a.Config.Server.Host, a.Config.Server.Port)

	return Run(address, a.Router.Engine)
}