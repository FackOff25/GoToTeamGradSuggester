package usecase

import "github.com/FackOff25/GoToTeamGradSuggester/internal/domain"

func sortPlacesForRoute(places []domain.ApiLocation) []domain.ApiLocation {
	return places
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
	Greq := domain.GrouteRequest{
		TravelMode:    req.TravelMode,
		Intermediates: make([]domain.Gwaypoint, 0),
	}

	for i, v := range sortedPlaces {
		if i == 0 {
			Greq.Origin = domain.Gwaypoint{Glocation: domain.Glocation{GlatLng: domain.GlatLng{Lat: v.Lat, Lng: v.Lng}}}
		} else if i == len(sortedPlaces)-1 {
			Greq.Destination = domain.Gwaypoint{Glocation: domain.Glocation{GlatLng: domain.GlatLng{Lat: v.Lat, Lng: v.Lng}}}
		} else {
			Greq.Intermediates = append(Greq.Intermediates, domain.Gwaypoint{Glocation: domain.Glocation{GlatLng: domain.GlatLng{Lat: v.Lat, Lng: v.Lng}}})
		}
	}

	return nil, nil
}
