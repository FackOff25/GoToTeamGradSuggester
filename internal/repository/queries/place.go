package queries

import "github.com/FackOff25/GoToTeamGradGoLibs/googleApi"

func ComparePlaces(first, second googleApi.Place) bool {
	return first.Rating < second.Rating
}
