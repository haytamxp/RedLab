package app

import (
	"fmt"

	"github.com/haytamxp/redlab/backend/internal/config"
	"github.com/haytamxp/redlab/backend/internal/database"
	"github.com/haytamxp/redlab/backend/internal/handlers"
	"github.com/haytamxp/redlab/backend/internal/logger"
	"github.com/haytamxp/redlab/backend/internal/repository"
	"github.com/haytamxp/redlab/backend/internal/router"
	"github.com/haytamxp/redlab/backend/internal/services"
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

	db, err := database.Connect(cfg)
	if err != nil {
		return nil, err
	}

	userRepository := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepository)

	authService := services.NewAuthService(
		userService,
		cfg.JWT.Secret,
		cfg.JWT.Expiration,
	)

	authHandler := handlers.NewAuthHandler(authService)

	r := router.New()
	r.RegisterRoutes(authHandler)

	return &App{
		Config: cfg,
		Router: r,
	}, nil
}

func (a *App) Start() error {

	logger.Log.Info("Starting RedLab Backend")

	address := fmt.Sprintf(
		"%s:%s",
		a.Config.Server.Host,
		a.Config.Server.Port,
	)

	return Run(address, a.Router.Engine)
}