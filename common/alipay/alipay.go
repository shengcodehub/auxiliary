package alipay

import (
	"context"
	"github.com/shengcodehub/auxiliary/utils"
	"github.com/smartwalle/alipay/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Conf struct {
	Default AuthConf
}

type AuthConf struct {
	AppId          string
	RsaPrivateKey  string
	RsaPublicKey   string
	AppCertPath    string
	AlipayCertPath string
	RootCertPath   string
	NotifyURL      string
	ReturnURL      string
}

var alipayClient *alipay.Client

func Setup(c Conf) {
	privateKey, err := os.ReadFile(c.Default.RsaPrivateKey)
	if err != nil {
		logx.Errorf("read alipay RsaPrivateKey failed %s", err.Error())
		return
	}
	client, err := alipay.New(c.Default.AppId, string(privateKey), true)
	if err != nil {
		logx.Errorf("open alipay failed %s", err.Error())
		return
	}

	// 加载应用公钥证书
	if err = client.LoadAppCertPublicKeyFromFile(c.Default.AppCertPath); err != nil {
		logx.Errorf("open alipay AppCertPath failed %s", err.Error())
		return
	}

	// 加载支付宝根证书
	if err = client.LoadAliPayRootCertFromFile(c.Default.RootCertPath); err != nil {
		logx.Errorf("open alipay AlipayCertPath failed %s", err.Error())
		return
	}

	// 加载支付宝公钥证书
	if err = client.LoadAlipayCertPublicKeyFromFile(c.Default.AlipayCertPath); err != nil {
		logx.Errorf("open alipay AlipayCertPath failed %s", err.Error())
		return
	}

	//// 加载内容密钥，可选
	//if err = client.SetEncryptKey("FtVd5SgrsUzYQRAPBmejHQ=="); err != nil {
	//	// 错误处理
	//}

	alipayClient = client
}

func GetAlipay() *alipay.Client {
	return alipayClient
}

func TradeWapPay(c Conf, orderID string, subject string, totalFee float64) (string, error) {
	var p = alipay.TradeWapPay{}
	p.NotifyURL = c.Default.NotifyURL
	// p.ReturnURL = c.Default.ReturnURL
	p.Subject = subject
	p.OutTradeNo = orderID
	p.TotalAmount = strconv.FormatFloat(totalFee, 'f', -1, 64)
	p.ProductCode = "QUICK_WAP_WAY"
	url, err := GetAlipay().TradeWapPay(p)
	if err != nil {
		logx.Errorf("alipay TradeWapPay failed %s", err.Error())
		return "", err
	}
	// 这个 payURL 即是用于打开支付宝支付页面的 URL，可将输出的内容复制，到浏览器中访问该 URL 即可打开支付页面。
	payURL := url.String()
	return payURL, nil
}

func RiskComplaintQuery(ctx context.Context) (result *alipay.SecurityRiskComplaintInfoBatchQueryRsp, err error) {
	p := alipay.SecurityRiskComplaintInfoBatchQueryReq{
		CurrentPageNum:    1,
		PageSize:          100,
		GmtComplaintStart: utils.GetTimestampByTime(time.Now().Unix() - 300),
		GmtComplaintEnd:   utils.GetTimestamp(time.Now()),
		StatusList:        []string{"WAIT_PROCESS", "PROCESSING"},
	}
	logx.Infof("alipay RiskComplaintQuery req:%+v", p)
	result, err = GetAlipay().SecurityRiskComplaintInfoBatchQuery(ctx, p)
	if err != nil {
		logx.Errorf("alipay RiskComplaintQuery failed %s", err.Error())
		return nil, err
	}
	return result, nil
}

func OrderQuery(ctx context.Context, orderID string) (result *alipay.TradeQueryRsp, err error) {
	p := alipay.TradeQuery{}
	p.OutTradeNo = orderID
	result, err = GetAlipay().TradeQuery(ctx, p)
	if err != nil {
		logx.Errorf("alipay OrderQuery failed %s, orderId:%s", err.Error(), orderID)
		return nil, err
	}
	return result, nil
}

func TradeClose(ctx context.Context, orderID string) (result *alipay.TradeCloseRsp, err error) {
	p := alipay.TradeClose{}
	p.OutTradeNo = orderID
	result, err = GetAlipay().TradeClose(ctx, p)
	if err != nil {
		logx.Errorf("alipay OrderClosed failed %s, orderId:%s", err.Error(), orderID)
		return nil, err
	}
	return result, nil
}

func VerifyNotify(w http.ResponseWriter, r *http.Request, callback func(data *alipay.Notification)) {
	err := r.ParseForm()
	if err != nil {
		logx.Errorf("alipay VerifyNotify ParseForm failed %s", err.Error())
		return
	}
	data, err := GetAlipay().DecodeNotification(r.Form)
	if err != nil {
		logx.Errorf("alipay VerifyNotify DecodeNotification failed %s", err.Error())
		return
	}
	callback(data)
	alipay.ACKNotification(w)
}
