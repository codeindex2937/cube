package core

import (
	"testing"

	"cube/lib/context"

	"github.com/stretchr/testify/assert"
)

func TestRegList(t *testing.T) {
	c := NewFake()
	description1 := "reg1"
	token1 := "tok1"
	description2 := "reg2"
	token2 := "tok2"
	ctx1 := context.ChatContext{
		UserID:   "1",
		UserName: "username1",
	}
	ctx2 := context.ChatContext{
		UserID:   "2",
		UserName: "username2",
	}

	resp := c.CreateReg(ctx1, description1, token1)
	if !assert.Equal(t, context.Response(`ID=1 "reg1" Token=tok1`), resp) {
		return
	}

	resp = c.CreateReg(ctx2, description2, token2)
	if !assert.Equal(t, context.Response(`ID=2 "reg2" Token=tok2`), resp) {
		return
	}

	resp = c.ListReg(ctx1)
	if !assert.Equal(t, context.Response("ID=1 \"reg1\" Token=tok1\n"), resp) {
		return
	}

	resp = c.ListReg(ctx2)
	if !assert.Equal(t, context.Response("ID=2 \"reg2\" Token=tok2\n"), resp) {
		return
	}
}
