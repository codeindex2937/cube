package context

type Response string

const Success Response = "success"

type ChatContext struct {
	Token       string `json:"token"`
	UserID      string `json:"user_id"`
	UserName    string `json:"username"`
	Text        string `json:"text"`
	ChannelId   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	PostId      string `json:"post_id"`
	Timestamp   int64  `json:"timestamp"`
	TriggerWord string `json:"trigger_wordd"`
}

type Context struct {
	Args []string
	Req  ChatContext
}
