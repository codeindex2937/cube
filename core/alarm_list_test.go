package core

import (
	"testing"

	"cube/lib/context"

	"github.com/stretchr/testify/assert"
)

func TestAlarmList(t *testing.T) {
	c := NewFake()
	pattern := "* * * * *"
	message := "surprise"
	ctx1 := context.ChatContext{
		UserID:   1,
		Username: "username1",
	}
	ctx2 := context.ChatContext{
		UserID:   2,
		Username: "username2",
	}

	resp := c.CreateAlarm(ctx1, pattern, message)
	if !assert.Equal(t, context.NewTextResponse("ID=1 \"* * * * *\" \"surprise\" Next=0001-01-01 00:00:00 +0000 UTC"), resp) {
		return
	}

	resp = c.CreateAlarm(ctx2, "* * * * *", "whatever")
	if !assert.Equal(t, context.NewTextResponse("ID=2 \"* * * * *\" \"whatever\" Next=0001-01-01 00:00:00 +0000 UTC"), resp) {
		return
	}

	resp = c.ListAlarm(ctx1)
	assert.Equal(t, context.NewTextResponse("ID=1 \"* * * * *\" \"surprise\" Next=0001-01-01 00:00:00 +0000 UTC\n"), resp)
}
