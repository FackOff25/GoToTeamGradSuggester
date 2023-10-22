package queries

import (
	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
)

type QueriesInterface interface {
	GetUserPlaces() (*domain.User, error) 
}