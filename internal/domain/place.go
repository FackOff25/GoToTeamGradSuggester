package domain

import (
	"github.com/google/uuid"
)

const (
	TypePlacePark   = "park"
	TypePlaceCafe   = "cafe"
	TypePlaceMuseum = "museum"

	ReactionLike      = "like"
	ReactionVisited   = "visited"
	ReactionRefuse    = "refuse"
	ReactionUnlike    = "unlike"
	ReactionUnvisited = "unvisited"
	ReactionUnrefuse  = "unrefuse"

	CategoryLikeRus      = "Избранное"
	CategoryUnvisitedRus = "Не посещенное"
)

type GetCategoriesResponse struct {
	Categories []string `json:"categories"`
}

type ApiLocation struct {
	Lat float64 `json:"lat,omitempty"`
	Lng float64 `json:"lng,omitempty"`
}

type ApiPlace struct {
	Id             uuid.UUID   `json:"id,omitempty"`
	Cover          string      `json:"cover,omitempty"`
	Rating         float32     `json:"rating,omitempty"`
	RatingCount    int         `json:"rating_count,omitempty"`
	Name           string      `json:"name,omitempty"`
	Location       ApiLocation `json:"location,omitempty"`
	PlaceId        string      `json:"place_id,omitempty"`
	ApiRatingCount int         `json:"user_ratings_total,omitempty"`
	// Types          []string  `json:"types,omitempty"`
}

type DbPlace struct {
	Place_id string
	Types    []string
}

type SuggestPlace struct {
	PlaceId     string      `json:"place_id,omitempty"`
	Name        string      `json:"name,omitempty"`
	Location    ApiLocation `json:"location,omitempty"`
	Cover       string      `json:"cover,omitempty"`
	Photos      []string    `json:"photos,omitempty"`
	Rating      float32     `json:"rating,omitempty"`
	RatingCount int         `json:"rating_count,omitempty"`
	Reaction    []string    `json:"reactions,omitempty"`
	SortValue   float32     `json:"-"`
}

func GetFilterReactionsMap() map[string]string {
	return map[string]string{
		CategoryLikeRus:      ReactionLike,
		CategoryUnvisitedRus: ReactionUnvisited,
	}
}
