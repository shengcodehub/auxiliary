package ding

import (
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/utils"
)

type Conf struct {
	Url         string
	AccessToken TokenConf
}

type TokenConf struct {
	FeedBack     string
	PayComplaint string
}

type DingMarkData struct {
	MsgType  string           `json:"msgtype"`
	Markdown DingMarkdownData `json:"markdown"`
	At       DingAtData       `json:"at"`
}

type DingMarkdownData struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type DingAtData struct {
	IsAtAll bool `json:"isAtAll"`
}

type RobotConf struct {
	ClientId     string
	ClientSecret string
	Topic        string
}

type TextData struct {
	MsgType string      `json:"msgtype"`
	Text    TextContent `json:"text"`
	At      DingAtData  `json:"at"`
}

type TextContent struct {
	Content string `json:"content"`
}

type DingTalkQueue struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	AtAll    bool   `json:"at_all"`
	Category int64  `json:"category"` // 1 问题反馈 2 支付投诉
}

func StartRobot(c RobotConf, onChatReceive chatbot.IChatBotMessageHandler) *client.StreamClient {
	logger.SetLogger(logger.NewStdTestLogger())
	cli := client.NewStreamClient(
		client.WithAppCredential(client.NewAppCredentialConfig(c.ClientId, c.ClientSecret)),
		client.WithUserAgent(client.NewDingtalkGoSDKUserAgent()),
		client.WithSubscription(utils.SubscriptionTypeKCallback, c.Topic, chatbot.NewDefaultChatBotFrameHandler(onChatReceive).OnEventReceived),
	)
	return cli
}
