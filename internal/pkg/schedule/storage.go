package schedule

import (
	"context"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
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
	SELECT j.id, j.date_time, j.action, j.state, j.last_heart_beat, q.id, q.queue_id
	FROM jobs AS j
	INNER JOIN queues AS q
		ON j.ref_queue_id = q.id
	WHERE j.date_time < $1 AND q.state = 'ready'::QUEUE_STATE
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

	if err = s.updateStateWithTx(ctx, tx, job.ID, JobStateRunning, job.Queue.ID, QueueStateBusy); err != nil {
		return Job{}, errors.Wrap(err, "set job into running state")
	}

	if err = tx.Commit(ctx); err != nil {
		return Job{}, errors.Wrap(err, "commit tx")
	}

	return job, nil
}

func (s *Storage) FinishJob(ctx context.Context, job Job) error {
	return s.updateState(ctx, job.ID, JobStateDone, job.Queue.ID, QueueStateReady)
}

func (s *Storage) RenewJob(ctx context.Context, job Job) error {
	return s.updateState(ctx, job.ID, JobStateNew, job.Queue.ID, QueueStateReady)
}

const updateJobStateQuery = `UPDATE jobs SET state = $1, last_heart_beat = now() WHERE id = $2`
const updateQueueStateQuery = `UPDATE queues SET state = $1 WHERE id = $2`

func (s *Storage) updateState(ctx context.Context,
	jobID int64, jobState JobState,
	queueID int64, queueState QueueState) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if err = s.updateStateWithTx(ctx, tx, jobID, jobState, queueID, queueState); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "commit tx")
	}
	return nil
}

func (s *Storage) updateStateWithTx(ctx context.Context, tx pgx.Tx,
	jobID int64, jobState JobState,
	queueID int64, queueState QueueState) error {
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		_, err := tx.Exec(egCtx, updateJobStateQuery, jobState, jobID)
		return errors.Wrap(err, "update job state")
	})
	eg.Go(func() error {
		_, err := tx.Exec(egCtx, updateQueueStateQuery, queueState, queueID)
		return errors.Wrap(err, "update queue state")
	})
	return errors.Wrap(eg.Wait(), "update state")
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

const getRunningJobsForTooLongQuery = `
	SELECT j.id, j.date_time, j.action, j.state, j.last_heart_beat, q.id, q.queue_id
	FROM jobs AS j
	INNER JOIN queues AS q
		ON j.ref_queue_id = q.id
	WHERE j.last_heart_beat IS NOT NULL AND j.last_heart_beat < $1 AND q.state = 'busy'::QUEUE_STATE`

func (s *Storage) GetRunningJobsForTooLong(ctx context.Context, dateTime time.Time) ([]Job, error) {
	var jobs []Job
	if err := pgxscan.Select(ctx, s.pool, &jobs, getRunningJobsForTooLongQuery, dateTime); err != nil {
		return nil, errors.Wrap(err, "pgxscan select")
	}
	return jobs, nil
}
