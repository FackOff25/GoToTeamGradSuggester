package domain

type User struct {
	Id string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	PlaceTypePreferences map[string]float64 `json:"preferences,omitempty"`
}
