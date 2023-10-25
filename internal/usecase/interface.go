package usecase

import (
	"github.com/FackOff25/GoToTeamGradGoLibs/googleApi"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
	"github.com/FackOff25/GoToTeamGradSuggester/pkg/config"
)

type UsecaseInterface interface {
	GetUser() (*domain.User, error)
	GetNearbyPlaces(cfg *config.Config, location string, radius int, placeType string) ([]googleApi.Place, error)
}
