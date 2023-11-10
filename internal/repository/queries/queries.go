package queries

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Queries struct {
	Ctx context.Context
	pgxpool.Pool
}

// func New(ctx context.Context, pool pgxpool.Pool) *repository.QueriesInterface {
// 	return *Queries{ctx: ctx, Pool: pool}
// }

func UnmarshalUserRatings(s string) (map[string]float32, error) {
	jsonbytes := []byte(s)
	m := make(map[string]float32)
	err := json.Unmarshal(jsonbytes, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func MarshalUserRatings(m map[string]float32) (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
