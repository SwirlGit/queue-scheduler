package schedule

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

var errAlreadyExists = errors.New("already exists")

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}

const inertJobQuery = `INSERT INTO jobs (ref_queue_id, date_time, action) VALUES ($1, $2, $3)`

func (s *Storage) InsertJob(ctx context.Context, job Job) error {
	internalQueueID, err := s.getQueueInternalIDOrCreate(ctx, job.QueueID)
	if err != nil {
		return errors.Wrap(err, "get queue internal id or create")
	}
	if _, err = s.pool.Exec(ctx, inertJobQuery, internalQueueID, job.DateTime, job.Action); err != nil {
		return errors.Wrap(err, "exec query")
	}
	return nil
}

func (s *Storage) getQueueInternalIDOrCreate(ctx context.Context, queueID string) (int64, error) {
	id, err := s.getQueueInternalID(ctx, queueID)
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return 0, errors.Wrap(err, "get queue internal id")
	}

	if id, err = s.createQueue(ctx, queueID); err == nil {
		return id, nil
	}
	if !errors.Is(err, errAlreadyExists) {
		return 0, errors.Wrap(err, "create queue")
	}

	if id, err = s.getQueueInternalID(ctx, queueID); err != nil {
		return 0, errors.Wrap(err, "get existed queue id")
	}
	return id, nil
}

const getQueueInternalIDQuery = `SELECT id FROM queues WHERE queue_id = $1`

func (s *Storage) getQueueInternalID(ctx context.Context, queueID string) (int64, error) {
	var id int64
	err := pgxscan.Get(ctx, s.pool, &id, getQueueInternalIDQuery, queueID)
	return id, errors.Wrap(err, "pgxscan get")
}

const createQueueQuery = `INSERT INTO queues (queue_id) VALUES ($1) ON CONFLICT (queue_id) DO NOTHING RETURNING id`

func (s *Storage) createQueue(ctx context.Context, queueID string) (int64, error) {
	var id int64
	err := pgxscan.Get(ctx, s.pool, &id, createQueueQuery, queueID)
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return 0, errors.Wrap(err, "pgxscan get")
	}
	return 0, errAlreadyExists
}
