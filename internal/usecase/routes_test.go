package usecase

import (
	"fmt"
	"testing"

	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
)

type TestPlace struct {
	location domain.ApiLocation
}

func convertToPlaceToSort(places []TestPlace) []PlaceToSort {
	result := make([]PlaceToSort, len(places))
	for i, v := range places {
		result[i] = v
	}
	return result
}

func (pl TestPlace) GetLocation() domain.ApiLocation {
	return pl.location
}

func TestMakingMatrix(t *testing.T) {
	places := []TestPlace{
		{
			location: domain.ApiLocation{
				Lat: 10,
				Lng: 10,
			},
		},
		{
			location: domain.ApiLocation{
				Lat: 5,
				Lng: 5,
			},
		},
		{
			location: domain.ApiLocation{
				Lat: 0,
				Lng: 0,
			},
		},
	}

	converted := convertToPlaceToSort(places)
	matrix := makeGraphMatrix(converted)

	fmt.Printf("%v", matrix)
}
