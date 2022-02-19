package schedule

import "time"

type JobState = string

const (
	JobStateNew     = JobState("new")
	JobStateRunning = JobState("running")
	JobStateDone    = JobState("done")
)

type Job struct {
	ID            int64      `db:"id"`
	QueueID       string     `db:"queue_id"`
	DateTime      time.Time  `db:"date_time"`
	Action        string     `db:"action"`
	State         JobState   `db:"state"`
	LastHeartBeat *time.Time `db:"last_heart_beat"`
}
