package usecase

import "github.com/FackOff25/GoToTeamGradSuggester/internal/domain"

const ratingWeight = 1

func getPlaceTypesWeight() map[string]float64 {
	return map[string]float64{
		domain.TypePlacePark: 1.5,
	}
}
