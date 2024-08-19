package main

import (
	"github.com/centrifugal/centrifuge-go"
	"github.com/zeromicro/go-zero/core/logx"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-jwt/jwt/v4"
)

func GetAccSub(uuid string) {
	channel := "acc:" + uuid
	ctx := gctx.New()

	token := connToken(uuid, 0)
	client := centrifuge.NewJsonClient(
		"ws://localhost:8000/connection/websocket",
		centrifuge.Config{
			Token:              token,
			ReadTimeout:        1 * time.Minute,
			WriteTimeout:       1 * time.Minute,
			HandshakeTimeout:   1 * time.Minute,
			MaxServerPingDelay: 3 * time.Minute,
		},
	)
	defer client.Close()

	client.OnConnecting(func(e centrifuge.ConnectingEvent) {
		logx.Infof("client.OnConnecting code:%d, reason:%s", e.Code, e.Reason)
	})
	client.OnConnected(func(e centrifuge.ConnectedEvent) {
		logx.Infof("client.OnConnected ID:%s", e.ClientID)
	})
	client.OnDisconnected(func(e centrifuge.DisconnectedEvent) {
		logx.Errorf("client.OnDisconnected code:%d, reason:%s", e.Code, e.Reason)
	})
	client.OnError(func(e centrifuge.ErrorEvent) {
		logx.Errorf("client.OnError fail:%s", e.Error)
	})

	err := client.Connect()
	if err != nil {
		logx.Errorf("client.Connect 失败 fail:%s uuid:%s", err.Error(), uuid)
		return
	}

	sub, err := client.NewSubscription(channel)
	if err != nil {
		logx.Errorf("client.NewSubscription 失败 fail:%s uuid:%s", err.Error(), uuid)
		return
	}
	sub.OnSubscribing(func(e centrifuge.SubscribingEvent) {
		logx.Infof("sub.OnSubscribing Channel:%s, code:%d, reason:%s", sub.Channel, e.Code, e.Reason)
	})
	sub.OnSubscribed(func(e centrifuge.SubscribedEvent) {
		logx.Infof("sub.OnSubscribed Channel:%s", sub.Channel)
	})
	sub.OnUnsubscribed(func(e centrifuge.UnsubscribedEvent) {
		logx.Infof("sub.OnUnsubscribed Channel:%s, code:%d, reason:%s", sub.Channel, e.Code, e.Reason)
	})

	sub.OnPublication(func(e centrifuge.PublicationEvent) {
		logx.Infof("sub.OnPublication 消息 uuid:%s, channel:%s, chatMessage:%s", uuid, sub.Channel, string(e.Data))
	})

	sub.OnError(func(e centrifuge.SubscriptionErrorEvent) {
		logx.Errorf("sub.OnError 消息 fail:%s", e.Error)
	})
	err = sub.Subscribe()
	if err != nil {
		logx.Errorf("sub.Subscribe 失败 fail:%s, uuid:%s", err.Error(), uuid)
		return
	}

	//推送开始匹配
	msg := AccRequest{
		Seq:  "1692613366219-776882",
		Cmd:  "heartbeat", // 心跳
		Data: g.Map{},
	}

	data := gjson.New(msg).String()
	_, err = sub.Publish(ctx, gconv.Bytes(data))
	if err != nil {
		logx.Errorf("发送失败 err:%v, uuid:%s", err.Error(), uuid)
	} else {
		logx.Infof("发送成功 uuid:%s", uuid)
	}

	time.Sleep(10 * time.Minute)
}

func connToken(uuid string, exp int64) string {
	// NOTE that JWT must be generated on backend side of your application!
	// Here we are generating it on client side only for example simplicity.
	claims := jwt.MapClaims{"sub": uuid}
	if exp > 0 {
		claims["exp"] = exp
	}
	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("6a5ade9c-a1c9-4131-a9a3-56c398cde55a"))
	if err != nil {
		panic(err)
	}
	return t
}

type AccRequest struct {
	Seq  string      `json:"seq"`  // 消息的唯一Id
	Cmd  string      `json:"cmd"`  // 请求命令字
	Data interface{} `json:"data"` // 数据 json
}
