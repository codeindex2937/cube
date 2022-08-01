package core

import (
	"testing"

	"cube/lib/context"
	"cube/lib/database"
	"cube/lib/fake"

	"github.com/stretchr/testify/assert"
)

func TestAlarmCreate(t *testing.T) {
	as := assert.New(t)
	c := NewFake()
	pattern := "* * * * *"
	message := "surprise"
	ctx := context.ChatContext{
		UserID: 1,
	}

	resp := c.CreateAlarm(ctx, pattern, message)
	if !as.Equal(context.NewTextResponse("ID=1 \"* * * * *\" \"surprise\" Next=0001-01-01 00:00:00 +0000 UTC"), resp) {
		return
	}

	sched := c.Schedule.(*fake.ScheduleService)
	for _, i := range []uint64{1} {
		as.True(sched.ExistTask(i))
	}

	var record database.Alarm
	tx := c.DB.First(&record)
	if !as.NoError(tx.Error) {
		return
	}

	as.Equal(pattern, record.Pattern)
	as.Equal(userID, record.UserID)
	as.Equal(message, record.Message)
}
