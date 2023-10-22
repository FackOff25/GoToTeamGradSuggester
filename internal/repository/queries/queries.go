package queries

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
)

type Queries struct {
	Ctx context.Context
	pgxpool.Pool
}


func (q *Queries) GetUserPlaces() (*domain.User, error) {
	return &domain.User{Places: []string{"Moscow"}}, nil
}
