package utils

import (
	"bytes"
	"fmt"
	"net/url"

	"cube/lib/context"

	"github.com/alexflint/go-arg"
)

func ParseRequest(body []byte) (ctx context.ChatContext, err error) {
	params, err := url.ParseQuery(string(body))
	if err != nil {
		return ctx, fmt.Errorf("ParseQuery: %w", err)
	}

	return context.ChatContext{
		Token:       params.Get("token"),
		ChannelId:   params.Get("channel_id"),
		ChannelName: params.Get("channel_name"),
		UserID:      params.Get("user_id"),
		UserName:    params.Get("user_name"),
		Text:        params.Get("text"),
	}, nil
}

func PrintHelp(cmd string, args interface{}) context.Response {
	p, _ := arg.NewParser(arg.Config{Program: cmd}, args)
	buf := new(bytes.Buffer)
	buf.WriteString("```\n")
	p.WriteHelp(buf)
	buf.WriteString("```\n")
	return context.Response(buf.String())
}

func PrintUsage(config arg.Config, args interface{}) context.Response {
	p, _ := arg.NewParser(config, args)
	buf := new(bytes.Buffer)
	buf.WriteString("```\n")
	p.WriteUsage(buf)
	buf.WriteString("```\n")
	return context.Response(buf.String())
}
