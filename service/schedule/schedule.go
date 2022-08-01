package schedule

import (
	"cube/lib/utils"
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

func Parse(s string, ts *utils.TimeService) (sched Schedule, err error) {
	if t, err := ts.Parse(s); err == nil {
		return Schedule{
			onceTime: t,
		}, nil
	}

	if cronSched, err := cron.ParseStandard(s); err == nil {
		return Schedule{
			cron: cronSched,
		}, nil
	}

	return
}
