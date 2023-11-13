package usecase

import "github.com/FackOff25/GoToTeamGradSuggester/internal/domain"

const ratingWeight = 5

func getPlaceTypesWeight() map[string]float32 {
	return map[string]float32{
		domain.TypePlacePark: 2.5,
	}
}
