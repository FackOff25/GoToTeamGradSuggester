package usecase

import (
	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
)

type UsecaseInterface interface {
	GetUser() (*domain.User, error)
}
