package queries

import "github.com/FackOff25/GoToTeamGradSuggester/internal/domain"

const (
	addUserQuery = "INSERT INTO users (id) VALUES ($1);"
)


func (q *Queries) GetUser(id string) (*domain.User, error) {
	return &domain.User{Username: "username", PlaceTypePreferences: map[string]float32{
		domain.TypePlaceCafe: 1.0, 
		domain.TypePlaceMuseum: 2.2,
		domain.TypePlacePark: 1.5,
	}}, nil
}

func (q *Queries) AddUser(id string) error {
	_, err := q.Pool.Query(q.Ctx, addUserQuery, id)
	return err
}

