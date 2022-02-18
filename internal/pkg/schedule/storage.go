package schedule

import "github.com/jackc/pgx/v4/pgxpool"

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}
