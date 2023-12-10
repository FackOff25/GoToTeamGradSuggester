package controller

import (
	"bytes"
	"encoding/json"
	"html"
	"net/http"
	"strconv"
	"strings"

	"github.com/FackOff25/GoToTeamGradGoLibs/categories"
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
		if radius > 6000 {
			radius = 6000
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

	types := []string{}
	if c.QueryParams().Has("types") {
		typesStr := c.QueryParam("types")
		typesSlice := strings.Split(typesStr, ",")
		cats := categories.GetReversedCategoryMap()
		for _, v := range typesSlice {
			placeType, ok := cats[v]
			if !ok {
				log.Errorf("Bad category: %s", v)
				return echo.ErrBadRequest
			}
			types = append(types, placeType)
		}
	}

	reactions := []string{}
	if c.QueryParams().Has("reactions") {
		reactionsStr := c.QueryParam("reactions")
		reactionsSlice := strings.Split(reactionsStr, ",")
		reactionsMap := domain.GetFilterReactionsMap()
		for _, v := range reactionsSlice {
			_, ok := reactionsMap[v]
			if ok {
				reactions = append(reactions, v)
			}
		}
	}

	places, _ := pc.Usecase.GetMergedNearbyPlaces(pc.Cfg, user, location, radius, limit, offset, types, reactions)

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

func (pc *Controller) GetCategoriesHandler(c echo.Context) error {
	defer c.Request().Body.Close()

	var result []string
	categoriesMap := categories.GetCategoriesMap()
	for _, v := range pc.Cfg.Categories {
		cat, ok := categoriesMap[v]
		if ok {
			result = append(result, cat)
		}

	}

	resp := domain.GetCategoriesResponse{
		Categories: result,
	}

	resBodyBytes := new(bytes.Buffer)
	encoder := json.NewEncoder(resBodyBytes)
	encoder.SetEscapeHTML(false)
	encoder.Encode(resp)

	return c.JSONBlob(http.StatusOK, resBodyBytes.Bytes())
}
