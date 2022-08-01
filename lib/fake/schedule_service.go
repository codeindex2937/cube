package fake

import (
	"time"

	"cube/lib/utils"
	"cube/service/schedule"
)

type ScheduleService struct {
	taskMap map[uint64]*schedule.Task
	ts      *utils.TimeService
}

func NewScheduleService(ts *utils.TimeService) *ScheduleService {
	return &ScheduleService{
		taskMap: make(map[uint64]*schedule.Task),
		ts:      ts,
	}
}

func (s *ScheduleService) AddTask(task *schedule.Task) {
	s.taskMap[task.ID] = task
}

func (s *ScheduleService) RemoveTasks(IDs []uint64) {
	for _, ID := range IDs {
		delete(s.taskMap, ID)
	}
}

func (s *ScheduleService) SearchTask(ID uint64) (*schedule.Task, time.Time) {
	return s.taskMap[ID], time.Time{}
}

func (s *ScheduleService) ExistTask(ID uint64) bool {
	_, ok := s.taskMap[ID]
	return ok
}

func (s *ScheduleService) Parse(t string) (sched schedule.Schedule, err error) {
	return schedule.Parse(t, s.ts)
}
