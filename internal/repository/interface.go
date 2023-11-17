package repository

import (
	"github.com/FackOff25/GoToTeamGradGoLibs/googleApi"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
)

type QueriesInterface interface {
	GetUser(id string) (*domain.User, error)
	AddUser(id string, categories []string) error
	UpdateUserPreferences(uuid string, placeId string, reaction string, types []string) error // update in users table [metrics]
	SavePlaces([]googleApi.Place) error
	AddPlace(gID string, types []string) error
	GetPlaceById(gID string) (*domain.DbPlace, error)
	SaveUserReaction(userId, placeUuid, reaction string) error // insert/update in users_places table [interconnection]
	GetPlaceUuid(gID string) (string, error)
	GetUserReaction(userId string, placeId string) (bool, bool, error)
}
