package schedule

import "time"

type Job struct {
	DateTime time.Time
	QueryID  int64
	Action   string
}
