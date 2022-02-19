package schedule

type QueueState = string

const (
	QueueStateReady = JobState("ready")
	QueueStateBusy  = JobState("busy")
)

type Queue struct {
	ID      int64      `db:"id"`
	QueueID string     `db:"queue_id"`
	State   QueueState `db:"state"`
}
