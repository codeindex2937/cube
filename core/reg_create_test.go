package core

import (
	"testing"

	"cube/lib/context"
	"cube/lib/database"

	"github.com/stretchr/testify/assert"
)

func TestRegCreate(t *testing.T) {
	as := assert.New(t)
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
	as.Equal(context.NewTextResponse(`ID=1 "desc1" Token=tok1`), resp)

	record := database.Registration{}
	tx := c.DB.First(&record, map[string]interface{}{"token": token1})
	if as.NoError(tx.Error) {
		as.Equal(token1, record.Token)
		as.Equal(ctx1.Username, record.UserName)
	}

	resp = c.CreateReg(ctx2, description2, token2)
	as.Equal(context.NewTextResponse(`ID=2 "desc2" Token=tok2`), resp)

	record = database.Registration{}
	tx = c.DB.First(&record, map[string]interface{}{"token": token2})
	if as.NoError(tx.Error) {
		as.Equal(token2, record.Token)
		as.Equal(ctx2.Username, record.UserName)
	}
}
