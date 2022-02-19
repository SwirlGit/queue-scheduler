package schedule

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

var ErrNoAvailableJobs = errors.New("no available jobs")

var errAlreadyExists = errors.New("already exists")

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}

const getJobForUpdateQuery = `
	SELECT j.id, q.queue_id, j.date_time, j.action, j.state, j.last_heart_beat
	FROM jobs AS j
	LEFT JOIN queues AS q
		ON j.ref_queue_id = q.id
	WHERE j.date_time < $1 AND state = 'new'
	LIMIT 1
	FOR UPDATE SKIP LOCKED
`

func (s *Storage) TakeJobIntoWork(ctx context.Context) (Job, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Job{}, errors.Wrap(err, "begin tx")
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var job Job
	if err = pgxscan.Get(ctx, tx, &job, getJobForUpdateQuery); errors.Is(err, pgx.ErrNoRows) {
		return Job{}, ErrNoAvailableJobs
	}
	if err != nil {
		return Job{}, errors.Wrap(err, "get job for update")
	}

	if _, err = tx.Exec(ctx, updateJobStateQuery, JobStateRunning, job.ID); err != nil {
		return Job{}, errors.Wrap(err, "set job into running state")
	}

	if err = tx.Commit(ctx); err != nil {
		return Job{}, errors.Wrap(err, "commit tx")
	}

	return job, nil
}

func (s *Storage) FinishJob(ctx context.Context, jobID int64) error {
	return s.updateJobState(ctx, jobID, JobStateDone)
}

const updateJobStateQuery = `UPDATE jobs SET state = $1, last_heart_beat = now() WHERE id = $2`

func (s *Storage) updateJobState(ctx context.Context, jobID int64, state JobState) error {
	_, err := s.pool.Exec(ctx, updateJobStateQuery, state, jobID)
	return errors.Wrap(err, "exec query")
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
