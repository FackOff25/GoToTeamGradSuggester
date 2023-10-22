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
	"github.com/FackOff25/GoToTeamGradSuggester/pkg/config/configReader"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func Run(configFilePath string) {
	configReader := configReader.NewConfigReader(configFilePath)
	
	config, err := configReader.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	serverAddress := config.ServerAddress + ":" + config.ServerPort
	fmt.Println(config)

	if err := configureServer(e); err != nil {
		log.Fatalf("error while configuring server: %s", err)
	}

	if err := e.Start(serverAddress); err != nil {
		log.Fatalf("error while starting server: %s", err)
	}
}

func configureServer(e *echo.Echo) error {

	ctx := context.Background()
	repo := repository.New(&queries.Queries{Ctx: ctx, Pool: pgxpool.Pool{}}, ctx)
	uc := usecase.New(*repo, ctx)

	controller := controller.Controller{Usecase: uc}

	e.GET("/api/v1/get", controller.Get)

	e.GET("/api/v1/dummy", handler.CreateNotImplementedResponse)

	return nil
}