package schedule

import (
	"context"

	"github.com/SwirlGit/queue-scheduler/internal/pkg/schedule"
	"github.com/pkg/errors"
)

type scheduleStorage interface {
	InsertJob(ctx context.Context, job schedule.Job) error
}

type Service struct {
	scheduleStorage scheduleStorage
}

func NewService(scheduleStorage scheduleStorage) *Service {
	return &Service{scheduleStorage: scheduleStorage}
}

func (s *Service) ScheduleJob(ctx context.Context, job Job) error {
	return errors.Wrap(s.scheduleStorage.InsertJob(ctx, schedule.Job{
		QueueID:  job.QueueID,
		DateTime: job.DateTime,
		Action:   job.Action,
	}), "insert job into storage")
}
