package core

import (
	"testing"

	"cube/core/helpers/alarm"
	"cube/core/helpers/registration"
	"cube/lib/context"
	"cube/lib/database"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type regData struct {
	registration.CreateArgs
	alarms []alarm.CreateArgs
}

func setupTestRegDelete(c *Core) {
	datas := []struct {
		UserID  string
		regArgs []regData
	}{
		{
			UserID: "1",
			regArgs: []regData{
				{
					registration.CreateArgs{"desc1", "tok1", []string{}},
					[]alarm.CreateArgs{
						{0, pattern, msg},
					},
				},
				{
					registration.CreateArgs{"desc1", "tok1", []string{}},
					[]alarm.CreateArgs{
						{0, pattern, msg},
					},
				},
				{
					registration.CreateArgs{"desc1", "tok1", []string{}},
					[]alarm.CreateArgs{
						{0, pattern, msg},
					},
				},
			},
		},
		{
			UserID: "2",
			regArgs: []regData{
				{
					registration.CreateArgs{"desc1", "tok1", []string{}},
					[]alarm.CreateArgs{
						{0, otherAlarm.Pattern, otherAlarm.Message},
					},
				},
			},
		},
	}

	var regID uint64 = 0
	for _, data := range datas {
		ctx := context.ChatContext{
			UserID: data.UserID,
		}

		for _, reg := range data.regArgs {
			c.CreateReg(ctx, reg.Description, reg.Token)
			regID++

			for _, alarm := range reg.alarms {
				c.CreateAlarmToChannel(ctx, regID, alarm.CronPattern, alarm.Message)
			}
		}
	}
}

func TestRegDelete(t *testing.T) {
	c := NewFake()
	ctx := context.ChatContext{
		UserID:   "1",
		UserName: "username1",
	}

	setupTestRegDelete(c)

	resp := c.DeleteReg(ctx, 1, 2)
	assert.Equal(t, context.Response("2 alarms deleted"), resp)

	record := database.Registration{}
	assert.Equal(t, gorm.ErrRecordNotFound, c.DB.First(&record, 1).Error)
}
