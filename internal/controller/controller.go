package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/FackOff25/GoToTeamGradSuggester/internal/usecase"
	"github.com/FackOff25/GoToTeamGradSuggester/pkg/config"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
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

	if !c.QueryParams().Has("location") {
		return echo.ErrBadRequest
	}

	location := c.QueryParam("location")

	radius := 5000
	var err error
	if c.QueryParams().Has("radius") {
		radius, err = strconv.Atoi(c.QueryParam("radius"))
		if err != nil {
			log.Errorf("Bad radius: %s", c.QueryParam("radius"))
			return echo.ErrBadRequest
		}
	}

	limit := 20
	if c.QueryParams().Has("limit") {
		limit, err = strconv.Atoi(c.QueryParam("limit"))
		if err != nil {
			log.Errorf("Bad limit: %s", c.QueryParam("limit"))
			return echo.ErrBadRequest
		}
	}

	offset := 0
	if c.QueryParams().Has("offset") {
		offset, err = strconv.Atoi(c.QueryParam("offset"))
		if err != nil {
			log.Errorf("Bad offset: %s", c.QueryParam("offset"))
			return echo.ErrBadRequest
		}
	}

	places, _ := pc.Usecase.GetMergedNearbyPlaces(pc.Cfg, location, radius, limit, offset)

	places = pc.Usecase.SortPlaces(places)[offset:]

	places = pc.Usecase.UniqPlaces(places)

	if len(places) > limit {
		places = places[:limit]
	}

	resBodyBytes := new(bytes.Buffer)
	encoder := json.NewEncoder(resBodyBytes)
	encoder.SetEscapeHTML(false)
	encoder.Encode(places)

	return c.JSONBlob(http.StatusOK, resBodyBytes.Bytes())
}
