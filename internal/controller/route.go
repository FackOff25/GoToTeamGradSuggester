package controller

import (
	"encoding/json"
	"net/http"

	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func (pc *Controller) GetRoute(c echo.Context) error {
	defer c.Request().Body.Close()

	uuid, ok := c.Request().Header[uuidHeader]
	if !ok {
		return echo.ErrUnauthorized
	}

	u, _ := pc.Usecase.GetUser(uuid[0])
	if u == nil {
		return echo.ErrUnauthorized
	}

	var requestBody domain.RouteReq
	err := json.NewDecoder(c.Request().Body).Decode(&requestBody)
	if err != nil {
		return echo.ErrBadRequest
	}

	resBody, err := pc.Usecase.GetRoute(&requestBody)
	if err != nil {
		log.Errorf("getting route error: %s", err.Error())
		return echo.ErrInternalServerError
	}

	resBodyBytes, err := json.Marshal(resBody)
	if err != nil {
		log.Errorf("marshaling error: %s", err.Error())
		return echo.ErrInternalServerError
	}

	return c.JSONBlob(http.StatusOK, resBodyBytes)
}
