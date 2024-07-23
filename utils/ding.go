package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func DingCheckSign(ctx *gin.Context) bool {
	sign := ctx.GetHeader("SIGN")
	ts := ctx.GetHeader("TIMESTAMP")
	if len(sign) == 0 || len(ts) == 0 {
		return false
	}
	msg := ts + "\n" + viper.GetString("Ding.AppSecret")
	if DingSign(msg) != sign {
		return false
	}
	return true
}

func DingSign(msg string) string {
	appSecret := viper.GetString("Ding.AppSecret")
	hash := hmac.New(sha256.New, []byte(appSecret))
	hash.Write([]byte(msg))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
