package queries

import (
	"fmt"

	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
)

const (
	addUserQuery     = "INSERT INTO users (id, ratings) VALUES ($1, $2);"
	getUserByIdQuery = `SELECT id, username, ratings FROM users WHERE id = $1;`
	updateUserQuery  = `UPDATE users SET username = $1, ratings = $2 WHERE id = $3;`
)

func (q *Queries) GetUser(id string) (*domain.User, error) {
	row := q.Pool.QueryRow(q.Ctx, getUserByIdQuery, id)

	user := domain.User{}
	var username, ratings *string

	err := row.Scan(&user.Id, &username, &ratings)
	if err != nil {
		return nil, err
	}

	if username != nil {
		user.Username = *username
	}

	if ratings == nil {
		return nil, fmt.Errorf("error: no user ratings in db")
	}
	ratingsString := *ratings

	m, err := UnmarshalUserRatings(ratingsString)
	if err != nil {
		return nil, err
	}

	user.PlaceTypePreferences = m
	return &user, nil
}

func (q *Queries) AddUser(id string, categories []string) error {
	m := make(map[string]float32)

	for _, v := range categories {
		m[v] = 1.0
	}

	s, err := MarshalUserRatings(m)
	if err != nil {
		return err
	}

	t, err := q.Pool.Exec(q.Ctx, addUserQuery, id, s)
	if err != nil {
		return err
	}

	if !t.Insert() {
		return fmt.Errorf("add user error: wrong sql result")
	}

	return nil
}

func (q *Queries) ApplyUserReactionToPlace(uuid string, placeId string, reaction string, types []string) error {
	u, err := q.GetUser(uuid)
	if err != nil {
		return err
	}

	var multiplier float32

	switch reaction {
	case domain.ReactionVisited:
		multiplier = 3
	case domain.ReactionLike:
		multiplier = 5
	case domain.ReactionRefuse:
		multiplier = 0.5
	default:
		multiplier = 1
	}

	for _, v := range types {
		_, ok := u.PlaceTypePreferences[v]
		if ok {
			u.PlaceTypePreferences[v] *= multiplier
		} else {
			u.PlaceTypePreferences[v] = 1
		}
	}

	err = q.UpdateUser(u)

	if err != nil {
		return fmt.Errorf("apply user reaction to place error: %s", err.Error())
	}

	return nil

}

func (q *Queries) UpdateUser(user *domain.User) error {
	s, err := MarshalUserRatings(user.PlaceTypePreferences)
	if err != nil {
		return err
	}

	t, err := q.Pool.Exec(q.Ctx, updateUserQuery, &user.Username, &s, &user.Id)
	if !t.Update() {
		return fmt.Errorf("update user error: wrong sql operation result")
	}

	return err
}
