package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/FackOff25/GoToTeamGradGoLibs/logger"
	"github.com/FackOff25/GoToTeamGradGoLibs/postgres"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/controller"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/controller/handler"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/repository"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/repository/queries"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/usecase"
	"github.com/FackOff25/GoToTeamGradSuggester/pkg/config"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func Run(configFilePath string) {
	cfg, err := config.GetConfig(configFilePath)

	configOutput := cfg.LogOutput
	if err := os.MkdirAll(filepath.Dir(configOutput), 0770); err != nil {
		panic(err)
	}
	logFile, err := os.OpenFile(configOutput, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("Error opening log file: %s", err)
	}
	defer func() {
		if err := logFile.Close(); err != nil {
			panic(err)
		}
	}()

	logOutput := io.MultiWriter(os.Stdout, logFile)

	logger.InitEx(logger.Options{
		Name:      cfg.LogAppName,
		LogLevel:  log.Level(cfg.LogLevel),
		LogFormat: cfg.LogFormat,
		Out:       logOutput,
	})

	log.Infof("%s", cfg.LogFormat)

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
	uc := usecase.New(*repo, ctx, config)

	c := controller.Controller{Usecase: uc, Cfg: config}

	e.GET("/api/v1/suggest/nearby", c.CreatePlacesListHandler)
	e.GET("/api/v1/suggest/user/", c.GetUser)
	e.POST("/api/v1/suggest/user/new", c.AddUser)

	e.POST("/api/v1/suggest/reaction", c.CreateNewReactionHandler)
	e.POST("/api/v1/suggest/reaction/initial", c.OnboardingTypesHandler)

	e.GET("/api/v1/suggest/dummy", handler.CreateNotImplementedResponse)
	e.GET("/api/v1/suggest/ping", c.Ping)

	e.POST("/api/v1/suggest/route", c.GetRoute)

	e.GET("/api/v1/suggest/categories/list", c.GetCategoriesHandler)

	e.POST("/api/v1/suggest/route/sortPlaces", c.SortPlaces)

	return nil
}
