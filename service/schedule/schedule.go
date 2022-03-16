package schedule

import (
	"time"

	"github.com/robfig/cron"
)

type Schedule struct {
	onceTime time.Time
	cron     cron.Schedule
}

func (sched Schedule) Next(offset time.Time) time.Time {
	if !sched.onceTime.IsZero() {
		return sched.onceTime
	}

	return sched.cron.Next(offset)
}

func (sched Schedule) IsOnce() bool {
	return !sched.onceTime.IsZero()
}

func Parse(s string) (sched Schedule, err error) {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	if err == nil {
		return Schedule{
			onceTime: t,
		}, nil
	}

	cronSched, err := cron.ParseStandard(s)
	if err == nil {
		return Schedule{
			cron: cronSched,
		}, nil
	}

	return
}
