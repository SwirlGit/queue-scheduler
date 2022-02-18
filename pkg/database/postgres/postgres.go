package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(cfg *Config, appName string) (*DB, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.fullURL(appName))
	if err != nil {
		return nil, errors.Wrap(err, "pgxpool parse config")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.ConnectConfig(ctx, poolCfg)
	if err != nil {
		return nil, errors.Wrap(err, "pgxpool connect config")
	}

	return &DB{pool: pool}, nil
}

func (d *DB) Pool() *pgxpool.Pool {
	return d.pool
}
