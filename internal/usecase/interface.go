package usecase

import (
	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
)

type UsecaseInterface interface {
	GetUserPlaces() (*domain.User, error)
}
