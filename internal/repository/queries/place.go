package queries

import (
	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
)

func (q *Queries) GetUserPlaces() (*domain.User, error) {
	return &domain.User{Username: "username"}, nil
}

