package repository

import (
	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
	"github.com/google/uuid"
)

func (r *Repo) GetUserPlaces(userID uuid.UUID) ([]domain.Place, error) {
	// postgeSQL call

	places := []domain.Place{}

	return places, nil
}
