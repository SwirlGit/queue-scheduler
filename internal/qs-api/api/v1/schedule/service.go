package schedule

type Service struct {

}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ScheduleJob(_ Job) error  {
	return nil
}