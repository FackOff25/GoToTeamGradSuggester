package repository

import (
	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
)

type QueriesInterface interface {
	GetUser(id string) (*domain.User, error) 
	AddUser(id string) error
}
