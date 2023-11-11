package domain

const (
	ErrorUserAlreadyExists = "user already exists"
)

type User struct {
	Id                   string             `json:"id,omitempty"`
	Username             string             `json:"username,omitempty"`
	PlaceTypePreferences map[string]float32 `json:"preferences,omitempty"`
}

type NewReactionRequest struct {
	PlaceId  string `json:"place_id"`
	Reaction string `json:"reaction"`
}
