package controller

import (
	"bytes"
	"encoding/json"
	"html"
	"net/http"
	"strconv"

	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/usecase"
	"github.com/FackOff25/GoToTeamGradSuggester/pkg/config"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

const uuidHeader = "X-Uuid"

type Controller struct {
	Usecase usecase.UsecaseInterface
	Cfg     *config.Config
}

func (pc *Controller) Ping(c echo.Context) error {
	defer c.Request().Body.Close()
	return c.JSON(http.StatusOK, nil)
}

func (pc *Controller) GetUser(c echo.Context) error {
	defer c.Request().Body.Close()

	id := html.EscapeString(c.Request().Header.Get("X-UUID"))

	u, err := pc.Usecase.GetUser(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return echo.ErrNotFound
		}
		log.Errorf("Repo error: %s; id: %s", err.Error(), id)
		return echo.ErrInternalServerError
	}

	u.PlaceTypePreferences = nil // на клиенте не нужны предпочтения
	b, err := json.Marshal(u)
	if err != nil {
		log.Errorf("Marshal error: %s; id: %s", err.Error(), id)
		return echo.ErrInternalServerError
	}

	return c.JSONBlob(http.StatusOK, b)
}

func (pc *Controller) AddUser(c echo.Context) error {
	defer c.Request().Body.Close()

	id := html.EscapeString(c.Request().Header.Get("X-UUID"))

	if id == "" {
		return echo.ErrBadRequest
	}

	err := pc.Usecase.AddUser(id)
	if err != nil {
		if err.Error() == domain.ErrorUserAlreadyExists {
			return echo.ErrConflict
		}
		log.Errorf("Repo error: %s; id: %s", err.Error(), id)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, nil)
}

func (pc *Controller) CreatePlacesListHandler(c echo.Context) error {
	defer c.Request().Body.Close()

	uuid, ok := c.Request().Header[uuidHeader]
	if !ok {
		return echo.ErrUnauthorized
	}

	user, err := pc.Usecase.GetUser(uuid[0])
	if err != nil {
		return echo.ErrUnauthorized
	}

	if !c.QueryParams().Has("location") {
		return echo.ErrBadRequest
	}

	location := c.QueryParam("location")

	radius := 5000
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

	places, err := pc.Usecase.GetMergedNearbyPlaces(pc.Cfg, user, location, radius, limit, offset)
	if err != nil {
		log.Error(err)
		return echo.ErrInternalServerError
	}

	places = pc.Usecase.UniqPlaces(places)

	places = pc.Usecase.SortPlaces(places)

	if len(places) > offset {
		places = places[offset:]
		if len(places) > limit {
			places = places[:limit]
		}
	} else {
		places = []domain.SuggestPlace{}
	}

	resBodyBytes := new(bytes.Buffer)
	encoder := json.NewEncoder(resBodyBytes)
	encoder.SetEscapeHTML(false)
	encoder.Encode(places)

	return c.JSONBlob(http.StatusOK, resBodyBytes.Bytes())
}
