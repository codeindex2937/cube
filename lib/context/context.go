package context

import (
	"github.com/gin-gonic/gin"
)

type IResponse interface {
	Normalize() gin.H
	Text() string
}

type TextResponse struct {
	text string
}

var Success = TextResponse{
	text: "success",
}

type ChatContext struct {
	Token       string `json:"token"`
	UserID      int    `json:"user_id"`
	Username    string `json:"username"`
	Text        string `json:"text"`
	ChannelId   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	PostID      string `json:"post_id"`
	Timestamp   int64  `json:"timestamp"`
	TriggerWord string `json:"trigger_word"`
}

type Context struct {
	Args []string
	Req  ChatContext
}

func (r TextResponse) Normalize() gin.H {
	return gin.H{
		"text": r.text,
	}
}

func (r TextResponse) Text() string {
	return r.text
}

func NewTextResponse(text string) IResponse {
	return TextResponse{
		text: text,
	}
}

func NewErrorResponse(err error) IResponse {
	return TextResponse{
		text: err.Error(),
	}
}
