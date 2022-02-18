package schedule

import "time"

type Job struct {
	QueueID  string
	DateTime time.Time
	Action   string
}
