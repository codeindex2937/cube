package schedule

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"cube/lib/database"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

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
	sched, _ := Parse(fmt.Sprintf("%v * * * *", minute))
	s.AddTask(&Task{
		Sched: sched,
		ID:    ID,
	})
}

func TestScheduleServiceAddAndRemoveTask(t *testing.T) {
	done := make(chan struct{})
	var wg sync.WaitGroup
	db, err := database.New(":memory:")
	if !assert.NoError(t, err) {
		return
	}

	s := NewService(db)
	assert.True(t, s.s.NextSchedule().IsZero())

	setupTest(&wg, s, done)

	minute := time.Now().Minute()
	addTask(s, 2, minute+2)
	assert.Less(t, 1*time.Minute, time.Until(s.s.NextSchedule()))
	assert.Greater(t, 2*time.Minute, time.Until(s.s.NextSchedule()))

	addTask(s, 1, minute+1)
	assert.Less(t, 0*time.Minute, time.Until(s.s.NextSchedule()))
	assert.Greater(t, 1*time.Minute, time.Until(s.s.NextSchedule()))

	s.RemoveTasks([]uint64{1})
	assert.Less(t, 1*time.Minute, time.Until(s.s.NextSchedule()))
	assert.Greater(t, 2*time.Minute, time.Until(s.s.NextSchedule()))

	close(done)
	wg.Wait()
}

func TestScheduleServiceSearchTask(t *testing.T) {
	done := make(chan struct{})
	var wg sync.WaitGroup
	db, err := database.New(":memory:")
	if !assert.NoError(t, err) {
		return
	}

	s := NewService(db)

	setupTest(&wg, s, done)

	now := time.Now()
	minute := now.Minute()
	addTask(s, 1, minute+1)
	addTask(s, 2, minute+2)

	task, nextSched := s.SearchTask(1)
	if assert.NotNil(t, task) {
		assert.Equal(t, task.ID, uint64(1))
		assert.Equal(t, minute+1, nextSched.Minute())
	}

	task, nextSched = s.SearchTask(2)
	if assert.NotNil(t, task) {
		assert.Equal(t, task.ID, uint64(2))
		assert.Equal(t, minute+2, nextSched.Minute())
	}

	task, nextSched = s.SearchTask(3)
	assert.Nil(t, task)
}

func TestScheduleOnceTask(t *testing.T) {
	done := make(chan struct{})
	exec := make(chan struct{})
	var wg sync.WaitGroup
	var isTaskExecuted bool
	db, err := database.New(":memory:")
	if !assert.NoError(t, err) {
		return
	}

	s := NewService(db)

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

	if tx := db.Save(record); !assert.NoError(t, tx.Error) {
		return
	}

	alarmID := record.AlarmID

	sched, err := Parse(now)
	if !assert.NoError(t, err) {
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
	assert.True(t, isTaskExecuted)

	tx := db.First(record, alarmID)
	assert.Error(t, gorm.ErrRecordNotFound, tx.Error)

	close(done)
	wg.Wait()
}
