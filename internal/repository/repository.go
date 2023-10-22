package repository

import (
	// "github.com/jackc/pgx/v5"
	"context"

	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/repository/queries"
)

type Repo struct {
	ctx context.Context
	queries.QueriesInterface
}

func New(queries queries.QueriesInterface, ctx context.Context) *Repo {
	return &Repo{
		QueriesInterface: queries,
		ctx: ctx,
	}
}

func (r *Repo) GetUserPlaces() (*domain.User, error) {
	return r.QueriesInterface.GetUserPlaces()
}
