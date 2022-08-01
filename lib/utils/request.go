package utils

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"

	"cube/lib/context"

	"github.com/alexflint/go-arg"
)

func ParseRequest(body []byte) (ctx context.ChatContext, err error) {
	params, err := url.ParseQuery(string(body))
	if err != nil {
		return ctx, fmt.Errorf("ParseQuery: %w", err)
	}

	userID, _ := strconv.Atoi(params.Get("user_id"))
	return context.ChatContext{
		Token:       params.Get("token"),
		ChannelId:   params.Get("channel_id"),
		ChannelName: params.Get("channel_name"),
		UserID:      userID,
		Username:    params.Get("user_name"),
		Text:        params.Get("text"),
	}, nil
}

func PrintHelp(cmd string, args interface{}) context.IResponse {
	p, _ := arg.NewParser(arg.Config{Program: cmd}, args)
	buf := new(bytes.Buffer)
	buf.WriteString("```\n")
	p.WriteHelp(buf)
	buf.WriteString("```\n")
	return context.NewTextResponse(buf.String())
}

func PrintUsage(config arg.Config, args interface{}) context.IResponse {
	p, _ := arg.NewParser(config, args)
	buf := new(bytes.Buffer)
	buf.WriteString("```\n")
	p.WriteUsage(buf)
	buf.WriteString("```\n")
	return context.NewTextResponse(buf.String())
}
