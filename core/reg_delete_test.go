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
		UserID  int
		regArgs []regData
	}{
		{
			UserID: 1,
			regArgs: []regData{
				{
					registration.CreateArgs{Description: "desc1", Token: "tok1", Dummy: []string{}},
					[]alarm.CreateArgs{
						{Channel: 0, CronPattern: pattern, Message: msg},
					},
				},
				{
					registration.CreateArgs{Description: "desc1", Token: "tok1", Dummy: []string{}},
					[]alarm.CreateArgs{
						{Channel: 0, CronPattern: pattern, Message: msg},
					},
				},
				{
					registration.CreateArgs{Description: "desc1", Token: "tok1", Dummy: []string{}},
					[]alarm.CreateArgs{
						{Channel: 0, CronPattern: pattern, Message: msg},
					},
				},
			},
		},
		{
			UserID: 2,
			regArgs: []regData{
				{
					registration.CreateArgs{Description: "desc1", Token: "tok1", Dummy: []string{}},
					[]alarm.CreateArgs{
						{Channel: 0, CronPattern: otherAlarm.Pattern, Message: otherAlarm.Message},
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
	as := assert.New(t)
	c := NewFake()
	ctx := context.ChatContext{
		UserID:   1,
		Username: "username1",
	}

	setupTestRegDelete(c)

	resp := c.DeleteReg(ctx, 1, 2)
	as.Equal(context.NewTextResponse("2 alarms deleted"), resp)

	record := database.Registration{}
	as.Equal(gorm.ErrRecordNotFound, c.DB.First(&record, 1).Error)
}
