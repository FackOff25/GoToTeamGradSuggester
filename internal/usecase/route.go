package usecase

import (
	"fmt"
	"math"

	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
	"github.com/FackOff25/GoToTeamGradSuggester/pkg/gamilton"
)

// first element is starting point, second is last (can be the same one)
func SortPlacesForRoute(places []domain.ApiLocation) []domain.ApiLocation {
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

func (uc *UseCase) PrepareGreq(req *domain.RouteReq) (*domain.GrouteRequest, error) {
	if req.TravelMode == "" {
		req.TravelMode = domain.TravelModeWalk
	}

	unsortedPlaces := make([]domain.ApiLocation, 0)
	unsortedPlaces = append(unsortedPlaces, req.Start)
	unsortedPlaces = append(unsortedPlaces, req.End)

	for _, v := range req.Waypoints {
		unsortedPlaces = append(unsortedPlaces, v.Location)
	}

	sortedPlaces := SortPlacesForRoute(unsortedPlaces)
	if len(sortedPlaces) == 0 {
		return nil, fmt.Errorf("error: empty route slice")
	}

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

	return &GreqBody, nil
}

func (uc *UseCase) GetRoute(gRoute *domain.GrouteResp, travelMode string) (*domain.Route, error) {
	clientResp := domain.Route{TravelMode: travelMode, Polylines: make([]domain.Polyline, 0)}

	for _, v := range gRoute.Routes {
		for _, val := range v.Legs {
			clientResp.Polylines = append(clientResp.Polylines, domain.Polyline{PolylineString: val.Polyline.EncodedPolyline})
		}
	}

	return &clientResp, nil
}
