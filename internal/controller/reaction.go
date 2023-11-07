package controller

import (
	"encoding/json"
	"net/http"

	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
	"github.com/labstack/echo/v4"
)

func (pc *Controller) CreateNewReactionHandler(c echo.Context) error {
	defer c.Request().Body.Close()

	uuid, ok := c.Request().Header[uuidHeader]
	if !ok {
		return echo.ErrUnauthorized
	}

	var requestBody domain.NewReactionRequest
	err := json.NewDecoder(c.Request().Body).Decode(&requestBody)
	if err != nil {
		return echo.ErrBadRequest
	}

	err = pc.Usecase.ApplyUserReactionToPlace(uuid[0], requestBody.PlaceId, requestBody.Reaction)
	if err != nil {
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, "OK")
}
