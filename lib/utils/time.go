package utils

import (
	"cube/lib/logger"
	"time"
)

var log = logger.Log()

type TimeService struct {
	loc *time.Location
}

func NewTimeService() *TimeService {
	return &TimeService{
		loc: time.Local,
	}
}

func (ts *TimeService) SetTimezone(tz string) {
	location, err := time.LoadLocation(tz)
	if err != nil {
		log.Error(err)
	} else {
		ts.loc = location
	}
}

func (ts TimeService) LocalTime() time.Time {
	return time.Now().In(ts.loc)
}

func (ts TimeService) LocalTimeString(t time.Time) string {
	return t.In(ts.loc).String()
}

func (ts TimeService) Parse(s string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", s, ts.loc)
}
