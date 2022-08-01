package core

import (
	"strconv"
	"testing"

	"cube/lib/context"
	"cube/lib/database"
	"cube/lib/fake"

	"github.com/stretchr/testify/assert"
)

func setupTestAlarmDelete(c *Core) {
	datas := []struct {
		UserID    string
		alarmArgs [][]string
	}{
		{
			UserID: "1",
			alarmArgs: [][]string{
				{pattern, msg},
				{pattern, msg},
				{pattern, msg},
			},
		},
		{
			UserID: "2",
			alarmArgs: [][]string{
				{otherAlarm.Pattern, otherAlarm.Message},
			},
		},
	}

	for _, data := range datas {
		uid, _ := strconv.Atoi(data.UserID)
		ctx := context.ChatContext{
			UserID: uid,
		}

		for _, args := range data.alarmArgs {
			c.CreateAlarm(ctx, args[0], args[1])
		}
	}
}

func TestAlarmDelete(t *testing.T) {
	c := NewFake()
	uid, _ := strconv.Atoi(userID)
	ctx := context.ChatContext{
		UserID: uid,
	}

	setupTestAlarmDelete(c)

	sched := c.Schedule.(*fake.ScheduleService)
	assert.True(t, sched.ExistTask(1))
	assert.True(t, sched.ExistTask(2))
	assert.True(t, sched.ExistTask(3))
	assert.True(t, sched.ExistTask(4))

	resp := c.DeleteAlarm(ctx, "2", "3")
	assert.Equal(t, context.NewTextResponse("2 alarms deleted"), resp)

	assert.True(t, sched.ExistTask(1))
	assert.False(t, sched.ExistTask(2))
	assert.False(t, sched.ExistTask(3))
	assert.True(t, sched.ExistTask(4))

	records := []database.Alarm{}
	tx := c.DB.Find(&records)
	if !assert.NoError(t, tx.Error) {
		return
	}

	if assert.Equal(t, 2, len(records)) {
		assert.Equal(t, pattern, records[0].Pattern)
		assert.Equal(t, userID, records[0].UserID)
		assert.Equal(t, msg, records[0].Message)
		//
		assert.Equal(t, otherAlarm.Pattern, records[1].Pattern)
		assert.Equal(t, otherAlarm.UserID, records[1].UserID)
		assert.Equal(t, otherAlarm.Message, records[1].Message)
	}
}
