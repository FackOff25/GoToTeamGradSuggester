package usecase

import (
	"fmt"
	"testing"

	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
)

func TestMakingMatrix(t *testing.T) {
	places := []domain.ApiLocation{
		domain.ApiLocation{
			Lat: 10,
			Lng: 10,
		},
		domain.ApiLocation{
			Lat: 5,
			Lng: 5,
		},
		domain.ApiLocation{
			Lat: 0,
			Lng: 0,
		},
	}

	matrix := makeGraphMatrix(places)

	fmt.Printf("%v", matrix)
}
