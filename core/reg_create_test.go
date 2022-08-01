package core

import (
	"testing"

	"cube/lib/context"
	"cube/lib/database"

	"github.com/stretchr/testify/assert"
)

func TestRegCreate(t *testing.T) {
	c := NewFake()
	description1 := "desc1"
	token1 := "tok1"
	ctx1 := context.ChatContext{
		UserID:   1,
		Username: "username1",
	}
	description2 := "desc2"
	token2 := "tok2"
	ctx2 := context.ChatContext{
		UserID:   2,
		Username: "username2",
	}

	resp := c.CreateReg(ctx1, description1, token1)
	assert.Equal(t, context.NewTextResponse(`ID=1 "desc1" Token=tok1`), resp)

	record := database.Registration{}
	tx := c.DB.First(&record, map[string]interface{}{"token": token1})
	if assert.NoError(t, tx.Error) {
		assert.Equal(t, token1, record.Token)
		assert.Equal(t, ctx1.Username, record.UserName)
	}

	resp = c.CreateReg(ctx2, description2, token2)
	assert.Equal(t, context.NewTextResponse(`ID=2 "desc2" Token=tok2`), resp)

	record = database.Registration{}
	tx = c.DB.First(&record, map[string]interface{}{"token": token2})
	if assert.NoError(t, tx.Error) {
		assert.Equal(t, token2, record.Token)
		assert.Equal(t, ctx2.Username, record.UserName)
	}
}
