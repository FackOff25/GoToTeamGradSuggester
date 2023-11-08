package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/FackOff25/GoToTeamGradGoLibs/googleApi"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/repository"
	"github.com/FackOff25/GoToTeamGradSuggester/pkg/config"
	log "github.com/sirupsen/logrus"
)

type UseCase struct {
	repo *repository.Repo
	ctx  context.Context
	cfg  *config.Config
}

func New(r repository.Repo, ctx context.Context, cfg *config.Config) UsecaseInterface {
	return &UseCase{repo: &r, ctx: ctx, cfg: cfg}
}

func comparePlaces(first, second domain.SuggestPlace) bool {
	return first.SortValue > second.SortValue
}

func isPlaceRight(place googleApi.Place) bool {
	return place.RatingCount > 100 && place.Rating > 2.0
}

func calculateSortValue(user *domain.User, place googleApi.Place) float64 {
	value := float64(0)
	weights := getPlaceTypesWeight()

	for _, placeType := range place.Types {
		pref, ok := user.PlaceTypePreferences[placeType]
		if !ok {
			pref = 1
		}

		weight, ok := weights[placeType]
		if !ok {
			weight = 1
		}

		value += pref * weight
	}

	value += float64(place.Rating * ratingWeight)
	return value
}

func formNearbyPlace(cfg *config.Config, user *domain.User, result googleApi.Place) (domain.SuggestPlace, error) {
	location := domain.ApiLocation{
		Lat: result.Geometry.Location.Lat,
		Lng: result.Geometry.Location.Lng,
	}

	var cover string
	if len(result.Photos) > 0 {
		reference := result.Photos[0].Reference
		cover = cfg.PlacesApiHost + "place/photo?maxwidth=" + strconv.FormatInt(result.Photos[0].Width, 10) + "&photo_reference=" + reference
		result.Photos = result.Photos[1:]
	}

	var photos []string
	for _, photo := range result.Photos {
		cover = cfg.PlacesApiHost + "place/photo?maxwidth=" + strconv.FormatInt(photo.Width, 10) + "&photo_reference=" + photo.Reference
	}

	return domain.SuggestPlace{
		PlaceId:     result.PlaceId,
		Name:        result.Name,
		Location:    location,
		Cover:       cover,
		Photos:      photos,
		Rating:      float64(result.Rating),
		RatingCount: int(result.RatingCount),
		SortValue:   calculateSortValue(user, result),
	}, nil
}

func (uc *UseCase) GetNearbyPlaces(cfg *config.Config, location string, radius int, placeType string, pageToken string) ([]googleApi.Place, string, error) {
	request := cfg.PlacesApiHost + "place/nearbysearch/" + "json"

	if pageToken != "" {
		request += "?pagetoken=" + pageToken
	} else {
		request += "?language=ru"

		request += "&location=" + location

		request += "&radius=" + fmt.Sprintf("%d", radius)

		if placeType != "" {
			request += "&type=" + placeType
		}
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", request, nil)
	if err != nil {
		return []googleApi.Place{}, "", errors.New("Error while creating request: " + err.Error())
	}

	req.Header.Set("Proxy-Header", "go-explore")
	resp, err := client.Do(req)
	if err != nil {
		return []googleApi.Place{}, "", err
	}

	data, _ := io.ReadAll(resp.Body)
	var result googleApi.NearbyPlacesAnswer
	json.Unmarshal(data, &result)

	if result.Status != googleApi.STATUS_OK && result.Status != googleApi.STATUS_ZERO_RESULTS {
		return []googleApi.Place{}, "", errors.New(result.Status)
	}

	return result.Result, result.NextPageToken, nil
}

func (uc *UseCase) GetMergedNearbyPlaces(cfg *config.Config, user *domain.User, location string, radius int, limit int, offset int) ([]domain.SuggestPlace, error) {
	waitGroup := new(sync.WaitGroup)
	result := []googleApi.Place{}

	goroutineFunc := func(placeType string) {
		defer waitGroup.Done()
		typeResult, _, err := uc.GetNearbyPlaces(cfg, location, radius, placeType, "")
		if err != nil {
			log.Errorf("Error during fetching nerby places of type \"%s\": %s", placeType, err)
		}
		result = append(result, typeResult...)
	}

	for _, placeType := range cfg.Categories {
		waitGroup.Add(1)
		go goroutineFunc(placeType)
	}
	waitGroup.Wait()

	return uc.proceedPlaces(cfg, user, result), nil
}

func (uc *UseCase) proceedPlaces(cfg *config.Config, user *domain.User, places []googleApi.Place) []domain.SuggestPlace {
	//TODO: make it parallel
	proceeded := []domain.SuggestPlace{}
	for _, place := range places {
		if isPlaceRight(place) {
			proceededPlace, err := formNearbyPlace(uc.cfg, user, place)
			if err == nil {
				proceeded = append(proceeded, proceededPlace)
			}
		}
	}
	return proceeded
}

func (uc *UseCase) SortPlaces(places []domain.SuggestPlace) []domain.SuggestPlace {
	sort.Slice(places, func(i, j int) bool {
		return comparePlaces(places[i], places[j])
	})
	return places
}

func (uc *UseCase) UniqPlaces(places []domain.SuggestPlace) []domain.SuggestPlace {
	allKeys := make(map[string]bool)
	list := []domain.SuggestPlace{}
	for _, item := range places {
		if _, value := allKeys[item.PlaceId]; !value {
			allKeys[item.PlaceId] = true
			list = append(list, item)
		}
	}
	return list
}
