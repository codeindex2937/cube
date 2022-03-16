package core

import (
	"testing"

	"cube/lib/context"
	"cube/lib/database"
	"cube/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestRegClear(t *testing.T) {
	c := NewFake()
	ctx := context.ChatContext{
		UserID: userID,
	}

	setupTestRegDelete(c)

	sched := c.Schedule.(*fake.ScheduleService)
	assert.True(t, sched.ExistTask(1))
	assert.True(t, sched.ExistTask(2))
	assert.True(t, sched.ExistTask(3))
	assert.True(t, sched.ExistTask(4))

	resp := c.ClearReg(ctx)
	assert.Equal(t, context.Response("3 registrations deleted, 3 alarms deleted"), resp)

	assert.False(t, sched.ExistTask(1))
	assert.False(t, sched.ExistTask(2))
	assert.False(t, sched.ExistTask(3))
	assert.True(t, sched.ExistTask(4))

	records := []database.Alarm{}
	tx := c.DB.Find(&records)
	if !assert.NoError(t, tx.Error) {
		return
	}

	if assert.Equal(t, 1, len(records)) {
		assert.Equal(t, "2", records[0].UserID)
		assert.Equal(t, uint64(4), records[0].RegID)
		assert.Equal(t, otherAlarm.Pattern, records[0].Pattern)
		assert.Equal(t, otherAlarm.Message, records[0].Message)
	}
}
