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

type Action struct {
	ActionType string `json:"type"`
	Name       string `json:"name"`
	Value      string `json:"value"`
	Text       string `json:"text"`
	Style      string `json:"style"`
}

type ActionResponse struct {
	Title      string
	Message    string
	CallbackID string
	Actions    []Action
}

type ActionCallback struct {
	Actions    []Action `json:"actions"`
	CallbackID string   `json:"callback_id"`
	PostID     int      `json:"post_id"`
	Token      string   `json:"token"`
	User       struct {
		UserID   int    `json:"user_id"`
		Username string `json:"username"`
	} `json:"user"`
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

func (r ActionResponse) Normalize() gin.H {
	return gin.H{
		"text": r.Title,
		"attachments": []map[string]interface{}{
			{
				"callback_id": r.CallbackID,
				"text":        r.Message,
				"actions":     r.Actions,
			},
		},
	}
}

func (r ActionResponse) Text() string {
	return r.Message
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
