package yidun

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/galaxy-book/captcha-golang-demo/sdk"
	jsoniter "github.com/json-iterator/go"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Conf struct {
	CaptchaId string
	SecretId  string
	SecretKey string
}
type TextConf struct {
	ImgBusinessId      string
	HallTextBusinessId string
	TextBusinessId     string
	SecretId           string
	SecretKey          string
}

func Verify(c Conf, validate string, account string) (bool, error) {
	verifier, err := sdk.New(c.CaptchaId, c.SecretId, c.SecretKey)
	if err != nil {
		return false, err
	}
	verifyResult, err := verifier.Verify(validate, account)
	if err != nil {
		return false, err
	}
	if verifyResult.Result {
		return true, nil
	} else {
		return false, nil
	}
}

func TextCheck(c TextConf, content string, uid int64) bool {
	params := url.Values{
		"content": []string{content},
	}
	ret := check(params, uid, c.SecretId, c.TextBusinessId, c.SecretKey)

	code, _ := ret.Get("code").Int()
	message, _ := ret.Get("msg").String()
	if code == 200 {
		result := ret.Get("result")
		antispam := result.Get("antispam")
		if antispam != nil {
			taskId, _ := antispam.Get("taskId").String()
			//dataId, _ := antispam.Get("dataId").String()
			suggestion, _ := antispam.Get("suggestion").Int()
			//suggestionLevel, _ := antispam.Get("suggestionLevel").Int()
			//resultType, _ := antispam.Get("resultType").Int()
			//censorType, _ := antispam.Get("censorType").Int()
			//strategyVersions, _ := antispam.Get("strategyVersions").Array()
			//isRelatedHit, _ := antispam.Get("isRelatedHit").Bool()
			labels, _ := antispam.Get("labels").Array()
			if suggestion == 0 {
				logx.Infof("yidun taskId: %s, 文本机器检测结果: 通过", taskId)
				return true
			} else if suggestion == 1 {
				logx.Infof("yidun taskId: %s, 文本机器检测结果: 嫌疑, 需人工复审, 分类信息如下: %s", taskId, labels)
			} else if suggestion == 2 {
				logx.Infof("yidun taskId=%s, 文本机器检测结果: 不通过, 分类信息如下: %s", taskId, labels)
			}
		}
	} else {
		logx.Errorf("yidun 文本机器检测结果 ERROR: code=%d, msg=%s", code, message)
	}
	return false
}

func check(params url.Values, uid int64, secretId string, textBusinessId string, secretKey string) *simplejson.Json {
	params["dataId"] = []string{generateUniqueIdentifier(uid)}
	params["secretId"] = []string{secretId}
	params["businessId"] = []string{textBusinessId}
	params["version"] = []string{"v5.3"}
	params["timestamp"] = []string{strconv.FormatInt(time.Now().UnixNano()/1000000, 10)}
	params["nonce"] = []string{strconv.FormatInt(rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(10000000000), 10)}
	params["signature"] = []string{genSignature(params, secretKey)}
	apiUrl := "http://as.dun.163.com/v5/text/check"
	str, _ := jsoniter.MarshalToString(params)
	resp, err := http.Post(apiUrl, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	logx.Infof("yidun 文本检测请求：%s", str)
	if err != nil {
		logx.Errorf("易盾文本验证调用API接口失败:%s, data:%s", err.Error(), str)
		return nil
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	contents, _ := io.ReadAll(resp.Body)
	result, _ := simplejson.NewJson(contents)
	return result
}

func generateUniqueIdentifier(uid int64) string {
	// 使用当前时间戳和随机数生成唯一标识符
	rand.NewSource(time.Now().UnixNano())
	timestamp := time.Now().Unix()
	randomNumber := rand.Intn(100000)
	uniqueID := fmt.Sprintf("%d%d%d", uid, timestamp, randomNumber)
	return uniqueID
}

func genSignature(params url.Values, secretKey string) string {
	var paramStr string
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		paramStr += key + params[key][0]
	}
	paramStr += secretKey
	md5Reader := md5.New()
	md5Reader.Write([]byte(paramStr))
	return hex.EncodeToString(md5Reader.Sum(nil))
}
