package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/FackOff25/GoToTeamGradGoLibs/googleApi"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/repository/queries"
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

func (pc *Controller) formNearbyPlace(result googleApi.Place) (domain.NearbyPlace, error) {
	location := domain.ApiLocation{
		Lat: result.Geometry.Location.Lat,
		Lng: result.Geometry.Location.Lng,
	}

	var cover string
	if len(result.Photos) > 0 {
		reference := result.Photos[0].Reference
		cover = pc.Cfg.PlacesApiHost + "place/photo?maxwidth=" + strconv.FormatInt(result.Photos[0].Width, 10) + "&photo_reference=" + reference
		result.Photos = result.Photos[1:]
	}

	var photos []string
	for _, photo := range result.Photos {
		cover = pc.Cfg.PlacesApiHost + "place/photo?maxwidth=" + strconv.FormatInt(photo.Width, 10) + "&photo_reference=" + photo.Reference
	}

	return domain.NearbyPlace{
		PlaceId:     result.PlaceId,
		Name:        result.Name,
		Location:    location,
		Cover:       cover,
		Photos:      photos,
		Rating:      float32(result.Rating),
		RatingCount: int(result.RatingCount),
	}, nil
}

func (pc *Controller) CreatePlacesListHandler(c echo.Context) error {
	defer c.Request().Body.Close()

	if !c.QueryParams().Has("location") {
		return echo.ErrBadRequest
	}

	location := c.QueryParam("location")

	radius := 10000
	var err error
	if c.QueryParams().Has("radius") {
		radius, err = strconv.Atoi(c.QueryParam("radius"))
		if err != nil {
			log.Errorf("Bad radius: %s", c.QueryParam("radius"))
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

	sort.Slice(places, func(i, j int) bool {
		return queries.ComparePlaces(places[i], places[j])
	})

	var result []domain.NearbyPlace

	for _, v := range places {
		place, _ := pc.formNearbyPlace(v)
		result = append(result, place)
	}

	resBodyBytes := new(bytes.Buffer)
	encoder := json.NewEncoder(resBodyBytes)
	encoder.SetEscapeHTML(false)
	encoder.Encode(result)

	return c.JSONBlob(http.StatusOK, resBodyBytes.Bytes())
}
