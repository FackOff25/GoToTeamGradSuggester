package usecase

import (
	"math"

	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
	"github.com/FackOff25/GoToTeamGradSuggester/pkg/gamilton"

	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/FackOff25/GoToTeamGradGoLibs/googleApi"
)

// first element is starting point, second is last (can be the same one)
func sortPlacesForRoute(places []domain.ApiLocation) []domain.ApiLocation {
	return places
	matrix := makeGraphMatrix(places)
	path := gamilton.HungryAlgorythm(matrix)

	var result []domain.ApiLocation
	for _, v := range path {
		result = append(result, places[v])
	}
	return result
}

func getDistanceBetweenPlaces(place1 domain.ApiLocation, place2 domain.ApiLocation) float64 {
	latDist := place1.Lat - place2.Lat
	lngDist := place1.Lng - place2.Lng
	return math.Sqrt(latDist*latDist + lngDist*lngDist)
}

func makeGraphMatrix(places []domain.ApiLocation) [][]float64 {
	matrix := make([][]float64, len(places))
	for i := 0; i < len(places); i++ {
		matrix[i] = make([]float64, len(places))
	}

	for i := range places {
		matrix[i][i] = gamilton.NEAR_INFINITE_NUMBER
		if i < len(places) {
			for j := i + 1; j < len(places); j++ {
				dist := getDistanceBetweenPlaces(places[i], places[j])
				matrix[i][j] = dist
				matrix[j][i] = dist
			}
		}
	}

	return matrix
}


func (uc *UseCase) GetRoute(req *domain.RouteReq) (*domain.Route, error) {
	if req.TravelMode == "" {
		req.TravelMode = domain.TravelModeWalk
	}

	unsortedPlaces := make([]domain.ApiLocation, 0)
	for _, v := range req.Waypoints {
		unsortedPlaces = append(unsortedPlaces, v.Location)
	}

	sortedPlaces := sortPlacesForRoute(unsortedPlaces)
	GreqBody := domain.GrouteRequest{
		TravelMode:    req.TravelMode,
		Intermediates: make([]domain.Gwaypoint, 0),
	}

	for i, v := range sortedPlaces {
		if i == 0 {
			GreqBody.Origin = domain.Gwaypoint{Glocation: domain.Glocation{GlatLng: domain.GlatLng(v)}}
		} else if i == len(sortedPlaces)-1 {
			GreqBody.Destination = domain.Gwaypoint{Glocation: domain.Glocation{GlatLng: domain.GlatLng(v)}}
		} else {
			GreqBody.Intermediates = append(GreqBody.Intermediates, domain.Gwaypoint{Glocation: domain.Glocation{GlatLng: domain.GlatLng(v)}})
		}
	}

	BytesGreqBody, err := json.Marshal(GreqBody)
	if err != nil {
		return nil, err
	}

	request := uc.cfg.RoutesApiHost + "directions/v2:computeRoutes"

	client := &http.Client{}
	Grequest, err := http.NewRequest(http.MethodPost, request, bytes.NewReader(BytesGreqBody))
	if err != nil {
		return nil, err
	}

	Grequest.Header.Set("Proxy-Header", "go-explore")
	resp, err := client.Do(Grequest)
	if err != nil {
		
	}

	data, _ := io.ReadAll(resp.Body)
	var result googleApi.NearbyPlacesAnswer
	json.Unmarshal(data, &result)

	if result.Status != googleApi.STATUS_OK && result.Status != googleApi.STATUS_ZERO_RESULTS {
	}

	return nil, nil
}
