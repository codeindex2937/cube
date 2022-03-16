package schedule

import (
	"context"
	"sync"
	"time"

	"cube/lib/database"
	"cube/lib/logger"

	"gorm.io/gorm"
)

var log = logger.Log

type IService interface {
	AddTask(task *Task)
	RemoveTasks(IDs []uint64)
	SearchTask(ID uint64) (*Task, time.Time)
}

type Task struct {
	Sched             Schedule
	ID                uint64
	Run               func()
	previousTriggered time.Time
}

type Service struct {
	m          sync.Mutex
	s          ScheduleList
	db         *gorm.DB
	reschedule chan bool
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db:         db,
		reschedule: make(chan bool),
	}
}

func (s *Service) AddTask(task *Task) {
	s.m.Lock()
	defer s.m.Unlock()

	s.s.Insert(task.Sched.Next(time.Now()), task.ID, task)
	s.reschedule <- true
}

func (s *Service) RemoveTasks(IDs []uint64) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, ID := range IDs {
		_ = s.s.Remove(ID)
	}
	s.reschedule <- true
}

func (s *Service) SearchTask(ID uint64) (*Task, time.Time) {
	s.m.Lock()
	defer s.m.Unlock()

	task, nextSched, err := s.s.Search(ID)
	if err != nil {
		if err == ErrNoItem {
			return nil, time.Time{}
		} else {
			log.Errorf("search task: %v", err)
		}
	}

	return task.(*Task), nextSched
}

func (s *Service) removeOverdueTask() (overdueTasks []*Task, err error) {
	overdueTasks = []*Task{}

	for !s.s.IsEmpty() && time.Until(s.s.NextSchedule()) < time.Minute {
		value, err := s.s.Pop()
		if err != nil {
			log.Errorf("unexpected nil task\n")
			break
		}

		overdueTasks = append(overdueTasks, value.(*Task))
	}
	return
}

func (s *Service) runOverdueTasks() {
	now := time.Now()
	nextTimeSlot := now.Add(time.Minute)
	overdueTasks, _ := s.removeOverdueTask()

	for _, task := range overdueTasks {
		if task.Sched.IsOnce() {
			s.db.Delete(&database.Alarm{}, task.ID)
			continue
		}

		nextSchedule := task.Sched.Next(nextTimeSlot)
		if task.previousTriggered == nextSchedule {
			log.Errorf("no schedule stepping? %v\n", task)
		} else {
			task.previousTriggered = nextSchedule
			s.s.Insert(nextSchedule, task.ID, task)
		}
	}

	for i := range overdueTasks {
		task := overdueTasks[i]
		go func() {
			task.Run()
		}()
	}
}

func (s *Service) Run(ctx context.Context) {
	running := true

	for running {
		nextSchedule := s.s.NextSchedule()
		for running && nextSchedule.IsZero() {
			select {
			case <-s.reschedule:
				nextSchedule = s.s.NextSchedule()
			case <-ctx.Done():
				running = false
			}
		}

		if !running {
			break
		}

		select {
		case <-s.reschedule:
		case <-ctx.Done():
			running = false
		case <-time.After(time.Until(nextSchedule)):
			s.runOverdueTasks()
		}
	}
}
