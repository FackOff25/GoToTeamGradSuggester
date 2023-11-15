package domain

const (
	TravelModeWalk       = "WALK"
	TravelModeDrive      = "DRIVE"
	TravelModeBicycle    = "BICYCLE"
	TravelModeTwoWheeler = "TWO_WHEELER"
	TravelModeTransit    = "TRANSIT"
)

type Waypoint struct {
	PlaceId  string      `json:"place_id,omitempty"`
	Location ApiLocation `json:"location,omitempty"`
}

type RouteReq struct {
	TravelMode string      `json:"travel_mode,omitempty"`
	Start      ApiLocation `json:"start,omitempty"`
	End        ApiLocation `json:"end,omitempty"`
	Waypoints  []Waypoint  `json:"waypoints,omitempty"`
}

type SortPlacesReq struct {
	Start     ApiLocation   `json:"start,omitempty"`
	End       ApiLocation   `json:"end,omitempty"`
	Waypoints []ApiLocation `json:"waypoints,omitempty"`
}

type SortPlacesResp struct {
	Start     ApiLocation   `json:"start,omitempty"`
	End       ApiLocation   `json:"end,omitempty"`
	Waypoints []ApiLocation `json:"waypoints,omitempty"`
}

type Polyline struct {
	PolylineString string `json:"polyline,omitempty"`
}

type Route struct {
	TravelMode string     `json:"travel_mode,omitempty"`
	Polylines  []Polyline `json:"route,omitempty"`
}

// G req

type GlatLng struct {
	Lat float64 `json:"latitude,omitempty"`
	Lng float64 `json:"longitude,omitempty"`
}

type Glocation struct {
	GlatLng `json:"latLng,omitempty"`
}

type Gwaypoint struct {
	Glocation `json:"location,omitempty"`
}

type GPlace struct {
	Gwaypoint `json:"waypoint,omitempty"`
}

type GrouteRequest struct {
	Origin        Gwaypoint   `json:"origin,omitempty"`
	Destination   Gwaypoint   `json:"destination,omitempty"`
	Intermediates []Gwaypoint `json:"intermediates,omitempty"`
	TravelMode    string      `json:"travelMode,omitempty"`
}

// G resp

type GpolyLine struct {
	EncodedPolyline string `json:"encodedPolyline,omitempty"`
}

type Gleg struct {
	Polyline GpolyLine `json:"polyline,omitempty"`
}

type Groute struct {
	Legs []Gleg `json:"legs,omitempty"`
}

type GrouteResp struct {
	Routes []Groute `json:"routes,omitempty"`
}
