package schedule

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseOnce(t *testing.T) {
	as := assert.New(t)
	sched, err := Parse("2021-12-07 01:02:03", ts)
	as.NoError(err)
	as.False(sched.onceTime.IsZero())
	as.Equal(2021, sched.onceTime.Year())
	as.Equal(time.Month(12), sched.onceTime.Month())
	as.Equal(7, sched.onceTime.Day())
	as.Equal(1, sched.onceTime.Hour())
	as.Equal(2, sched.onceTime.Minute())
	as.Equal(3, sched.onceTime.Second())
	as.Equal(time.Duration(0), sched.Next(time.Now()).Sub(sched.onceTime))
	as.True(sched.IsOnce())
}

func TestParseCron(t *testing.T) {
	as := assert.New(t)
	now := time.Now()
	sched, err := Parse("* * * * *", ts)
	as.NoError(err)
	as.True(sched.onceTime.IsZero())
	as.Greater(time.Duration(0), now.Sub(sched.cron.Next(now)))
	as.Less(time.Duration(0), now.Add(time.Minute).Sub(sched.cron.Next(now)))
	as.False(sched.IsOnce())
}
