package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/FackOff25/GoToTeamGradGoLibs/googleApi"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/repository"
	"github.com/FackOff25/GoToTeamGradSuggester/pkg/config"
)

type UseCase struct {
	repo *repository.Repo
	ctx  context.Context
}

func New(r repository.Repo, ctx context.Context) UsecaseInterface {
	return &UseCase{repo: &r, ctx: ctx}
}

func (uc *UseCase) GetNearbyPlaces(cfg *config.Config, location string, radius int, placeType string) ([]googleApi.Place, error) {
	request := cfg.PlacesApiHost + "place/nearbysearch/" + "json"
	request += "?language=ru"

	request += "&location=" + location

	request += "&radius=" + fmt.Sprintf("%d", radius)

	if placeType != "" {
		request += "&type=" + placeType
	}

	resp, err := http.Get(request)
	if err != nil {
		return []googleApi.Place{}, err
	}

	data, _ := io.ReadAll(resp.Body)
	var result googleApi.NearbyPlacesAnswer
	json.Unmarshal(data, &result)

	if result.Status != googleApi.STATUS_OK {
		return []googleApi.Place{}, errors.New(result.Status)
	}

	return result.Result, nil
}
