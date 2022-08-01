package core

import (
	"strconv"
	"testing"

	"cube/lib/context"
	"cube/lib/database"
	"cube/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestRegClear(t *testing.T) {
	as := assert.New(t)
	s := assert.New(t)
	c := NewFake()
	uid, _ := strconv.Atoi(userID)
	ctx := context.ChatContext{
		UserID: uid,
	}

	setupTestRegDelete(c)

	sched := c.Schedule.(*fake.ScheduleService)
	s.True(sched.ExistTask(1))
	s.True(sched.ExistTask(2))
	s.True(sched.ExistTask(3))
	s.True(sched.ExistTask(4))

	resp := c.ClearReg(ctx)
	as.Equal(context.NewTextResponse("3 registrations deleted, 3 alarms deleted"), resp)

	s.False(sched.ExistTask(1))
	s.False(sched.ExistTask(2))
	s.False(sched.ExistTask(3))
	s.True(sched.ExistTask(4))

	records := []database.Alarm{}
	if !s.NoError(c.DB.Find(&records).Error) {
		return
	}

	if s.Equal(1, len(records)) {
		s.Equal("2", records[0].UserID)
		s.Equal(uint64(4), records[0].RegID)
		s.Equal(otherAlarm.Pattern, records[0].Pattern)
		s.Equal(otherAlarm.Message, records[0].Message)
	}
}
