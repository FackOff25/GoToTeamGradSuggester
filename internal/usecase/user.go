package usecase

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
	"github.com/jackc/pgx/v5"
)

func (uc *UseCase) AddUser(id string) error {
	_, err := uc.repo.GetUser(id)

	if err != nil {
		if err == pgx.ErrNoRows {
			err := uc.repo.AddUser(id)
			if err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}

	return nil
}

func (uc *UseCase) GetUser(uuid string) (*domain.User, error) {
	return uc.repo.GetUser(uuid)
}

func (uc *UseCase) ApplyUserReactionToPlace(uuid string, placeId string, reaction string) error {
	getPlaceTypes(placeId)
	return nil
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
