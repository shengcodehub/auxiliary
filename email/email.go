package email

import (
	"bytes"
	"github.com/zeromicro/go-zero/core/logx"
	"gopkg.in/gomail.v2"
	"io"
)

type Conf struct {
	Host     string
	Port     int
	Username string
	Password string
}

func Send(c Conf, to []string, subject string, body string, zipBuf *bytes.Buffer, fileName string, from string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	if zipBuf != nil {
		// 添加附件
		m.Attach(fileName, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(zipBuf.Bytes())
			return err
		}))
	}
	// 设置SMTP服务器信息
	d := gomail.NewDialer(c.Host, c.Port, c.Username, c.Password)
	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		logx.Errorf("Failed to send email: %v", err)
		return err
	}
	return nil
}
