package queries

import (
	"encoding/json"
	"fmt"

	"github.com/FackOff25/GoToTeamGradGoLibs/googleApi"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
	"github.com/jackc/pgx/v5"
)

const (
	addPlaceQuery               = `INSERT INTO places (place_id, types) VALUES ($1, $2);`
	getPlaceTypesByPlaceIdQuery = `SELECT types FROM places WHERE place_id = $1;`
	getPlaceUuidQuery           = `SELECT id FROM places WHERE place_id = $1;`
	updateReactionQuery         = `UPDATE users_places SET %s = $1 WHERE user_id = $2 AND place_id = $3;`
	insertReactionQuery         = `UPDATE INTO users_places(user_id, place_id, ) SET %s = $1 WHERE user_id = $2 AND place_id = $3;`
	getUserPlaceReaction        = `SELECT like_mark, visited_mark FROM users_places WHERE place_id = $1 AND user_id = $2;`
	insertDefaultReactionQuery  = `INSERT INTO users_places(place_id, user_id, like_mark, visited_mark) VALUES ($1, $2, false, false);`
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

func (q *Queries) GetPlaceUuid(gID string) (string, error) {
	var uuid string
	row := q.Pool.QueryRow(q.Ctx, getPlaceUuidQuery, gID)

	err := row.Scan(&uuid)
	if err != nil {
		return "", err
	}

	return uuid, nil
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

func (q *Queries) GetUserReaction(userId string, placeId string) (bool, bool, error) {
	var likeFlag, visitedFlag bool
	r := q.Pool.QueryRow(q.Ctx, getUserPlaceReaction, &placeId, &userId)

	err := r.Scan(&likeFlag, &visitedFlag)
	return likeFlag, visitedFlag, err
}

func (q *Queries) SaveUserReaction(userId, placeUuid, reaction string) error {
	var columnName string
	reactionFlag := false

	if reaction == domain.ReactionLike || reaction == domain.ReactionUnlike {
		columnName = "like_mark"
	} else if reaction == domain.ReactionVisited || reaction == domain.ReactionUnvisited {
		columnName = "visited_mark"
	} else {
		return nil
	}

	if reaction == domain.ReactionLike || reaction == domain.ReactionVisited {
		reactionFlag = true
	}

	_, _, err := q.GetUserReaction(userId, placeUuid)
	if err != nil {
		if err == pgx.ErrNoRows {
			_, err := q.Pool.Exec(q.Ctx, "INSERT INTO users_places(place_id, user_id, like_mark, visited_mark) VALUES ($1, $2, false, false);", placeUuid, userId)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	_, err = q.Pool.Exec(q.Ctx, fmt.Sprintf(updateReactionQuery, columnName), reactionFlag, userId, placeUuid)
	return err
}
