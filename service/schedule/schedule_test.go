package schedule

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseOnce(t *testing.T) {
	sched, err := Parse("2021-12-07 01:02:03", ts)
	assert.NoError(t, err)
	assert.False(t, sched.onceTime.IsZero())
	assert.Equal(t, 2021, sched.onceTime.Year())
	assert.Equal(t, time.Month(12), sched.onceTime.Month())
	assert.Equal(t, 7, sched.onceTime.Day())
	assert.Equal(t, 1, sched.onceTime.Hour())
	assert.Equal(t, 2, sched.onceTime.Minute())
	assert.Equal(t, 3, sched.onceTime.Second())
	assert.Equal(t, time.Duration(0), sched.Next(time.Now()).Sub(sched.onceTime))
	assert.True(t, sched.IsOnce())
}

func TestParseCron(t *testing.T) {
	now := time.Now()
	sched, err := Parse("* * * * *", ts)
	assert.NoError(t, err)
	assert.True(t, sched.onceTime.IsZero())
	assert.Greater(t, time.Duration(0), now.Sub(sched.cron.Next(now)))
	assert.Less(t, time.Duration(0), now.Add(time.Minute).Sub(sched.cron.Next(now)))
	assert.False(t, sched.IsOnce())
}
