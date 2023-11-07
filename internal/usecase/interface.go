package usecase

import (
	"github.com/FackOff25/GoToTeamGradGoLibs/googleApi"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
	"github.com/FackOff25/GoToTeamGradSuggester/pkg/config"
)

type UsecaseInterface interface {
	GetUser() (*domain.User, error)
	GetNearbyPlaces(cfg *config.Config, location string, radius int, placeType string, pageToken string) ([]googleApi.Place, string, error)
	GetMergedNearbyPlaces(cfg *config.Config, location string, radius int, limit int, offset int) ([]domain.SuggestPlace, error)
	SortPlaces(places []domain.SuggestPlace) []domain.SuggestPlace
	UniqPlaces(places []domain.SuggestPlace) []domain.SuggestPlace
}
