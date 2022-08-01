package schedule

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"cube/lib/database"
	"cube/lib/utils"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var ts = utils.NewTimeService()

func setupTest(wg *sync.WaitGroup, s *Service, done chan struct{}) {
	wg.Add(1)
	go func() {
		var running = true
		for running {
			select {
			case <-done:
				running = false
			case <-s.reschedule:
			}
		}
		wg.Done()
	}()
}

func addTask(s *Service, ID uint64, minute int) {
	sched, _ := Parse(fmt.Sprintf("%v * * * *", minute), ts)
	s.AddTask(&Task{
		Sched: sched,
		ID:    ID,
	})
}

func TestScheduleServiceAddAndRemoveTask(t *testing.T) {
	as := assert.New(t)
	done := make(chan struct{})
	var wg sync.WaitGroup
	db, err := database.New(":memory:")
	if !as.NoError(err) {
		return
	}

	s := NewService(db, ts)
	as.True(s.s.NextSchedule().IsZero())

	setupTest(&wg, s, done)

	minute := time.Now().Minute()
	addTask(s, 2, minute+2)
	as.Less(1*time.Minute, time.Until(s.s.NextSchedule()))
	as.Greater(2*time.Minute, time.Until(s.s.NextSchedule()))

	addTask(s, 1, minute+1)
	as.Less(0*time.Minute, time.Until(s.s.NextSchedule()))
	as.Greater(1*time.Minute, time.Until(s.s.NextSchedule()))

	s.RemoveTasks([]uint64{1})
	as.Less(1*time.Minute, time.Until(s.s.NextSchedule()))
	as.Greater(2*time.Minute, time.Until(s.s.NextSchedule()))

	close(done)
	wg.Wait()
}

func TestScheduleServiceSearchTask(t *testing.T) {
	as := assert.New(t)
	done := make(chan struct{})
	var wg sync.WaitGroup
	db, err := database.New(":memory:")
	if !as.NoError(err) {
		return
	}

	s := NewService(db, ts)

	setupTest(&wg, s, done)

	now := time.Now()
	minute := now.Minute()
	addTask(s, 1, minute+1)
	addTask(s, 2, minute+2)

	task, nextSched := s.SearchTask(1)
	if as.NotNil(task) {
		as.Equal(task.ID, uint64(1))
		as.Equal(minute+1, nextSched.Minute())
	}

	task, nextSched = s.SearchTask(2)
	if as.NotNil(task) {
		as.Equal(task.ID, uint64(2))
		as.Equal(minute+2, nextSched.Minute())
	}

	task, nextSched = s.SearchTask(3)
	as.Nil(task)
}

func TestScheduleOnceTask(t *testing.T) {
	as := assert.New(t)
	done := make(chan struct{})
	exec := make(chan struct{})
	var wg sync.WaitGroup
	var isTaskExecuted bool
	db, err := database.New(":memory:")
	if !as.NoError(err) {
		return
	}

	s := NewService(db, ts)

	wg.Add(1)
	go func() {
		var running = true
		for running {
			select {
			case <-done:
				running = false
			case <-s.reschedule:
			}
		}
		wg.Done()
	}()

	now := time.Now().Add(-1 * time.Minute).Format("2006-01-02 15:04:05")
	record := &database.Alarm{
		UserID:  "1",
		Pattern: now,
		Message: "",
	}

	if tx := db.Save(record); !as.NoError(tx.Error) {
		return
	}

	alarmID := record.AlarmID

	sched, err := Parse(now, ts)
	if !as.NoError(err) {
		return
	}

	s.AddTask(&Task{
		Sched: sched,
		ID:    alarmID,
		Run: func() {
			isTaskExecuted = true
			close(exec)
		},
	})

	s.runOverdueTasks()

	<-exec
	as.True(isTaskExecuted)

	tx := db.First(record, alarmID)
	as.Error(gorm.ErrRecordNotFound, tx.Error)

	close(done)
	wg.Wait()
}
