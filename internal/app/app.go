package app

import (
	"context"
	"log"

	"github.com/FackOff25/GoToTeamGradSuggester/internal/controller"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/controller/handler"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/repository"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/repository/queries"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/FackOff25/GoToTeamGradSuggester/internal/usecase"
	"github.com/FackOff25/GoToTeamGradSuggester/pkg/config"
	"github.com/labstack/echo/v4"
)

func Run(configFilePath string) {
	cfg, err := config.GetConfig(configFilePath)

	if err != nil {
		log.Fatalf("error while reading config: %s", err)
	}

	serverAddress := cfg.ServerAddress + ":" + cfg.ServerPort

	e := echo.New()

	if err := configureServer(e, cfg); err != nil {
		log.Fatalf("error while configuring server: %s", err)
	}

	if err := e.Start(serverAddress); err != nil {
		log.Fatalf("error while starting server: %s", err)
	}
}

func configureServer(e *echo.Echo, config *config.Config) error {
	ctx := context.Background()
	repo := repository.New(&queries.Queries{Ctx: ctx, Pool: pgxpool.Pool{}}, ctx)
	uc := usecase.New(*repo, ctx)

	controller := controller.Controller{Usecase: uc, Cfg: config}

	e.GET("/api/v1/suggset/get", controller.Get)

	e.GET("/api/v1/suggset/nearby", controller.CreatePlacesListHandler)

	e.GET("/api/v1/suggest/dummy", handler.CreateNotImplementedResponse)

	return nil
}
