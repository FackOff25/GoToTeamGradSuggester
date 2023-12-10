package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/FackOff25/GoToTeamGradGoLibs/categories"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
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
	if requestBody.PlaceId == "" ||
		(requestBody.Reaction != domain.ReactionVisited &&
			requestBody.Reaction != domain.ReactionLike &&
			requestBody.Reaction != domain.ReactionRefuse &&
			requestBody.Reaction != domain.ReactionUnvisited &&
			requestBody.Reaction != domain.ReactionUnlike &&
			requestBody.Reaction != domain.ReactionUnrefuse) {
		return echo.ErrBadRequest
	}

	err = pc.Usecase.ApplyUserReactionToPlace(uuid[0], requestBody.PlaceId, requestBody.Reaction)
	if err != nil {
		if strings.Contains(err.Error(), pgx.ErrNoRows.Error()) {
			return echo.ErrNotFound
		}
		log.Errorf("applying error: %s", err.Error())
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, "OK")
}

func (pc *Controller) OnboardingTypesHandler(c echo.Context) error {
	defer c.Request().Body.Close()

	userId, ok := c.Request().Header[uuidHeader]
	if !ok {
		return echo.ErrUnauthorized
	}

	var requestBody domain.OnboardingReactions
	err := json.NewDecoder(c.Request().Body).Decode(&requestBody)
	if err != nil {
		return echo.ErrBadRequest
	}

	types := []string{}
	cats := categories.GetReversedCategoryMap()
	for _, v := range requestBody.Primary {
		placeType, ok := cats[v]
		if !ok {
			log.Errorf("Bad category: %s", v)
			return echo.ErrBadRequest
		}
		types = append(types, placeType)
	}

	err = pc.Usecase.ApplyOnboardingReactions(userId[0], types)
	if err != nil {
		log.Errorf("applying error: %s", err.Error())
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, "OK")
}
