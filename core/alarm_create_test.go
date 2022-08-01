package core

import (
	"testing"

	"cube/lib/context"
	"cube/lib/database"
	"cube/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestAlarmCreate(t *testing.T) {
	c := NewFake()
	pattern := "* * * * *"
	message := "surprise"
	ctx := context.ChatContext{
		UserID: "1",
	}

	resp := c.CreateAlarm(ctx, pattern, message)
	if !assert.Equal(t, context.NewTextResponse("ID=1 \"* * * * *\" \"surprise\" Next=0001-01-01 00:00:00 +0000 UTC"), resp) {
		return
	}

	sched := c.Schedule.(*fake.ScheduleService)
	for _, i := range []uint64{1} {
		assert.True(t, sched.ExistTask(i))
	}

	var record database.Alarm
	tx := c.DB.First(&record)
	if !assert.NoError(t, tx.Error) {
		return
	}

	assert.Equal(t, pattern, record.Pattern)
	assert.Equal(t, userID, record.UserID)
	assert.Equal(t, message, record.Message)
}
