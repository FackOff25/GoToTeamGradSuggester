package repository

import (
	"github.com/FackOff25/GoToTeamGradGoLibs/googleApi"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
)

type QueriesInterface interface {
	GetUser(id string) (*domain.User, error)
	AddUser(id string, categories []string) error
	ApplyUserReactionToPlace(uuid string, placeId string, reaction string, types []string) error
	SavePlaces([]googleApi.Place) error
	AddPlace(gID string, types []string) error
	GetPlaceById(gID string) (*domain.DbPlace, error)
}
