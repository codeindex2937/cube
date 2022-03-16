package core

import (
	"testing"

	"cube/lib/context"
	"cube/lib/database"
	"cube/lib/fake"

	"github.com/stretchr/testify/assert"
)

var userID = "1"
var pattern = "* * * * *"
var msg = "surprise"
var otherAlarm = database.Alarm{
	UserID:  "2",
	Pattern: "* * * * *",
	Message: "what ever",
}

func TestAlarmClear(t *testing.T) {
	c := NewFake()
	ctx := context.ChatContext{
		UserID: userID,
	}

	setupTestAlarmDelete(c)

	sched := c.Schedule.(*fake.ScheduleService)
	assert.True(t, sched.ExistTask(1))
	assert.True(t, sched.ExistTask(2))
	assert.True(t, sched.ExistTask(3))

	resp := c.ClearAlarm(ctx)
	assert.Equal(t, context.Response("3 alarms deleted"), resp)

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
		assert.Equal(t, otherAlarm.Pattern, records[0].Pattern)
		assert.Equal(t, otherAlarm.UserID, records[0].UserID)
		assert.Equal(t, otherAlarm.Message, records[0].Message)
	}
}
