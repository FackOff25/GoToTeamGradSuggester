package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/FackOff25/GoToTeamGradSuggester/internal/usecase"
	"github.com/FackOff25/GoToTeamGradSuggester/pkg/config"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	Usecase usecase.UsecaseInterface
	Cfg     *config.Config
}

func (pc *Controller) Get(c echo.Context) error {
	defer c.Request().Body.Close()

	user, _ := pc.Usecase.GetUser()

	return c.JSON(http.StatusOK, user)
}

func (pc *Controller) CreatePlacesListHandler(c echo.Context) error {
	defer c.Request().Body.Close()

	log.Printf("Got request")

	if !c.QueryParams().Has("location") {
		return echo.ErrBadRequest
	}

	location := c.QueryParam("location")

	radius := 1000
	var err error
	if c.QueryParams().Has("radius") {
		radius, err = strconv.Atoi(c.QueryParam("radius"))
		if err != nil {
			log.Printf("Bad radius: %s", c.QueryParam("radius"))
			return echo.ErrBadRequest
		}
	}
	/*
		types := []string{
			"aquarium",
			"art_gallery",
			"cafe",
			"church",
			"museum",
			"park",
		}
	*/
	places, _ := pc.Usecase.GetNearbyPlaces(pc.Cfg, location, radius, "park")

	return c.JSON(http.StatusOK, places)
}
