package schedule

import (
	"context"
	"sync"
	"time"

	"github.com/SwirlGit/queue-scheduler/internal/pkg/schedule"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const defaultCheckDuration = 30 * time.Second

type scheduleStorage interface {
	TakeJobIntoWork(ctx context.Context) (schedule.Job, error)
	FinishJob(ctx context.Context, job schedule.Job) error
}

type Service struct {
	logger          *zap.Logger
	scheduleStorage scheduleStorage

	checkDuration time.Duration

	currentWorkers int
	doneChan       chan struct{}
	stopWaitGroup  sync.WaitGroup
}

func NewService(logger *zap.Logger, scheduleStorage scheduleStorage, checkDuration time.Duration) *Service {
	if checkDuration == 0 {
		checkDuration = defaultCheckDuration
	}
	return &Service{
		logger:          logger,
		scheduleStorage: scheduleStorage,
		checkDuration:   checkDuration,
		doneChan:        make(chan struct{}),
	}
}

func (s *Service) Start(workersAmount int) error {
	if s.currentWorkers > 0 {
		return errors.New("already started")
	}
	s.AddWorkers(workersAmount)
	return nil
}

func (s *Service) AddWorkers(amount int) {
	for i := 0; i < amount; i++ {
		s.currentWorkers++
		s.stopWaitGroup.Add(1)
		go func() {
			defer s.stopWaitGroup.Done()
			s.doUntilStop()
		}()
	}
}

func (s *Service) RemoveWorkers(amount int) error {
	if amount > s.currentWorkers {
		return errors.Errorf("current workers amount = %v is less than requested = %v to stop",
			s.currentWorkers, amount)
	}
	for i := 0; i < amount; i++ {
		s.doneChan <- struct{}{}
	}
	return nil
}

func (s *Service) Stop() {
	close(s.doneChan)
	s.stopWaitGroup.Wait()
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
	job, err := s.scheduleStorage.TakeJobIntoWork(ctx)
	if errors.Is(err, schedule.ErrNoAvailableJobs) {
		return
	}
	if err != nil {
		s.logger.Error("failed to take job into work", zap.Error(err))
		return
	}

	s.doJob(job)

	if err = s.scheduleStorage.FinishJob(ctx, job); err != nil {
		s.logger.Error("failed to finish job", zap.Error(err))
		return
	}
}

func (s *Service) doJob(job schedule.Job) {
	s.logger.Info("job is done", zap.String("action", job.Action))
}
