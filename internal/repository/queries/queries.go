package queries

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Queries struct {
	Ctx context.Context
	pgxpool.Pool
}

// func New(ctx context.Context, pool pgxpool.Pool) *repository.QueriesInterface {
// 	return *Queries{ctx: ctx, Pool: pool}
// }
