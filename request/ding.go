package request

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/shengcodehub/auxiliary/common/ding"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"strings"
	"time"
)

type DingRequestLogic struct {
	baseUrl string
}

func NewDingRequestLogic(c ding.Conf) *DingRequestLogic {
	return &DingRequestLogic{
		baseUrl: c.Url,
	}
}

func (u *DingRequestLogic) Post(url string, bodyStr string) ([]byte, error) {
	url = fmt.Sprintf("%s%s", u.baseUrl, url)
	req, err := http.NewRequest("POST", url, strings.NewReader(bodyStr))
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Do(req)

	logx.Infof("【DingRequest】请求url： %v", url)

	if err != nil {
		return nil, err
	}

	return u.parseRes(res)
}

func (u *DingRequestLogic) parseRes(res *http.Response) ([]byte, error) {

	defer func() {
		if res != nil && res.Body != nil {
			_ = res.Body.Close()
		}
	}()

	js, err := io.ReadAll(res.Body)

	logx.Infof("【DingRequest】请求结果 err: %v, %v", err, string(js))

	if err != nil {
		return nil, err
	}

	return js, nil
}

func (u *DingRequestLogic) SendDingMarkMsg(c ding.Conf, val *ding.DingTalkQueue) {
	accessToken := c.AccessToken.FeedBack
	if val.Category == 2 {
		accessToken = c.AccessToken.PayComplaint
	}
	data := &ding.DingMarkData{
		MsgType:  "markdown",
		Markdown: ding.DingMarkdownData{Title: val.Title, Text: val.Content},
		At:       ding.DingAtData{IsAtAll: val.AtAll},
	}
	bodyStr, err := jsoniter.MarshalToString(data)
	if err != nil {
		logx.Errorf("[DingRequestLogic MarshalToString] err: %v", err)
		return
	}
	res, err := u.Post(accessToken, bodyStr)
	if err != nil {
		logx.Errorf("[DingRequestLogic] %s，钉钉返回error：%s", val.Content, err.Error())
		return
	}
	logx.Infof("[DingRequestLogic] %s，钉钉返回：%s", val.Content, res)
}

func (u *DingRequestLogic) SendBackTextMsg(url string, content string) {
	data := &ding.TextData{
		MsgType: "text",
		Text:    ding.TextContent{Content: content},
		At:      ding.DingAtData{IsAtAll: false},
	}
	bodyStr, err := jsoniter.MarshalToString(data)
	if err != nil {
		logx.Errorf("[DingRequestLogic SendBackTextMsg MarshalToString] err: %v", err)
		return
	}
	res, err := u.Post(url, bodyStr)
	if err != nil {
		logx.Errorf("[DingRequestLogic SendBackTextMsg] %s，钉钉返回error：%s", content, err.Error())
		return
	}
	logx.Infof("[DingRequestLogic SendBackTextMsg] %s，钉钉返回：%s", content, res)
}
