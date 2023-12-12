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

func calculateSortValue(user *domain.User, place googleApi.Place) float32 {
	value := float32(0)
	weights := getPlaceTypesWeight()

	parameterCounter := 1
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
		parameterCounter++
	}

	ratingWeight := float32(place.Rating * ratingWeight)
	if ratingWeight > 100 {
		ratingWeight = 100
	}
	value += ratingWeight
	parameterCounter++

	value = value / float32(parameterCounter)
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
		Rating:      float32(result.Rating),
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

func (uc *UseCase) GetMergedNearbyPlaces(cfg *config.Config, user *domain.User, location string, radius int, limit int, offset int, types []string, reactions []string) ([]domain.SuggestPlace, error) {
	waitGroup := new(sync.WaitGroup)
	result := []googleApi.Place{}

	goroutineFunc := func(placeType string) { // TODO: fix race condition
		defer waitGroup.Done()
		typeResult, _, err := uc.GetNearbyPlaces(cfg, location, radius, placeType, "")
		if err != nil {
			log.Errorf("Error during fetching nerby places of type \"%s\": %s", placeType, err)
		}
		result = append(result, typeResult...)
	}

	if len(types) == 0 {
		types = cfg.Categories
	}

	for _, placeType := range types {
		waitGroup.Add(1)
		go goroutineFunc(placeType)
	}
	waitGroup.Wait()

	err := uc.repo.SavePlaces(result)
	if err != nil {
		return nil, err
	}

	return uc.proceedPlaces(cfg, user, result, reactions), nil
}

func (uc *UseCase) proceedPlaces(cfg *config.Config, user *domain.User, places []googleApi.Place, reactions []string) []domain.SuggestPlace {
	//TODO: make it parallel

	likeNeededFlag := false
	unvisitedNeededFlag := false
	if contains(reactions, domain.ReactionLike) {
		likeNeededFlag = true
	}
	if contains(reactions, domain.ReactionUnvisited) {
		unvisitedNeededFlag = true
	}

	proceeded := []domain.SuggestPlace{}
	for _, place := range places {
		if isPlaceRight(place) {
			proceededPlace, err := formNearbyPlace(uc.cfg, user, place)
			if err == nil {
				placeUuid, err := uc.repo.GetPlaceUuid(proceededPlace.PlaceId)
				if placeUuid != "" && err == nil {
					likeFlag, visitedFlag, err := uc.repo.GetUserReaction(user.Id, placeUuid)
					if err == nil {
						reactions := make([]string, 0)
						if likeFlag {
							reactions = append(reactions, domain.ReactionLike)
							proceededPlace.SortValue *= 2
						}
						if visitedFlag {
							reactions = append(reactions, domain.ReactionVisited)
							proceededPlace.SortValue /= 1.5
						}
						proceededPlace.Reaction = reactions
					}
				}
				if (contains(proceededPlace.Reaction, domain.ReactionLike) && likeNeededFlag) ||
					(!contains(proceededPlace.Reaction, domain.ReactionVisited) && unvisitedNeededFlag) || (len(reactions) == 0) {
					proceeded = append(proceeded, proceededPlace)
				}
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

func contains(a []string, s string) bool {
	for i := 0; i < len(a); i++ {
		if a[i] == s {
			return true
		}
	}
	return false
}
