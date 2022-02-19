package checker

import (
	"context"
	"sync"
	"time"

	"github.com/SwirlGit/queue-scheduler/internal/pkg/schedule"
)

const (
	defaultCheckDuration = 1 * time.Minute
	maxRunningDuration   = 5 * time.Minute
)

type scheduleStorage interface {
	GetRunningJobsForTooLong(ctx context.Context, dateTime time.Time) ([]schedule.Job, error)
	RenewJob(ctx context.Context, job schedule.Job) error
}

type Service struct {
	scheduleStorage scheduleStorage
	checkDuration   time.Duration
	doneChan        chan struct{}
}

func NewService(scheduleStorage scheduleStorage, checkDuration time.Duration) *Service {
	if checkDuration == 0 {
		checkDuration = defaultCheckDuration
	}
	return &Service{
		scheduleStorage: scheduleStorage,
		checkDuration:   checkDuration,
		doneChan:        make(chan struct{}),
	}
}

func (s *Service) Start() {
	go s.doUntilStop()
}

func (s *Service) Stop() {
	close(s.doneChan)
}

func (s *Service) doUntilStop() {
	ticker := time.NewTicker(s.checkDuration)
	defer ticker.Stop()

	for {
		select {
		case <-s.doneChan:
			return
		case <-ticker.C:
			s.do()
		}
	}
}

func (s *Service) do() {
	ctx := context.Background()
	dateTime := time.Now().Add(maxRunningDuration)
	jobs, err := s.scheduleStorage.GetRunningJobsForTooLong(ctx, dateTime)
	if err != nil {
		// TODO; log
		return
	}

	if len(jobs) == 0 {
		return
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10)
	for i := range jobs {
		i := i
		semaphore <- struct{}{}
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
				<-semaphore
			}()
			if err := s.scheduleStorage.RenewJob(ctx, jobs[i]); err != nil {
				// TODO: log
			}
		}()
	}
	wg.Wait()
	close(semaphore)
}
