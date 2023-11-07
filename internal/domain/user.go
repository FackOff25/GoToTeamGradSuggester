package domain

import "github.com/google/uuid"

type User struct {
	Id                   uuid.UUID
	Username             string
	PlaceTypePreferences map[string]float32
}

type NewReactionRequest struct {
	PlaceId  string `json:"place_id"`
	Reaction string `json:"reaction"`
}
