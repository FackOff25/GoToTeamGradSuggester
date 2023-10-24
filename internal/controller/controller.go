package controller

import (
	"net/http"

	"github.com/FackOff25/GoToTeamGradSuggester/internal/usecase"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	Usecase usecase.UsecaseInterface
}

func (pc *Controller) Get(c echo.Context) error {
	defer c.Request().Body.Close()

	user, _ := pc.Usecase.GetUser()

	return c.JSON(http.StatusOK, user)
}