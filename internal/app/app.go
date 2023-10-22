package app

import (
	"log"

	config "github.com/FackOff25/GoToTeamGradSuggester/pkg"
	"github.com/labstack/echo/v4"
)

func Run(configFilePath string) {

	cfg, err := config.GetConfig(configFilePath)

	if err != nil {
		log.Fatalf("error while reading config: %s", err)
	}

	serverAddress := cfg.ServerAddress + ":" + cfg.ServerPort

	e := echo.New()

	if err := configureServer(e); err != nil {
		log.Fatalf("error while configuring server: %s", err)
	}

	if err := e.Start(serverAddress); err != nil {
		log.Fatalf("error while starting server: %s", err)
	}
}

func configureServer(e *echo.Echo) error {

	//e.GET("/api/v1/suggset", placesController.CreatePlacesListHandler)

	return nil
}
