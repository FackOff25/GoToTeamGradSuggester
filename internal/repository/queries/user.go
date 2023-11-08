package queries

import "github.com/FackOff25/GoToTeamGradSuggester/internal/domain"

func (q *Queries) GetUser() (*domain.User, error) {
	return &domain.User{Username: "username", PlaceTypePreferences: map[string]float32{
		domain.TypePlaceCafe:   1.0,
		domain.TypePlaceMuseum: 2.2,
		domain.TypePlacePark:   1.5,
	}}, nil
}
