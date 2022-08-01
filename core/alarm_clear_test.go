package core

import (
	"strconv"
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
	as := assert.New(t)
	c := NewFake()
	uid, _ := strconv.Atoi(userID)
	ctx := context.ChatContext{
		UserID: uid,
	}

	setupTestAlarmDelete(c)

	sched := c.Schedule.(*fake.ScheduleService)
	as.True(sched.ExistTask(1))
	as.True(sched.ExistTask(2))
	as.True(sched.ExistTask(3))

	resp := c.ClearAlarm(ctx)
	as.Equal("3 alarms deleted", resp.Text())

	as.False(sched.ExistTask(1))
	as.False(sched.ExistTask(2))
	as.False(sched.ExistTask(3))
	as.True(sched.ExistTask(4))

	records := []database.Alarm{}
	tx := c.DB.Find(&records)
	if !as.NoError(tx.Error) {
		return
	}

	if as.Equal(1, len(records)) {
		as.Equal(otherAlarm.Pattern, records[0].Pattern)
		as.Equal(otherAlarm.UserID, records[0].UserID)
		as.Equal(otherAlarm.Message, records[0].Message)
	}
}
