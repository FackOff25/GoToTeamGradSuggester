package queries

import (
	"encoding/json"
	"fmt"
)

const (
	addPlaceQuery          = `INSERT INTO places (place_id, types) VALUES ($1, $2) RETURNING id;`
	getPlaceQueryByPlaceId = `SELECT id, place_id, types FROM places WHERE place_id = $1;`
	getPlaceQueryByUUID    = `SELECT id, place_id, types FROM places WHERE id = $1;`
)

func (q *Queries) AddPlace(gID string, types []string) (string, error) {
	typesBytes, _ := json.Marshal(types)
	fmt.Println(string(typesBytes))
	typesString := string(typesBytes)
	var id string
	row := q.Pool.QueryRow(q.Ctx, addPlaceQuery, gID, typesString)

	err := row.Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (q *Queries) LikePlace(gId, userId, reaction string) error {
	return nil
}
