package usecase

import (
	"math"

	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
	"github.com/FackOff25/GoToTeamGradSuggester/pkg/gamilton"
)

type PlaceToSort interface {
	GetLocation() domain.ApiLocation
}

// first element is starting point, second is last (can be the same one)
func sortPlacesForRoute(places []PlaceToSort) []PlaceToSort {
	matrix := makeGraphMatrix(places)
	path := gamilton.HungryAlgorythm(matrix)

	var result []PlaceToSort
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

func makeGraphMatrix(places []PlaceToSort) [][]float64 {
	matrix := make([][]float64, len(places))
	for i := 0; i < len(places); i++ {
		matrix[i] = make([]float64, len(places))
	}

	for i := range places {
		matrix[i][i] = gamilton.NEAR_INFINITE_NUMBER
		if i < len(places) {
			for j := i + 1; j < len(places); j++ {
				dist := getDistanceBetweenPlaces(places[i].GetLocation(), places[j].GetLocation())
				matrix[i][j] = dist
				matrix[j][i] = dist
			}
		}
	}

	return matrix
}
