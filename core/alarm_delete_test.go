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
	as.True(sched.ExistTask(4))

	resp := c.DeleteAlarm(ctx, "2", "3")
	as.Equal(context.NewTextResponse("2 alarms deleted"), resp)

	as.True(sched.ExistTask(1))
	as.False(sched.ExistTask(2))
	as.False(sched.ExistTask(3))
	as.True(sched.ExistTask(4))

	records := []database.Alarm{}
	tx := c.DB.Find(&records)
	if !as.NoError(tx.Error) {
		return
	}

	if as.Equal(2, len(records)) {
		as.Equal(pattern, records[0].Pattern)
		as.Equal(userID, records[0].UserID)
		as.Equal(msg, records[0].Message)
		//
		as.Equal(otherAlarm.Pattern, records[1].Pattern)
		as.Equal(otherAlarm.UserID, records[1].UserID)
		as.Equal(otherAlarm.Message, records[1].Message)
	}
}
