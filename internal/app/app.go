package app

import (
	"context"
	"fmt"
	"log"

	"github.com/FackOff25/GoToTeamGradSuggester/internal/controller"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/controller/handler"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/repository"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/repository/queries"	
	"github.com/FackOff25/GoToTeamGradSuggester/internal/usecase"
	"github.com/FackOff25/GoToTeamGradSuggester/pkg/config"
	"github.com/FackOff25/GoToTeamGradSuggester/pkg/postgres"
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
	pg, err := postgres.New(fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s", config.DBuser, config.DBpassword, config.DBurl, config.DBport, config.DBname), ctx)
	if err != nil {
		return err
	}
	p, err := pg.Connect()
	if err != nil {
		return err
	}
	repo := repository.New(&queries.Queries{Ctx: ctx, Pool: *p}, ctx)
	uc := usecase.New(*repo, ctx)

	c := controller.Controller{Uc: uc, Cfg: config}

	e.GET("/api/v1/suggest/nearby", c.CreatePlacesListHandler)
	e.POST("/api/v1/user/", c.GetUser)

	e.GET("/api/v1/suggest/dummy", handler.CreateNotImplementedResponse)


	return nil
}
