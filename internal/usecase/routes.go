package usecase

import "github.com/FackOff25/GoToTeamGradSuggester/internal/domain"

type PlaceToSort interface {
	GetLocation() domain.ApiLocation
}

func sortPlacesForRoute(places []PlaceToSort) []PlaceToSort {
	return places
}
