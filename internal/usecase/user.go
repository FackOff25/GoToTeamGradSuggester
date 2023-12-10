package usecase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
	"github.com/jackc/pgx/v5"
)

func (uc *UseCase) AddUser(id string) error {
	u, err := uc.repo.GetUser(id)

	if err != nil {
		if err == pgx.ErrNoRows {
			err := uc.repo.AddUser(id, uc.cfg.Categories)
			if err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}

	if u != nil {
		return fmt.Errorf(domain.ErrorUserAlreadyExists)
	}

	return nil
}

func (uc *UseCase) GetUser(uuid string) (*domain.User, error) {
	return uc.repo.GetUser(uuid)
}

func (uc *UseCase) ApplyUserReactionToPlace(uuid string, placeId string, reaction string) error {
	var types []string
	p, err := uc.repo.GetPlaceById(placeId)
	types = p.Types
	if err != nil {
		if err == pgx.ErrNoRows {
			types = getPlaceTypes(placeId)
		} else {
			return err
		}
	}

	if len(types) == 0 {
		return fmt.Errorf("no place types for place with id: %s", placeId)
	}

	err = uc.repo.UpdateUserPreferences(uuid, placeId, reaction, types)
	if err != nil {
		return err
	}

	err = uc.SaveUserReaction(uuid, placeId, reaction)

	return err
}

func (uc *UseCase) SaveUserReaction(userUuid, placeId, reaction string) error {
	_, err := uc.GetUser(userUuid)
	if err != nil {
		return err
	}

	placeUuid, err := uc.repo.GetPlaceUuid(placeId)
	if err != nil {
		return err
	}

	err = uc.repo.SaveUserReaction(userUuid, placeUuid, reaction)
	return err
}

func (uc *UseCase) ApplyOnboardingReactions(userId string, types []string) error {
	u, err := uc.repo.GetUser(userId) // current user's ratings
	if err != nil {
		return err
	}

	for i := 0; i < len(types); i++ {
		_, ok := u.PlaceTypePreferences[types[i]]
		if ok {
			u.PlaceTypePreferences[types[i]] = 2
		}
	}

	err = uc.repo.UpdateUser(u) // saving new ratings
	return err
}

//
// TEMPORARY
//

type PlaceInfo struct {
	Tags []string `json:"tags,omitempty"`
}

func getPlaceTypes(placeId string) []string {
	request := "http://go-explore.online/api/v1/places/info?place_id=" + placeId

	client := &http.Client{}
	req, err := http.NewRequest("GET", request, nil)
	if err != nil {
		return []string{}
	}

	req.Header.Set("Proxy-Header", "go-explore")
	resp, err := client.Do(req)
	if err != nil {
		return []string{}
	}

	data, _ := io.ReadAll(resp.Body)
	var result PlaceInfo
	json.Unmarshal(data, &result)

	return result.Tags
}
