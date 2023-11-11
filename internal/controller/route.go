package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
	"github.com/jackc/pgx/v5"
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

	// err = pc.Usecase.ApplyUserReactionToPlace(uuid[0], requestBody.PlaceId, requestBody.Reaction)
	if err != nil {
		if strings.Contains(err.Error(), pgx.ErrNoRows.Error()) {
			return echo.ErrNotFound
		}
		log.Errorf("applying error: %s", err.Error())
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, "OK")
}
