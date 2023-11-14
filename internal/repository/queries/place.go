package queries

import (
	"encoding/json"

	"github.com/FackOff25/GoToTeamGradGoLibs/googleApi"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
	"github.com/jackc/pgx/v5"
)

const (
	addPlaceQuery               = `INSERT INTO places (place_id, types) VALUES ($1, $2);`
	getPlaceTypesByPlaceIdQuery = `SELECT types FROM places WHERE place_id = $1;`
	getPlaceByUUIDQuery         = `SELECT id, place_id, types FROM places WHERE id = $1;`
)

func (q *Queries) AddPlace(gID string, types []string) error {
	typesBytes, _ := json.Marshal(types)
	typesString := string(typesBytes)
	_, err := q.Pool.Exec(q.Ctx, addPlaceQuery, gID, typesString)

	if err != nil {
		return err
	}

	return nil
}

func (q *Queries) GetPlaceById(gID string) (*domain.DbPlace, error) {
	var types string
	row := q.Pool.QueryRow(q.Ctx, getPlaceTypesByPlaceIdQuery, gID)

	err := row.Scan(&types)
	if err != nil {
		return nil, err
	}

	s := make([]string, 0)
	err = json.Unmarshal([]byte(types), &s)
	if err != nil {
		return nil, err
	}

	return &domain.DbPlace{Place_id: gID, Types: s}, nil
}

func (q *Queries) LikePlace(gId, userId, reaction string) error {
	return nil
}

func (q *Queries) SavePlaces(p []googleApi.Place) error {
	for _, v := range p {
		_, err := q.GetPlaceById(v.PlaceId)
		if err != nil {
			if err == pgx.ErrNoRows {
				err := q.AddPlace(v.PlaceId, v.Types)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}
