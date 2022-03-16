package fake

import (
	"time"

	"cube/service/schedule"
)

type ScheduleService struct {
	taskMap map[uint64]*schedule.Task
}

func NewScheduleService() *ScheduleService {
	return &ScheduleService{
		taskMap: make(map[uint64]*schedule.Task),
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
